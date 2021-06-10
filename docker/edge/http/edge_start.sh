if [ $# -eq 0 ]; then
	htype='3.0'
else
	htype='arm'
fi

docker run -d --name=edgebroker -v $(pwd)/config.json:/config.json -p 8060:8060  fogflow/broker:$htype
docker run -d --name=edgeworker -v $(pwd)/config.json:/config.json -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock fogflow/worker:$htype

