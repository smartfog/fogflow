package main

import (
	"bytes"
	"encoding/json"
	"gonum.org/v1/gonum/mat"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	. "github.com/smartfog/fogflow/common/communicator"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"

	. "github.com/smartfog/fogflow/common/config"
	"math"
	url2 "net/url"
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


	edgeUtilization  map[string] []float64
	edgeUtilLock     sync.RWMutex

	METRIC_COUNT     int
	metricWeights    []float64

	workerStats      map[string]*WorkerStat
	workerStats_lock sync.RWMutex

	// The moving average window for updating logs
	statAlpha float32

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
	master.statAlpha = 0.2
	master.METRIC_COUNT = 4
	master.messageBus = configuration.GetMessageBus()
	master.discoveryURL = configuration.GetDiscoveryURL()

	master.workers = make(map[string]*WorkerProfile)

	master.workerStats = make(map[string]*WorkerStat)

	master.operatorList = make(map[string][]DockerImage)

	master.subID2Type = make(map[string]string)
	master.edgeUtilLock.Lock()
	master.edgeUtilization = make( map[string] []float64)
	master.edgeUtilLock.Unlock()

	//Build the Metrics Vector from Preferences Matrix Using AHP Method
	master.InitMetricWeights()
	// find a nearby IoT Broker
	for {
		nearby := NearBy{}
		nearby.Latitude = master.cfg.PLocation.Latitude
		nearby.Longitude = master.cfg.PLocation.Longitude
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
	master.functionMgr = NewFogFunctionMgr(master)
	master.functionMgr.Init()

	master.topologyMgr = NewTopologyMgr(master)
	master.topologyMgr.Init()

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
	master.InitPrometheus(configuration)
}

func (master *Master) InitPrometheus(configuration *Config) {
	INFO.Printf("Initializing Prometheus on %v with AdminPort: %v, DataProt: %v \n",
		configuration.Prometheus.Address,
		configuration.Prometheus.AdminPort,configuration.Prometheus.DataPort)
	if configuration.Prometheus.Address=="auto" {
		configuration.Prometheus.Address = "prometheus"
	}
	//TODO check the address to see if prometheus is up and running
}

func (master *Master) onTimer() {
	master.UpdateUtilization()
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

	ctxObj.Entity.ID = master.id
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
		taskMsg := SendMessage{Type: "prefetch_image", RoutingKey: workerID + ".", From: master.id, PayLoad: imageList}
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
	case "heart_stat":
		//INFO.Println("AMIR: recieved Heart stats:",msg.PayLoad)
		stat := WorkerStat{}
		err := json.Unmarshal(msg.PayLoad,&stat)
		if err == nil {
			master.onNewStat(msg.From, &stat)
		}

	}

	return nil
}

func (master *Master) onHeartbeat(from string, profile *WorkerProfile) {
	//INFO.Printf("HEARTBEAT: The edge address of worker %v is: %v ", profile.WID,profile.EdgeAddress)
	if _,exists := master.workers[profile.WID]; !exists {
		//A new worker registration!



		//Create Utilization Entry for the worker
		DEBUG.Printf("Registering new Edge: %v",profile.EdgeAddress)
		master.edgeUtilLock.Lock()
		master.edgeUtilization[profile.EdgeAddress+ ":" + strconv.Itoa(profile.CAdvisorPort)]=
			make([]float64,master.METRIC_COUNT,master.METRIC_COUNT)
			master.edgeUtilLock.Unlock()
		//Add the worker to the workers
		master.workerList_lock.Lock()
		master.workers[profile.WID] = profile
		master.workerList_lock.Unlock()

		// Notify Prometheus
		master.UpdatePrometheusConfig()

	}else {
		master.workerList_lock.Lock()
		master.workers[profile.WID] = profile
		master.workerList_lock.Unlock()
	}
	master.UpdateUtilization()
	//DEBUG.Printf("List of workers: %v",master.workers)
}

func (master *Master) onNewStat(from string, newStat *WorkerStat){
	if oldStat, ok := master.workerStats[newStat.WID]; ok {
		//update statistics
		//INFO.Println("not the first time, old stat is: ",oldStat)
		master.workerStats_lock.Lock()
		oldStat.UtilMemory = (newStat.UtilMemory * master.statAlpha)+( (1-master.statAlpha)*oldStat.UtilMemory )
		oldStat.UtilCPU = (newStat.UtilCPU * master.statAlpha)+( (1-master.statAlpha)*oldStat.UtilCPU )
		master.workerStats[newStat.WID]=oldStat
		master.workerStats_lock.Unlock()
		//INFO.Println("and the new stat is: ",master.workerStats[newStat.WID])

	} else {
		// First time that we have this value. Added to map
		master.workerStats[newStat.WID]=newStat
		INFO.Println("first time....")
		INFO.Println("Updating:",newStat)
	}

	if newStat.UtilMemory>0.96 {
		INFO.Println("received memory utilization more than 95%!")
		master.MoveAFunction()
	}

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

		taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: pScheduledTaskInstance.WorkerID + ".", From: master.id, PayLoad: *pScheduledTaskInstance}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
}

