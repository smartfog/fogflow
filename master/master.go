package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	. "github.com/smartfog/fogflow/common/communicator"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"

	. "github.com/smartfog/fogflow/common/config"
)

type Master struct {
	cfg *Config

	BrokerURL string

	id           string
	myURL        string
	messageBus   string
	discoveryURL string

	communicator *Communicator
	ticker       *time.Ticker
	agent        *NGSIAgent

	//list of all workers
	workers         map[string]*WorkerProfile
	workerList_lock sync.RWMutex

	//list of all operators
	operatorList      map[string]Operator
	operatorList_lock sync.RWMutex

	//list of all docker images
	dockerImageList      map[string][]DockerImage
	dockerImageList_lock sync.RWMutex

	//list of all submitted topologies
	topologyList      map[string]*Topology
	topologyList_lock sync.RWMutex

	//list of all submitted topologies
	fogfunctionList      map[string]*FogFunction
	fogfunctionList_lock sync.RWMutex

	//to manage the orchestration of service topology
	serviceMgr *ServiceMgr

	//to manage the orchestration of tasks
	taskMgr *TaskMgr

	//number of deployed task
	curNumOfTasks int
	prevNumOfTask int
	counter_lock  sync.RWMutex

	//type of subscribed entities
	subID2Type map[string]string
}

func (master *Master) Start(configuration *Config) {
	master.cfg = configuration

	master.messageBus = configuration.GetMessageBus()
	master.discoveryURL = configuration.GetDiscoveryURL()

	master.workers = make(map[string]*WorkerProfile)

	master.operatorList = make(map[string]Operator)
	master.dockerImageList = make(map[string][]DockerImage)
	master.topologyList = make(map[string]*Topology)
	master.fogfunctionList = make(map[string]*FogFunction)

	master.subID2Type = make(map[string]string)

	// find a nearby IoT Broker
	for {
		nearby := NearBy{}
		nearby.Latitude = master.cfg.Location.Latitude
		nearby.Longitude = master.cfg.Location.Longitude
		nearby.Limit = 1

		client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL()}
		selectedBroker, err := client.DiscoveryNearbyIoTBroker(nearby)

		if err == nil && selectedBroker != "" {
			master.BrokerURL = selectedBroker
			break
		} else {
			if err != nil {
				ERROR.Println(err)
			}

			INFO.Println("continue to look up a nearby IoT broker")
			time.Sleep(5 * time.Second)
		}
	}

	// initialize the manager for both fog function and service topology
	master.taskMgr = NewTaskMgr(master)
	master.taskMgr.Init()

	master.serviceMgr = NewServiceMgr(master)
	master.serviceMgr.Init()

	// announce myself to the nearby IoT Broker
	master.registerMyself()

	// start the NGSI agent
	master.agent = &NGSIAgent{Port: configuration.Master.AgentPort}
	master.myURL = "http://" + configuration.InternalIP + ":" + strconv.Itoa(configuration.Master.AgentPort)
	master.agent.Start()
	master.agent.SetContextNotifyHandler(master.onReceiveContextNotify)
	master.agent.SetContextAvailabilityNotifyHandler(master.onReceiveContextAvailability)

	// start the message consumer
	go func() {
		cfg := MessageBusConfig{}
		cfg.Broker = configuration.GetMessageBus()
		cfg.Exchange = "fogflow"
		cfg.ExchangeType = "topic"
		cfg.DefaultQueue = master.id
		cfg.BindingKeys = []string{master.id + ".", "heartbeat.*"}

		// create the communicator with the broker info and topics
		master.communicator = NewCommunicator(&cfg)
		for {
			retry, err := master.communicator.StartConsuming(master.id, master)
			if retry {
				INFO.Printf("Going to retry launching the rabbitmq. Error: %v", err)
			} else {
				INFO.Printf("stop retrying")
				break
			}
		}
	}()

	master.prevNumOfTask = 0
	master.curNumOfTasks = 0

	// start a timer to do something periodically
	master.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			<-master.ticker.C
			//master.onTimer()
		}
	}()

	// subscribe to the update of required context information
	master.triggerInitialSubscriptions()
}

func (master *Master) onTimer() {
	master.counter_lock.Lock()
	delta := master.curNumOfTasks - master.prevNumOfTask
	fmt.Printf("# of orchestrated tasks = %d, throughput = %d/s\r\n", master.curNumOfTasks, delta)
	master.prevNumOfTask = master.curNumOfTasks
	master.counter_lock.Unlock()
}

