#!/bin/sh
echo "Configuring IoTAgent....."

sed -i '134s/localhost/'$1'/' ./config.js

sed -i  "s/        port: '1026'/        port: '$2'/" ./config.js

echo "Starting MongoDB....."
cd /usr/local/bin/
./docker-entrypoint.sh mongod &


echo "Starting IoTAgent....."
cd /opt/iotajson
./docker/entrypoint.sh config.js
