#This script is used to launch the docker containers for edge components i.e. Fogflow Edge Broker and Fogflow Edge Worker
#The input to script decides which docker image is to be run in docker conatiners
#For e.g. 
# ./edge_start.sh : this command will run the docker image with "3.2.6" tag i.e. fogflow/broker:3.2.6
# ./edge_start.sh arm : this command will run the docker image with "arm" tag  i.e. fogflow/broker:arm

if [ $# -eq 0 ]; then
	htype='3.2.6'
else
	htype='arm'
fi

docker run -d --name=edgebroker -v $(pwd)/config.json:/config.json -p 8060:8060  fogflow/broker:$htype
docker run -d --name=edgeworker -v $(pwd)/config.json:/config.json -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock fogflow/worker:$htype

