package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8s.io/client-go/rest"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type Kubernetes struct {
	workerCfg            *Config
	clientset            *kubernetes.Clientset
	applicationNameSpace string
}

func (k8s *Kubernetes) Init(cfg *Config) bool {
	k8s.workerCfg = cfg
	k8s.applicationNameSpace = cfg.Worker.AppNameSpace

	config, err := rest.InClusterConfig()
	if err != nil {
		ERROR.Println(err.Error())
		return false
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		ERROR.Println(err.Error())
		return false
	}

	k8s.clientset = clientset

	return true
}

func readFileConfig() *rest.Config {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		ERROR.Println(err.Error())
		return nil
	}

	return config
}

// creates the in-cluster config
func fetchClusterConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		ERROR.Println(err.Error())
		return nil
	}

	return config
}

func (k8s *Kubernetes) PullImage(dockerImage string, tag string) (string, error) {
	return "test", nil
}

// Ask the kernel for a free open port that is ready to use
func (k8s *Kubernetes) findFreePortNumber() int {
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

func (k8s *Kubernetes) StartTask(task *ScheduledTaskInstance, brokerURL string) (string, string, error) {
	dockerImage := task.DockerImage
	taskName := strings.ToLower(task.OperatorName)

	// find a free listening port number available on the host machine
	freePort := strconv.Itoa(k8s.findFreePortNumber())

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

	INFO.Println("[namespace for applications]", k8s.applicationNameSpace)

	deploymentsClient := k8s.clientset.AppsV1().Deployments(k8s.applicationNameSpace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fogflow-task-" + freePort,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": taskName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": taskName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  taskName,
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
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	coreV1Client := k8s.clientset.CoreV1()

	serviceSpec := &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: k8s.applicationNameSpace,
			Name:      "fogflow-task-" + freePort,
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": taskName,
			},
			Ports: []coreV1.ServicePort{
				{
					Port: pport,
				},
			},
		},
	}

	service, err := coreV1Client.Services(k8s.applicationNameSpace).Create(context.TODO(), serviceSpec, metaV1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	refURL := "http://" + service.Spec.ClusterIP + ":" + freePort
	fmt.Printf("Created service at %s\n", refURL)

	return result.GetObjectMeta().GetName(), refURL, err
}

func (k8s *Kubernetes) StopTask(podId string) {
	deploymentsClient := k8s.clientset.AppsV1().Deployments(k8s.applicationNameSpace)

	fmt.Println("Deleting Deployment ", podId)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), podId, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deployment Deleted : ", podId)

	coreV1Client := k8s.clientset.CoreV1()
	err2 := coreV1Client.Services(k8s.applicationNameSpace).Delete(context.TODO(), podId, metaV1.DeleteOptions{})
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Deleted service %s\n", podId)
}
