*****************************************
Integrate FogFlow with Scorpio Broker
*****************************************

NGSI-LD Adapter is built to enable FogFlow Ecosystem to provide Linked Data to the users. `Scorpio Broker`_ being the first reference implementation of NGSI-LD Specification, is being used here for receiving the Linked-Data from Fogflow.

.. _`Scorpio Broker`: https://scorpio.readthedocs.io/en/latest/

The figure below shows how NGSI-LD Adapter works in transforming the NGSIv1 data from Fogflow into NGSI-LD data to Scorpio Broker.

.. figure:: figures/ngsi-ld-adapter.png

1. User sends a subscription request to the adapter. 
2. The adapter then forwards this request to the Fogflow broker, to subscribe itself for the Context Data specified in its request.
3. Context data update is received at Fogflow broker.
4. Adapter receives notification from the Fogflow broker for the subscribed data.
5. Adapter converts the received data into NGSI-LD data format and forwards it to the Scorpio broker. 


Running NGSI-LD Adapter
===============================================

**Pre-Requisites:**

* Fogflow should be up and running with atleast one node.
* Scorpio broker should be up and running.

NGSI-LD Adapter can be run under Fogflow ecosystem using Fogflow Dashboard as given below. 

**Register an Operator:** Go to "Operator" in Operator Registry on Fogflow Dashboard. Register a new Operator with a Parameter Element as given below.
   
   Name: service_port ; Value: 8888
   
   (Is is assumed that the user has already gone through "REGISTER YOUR TASK OPERATORS" in `this`_ tutorial.)

.. _`this`: https://fogflow.readthedocs.io/en/latest/example1.html
   
**Register a Docker Image:** Go to "DockerImage" in Operator Registry and register an image fogflow/ngsildadapter:latest. Associate it with the above operator by choosing the operator from DropDown. Users can also build their image for NGSI-LD-Adapter by editing and running `build`_ file.

.. _`build`: https://github.com/smartfog/fogflow/blob/document-update/application/operator/NGSI-LD-Adapter/build

**Register a Fog Function** as shown in the figure below. In "SelectedType", provide the Entity Type (say "LD") of the Context Data that will be used to trigger this Fog Function. Choose the operator registered in Step#1 as the operator in Fog Function.

.. figure:: figures/fogfunction_ngsi-ld-adapter.png


**Trigger the Fog Function** by sending an update request to Fogflow Broker with the Entity Type as "LD" (or whatever is specified in Step#3 as the SelectedType). It should include fogflowIP and ngbIP in the attributes along with location metadata. Example request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<Fogflow-Broker-IP>:8070/ngsi10/updateContext' \
      -H 'Content-Type: application/json' \
      -d '
      {
        "contextElements": [
        {
            "entityId": {
            "id": "LD001",
            "type": "LD",
            "isPattern": false
            },
            "attributes": [
                 {
                     "name": "fogflowIP",
                     "type": "string",
                     "value": "<IP>"
                 },
                 {
                     "name": "ngbIP",
                     "type": "string",
                     "value": "<IP>"
                 }
             ],
             "domainMetadata": [
                 {
                     "name": "location",
                     "type": "point",
                     "value": {
                                  "latitude": 52,
                                  "longitude": 67
                     }
                 }
             ]
        }
        ],
        "updateAction": "UPDATE"
       }'


NGSI-LD-Adapter task will be created and it will be listening on port 8888. Users can list it in the tasks running on either the cloud node or the edge node, whichever is nearest to the location provided in the metadata of the above request. 


Using NGSI-LD Adapter
===============================================

To use the NGSI-LD-Adapter for context data transformation, follow the below steps.


**Send subscription request** to LD-Adapter, it will forward the same request to Fogflow Broker. This is because the access to Fogflow broker will not be available directly to the user. Examle Subscription request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<LD-Adapter-Host-IP>:8888/subscribeContext' \
      -H 'Content-Type: application/json' \
      -d '
    {
      "entities": [
        {
          "id": "Temperature.*",
          "type": "Temperature",
          "isPattern": true
        }
      ],
      "attributes": [
        "temp"
      ],
      "restriction": {
        "scopes": [
          {
            "scopeType": "circle",
            "scopeValue": {
              "centerLatitude": 49.406393,
              "centerLongitude": 8.684208,
              "radius": 2000
            }
          }
        ]
      },
      "reference": "http://<LD-Adapter-Host-IP>:8888"
    }'


**Send update request** to Fogflow Broker with an entity of type and attributes defined in the above subscription. An example request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<Fogflow-Broker-IP>:8070/ngsi10/updateContext' \
      -H 'Content-Type: application/json' \
      -d '
      {
        "contextElements": [
          {
            "entityId": {
              "id": "Temperature001",
              "type": "Temperature",
              "isPattern": false
            },
            "attributes": [
              {
                "name": "temp",
                "type": "float",
                "value": 34
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
              }
             ]
          }
        ],
        "updateAction": "UPDATE"
      }'


Check if the entity in NGSI-LD format has been updated on Scorpio Broker by visiting URL:  http://<Scorpio-Broker-IP:Port>/ngsi-ld/v1/entities?type=http://example.org/Temperature

Following code block shows the trasformed context data.

.. code-block:: console

    {"@context": ["https://schema.lab.fiware.org/ld/context", "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
    {"Temperature": "http://example.org/Temperature", "temp": "http://example.org/temp"}], "type": "Temperature", 
    "id": "urn:ngsi-ld:Temperature001", "temp": {"type": "Property", "value": 34}, "location": {"type": "GeoProperty", 
    "value": "{\"type\": \"point\", \"coordinates\": [49.406393, 8.684208]}"}}
