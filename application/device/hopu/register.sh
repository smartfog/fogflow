curl -iX POST \
  'http://127.0.0.1:8080/device' \
  -H 'Content-Type: application/json' \
  -d '
    {
        "id": "urn:Device.12345",
        "type": "HOPU",
        "attributes": {
            "protocol": {
                "type": "string",
                "value": "MQTT"
                },
            "mqttbroker":  {
                "type": "string",
                "value": "mqtt://mqtt.cdtidev.nec-ccoc.com:1883"
                },
            "topic": {
                "type": "string",
                "value": "/api/12345/attrs"
                },
            "mappings": {
                "type": "object",
                "value": {
                    "temp8": {
                        "name": "temperature",
                        "type": "float",
                        "entity_type": "AirQualityObserved"
                    },
                    "hum8": {
                        "name": "humidity",
                        "type": "float",
                        "entity_type": "AirQualityObserved"
                    }
                }
            }          
        },
        "metadata": {
            "location": {
                    "type": "point",
                    "value": {
                        "latitude": 49.406393,
                        "longitude": 8.684208
                    }
            }
        }
    }'
