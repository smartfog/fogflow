#!/bin/sh

echo "Configuring Transformer Function..."

sed -i  "s/fogflow_subscription_endpoint=.*/fogflow_subscription_endpoint=$1:8070/" ./config/config.ini
sed -i  "s/ngsi-ld-broker=.*/ngsi-ld-broker=$2:9090/" ./config/config.ini

python main.py
