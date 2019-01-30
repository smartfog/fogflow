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

```bash
git clone git@github.com:

```

# Start Up Orion

follow the orion tutorial to set up a Orion Context Broker instance


# Programe a simple fog function via FogFlow Dashboard


# Simulate an IoT device

take the powerpanel sensor as an example and then write a simple temperature sensor, which can publish the temperature entity every 5 seconds. 


# Check if the fog function is triggered


# Check if the fog function is triggered


# Issue a subscription to forward the generated result to Orion Context Broker


# Query the result from Orion Context Broker







