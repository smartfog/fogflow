#!/bin/bash

#Check for command line argument
if (( $# != 1 )); then
	echo "Illegal number of parameters"
	echo "usage: ./install.sh [externalIPs] e.g. :  ./install.sh 172.30.48.24"
	exit 1
fi

command="$1"

mkdir -p dgraph

echo "********************************************"
echo "The serving cloud IP address is "$command
echo "********************************************"

kubectl  create namespace fogflow

kubectl create -f serviceaccount.yaml

kubectl create -f configmap.yaml

#Configuring externalIPs in yaml files for deployment
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" discovery.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" broker.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" dgraph.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" rabbitmq.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" master.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" designer.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" nginx.yaml

#Deploying yaml files
kubectl create -f discovery.yaml
sleep 10s
kubectl create -f broker.yaml
kubectl create -f dgraph.yaml
kubectl create -f rabbitmq.yaml
sleep 25s
kubectl create -f master.yaml
sleep 10s
kubectl create -f worker.yaml
sleep 10s
kubectl create -f designer.yaml
sleep 10s
kubectl create -f nginx.yaml



