curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
  -H 'Content-Type: application/json' \
  -H 'Destination: orion-broker' \
  -d '{
        "entities":[{"type":"AirQualityObserved","isPattern":true}], 
        "reference": "http://host.docker.internal:9090"
      }'

