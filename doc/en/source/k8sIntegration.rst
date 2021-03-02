**********************
Kubernetes Integration
**********************

The components of FogFlow can be built via source code as well as in docker environment using docker-compose tool. In docker environment each component of FogFlow is running running as single instance. Whole FogFlow system will have to re-start in case any single component container goes down and if any single service is overloaded it cannot scale to handle the load.  
To overcome these issues FogFlow has migrated to Kubernetes. FogFlow components will be deployed in Kubernetes cluster environment based on end user requirement. Various cluster configuration can be deployed:

1.	Master and Worker on same Node
2.	Single Master and Single Worker
3.	Single Master and Multi Worker Node
4.	Multiple Master and Multi Worker Node


Along with cluster following features of K8s are implemented in FogFlow:

1. **High Availability and Load Balancing**:
High Availability is about setting up Kubernetes, along with its supporting components in a way that there is no single point of failure. If the environment setup has multiple applications running on Single containers that container can easily fail. Same as the virtual machines for high availability in Kubernetes multiple replicas of containers can be run. Load balancing is efficient in distributing incoming network traffic across a group of backend servers.
A load balancer is a device that distributes network or application traffic across a cluster of servers. The load balancer has a big role to achieve high availability and performance increase of cluster.
 
2. **Self-healing**: if any of the pod are deleted manually or a pod got deleted accidentally or restarted. The deployment will make sure that it brings back the pod because Kubernetes has a feature to auto-heal the pods.

3. **Automated Rollouts & Rollback**: can be achieved by rolling update. Rolling updates are the default strategy to update the running version of your app. It updates cycles previous Pod out and bring newer Pod in incrementally.
When any introduced change that breaks production, then there should have a plan to roll back that change Kubernetes and kubectl offer a simple mechanism to roll back changes to resources such as Deployments.

4. **Ease the deployment with Helm Support**: Helm is a tool that streamlines installing and managing Kubernetes applications. It helps in managing Kubernetes applications. 
Helm Chart helps to define, install, and upgrade even the most complex Kubernetes application.
FogFlow document would be updated with the functioning details of above features to understand and access the Kubernetes environment well.


**Limitation of FogFlow K8s Integration**

Below are few limitations of FogFlow Kubernetes Integration. These limitation will be implemented with FogFlow in future.


1. Security and Network Policy in K8s environment.

2. Taints and Trait

3. Performance Evaluation

4. Other Functionality


FogFlow Cloud architecture diagram on Kubernetes
----------------------------------------------




.. figure:: figures/k8s-architecture.png





FogFlow cloud node components such as Dgraph, Discovery, Broker, Designer, Master, Worker, Rabbitmq are distributed in cluster nodes. The communication between FogFlow components and their behaviour are as previous and the worker node will launch task instances on docker container. 



Follow the link `here`_ to know how Kubernetes component works

.. _`here`: https://kubernetes.io/docs/concepts/overview/components/



Here are the prerequisite commands for running FogFlow on K8s:

1. docker
2. Kubernetes
3. Helm

.. important:: 
	**please also allow your user to execute the Docker Command without Sudo**
	
To install Kubernetes, please refer to  `Kubernetes Official Site`_ or Check alternate `Install Kubernetes`_,

To install Helm, please refer `Install Helm`_

.. _`Kubernetes Official Site`: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

.. _`Install Kubernetes`: https://medium.com/@vishal.sharma./installing-configuring-kubernetes-cluster-on-ubuntu-18-04-lts-hosts-f37b959c8410

.. _`Install Helm`: https://helm.sh/docs/intro/install/


Deploy FogFlow Cloud Components on K8s Environment
--------------------------------------------------

FogFlow cloud node components such as Dgraph, Discovery, Broker, Designer, Master, Worker, Rabbitmq are distributed in cluster nodes. The communication between FogFlow components and their behaviour are as usual and the worker node will launch task instances on docker container. 


**Fetch all required scripts**

Download the Kubernetes file and the configuration files as below.

.. code-block:: console    

	# the Kubernetes yaml file to start all FogFlow components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/helm/fogflow-chart
	
	# the configuration file used by all FogFlow components
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/yaml/config.json

	# the configuration file used by the nginx proxy
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/yaml/nginx.conf

	
   
Change the IP configuration accordingly
-------------------------------------------------------------

