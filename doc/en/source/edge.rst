Start an edge node
==========================

Install Docker Engine 
------------------------

To install docker engine, please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

.. note:: Docker engine must be installed on each fog node, because all task instances in FogFlow will be instantiated within a docker container


Typically, a fog node needs to deploy a Worker and an IoT broker. 
The IoT Broker at the fog node can establish the data flows between all task instances launched on the same edge node. 
However, this Edge IoT Broker is optional, 
especially when the fog node is very light-weight and can only support a few tasks without any data dependency. 


Run both Worker and IoT Broker at the fog node
-------------------------------------------------

	- start dockerized Worker and IoT Broker by using a docker compose script
	
		#. before starting the script (fogflow/deployment/fog/docker-compose.yml), 
		 	please change its configuration accordingly, as shown by the following picture. 
		  	The main required change is the IP address of FogFlow core services in the cloud, 
		  	such as IoT Discovery, RabbitMQ

			.. figure:: figures/fog-node.png
 			  :scale: 100 %
 			  :alt: map to buried treasure
	
		#. start the docker-compose.yml script to launch Worker and IoT Broker(Edge)
		
			.. code-block:: bash
			
				cd fogflow/deployment/fog 
  				docker-compose up
	
	- start native Worker and IoT Broker on the fog node
	
		#. if you have not compiled them from the source code, 
			you need to download the executable files, currently ARM and Linux based version are provided
		 	for a ARM device like raspberry pi, please use fog-arm6.zip; 
			for any x86 Linux machine, please use fog-linux64.zip
			
			.. code-block:: bash
			
				// for ARM-based fog node
				wget https://github.com/smartfog/fogflow/blob/master/deployment/fog/download/arm6/fog-arm6.zip
				
				// for Linux-based fog node (64bits, x86 processor architecture)
				wget https://github.com/smartfog/fogflow/blob/master/deployment/fog/download/linux64/fog-linux64.zip
				
			
		#. unzip the download zip file
		
			
			.. code-block:: bash
			
				unzip  fog-x.zip	
	
		#. change the configuration file of IoT Broker(Edge): 
		
			- "host": to be the IP address of the fog node 
			- "discoveryURL": change it to the accessible IP address of IoT Discovery in the cloud
			- physical_location: set the geo-location of the fog node
	
		#. start IoT Broker(Edge)
		
			.. code-block:: bash
			
				cd fog-arm6/broker 
				./broker

		#. change the configuration file of Worker: 
		
			- "my_ip": to be the IP address of the fog node 
			- "message_bus": to be the HOST_IP address of the RabbitMQ in the cloud
			- "iot_discovery_url": change it to the accessible IP address of IoT Discovery in the cloud
			- physical_location: set the geo-location of the fog node


		#. start Worker(Edge)
		
			.. code-block:: bash
			
				cd fog-arm6/worker
				./worker
		

Run only Worker at the light-weight fog node
-------------------------------------------------

		#.  if you have not compiled them from the source code, 
			you need to download the executable files, currently ARM and Linux based version are provided
		 	for a ARM device like raspberry pi, please use fog-arm6.zip; 
			for any x86 Linux machine, please use fog-linux64.zip
			
			.. code-block:: bash
			
				wget
				
			
		#. unzip the download zip file
		
			
			.. code-block:: bash
			
				unzip  fog-x.zip	
		
	
		#. change the configuration file of Worker: 
		
			- "my_ip": to be the IP address of the fog node 
			- "message_bus": to be the HOST_IP address of the RabbitMQ in the cloud
			- "iot_discovery_url": change it to the accessible IP address of IoT Discovery in the cloud
			- physical_location: set the geo-location of the fog node


		#. start Worker(Edge)
		
			.. code-block:: bash
			
				cd fog-arm6/worker
				./worker