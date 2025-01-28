Quick Start
===========================================================


This is an one-page introductory tutorial to FogFlow.
In the FIWARE-based architecture, FogFlow can be used to dynamically trigger data processing functions 
between IoT devices and Orion Context Broker, 
for the purpose of transforming and preprocessing raw data at edge nodes (e.g., IoT gateways or Raspberry Pis).

The tutorial introduces a typical FogFlow system setup with a simple example to do anomaly detection at edges for temperature sensor 
data.
It explains an example usecase implementation using FogFlow and FIWARE Orion in integration with each other. 

Every time to implement a usecase FogFlow creates some internal NGSI entities such as operators, Docker image, Fog Function and service topology.
So these entity data are very important for FogFlow system and these are need to be stored somewhere. An entity data can not be stored in FogFlow memory
because memory is volatile and it will lose content when power is lost. To solve this issue FogFlow introduces `Dgraph`_  a persistent storage.
The persistent storage will store FogFlow entity data in the form of graph.



.. _`Dgraph`: https://dgraph.io/docs/get-started/




As shown in the following diagram, in this use case a connected temperature sensor sends an update message to the FogFlow system, 
which triggers some running task instance of a pre-defined fog function to generate some analytics result. 
The fog function is specified in advance via the FogFlow dashboard, 
however, it is triggerred only when the temperature sensor joins the sytem. In a real distributed setup, 
the running task instance will be deployed at the edge node closed to the temperature sensor. 
Once the generated analytics result is generated, 
it will be forwarded from the FogFlow system to Orion Context Broker. 
This is because a subscription with Orion Context Broker as the reference URL has been issued.  


.. figure:: figures/systemview.png



Here are the prerequisite commands for running FogFlow:

1. docker

2. docker-compose

For ubuntu-16.04, you need to install docker-ce and docker-compose.

To install Docker CE, please refer to `Install Docker CE`_, required version > 18.03.1-ce;


.. important:: 
	**please also allow your user to execute the Docker Command without Sudo `Docker Post-Install`_**



To install Docker Compose, please refer to `Install Docker Compose`_, 
required version 18.03.1-ce, required version > 2.4.2

.. _`Install Docker CE`: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
.. _`Install Docker Compose`: https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04
.. _`Docker Post-Install`: https://docs.docker.com/engine/install/linux-postinstall/






Fetch all required scripts
-------------------------------------------------------------

Download the docker-compose file and the configuration files as below.

.. code-block:: console    

	# the docker-compose file to start all FogFlow components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/refs/heads/development/release/latest/cloud/docker-compose.yml
	
	# the template of the configuration file used by all FogFlow components
	wget https://github.com/smartfog/fogflow/blob/development/release/latest/cloud/config-template.json

	# the configuration file for ngsix server
	wget https://raw.githubusercontent.com/smartfog/fogflow/refs/heads/development/release/latest/cloud/nginx.conf

	
Configure FogFlow for the tutorial
-------------------------------------------------------------

Generate your config.json file starting from the download template

.. code-block:: console    

	cp config-template.json config.json


In order to have the hello world running configure the config.json file as following:

.. code-block:: console  

	{
		...
		"designer": {
			...
			"doNotInitApplications": false
		},
		...
	}

The overall config.json might look as the following.

.. code-block:: console  

	{
		"my_hostip": "10.1.99.99",
		"physical_location":{
			"longitude": 139,
			"latitude": 35
		},
		"site_id": "001",
		"logging":{
			"info":"stdout",
			"error":"stdout",
			"protocol": "stdout",
			"debug": "discard"
		},
		"discovery": {
			"http_port": 8090,
			"storeOnDisk": false,
			"delayStoreOnFile" : 3
		},
		"broker": {
			"http_port": 8070,
			"heartbeat_interval": 30
		},     
		"master": {
			"ngsi_agent_port": 1060,
			"rest_api_port": 8010,
			"infinite_reconnection_tries": true    
		},
		"worker": {
			"container_autoremove": false,
			"start_actual_task": true,
			"capacity": 8,
			"heartbeat_interval": 30,
			"detection_duration": 10,
			"infinite_reconnection_tries": true
		},
		"designer": {
			"webSrvPort": 8080,
			"agentPort": 1030,
			"ldAgentPort":1090,
			"doNotInitApplications": false
		},    
		"rabbitmq": {
			"port": 5672,
			"username": "admin",
			"password":"mypass"
		},
		"https": {
			"enabled" : false
		},
		"persistent_storage": {
			"port": 9082
		}     
	}


Other configuration that might be useful:

- **site_id**: each FogFlow node (either cloud node or edge node) requires to have a unique string-based ID to identify itself in the system;
- **physical_location**: the geo-location of the FogFlow node;
- **worker.capacity**: it means the maximal number of docker containers that the FogFlow node can invoke;  

Change the IP configuration in the configuration file
-------------------------------------------------------------

