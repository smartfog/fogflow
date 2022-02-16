#!/bin/bash

#Check for command line argument
if (( $# != 1 )); then
	echo "Illegal number of parameters"
	echo "usage: ./configurate.sh [externalIPs] e.g. :  ./configurate.sh 172.30.48.24"
	exit 1
fi

command="$1"

echo "********************************************"
echo "The external IP address is "$command
echo "********************************************"

sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" nginx.yaml
sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" rabbitmq.yaml

kubectl  create namespace fogflow

kubectl create -f serviceaccount.yaml


#Configuring externalIPs in yaml files for deployment
#sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" discovery.yaml
#sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" broker.yaml
#sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" rabbitmq.yaml
#sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" master.yaml
#sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" designer.yaml


##Deploying yaml files
#kubectl create -f rabbitmq.yaml
#kubectl create -f nginx.yaml

#kubectl create -f discovery.yaml
#kubectl create -f broker.yaml
#kubectl create -f master.yaml
#kubectl create -f worker.yaml
#kubectl create -f designer.yaml



