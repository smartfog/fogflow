# Prometheus in FogFlow
## Intro
This section is about integrating prometheus as a component in the FogFlow system. FogFlow uses the prometheus for distributed monitroing. You can use prometheus to check CPU/Memory/Network/IO usage of each container. The API of Prometheus can be used for QoS-aware orchestration.

## Basic Flow
Prometheus has "target"s and it scrape metric data from the targets. Each FogFlow component ( Edge Node, Cloud Node) would be a target for prometheus. 
Every time an edge `worker` is registered to master via RabbitMQ protocol, the master should add the address of the edge node to prometheus config file. 

## How to do it
There are two ways to update `prometheus.yml` config file. 
- Using the config file as a shared volume between master and prometheus. This method might be challenging to scale out. 
- An HTTP call to prometheus (From whoever needs it, in this case: `master`). Prometheus does not support this by default and we rebuild the prometheus docker image for fogflow to support the feature of "remotely updating the configuration file".

## Configuration
- The `refresh_interval` in the `prometheus.yml` indicated the frequency of looking for new targets. The default value is 10 seconds. 
- k
## TODO
- The security concern: Some way of authenticating the caller should be provisioned. For now, the use case is at proto-type level, and community can decide on the method later on with more insight. 
