#!/bin/bash

#Check if Curl command is present or not, If not present then install
if command -v curl >/dev/null
then
   continue
else
    echo "curl could not be found"
    apt-get install curl -y &> /dev/null;
fi

#Check if jq command is present or not, If not present then install
if command -v jq >/dev/null
then
   continue
else
    echo "jq could not be found"
    apt-get install jq -y &> /dev/null
fi

#Read the Keyrock IP from oauth_config.js file

IP_ADDRESS=`cat $(pwd)/oauth_config.js | grep "IDM_IP" | cut -f 2 -d ":" | sed 's/"//g' | tr -d '\n'`

#To fetch Application ID
for connect1 in 1 2
do
  curl --include \
     --request GET \
     --header "X-Auth-token: $1" \
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

#Fetch secret of the Application to authenticate the token shared by IOT devices.
for connect2 in 1 2 3
do
  curl --include \
     --header "X-Auth-token: $1" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID" > .accessAPI.txt
  if [ $? -eq 0 ]; then
      break
  else
      echo reconnecting...
      continue
  fi
done

#Extracting SECRET from hidden file
SECRET=`cat .accessAPI.txt |  tail -1 |jq . | grep -m 1 "secret" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`
if [ -z $SECRET ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Register a IoT Device for generating access token used in request made to edge
for a in 1 2
do
     curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --header "X-Auth-token: $1" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/iot_agents" > .iot_device.js

    id1=`cat .iot_device.js | tail -1 | jq -r '.iot_agent.id' | tr -d '\n'`
    pass=`cat .iot_device.js | tail -1 | jq -r '.iot_agent.password' | tr -d '\n'`
    echo ---------------------------
    echo "ID is $id1 and PASSWORD is $pass"
    echo ---------------------------

    curl -iX POST \
    "http://$IP_ADDRESS:3000/oauth2/token" \
    -H 'Accept: application/json' \
    -H "Authorization: Basic $(echo -n $ID:$SECRET | base64 -w 0)" \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    --data "username=$id1&password=$pass&grant_type=password" > .access.js
     
    access_token=`cat .access.js | tail -1 | jq -r '.access_token' | tr -d '\n'`

    echo ------------------------------------
    echo "Device access token is $access_token"
    echo -------------------------------------
    
    if [ $? -eq 0 ]; then
         break
     else
         echo reconnecting...
         continue
     fi
done

rm -rf .access.js .iot_device.js .accessAPI.txt .API_detail.js
