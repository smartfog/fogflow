[![FIWARE FogFlow](https://nexus.lab.fiware.org/repository/raw/public/badges/chapters/processing.svg)](https://www.fiware.org/developers/catalogue/)
[![FogFlow 2.0](https://img.shields.io/badge/FogFlow-2.0.svg)](http://fogflow.readthedocs.io/)

**Description:** 
This is an Introductory Tutorial to [FogFlow](https://eprosima-fast-rtps.readthedocs.io).
In the FIWARE-based architecture, FogFlow can be used to dynamically trigger data processing functions 
between IoT devices and Orion Context Broker, for the purpose of transforming and preprocessing raw data at edge nodes, 
which can be deployed at IoT gateways or directly on the devices like Raspberry Pis.

The tutorial introduces a typical FogFlow system setup with a simple example to do anomaly detection at edges for temperature sensor data. 

---

# What is FogFlow?

[FogFlow](https://fogflow.readthedocs.io) is an IoT edge computing framework 
to automatically orchestrate dynamic data processing flows over cloud and edges driven by context, 
including system context on the available system resources from all layers, 
data context on the registered metadata of all available data entities, 
and also usage context on the expected QoS defined by users.

---

# System View

describe the entire scenario and also use a figure to show the system view

# Start Up FogFlow

Here are the prerequisite commands for starting Fogflow:
1. docker
2. docker-compose

For ubuntu-16.04, you need to install docker-ce and docker-compose.

To install Docker CE, please refer to [Install Docker CE](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04), required version 18.03.1-ce;
*please also allow your user to execute the Docker Command without Sudo*

To install Docker Compose, please refer to [Install Docker Compose](https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04), required version 18.03.1-ce, required version 2.4.2

**Setup Fogflow:**

Download the docker-compose file and the config.json file to setup flogflow.

```bash
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/docker-compose.yml

wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/config.json
```
you need to change the following addresses in config.json according to your own environment.

- **webportal_ip**: this is the IP address to access the FogFlow web portal provided by Task Designer. It must be accessible from outside by user's browser.  
- **coreservice_ip**: it is used by all edge nodes to access the FogFlow core services, including Discovery, Broker(Cloud), and RabbitMQ;
- **external_hostip**: this is the same as coreservice_ip, for the cloud part of FogFlow;        
- **internal_hostip**: is the IP of your default docker bridge, which is the "docker0" network interface on your host

**Firewall rules:** to make your FogFlow web portal accessible via the external_ip; the following ports must be open as well: 80, 443, 8080, and 5672 for TCP
    
Pull the docker images of all FogFlow components and start the FogFlow system
```bash
docker-compose pull

docker-compose up -d
```
Check all the containers are Up and Running using "docker ps -a"
```bash
root@fffog-ynomljrk3y7bs23:~# docker ps -a
CONTAINER ID  IMAGE              COMMAND           CREATED      STATUS       PORTS                                       NAMES
122a61ece2ce  fogflow/master     "/master"         26 hours ago Up 26 hours  0.0.0.0:1060->1060/tcp                      fogflow_master_1
e625df7a1e51  fogflow/designer   "node main.js"    26 hours ago Up 26 hours  0.0.0.0:80->80/tcp, 0.0.0.0:1030->1030/tcp  fogflow_designer_1
42ada6ee39ae  fogflow/broker     "/broker"         26 hours ago Up 26 hours  0.0.0.0:8080->8080/tcp                      fogflow_broker_1
39b166181acc  fogflow/discovery  "/discovery"      26 hours ago Up 26 hours  0.0.0.0:443->443/tcp                        fogflow_discovery_1
8951aaac0049  tutum/rabbitmq     "/run.sh"         26 hours ago Up 26 hours  0.0.0.0:5672->5672/tcp, 15672/tcp           fogflow_rabbitmq_1
7f32d441c54a  mdillon/postgis    "docker-entry…"   26 hours ago Up 26 hours  0.0.0.0:5432->5432/tcp                      fogflow_postgis_1
53bf689d3db6  fogflow/worker     "/worker"         26 hours ago Up 26 hours                                              fogflow_cloud_worker_1
root@fffog-ynomljrk3y7bs23:~# 

```

**Test the Fogflow Dashboard:**

Open the link “http://<webportal_ip>” in your browser to check the status of all FogFlow running components in the cloud.

So now you can also check all the components using dashboard.


# Start Up Orion

You may follow the orion docs to set up a Orion Context Broker instance from here: [Installing Orion](https://fiware-orion.readthedocs.io/en/master/admin/install/index.html)

You may also setup Orion on docker using below commands.(docker and docker-compose are required for these method)
Note: Orion container has a dependency on MongoDB database.

## Easiest way to setup Orion (docker is ):

**Prerequisite:** Docker should be installed.

First launch MongoDB container using below command:
```bash
sudo docker run --name mongodb -d mongo:3.4
```

And then run Orion with this command
```bash
sudo docker run -d --name orion1 --link mongodb:mongodb -p 1026:1026 fiware/orion -dbhost mongodb
```

Check that everything works with
```bash
curl localhost:1026/version
```

# Program a simple fog function via FogFlow Dashboard
**1. Creating a FogFunction**
- Click on FogFuntion from the top menu on Fogflow portal.
- Select Editor to start creating the Fog Function. Editor consists of a graphical editor and a text (or code) editor.

**2. Adding FogFunction element**
- Right click in graphical editor area and select "FogFunction".
- Right click in graphical editor area and select "InputTrigger". Use connectors Selector-Selectors to connect FogFunction and InputTrigger elements.
- Right click in graphical editor area and select "SelectCondition". Use connectors Condition-Conditions to connect InputTrigger and SelectCondition elements.

**3. Configuring FogFunction element**
- Click on Configure button of "FogFunction" element on its top right corner and provide FogFunction name and user.

**4. Configuring InputTrigger element**
- Click on Configure button of "InputTrigger" element on its top right corner and choose SelectedAttributes and GroupBy values.

**5. Configuring SelectCondition element**
- Click on Configure button of "SelectCondition" element on its top right corner and choose the condition for triggering the FogFuntion.

**6. Customizing FogFuntion Code**
- Edit the FogFunction to customize the logic of FogFunction. Currently FogFlow allows developers to specify the function code in Javascript, python and docker image.

**7. Submitting FogFuntion**
- Click on "Create a Fog Function" button once customization is complete. You can also edit it again by selecting the Fog Function from the list of registered Fog Functions.

# Simulate an IoT device

There are three inbuilt usecases in Fogflow:
-  **Anomaly Detector for retails**
-  **Lost Child finder for public safety**
-  **Smart Parking for smart cities**

You may run any of these usecases, provided the following prerquisites are fulfilled:
1. Fogflow should be installed and running well.
2. Edge Devices should be simulated and running.
Simulated devices will feed the Fogflow System with Context Data on regular basis (say 5 seconds).
Follow these steps to get the devices running:

## 1. Simulate devices:
-  Install python2, pip for python2, nodejs, and npm in order to run the simulated devices:
```bash
apt install python2.7 python-pip
curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
apt-get install -y nodejs
node -v
npm -v
```
-  Download the code repository:
```bash
git clone https://github.com/smartfog/fogflow.git
```

## 2. Run the simulated devices:
-  Start the simulated powerpanel device for "anomaly detection"
```bash
cd  fogflow/application/device/powerpanel
npm install
```            
(Note: Please change the "discoveryURL": "http://<Fogflow_Discovery_Ip>:443/ngsi9" in the following profile.json files before proceeding.)      
```bash
node powerpanel profile1.json
node powerpanel profile2.json
node powerpanel profile3.json
```
-  Start the simulated camera device for "lost child finder"
```bash
cd  fogflow/application/device/camera1
pip install –r requirements.txt
```            
(Note: Please change the "discoveryURL": "http://<Fogflow_Discovery_Ip>:443/ngsi9" in the following profile.json files before proceeding.)            
```bash
Python fakecamera.py profile.json
```          

# Check if the fog function is triggered

**There are two way to trigger the fog function**
### Step 1
**1.create a “Temperature” sensor entity by filling the following element**
 - **Device ID:** to specify a unique entity ID.
 - **Device Type:** use “Temperature” as the entity type.
 - **Location:** use “Temperature” as the entity type.
 ![Temperature](https://fogflow.readthedocs.io/en/latest/_images/device-registration.png)
 
 - Once the device profile is registered, a new “Temperature” sensor entity will be created and it will trigger the “HelloWorld” fog function automatically.
 ### Step 2
 - To trigger the “HelloWorld” fog function is to send a NGSI entity update to create the “Temperature” sensor entity.
 -Send the post request to the FogFlow broker.
 ```bash
    curl -iX POST \
  'http://localhost:8080/ngsi10/updateContext' \
  -H 'Content-Type: application/json' \
  -d '
{
    "contextElements": [
        {
            "entityId": {
                "id": "Device.temp001",
                "type": "Temperature",
                "isPattern": false
            },
            "attributes": [
            {
              "name": "temp",
              "type": "integer",
              "contextValue": 10
            }
            ],
            "domainMetadata": [
            {
                "name": "location",
                "type": "point",
                "value": {
                    "latitude": 49.406393,
                    "longitude": 8.684208
                }
            }
            ]
        }
    ],
    "updateAction": "UPDATE"
}'
 ``` 
 **You can check whether the fog function is triggered or not in the following way**
 - check the task instance of this fog function, as shown in the following picture
 ![check the task instance of this fog function, as shown in the following picture](https://fogflow.readthedocs.io/en/latest/_images/fog-function-task-instance.png)
 -check the result generated by its running task instance, as shown in the following picture
 ![check the result generated by its running task instance, as shown in the following picture](https://fogflow.readthedocs.io/en/latest/_images/fog-function-result.png)


# Check if the fog function is triggered


# Issue a subscription to forward the generated result to Orion Context Broker


# Query the result from Orion Context Broker







