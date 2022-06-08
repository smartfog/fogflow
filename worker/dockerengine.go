package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"

	docker "github.com/fsouza/go-dockerclient"
)

type DockerEngine struct {
	client *docker.Client

	workerCfg *Config
}

func (dockerengine *DockerEngine) Init(cfg *Config) bool {
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

		dockerengine.client = client
	} else {
		// for Linux
		endpoint := "unix:///var/run/docker.sock"
		client, err := docker.NewClient(endpoint)

		if err != nil || client == nil {
			INFO.Println("Couldn't connect to docker: %v", err)
			return false
		}

		dockerengine.client = client
	}

	dockerengine.workerCfg = cfg

	return true
}

func (dockerengine *DockerEngine) PullImage(dockerImage string, tag string) (string, error) {
	auth := docker.AuthConfiguration{}

	if dockerengine.workerCfg.Worker.Registry.IsConfigured() == true {
		auth.Username = dockerengine.workerCfg.Worker.Registry.Username
		auth.Password = dockerengine.workerCfg.Worker.Registry.Password
		auth.Email = dockerengine.workerCfg.Worker.Registry.Email
		auth.ServerAddress = dockerengine.workerCfg.Worker.Registry.ServerAddress
		dockerImage = dockerImage
	}

	DEBUG.Printf("options : %+v\r\n", auth)

	opts := docker.PullImageOptions{
		Repository: dockerImage,
		Tag:        tag,
	}

	DEBUG.Printf("options : %+v\r\n", opts)

	err := dockerengine.client.PullImage(opts, auth)
	if err != nil {
		ERROR.Printf("failed to pull this image %s, error %v\r\n", dockerImage, err)
		return "", err
	}

	// check if the image exists now
	resp, err := dockerengine.client.InspectImage(dockerImage)
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

func (dockerengine *DockerEngine) InspectImage(dockerImage string) bool {
	_, err := dockerengine.client.InspectImage(dockerImage)
	if err != nil {
		INFO.Printf("operator image %s does not exist locally\r\n", dockerImage)
		return false
	} else {
		INFO.Printf("operator image %s exists locally\r\n", dockerImage)
		return true
	}
}

// Ask the kernel for a free open port that is ready to use
func (dockerengine *DockerEngine) findFreePortNumber() int {
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

//functionCode string, taskID string, adminCfg []interface{}, servicePorts []string)
func (dockerengine *DockerEngine) StartTask(task *ScheduledTaskInstance, brokerURL string) (string, string, error) {
	dockerImage := task.DockerImage
	INFO.Println("to execute Task [", task.ID, "] to perform Operation [",
		dockerImage, "] with parameters [", task.Parameters, "]")

	// function code
	functionCode := task.FunctionCode

	// find a free listening port number available on the host machine
	freePort := strconv.Itoa(dockerengine.findFreePortNumber())

	// configure the task with its output streams via its admin interface
	commands := make([]interface{}, 0)

	// set broker URL
	setBrokerCmd := make(map[string]interface{})
	setBrokerCmd["command"] = "CONNECT_BROKER"
	setBrokerCmd["brokerURL"] = brokerURL
	commands = append(commands, setBrokerCmd)

	// pass the reference URL to the task so that the task can issue context subscription as well
	setReferenceCmd := make(map[string]interface{})
	setReferenceCmd["command"] = "SET_REFERENCE"
	setReferenceCmd["url"] = "http://" + dockerengine.workerCfg.InternalIP + ":" + freePort
	commands = append(commands, setReferenceCmd)

	// set output stream
	for _, outStream := range task.Outputs {
		setOutputCmd := make(map[string]interface{})
		setOutputCmd["command"] = "SET_OUTPUTS"
		setOutputCmd["type"] = outStream.Type
		setOutputCmd["id"] = outStream.StreamID
		commands = append(commands, setOutputCmd)
	}

	// check if it is required to set up the portmapping for its endpoint services
	servicePorts := make([]string, 0)

	for _, parameter := range task.Parameters {
		// deal with the service port
		if parameter.Name == "service_port" {
			servicePorts = strings.Split(parameter.Value, ";")
		}
	}

	// prepare the configuration for a docker container, host mode for the container network
	evs := make([]string, 0)
	evs = append(evs, fmt.Sprintf("myport=%s", freePort))

	// pass the initial configuration as the environmental variable
	jsonString, _ := json.Marshal(commands)
	evs = append(evs, fmt.Sprintf("adminCfg=%s", jsonString))

	config := docker.Config{Image: dockerImage, Env: evs}

	hostConfig := docker.HostConfig{}

	//if runtime.GOOS == "darwin" {   already use the bridge model
	internalPort := docker.Port(freePort + "/tcp")
	portBindings := map[docker.Port][]docker.PortBinding{
		internalPort: []docker.PortBinding{docker.PortBinding{HostIP: "0.0.0.0", HostPort: freePort}}}

	config.ExposedPorts = map[docker.Port]struct{}{internalPort: {}}

	// add the other listening ports into the exposed port list
	for _, port := range servicePorts {
		internalPort := docker.Port(port + "/tcp")
		portBindings[internalPort] = []docker.PortBinding{docker.PortBinding{HostIP: "0.0.0.0", HostPort: port}}
	}

	hostConfig.PortBindings = portBindings

	// to configure if the container will be removed once it is terminated
	hostConfig.AutoRemove = dockerengine.workerCfg.Worker.ContainerAutoRemove

	if functionCode != "" {
		fileName := "/tmp/" + task.ID
		dockerengine.writeTempFile(fileName, functionCode)

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
	container, err := dockerengine.client.CreateContainer(containerOptions)
	if err != nil {
		ERROR.Println(err)
		return "", freePort, err
	}

	// start the new container
	err = dockerengine.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		ERROR.Println(err)
		return "", freePort, err
	}

	refURL := "http://" + dockerengine.workerCfg.InternalIP + ":" + freePort

	return container.ID, refURL, nil
}

func (dockerengine *DockerEngine) StopTask(containerID string) {
	go dockerengine.client.StopContainer(containerID, 1)
}

func (dockerengine *DockerEngine) writeTempFile(fileName string, fileContent string) {
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
