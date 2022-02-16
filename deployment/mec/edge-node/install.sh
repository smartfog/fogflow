#!/bin/bash

#Check for command line argument
if (( $# != 1 )); then
        echo "Illegal number of parameters"
        echo "usage: ./install.sh [externalIPs] e.g. :  ./install.sh 172.30.48.46"
        exit 1
fi

command="$1"

echo "************************************************"
echo "The serving IP address of edge is "$command
echo "************************************************"

#creating deployments 
microk8s.kubectl  create namespace fogflow

microk8s.kubectl create serviceaccount edge -n fogflow

microk8s.kubectl create -f edge-serviceaccount.yaml

microk8s.kubectl create -f edge-configmap.yaml

#configuring external IP for edge broker
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" edge-broker.yaml

microk8s.kubectl create -f edge-broker.yaml
microk8s.kubectl create -f edge-worker.yaml