func (master *Master) Quit() {
	INFO.Println("to stop the master")
	master.unregisterMyself()
	INFO.Println("unregister myself")
	master.ticker.Stop()
	INFO.Println("stop the timer")
	master.communicator.StopConsuming()
	INFO.Println("stop consuming the messages")
}

func (master *Master) registerMyself() {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = master.id
	ctxObj.Entity.Type = "Master"
	ctxObj.Entity.IsPattern = false

	ctxObj.Metadata = make(map[string]ValueObject)

	mylocation := Point{}
	mylocation.Latitude = master.cfg.Location.Latitude
	mylocation.Longitude = master.cfg.Location.Longitude
	ctxObj.Metadata["location"] = ValueObject{Type: "point", Value: mylocation}

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}
}

func (master *Master) unregisterMyself() {
	entity := EntityId{}
	entity.ID = master.id
	entity.Type = "Master"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	err := client.DeleteContext(&entity)
	if err != nil {
		ERROR.Println(err)
	}
}

func (master *Master) triggerInitialSubscriptions() {
	master.subscribeContextEntity("Operator")
	master.subscribeContextEntity("DockerImage")
	master.subscribeContextEntity("Topology")
	master.subscribeContextEntity("FogFunction")
	master.subscribeContextEntity("ServiceIntent")
	master.subscribeContextEntity("TaskIntent")
}

func (master *Master) subscribeContextEntity(entityType string) {
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}
	newEntity.Type = entityType
	newEntity.IsPattern = true
	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)
	subscription.Reference = master.myURL

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		ERROR.Println(err)
	}
	INFO.Println(sid)

	master.subID2Type[sid] = entityType
}

func (master *Master) onReceiveContextNotify(notifyCtxReq *NotifyContextRequest) {
	sid := notifyCtxReq.SubscriptionId
	stype := master.subID2Type[sid]

	if len(notifyCtxReq.ContextResponses) == 0 {
		return
	}

	contextObj := CtxElement2Object(&(notifyCtxReq.ContextResponses[0].ContextElement))

	switch stype {
	// registry of an operator
	case "Operator":
		master.handleOperatorRegistration(contextObj)

	// registry of a docker image
	case "DockerImage":
		master.handleDockerImageRegistration(contextObj)

	// topology to define service template
	case "Topology":
		master.handleTopologyUpdate(contextObj)

	// fog function that includes a pair of topology and intent
	case "FogFunction":
		master.handleFogFunctionUpdate(contextObj)

	// service orchestration
	case "ServiceIntent":
		master.serviceMgr.handleServiceIntentUpdate(contextObj)

	// task orchestration
	case "TaskIntent":
		master.taskMgr.handleTaskIntentUpdate(contextObj)
	}
}

//
// to handle the registry of operator
//
func (master *Master) handleOperatorRegistration(operatorCtxObj *ContextObject) {
	INFO.Println(operatorCtxObj)

	if operatorCtxObj.IsEmpty() {
		// does not handle the removal of operator
		return
	}

	var operator = Operator{}
	jsonText, _ := json.Marshal(operatorCtxObj.Attributes["operator"].Value.(map[string]interface{}))
	err := json.Unmarshal(jsonText, &operator)
	if err != nil {
		ERROR.Println("failed to read the given operator")
	} else {
		master.operatorList_lock.Lock()
		master.operatorList[operator.Name] = operator
		master.operatorList_lock.Unlock()
	}
}

//
// to handle the management of docker images
//
func (master *Master) handleDockerImageRegistration(dockerImageCtxObj *ContextObject) {
	INFO.Println(dockerImageCtxObj)

	if dockerImageCtxObj.IsEmpty() {
		// does not handle the removal of operator
		return
	}

	dockerImage := DockerImage{}
	dockerImage.OperatorName = dockerImageCtxObj.Attributes["operator"].Value.(string)
	dockerImage.ImageName = dockerImageCtxObj.Attributes["image"].Value.(string)
	dockerImage.ImageTag = dockerImageCtxObj.Attributes["tag"].Value.(string)
	dockerImage.TargetedHWType = dockerImageCtxObj.Attributes["hwType"].Value.(string)
	dockerImage.TargetedOSType = dockerImageCtxObj.Attributes["osType"].Value.(string)
	dockerImage.Prefetched = dockerImageCtxObj.Attributes["prefetched"].Value.(bool)

	master.dockerImageList_lock.Lock()
	master.dockerImageList[dockerImage.OperatorName] = append(master.dockerImageList[dockerImage.OperatorName], dockerImage)
	master.dockerImageList_lock.Unlock()

	if dockerImage.Prefetched == true {
		// inform all workers to prefetch this docker image in advance
		master.prefetchDockerImages(dockerImage)
	}
}

