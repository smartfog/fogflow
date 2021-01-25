.. _cloud-setup:

*****************************************
System Setup
*****************************************

Prerequisite
=================

Here are the prerequisite commands for starting FogFlow:

1. docker

2. docker-compose

For ubuntu-16.04, you need to install docker-ce and docker-compose.

To install Docker CE, please refer to `Install Docker CE`_, required version > 18.03.1-ce;

.. important:: 
	**please also allow your user to execute the Docker Command without Sudo**


To install Docker Compose, please refer to `Install Docker Compose`_, 
required version 18.03.1-ce, required version > 2.4.2

.. _`Install Docker CE`: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
.. _`Install Docker Compose`: https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04



Start FogFlow Cloud node
=============================

Fetch all required scripts
---------------------------------

Download the docker-compose file and the configuration files as below.

.. code-block:: console    

	# the docker-compose file to start all FogFlow components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/docker-compose.yml

	# the configuration file used by all FogFlow components
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/config.json

	# the configuration file used by the nginx proxy
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/nginx.conf
	

Change the IP configuration accordingly
---------------------------------------------


You need to change the following IP addresses in config.json according to your own environment.


- **my_hostip**: the IP of the FogFlow cloud node and this IP address should be accessible to the FogFlow edge node. Please DO NOT use "127.0.0.1" for this. 
- **site_id**: each FogFlow node (either cloud node or edge node) requires to have a unique string-based ID to identify itself in the system;
- **physical_location**: the geo-location of the FogFlow node;
- **worker.capacity**: it means the maximal number of docker containers that the FogFlow node can invoke;  


Change the IP configuration of elasticsearch and metricbeat accordingly
---------------------------------------------------------------------------

You need to change the following IP addresses in docker-compose.yml according to your own environment.

- **output.elasticsearch.hosts**: it is the host location of elasticsearch on which metricbeat shares data in csv format.

Also need to change the following IP addresses in metricbeat.docker.yml according to your own environment.

- **name**: It is the name given for uniqueness for cloud nodes from egde nodes on grafana metric dashboard. You can mention any name in place of IP address.

- **hosts**: It is the host location of elasticsearh database, where metricbeat is going to share metric data.


.. important:: 

	please DO NOT use "127.0.0.1" as the IP address of **my_hostip** , because they will be used by a running task inside a docker container. 
	
	**Firewall rules:** to make your FogFlow web portal accessible via the external_ip; the following ports must be open as well: 80 and 5672 for TCP



Start all components on the FogFlow Cloud Node
------------------------------------------------------


Pull the docker images of all FogFlow components and start the FogFlow system

.. code-block:: console    

    # if you already download the docker images of FogFlow components, this command can fetch the updated images
	docker-compose pull  

	docker-compose up -d


Validate your setup
----------------------------------


There are two ways to check if the FogFlow cloud node is started correctly: 


- Check all the containers are Up and Running using "docker ps -a"

.. code-block:: console    

	docker ps -a
	
	CONTAINER ID      IMAGE                       COMMAND                  CREATED             STATUS              PORTS                                                 NAMES
	90868b310608      nginx:latest            "nginx -g 'daemon of…"   5 seconds ago       Up 3 seconds        0.0.0.0:80->80/tcp                                       fogflow_nginx_1
	d4fd1aee2655      fogflow/worker          "/worker"                6 seconds ago       Up 2 seconds                                                                 fogflow_cloud_worker_1
	428e69bf5998      fogflow/master          "/master"                6 seconds ago       Up 4 seconds        0.0.0.0:1060->1060/tcp                               fogflow_master_1
	9da1124a43b4      fogflow/designer        "node main.js"           7 seconds ago       Up 5 seconds        0.0.0.0:1030->1030/tcp, 0.0.0.0:8080->8080/tcp       fogflow_designer_1
	bb8e25e5a75d      fogflow/broker          "/broker"                9 seconds ago       Up 7 seconds        0.0.0.0:8070->8070/tcp                               fogflow_cloud_broker_1
	7f3ce330c204      rabbitmq:3              "docker-entrypoint.s…"   10 seconds ago      Up 6 seconds        4369/tcp, 5671/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp     fogflow_rabbitmq_1
	9e95c55a1eb7      fogflow/discovery       "/discovery"             10 seconds ago      Up 8 seconds        0.0.0.0:8090->8090/tcp                               fogflow_discovery_1
        399958d8d88a      grafana/grafana:6.5.0   "/run.sh"                29 seconds ago      Up 27 seconds       0.0.0.0:3003->3000/tcp                               fogflow_grafana_1
        9f99315a1a1d      fogflow/elasticsearch:7.5.1 "/usr/local/bin/dock…" 32 seconds ago    Up 29 seconds       0.0.0.0:9200->9200/tcp, 0.0.0.0:9300->9300/tcp       fogflow_elasticsearch_1
        57eac616a67e      fogflow/metricbeat:7.6.0 "/usr/local/bin/dock…"   32 seconds ago     Up 29 seconds                                                                  fogflow_metricbeat_1
	
	
