#start a container for test
#docker run -p 8009:8080 -t -i facefinder /bin/bash

#configurate
curl -X POST "http://192.168.1.80:8009/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://192.168.1.80:8091/ngsi10/subscribeContext" -d @subscriptionCamera.json --header "Content-Type:application/json" --header "Accept:application/json"
curl -X POST "http://192.168.1.80:8091/ngsi10/subscribeContext" -d @subscriptionChildLost.json --header "Content-Type:application/json" --header "Accept:application/json"

