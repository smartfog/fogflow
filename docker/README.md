# FIWARE FogFlow for IoT Edge Computing


This folder includes all docker-compose files to deploy FogFlow, including the following two parts.  

## FogFlow Core

- deployment all FogFlow core components on a cloud node; the detailed steps are introduced [here](https://fogflow.readthedocs.io/en/latest/setup.html)

For a simple demonstration or test, this is already enough to try out the FogFlow system, because this docker compose file will also launch an edge node as part of the FogFlow Core. 

## FogFlow Edge Node

- deployment a FogFlow edge node; the detailed steps are introduced [here](https://fogflow.readthedocs.io/en/latest/edge.html)

## How to build each FogFlow component

the dockerfile files to build FogFlow components are located at the following folders

- discovery: /discovery/Dockerfile
	
- broker: /broker/Dockerfile
	
- master: /master/Dockerfile
	
- worker: /worker/Dockerfile
	
- designer: /designer/Dockerfile

