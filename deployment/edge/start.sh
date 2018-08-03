if [ $# -eq 0 ]; then
	htype='latest'
else
	htype='arm'
fi

docker run -d --name=edgebroker -v $(pwd)/config.json:/config.json -p 80:80  fogflow/broker:$htype
docker run -d --name=edgeworker -v $(pwd)/config.json:/config.json -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock fogflow/worker:$htype

