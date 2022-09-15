package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type EdgeController struct {
	workerCfg         *Config
	edgeControllerURL string
}

func (mec *EdgeController) Init(cfg *Config) bool {
	// retrieve the accessible URL of the edge controller
	mec.edgeControllerURL = cfg.GetEdgeControllerURL()
	mec.workerCfg = cfg

	return true
}

func (mec *EdgeController) PullImage(dockerImage string, tag string) (string, error) {
	return "test", nil
}

// Ask the kernel for a free open port that is ready to use
func (mec *EdgeController) findFreePortNumber() int {
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

func (mec *EdgeController) StartTask(task *ScheduledTaskInstance, brokerURL string) (string, string, error) {
	dockerImage := task.DockerImage

	// find a free listening port number available on the host machine
	freePort := strconv.Itoa(mec.findFreePortNumber())

	// configure the task with its output streams via its admin interface
	commands := make([]interface{}, 0)

	// set broker URL
	setBrokerCmd := make(map[string]interface{})
	setBrokerCmd["command"] = "CONNECT_BROKER"
	setBrokerCmd["brokerURL"] = brokerURL
	commands = append(commands, setBrokerCmd)

	// set CorrelatorID
	setCorrelatorCmd := make(map[string]interface{})
	setCorrelatorCmd["command"] = "SET_CORRELATORID"
	setCorrelatorCmd["correlatorID"] = task.ID
	commands = append(commands, setCorrelatorCmd)

	// pass the reference URL to the task so that the task can issue context subscription as well
	setReferenceCmd := make(map[string]interface{})
	setReferenceCmd["command"] = "SET_REFERENCE"
	setReferenceCmd["url"] = "http://fogflow-deployment-" + freePort + ":" + freePort
	commands = append(commands, setReferenceCmd)

	// set output stream
	for _, outStream := range task.Outputs {
		setOutputCmd := make(map[string]interface{})
		setOutputCmd["command"] = "SET_OUTPUTS"
		setOutputCmd["type"] = outStream.Type
		setOutputCmd["id"] = outStream.StreamID
		commands = append(commands, setOutputCmd)
	}

	jsonString, _ := json.Marshal(commands)

	iport, err := strconv.ParseInt(freePort, 10, 32)
	pport := int32(iport)
	if err != nil {
		panic(err.Error())
	}

	taskId := "fogflow-deployment-" + freePort
	appNameSpace := mec.workerCfg.Worker.AppNameSpace

	deployment := &appsv1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: appNameSpace,
			Name:      taskId,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: dockerImage,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: pport,
									HostPort:      pport,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "myport",
									Value: freePort,
								},
								{
									Name:  "adminCfg",
									Value: string(jsonString),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	jsonPayload, err := json.Marshal(deployment)
	if err != nil {
		ERROR.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}

	requestURL := mec.edgeControllerURL + "/api/v1/create/deployment/" + appNameSpace
	mec.sendRequest("POST", requestURL, task.Parameters, jsonPayload)

	serviceSpec := &coreV1.Service{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: mec.workerCfg.Worker.AppNameSpace,
			Name:      taskId,
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": "demo",
			},
			Ports: []coreV1.ServicePort{
				{
					Port: pport,
				},
			},
			Type: "ClusterIP",
		},
	}

	// create a k8s service
	jsonPayload, err = json.Marshal(serviceSpec)
	if err != nil {
		ERROR.Fatalf("Error occured during marshaling. Error: %s", err.Error())
		return "", "", err
	}

	requestURL = mec.edgeControllerURL + "/api/v1/create/service/" + appNameSpace

	resp, err := mec.sendRequest("POST", requestURL, task.Parameters, jsonPayload)
	if err != nil {
		ERROR.Fatalf("NOT able to interact with Edge Controller: %s", err.Error())
		return "", "", err
	}

	msg := make(map[string]string)
	if err := json.Unmarshal(resp, &msg); err != nil {
		ERROR.Fatalf("fail to extrac the response from Edge Controller: %s", err.Error())
		return "", "", err
	}

	refURL := "http://" + msg["cluster_ip"] + ":" + freePort
	fmt.Printf("Created service at %s\n", refURL)

	return taskId, refURL, nil
}

func (mec *EdgeController) StopTask(taskId string) {
	appNameSpace := mec.workerCfg.Worker.AppNameSpace

	deploymentName := taskId
	requestURL := mec.edgeControllerURL + "/api/v1/delete/deployment/" + appNameSpace + "/" + deploymentName
	mec.sendRequest("DELETE", requestURL, nil, nil)

	serviceName := taskId
	requestURL = mec.edgeControllerURL + "/api/v1/delete/service/" + appNameSpace + "/" + serviceName
	mec.sendRequest("DELETE", requestURL, nil, nil)
}

func (mec *EdgeController) sendRequest(method string, url string, parameters []Parameter, payload []byte) ([]byte, error) {
	INFO.Println(method, url, string(payload))

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// set the header based on the parameters of the operator
	for _, parameter := range parameters {
		request.Header.Set(parameter.Name, parameter.Value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ERROR.Println("Failed to interact with the MEC edge controller ", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return body, err
}
