#!/bin/sh

echo "Configuring Transformer Function..."

echo $1
echo $2

sed -i  "s/fogflow_subscription_endpoint=.*/fogflow_subscription_endpoint=$1:8080/" ./config/config.ini
sed -i  "s/ngsi-ld-broker=.*/ngsi-ld-broker=$2:9090/" ./config/config.ini



python main.py
