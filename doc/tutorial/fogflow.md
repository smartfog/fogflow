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

for simplicity, we use a single docker compose file to launch all necessary FogFlow components, including the core FogFlow services and also one edge node. The edge node means one FogFlow Broker and one FogFlow Worker and they are deployed on a physical edge node, such as an IoT Gateway or a raspberry Pi. 

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

You may follow the orion docs to set up a Orion Context Broker instance from here: [Install Orion](https://fiware-orion.readthedocs.io/en/master/admin/install/index.html)

You may also setup Orion on docker using below commands.(docker is required for this method)




# Programe a simple fog function via FogFlow Dashboard


# Simulate an IoT device

take the powerpanel sensor as an example and then write a simple temperature sensor, which can publish the temperature entity every 5 seconds. 


# Check if the fog function is triggered


# Check if the fog function is triggered


# Issue a subscription to forward the generated result to Orion Context Broker


# Query the result from Orion Context Broker







