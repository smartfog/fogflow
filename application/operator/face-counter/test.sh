#start a container for test
#docker run -p 8000:8080 -t -i facecounter /bin/bash

#configurate
curl -X POST "http://192.168.1.102:8133/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://192.168.1.102:8071/ngsi10/subscribeContext" -d @subscription.json --header "Content-Type:application/json" --header "Accept:application/json"

