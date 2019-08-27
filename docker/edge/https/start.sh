if [ $# -eq 0 ]; then
	htype='latest'
else
	htype='arm'
fi

#docker run -d --name=cadvisor -v /:/rootfs:ro -v /var/run:/var/run:rw -v /sys:/sys:ro -v /var/lib/docker/:/var/lib/docker:ro  -p 9092:8080  google/cadvisor 
docker run -d --name=edgebroker -v $(pwd)/config.json:/config.json -p 8082:8080  fogflow/broker:$htype
docker run -d --name=edgeworker -v $(pwd)/config.json:/config.json -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock fogflow/worker:$htype

