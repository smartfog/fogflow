curl -iX POST \
  'http://10.146.0.2:8080/ngsi10/subscribeContext' \
  -H 'Content-Type: application/json'  \
  -H 'Destination: orion-broker'  \
  -d '
{
  "entities": [
    {
      "id": ".*",
      "type": "Result",
      "isPattern": true
    }
  ],
  "reference": "http://10.146.0.2:1026/v2"
}'
