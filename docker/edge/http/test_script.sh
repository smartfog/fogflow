#!/bin/bash
IP_ADDRESS=180.179.214.199
Edge_IP=180.179.214.211
Notifier=180.179.214.208

#To generate User Access token
for b in 1 2
do
   curl --include \
        --request POST \
        --header "Content-Type: application/json" \
        --data-binary "{
     \"name\": \"admin@test.com\",
     \"password\": \"1234\"
   }" \
   "http://$IP_ADDRESS:3000/v1/auth/tokens" > .generate_token.txt
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done

#Extracting token from hidden file
token=`grep "X-Subject-Token" .generate_token.txt | cut -f 2 -d ":" |  sed 's/\r$//g' | tr -d '\n'`
if [ -z $token ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#To fetch Application ID
for c in 1 2
do
  curl --include \
     --request GET \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications" > .API_detail.js
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
  fi
done

#Extracting Application ID from hidden file
ID=`cat .API_detail.js | tail -1| jq -r '.[] | .[] | .id' | tr -d '\n'`
if [ -z $ID ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Register a IoT Device for generating access token used in request made to edge
for a in 1 2
do
     curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/iot_agents" > .iot_device.js

    cat .iot_device.js

    id1=`cat .iot_device.js | tail -1 | jq -r '.iot_agent.id' | tr -d '\n'`
    pass=`cat .iot_device.js | tail -1 | jq -r '.iot_agent.password' | tr -d '\n'`
    echo ---------------------------
    echo "ID is $id1 and PASSWORD is $pass"
    echo ---------------------------
    
    curl -iX POST \
    "http://$IP_ADDRESS:3000/oauth2/token" \
    -H 'Accept: application/json' \
    -H 'Authorization: Basic NmUzOTZkZWYtM2ZhOS00ZmY5LTg0ZWItMjY2YzEzZTkzOTY0OjFlYmZiZmY4LWExNGUtNDFmOS1iYTMzLTI3MTRmNGIyNDkwNQ==' \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    --data "username=$id1&password=$pass&grant_type=password" > .access.js
 
    access_token=`cat .access.js | tail -1 | jq -r '.access_token' | tr -d '\n'`

    echo ------------------------------------
    echo "Access token is $access_token"
    echo -------------------------------------

    rm -rf .access.js
done

#To create and delete multiple sensors from IDM
for try in {1..100}
do 
     curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/iot_agents" > .iot.js

    cat .iot.js

    id=`cat .iot.js | tail -1 | jq -r '.iot_agent.id' | tr -d '\n'`
    echo ---------------------------
    echo "ID is $id"
    echo ---------------------------

    curl --include \
     --request DELETE \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/iot_agents/$id"

    rm -rf .iot.js
done

#Script to fire multiple registration requests
for connect in {1..100} 
do
	curl -iX  POST \
	"http://$Edge_IP:5556/NGSI9/registerContext" \
	-H 'Content-Type: application/json' \
	-H "X-Auth-token: $access_token" \
	-H 'fiware-service: openiot' \
	-H 'fiware-servicepath: /' \
	-d '
	{
        	"contextRegistrations": [
            	{	
                	"entities": [
                    	{		
                        	"type": "Lamp",
                        	"isPattern": "false",
                        	"id": "Lamp.00'$connect'"
                    	}
                	],
                	"attributes": [
                    	{
                        	"name": "on",
                        	"type": "command"
                    	},
                    	{
                        	"name": "off",
                        	"type": "command"
                    	}
                	],
                	"providingApplication": "http://'$Notifier':8888"
            	}
        	],
    	"duration": "P1Y"
	}'
done


#Making Update Request
for i in {1..100}
do 
	curl -iX POST \
	"http://$Edge_IP:5556/ngsi10/updateContext" \
	-H 'Content-Type: application/json' \
	-H "X-Auth-token: $access_token" \
	-H 'fiware-service: openiot' \
	-H 'fiware-servicepath: /' \
	-d '{
		"contextElements": [
		{
			"entityId": {
			"id": "Lamp.00'$i'",
			"type": "Lamp",
			"isPattern": false
			},
			"attributes": [
				 {
					 "name": "on",
					 "type": "command",
					 "value": ""
				 }
			 ]
		}
		],
		"updateAction": "UPDATE"
	}'
done
rm -rf .generate_token.txt .API_detail.js

