package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/sethgrid/pester"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type taskContext struct {
	ListeningPort      string
	EndPointServiceIDs []EntityId
	Subscriptions      []string
	EntityID2SubID     map[string]string
	OutputStreams      []EntityId
	ContainerID        string
}

type Executor struct {
	client Engine

	workerCfg *Config
	brokerURL string

	taskInstances map[string]*taskContext
	taskMap_lock  sync.RWMutex
}

func (e *Executor) Init(cfg *Config, selectedBrokerURL string) bool {
	e.workerCfg = cfg
	e.brokerURL = selectedBrokerURL
	e.taskInstances = make(map[string]*taskContext)

	if strings.EqualFold(cfg.Worker.ContainerManagement, "docker") {
		e.client = &DockerEngine{}
	} else if strings.EqualFold(cfg.Worker.ContainerManagement, "kubernetes") {
		e.client = &Kubernetes{}
	} else if strings.EqualFold(cfg.Worker.ContainerManagement, "mec") {
		e.client = &EdgeController{}
	} else {
		e.client = &DockerEngine{}
	}

	return e.client.Init(cfg)
}

func (e *Executor) Shutdown() {
	e.terminateAllTasks()
}

func (e *Executor) GetNumOfTasks() int {
	e.taskMap_lock.RLock()
	defer e.taskMap_lock.RUnlock()

	return len(e.taskInstances)
}

func (e *Executor) PullImage(dockerImage string, tag string) (string, error) {
	return e.client.PullImage(dockerImage, tag)
}

func (e *Executor) LaunchTask(task *ScheduledTaskInstance) bool {
	if e.workerCfg.Worker.StartActualTask == false {
		// just for the performance evaluation of Topology Master
		taskCtx := taskContext{}

		e.taskMap_lock.Lock()
		e.taskInstances[task.ID] = &taskCtx
		e.taskMap_lock.Unlock()

		INFO.Printf("register this task")

		// register this new task entity to IoT Broker
		e.registerTask(task, "000", "000")

		return true
	}

	taskCtx := taskContext{}
	taskCtx.EntityID2SubID = make(map[string]string)
	taskCtx.EndPointServiceIDs = make([]EntityId, 0)

	// set output stream
	for _, outStream := range task.Outputs {
		// record its outputs
		var eid EntityId
		eid.ID = outStream.StreamID
		eid.Type = outStream.Type
		eid.IsPattern = false
		taskCtx.OutputStreams = append(taskCtx.OutputStreams, eid)
	}

	// start a container to run the scheduled task instance
	containerId, freePort, err := e.client.StartTask(task, e.brokerURL)
	if err != nil {
		ERROR.Println(err)
		return false
	}
	INFO.Printf(" task %s  started within container = %s\n", task.ID, containerId)

	taskCtx.ListeningPort = freePort
	taskCtx.ContainerID = containerId

	// register the service ports of uservices
	// check if it is required to set up the portmapping for its endpoint services
	servicePorts := make([]string, 0)

	for _, parameter := range task.Parameters {
		// deal with the service port
		if parameter.Name == "service_port" {
			servicePorts = append(servicePorts, parameter.Values...)
		}
	}

	if len(servicePorts) > 0 {
		// currently, we assume that each task will only provide one end-point service
		eid := e.registerEndPointService(task.TopologyName, task.ID, task.OperatorName, e.workerCfg.ExternalIP, servicePorts[0], e.workerCfg.Location)
		taskCtx.EndPointServiceIDs = append(taskCtx.EndPointServiceIDs, eid)
	}

	INFO.Printf("subscribe its input streams")

	// subscribe input streams on behalf of the launched task
	taskCtx.Subscriptions = make([]string, 0)

	for _, inputStream := range task.Inputs {
		if inputStream.MsgFormat == "NGSIV1" {
			DEBUG.Println("========Subscription for NGSIV1 task==========")
			subID, err := e.subscribeInputStream(freePort, &inputStream)
			if err == nil {
				DEBUG.Println("===========subID = ", subID)
				taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
				taskCtx.EntityID2SubID[inputStream.ID] = subID
			} else {
				ERROR.Println(err)
			}
		} else if inputStream.MsgFormat == "NGSILD" {
			subID, err := e.subscribeLdInputStream(freePort, &inputStream)
			if err == nil {
				DEBUG.Println("===========subID = ", subID)
				taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
				taskCtx.EntityID2SubID[inputStream.ID] = subID
			} else {
				ERROR.Println(err)
			}
		} else {
			ERROR.Println("unsupported message protocol")
		}
	}

	// update the task list
	e.taskMap_lock.Lock()
	e.taskInstances[task.ID] = &taskCtx
	e.taskMap_lock.Unlock()

	INFO.Printf("register this task")

	// register this new task entity to IoT Broker
	e.registerTask(task, freePort, containerId)

	return true
}

