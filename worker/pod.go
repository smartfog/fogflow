package main

import (
        "context"
        "fmt"

        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
        appsv1 "k8s.io/api/apps/v1"
        apiv1 "k8s.io/api/core/v1"
//        CreateOptions "k8s.io/apimachinery/pkg/apis/meta/v1"
//        Deployment "k8s.io/api/apps/v1"
        //
        // Uncomment to load all auth plugins
        // _ "k8s.io/client-go/plugin/pkg/client/auth"
        //
        // Or uncomment to load specific auth plugins
        // _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
        // _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
        // _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
        // _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type pod struct{}

func (p *pod) pod(dockerImage string, port string) (string, error){
        // creates the in-cluster config
        config, err := rest.InClusterConfig()
        if err != nil {
                panic(err.Error())
        }
        // creates the clientset
	fmt.Println("******port num is**********",port)
        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err.Error())
        }
                deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

        deployment := &appsv1.Deployment{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "fogflow-deployment",
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
                                                                        ContainerPort: 80,
                                                                },
                                                        },
                                                },
                                        },
                                },
                        },
                },
        }

        // Create Deployment
        fmt.Println("Creating deployment...")
        result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
        if err != nil {
                panic(err)
        }
        fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return dockerImage, err
}

