Pre-requisite: GPA should be running as a Fog Function. Refer [README for GPA](https://github.com/smartfog/fogflow/blob/master/application/operator/GeneralPurposeAdapter/README.md).

Follow these steps to run IoT Agent JSON as a Fog Function:

1. Register a new operator with image "fogflow/iotagent" in Operator Registry in FogFlow.

2. Register a FogFunction and set an Entity Type (say IOTA) as SelectedType. Choose the operator registered in Step#1 as Fog Function Operator.

3. Trigger the Fog Function by sending an update request to Fogflow Broker with an Entity Type (say IOTA). It should include GPA Broker IP and Port. Example request is given below:

```bash
curl -iX POST \
  'http://<Fogflow_Broker_IP>:<Fogflow_Broker_Port>/ngsi10/updateContext' \
  -H 'Content-Type: application/json' \
  -d '
{
    "contextElements": [
    {
        "entityId": {
        "id": "IOTA001",
        "type": "IOTA",
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

IOTA task will be created soon on the nearest edge and it will be listening on port 4041 for registrations and port 7896 for data. 
Both IOTA and GPA can be verified by pushing data to Fogflow through IoT Agent.

Steps to push data to IoT Agent are given below with example requests:

- **Register a service**:

```bash
curl -iX POST \
  'http://<IoT_Agent_IP>:4041/iot/services' \
  -H 'Content-Type: application/json' \
  -H 'fiware-service: iota' \
  -H 'fiware-servicepath: /' \
  -d '{
"services": [
   {
     "apikey":      "FFNN1111",
     "entity_type": "Thing",
     "resource":    "/iot/json"
   }
]
}'
```

- **Register a device within service**:

```bash
curl -X POST \
  http://<IoT_Agent_IP>:4041/iot/devices \
  -H 'content-type: application/json' \
  -H 'fiware-service: iota' \
  -H 'fiware-servicepath: /' \
  -d '{
        "devices": [{
                "device_id": "Device1111",
                "entity_name": "Thing1111",
                "entity_type": "Thing",
                "attributes": [{
                        "object_id":"locationName",
                        "name": "locationName",
                        "type": "string"
                },{
                        "object_id": "locationId",
                        "name": "locationId",
                        "type": "string"
                },{
                        "object_id": "Temperature",
                        "name": "Temperature",
                        "type": "int"
                }
                ]}]
}'
```

- **Push data using the registered device**:

```bash
curl -X POST \
  'http://<IoT_Agent_IP>:7896/iot/json?i=Device1111&k=FFNN1111' \
  -H 'content-type: application/json' \
  -H 'fiware-service: iota' \
  -H 'fiware-servicepath: /' \
  -d '{ 
    "locationName":"Delhi",
    "locationId":"0011",
    "Temperature":45
}'
```

Verify the data on Fogflow Broker by visiting http://<Fogflow_Broker_IP>:<Fogflow_Broker_Port>/ngsi10/entity
