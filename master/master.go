package main

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	. "fogflow/common/communicator"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type Master struct {
	cfg *Config

	BrokerURL string

	myID         string
	myURL        string
	messageBus   string
	discoveryURL string

	communicator *Communicator
	ticker       *time.Ticker
	agent        *NGSIAgent

	//list of all workers
	workers         map[string]*WorkerProfile
	workerList_lock sync.RWMutex

	//list of all dockerized operators
	operatorList      map[string][]DockerImage
	operatorList_lock sync.RWMutex

	//to manage the orchestration of fog functions
	functionMgr *FunctionMgr

	//to manage the orchestration of service topology
	topologyMgr *TopologyMgr

	//type of subscribed entities
	subID2Type map[string]string
}

func (master *Master) Start(configuration *Config) {
	master.cfg = configuration

	master.messageBus = configuration.MessageBus
	master.discoveryURL = configuration.IoTDiscoveryURL

	master.workers = make(map[string]*WorkerProfile)

	master.operatorList = make(map[string][]DockerImage)

	master.subID2Type = make(map[string]string)

	// find a nearby IoT Broker
	for {
		nearby := NearBy{}
		nearby.Latitude = master.cfg.PLocation.Latitude
		nearby.Longitude = master.cfg.PLocation.Longitude
		nearby.Limit = 1

		client := NGSI9Client{IoTDiscoveryURL: master.cfg.IoTDiscoveryURL}
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
	master.functionMgr = NewFogFunctionMgr(master)
	master.functionMgr.Init()

	master.topologyMgr = NewTopologyMgr(master)
	master.topologyMgr.Init()

	// announce myself to the nearby IoT Broker
	master.registerMyself()

	// start the NGSI agent
	master.agent = &NGSIAgent{Port: configuration.AgentPort}
	master.myURL = "http://" + configuration.MyIP + ":" + strconv.Itoa(configuration.AgentPort)
	master.agent.Start()
	master.agent.SetContextNotifyHandler(master.onReceiveContextNotify)
	master.agent.SetContextAvailabilityNotifyHandler(master.onReceiveContextAvailability)

	// start the message consumer
	go func() {
		cfg := MessageBusConfig{}
		cfg.Broker = configuration.MessageBus
		cfg.Exchange = "fogflow"
		cfg.ExchangeType = "topic"
		cfg.DefaultQueue = "master" + master.myID
		cfg.BindingKeys = []string{"master." + master.myID + ".", "heartbeat.*"}

		// create the communicator with the broker info and topics
		master.communicator = NewCommunicator(&cfg)
		for {
			retry, err := master.communicator.StartConsuming("master"+master.myID, master)
			if retry {
				INFO.Printf("Going to retry launching the rabbitmq. Error: %v", err)
			} else {
				INFO.Printf("stop retrying")
				break
			}
		}
	}()

	// start a timer to do something periodically
	master.ticker = time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-master.ticker.C
			master.onTimer()
		}
	}()

	// subscribe to the update of required context information
	master.triggerInitialSubscriptions()
}

func (master *Master) onTimer() {

}

func (master *Master) Quit() {
	INFO.Println("to stop the master")
	master.unregisterMyself()
	master.communicator.StopConsuming()
	master.ticker.Stop()
	INFO.Println("stop consuming the messages")
}

func (master *Master) registerMyself() {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = "SysComponent.Master." + master.myID
	ctxObj.Entity.Type = "Master"
	ctxObj.Entity.IsPattern = false

	ctxObj.Metadata = make(map[string]ValueObject)

	mylocation := Point{}
	mylocation.Latitude = master.cfg.PLocation.Latitude
	mylocation.Longitude = master.cfg.PLocation.Longitude
	ctxObj.Metadata["location"] = ValueObject{Type: "point", Value: mylocation}

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	err := client.UpdateContext(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}
}

func (master *Master) unregisterMyself() {
	entity := EntityId{}
	entity.ID = "Master." + master.myID
	entity.Type = "Master"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	err := client.DeleteContext(&entity)
	if err != nil {
		ERROR.Println(err)
	}
}

func (master *Master) triggerInitialSubscriptions() {
	master.subscribeContextEntity("Topology")
	master.subscribeContextEntity("Requirement")
	master.subscribeContextEntity("DockerImage")
	master.subscribeContextEntity("FogFunction")
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

	DEBUG.Println("NGSI10 NOTIFY ", sid, " , ", stype)

	switch stype {
	case "DockerImage":
		master.handleDockerImageRegistration(notifyCtxReq.ContextResponses, sid)

	//output-driven service orchestration for service topology
	case "Topology":
		master.topologyMgr.handleTopologyUpdate(notifyCtxReq.ContextResponses, sid)
	case "Requirement":
		master.topologyMgr.handleRequirementUpdate(notifyCtxReq.ContextResponses, sid)

	//input-driven service orchestration for serverless function
	case "FogFunction":
		master.functionMgr.handleFogFunctionUpdate(notifyCtxReq.ContextResponses, sid)
	}
}

