package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	. "fogflow/common/communicator"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"

	. "fogflow/common/config"
)

const MAX_HEARTBEAT_DURATION = 60 // in seconds

type Master struct {
	cfg *Config

	BrokerURL string

	id           string
	myURL        string
	messageBus   string
	discoveryURL string
	designerURL  string

	SecurityCfg *HTTPS

	communicator  *Communicator
	communicator2 *Communicator
	ticker        *time.Ticker
	agent         *NGSIAgent

	//list of all workers
	workers         map[string]*WorkerProfile
	workerList_lock sync.RWMutex

	//list of all operators
	operatorList      map[string]Operator
	operatorList_lock sync.RWMutex

	//list of all docker images
	// dockerImageList      map[string][]DockerImage
	// dockerImageList_lock sync.RWMutex

	// //list of all submitted topologies
	// topologyList      map[string]*Topology
	// topologyList_lock sync.RWMutex

	// //list of all submitted topologies
	// fogfunctionList      map[string]*FogFunction
	// fogfunctionList_lock sync.RWMutex

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
	master.SecurityCfg = &configuration.HTTPS

	master.messageBus = configuration.GetMessageBus()
	master.discoveryURL = configuration.GetDiscoveryURL()
	master.designerURL = configuration.GetDesignerURL()

	master.workers = make(map[string]*WorkerProfile)

	master.operatorList = make(map[string]Operator)
	// master.dockerImageList = make(map[string][]DockerImage)
	// master.topologyList = make(map[string]*Topology)
	// master.fogfunctionList = make(map[string]*FogFunction)

	master.subID2Type = make(map[string]string)

	// communicate with the cloud_broker
	master.BrokerURL = configuration.GetBrokerURL()
	INFO.Println("communicate with the cloud broker via ", master.BrokerURL)

	// initialize the manager for both fog function and service topology
	master.taskMgr = NewTaskMgr(master)
	master.taskMgr.Init()

	master.serviceMgr = NewServiceMgr(master)
	master.serviceMgr.Init()

	master.myURL = "http://" + configuration.GetMasterIP() + ":" + strconv.Itoa(configuration.Master.AgentPort)

	// start the NGSI agent
	master.agent = &NGSIAgent{Port: configuration.Master.AgentPort, SecurityCfg: master.cfg.HTTPS}
	master.agent.Start()
	master.agent.SetContextAvailabilityNotifyHandler(master.onReceiveContextAvailability)

	cfg := MessageBusConfig{}
	cfg.Broker = configuration.GetMessageBus()
	cfg.Exchange = "fogflow"
	cfg.ExchangeType = "topic"
	cfg.DefaultQueue = master.id
	cfg.BindingKeys = []string{master.id + ".", "heartbeat.*", "orchestration.*"}

	// create the communicator with the broker info and topics
	master.communicator = NewCommunicator(&cfg)

	// start the message consumer
	go func() {
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
			master.onTimer()
		}
	}()

	master.registerMyself()
}

func (master *Master) onTimer() {
	master.counter_lock.Lock()
	// delta := master.curNumOfTasks - master.prevNumOfTask
	// fmt.Printf("# of orchestrated tasks = %d, throughput = %d/s\r\n", master.curNumOfTasks, delta)
	master.prevNumOfTask = master.curNumOfTasks
	master.counter_lock.Unlock()

	// check the liveness of each worker
	master.workerList_lock.Lock()
	for k, w := range master.workers {
		if w.IsLive(MAX_HEARTBEAT_DURATION) == false {
			workerID := w.WID
			delete(master.workers, k)
			INFO.Println("REMOVE worker " + workerID + " from the list")
		}
	}

	master.workerList_lock.Unlock()
}

func (master *Master) Quit() {
	INFO.Println("to stop the master")
	master.unregisterMyself()
	INFO.Println("unregister myself")
	master.ticker.Stop()
	INFO.Println("stop the timer")
	master.communicator.StopConsuming()
	INFO.Println("stop consuming the message")
}

