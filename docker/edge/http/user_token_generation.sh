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

#Script to generate user token.
#user credentials are fetched while executing the script as command line argument

for connect in 1 2 3 
do 
   curl --include \
        --request POST \
        --header "Content-Type: application/json" \
        --data-binary "{
     \"name\": \"$1\",
     \"password\": \"$2\"
   }" \
   "http://$IP_ADDRESS:3000/v1/auth/tokens" > .generate_user_token.txt
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done

#Extracting token from hidden file
token=`grep "X-Subject-Token" .generate_user_token.txt | cut -f 2 -d ":" |  sed 's/\r$//g' | tr -d '\n'`
if [ -z $token ]; then
   echo "Recheck and try again"
   exit 1
fi

#Displaying token
echo -----------------------------
echo IDM token is $token
echo -----------------------------

rm -rf .generate_user_token.txt



