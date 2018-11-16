Start an edge node
==========================

Typically, an FogFlow edge node needs to deploy a Worker and an IoT broker. 
The Edge IoT Broker at the edge node can establish the data flows between all task instances launched on the same edge node. 
However, this Edge IoT Broker is optional, 
especially when the edge node is a very constrained device that can only support a few tasks without any data dependency. 

Here are the steps to start an FogFlow edge node: 

Install Docker Engine 
------------------------

To install Docker CE, please refer to |install_docker|, required version 18.03.1-ce, *please also allow your user to execute the Docker Command without Sudo*
    .. |install_docker| raw:: html
        <a href="https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04" target="_blank">How to install Docker</a>


.. note:: Docker engine must be installed on each edge node, because all task instances in FogFlow will be launched within a docker container.


Download the deployment script 
-------------------------------------------------

.. code-block:: console    
         
	#download the deployment scripts
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/start.sh
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/stop.sh 
	
	#make them executable
	chmod +x start.sh  stop.sh       


Download the default configuration file 
-------------------------------------------------

.. code-block:: console   
         	
	#download the configuration file          
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/config.json


Change the configuration file accordingly
-------------------------------------------------

You can use the default setting for a simple test, but you need to change the following addresses according to your own environment: 
        
- **coreservice_ip**: please refer to the configuration of the cloud part. This is the accessible address of your FogFlow core services running in the cloud node;
- **external_hostip**: this is the external IP address, accessible for the cloud broker. It is useful when your edge node is behind NAT;
- **internal_hostip** is the IP of your default docker bridge, which is the "docker0" network interface on your host. 

.. code-block:: json

    //you can see the following part in the default configuration file
    { 
        "coreservice_ip": "155.54.239.141", 
        "external_hostip": "35.234.116.177", 
        "internal_hostip": "172.17.0.1", 
        …
    } 


Start both Edge IoT Broker and FogFlow Worker
-------------------------------------------------

.. note:: if the edge node is ARM-basd, please attach arm as the command parameter

.. code-block:: console    

      #start both components in the same script
      ./start.sh 
    
      # if the edge node is ARM-basd, please attach arm as the command parameter
      #./start.sh  arm
      

Stop both Edge IoT Broker and FogFlow Worker
-------------------------------------------------

.. code-block:: console    

	#stop both components in the same script
	./stop.sh 


        
