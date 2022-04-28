curl -iX POST \
  'http://localhost:8070/ngsi10/updateContext' \
  -H 'Content-Type: application/json' \
  -d '
    {
        "contextElements": [
            {
                "entityId": {
                    "id": "Device.12345",
                    "type": "HOPUDevice",
                    "isPattern": false
                },
                "attributes": [
                {
                  "name": "mqttbroker",
                  "type": "string",
                  "value": "mqtt://mqtt.cdtidev.nec-ccoc.com:1883"
                },
                {
                  "name": "topic",
                  "type": "string",
                  "value": "/api/12345/attrs"
                },                
                {
                  "name": "mappings",
                  "type": "object",
                  "value": {
                      "temp8": {
                            "name": "temperature",
                            "type": "Number",
                            "entity_type": "AirQualityObserved"
                           },
                       "hum8": {
                            "name": "humidity",
                            "type": "Number",
                            "entity_type": "AirQualityObserved"
                        }
                    }
                }
                ],
                "domainMetadata": [
                {
                    "name": "location",
                    "type": "point",
                    "value": {
                        "latitude": 49.406393,
                        "longitude": 8.684208
                    }
                },{
                    "name": "city",
                    "type": "string",
                    "value": "Heidelberg"
                }
                ]
            }
        ],
        "updateAction": "UPDATE"
    }'
