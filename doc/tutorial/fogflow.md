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

For simplicity, we use a single docker compose file to launch all necessary FogFlow components, including the core FogFlow services and also one edge node. The edge node means one FogFlow Broker and one FogFlow Worker and they are deployed on a physical edge node, such as an IoT Gateway or a raspberry Pi. 

**Prerequisite:** (For both Cloud and Edge node)

These two commands should be present on your system.
1. docker
2. docker-compose

For ubuntu-16.04, you need to install docker-ce and docker-compose.

To install Docker CE, please refer to [Install Docker CE](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04), required version 18.03.1-ce;
*please also allow your user to execute the Docker Command without Sudo*

To install Docker Compose, please refer to [Install Docker Compose](https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04), required version 18.03.1-ce, required version 2.4.2

**Setup Fogflow Cloud Node:**

Download the docker-compose file and the config.json file to setup flogflow cloud node

```bash
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/docker-compose.yml

wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/config.json
```
you need to change the following addresses in config.json according to your own environment.

-  **webportal_ip**: this is the IP address to access the FogFlow web portal provided by Task Designer. It must be accessible from outside by user's browser.  
-  **coreservice_ip**: it is used by all edge nodes to access the FogFlow core services, including Discovery, Broker(Cloud), and RabbitMQ;
-  **external_hostip**: this is the same as coreservice_ip, for the cloud part of FogFlow;        
-  **internal_hostip**: is the IP of your default docker bridge, which is the "docker0" network interface on your host

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

**Setup Fogflow Edge Node:** (optional for basic setup)

An FogFlow edge node needs to deploy a Worker and an IoT broker. An Edge node can be on physical device or raspberry Pi. It can be setup in cloud/VM for testing/simulation environment. 

Download shell scripts to start/stop edge node.

Note: Docker and docker-compose should be installed before running these script as specified in Prerequisite section.

```bash
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/start.sh
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/stop.sh
```

Download the configuration file for Edge node and change it accordingly.
```bash
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/edge/config.json
```
You need to change the following addresses according to your own environment in config.json file: 
        
- **coreservice_ip**: please refer to the configuration of the cloud part. This is the accessible address of your FogFlow core services running in the cloud node;
- **external_hostip**: this is the external IP address, accessible for the cloud broker. It is useful when your edge node is behind NAT;
- **internal_hostip** is the IP of your default docker bridge, which is the "docker0" network interface on your host.

To start Edge node components(Broker and worker), run the script 'start.sh'.
```bash
./start.sh
```

To stop Edge node components(Broker and worker), run the script 'stop.sh'
```bash
./stop.sh
```

if the edge node is ARM-basd, please attach arm as the command parameter
```bash
./start.sh  arm
```

# Start Up Orion

