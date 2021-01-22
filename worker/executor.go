package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/sethgrid/pester"

	docker "github.com/fsouza/go-dockerclient"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type taskContext struct {
	ListeningPort      string
	EndPointServiceIDs []EntityId
	Subscriptions      []string
	EntityID2SubID     map[string]string
	OutputStreams      []EntityId
	ContainerID        string
}

type pullResult struct {
	imageRef string
	err      error
}

type Executor struct {
	client    *docker.Client
	workerCfg *Config
	brokerURL string

	taskInstances map[string]*taskContext
	taskMap_lock  sync.RWMutex
}

func (e *Executor) Init(cfg *Config, selectedBrokerURL string) bool {
	// for Windows
	if runtime.GOOS == "windows" {
		endpoint := os.Getenv("DOCKER_HOST")
		path := os.Getenv("DOCKER_CERT_PATH")
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)
		client, err := docker.NewTLSClient(endpoint, cert, key, ca)

		if err != nil || client == nil {
			INFO.Println("Couldn't connect to docker: %v", err)
			return false
		}

		e.client = client
	} else {
		// for Linux
		endpoint := "unix:///var/run/docker.sock"
		client, err := docker.NewClient(endpoint)

		if err != nil || client == nil {
			INFO.Println("Couldn't connect to docker: %v", err)
			return false
		}

		e.client = client
	}

	e.workerCfg = cfg
	e.brokerURL = selectedBrokerURL

	e.taskInstances = make(map[string]*taskContext)
	return true
}

func (e *Executor) Shutdown() {
	e.terminateAllTasks()
}

func (e *Executor) GetNumOfTasks() int {
	e.taskMap_lock.RLock()
	defer e.taskMap_lock.RUnlock()

	return len(e.taskInstances)
}

func (e *Executor) InspectImage(dockerImage string) bool {
	_, err := e.client.InspectImage(dockerImage)
	if err != nil {
		INFO.Printf("operator image %s does not exist locally\r\n", dockerImage)
		return false
	} else {
		INFO.Printf("operator image %s exists locally\r\n", dockerImage)
		return true
	}
}

func (e *Executor) PullImage(dockerImage string, tag string) (string, error) {
	auth := docker.AuthConfiguration{}

	if e.workerCfg.Worker.Registry.IsConfigured() == true {
		auth.Username = e.workerCfg.Worker.Registry.Username
		auth.Password = e.workerCfg.Worker.Registry.Password
		auth.Email = e.workerCfg.Worker.Registry.Email
		auth.ServerAddress = e.workerCfg.Worker.Registry.ServerAddress
		dockerImage = dockerImage
	}

	DEBUG.Printf("options : %+v\r\n", auth)

	opts := docker.PullImageOptions{
		Repository: dockerImage,
		Tag:        tag,
	}

	DEBUG.Printf("options : %+v\r\n", opts)

	err := e.client.PullImage(opts, auth)
	if err != nil {
		ERROR.Printf("failed to pull this image %s, error %v\r\n", dockerImage, err)
		return "", err
	}

	// check if the image exists now
	resp, err := e.client.InspectImage(dockerImage)
	if err != nil {
		ERROR.Printf("the image %s does not exist, even throug it has been pulled, error %v\r\n", dockerImage, err)
		return "", err
	}

	if resp == nil {
		return "", nil
	}

	imageRef := resp.ID
	if len(resp.RepoDigests) > 0 {
		imageRef = resp.RepoDigests[0]
	}

	INFO.Println("fetched image ", dockerImage)
	return imageRef, nil
}

func (e *Executor) ListContainers() {
	containers, _ := e.client.ListContainers(docker.ListContainersOptions{All: true})
	for _, container := range containers {
		DEBUG.Println("Name: ", container.Names)
	}
}