.. important:: 

	if you see any container is missing, you can run "docker ps -a" to check if any FogFlow component is terminated with some problem. If there is, you can further check its output log by running "docker logs [container ID]"


- Check the system status from the FogFlow DashBoard

You can open the FogFlow dashboard in your web browser to see the current system status via the URL: http://<coreservice_ip>/index.html

.. important:: 

	If the FogFlow cloud node is behind a gateway, you need to create a mapping from the gateway IP to the coreservice_ip and then access the FogFlow dashboard via the gateway IP;
	If the FogFlow cloud node is a VM in a public cloud like Azure Cloud, Google Cloud, or Amazon Cloud, you need to access the FogFlow dashboard via the public IP of your VM;

Once you are able to access the FogFlow dashboard, you can see the following web page

.. figure:: figures/dashboard.png


Configure Elasticsearch on Grafana Dashboard
-------------------------------------------------------------

Grafana dashboard can be accessible on web browser to see the current system status via the URL: 
http://<output.elasticsearch.hosts>:3003/. The default username and password for Grafana login are admin and admin respectively.


- After successful login to grafana, click on "Create your first data source" on Home Dashboard to setup the source of data.
- Select Elasticsearch from Add Data Sourch page. Now you are on page Data Sources/Elasticsearch same as below figure.


.. figure:: figures/Elastic_config.png


1. Put a name for the Data Source.
2. In HTTP detail ,mention URL of your elasticsearch and Port. URL shall include HTTP. 
3. In Access select Server(default). URL needs to be accessible from the Grafana backend/server.
4. In Elasticsearch details, put @timestamp for Time field name. Here a default for the time field can be specified with the name of your Elasticsearch index. Use a time pattern for the index name or a wildcard.
5. Select Elasticsearch Version.

Then click on "Save & Test" button.

Set up the Metricbeat
---------------------------------------------


- Change the details of Elasticsearch in metricbeat.docker.yml file as below:


.. code-block:: json

        name: "<155.54.239.141_cloud>"
        metricbeat.modules:
        - module: docker
          #Docker module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","memory","network"]
          hosts: ["unix:///var/run/docker.sock"]
          period: 10s
          enabled: true
        - module: system
          #System module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","load","memory","network"]
          period: 10s

        output.elasticsearch:
          hosts: '155.54.239.141:9200'
	  
	  
	  
Try out existing IoT services
-------------------------------------


Once the FogFlow cloud node is set up, you can try out some existing IoT services without running any FogFlow edge node.
For example, you can try out a simple fog function as below.  

- Click "Operator Registry" in the top navigator bar to triger the initialization of pre-defined operators. 

After you first click "Operator Registry", a list of pre-defined operators will be registered in the FogFlow system. 
With a second click, you can see the refreshed list as shown in the following figure.

.. figure:: figures/operator-list.png


- Click "Service Topology" in the top navigator bar to triger the initialization of pre-defined service topologies. 

After you first click "Service Topology", a list of pre-defined topologies will be registered in the FogFlow system. 
With a second click, you can see the refreshed list as shown in the following figure.

.. figure:: figures/topology-list.png


- Click "Fog Function" in the top navigator bar to triger the initialization of pre-defined fog functions. 

After you first click "Fog Function", a list of pre-defined functions will be registered in the FogFlow system. 
With a second click, you can see the refreshed list as shown in the following figure.

.. figure:: figures/function-list.png


- Create an IoT device entity to trigger the Fog Function

