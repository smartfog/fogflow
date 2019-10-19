#start a container for test
docker run --name legotpu -d --privileged -p 8008:8008 -v /dev/bus/usb:/dev/bus/usb --env-file env.list fogflow/tpu

sleep 5

#issue a subscription to get the input data
curl -X POST "http://192.168.1.100/ngsi10/subscribeContext" -d @subscription.json --header "Content-Type:application/json" --header "Accept:application/json"

