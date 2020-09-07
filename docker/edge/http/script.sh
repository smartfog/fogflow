#!/bin/bash

#Check if Curl command is present or not, If not present then install
if ! command -v curl &> /dev/null
then
    echo "curl could not be found"
    sudo apt-get install curl -y
    exit
fi

#Check if jq command is present or not, If not present then install
if ! command -v curl &> /dev/null
then
    echo "jq could not be found"
    sudo apt-get install jq -y
    exit
fi

#Read the Keyrock IP, Application name, redirect_uri and application URL from oauth_config.js file
IP_ADDRESS=`cat $(pwd)/oauth_config.js | grep "IDM_IP" | cut -f 2 -d ":" | sed 's/"//g' | tr -d '\n'`

appName=`cat $(pwd)/oauth_config.js | grep "Application_Name" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`


callbackUrl=`cat $(pwd)/oauth_config.js | grep "Redirect_uri" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`


appUri=`cat $(pwd)/oauth_config.js | grep "Url" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`



#Generate API token, It will return X-Subject-Token taht will further use for application register
curl --connect-timeout 30 --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
  \"name\": \"admin@test.com\",
  \"password\": \"1234\"
}" \
"http://$IP_ADDRESS:3000/v1/auth/tokens" > generate_token.txt

token=`grep "X-Subject-Token" generate_token.txt | cut -f 2 -d ":" |  sed 's/\r$//g' | tr -d '\n'`
echo IDM token is $token
echo -----------------------------

#Register application, it will return Application ID and Password. ID and Password will use to genenrate access token
curl --connect-timeout 30 -sb --include --request POST  -H "Content-Type: application/json"  -H "X-Auth-token: $token" --data-binary "{
  \"application\": {
    \"name\": \"$appName\",
    \"description\": \"description\",
    \"redirect_uri\": \"$callbackUrl\",
    \"url\": \"$appUri\",
    \"grant_type\": [
      \"authorization_code\",
      \"implicit\",
      \"password\"
    ],
    \"token_types\": [
        \"jwt\",
        \"permanent\"
    ]
  }
}" \
"http://$IP_ADDRESS:3000/v1/applications" > accessAPI.txt

ID=`cat accessAPI.txt | tail -1 | jq . | grep "id" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`
SECRET=`cat accessAPI.txt |  tail -1 |jq . | grep -m 1 "secret" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`


#Create access token with Resource Owner Password credentials, received from above step
echo App ID and Passwords are $ID,$SECRET
echo -----------------------------


curl --connect-timeout 30 -X POST -H "Authorization: Basic $(echo -n $ID:$SECRET | base64 -w 0)"   -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password&username=admin@test.com&password=1234" "http://$IP_ADDRESS:3000/oauth2/token" > access_token.txt
#cat access_token.txt 
access=`cat access_token.txt |jq . | grep "access_token" | cut -d ":" -f2 | sed -e 's/ //g' -e 's/"//g' -e 's/,//g' -e 's/\\n//g' | tr -d '\n'`
echo " access token is" $access

echo "end of script"
rm -rf accessAPI.txt access_token.txt generate_token.txt