func (master *Master) TerminateTasks(instances []*ScheduledTaskInstance) {
	INFO.Println("to terminate all scheduled tasks, ", len(instances))
	for _, instance := range instances {
		taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: instance.WorkerID + ".", From: master.id, PayLoad: *instance}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
}

func (master *Master) DeployTask(taskInstance *ScheduledTaskInstance) {
	// convert the operator name into the name of a proper docker image for the assigned worker
	operatorName := (*taskInstance).DockerImage
	assignedWorkerID := (*taskInstance).WorkerID
	(*taskInstance).DockerImage = master.DetermineDockerImage(operatorName, assignedWorkerID)

	taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
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
	INFO.Println("select a suitable image to execute on the selected worker")

	selectedDockerImageName := ""

	wProfile := master.workers[wID]
	master.operatorList_lock.RLock()
	for _, image := range master.operatorList[operatorName] {
		DEBUG.Println(image.TargetedOSType, image.TargetedHWType)
		DEBUG.Println(wProfile.OSType, wProfile.HWType)

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

	master.operatorList_lock.RUnlock()

	DEBUG.Println(selectedDockerImageName)

	return selectedDockerImageName
}

//
// to select the worker that is closest to the given points
//

func (master *Master) SelectWorkerAHP() string {
	var selectedWorkerID string = "0"
	var highestUtility float64 =0
	for _,worker := range master.workers{
		var utility  float64
		utility = 0
		instanceID:= worker.EdgeAddress+":"+strconv.Itoa(worker.CAdvisorPort)
		for i:=0; i<master.METRIC_COUNT ; i++ {
			utility = utility + master.metricWeights[i]*master.edgeUtilization[instanceID][i]
		}
		if utility >=highestUtility {
			selectedWorkerID=worker.WID
		}

	}

	return selectedWorkerID
	//USEFUL VARIABLES:
	//master.metricWeights
	//master.edgeUtilization
}

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
			if location.Latitude == 0 && location.Longitude == 0 {
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

	return closestWorkerID
}


func (master *Master)MoveAFunction() bool{
	for _,fogflow := range master.functionMgr.functionFlows{
		for _, taskinstance:= range fogflow.DeploymentPlan{
			taskinstance.WorkerID=master.GetAnotherWorker(taskinstance.WorkerID)
			master.DeployTask(taskinstance)
			INFO.Println("AMIR: moved the function: +%v......",taskinstance.TaskName)
		}
	}
	return false
}

func (master *Master)GetAnotherWorker(oldWorkerId string) string {

		for _, worker := range master.workers {
			if worker.WID != oldWorkerId {
				INFO.Println("changed from worker: %+v to new worker: %+v",oldWorkerId,worker.WID)
				return worker.WID
			}

		}

		INFO.Println("Could not find another worker")
		return oldWorkerId
}

func (master *Master) UpdatePrometheusConfig(){

	//TODO: Make this code more robust. If the request fails either prometheus is not ready or not running.
	// We should retry several times and send an proper error message if prometheus is down.
	baseConfig := `[{
	    "targets": ["127.0.0.1:8092"],
	    "labels": {
	      "job": "fogflow"
	    }
	}]`
	var configFile []PrometheusConfig

	if err := json.Unmarshal([]byte(baseConfig), &configFile); err != nil {
		panic(err)
	}
	for _, worker := range master.workers {
		configFile[0].Targets =
			append(configFile[0].Targets, worker.EdgeAddress+ ":" + strconv.Itoa(worker.CAdvisorPort))
	}


	//send update to prometheus
	body, err := json.Marshal(configFile)
	if err != nil {
		panic(err)
	}
	url:=master.cfg.Prometheus.Address+":" + strconv.Itoa(master.cfg.Prometheus.AdminPort)
	DEBUG.Printf("sending new config to premetheus[%v]:\n%v", url, string(body))

	req, err := http.NewRequest("POST", "http://"+url+"/config", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


}

type PromReply struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
			} `json:"metric"`
			Value []string `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func (master *Master)InitMetricWeights() {
	preferenceMatrix := mat.NewDense(master.METRIC_COUNT, master.METRIC_COUNT, []float64{
		1,0.5, 10, 5,
		2, 1 , 2 , 7,
		1,0.5, 1 ,0.3,
		2,  4, 5 , 3,
	})
	DEBUG.Printf("Matrix of Preference of Metrics is =\n %v\n\n", mat.Formatted(preferenceMatrix, mat.Prefix("    ")))

	var eig mat.Eigen
	ok := eig.Factorize(preferenceMatrix, mat.EigenRight)
	if !ok {
		log.Fatal("Eigendecomposition failed")
	}

	eigenValues := eig.Values(nil)
	var maxIndex,maxValue float64
	maxIndex = 0
	maxValue = -1

	//Find out the biggest Real EigenValue
	for i := 0; i < len(eigenValues); i++ {
		if real(eigenValues[i])>maxValue && imag(eigenValues[i])==0 {
			maxIndex=float64(i)
			maxValue=real(eigenValues[i])
		}
	}

	eigenVectors := eig.VectorsTo(nil)
	r,_:=eigenVectors.Dims()
	master.metricWeights = make([]float64, r)

	for i := 0; i < r; i++ {
		master.metricWeights[i]=real(eigenVectors.At(i,int(maxIndex)))
	}

	DEBUG.Printf("Decision Making Metric Weights : %v\n", master.metricWeights)

	// It would be called in the Timer:
	//master.UpdateUtilization()

}



func (master *Master)UpdateUtilization () {
	//DEBUG.Printf("Utilization of Edge nodes: \n %v \n",master.edgeUtilization)
	MetricQueries := make([]string,master.METRIC_COUNT)
	MetricQueries[0] = `sum(rate(container_cpu_usage_seconds_total{job="fogflow",id="/"}[2m] )) by (instance)`
	MetricQueries[1] = `sum(container_memory_working_set_bytes{job="fogflow",id=~"/docker.*"})by(instance)`
	MetricQueries[2] = `sum(container_memory_working_set_bytes{job="fogflow",id=~"/docker.*"}) by(instance) / sum(machine_memory_bytes) by (instance)`
	MetricQueries[3] = `sum(rate(container_memory_failures_total{failure_type="pgmajfault"}[20m])) by (instance)`

	//for address, metric := range *util {
	//	//AAA fmt.Printf("address %v is: %v \n", address, metric)
	//	//fmt.Printf("the value of metrics are: %v",util[address][0])
	//	//return worker.WID
	//}

	for metricNumber, metricQuery := range MetricQueries {

		promAddress := master.cfg.Prometheus.Address+":"+ strconv.Itoa(master.cfg.Prometheus.DataPort)
		url := "http://" + promAddress + "/api/v1/query?query=" + url2.QueryEscape(metricQuery)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			DEBUG.Printf("Error Talking to Prometheus Server for Orchestration")
			panic(err)
		}
		defer resp.Body.Close()
		text, _ := ioutil.ReadAll(resp.Body)
		var data PromReply
		json.Unmarshal(text, &data)

		for _, d := range data.Data.Result {
			f, _ := strconv.ParseFloat(d.Value[1], 64)
			if _,exists := master.edgeUtilization[d.Metric.Instance]; !exists {
				master.edgeUtilLock.Lock()
				master.edgeUtilization[d.Metric.Instance] =make([]float64,master.METRIC_COUNT,master.METRIC_COUNT)
				master.edgeUtilLock.Unlock()
			}
			master.edgeUtilLock.Lock()
			master.edgeUtilization[d.Metric.Instance][metricNumber] = f
			master.edgeUtilLock.Unlock()
		}
	}

	//Normalize the Utlization data
	for metricNumber, _ := range MetricQueries {
		//Normalize data of Each Metric (metricNumber)
		//for all hosts

		// DEBUG.Printf("%v: This is host: %v, and this is data: %v\n",metricNumber,host,data[metricNumber])

		//Get the sum
		var sumMetric float64 = 0
		for _,data :=range master.edgeUtilization{
			sumMetric += data[metricNumber]
		}

		//Divide each value to sum --> Normalize
		for _,data :=range master.edgeUtilization{
			master.edgeUtilLock.Lock()
			data[metricNumber] = divide(data[metricNumber],sumMetric)
			master.edgeUtilLock.Unlock()
		}
	}
}

func divide(a,b float64) float64{
	if math.IsNaN(a/b){
		return 0
	} else {
		return a/b

	}
}