func (master *Master) prefetchDockerImages(image DockerImage) {
	workers := master.queryWorkers()

	for _, worker := range workers {
		workerID := worker.Entity.ID
		taskMsg := SendMessage{Type: "PREFETCH_IMAGE", RoutingKey: workerID + ".", From: master.id, PayLoad: image}
		master.communicator.Publish(&taskMsg)
	}
}

//
// to update the fog function list
//
func (master *Master) handleFogFunctionUpdate(fogfunctionCtxObj *ContextObject) {
	INFO.Println(fogfunctionCtxObj)

	// the fog function is going to be deleted
	if fogfunctionCtxObj.IsEmpty() {
		var eid = fogfunctionCtxObj.Entity.ID

		master.fogfunctionList_lock.RLock()
		fogfunction := master.fogfunctionList[eid]
		master.fogfunctionList_lock.RUnlock()

		DEBUG.Printf("%+v\r\n", fogfunction)

		// remove the service intent
		master.serviceMgr.removeServiceIntent(fogfunction.Intent.ID)

		// remove the service topology
		topology := fogfunction.Topology
		master.topologyList_lock.Lock()
		master.topologyList[topology.Name] = &topology
		master.topologyList_lock.Unlock()

		// remove this fog function entity
		master.fogfunctionList_lock.Lock()
		delete(master.fogfunctionList, eid)
		master.fogfunctionList_lock.Unlock()

		return
	}

	topology := Topology{}

	topologyJsonText, err := json.Marshal(fogfunctionCtxObj.Attributes["topology"].Value.(map[string]interface{}))
	if err != nil {
		ERROR.Println("the topology object is not defined")
		return
	}
	err = json.Unmarshal(topologyJsonText, &topology)
	if err != nil {
		ERROR.Println("the topology object is not correctly defined")
		return
	}

	intent := ServiceIntent{}

	intentJsonText, err := json.Marshal(fogfunctionCtxObj.Attributes["intent"].Value.(map[string]interface{}))
	if err != nil {
		ERROR.Println("the intent object is not defined")
		return
	}
	err = json.Unmarshal(intentJsonText, &intent)
	if err != nil {
		ERROR.Println("the intent object is not correctly defined")
		return
	}

	// allow the ID of this service intent
	intent.ID = fogfunctionCtxObj.Entity.ID

	fogfunction := FogFunction{}

	fogfunction.Id = fogfunctionCtxObj.Entity.ID
	fogfunction.Name = fogfunctionCtxObj.Attributes["name"].Value.(string)
	fogfunction.Topology = topology
	fogfunction.Intent = intent

	// add the service topology
	master.topologyList_lock.Lock()
	master.topologyList[topology.Name] = &topology
	master.topologyList_lock.Unlock()

	// handle the associated service intent
	master.serviceMgr.handleServiceIntent(&fogfunction.Intent)

	// create or update this fog function
	master.fogfunctionList_lock.Lock()
	master.fogfunctionList[fogfunction.Id] = &fogfunction
	master.fogfunctionList_lock.Unlock()

	INFO.Println(fogfunction)
}

//
// to update the topology list
//
func (master *Master) handleTopologyUpdate(topologyCtxObj *ContextObject) {
	INFO.Println(topologyCtxObj)

	if topologyCtxObj.IsEmpty() {
		// remove this service topology entity
		master.topologyList_lock.Lock()

		var eid = topologyCtxObj.Entity.ID

		// find which one has this id
		for _, topology := range master.topologyList {
			if topology.Id == eid {
				var name = topology.Name
				delete(master.topologyList, name)
				break
			}
		}

		master.topologyList_lock.Unlock()

		return
	}

	// create or update this service topology
	topology := Topology{}
	jsonText, _ := json.Marshal(topologyCtxObj.Attributes["template"].Value.(map[string]interface{}))
	err := json.Unmarshal(jsonText, &topology)
	if err == nil {
		INFO.Println(topology)

		topology.Id = topologyCtxObj.Entity.ID

		master.topologyList_lock.Lock()
		master.topologyList[topology.Name] = &topology
		master.topologyList_lock.Unlock()

		INFO.Println(topology)
	}

}