func (e *Executor) registerEndPointService(serviceName string, taskID string, operateName string, ipAddr string, port string, location PhysicalLocation) EntityId {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = "uService." + serviceName + "." + taskID
	ctxObj.Entity.Type = "uService"
	ctxObj.Entity.IsPattern = false

	ctxObj.Metadata = make(map[string]ValueObject)
	ctxObj.Metadata["service"] = ValueObject{Type: "string", Value: serviceName}
	ctxObj.Metadata["taskID"] = ValueObject{Type: "string", Value: taskID}
	ctxObj.Metadata["operator"] = ValueObject{Type: "string", Value: operateName}
	ctxObj.Metadata["IP"] = ValueObject{Type: "string", Value: ipAddr}
	ctxObj.Metadata["port"] = ValueObject{Type: "object", Value: port}
	ctxObj.Metadata["location"] = ValueObject{Type: "string", Value: location}

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}

	return ctxObj.Entity
}

func (e *Executor) deRegisterEndPointService(eid EntityId) {
	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.DeleteContext(&eid)
	if err != nil {
		ERROR.Println(err)
	}
}

func (e *Executor) configurateTask(port string, commands []interface{}) bool {
	taskAdminURL := fmt.Sprintf("http://%s:%s/admin", e.workerCfg.InternalIP, port)

	jsonText, _ := json.Marshal(commands)

	INFO.Println(taskAdminURL)
	INFO.Printf("configuration: %s\r\n", string(jsonText))

	req, _ := http.NewRequest("POST", taskAdminURL, bytes.NewBuffer(jsonText))
	req.Header.Set("Content-Type", "application/json")

	client := pester.New()
	client.MaxRetries = 30
	client.Backoff = pester.LinearBackoff

	resp, err := client.Do(req)
	if err != nil {
		ERROR.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	INFO.Println("task on port ", port, " has been configured with parameters ", jsonText)
	INFO.Println("response Body:", string(body))

	return true
}

func (e *Executor) registerTask(task *ScheduledTaskInstance, portNum string, containerID string) {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = task.ID
	ctxObj.Entity.Type = "Task"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["id"] = ValueObject{Type: "string", Value: task.ID}
	ctxObj.Attributes["port"] = ValueObject{Type: "string", Value: portNum}
	ctxObj.Attributes["status"] = ValueObject{Type: "string", Value: task.Status}
	ctxObj.Attributes["worker"] = ValueObject{Type: "string", Value: task.WorkerID}

	ctxObj.Attributes["task"] = ValueObject{Type: "string", Value: task.TaskName}
	ctxObj.Attributes["operator"] = ValueObject{Type: "string", Value: task.OperatorName}
	ctxObj.Attributes["service"] = ValueObject{Type: "string", Value: task.TopologyName}

	ctxObj.Metadata = make(map[string]ValueObject)
	ctxObj.Metadata["topology"] = ValueObject{Type: "string", Value: task.TopologyName}
	ctxObj.Metadata["worker"] = ValueObject{Type: "string", Value: task.WorkerID}

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}
}

func (e *Executor) updateTask(taskID string, status string) {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = taskID
	ctxObj.Entity.Type = "Task"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["status"] = ValueObject{Type: "string", Value: status}

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}
}