func (master *Master) registerMyself() {
	profile := MasterProfile{}
	profile.WID = master.id
	profile.PLocation = master.cfg.Location
	profile.AgentURL = master.myURL

	taskMsg := SendMessage{Type: "MASTER_JOIN", RoutingKey: "designer.", From: master.id, PayLoad: profile}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) unregisterMyself() {
	profile := MasterProfile{}
	profile.WID = master.id

	taskMsg := SendMessage{Type: "MASTER_LEAVE", RoutingKey: "designer.", From: master.id, PayLoad: profile}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) prefetchDockerImages(image DockerImage) {
	master.workerList_lock.RLock()
	defer master.workerList_lock.RUnlock()

	for _, worker := range master.workers {
		workerID := worker.WID
		taskMsg := SendMessage{Type: "PREFETCH_IMAGE", RoutingKey: workerID + ".", From: master.id, PayLoad: image}
		INFO.Println(taskMsg)
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

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
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
		//entityRegistration := EntityRegistration{}
		for _, entity := range registration.EntityIdList {
			// convert context registration to entity registration
			fmt.Println("entity.MsgFormat", entity.MsgFormat)
			if entity.MsgFormat == "NGSILD" {
				entityRegistration := master.ldContextRegistration2EntityRegistration(&entity, &registration)
				go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
			} else {
				entityRegistration := master.contextRegistration2EntityRegistration(&entity, &registration)
				go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
			}
			//go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
		}
	}
}

func (master *Master) RetrieveContextLdEntity(eid string, fsp string) interface{} {
	query := LDQueryContextRequest{}

	query.Entities = make([]EntityId, 0)
	query.Type = "Query"
	entity := EntityId{}
	idSplit := strings.Split(eid, "@")
	entity.ID = idSplit[0]
	entity.IsPattern = false
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	ctxObjects, err := client.QueryLdContext(&query, idSplit[1], fsp)
	if err == nil && ctxObjects != nil && len(ctxObjects) > 0 {
		return ctxObjects[0]
	} else {
		if err != nil {
			ERROR.Println("error occured when retrieving a context entity :", err)
		}

		return nil
	}
}

func (master *Master) ldContextRegistration2EntityRegistration(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	ctxObj := master.RetrieveContextLdEntity(entityId.ID, entityId.FiwareServicePath)
	if ctxObj == nil {
		entityRegistration.ID = entityId.ID
		entityRegistration.Type = entityId.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
	} else {
		ldCtcObj := ctxObj.(map[string]interface{})
		entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)
		entityRegistration.MetadataList = make(map[string]ContextMetadata)
		entityRegistration.MsgFormat = entityId.MsgFormat
		for key, attr := range ldCtcObj {
			if key != "modifiedAt" && key != "createdAt" && key != "observationSpace" && key != "operationSpace" && key != "@context" && key != "fiwareServicePath" {
				if key == "id" {
					entityRegistration.ID = entityId.ID
				} else if key == "type" {
					entityRegistration.Type = ldCtcObj[key].(string)
				} else if key == "FiwareServicePath" {
					entityRegistration.FiwareServicePath = ldCtcObj[key].(string)
				} else {
					attrmap := attr.(map[string]interface{})
					if attrmap["type"] != "GeoProperty" {
						attributeRegistration := ContextRegistrationAttribute{}
						attributeRegistration.Name = key
						attributeRegistration.Type = attrmap["type"].(string)
						entityRegistration.AttributesList[key] = attributeRegistration
					} else {
						metaData := attr.(map[string]interface{})
						cm := ContextMetadata{}
						cm.Name = key
						matadataCordinate := metaData["value"].(map[string]interface{})
						typ, points := GetNGSIV1DomainMetaData(matadataCordinate["type"].(string), matadataCordinate["coordinates"])
						cm.Type = typ
						cm.Value = points
						entityRegistration.MetadataList[key] = cm
					}
				}
			}
		}
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) contextRegistration2EntityRegistration(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	ctxObj := master.RetrieveContextEntity(entityId.ID)
	if ctxObj == nil {
		entityRegistration.ID = entityId.ID
		entityRegistration.Type = entityId.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
	} else {
		entityRegistration.ID = ctxObj.Entity.ID
		entityRegistration.Type = ctxObj.Entity.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
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

	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL(), SecurityCfg: &master.cfg.HTTPS}
	subscriptionId, err := client.SubscribeContextAvailability(availabilitySubscription)
	if err != nil {
		ERROR.Println(err)
		return ""
	}

	return subscriptionId
}

func (master *Master) unsubscribeContextAvailability(sid string) {
	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL(), SecurityCfg: &master.cfg.HTTPS}
	err := client.UnsubscribeContextAvailability(sid)
	if err != nil {
		ERROR.Println(err)
	}
}

//
// to deal with the communication between master and workers via rabbitmq
//
func (master *Master) Process(msg *RecvMessage) error {
	switch msg.Type {
	case "WORKER_HEARTBEAT":
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

	// case "FogFunction":
	// 	master.handleFogFunctionUpdate(msg.PayLoad)

	case "ServiceIntent":
		serviceIntent := ServiceIntent{}
		err := json.Unmarshal(msg.PayLoad, &serviceIntent)
		if err == nil {
			master.serviceMgr.handleServiceIntentUpdate(&serviceIntent)
		}
	}

	return nil
}

