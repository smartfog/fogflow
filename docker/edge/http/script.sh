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
#Read the Keyrock IP, Application name, redirect_uri and application URL from oauth_config.js file
IP_ADDRESS=`cat $(pwd)/oauth_config.js | grep "IDM_IP" | cut -f 2 -d ":" | sed 's/"//g' | tr -d '\n'`

appName=`cat $(pwd)/oauth_config.js | grep "Application_Name" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`

callbackUrl=`cat $(pwd)/oauth_config.js | grep "Redirect_uri" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`

appUri=`cat $(pwd)/oauth_config.js | grep "Url" | cut -d ":" -f 2- | sed 's/"//g' | tr -d '\n'`

#Generate API token, It will return X-Subject-Token taht will further use for application register
for connect in 1 2 3
do
   curl --include \
        --request POST \
        --header "Content-Type: application/json" \
        --data-binary "{
     \"name\": \"admin@test.com\",
     \"password\": \"1234\"
   }" \
   "http://$IP_ADDRESS:3000/v1/auth/tokens" > generate_token.txt
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done

token=`grep "X-Subject-Token" generate_token.txt | cut -f 2 -d ":" |  sed 's/\r$//g' | tr -d '\n'`
echo -----------------------------
echo IDM token is $token 
echo -----------------------------

#Register application, it will return Application ID and Password. ID and Password will use to genenrate access token
for connect1 in 1 2 3
do
  curl -sb --include --request POST  -H "Content-Type: application/json"  -H "X-Auth-token: $token" --data-binary "{
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
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done

ID=`cat accessAPI.txt | tail -1 | jq . | grep "id" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`
SECRET=`cat accessAPI.txt |  tail -1 |jq . | grep -m 1 "secret" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`


#Create access token with Resource Owner Password credentials, received from above step
echo -----------------------------
echo App ID and Passwords are $ID,$SECRET 
echo -----------------------------

for connect2 in 1 2 3
do
   curl -X POST -H "Authorization: Basic $(echo -n $ID:$SECRET | base64 -w 0)"   -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password&username=admin@test.com&password=1234" "http://$IP_ADDRESS:3000/oauth2/token" > access_token.txt
if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done
#cat access_token.txt 
access=`cat access_token.txt |jq . | grep "access_token" | cut -d ":" -f2 | sed -e 's/ //g' -e 's/"//g' -e 's/,//g' -e 's/\\n//g' | tr -d '\n'`
echo -----------------------------
echo " access token is" $access 
echo -----------------------------

echo "end of script"
rm -rf accessAPI.txt access_token.txt generate_token.txt