func (e *Executor) deregisterTask(taskID string) {
	entity := EntityId{}
	entity.ID = taskID
	entity.Type = "Task"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.DeleteContext(&entity)
	if err != nil {
		ERROR.Println(err)
	}
}

// Subscribe for NGSILD input stream
func (e *Executor) subscribeLdInputStream(agentPort string, inputStream *InputStream) (string, error) {
	LdSubscription := LDSubscriptionRequest{}

	newEntity := EntityId{}
	var Fs, Fsp, ID string
	if len(inputStream.ID) > 0 { // for a specific context entity
		newEntity.Type = inputStream.Type
		ID, Fs = FiwareId(inputStream.ID)
		if Fs == "default" {
			Fs = ""
		}
		newEntity.ID = ID
		Fsp = inputStream.FiwareServicePath
	} else {
		newEntity.Type = inputStream.Type
	}
	fmt.Println("FS,FSP", Fs, Fsp)
	LdSubscription.Entities = make([]EntityId, 0)
	LdSubscription.Entities = append(LdSubscription.Entities, newEntity)
	LdSubscription.Type = "Subscription"
	LdSubscription.WatchedAttributes = inputStream.AttributeList

	LdSubscription.Notification.Endpoint.URI = "http://" + e.workerCfg.InternalIP + ":" + agentPort + "/notifyContext"

	DEBUG.Printf(" =========== issue the following subscription =========== %+v\r\n", LdSubscription)
	brokerURL := e.brokerURL
	brokerURL = strings.TrimSuffix(brokerURL, "/ngsi10")
	client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	sid, err := client.SubscribeLdContext(&LdSubscription, true, Fs, Fsp)
	fmt.Println("sid", sid)
	if err != nil {
		ERROR.Println(err)
		return "", err
	} else {
		return sid, nil
	}
}

//Subscribe for NGSIV1 input stream
func (e *Executor) subscribeInputStream(agentPort string, inputStream *InputStream) (string, error) {
	fmt.Println("====================Subscription here ===================")
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}

	if len(inputStream.ID) > 0 { // for a specific context entity
		newEntity.IsPattern = false
		newEntity.Type = inputStream.Type
		newEntity.ID = inputStream.ID
	} else { // for all context entities with a specific type
		newEntity.Type = inputStream.Type
		newEntity.IsPattern = true
	}

	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)

	subscription.Attributes = inputStream.AttributeList

	subscription.Reference = "http://" + e.workerCfg.InternalIP + ":" + agentPort

	DEBUG.Printf(" =========== issue the following subscription =========== %+v\r\n", subscription)

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		ERROR.Println(err)
		return "", err
	} else {
		return sid, nil
	}
}

func (e *Executor) unsubscribeInputStream(sid string) error {
	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.UnsubscribeContext(sid)
	if err != nil {
		ERROR.Println(err)
		return err
	} else {
		return nil
	}
}

func (e *Executor) createOuputStream(eID string, eType string) error {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = eID
	ctxObj.Entity.Type = eType
	ctxObj.Entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
		return err
	} else {
		return nil
	}
}

func (e *Executor) deleteOuputStream(eid *EntityId) error {
	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	err := client.DeleteContext(eid)
	if err != nil {
		ERROR.Println(err)
		return err
	} else {
		return nil
	}
}

