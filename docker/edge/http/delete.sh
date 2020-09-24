#!/bin/bash

#Read Application ID and IDM Token from delete_config.js file
App_Id=`cat $(pwd)/delete_config.js | grep "APP_ID" | cut -f 2 -d "=" | sed 's/ //g' | tr -d '\n'`
App_token=`cat $(pwd)/delete_config.js | grep "IDM_TOKEN" | cut -f 2 -d "=" | sed 's/ //g' | tr -d '\n'`

#curl request to delete application from Keyrock
curl -X DELETE "http://localhost:3000/v1/applications/$App_Id" -H "X-Auth-token:$App_token"
