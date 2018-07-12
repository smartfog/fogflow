Start an edge node
==========================

Typically, an FogFlow edge node needs to deploy a Worker and an IoT broker. 
The Edge IoT Broker at the edge node can establish the data flows between all task instances launched on the same edge node. 
However, this Edge IoT Broker is optional, 
especially when the edge node is a very constrained device that can only support a few tasks without any data dependency. 

Here are the steps to start an FogFlow edge node: 

Install Docker Engine 
------------------------

.. note:: Docker engine must be installed on each edge node, because all task instances in FogFlow will be launched within a docker container.

    To install Docker CE, please refer to |install_docker|, required version 18.03.1-ce, *please also allow your user to execute the Docker Command without Sudo*

      .. |install_docker| raw:: html

         <a href="https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04" target="_blank">How to install Docker</a>


Download the deployment script 
-------------------------------------------------

    .. parsed-literal::
         
          #download the deployment scripts
          wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/edge/start.sh
          wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/edge/stop.sh 
          
          #make them executable
          chmod +x start.sh  stop.sh       
          

Download the default configuration file 
-------------------------------------------------

    .. parsed-literal::
         
         
          #download the configuration file          
          wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/edge/config.json


Change the configuration file accordingly
-------------------------------------------------

The following picture shows the configurations you need to chane accordingly to your own environment. 

    .. figure:: figures/edgecfg.png
       :scale: 100 %

Start both Edge IoT Broker and FogFlow Worker
-------------------------------------------------

    .. note:: if the edge node is ARM-basd, please attach arm as the command parameter

    .. parsed-literal::

          #start both components in the same script
          ./start.sh 
        
          # if the edge node is ARM-basd, please attach arm as the command parameter
          #./start.sh  arm
          


Stop both Edge IoT Broker and FogFlow Worker
-------------------------------------------------


     .. parsed-literal::

          #stop both components in the same script
          ./stop.sh 


        