You need to change the following IP addresses in config.json according to your own environment.

- **coreservice_ip**: it is used by all FogFlow edge nodes to access the core services (e.g., nginx on port 80 and rabbitmq on port 5672) on the FogFlow cloud node; usually this will be the public IP of the FogFlow cloud node.
- **external_hostip**: for the configuration of the FogFlow cloud node, this is the same as coreservice_ip used by the components (Cloud Worker and Cloud Broker) to access the running FogFlow core services;        
- **internal_hostip**: this is the IP of your default K8s network Interface, which is the "cni0" network interface on your Linux host.

- **site_id**: each FogFlow node (either cloud node or edge node) requires to have a unique string-based ID to identify itself in the system;
- **physical_location**: the geo-location of the FogFlow node;
- **worker.capacity**: it means the maximal number of docker containers that the FogFlow node can invoke;  


Change values.yaml file
---------------------------

-Edit namespace as per requirement. Add the no. of replicaCount required.

-Change dgraph, configJson and nginxConf path in values.yaml  file as per the environment hostPath.

-Change externalIPs as per the environment.

.. code-block:: console

      #Kubernetes namespace of FogFlow components
      namespace: default

      #replicas will make sure that no. of replicaCount mention in values.yaml
      #are running all the time for the deployment
      replicaCount: 1

      serviceAccount:
      #Specifies whether a service account should be created
        create: true
      #Annotations to add to the service account
        annotations: {}
      #The name of the service account to use.
      #If not set and create is true, a name is generated using the fullname template
        name: ""

      #hostPath for dgraph volume mount
      dgraph:
        hostPath:
          path: /mnt/dgraph

      #hostPath for config.json
      configJson:
        hostPath:
          path: /home/necuser/fogflow/helm/files/fogflow-chart/config.json

      #hostPath for nginx.conf
      nginxConf:
        hostPath:
          path: /home/necuser/fogflow/fogflow/yaml/nginx.conf

      #External IP to expose cluster
      Service:
       spec:
        externalIPs:
        - XXX.XX.48.24

	  
Start all Fogflow components with Helm Chart
-------------------------------------------------------------

Execute Helm command from outside the Helm-Chart folder to start FogFlow Components, here helm-chart name is "fogflow-chart". 

Add "--set" flag with helm install command to pass configuration from command line.

.. code-block:: console
 
          helm install ./fogflow-chart --set externalIPs={XXX.XX.48.24} --generate-name


Refer Helm official `link`_ for more details

.. _`link`: https://helm.sh/docs/helm/

Validate the setup
-------------------------------------------------------------

There are two ways to check if the FogFlow cloud node is started correctly: 

- Check all the Pods are Up and Running using "kubectl get pods --namespace=<namespace_name>"

.. code-block:: console  

         kubectl get pods --namespace=fogflow
		 
		 
        NAME                           READY   STATUS              RESTARTS   AGE
        cloud-broker-c78679dd8-gx5ds   1/1     Running             0          8s
        cloud-worker-db94ff4f7-hwx72   1/1     Running             0          8s
        designer-bf959f7b7-csjn5       1/1     Running             0          8s
        dgraph-869f65597c-jrlqm        1/1     Running             0          8s
        discovery-7566b87d8d-hhknd     1/1     Running             0          8s
        master-86976888d5-drfz2        1/1     Running             0          8s
        nginx-69ff8d45f-xmhmt          1/1     Running             0          8s
        rabbitmq-85bf5f7d77-c74cd      1/1     Running             0          8s

		
- Check the system status from the FogFlow DashBoard

System status can also be verified from FogFlow dashboard on web browser to see the current system status via the URL: http://<coreservice_ip>/index.html


**Launch FogFlow Worker task instances on k8s pods**



Previously, task instances was launching on Docker containers. In that case, FogFlow worker used to call go-dockerclient and with the help of this client was launching task instances on docker containers.



.. figure:: figures/dockerTaskInstance.png





Now, An interface pod.go is being used. This interface will receive the dockerimage name, port and configuration details from worker and launch the task inside a pod.



.. figure:: figures/podTaskInstance.png 





Create a task using link `task_Instance`_

.. _`task_Instance`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html#define-a-dummy-fog-function 




**FogFlow Edge node Kubernetes support**