You need to change the following IP addresses in config.json according to your own environment and also check if the used port nubmers are blocked by your firewall. 

- **my_hostip**: this is the IP of your host machine, which should be accessible for both the web browser on your host machine and docker containers. Please DO NOT use "127.0.0.1" for this. 


.. important:: 

	please DO NOT use "127.0.0.1" as the IP address of **my_hostip**, because it is only accessible to a 
	running task inside a docker container. 
	
	**Firewall rules:** To make FogFlow web portal accessible and for its proper functioning, the following ports must be free and open over TCP in host machine. 
	
	.. code-block:: console

		Component		Port

		Discovery		8090 
		Broker			8070
		Designer		8080
		Nginx			  80
		Rabbitmq		5672
		Task			launched over any port internally


	**Note : Task Instance is launched over dynamically assigned port, which is not predefined. So, users can possibly allow local ports using rule in his firewall. This will result in smooth functioning of Task Instances.**

	**Above mentioned port number(s) are default port number(s)**. If user needs to change the port number(s), please make sure the change is consistence in all the configuration files named as **"config.json"**.

	**Mac Users:** if you like to test FogFlow on your Macbook, please install Docker Desktop and also use "host.docker.internal" 
	as my_hostip in the configuration file.


Start all Fogflow components 
-------------------------------------------------------------


Pull the docker images of all FogFlow components and start the FogFlow system


.. code-block:: console    

	#if you already download the docker images of FogFlow components, this command can fetch the updated images
	docker-compose pull  

	docker-compose up -d


Validate your setup
-------------------------------------------------------------


There are two ways to check if the FogFlow cloud node is started correctly: 


- Check all the containers are Up and Running using "docker ps -a"


.. code-block:: console    

	docker-compose ps
	
	795e6afe2857   nginx:latest            "/docker-entrypoint.…"   About a minute ago   Up About a minute   0.0.0.0:80->80/tcp                                                                               fogflow_nginx_1
	33aa34869968   fogflow/worker:3.2.8      "/worker"                About a minute ago   Up About a minute                                                                                                    fogflow_cloud_worker_1
	e4055b5cdfe5   fogflow/master:3.2.8      "/master"                About a minute ago   Up About a minute   0.0.0.0:1060->1060/tcp                                                                           fogflow_master_1
	cdf8d4068959   fogflow/designer:3.2.8    "node main.js"           About a minute ago   Up About a minute   0.0.0.0:1030->1030/tcp, 0.0.0.0:8080->8080/tcp                                                   fogflow_designer_1
	56daf7f078a1   fogflow/broker:3.2.8      "/broker"                About a minute ago   Up About a minute   0.0.0.0:8070->8070/tcp                                                                           fogflow_cloud_broker_1
	51901ce6ee5f   fogflow/discovery:3.2.8   "/discovery"             About a minute ago   Up About a minute   0.0.0.0:8090->8090/tcp                                                                           fogflow_discovery_1
	eb31cd255fde   rabbitmq:3              "docker-entrypoint.s…"   About a minute ago   Up About a minute   4369/tcp, 5671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp                           fogflow_rabbitmq_1

.. important:: 

	if you see any container is missing, you can run "docker ps -a" to check if any FogFlow component is terminated with some 
	problem. If there is, you can further check its output log by running "docker logs [container ID]"


- Check the system status from the FogFlow DashBoard

You can open the FogFlow dashboard in your web browser to see the current system status via the URL: http://<my_hostip>/index.html


.. important:: 

	If the FogFlow cloud node is behind a gateway, you need to create a mapping from the gateway IP to the my_hostip and then 
	access the FogFlow dashboard via the gateway IP;
	If the FogFlow cloud node is a VM in a public cloud like Azure Cloud, Google Cloud, or Amazon Cloud, you need to access the 
	FogFlow dashboard via the public IP of your VM;
	

Once you are able to access the FogFlow dashboard, you can see the following web page


.. figure:: figures/dashboard.png


Hello World Example - NGSI-LD
===========================================================

Once the FogFlow cloud node is set up, you can try out some existing IoT services without running any FogFlow edge node.
For example, you can try out a simple fog function as below.  


Initialize all defined services with three clicks
-------------------------------------------------------------

.. important:: 

	Make sure to have correctly configured the following parameter in the config.json: `"doNotInitApplications": false`


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


Send an NGSI-LD entity to FogFlow
-------------------------------------------------------------

 
Send a curl request to the FogFlow broker for entity update:

.. code-block:: console    

	
	curl --location 'http://localhost:8070/ngsi-ld/v1/entities' \
	--header 'Content-Type: application/ld+json' \
	--data-raw '{
	"id": "house2:smartrooms:Temperature:temp006",
	"type": "Temperature",
	"temperature": {
			"value": 23,
			"unitCode": "CEL",
			"type": "Property",
			"providedBy": {
					"type": "Relationship",
					"object": "smartbuilding:house2:sensor0815"
			}
	},
	"isPartOf": {
			"type": "Relationship",
			"object": "smartcity:houses:house2"
	},
	"@context": [
			{"Room": "urn:mytypes:room", "temperature": "myuniqueuri:temperature", "isPartOf": "myuniqueuri:isPartOf"},
			"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"]
	}'