func (master *Master) getTopologyByName(name string) *Topology {
	// find the required topology object
	master.topologyList_lock.RLock()
	defer master.topologyList_lock.RUnlock()

	topology := master.topologyList[name]
	return topology
}

func (master *Master) queryWorkers() []*ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.Type = "Worker"
	entity.IsPattern = true
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	ctxObjects, err := client.QueryContext(&query)
	if err != nil {
		ERROR.Println(err)
	}

	return ctxObjects
}

func (master *Master) onReceiveContextAvailability(notifyCtxAvailReq *NotifyContextAvailabilityRequest) {
	INFO.Println("===========RECEIVE CONTEXT AVAILABILITY=========")
	DEBUG.Println(notifyCtxAvailReq)

	subID := notifyCtxAvailReq.SubscriptionId

	var action string
	switch notifyCtxAvailReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	}

	for _, registrationResp := range notifyCtxAvailReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration
		for _, entity := range registration.EntityIdList {
			// convert context registration to entity registration
			entityRegistration := master.contextRegistration2EntityRegistration(&entity, &registration)
			go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
		}
	}
}

func (master *Master) contextRegistration2EntityRegistration(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	ctxObj := master.RetrieveContextEntity(entityId.ID)

	if ctxObj == nil {
		entityRegistration.ID = entityId.ID
		entityRegistration.Type = entityId.Type
	} else {
		entityRegistration.ID = ctxObj.Entity.ID
		entityRegistration.Type = ctxObj.Entity.Type

		entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)
		for attrName, attrValue := range ctxObj.Attributes {
			attributeRegistration := ContextRegistrationAttribute{}
			attributeRegistration.Name = attrName
			attributeRegistration.Type = attrValue.Type

			entityRegistration.AttributesList[attrName] = attributeRegistration
		}

		entityRegistration.MetadataList = make(map[string]ContextMetadata)
		for metaname, ctxmeta := range ctxObj.Metadata {
			cm := ContextMetadata{}
			cm.Name = metaname
			cm.Type = ctxmeta.Type
			cm.Value = ctxmeta.Value

			entityRegistration.MetadataList[metaname] = cm
		}
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) contextRegistration2EntityRegistrationNew(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	entityRegistration.ID = entityId.ID
	entityRegistration.Type = entityId.Type

	entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)

	for _, attribute := range ctxRegistration.ContextRegistrationAttributes {
		attributeRegistration := ContextRegistrationAttribute{}
		attributeRegistration.Name = attribute.Name
		attributeRegistration.Type = attribute.Type

		entityRegistration.AttributesList[attribute.Name] = attributeRegistration
	}

	entityRegistration.MetadataList = make(map[string]ContextMetadata)
	for _, ctxmeta := range ctxRegistration.Metadata {
		cm := ContextMetadata{}
		cm.Name = ctxmeta.Name
		cm.Type = ctxmeta.Type
		cm.Value = ctxmeta.Value

		entityRegistration.MetadataList[ctxmeta.Name] = cm
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) subscribeContextAvailability(availabilitySubscription *SubscribeContextAvailabilityRequest) string {

	availabilitySubscription.Reference = master.myURL + "/notifyContextAvailability"

	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL()}
	subscriptionId, err := client.SubscribeContextAvailability(availabilitySubscription)
	if err != nil {
		ERROR.Println(err)
		return ""
	}

	return subscriptionId
}

func (master *Master) unsubscribeContextAvailability(sid string) {
	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL()}
	err := client.UnsubscribeContextAvailability(sid)
	if err != nil {
		ERROR.Println(err)
	}
}

//
// to deal with the communication between master and workers via rabbitmq
//
func (master *Master) Process(msg *RecvMessage) error {
	//INFO.Println("type ", msg.Type)

	switch msg.Type {
	case "heart_beat":
		profile := WorkerProfile{}
		err := json.Unmarshal(msg.PayLoad, &profile)
		if err == nil {
			master.onHeartbeat(msg.From, &profile)
		}

	case "task_update":
		update := TaskUpdate{}
		err := json.Unmarshal(msg.PayLoad, &update)
		if err == nil {
			master.onTaskUpdate(msg.From, &update)
		}
	}

	return nil
}

func (master *Master) onHeartbeat(from string, profile *WorkerProfile) {
	master.workerList_lock.Lock()

	workerID := profile.WID
	if worker, exist := master.workers[workerID]; exist {
		worker.Capacity = profile.Capacity
	} else {
		profile.Workload = 0
		master.workers[workerID] = profile
	}

	master.workerList_lock.Unlock()
}

