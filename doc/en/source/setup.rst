*****************************************
Set up everything on a single machine
*****************************************


To check the basic features of FogFlow, you can set up the entire FogFlow system on a single Linux machine. 
If you already have docker and docker-composer installed on your Linux machine, 
the setup can be quickly finished in a few minutes. 

Here are the steps: 


#. install docker and docker-composer on your linux machine

	To install Docker CE, please refer to |install_docker|

		.. |install_docker| raw:: html

			<a href="https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04" target="_blank">Docker</a>

	To install Docker Compose, please refer to |install_docker_compose|

		.. |install_docker_compose| raw:: html

			<a href="https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04" target="_blank">Docker Compose</a>


#. download the deployment script and the configuration file

	.. code-block:: bash
		
		#download the deployment script
  		wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/core/docker-compose.yml
		
		#download the configuration file		
		wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/core/config.json


#. change the configuration file

	.. code-block:: bash
	
		# set the environment variable HOST_IP, which is the external IP address of your current machine
		export HOST_IP=AAA.BBB.CCC.DDD

		# go to the folder where the docker-compose.ymal is located
		cd fogflow/deployment/core 
  		docker-compose up

#. run the downloaded script

	.. code-block:: bash
		
		#start the FogFlow system 
  		docker-compose start


#. test the FogFlow dashboard

Open the link http://HOST_IP:8080 in your browser to check the status of all FogFlow running components in the cloud. 

If everything goes well, you can see the following page from this link. 

.. figure:: figures/designer.png
   :scale: 100 %
   :alt: map to buried treasure

Further more, you should be able to see the status of all core components running in the cloud, 
from the menu items on the left side of the System Management page. 

.. figure:: figures/status.png
   :scale: 100 %
   :alt: map to buried treasure












