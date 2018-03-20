#start a container for test
#docker run -p 8001:8080 -t -i awsiot

#configurate
curl -X POST "http://127.0.0.1:8080/admin" -d @config.json --header "Content-Type:application/json" --header "Accept:application/json"

#issue a subscription to get the input data
curl -X POST "http://127.0.0.1:8070/ngsi10/subscribeContext" -d @subscription.json --header "Content-Type:application/json" --header "Accept:application/json"

