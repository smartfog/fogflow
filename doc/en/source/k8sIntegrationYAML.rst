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

**Step 2** : Now, traverse to yaml folder in Fogflow repository using the **"fogflow/yaml"** path.

**Step 3** : Now create a namespace in kubernetes, using following command :

.. code-block:: console

    $kubectl create ns <user-specified> //E.g. kubectl create ns fogflow

**Step 4** : Now, in order to deploy cloud components, traverse to cloud-yaml folder and configure serviceaccount.yaml file. Edit the namespace with the one created in previous step and provide a serviceaccount name as shown below.

.. code-block:: console

    apiVersion: v1
    kind: ServiceAccount
    metadata:
    namespace: fogflow //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns  //User can provide this as per his/her choice
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-role
    rules:
    - apiGroups: [""]
    resources: ["services"]
    verbs: ["get","watch","list","create"]
    - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get","watch","list","create"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-viewer
    roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: fogflow-dns-role
    subjects:
    - kind: ServiceAccount
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-role

**Step 5** : Configure config.json as per user's environnment like shown below:

- **my_hostip**: this is the IP of your host machine, which should be accessible for both the web browser on your host machine and docker containers. Please DO NOT use "127.0.0.1" for this.
- **site_id**: each FogFlow node (either cloud node or edge node) requires to have a unique string-based ID to identify itself in the system;
- **physical_location**: the geo-location of the FogFlow node;
- **worker.capacity**: it means the maximal number of docker containers that the FogFlow node can invoke;  

.. code-block:: console

    "my_hostip": "172.30.48.24", //User should update the IP as per his/her environment
    "physical_location":{
        "longitude": 139.709059,
        "latitude": 35.692221
    },
    "site_id": "001",
    "worker": {
    "container_autoremove": false,
    "start_actual_task": true,
    "capacity": 8
    }

**Step 6** : Edit the namespace, serviceaccount value and configjson in dgraph-deployment.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f dgraph-deployment.yaml 

**Step 7** : Edit the namespace, serviceaccount value and configjson path in discovery.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f discovery.yaml 

**Step 8** : Edit the namespace, serviceaccount value and configjson path in cloud-broker.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f cloud-broker.yaml 

**Step 9** : Edit the namespace, serviceaccount value and configjson path in designer.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f designer.yaml 
    
**Step 10** : Edit the namespace, serviceaccount value and nginxConf path in nginx.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f nginx.yaml 

**Step 11** : Edit the namespace, serviceaccount value and configjson path in rabbitmq.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f rabbitmq.yaml 

**Step 12** : Edit the namespace, serviceaccount value and configjson path in master.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f master.yaml 

**Step 13** : Edit the namespace, serviceaccount value and configjson path in worker.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $kubectl create -f worker.yaml 


Now verify the deployments using, 

1. Fogflow dashboard : In your browser, type for http://<my_hostip>:80

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

.. important:: 

    Inorder to setup RBAC, use RBAC_setup.yaml file. Configure the namespace as per user (previously created namespace in step 3). Then use below command :

    $kubectl create -f RBAC_setup.yaml


Deploy FogFlow Edge Components on MicroK8s Environment Using YAML Files
-----------------------------------------------------------------------------

Fogflow Edge can be deployed on Microk8s cluster using individual YAML files if user whish to do so. To accomplish that following are the prerequisites :

Microk8s cluster

Important

To install microk8s, please refer to these `steps`_.

.. _`steps` : https://github.com/smartfog/fogflow/blob/k8s_manual_update/doc/en/source/k8sIntegration.rst#microk8s-installation-and-setup

To deploy edge components, follow below steps:

**Step 1** : Clone the github repository of Fogflow using this `repository link`_ if not already present.

.. _`repository link` : https://github.com/smartfog/fogflow

**Step 2** : Now, traverse to yaml folder in Fogflow repository using the **"fogflow/yaml/"** path.

**Step 3** : Now create a namespace in microk8s, using following command :

.. code-block:: console

    $microk8s.kubectl create ns <user-specified> //E.g. microk8s.kubectl create ns fogflow

**Step 4** : Now, in order to deploy edge components, traverse to edge-yaml folder and configure serviceaccount.yaml file. Edit the namespace with the one created in previous step and provide a serviceaccount name as shown below.

.. code-block:: console

    apiVersion: v1
    kind: ServiceAccount
    metadata:
    namespace: fogflow //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns  //User can provide this as per his/her choice
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-role
    rules:
    - apiGroups: [""]
    resources: ["services"]
    verbs: ["get","watch","list","create"]
    - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get","watch","list","create"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-viewer
    roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: fogflow-dns-role
    subjects:
    - kind: ServiceAccount
    namespace: fogflow  //Edit this as per previous step , for example namespace : fogflow
    name: fogflow-dns-role

**Step 5** : Configure config.json as per user's environnment like shown below:

- **coreservice_ip**: this is the IP where cloud node is running.
- **my_hostip**: this is the IP of your host machine, which should be accessible for both the web browser on your host machine and docker containers. Please DO NOT use "127.0.0.1" for this.
- **site_id**: each FogFlow node (either cloud node or edge node) requires to have a unique string-based ID to identify itself in the system;
- **physical_location**: the geo-location of the FogFlow node;
- **worker.capacity**: it means the maximal number of docker containers that the FogFlow node can invoke;  

.. code-block:: console

    "coreservice_ip": "172.30.48.46", //User should update the IP as per his/her environment where cloud is running
    "my_hostip": "172.30.48.24", //User should update the IP as per his/her environment i.e. the IP of host machine where edge is running
    "physical_location":{
        "longitude": 139.709059,
        "latitude": 35.692221
    },
    "site_id": "002",
    "worker": {
    "container_autoremove": false,
    "start_actual_task": true,
    "capacity": 4
    }

**Step 6** : Edit the namespace, serviceaccount value and configjson path in edge-broker.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $microk8s.kubectl create -f edge-broker.yaml 

**Step 7** : Edit the namespace, serviceaccount value and configjson path in worker.yaml as per user's environment and use below command to launch the deployments.

.. code-block:: console

    $microk8s.kubectl create -f worker.yaml

 
Now verify the deployments using, 

1. Check for pods status, using **microk8s.kubectl get pods --namespace=fogflow**

.. code-block:: console

    NAME                           READY   STATUS              RESTARTS   AGE
    edge-broker-c78679dd8-gx5ds    1/1     Running             0          8s
    worker-db94ff4f7-hwx72         1/1     Running             0          8s
    

.. important:: 

    Inorder to setup RBAC, use RBAC_setup.yaml file. Configure the namespace as per user (previously created namespace in step 3). Then use below command :

    $microk8s.kubectl create -f RBAC_setup.yaml
