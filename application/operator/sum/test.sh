#start a container for test
docker run -p 8001:8080 -t -i sum

#configurate
curl -X POST "http://192.168.1.80:8001/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://192.168.1.80:8091/ngsi10/subscribeContext" -d @subscription.json --header "Content-Type:application/json" --header "Accept:application/json"