Check if the fog function is triggered
-------------------------------------------------------------

With the initialization process, a fogfunction is registered to wait for "Temperature" entities. As soon a first occurrence of such entity appear in FogFlow, the task is started.
Since we have push such entity with the previous NGSI-LD request, the task is started. Check if a task is created under "Task" in System Management.**

.. figure:: figures/fog-function-task-running.png

This dummy task is simply taking as input the pushed entity, change the id to "Result+<inputEntityId>" and send it back to FogFlow. Thus, we can check if everything worked coorectly by looking at the entities:
System Management.**

.. figure:: figures/fog-function-entities-dummyresult.png


Connect a NGSI-LD context broker to FogFlow
===========================================================

FogFlow can be easily integrated in the FIWARE ecosystem by establishing data streams from and to a NGSI-LD context broker.
Let's first connect FogFlow to get the data from a NGSI-LD context broker. The idea is issue a NGSI-LD subscription to the context broker with a reference back to the FogFlow broker as the following image:

.. figure:: figures/ngsildbroker2fogflow.png

In order to do so, send the following request:

.. code-block:: console  
	curl --silent --output /dev/null  --location "http://<SCORPIO_HOST>:9090/ngsi-ld/v1/subscriptions" \
        --header 'Content-Type: application/ld+json' \
        --data-raw '{ 
                    "type": "Subscription",
                    "entities": [{
                            "type": "Temperature"
                    }],
                    "notification": {
                            "endpoint": {
                                    "uri": "http://<FOGFLOW_BROKER>:8070/ngsi-ld/v1/notifyContext",
                                    "accept": "application/json"
                            }
                    },
                    "notificationTrigger" : ["entityCreated", "entityUpdated"] ,
                    "@context": ["https://pastebin.com/raw/hFbLejdG"]
            }''

At this point every notification will flow inside FogFlow and to respective register fogfunctions.

In order to connect FogFlow to forward data to the NGSI-LD context broker, we need to do a similar setting as the following picture:

.. figure:: figures/fogflow2ngsildbroker.png

with a request that looks like the following (note that to set the destination as NGSI-LD we use a header):

.. code-block:: console  
	curl --location 'http://<FOGFLOW_BROKER>:8070/ngsi10/subscribeContext' \
	--header 'Destination: NGSI-LD' \
	--header 'Fiware-Correlator: http://<SCORPIO_HOST>:9090/' \
	--header 'Content-Type: text/plain' \
	--data '{
	"entities": [
		{
		"type": "Temperature",
		"isPattern": true
		}
	],
	"reference": "http://<SCORPIO_HOST>:9090/"
	}'


Hello World Example - NGSI-v1
===========================================================

Once the FogFlow cloud node is set up, you can try out some existing IoT services without running any FogFlow edge node.
For example, you can try out a simple fog function as below.  


Initialize all defined services with three clicks
-------------------------------------------------------------

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


Simulate an IoT device to trigger the Fog Function
-------------------------------------------------------------

There are two ways to trigger the fog function:

**1. Create a “Temperature” sensor entity via the FogFlow dashboard**


You can register a device entity via the device registration page: "System Status" -> "Device" -> "Add". 
Then you can create a “Temperature” sensor entity by filling the following element:
- **Device ID:** to specify a unique entity ID
- **Device Type:** use “Temperature” as the entity type
- **Location:** select a location on the map
 

.. figure:: figures/device-registration.png

**2. Send an NGSI entity update to create the “Temperature” sensor entity**
 
Send a curl request to the FogFlow broker for entity update:

.. code-block:: console    

	
	curl -iX POST \
		  'http://my_hostip/ngsi10/updateContext' \
		  -H 'Content-Type: application/json' \
		  -d '
		{
		    "contextElements": [
		        {
		            "entityId": {
		                "id": "Device.Temp001",
		                "type": "Temperature",
		                "isPattern": false
		                },
		            "attributes": [
		                    {
		                    "name": "temperature",
		                    "type": "float",
		                    "value": 73
		                    },
		                    {
		                    "name": "pressure",
		                    "type": "float",
		                    "value": 44
		                    }
		                ],
		            "domainMetadata": [
		                    {
		                    "name": "location",
		                    "type": "point",
		                    "value": {
		                    "latitude": -33.1,
		                    "longitude": -1.1
		                    }}
		                ]
		        }
		    ],
		    "updateAction": "UPDATE"
		}'


Check if the fog function is triggered
-------------------------------------------------------------

Check if a task is created under "Task" in System Management.**

.. figure:: figures/fog-function-task-running.png

Check if a Stream is created under "Stream" in System Management.**

.. figure:: figures/fog-function-streams.png




