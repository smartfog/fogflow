if [ $# -eq 0 ]; then
	htype='latest'
else
	htype='arm'
fi

docker run -d --name=broker -v $(pwd)/config.json:/config.json -p 8080  fogflow/broker:$htype
docker run -d --name=worker -v $(pwd)/config.json:/config.json fogflow/worker:$htype

