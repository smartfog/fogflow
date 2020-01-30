*****************************************************************
Secure the cloud-edge communication
*****************************************************************

To secure the communication between the FogFlow cloud node and the FogFlow edge nodes, 
FogFlow can be configured to use HTTPs for the NGSI9 and NGSI10 communication, 
which is mainly for data exchange between cloud node and edge nodes, or between two edge nodes. 
Also, the control channel between Topology Master and Worker can be secured by enabling TLS in RabbitMQ. 
Here we introduce the steps to secure the data exchange between one FogFlow cloud node and one FogFlow edge node. 



Configure your DNS server
===========================================================

As illustrated by the following picture, in order to set up FogFlow to support the HTTPs-based communication, 
the FogFlow cloud node and the FogFlow edge node are required to have their own domain names, 
because their signed certificates must be associated with their domain namers.
Therefore, you need to use a DNS service to resolve the domain names for both the cloud node and the edge node. 
For example, you can use `freeDNS`_ for this purpose. 

.. _`freeDNS`: https://freedns.afraid.org


.. figure:: figures/https-setup.png


.. important:: 

	please make sure that the domain names of the cloud node and the edge node can be properly resolved
	and you can see the correct IP address.  
	

Set up the FogFlow cloud node
===========================================================

Fetch all required scripts
--------------------------------------------

Download the docker-compose file and the configuration files as below.

.. code-block:: console    

	# download the script that can fetch all required files
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/fetch.sh
	
	# make this script executable
	chmod +x fetch.sh

	# run this script to fetch all required files
	./fetch.sh



Change the configuration file
--------------------------------------------

.. code-block:: console    
	
	{
	    "coreservice_ip": "cloudnode.fogflow.io",   #change this to the domain name of your own cloud node 
	    "external_hostip": "cloudnode.fogflow.io",  #change this to the domain name of your own cloud node 
		...
	}

Generate the key and certificate files
--------------------------------------------

.. code-block:: console    

	# make this script executable
	chmod +x key4cloudnode.sh

	# run this script to fetch all required files
	./key4cloudnode.sh  cloudnode.fogflow.io


Start the FogFlow components on the cloud node
--------------------------------------------

.. code-block:: console    

	docker-compose up -d 


Validate your setup
--------------------------------------------

.. code-block:: console    

    docker ps 

	CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                                   NAMES
	171fb42a0cd8        nginx               "nginx -g 'daemon of…"   6 seconds ago       Up 2 seconds        80/tcp, 0.0.0.0:443->443/tcp                            https_nginx_1
	739f31e8bc23        fogflow/master      "/master"                8 seconds ago       Up 3 seconds        0.0.0.0:1060->1060/tcp                                  https_master_1
	da2ebd3ae351        fogflow/worker      "/worker"                8 seconds ago       Up 4 seconds                                                                https_cloud_worker_1
	ea475cc8d696        fogflow/designer    "node main.js"           8 seconds ago       Up 5 seconds        0.0.0.0:1030->1030/tcp, 0.0.0.0:8080->8080/tcp          https_designer_1
	d35e00371bdb        fogflow/broker      "/broker"                10 seconds ago      Up 7 seconds        0.0.0.0:8070->8070/tcp, 0.0.0.0:8072->8072/tcp          https_cloud_broker_1
	c06da5d41e65        rabbitmq:3          "docker-entrypoint.s…"   12 seconds ago      Up 9 seconds        4369/tcp, 5671/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp   https_rabbitmq_1
	79c1464fa6ff        fogflow/discovery   "/discovery"             12 seconds ago      Up 10 seconds       0.0.0.0:8090->8090/tcp, 0.0.0.0:8092->8092/tcp          https_discovery_1
	


Set up the FogFlow edge node
===========================================================


Fetch all required scripts
--------------------------------------------

Download the docker-compose file and the configuration files as below.

.. code-block:: console    

	# download the script that can fetch all required files
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/https/fetch.sh
	
	# make this script executable
	chmod +x fetch.sh

	# run this script to fetch all required files
	./fetch.sh



Change the configuration file
--------------------------------------------

.. code-block:: console    
	
	{
	    "coreservice_ip": "cloudnode.fogflow.io",   #change this to the domain name of your own cloud node 
	    "external_hostip": "edgenode1.fogflow.io",  #change this to the domain name of your own edge node 
		...
	}


Generate the key and certificate files
--------------------------------------------

.. code-block:: console    

	# make this script executable
	chmod +x key4edgenode.sh

	# run this script to fetch all required files
	./key4edgenode.sh  edgenode1.fogflow.io


Start the FogFlow components on the edge node
--------------------------------------------

.. code-block:: console    

	docker-compose up -d 


Validate your setup
--------------------------------------------

.. code-block:: console    

	docker ps 

	CONTAINER ID        IMAGE               COMMAND             CREATED              STATUS              PORTS                                      NAMES
	16af186fb54e        fogflow/worker      "/worker"           About a minute ago   Up About a minute                                              https_edge_worker_1
	195bb8e44f5b        fogflow/broker      "/broker"           About a minute ago   Up About a minute   0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   https_edge_broker_1
	


Check system status via FogFlow Dashboard
===========================================================

You can open the FogFlow dashboard in your web browser to see the current system status via the URL: https://cloudnode.fogflow.io/index.html

.. important:: 

	please make sure that the domain names of the cloud node can be properly resolved. 
	If you use a self-signed SSL certificate, you will see a browser warning indicating that the certificate should not be trusted.
	You can proceed past this warning to view the FogFlow dashboard web page via https.