[![FIWARE FogFlow](https://nexus.lab.fiware.org/repository/raw/public/badges/chapters/robotics.svg)](https://www.fiware.org/developers/catalogue/)
[![FogFlow 2.0](https://img.shields.io/badge/FogFlow-1.6-5dc0cf.svg)](http://fogflow.readthedocs.io/)

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



# Start Up

To start the installation, do the following:

```bash
git clone git@github.com:

```


```bash
cmake .
make
make install
```
