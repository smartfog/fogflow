#!/bin/bash
#apt-get install jq -y

#Read the Keyrock IP, Application name, redirect_uri and application URL from user
echo Enter the Keyrock IDM IP:
read IP
IP_ADDRESS=`echo $IP | tr -d '\n'`
echo $IP | tr -d '\n'|od -c
echo IDM IP is $IP_ADDRESS
echo -----------------------------
echo Enter name of Application:
read name
appName=`echo $name | tr -d '\n'`
echo Name of application is $appName
echo -----------------------------
echo Enter redirect_uri/Callback_uri:
read callback
callbackUrl=`echo $callback | tr -d '\n'`
echo -----------------------------

echo Enter application url:
read Uri
appUri=`echo $Uri | tr -d '\n'`
echo -----------------------------

#Generate API token, It will return X-Subject-Token taht will further use for application register
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
  \"name\": \"admin@test.com\",
  \"password\": \"1234\"
}" \
"http://$IP_ADDRESS:3000/v1/auth/tokens" > generate_token.txt

token=`grep "X-Subject-Token" generate_token.txt | cut -f 2 -d ":" |  sed 's/\r$//g' | tr -d '\n'`
#token_id=`echo $token | sed 's/ //g'`
#echo $token | sed 's/ //g' > token.txt
echo IDM token is $token
echo -----------------------------
#echo "######## 1"

#Register application, it will return Application ID and Password. ID and Password will use to genenrate access token
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
#cat accessAPI.txt
#echo ""
# echo "######## 2"

ID=`cat accessAPI.txt | tail -1 | jq . | grep "id" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`
SECRET=`cat accessAPI.txt |  tail -1 |jq . | grep -m 1 "secret" | cut -f 2 -d ":" | sed -e 's/"//g' -e 's/,//g' -e 's/\r$//g' -e 's/ //g' | tr -d '\n'`

#echo $ID | od -c
#echo $SECRET | od -c

#Create access token with Resource Owner Password credentials, received from above step
echo App ID and Passwords are $ID,$SECRET
echo -----------------------------

#sh third.sh $ID $SECRET

curl -X POST -H "Authorization: Basic $(echo -n $ID:$SECRET | base64 -w 0)"   -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password&username=admin@test.com&password=1234" "http://$IP_ADDRESS:3000/oauth2/token" > access_token.txt
#cat access_token.txt 
access=`cat access_token.txt |jq . | grep "access_token" | cut -d ":" -f2 | sed -e 's/ //g' -e 's/"//g' -e 's/,//g' -e 's/\\n//g' | tr -d '\n'`
echo " access token is" $access

echo "end of script"



