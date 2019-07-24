#!/bin/sh
echo "Creating config.js for IoTAgent....."

sed -i  "s/        port: '27017'/        port: '$2'/" ./config.js

sed -i  "s/        port: '1026'/        port: '$4'/" ./config.js

sed -i '122s/localhost/'$3'/' ./config.js

sed -i '205s/localhost/'$1'/' ./config.js

echo "Starting IoTAgent....."

./docker/entrypoint.sh ./config.js
