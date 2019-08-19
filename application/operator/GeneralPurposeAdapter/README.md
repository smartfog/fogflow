Follow these steps to run General Purpose Adapter as a FogFunction:

1. Register a new operator with image "fogflow/generalpurposeadapter" in Operator Registry in FogFlow.

2. Register a FogFunction and set an Entity Type (say GPA) as SelectedType. Choose the operator registered in Step#1 as Fog Function Operator.

3. Trigger the Fog Function by sending an update request to Fogflow Broker with an Entity Type (say GPA). It should include Fogflow Broker IP and Port. Example request is given below:

```bash
curl -iX POST \
  'http://<Fogflow_Broker_IP>:8080/ngsi10/updateContext' \
  -H 'Content-Type: application/json' \
  -d '
{
    "contextElements": [
    {
        "entityId": {
        "id": "GPA001",
        "type": "GPA",
        "isPattern": false
        },
        "attributes": [
             {
                 "name": "brokerIP",
                 "type": "string",
                 "contextValue": "<Broker_IP>"
             },
             {
                 "name": "brokerPort",
                 "type": "string",
                 "contextValue": "<Broker_Port>"
             }
         ],
         "domainMetadata": [
             {
                 "name": "location",
                 "type": "point",
                 "value": {
                              "latitude": 37,
                              "longitude": 138
                 }
             }
         ]
    }
    ],
    "updateAction": "UPDATE"
}'
```

GPA task will be created soon on the nearest edge and it will be listening on port 1026. GPA can be verified after running IoT Agent as a Fog Function such that GPA serves as Context Broker to IoT Agent. This can be done by providing the GPA IP and Port in IoT Agent configuration. Refer the IoT Agent Operator README for more detail.
