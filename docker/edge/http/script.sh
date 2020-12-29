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

#Read the Keyrock IP and Edge URL from oauth_config.js file

IP_ADDRESS=`cat $(pwd)/oauth_config.js | grep "IDM_IP" | cut -f 2 -d ":" | sed 's/"//g' | tr -d '\n'`

EDGE_IP_ADDRESS=`cat $(pwd)/oauth_config.js | grep "Edge_IP" | cut -f 2 -d ":" | sed 's/"//g' | tr -d '\n'`

#Generate API token, It will return X-Subject-Token that will further use for fetching details
#The curl command will retry thrice to reach to server, if packet drops in previous attempts
for connect in 1 2 3
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

#Displaying token
echo -----------------------------
echo IDM token is $token 
echo -----------------------------

#Fetch Application ID to authenticate the token shared by IOT devices.
for connect1 in 1 2 3
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

#Extracting ID from hidden file
ID=`cat .API_detail.js | tail -1| jq -r '.[] | .[] | .id' | tr -d '\n'`
if [ -z $ID ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Fetch secret of the Application to authenticate the token shared by IOT devices.
for connect2 in 1 2 3 
do 
  curl --include \
     --header "X-Auth-token: $token" \
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

#Displaying Applicatiion ID and Shared Secret 
echo -----------------------------
echo App ID and Passwords are $ID,$SECRET 
echo -----------------------------

#Configure PEP-PROXY at Edge
#Fetch PEP-Proxy ID
for connect3 in 1 2 3
do 
   curl --include \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/pep_proxies" > .PEP_Detail.js
   if [ $? -eq 0 ]; then
      break
   else
      echo reconnecting...
      continue
    fi
done

#Extracting PEP-PROXY ID from hidden file
Pep_ID=`cat .PEP_Detail.js | tail -1 | jq -r '.pep_proxy.id' | tr -d '\n'`
if [ -z $Pep_ID ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Fetch PEP-PROXY Password
for connect4 in 1 2 3 
do 
   curl --include \
     --request PATCH \
     --header "Content-Type: application/json" \
     --header "X-Auth-token: $token" \
  "http://$IP_ADDRESS:3000/v1/applications/$ID/pep_proxies" > .PEP_Details.txt 
  if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
  fi
done

#Extracting PEP-PROXY Password from hidden file
Pep_password=`cat .PEP_Details.txt | tail -1 | jq . | grep "new_password" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '/n'`
if [ -z $Pep_password ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Displaying PEP-PROXY ID and PEP-PROXY Password
echo -----------------------------
echo PEP PROXY ID is $Pep_ID and PEP PROXY PASSWORD is $Pep_password
echo --------------------------------

#setting up of pep-config.js
sed -i "s/PEP_PROXY_IDM_HOST.*/PEP_PROXY_IDM_HOST \|\| '${IP_ADDRESS}',/" pep-config.js
sed -i "s/PEP_PROXY_APP_HOST.*/PEP_PROXY_APP_HOST \|\| '${EDGE_IP_ADDRESS}',/" pep-config.js
sed -i "s/PEP_PROXY_APP_PORT.*/PEP_PROXY_APP_PORT \|\| '8060',/" pep-config.js
sed -i "s/PEP_PROXY_APP_ID.*/PEP_PROXY_APP_ID \|\| '${ID}',/" pep-config.js
sed -i "s/PEP_PROXY_USERNAME.*/PEP_PROXY_USERNAME \|\| '${Pep_ID}',/" pep-config.js 
sed -i "s/PEP_PASSWORD.*/PEP_PASSWORD \|\| '${Pep_password}',/" pep-config.js

#Generate X-Auth-token used to make request to keyrock
for connect5 in 1 2 3
do
   curl -X POST -H "Authorization: Basic $(echo -n $ID:$SECRET | base64 -w 0)"   -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password&username=admin@test.com&password=1234" "http://$IP_ADDRESS:3000/oauth2/token" > .access_token.txt
if [ $? -eq 0 ]; then
     break
  else
     echo reconnecting...
     continue
   fi
done

#Extracting X-Auth-token from hidden file
access=`cat .access_token.txt |jq . | grep "access_token" | cut -d ":" -f2 | sed -e 's/ //g' -e 's/"//g' -e 's/,//g' -e 's/\\n//g' | tr -d '\n'`
if [ -z $access ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Displaying Token 
echo -----------------------------
echo " access token is" $access 
echo -----------------------------

#To generate Authorization Token
AUTH=`echo -n $ID:$SECRET | base64 | tr -d "\t\r\n"` 
if [ -z $AUTH ]; then
   echo "Configure IDM .... Recheck and try again"
   exit 1
fi

#Displaying Authorization token used to register IOT Sensor with keyrock
echo ------------------------------------------
echo "Authorization : Basic "$AUTH 
echo ------------------------------------------

echo "end of automated script to configure PEP-Proxy at Edge Node"
rm -rf .access_token.txt .generate_token.txt .PEP_Details.txt .PEP_Detail.js .accessAPI.txt .API_detail.js 