//
// to handle the management of docker images
//
func (master *Master) handleDockerImageRegistration(responses []ContextElementResponse, sid string) {
	fetchedImageList := make([]DockerImage, 0)

	for _, response := range responses {
		dockerImageCtxObj := CtxElement2Object(&(response.ContextElement))
		//INFO.Printf("%+v\r\n", dockerImageCtxObj)

		dockerImage := DockerImage{}
		dockerImage.OperatorName = dockerImageCtxObj.Attributes["operator"].Value.(string)
		dockerImage.ImageName = dockerImageCtxObj.Attributes["image"].Value.(string)
		dockerImage.ImageTag = dockerImageCtxObj.Attributes["tag"].Value.(string)
		dockerImage.TargetedHWType = dockerImageCtxObj.Attributes["hwType"].Value.(string)
		dockerImage.TargetedOSType = dockerImageCtxObj.Attributes["osType"].Value.(string)
		dockerImage.Prefetched = dockerImageCtxObj.Attributes["prefetched"].Value.(bool)

		master.operatorList_lock.Lock()
		master.operatorList[dockerImage.OperatorName] = append(master.operatorList[dockerImage.OperatorName], dockerImage)
		master.operatorList_lock.Unlock()

		if dockerImage.Prefetched == true {
			// inform all workers to prefetch this docker image in advance
			fetchedImageList = append(fetchedImageList, dockerImage)
		}
	}

	if len(fetchedImageList) > 0 {
		master.prefetchDockerImages(fetchedImageList)
	}
}

func (master *Master) prefetchDockerImages(imageList []DockerImage) {
	workers := master.queryWorkers()

	for _, worker := range workers {
		workerID := worker.Entity.ID
		taskMsg := SendMessage{Type: "prefetch_image", RoutingKey: workerID + ".", From: master.myID, PayLoad: imageList}
		master.communicator.Publish(&taskMsg)
	}
}

func (master *Master) queryWorkers() []*ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.Type = "Worker"
	entity.IsPattern = true
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	ctxObjects, err := client.QueryContext(&query, nil)
	if err != nil {
		ERROR.Println(err)
	}

	return ctxObjects
}

func (master *Master) onReceiveContextAvailability(notifyCtxAvailReq *NotifyContextAvailabilityRequest) {
	INFO.Println("===========RECEIVE CONTEXT AVAILABILITY=========")

	DEBUG.Print("received raw availability notify: %+v\r\n", notifyCtxAvailReq)

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
			master.functionMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
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

	DEBUG.Print("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) subscribeContextAvailability(availabilitySubscription *SubscribeContextAvailabilityRequest) string {

	availabilitySubscription.Reference = master.myURL + "/notifyContextAvailability"

	client := NGSI9Client{IoTDiscoveryURL: master.cfg.IoTDiscoveryURL}
	subscriptionId, err := client.SubscribeContextAvailability(availabilitySubscription)
	if err != nil {
		ERROR.Println(err)
		return ""
	}

	return subscriptionId
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
	master.workers[profile.WID] = profile
	master.workerList_lock.Unlock()
}

func (master *Master) onTaskUpdate(from string, update *TaskUpdate) {
	INFO.Println("==task update=========")
	INFO.Println(update)
}

//
// to carry out deployment actions given by the orchestrators of fog functions and service topologies
//
func (master *Master) DeployTasks(taskInstances []*ScheduledTaskInstance) {
	for _, pScheduledTaskInstance := range taskInstances {
		// convert the operator name into the name of a proper docker image for the assigned worker
		operatorName := (*pScheduledTaskInstance).DockerImage
		assignedWorkerID := (*pScheduledTaskInstance).WorkerID
		(*pScheduledTaskInstance).DockerImage = master.DetermineDockerImage(operatorName, assignedWorkerID)

		taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: pScheduledTaskInstance.WorkerID + ".", From: master.myID, PayLoad: *pScheduledTaskInstance}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
}

func (master *Master) TerminateTasks(instances []*ScheduledTaskInstance) {
	INFO.Println("to terminate all scheduled tasks, ", len(instances))
	for _, instance := range instances {
		taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: instance.WorkerID + ".", From: master.myID, PayLoad: *instance}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
}

func (master *Master) DeployTask(taskInstance *ScheduledTaskInstance) {
	taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.myID, PayLoad: *taskInstance}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) TerminateTask(taskInstance *ScheduledTaskInstance) {
	taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.myID, PayLoad: *taskInstance}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) AddInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "ADD_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.myID, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) RemoveInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "REMOVE_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.myID, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

//
// the shared functions for function manager and topology manager to call
//
func (master *Master) RetrieveContextEntity(eid string) *ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.ID = eid
	entity.IsPattern = false
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL}
	ctxObjects, err := client.QueryContext(&query, nil)
	if err == nil && ctxObjects != nil && len(ctxObjects) > 0 {
		return ctxObjects[0]
	} else {
		if err != nil {
			ERROR.Println("error occured when retrieving a context entity :", err)
		}

		return nil
	}
}

//
// to select the right docker image of an operator for the selected worker
//
func (master *Master) DetermineDockerImage(operatorName string, wID string) string {
	selectedDockerImageName := ""

	wProfile := master.workers[wID]
	master.operatorList_lock.RLock()
	for _, image := range master.operatorList[operatorName] {
		DEBUG.Println(image.TargetedOSType, image.TargetedHWType)
		DEBUG.Println(wProfile.OSType, wProfile.HWType)
		if image.TargetedOSType == wProfile.OSType && image.TargetedHWType == wProfile.HWType {
			selectedDockerImageName = image.ImageName + ":" + image.ImageTag
		}
	}

	master.operatorList_lock.RUnlock()

	return selectedDockerImageName
}

//
// to select the worker that is closest to the given points
//
func (master *Master) SelectWorker(locations []Point) string {
	if len(locations) == 0 {
		for _, worker := range master.workers {
			return worker.WID
		}

		return ""
	}

	closestWorkerID := ""
	closestTotalDistance := uint64(18446744073709551615)
	for _, worker := range master.workers {
		INFO.Printf("check worker %+v\r\n", worker)

		wp := Point{}
		wp.Latitude = worker.PLocation.Latitude
		wp.Longitude = worker.PLocation.Longitude

		totalDistance := uint64(0)

		for _, location := range locations {
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

	return closestWorkerID
}
