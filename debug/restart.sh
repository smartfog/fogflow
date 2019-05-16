#!/bin/bash


#kill fogflow
docker kill $(docker ps -q)

#login to docker hub
docker login

#rebuild  the master
echo "rebuiling master..."
cd ../master
docker login
./build


#rebuild the worker
echo "rebuilding worker"
cd ../worker
docker login
./build

cd ../debug
docker-compose up -d 