func (e *Executor) startContainerWithBridge(dockerImage string, portNum string) (string, error) {
	// prepare the configuration for a docker container
	config := docker.Config{Image: dockerImage}
	portBindings := map[docker.Port][]docker.PortBinding{
		"8080/tcp": []docker.PortBinding{docker.PortBinding{HostIP: "0.0.0.0", HostPort: portNum}}}
	hostConfig := docker.HostConfig{PortBindings: portBindings}
	containerOptions := docker.CreateContainerOptions{Config: &config,
		HostConfig: &hostConfig}

	// create a new docker container
	container, err := e.client.CreateContainer(containerOptions)
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	// start the new container
	err = e.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	return container.ID, nil
}

func (e *Executor) writeTempFile(fileName string, fileContent string) {
	content := []byte(fileContent)
	tmpfile, err := os.Create(fileName)
	if err != nil {
		ERROR.Println(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		ERROR.Println(err)
	}
	if err := tmpfile.Close(); err != nil {
		ERROR.Println(err)
	}
}

func (e *Executor) startContainer(dockerImage string, portNum string, functionCode string, taskID string, adminCfg []interface{}, servicePorts []string) (string, error) {
	// prepare the configuration for a docker container, host mode for the container network
	evs := make([]string, 0)
	evs = append(evs, fmt.Sprintf("myport=%s", portNum))

	// pass the initial configuration as the environmental variable
	jsonString, _ := json.Marshal(adminCfg)
	evs = append(evs, fmt.Sprintf("adminCfg=%s", jsonString))

	config := docker.Config{Image: dockerImage, Env: evs}

	hostConfig := docker.HostConfig{}

	//if runtime.GOOS == "darwin" {   already use the bridge model
	internalPort := docker.Port(portNum + "/tcp")
	portBindings := map[docker.Port][]docker.PortBinding{
		internalPort: []docker.PortBinding{docker.PortBinding{HostIP: "0.0.0.0", HostPort: portNum}}}

	config.ExposedPorts = map[docker.Port]struct{}{internalPort: {}}

	// add the other listening ports into the exposed port list
	for _, port := range servicePorts {
		internalPort := docker.Port(port + "/tcp")
		portBindings[internalPort] = []docker.PortBinding{docker.PortBinding{HostIP: "0.0.0.0", HostPort: port}}
	}

	hostConfig.PortBindings = portBindings
	//} else {
	//	hostConfig.NetworkMode = "host"
	//}

	// to configure if the container will be removed once it is terminated
	hostConfig.AutoRemove = e.workerCfg.Worker.ContainerAutoRemove

	if functionCode != "" {
		fileName := "/tmp/" + taskID
		e.writeTempFile(fileName, functionCode)

		mount := docker.HostMount{}
		mount.Source = fileName
		mount.Target = "/app/function.js"
		mount.ReadOnly = true
		mount.Type = "bind"

		DEBUG.Println("mounting configuration ", mount)

		hostConfig.Mounts = make([]docker.HostMount, 0)
		hostConfig.Mounts = append(hostConfig.Mounts, mount)
	}

	containerOptions := docker.CreateContainerOptions{Config: &config,
		HostConfig: &hostConfig}

	// create a new docker container
	container, err := e.client.CreateContainer(containerOptions)
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	// start the new container
	err = e.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	return container.ID, nil
}

// Ask the kernel for a free open port that is ready to use
func (e *Executor) findFreePortNumber() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func (e *Executor) LaunchTask(task *ScheduledTaskInstance) bool {
	dockerImage := task.DockerImage
	fmt.Println("=========This is task============")
	fmt.Println(task)
	INFO.Println("to execute Task ", task.ID, " to perform Operation ", dockerImage)
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

	// first check the image locally
	if e.InspectImage(dockerImage) == false {
		// if the image does not exist locally, try to fetch it from docker hub
		_, pullError := e.PullImage(dockerImage, "latest")
		if pullError != nil {
			ERROR.Printf("failed to fetch the image %s\r\n", task.DockerImage)
			return false
		}
	}

	taskCtx := taskContext{}
	taskCtx.EntityID2SubID = make(map[string]string)
	taskCtx.EndPointServiceIDs = make([]EntityId, 0)

	// find a free listening port number available on the host machine
	freePort := strconv.Itoa(e.findFreePortNumber())

	// function code
	functionCode := task.FunctionCode

	// configure the task with its output streams via its admin interface
	commands := make([]interface{}, 0)

	// set broker URL
	setBrokerCmd := make(map[string]interface{})
	setBrokerCmd["command"] = "CONNECT_BROKER"
	setBrokerCmd["brokerURL"] = e.brokerURL
	commands = append(commands, setBrokerCmd)

	// pass the reference URL to the task so that the task can issue context subscription as well
	setReferenceCmd := make(map[string]interface{})
	setReferenceCmd["command"] = "SET_REFERENCE"
	setReferenceCmd["url"] = "http://" + e.workerCfg.InternalIP + ":" + freePort
	commands = append(commands, setReferenceCmd)

	// set output stream
	for _, outStream := range task.Outputs {
		setOutputCmd := make(map[string]interface{})
		setOutputCmd["command"] = "SET_OUTPUTS"
		setOutputCmd["type"] = outStream.Type
		setOutputCmd["id"] = outStream.StreamID
		commands = append(commands, setOutputCmd)

		// record its outputs
		var eid EntityId
		eid.ID = outStream.StreamID
		eid.Type = outStream.Type
		eid.IsPattern = false
		taskCtx.OutputStreams = append(taskCtx.OutputStreams, eid)
	}

	// check if it is required to set up the portmapping for its endpoint services
	servicePorts := make([]string, 0)

	for _, parameter := range task.Parameters {
		// deal with the service port
		if parameter.Name == "service_port" {
			servicePorts = append(servicePorts, parameter.Values...)
		}
	}

	// start a container to run the scheduled task instance
	containerId, err := e.startContainer(dockerImage, freePort, functionCode, task.ID, commands, servicePorts)
	if err != nil {
		ERROR.Println(err)
		return false
	}
	INFO.Printf(" task %s  started within container = %s\n", task.ID, containerId)

	taskCtx.ListeningPort = freePort
	taskCtx.ContainerID = containerId

	// register the service ports of uservices
	if len(servicePorts) > 0 {
		// currently, we assume that each task will only provide one end-point service
		eid := e.registerEndPointService(task.ServiceName, task.ID, task.OperatorName, e.workerCfg.ExternalIP, servicePorts[0], e.workerCfg.Location)
		taskCtx.EndPointServiceIDs = append(taskCtx.EndPointServiceIDs, eid)
	}

	INFO.Printf("subscribe its input streams")

	// subscribe input streams on behalf of the launched task
	taskCtx.Subscriptions = make([]string, 0)

	for _, inputStream := range task.Inputs {
		NGSILD := e.queryForNGSILdEntity(inputStream.ID)
		if NGSILD == 200 {
			fmt.Println(&inputStream)
			subID, err := e.subscribeLdInputStream(freePort, &inputStream)
			if err == nil {
				DEBUG.Println("===========subID = ", subID)
				taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
				taskCtx.EntityID2SubID[inputStream.ID] = subID
			} else {
				ERROR.Println(err)
			}
		}
		NGSIV1 := e.queryForNGSIV1Entity(inputStream.ID)
		if NGSIV1 == 200 {
			subID, err := e.subscribeInputStream(freePort, &inputStream)
			if err == nil {
				DEBUG.Println("===========subID = ", subID)
				taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
				taskCtx.EntityID2SubID[inputStream.ID] = subID
			} else {
				ERROR.Println(err)
			}
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

//Query for NGSILD Entity with entityId
func (e *Executor) queryForNGSILdEntity(eid string) int {
	if eid == "" {
		return 404
	}
	brokerURL := e.brokerURL
	brokerURL = strings.TrimSuffix(brokerURL, "/ngsi10")
	client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	statusCode := client.QueryForNGSILDEntity(eid)
	fmt.Println(statusCode)
	return statusCode
}

// Query for NGSIV1 Entity with EntityId
func (e *Executor) queryForNGSIV1Entity(eid string) int {
	if eid == "" {
		return 200
	}
	fmt.Println(e.brokerURL)
	client := NGSI10Client{IoTBrokerURL: e.brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	statusCode := client.QueryForNGSIV1Entity(eid)
	fmt.Println(statusCode)
	return statusCode
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
	ctxObj.Attributes["service"] = ValueObject{Type: "string", Value: task.ServiceName}

	ctxObj.Metadata = make(map[string]ValueObject)
	ctxObj.Metadata["topology"] = ValueObject{Type: "string", Value: task.ServiceName}
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

	if len(inputStream.ID) > 0 { // for a specific context entity
		newEntity.Type = inputStream.Type
		newEntity.ID = inputStream.ID
	} else { // for all context entities with a specific type
		newEntity.Type = inputStream.Type
	}

	LdSubscription.Entities = make([]EntityId, 0)
	LdSubscription.Entities = append(LdSubscription.Entities, newEntity)
	LdSubscription.Type = "Subscription"
	LdSubscription.WatchedAttributes = inputStream.AttributeList

	LdSubscription.Notification.Endpoint.URI = "http://" + e.workerCfg.InternalIP + ":" + agentPort + "/notifyContext"

	DEBUG.Printf(" =========== issue the following subscription =========== %+v\r\n", LdSubscription)
	brokerURL := e.brokerURL
	brokerURL = strings.TrimSuffix(brokerURL, "/ngsi10")
	client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: &e.workerCfg.HTTPS}
	sid, err := client.SubscribeLdContext(&LdSubscription, true)
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

func (e *Executor) PauseTask(taskID string) {
	e.taskMap_lock.RLock()
	defer e.taskMap_lock.RUnlock()

	containerID := e.taskInstances[taskID].ContainerID
	err := e.client.PauseContainer(containerID)
	if err != nil {
		ERROR.Println(err)
	}
}

func (e *Executor) ResumeTask(taskID string) {
	e.taskMap_lock.RLock()
	defer e.taskMap_lock.RUnlock()

	containerID := e.taskInstances[taskID].ContainerID
	err := e.client.UnpauseContainer(containerID)
	if err != nil {
		ERROR.Println(err)
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
	go e.client.StopContainer(containerID, 1)

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

func (e *Executor) onAddInput(flow *FlowInfo) {
	e.taskMap_lock.Lock()
	defer e.taskMap_lock.Unlock()

	taskCtx := e.taskInstances[flow.TaskInstanceID]

	if taskCtx == nil {
		return
	}

	if e.workerCfg.Worker.StartActualTask == false {
		return
	}
	Id := flow.InputStream.ID
	NGSILD := e.queryForNGSILdEntity(Id)
	if NGSILD == 200 {
		subID, err := e.subscribeLdInputStream(taskCtx.ListeningPort, &flow.InputStream)
		if err == nil {
			DEBUG.Println("===========subscribe new input = ", flow, " , subID = ", subID)
			taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
			taskCtx.EntityID2SubID[flow.InputStream.ID] = subID
		} else {
			ERROR.Println(err)
		}
	}
	NGSIV1 := e.queryForNGSIV1Entity(Id)
	if NGSIV1 == 200 {
		subID, err := e.subscribeInputStream(taskCtx.ListeningPort, &flow.InputStream)
		if err == nil {
			DEBUG.Println("===========subscribe new input = ", flow, " , subID = ", subID)
			taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
			taskCtx.EntityID2SubID[flow.InputStream.ID] = subID
		} else {
			ERROR.Println(err)
		}
	}
}

func (e *Executor) onRemoveInput(flow *FlowInfo) {
	e.taskMap_lock.Lock()
	defer e.taskMap_lock.Unlock()

	if e.workerCfg.Worker.StartActualTask == false {
		return
	}

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