You can register a device entity via the device registration page: 
1) click "System Status"; 
2) click "Device";
3) click "Add";

Then you will see the following device registration page. 

.. figure:: figures/device-registration.png

- Check if the fog function is triggered


Check if a task is created under "Task" in System Management.**

.. figure:: figures/fog-function-task-running.png


Check if a Stream is created under "Stream" in System Management.**

.. figure:: figures/fog-function-streams.png



Start FogFlow edge node
==========================

Typically, an FogFlow edge node needs to deploy a Worker, an IoT broker and a system monitoring agent metricbeat. 
The Edge IoT Broker at the edge node can establish the data flows between all task instances launched on the same edge node. 
However, this Edge IoT Broker is optional, 
especially when the edge node is a very constrained device that can only support a few tasks without any data dependency. 

Here are the steps to start an FogFlow edge node: 

Install Docker Engine 
------------------------

To install Docker CE and Docker Compose, please refer to `Install Docker CE and Docker Compose on Respberry Pi`_. 

.. _`Install Docker CE and Docker Compose on Respberry Pi`: https://withblue.ink/2019/07/13/yes-you-can-run-docker-on-raspbian.html


.. note:: Docker engine must be installed on each edge node, because all task instances in FogFlow will be launched within a docker container.


Download the deployment script 
-------------------------------------------------

.. code-block:: console    
         
	#download the deployment scripts
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/http/start.sh
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/http/stop.sh 
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/http/metricbeat.docker.yml
	
	#make them executable
	chmod +x start.sh  stop.sh       


Download the default configuration file 
-------------------------------------------------

.. code-block:: console   
         	
	#download the configuration file          
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/http/config.json


Change the configuration file accordingly
-------------------------------------------------

You can use the default setting for a simple test, but you need to change the following addresses according to your own environment: 
        
- **coreservice_ip**: please refer to the configuration of the cloud part. This is the accessible address of your FogFlow core services running in the cloud node;
- **external_hostip**: this is the external IP address, accessible for the cloud broker. It is useful when your edge node is behind NAT;
- **my_hostip** is the IP of your default docker bridge, which is the "docker0" network interface on your host.
- **site_id** is the user-defined ID for the edge Node. Broker and Worker IDs on that node will be formed according to this Site ID.
- **container_autoremove** is used to configure that the container associated with a task will be removed once its processing is complete.
- **start_actual_task** configures the Fogflow worker to include all those activities that are required to start or terminate a task or maintain a running task along with task configurations instead of performing the minimal effort. It is recommended to keep it true.
- **capacity** is the maximum number of docker containers that the FogFlow node can invoke. The user can set the limit by considering resource availability on a node.

.. code-block:: json

    //you can see the following part in the default configuration file
    { 
        "coreservice_ip": "155.54.239.141", 
        "external_hostip": "35.234.116.177", 
        "my_hostip": "172.17.0.1", 
        
	
	"site_id": "002",
	
	
	"worker": {
        "container_autoremove": false,
        "start_actual_task": true,
        "capacity": 4
	}
	
	
    } 

Change the Metricbeat configuration file accordingly
-----------------------------------------------------------

you need to change the following addresses in start.sh according to your own environment:

- **output.elasticsearch.hosts**: It is the elasticsearch host IP address on which metricbeat will share the metric data.

- change the details of Elasticsearch in metricbeat.docker.yml file as below:

.. code-block:: json

        name: "<155.54.239.141/edge02>"
        metricbeat.modules:
        - module: docker
          #Docker module parameters to monitor based on user requirement,example as below
          metricsets: ["cpu","memory","network"]
          hosts: ["unix:///var/run/docker.sock"]
          period: 10s
          enabled: true
        - module: system
          #System module parameters to monitor based on user requirement, example as below
          metricsets: ["cpu","load","memory","network"]
          period: 10s

        output.elasticsearch:
          hosts: '155.54.239.141:9200'
	  

Start Edge node components
-------------------------------------------------

.. note:: if the edge node is ARM-basd, please attach arm as the command parameter

.. code-block:: console    

      #start both components in the same script
      ./start.sh 
    
      #if the edge node is ARM-basd, please attach arm as the command parameter
      #./start.sh  arm
      

Stop Edge node components
-------------------------------------------------

.. code-block:: console    

	#stop both components in the same script
	./stop.sh 


     