func (e *Executor) TerminateTask(taskID string, paused bool) {
	INFO.Println("================== terminate task ID ============ ", taskID)

	if e.workerCfg.Worker.StartActualTask == false {
		// just for the performance evaluation of Topology Master
		e.taskMap_lock.Lock()

		if _, ok := e.taskInstances[taskID]; ok == true {
			delete(e.taskInstances, taskID)
		}

		e.taskMap_lock.Unlock()

		INFO.Printf("deregister this task")
		go e.deregisterTask(taskID)

		return
	}

	e.taskMap_lock.Lock()
	if _, ok := e.taskInstances[taskID]; ok == false {
		e.taskMap_lock.Unlock()
		return
	}

	containerID := e.taskInstances[taskID].ContainerID
	e.taskMap_lock.Unlock()

	//stop the container first
	e.client.StopTask(containerID)
	INFO.Printf(" task %s  terminate from the container = %s\n", taskID, containerID)

	e.taskMap_lock.Lock()

	// issue unsubscribe
	for _, subID := range e.taskInstances[taskID].Subscriptions {
		INFO.Println("issued subscription: ", subID)
		err := e.unsubscribeInputStream(subID)
		if err != nil {
			ERROR.Println(err)
		}
		INFO.Printf(" subscriptions (%s) have been canceled\n", subID)
	}

	// delete the output streams of the terminated task
	for _, outStream := range e.taskInstances[taskID].OutputStreams {
		e.deleteOuputStream(&outStream)
	}

	// deregister the associated end point service
	for _, serviceEntityID := range e.taskInstances[taskID].EndPointServiceIDs {
		go e.deRegisterEndPointService(serviceEntityID)
	}

	delete(e.taskInstances, taskID)

	e.taskMap_lock.Unlock()

	if paused == true {
		// only update its status
		go e.updateTask(taskID, "paused")
	} else {
		// deregister this task entity
		go e.deregisterTask(taskID)
	}
}

// stop all running tasks
func (e *Executor) terminateAllTasks() {
	idList := make([]string, 0)
	e.taskMap_lock.RLock()
	for id, _ := range e.taskInstances {
		idList = append(idList, id)
	}
	e.taskMap_lock.RUnlock()

	var wg sync.WaitGroup
	wg.Add(len(idList))

	for _, taskID := range idList {
		go func(tID string) {
			defer wg.Done()
			e.TerminateTask(tID, false)
		}(taskID)
	}

	wg.Wait()
}

// add the specified input for an existing task
func (e *Executor) onAddInput(flow *FlowInfo) {
	if e.workerCfg.Worker.StartActualTask == false {
		return
	}

	e.taskMap_lock.Lock()
	defer e.taskMap_lock.Unlock()

	taskCtx := e.taskInstances[flow.TaskInstanceID]
	if taskCtx == nil {
		ERROR.Println("the requested task does not exist")
		return
	}

	if flow.InputStream.MsgFormat == "NGSIV1" {
		subID, err := e.subscribeInputStream(taskCtx.ListeningPort, &flow.InputStream)
		if err == nil {
			DEBUG.Println("===========subscribe new input = ", flow, " , subID = ", subID)
			taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
			taskCtx.EntityID2SubID[flow.InputStream.ID] = subID
		} else {
			ERROR.Println(err)
		}
	} else if flow.InputStream.MsgFormat == "NGSILD" {
		subID, err := e.subscribeLdInputStream(taskCtx.ListeningPort, &flow.InputStream)
		if err == nil {
			DEBUG.Println("===========subscribe new input = ", flow, " , subID = ", subID)
			taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
			taskCtx.EntityID2SubID[flow.InputStream.ID] = subID
		} else {
			ERROR.Println(err)
		}
	} else {
		ERROR.Println("not supported")
	}
}

// remove the specified input for an existing task
func (e *Executor) onRemoveInput(flow *FlowInfo) {
	if e.workerCfg.Worker.StartActualTask == false {
		return
	}

	e.taskMap_lock.Lock()
	defer e.taskMap_lock.Unlock()

	taskCtx := e.taskInstances[flow.TaskInstanceID]
	subID := taskCtx.EntityID2SubID[flow.InputStream.ID]

	err := e.unsubscribeInputStream(subID)
	if err != nil {
		ERROR.Println(err)
	}

	for i, sid := range taskCtx.Subscriptions {
		if sid == subID {
			taskCtx.Subscriptions = append(taskCtx.Subscriptions[:i], taskCtx.Subscriptions[i+1:]...)
			break
		}
	}

	delete(taskCtx.EntityID2SubID, flow.InputStream.ID)
}