func (master *Master) onTaskUpdate(from string, update *TaskUpdate) {
	INFO.Println("==task update=========")
	INFO.Println(update)

}

func (master *Master) DeployTask(taskInstance *ScheduledTaskInstance) {
	master.counter_lock.Lock()
	master.curNumOfTasks = master.curNumOfTasks + 1
	master.counter_lock.Unlock()

	taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
	INFO.Println(taskMsg)

	// update the workload of this worker
	workerID := taskInstance.WorkerID

	master.workerList_lock.Lock()
	workerProfile := master.workers[workerID]
	workerProfile.Workload = workerProfile.Workload + 1
	master.workerList_lock.Unlock()

	go master.communicator.Publish(&taskMsg)
}

func (master *Master) TerminateTask(taskInstance *ScheduledTaskInstance) {
	taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) AddInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "ADD_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.id, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) RemoveInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "REMOVE_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.id, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

//
// the shared functions for function manager and topology manager to call
//
func (master *Master) RetrieveContextEntity(eid string) *ContextObject {
	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	ctxObj, err := client.GetEntity(eid)

	if err != nil {
		return nil
	}

	return ctxObj
}

//
// to select the right docker image of an operator for the selected worker
//
func (master *Master) DetermineDockerImage(operatorName string, wID string) string {
	INFO.Println("select a suitable image to execute on the selected worker")

	master.workerList_lock.RLock()
	wProfile := master.workers[wID]
	master.workerList_lock.RUnlock()

	if wProfile == nil {
		ERROR.Println("could not find this worker from the curent worker list: ", wID)
		return ""
	}

	//select a suitable image to execute on the selected worker
	selectedDockerImageName := ""

	master.dockerImageList_lock.RLock()
	for _, image := range master.dockerImageList[operatorName] {
		DEBUG.Println(image)
		DEBUG.Println(wProfile)

		hwType := "X86"
		osType := "Linux"

		if wProfile.HWType == "arm" {
			hwType = "ARM"
		}

		if wProfile.OSType == "linux" {
			osType = "Linux"
		}

		if image.TargetedOSType == osType && image.TargetedHWType == hwType {
			selectedDockerImageName = image.ImageName + ":" + image.ImageTag
			break
		}
	}

	master.dockerImageList_lock.RUnlock()

	DEBUG.Println(selectedDockerImageName)

	return selectedDockerImageName
}

func (master *Master) GetOperatorParamters(operatorName string) []Parameter {
	master.operatorList_lock.RLock()

	operator := master.operatorList[operatorName]
	parameters := make([]Parameter, len(operator.Parameters))
	copy(parameters, operator.Parameters)

	master.operatorList_lock.RUnlock()

	return parameters
}

//
// to select the worker that is closest to the given points
//
func (master *Master) SelectWorker(locations []Point) string {
	master.workerList_lock.RLock()
	defer master.workerList_lock.RUnlock()

	if len(locations) == 0 {
		for _, worker := range master.workers {
			return worker.WID
		}
		return ""
	}

	DEBUG.Printf("points: %+v\r\n", locations)

	// select the workers with the closest distance and also the worker is currently not overloaded
	closestWorkerID := ""
	closestTotalDistance := uint64(18446744073709551615)
	for _, worker := range master.workers {
		INFO.Printf("check worker %+v\r\n", worker)

		// if this worker is already overloaded, check the next one
		if worker.IsOverloaded() == true {
			continue
		}

		wp := Point{}
		wp.Latitude = worker.PLocation.Latitude
		wp.Longitude = worker.PLocation.Longitude

		totalDistance := uint64(0)

		for _, location := range locations {
			if location.IsEmpty() == true {
				continue
			}

			distance := Distance(wp, location)
			totalDistance += distance
			INFO.Printf("distance = %d between %+v, %+v\r\n", distance, wp, location)
		}

		if totalDistance < closestTotalDistance {
			closestWorkerID = worker.WID
			closestTotalDistance = totalDistance
		}

		INFO.Println("closest worker ", closestWorkerID, " with the closest distance ", closestTotalDistance)
	}

	// select the one with lowest capacity if there are more than one with the closest distance

	return closestWorkerID
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(p1 Point, p2 Point) uint64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = p1.Latitude * math.Pi / 180
	lo1 = p1.Longitude * math.Pi / 180
	la2 = p2.Latitude * math.Pi / 180
	lo2 = p2.Longitude * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return uint64(2 * r * math.Asin(math.Sqrt(h)))
}