func (master *Master) onHeartbeat(from string, profile *WorkerProfile) {
	master.workerList_lock.Lock()

	workerID := profile.WID
	fmt.Println("**** workerID and profile ******", workerID, profile)
	if worker, exist := master.workers[workerID]; exist {
		worker.Capacity = profile.Capacity
		worker.Last_Heartbeat_Update = time.Now()
	} else {
		profile.Workload = 0
		profile.Last_Heartbeat_Update = time.Now()
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

	go master.communicator.Publish(&taskMsg)

	// update the workload of this worker
	workerID := taskInstance.WorkerID

	master.workerList_lock.Lock()
	workerProfile := master.workers[workerID]
	workerProfile.Workload = workerProfile.Workload + 1
	master.workerList_lock.Unlock()
}

func (master *Master) TerminateTask(taskInstance *ScheduledTaskInstance) {
	master.workerList_lock.Lock()
	defer master.workerList_lock.Unlock()

	// update the workload of this worker
	workerID := taskInstance.WorkerID

	workerProfile := master.workers[workerID]
	if workerProfile != nil {
		workerProfile.Workload = workerProfile.Workload - 1

		taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
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
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.ID = eid
	entity.IsPattern = false
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	ctxObjects, err := client.QueryContext(&query)
	if err == nil && ctxObjects != nil && len(ctxObjects) > 0 {
		return ctxObjects[0]
	} else {
		if err != nil {
			ERROR.Println("error occured when retrieving a context entity :", err)
		}

		return nil
	}
}

func (master *Master) GetWorkerList(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(master.workers)
}

func (master *Master) GetTaskList(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(master.workers)
}

func (master *Master) GetStatus(w rest.ResponseWriter, r *rest.Request) {
	profile := MasterProfile{}
	profile.WID = master.id
	profile.PLocation = master.cfg.Location
	profile.AgentURL = master.myURL

	w.WriteJson(profile)
}

//
// to select the worker that is closest to the given points
//
func (master *Master) SelectWorker(locations []Point) string {
	master.workerList_lock.RLock()
	defer master.workerList_lock.RUnlock()
	fmt.Println("&&&& len(locations) &&&&&&&&&", len(locations))
	if len(locations) == 0 {
		for _, worker := range master.workers {
			return worker.WID
		}
		return ""
	}

	DEBUG.Printf("points: %+v\r\n", locations)
	fmt.Println("&&&& master.workers &&&&&&", master.workers)

	// select the workers with the closest distance and also the worker is currently not overloaded
	closestWorkerID := ""
	closestTotalDistance := uint64(18446744073709551615)
	for _, worker := range master.workers {
		fmt.Println("***** master.worker *******", worker)
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

			distance := Distance(&wp, &location)
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

//
// query the topology from Designer based on the given name
//
func (master *Master) getTopologyByName(name string) *Topology {
	designerURL := fmt.Sprintf("%s/topology/%s", master.cfg.GetDesignerURL(), name)
	fmt.Println(designerURL)

	req, err1 := http.NewRequest(http.MethodGet, designerURL, nil)
	if err1 != nil {
		ERROR.Println(err1)
		return nil
	}

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	resp, err2 := client.Do(req)
	if err2 != nil {
		ERROR.Println(err2)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	INFO.Println("response Body:", string(body))

	topology := Topology{}
	jsonErr := json.Unmarshal(body, &topology)
	if jsonErr != nil {
		ERROR.Println(jsonErr)
		return nil
	}

	INFO.Printf("%+v", topology)

	// // update the list of operators
	master.operatorList_lock.Lock()
	for _, operator := range topology.Operators {
		master.operatorList[operator.Name] = operator
	}
	master.operatorList_lock.Unlock()

	return &topology
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

	master.operatorList_lock.RLock()
	defer master.operatorList_lock.RUnlock()

	operator := master.operatorList[operatorName]

	dockerimages := operator.DockerImages

	for _, image := range dockerimages {
		fmt.Println("*****image*******", image)
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

func (master *Master) subscribeContextEntity(entityType string) {
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}
	newEntity.Type = entityType
	newEntity.IsPattern = true
	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)
	subscription.Reference = master.myURL

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		ERROR.Println(err)
	}
	INFO.Println(sid)

	master.subID2Type[sid] = entityType
}