You may follow the orion docs to set up a Orion Context Broker instance from here: [Installing Orion](https://fiware-orion.readthedocs.io/en/master/admin/install/index.html)

You may also setup Orion on docker using below commands.(docker and docker-compose are required for these method)
Note: Orion container has a dependency on MongoDB database.

Their are two methods to setup Orion using Docker:

## 1. The Fastest Way Using Docker-compose:

Docker Compose allows you to link an Orion Context Broker container to a MongoDB container in a few minutes. This method requires that install [Docker Compose](https://docs.docker.com/compose/install/)

Follow these steps:
1. Create a directory on your system on which to work( for example, ```bash ~/fiware ```).
2. Create a new file called docker-compose.yaml inside your directory with the following contents:

```bash
mongo:
  image: mongo:3.4
  command: --nojournal
orion:
  image: fiware/orion
  links:
    - mongo
  ports:
    - "1026:1026"
  command: -dbhost mongo
```
3. Using the command-line and within directory you created type:

```bash
sudo docker-compose up
```
Note: Regarding --nojournal it is not recommended for production, but it speeds up mongo container startup and avoids some race    conditions problems if Orion container is faster and doesn't find the DB up and ready.

After a few seconds you should have your Context Broker running and listening on port 1026.

Check that everything works with
```bash
curl localhost:1026/version
```
## 2. Setup Orion using ``` docker run ``` command:

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
## 3. Install Orion If MongoDB already exists on different host:

If you want to connect to a different MongoDB instance do the following commnad
```bash
sudo docker run -d --name orion1 -p 1026:1026 fiware/orion -dbhost <MongoDB Host> 
```

Check that everything works with
```bash
curl localhost:1026/version
```

# Program a simple fog function via FogFlow Dashboard
**1. Creating a FogFunction**
- Click on FogFuntion from the top menu on Fogflow portal.
- Select Editor to start creating the Fog Function. Editor consists of a graphical editor and a text (or code) editor.

![Creating a FogFunction](https://fogflow.readthedocs.io/en/latest/_images/fog-function-menu.png)

**2. Adding FogFunction element**
- Right click in graphical editor area and select "FogFunction".
- Right click in graphical editor area and select "InputTrigger". Use connectors Selector-Selectors to connect FogFunction and InputTrigger elements.
- Right click in graphical editor area and select "SelectCondition". Use connectors Condition-Conditions to connect InputTrigger and SelectCondition elements.

![Adding FogFunction element](https://fogflow.readthedocs.io/en/latest/_images/fog-function-selected.png)

**3. Configuring FogFunction element**
- Click on Configure button of "FogFunction" element on its top right corner and provide FogFunction name and user.

![Configuring FogFunction element](https://fogflow.readthedocs.io/en/latest/_images/fog-function-configuration.png)

**4. Configuring InputTrigger element**
- Click on Configure button of "InputTrigger" element on its top right corner and choose SelectedAttributes and GroupBy values.

![Configuring SelectCondition element](https://fogflow.readthedocs.io/en/latest/_images/fog-function-filter.png)

**5. Configuring SelectCondition element**
- Click on Configure button of "SelectCondition" element on its top right corner and choose the condition for triggering the FogFuntion.

![Configuring InputTrigger element](https://fogflow.readthedocs.io/en/latest/_images/fog-function-granularity.png)

**6. Customizing FogFuntion Code**
- Edit the FogFunction to customize the logic of FogFunction.

![Customizing FogFuntion Code](https://fogflow.readthedocs.io/en/latest/_images/fog-function-code.png)

Currently FogFlow allows developers to specify the function code, either by directly overwritting the following handler function (in Javascript or Python) or by selecting a registered operator (or docker image).

```bash
exports.handler = function(contextEntity, publish, query, subscribe) {
    console.log("enter into the user-defined fog function");

    var entityID = contextEntity.entityId.id;

    if (contextEntity == null) {
        return;
    }
    if (contextEntity.attributes == null) {
        return;
    }

    var updateEntity = {};
    updateEntity.entityId = {
        id: "Stream.result." + entityID,
        type: 'result',
        isPattern: false
    };
    updateEntity.attributes = {};
    updateEntity.attributes.city = {
        type: 'string',
        value: 'Heidelberg'
    };

    updateEntity.metadata = {};
    updateEntity.metadata.location = {
        type: 'point',
        value: {
            'latitude': 33.0,
            'longitude': -1.0
        }
    };

    console.log("publish: ", updateEntity);
    publish(updateEntity);
};
```

You can take the example Javascript code above as the implementation of your “HelloWorld” fog function. This example fog function simply writes a fixed entity by calling the “publish” callback function.

The input parameters of a fog function are predefined and fixed, including:
**- contextEntity:** representing the received entity data
**- publish:** the callback function to publish your generated result back to the FogFlow system
**- query:** optional, this is used only when your own internal function logic needs to query some extra entity data from the FogFlow context management system.
**- subscribe:** optional, this is used only when your own internal function logic needs to subscribe some extra entity data from the FogFlow context management system.

**example usage of publish:**
```bash
var updateEntity = {};
updateEntity.entityId = {
       id: "Stream.Temperature.0001",
       type: 'Temperature',
       isPattern: false
};
updateEntity.attributes = {};
updateEntity.attributes.city = {type: 'string', value: 'Heidelberg'};

updateEntity.metadata = {};
updateEntity.metadata.location = {
    type: 'point',
    value: {'latitude': 33.0, 'longitude': -1.0}
};

publish(updateEntity);
```

**example usage of query:**
```bash
var queryReq = {}
queryReq.entities = [{type:'Temperature', isPattern: true}];
var handleQueryResult = function(entityList) {
    for(var i=0; i<entityList.length; i++) {
        var entity = entityList[i];
        console.log(entity);
    }
}

query(queryReq, handleQueryResult);
```

**example usage of subscribe:**
```bash
var subscribeCtxReq = {};
subscribeCtxReq.entities = [{type: 'Temperature', isPattern: true}];
subscribeCtxReq.attributes = ['avg'];

subscribe(subscribeCtxReq);
```

**7. Submitting FogFuntion**
- Click on "Create a Fog Function" button once customization is complete. You can also edit it again by selecting the Fog Function from the list of registered Fog Functions.

![Submitting FogFuntion](https://fogflow.readthedocs.io/en/latest/_images/fog-function-submit.png)


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


# Check if the fog function is triggered


# Issue a subscription to forward the generated result to Orion Context Broker


# Query the result from Orion Context Broker







