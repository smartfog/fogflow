*****************************************************************
Security
*****************************************************************

Secure the cloud-edge communication
===============================================
To secure the communication between the FogFlow cloud node and the FogFlow edge nodes, 
FogFlow can be configured to use HTTPs for the NGSI9 and NGSI10 communication, 
which is mainly for data exchange between cloud node and edge nodes, or between two edge nodes. 
Also, the control channel between Topology Master and Worker can be secured by enabling TLS in RabbitMQ. 
The introduction steps to secure the data exchange between one FogFlow cloud node and one FogFlow edge node. 



Configure DNS server
===========================================================

As illustrated by the following picture, in order to set up FogFlow to support the HTTPs-based communication, 
the FogFlow cloud node and the FogFlow edge node are required to have their own domain names, 
because their signed certificates must be associated with their domain namers.
Therefore, DNS service is needed to be used to resolve the domain names for both the cloud node and the edge node. 
For example, `freeDNS`_ can be used for this purpose. 

.. _`freeDNS`: https://freedns.afraid.org


.. figure:: figures/https-setup.png


.. important:: 

	please make sure that the domain names of the cloud node and the edge node can be properly resolved
	and correct IP address can be seen.  
	

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
-----------------------------------------------

.. code-block:: console    

	docker-compose up -d 


Validate setup
--------------------------------------------

.. code-block:: console    

    docker ps 

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
------------------------------------------------

.. code-block:: console    

	docker-compose up -d 


Validate setup
--------------------------------------------

.. code-block:: console    

	docker ps 

	CONTAINER ID        IMAGE               COMMAND             CREATED              STATUS              PORTS                                      NAMES
	16af186fb54e        fogflow/worker      "/worker"           About a minute ago   Up About a minute                                              https_edge_worker_1
	195bb8e44f5b        fogflow/broker      "/broker"           About a minute ago   Up About a minute   0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   https_edge_broker_1
	


Check system status via FogFlow Dashboard
===========================================================

FogFlow dashboard can be opened in web browser to see the current system status via the URL: https://cloudnode.fogflow.io/index.html

.. important:: 

	please make sure that the domain names of the cloud node can be properly resolved. 

	If self-signed SSL certificate is being used, a browser warning indication can be seen that the crtificate should not be trusted.
	It can be proceeded past this warning to view the FogFlow dashboard web page via https.
