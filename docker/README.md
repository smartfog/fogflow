# FIWARE FogFlow for IoT Edge Computing


This folder includes all docker-compose files to deploy FogFlow, including the following two parts.  

## How to deploy the FogFlow system with the provided docker images

### FogFlow Core

- deployment all FogFlow core components on a cloud node; the detailed steps are introduced [here](https://fogflow.readthedocs.io/en/latest/setup.html)

For a simple demonstration or test, this is already enough to try out the FogFlow system, because this docker compose file will also launch an edge node as part of the FogFlow Core. 

### FogFlow Edge Node

- deployment a FogFlow edge node; the detailed steps are introduced [here](https://fogflow.readthedocs.io/en/latest/edge.html)

## How to build each FogFlow component

the dockerfile files to build FogFlow components are located at the following folders

- discovery: /discovery/Dockerfile, please check the [readme](https://github.com/smartfog/fogflow/tree/master/discovery) to see the detail instruction
	
- broker: /broker/Dockerfile, lease check the [readme](https://github.com/smartfog/fogflow/tree/master/broker) to see the detail instruction
	
- master: /master/Dockerfile, lease check the [readme](https://github.com/smartfog/fogflow/tree/master/master) to see the detail instruction
	
- worker: /worker/Dockerfile, lease check the [readme](https://github.com/smartfog/fogflow/tree/master/worker) to see the detail instruction
	
- designer: /designer/Dockerfile, lease check the [readme](https://github.com/smartfog/fogflow/tree/master/designer) to see the detail instruction

A bash script is provided to build the images of all FogFlow components. 

```console
./build.sh

Once you log in to your own docker hub account, you can publish all generated docker images to your own docker registry. 

```console
./publish.sh


