.. _cloud-setup:

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

    You can use the default setting for a simple test, but you need to change the following addresses according to your own environment: 
        
        * **coreservice_ip**: it is used by users to access the FogFlow web portal and must be accessible from your browser; Also, it is used by all edge nodes to access the FogFlow core services, including Discovery, Broker(Cloud), and RabbitMQ;
        * **external_hostip**: this is the same as coreservice_ip, for the cloud part of FogFlow;        
        * **internal_hostip** is the IP of your default docker bridge, which is the "docker0" network interface on your host. 

    .. code-block:: json
    
        //you can see the following part in the default configuration file
        { 
            "coreservice_ip": "155.54.239.141", 
            "external_hostip": "155.54.239.141", 
            "internal_hostip": "172.17.0.1", 
            …
        } 


    .. important:: 
        * **firewall rules**: to make your FogFlow web portal accessible via the external_ip; the following ports must be open as well: 80, 443, 8080, and 5672 for TCP

    
    We also assume that you can use the default port numbers for various FogFlow components. 
    More specially, the following ports are required.    
        - 80: for FogFlow web portal to be accessible at the external IP    
        - 443: for Discovery to be accessible at the external IP    
        - 8080: for Broker to be accessible at the external IP    
        - 5672: for RabbitMQ, used only internally between Master and Worker(s) 
  

Run the downloaded script
===============================================================

     .. parsed-literal::
         
          #pull the docker images of all FogFlow components
          docker-compose pull 
        
          #start the FogFlow system 
          docker-compose up -d 

Test the FogFlow dashboard
===============================================================

    Open the link "http://external_ip" in your browser to check the status of all FogFlow running components in the cloud. 

    If everything goes well, you should be able to see the following page from this link. 

    .. figure:: figures/designer.png
       :width: 100 %

    Furthermore, you should be able to see the status of all core components running in the cloud, 
    from the menu items on the left side of the System Management page. 

    .. figure:: figures/status.png
       :width: 100 %












