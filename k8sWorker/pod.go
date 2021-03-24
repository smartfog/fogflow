package main

import (
	"context"
	"fmt"
	"strconv"
	docker "github.com/fsouza/go-dockerclient"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type pod struct {
	client *docker.Client
}

func (p *pod) pod(dockerImage string, port string, adminCfg []interface{}) (string, error) {

	var brokerURL, url, id, Etype string
	for _, commandValue := range adminCfg {
		commandValueMap := commandValue.(map[string]interface{})
		for key, value := range commandValueMap {
			fmt.Println("key", key)
			if key == "brokerURL" {
				brokerURL = value.(string)
			}
			if key == "url" {
				url = value.(string)
			}
			if key == "id" {
				id = value.(string)
			}
			if key == "type" {
				Etype = value.(string)
			}
		}
	}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	iport, err := strconv.ParseInt(port, 10, 32)
	pport := int32(iport)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fogflow-deployment" + port,
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
									Value: port,
								},
								{
									Name:  "brokerURL",
									Value: brokerURL,
								},
								{
									Name:  "url",
									Value: url,
								},
								{
									Name:  "id",
									Value: id,
								},
								{
									Name:  "type",
									Value: Etype,
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

	coreV1Client := clientset.CoreV1()

	serviceSpec := &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "fogflow-deployment" + port,
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": "demo",
			},
			Type: "LoadBalancer",
			Ports: []coreV1.ServicePort{
				{
					Port: pport,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: pport,
					},
				},
			},
		},
	}

	service, err := coreV1Client.Services(apiv1.NamespaceDefault).Create(context.TODO(), serviceSpec, metaV1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %s\n", service.ObjectMeta.Name)

	return result.GetObjectMeta().GetName(), err
}
