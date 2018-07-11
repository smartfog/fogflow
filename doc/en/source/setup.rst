*****************************************
Set up everything on a single machine
*****************************************


To check the basic features of FogFlow, you can set up the entire FogFlow system on a single Linux machine (e.g., Ubuntu 16.04.4 LTS). 
If you already have Docker-CE and Docker Compose installed on your machine, 
the setup can be quickly finished in just a few minutes. 

Here are the steps to follow: 


Install Docker CE and Docker Compose on your Linux machine
===============================================================

    .. important::
    
        To install Docker CE, please refer to |install_docker|, required version 18.03.1-ce, *please also allow your user to execute the Docker Command without Sudo*

            .. |install_docker| raw:: html

                <a href="https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04" target="_blank">How to install Docker</a>

        To install Docker Compose, please refer to |install_docker_compose|, required version 2.4.2

            .. |install_docker_compose| raw:: html

                <a href="https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04" target="_blank">How to install Docker Compose</a>

Download the deployment script and the configuration file
===============================================================

    .. parsed-literal::
         
          #download the deployment script
          wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/core/docker-compose.yml
          
          #download the configuration file          
          wget https://raw.githubusercontent.com/smartfog/fogflow/master/deployment/core/config.json


Change the configuration file according to your local environment
====================================================================

    As a simple change, you only need to replace the IP address with your own host IP in three places, as illustrated in the following figure. 

    .. important:: **HOST_IP** is the IP address of your Linux machine. 
    We also assume that you can use the default port numbers for various FogFlow components. 
    More specially, the following ports are required.    
        - 80: for Designer to provide the FogFlow web portal
        - 443: for Discovery
        - 8080: for Broker   
        - 5672: for RabbitMQ 
  
    .. figure:: figures/configuration.png
       :scale: 100 %


Run the downloaded script
===============================================================

     .. parsed-literal::

          #pull the required docker images and create their containers, only required for the first time
          docker-compose create
          
          #start the FogFlow system 
          docker-compose start

          #command to stop the FogFlow system
          docker-compose stop  #no need for the following steps


Test the FogFlow dashboard
===============================================================

    Open the link "http://HOST_IP" in your browser to check the status of all FogFlow running components in the cloud. 

    If everything goes well, you should be able to see the following page from this link. 

    .. figure:: figures/designer.png
       :scale: 100 %

    Furthermore, you should be able to see the status of all core components running in the cloud, 
    from the menu items on the left side of the System Management page. 

    .. figure:: figures/status.png
       :scale: 100 %












