#start a container for test
docker run -p 8009:8080 -t -i local/facefinder #/bin/bash

#configurate
curl -X POST "http://192.168.126.128:8009/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://192.168.126.128:8091/ngsi10/subscribeContext" -d @subscriptionCamera.json --header "Content-Type:application/json" --header "Accept:application/json"
curl -X POST "http://192.168.126.128:8091/ngsi10/subscribeContext" -d @subscriptionChildLost.json --header "Content-Type:application/json" --header "Accept:application/json"

