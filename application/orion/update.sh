curl -iX POST \
 'http://10.146.0.2:8080/ngsi10/updateContext' \
 -H 'Content-Type: application/json' \
 -d '
{
   "contextElements": [
       {
           "entityId": {
               "id": "Device.Temperature.100",
               "type": "Temperature",
               "isPattern": false
               },
           "attributes": [
                   {
                   "name": "temperature",
                   "type": "float",
                   "contextValue": 73
                   },
                   {
                   "name": "pressure",
                   "type": "float",
                   "contextValue": 44
                   }
               ],
           "domainMetadata": [
                   {
                   "name": "location",
                   "type": "point",
                   "value": {
                   "latitude": -33.1,
                   "longitude": -1.1
                   }}
               ]
       }
   ],
   "updateAction": "UPDATE"
}'

