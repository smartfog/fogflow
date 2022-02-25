#!/bin/bash

#Check for command line argument
if (( $# == 1 )); then
    command="$1"
    
    echo "********************************************"
    echo "The external IP address is "$command
    echo "********************************************"
    
    sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" nginx.yaml
    sed -i "s/externalIPs:\ \[.*/externalIPs:\ \[${command}]/" rabbitmq.yaml
    sed -i "s/my_hostip\": \".*\"/my_hostip\": \"${command}\"/" configmap.yaml
fi

DIR="/tmp/fogflow/"
if [ ! -d "$DIR" ]; then
    # create a folder to save the FogFlow metadata, to be used by FogFlow designer
    echo "create a folder for the persistent storage at ${DIR}"
    mkdir $DIR
fi


kubectl create namespace fogflow


kubectl apply -f designer-pv.yaml

kubectl -n fogflow create -f designer-pvc.yaml




