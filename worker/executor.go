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
	"sync"

	"github.com/sethgrid/pester"

	docker "github.com/fsouza/go-dockerclient"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type taskContext struct {
	ListeningPort  string
	Subscriptions  []string
	EntityID2SubID map[string]string
	OutputStreams  []EntityId
	ContainerID    string
}

type pullResult struct {
	imageRef string
	err      error
}

type Executor struct {
	client        *docker.Client
	workerCfg     *Config
	brokerURL     string
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

func (e *Executor) ListImages() {
	imgs, _ := e.client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

func (e *Executor) InspectImage(dockerImage string) {
	img, err := e.client.InspectImage(dockerImage)
	if err != nil {
		fmt.Printf("failed to access this image %+v", err)
		return
	}

	fmt.Printf("cfg : %+v", img.Config)
}

func (e *Executor) PullImage(dockerImage string, tag string) (string, error) {
	auth := docker.AuthConfiguration{}

	if e.workerCfg.Worker.Registry.IsConfigured() == true {
		auth.Username = e.workerCfg.Worker.Registry.Username
		auth.Password = e.workerCfg.Worker.Registry.Password
		auth.Email = e.workerCfg.Worker.Registry.Email
		auth.ServerAddress = e.workerCfg.Worker.Registry.ServerAddress
		dockerImage = e.workerCfg.Worker.Registry.ServerAddress + "/" + dockerImage
	}

	fmt.Printf("options : %+v\r\n", auth)

	opts := docker.PullImageOptions{
		Repository: dockerImage,
		Tag:        tag,
	}

	fmt.Printf("options : %+v\r\n", opts)

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
		fmt.Println("Name: ", container.Names)
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

func (e *Executor) startContainer(dockerImage string, portNum string, functionCode string, taskID string) (string, error) {
	// prepare the configuration for a docker container, host mode for the container network
	evs := make([]string, 0)
	evs = append(evs, fmt.Sprintf("myport=%s", portNum))

	config := docker.Config{Image: dockerImage, Env: evs}

	hostConfig := docker.HostConfig{}

	hostConfig.NetworkMode = "host"
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

	INFO.Println("to execute Task ", task.ID, " to perform Operation ", dockerImage)

	//	if e.workerCfg.Worker.Registry.IsConfigured() == true {
	// to fetch the docker image
	_, pullError := e.PullImage(dockerImage, "latest")
	if pullError != nil {
		ERROR.Printf("failed to fetch the image %s\r\n", task.DockerImage)
		return false
	}
	//	} else {
	//		// assume the requested image is available
	//		INFO.Println("no docker registery is configured, therefore we assume the requested docker image is already built locally and available")
	//	}

	taskCtx := taskContext{}
	taskCtx.EntityID2SubID = make(map[string]string)

	// find a free listening port number available on the host machine
	freePort := strconv.Itoa(e.findFreePortNumber())

	// function code
	functionCode := task.FunctionCode

	// start a container to run the scheduled task instance
	containerId, err := e.startContainer(dockerImage, freePort, functionCode, task.ID)
	if err != nil {
		ERROR.Println(err)
		return false
	}
	INFO.Printf(" task %s  started within container = %s\n", task.ID, containerId)

	taskCtx.ListeningPort = freePort
	taskCtx.ContainerID = containerId

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

	INFO.Printf("configure the task with %+v, via port %s\r\n", commands, freePort)

	if e.configurateTask(freePort, commands) == false {
		ERROR.Println("failed to configure the task instance")
		return false
	}

	INFO.Printf("subscribe its input streams")

	// subscribe input streams on behalf of the launched task
	taskCtx.Subscriptions = make([]string, 0)

	for _, streamType := range task.Inputs {
		for _, streamId := range streamType.Streams {
			subID, err := e.subscribeInputStream(freePort, streamType.Type, streamId)
			if err == nil {
				fmt.Println("===========subID = ", subID)
				taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
				taskCtx.EntityID2SubID[streamId] = subID
			} else {
				fmt.Println(err)
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
		panic(err)
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

	ctxObj.Entity.ID = "Task." + task.ID
	ctxObj.Entity.Type = "Task"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["id"] = ValueObject{Type: "string", Value: task.ID}
	ctxObj.Attributes["port"] = ValueObject{Type: "string", Value: portNum}
	ctxObj.Attributes["status"] = ValueObject{Type: "string", Value: task.Status}
	ctxObj.Attributes["worker"] = ValueObject{Type: "string", Value: task.WorkerID}

	ctxObj.Metadata = make(map[string]ValueObject)
	ctxObj.Metadata["topology"] = ValueObject{Type: "string", Value: task.ServiceName}
	ctxObj.Metadata["worker"] = ValueObject{Type: "string", Value: task.WorkerID}

	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	err := client.UpdateContext(&ctxObj)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *Executor) updateTask(taskID string, status string) {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = "Task." + taskID
	ctxObj.Entity.Type = "Task"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["status"] = ValueObject{Type: "string", Value: status}

	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	err := client.UpdateContext(&ctxObj)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *Executor) deregisterTask(taskID string) {
	entity := EntityId{}
	entity.ID = "Task." + taskID
	entity.Type = "Task"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	err := client.DeleteContext(&entity)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *Executor) subscribeInputStream(agentPort string, streamType string, streamId string) (string, error) {
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}

	if len(streamId) > 0 { // for a specific context entity
		newEntity.IsPattern = false
		newEntity.Type = streamType
		newEntity.ID = streamId
	} else { // for all context entities with a specific type
		newEntity.Type = streamType
		newEntity.IsPattern = true
	}

	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)

	subscription.Reference = "http://" + e.workerCfg.InternalIP + ":" + agentPort

	fmt.Printf(" =========== issue the following subscription =========== %+v\r\n", subscription)

	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		fmt.Println(err)
		return "", err
	} else {
		return sid, nil
	}
}

func (e *Executor) unsubscribeInputStream(sid string) error {
	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	err := client.UnsubscribeContext(sid)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (e *Executor) deleteOuputStream(eid *EntityId) error {
	client := NGSI10Client{IoTBrokerURL: e.brokerURL}
	err := client.DeleteContext(eid)
	if err != nil {
		fmt.Println(err)
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
	var wg sync.WaitGroup
	wg.Add(len(e.taskInstances))

	for taskID, _ := range e.taskInstances {
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

	subID, err := e.subscribeInputStream(taskCtx.ListeningPort, flow.EntityType, flow.EntityID)
	if err == nil {
		fmt.Println("===========subscribe new input = ", flow, " , subID = ", subID)
		taskCtx.Subscriptions = append(taskCtx.Subscriptions, subID)
		taskCtx.EntityID2SubID[flow.EntityID] = subID
	} else {
		ERROR.Println(err)
	}
}

func (e *Executor) onRemoveInput(flow *FlowInfo) {
	e.taskMap_lock.Lock()
	defer e.taskMap_lock.Unlock()

	taskCtx := e.taskInstances[flow.TaskInstanceID]
	subID := taskCtx.EntityID2SubID[flow.EntityID]

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

	delete(taskCtx.EntityID2SubID, flow.EntityID)
}
