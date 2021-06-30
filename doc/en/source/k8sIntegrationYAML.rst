******************************************
Kubernetes Integration Using YAML Files
******************************************

Fogflow can be deployed on kubernetes cluster using individual YAML files if user whish to do so. To accomplish that, following are the prerequisites :

1. docker
2. Kubernetes

.. important:: 
	**please also allow your user to execute the Docker Command without Sudo**
	
To install Kubernetes, please refer to  `Kubernetes Official Site`_ or Check alternate `Install Kubernetes`_,


.. _`Kubernetes Official Site`: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

.. _`Install Kubernetes`: https://medium.com/@vishal.sharma./installing-configuring-kubernetes-cluster-on-ubuntu-18-04-lts-hosts-f37b959c8410

Deploy FogFlow Cloud Components on K8s Environment Using YAML Files
--------------------------------------------------------------------

FogFlow cloud node components such as Dgraph, Discovery, Broker, Designer, Master, Worker, Rabbitmq are distributed in cluster nodes. The communication between FogFlow components and their behaviour are as usual and the worker node will launch task instances on kubernetes pod.

Inorder to setup the components, please refer the steps below:

**Step 1** : Clone the github repository of Fogflow using this `link`_.

.. _`link` : https://github.com/smartfog/fogflow

**Step 2** : Now, traverse to **"deployment/kubernetes/cloud-node"** folder in Fogflow repository.
  
**Step 3** : Edit the **externalIPs** value in nginx.yaml as per user's environment.

.. code-block:: console

    apiVersion: v1
    kind: Service
    metadata:
    namespace: fogflow                      
    name: nginx
    labels:
        run: nginx
    spec:
    type: LoadBalancer
    ports:
        - port: 80
        targetPort: 80
    selector:
        run: nginx
    externalIPs: [172.30.48.24]  //edit this
   
**Step 4** : Now, in order to deploy cloud-node components, use below command.

.. code-block:: console

    ./install.sh

Now verify the deployments using, 

1. Fogflow dashboard : In your browser, type for http://<externalIPs>:80 (externalIPs is the one mentioned in nginx.yaml file).

2. Check for pods status, using **kubectl get pods --namespace=fogflow**

.. code-block:: console

    NAME                           READY   STATUS              RESTARTS   AGE
    cloud-broker-c78679dd8-gx5ds   1/1     Running             0          8s
    cloud-worker-db94ff4f7-hwx72   1/1     Running             0          8s
    designer-bf959f7b7-csjn5       1/1     Running             0          8s
    dgraph-869f65597c-jrlqm        1/1     Running             0          8s
    discovery-7566b87d8d-hhknd     1/1     Running             0          8s
    master-86976888d5-drfz2        1/1     Running             0          8s
    nginx-69ff8d45f-xmhmt          1/1     Running             0          8s
    rabbitmq-85bf5f7d77-c74cd      1/1     Running             0          8s

In order to stop the deployments of Fogflow system, follow below command:

.. code-block:: console

    ./uninstall.sh