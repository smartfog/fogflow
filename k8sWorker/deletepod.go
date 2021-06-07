package main

import (
	"context"
	//"encoding/json"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	//appsv1 "k8s.io/api/apps/v1"
	//apiv1 "k8s.io/api/core/v1"
	//coreV1 "k8s.io/api/core/v1"
        metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type delpod struct {
	client *docker.Client
}

func (p *delpod) deletepod(podId string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	/*iport, err := strconv.ParseInt(port, 10, 32)
	pport := int32(iport)
	if err != nil {
		panic(err.Error())
	}*/
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1().Deployments("fogflow")

	fmt.Println("Deleting Deployment ",podId)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), podId, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil{
		panic(err)
	}
	fmt.Println("Deployment Deleted : ",podId)

	coreV1Client := clientset.CoreV1()

	err2 := coreV1Client.Services("fogflow").Delete(context.TODO(), podId, metaV1.DeleteOptions{})
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Deleted service %s\n", podId)
}
