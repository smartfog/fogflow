#start a container for test
docker run -d --privileged -p 8008:8008 -v /dev/bus/usb:/dev/bus/usb  --name lego_tpu fogflow/tpu

#configurate
curl -X POST "http://192.168.1.100:8008/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://192.168.1.100/ngsi10/subscribeContext" -d @subscription.json --header "Content-Type:application/json" --header "Accept:application/json"

