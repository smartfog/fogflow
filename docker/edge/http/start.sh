if [ $# -eq 0 ]; then
	htype='latest'
else
	htype='arm'
fi

docker run -d --name=edgebroker -v $(pwd)/config.json:/config.json -p 8060:8060  fogflow/broker:$htype
docker run -d --name=edgeworker -v $(pwd)/config.json:/config.json -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock fogflow/worker:$htype

# Edit <Edge_public_IP> in following command and uncomment it to run IoT Agent on Fogflow Edge Node.
# IoT Agent will use embedded mongodb i.e., mongodb will be running on localhost.
#docker run -d --name=iot-agent-json --env IOTA_CB_HOST=<Edge_public_IP> --env IOTA_CB_PORT=8070 --env IOTA_CB_NGSI_VERSION=v1 --env IOTA_MONGO_HOST=localhost --env IOTA_PROVIDER_URL=http://<Edge_public_IP>:4041 -p 4041:4041 -p 7896:7896 fogflow/iotajson-mongo:$htype
