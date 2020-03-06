*****************************************
Integrate FogFlow with WireCloud
*****************************************

`WireCloud`_ builds on cutting-edge end-user development, RIA and semantic technologies to offer a next-generation end-user centred web application mashup platform aimed at leveraging the long tail of the Internet of Services.

.. _`WireCloud`: https://wirecloud.readthedocs.io/en/stable/

**Send update request**  to Fogflow Broker with an entity of type and attributes defined in the above subscription.
An example request is given below:

.. code-block:: console

        curl -iX POST \
        'http://<Fogflow broker>:8070/ngsi10/updateContext' \
         -H 'Content-Type: application/json' \
        -d '
      {
        "contextElements": [{
                "entityId": {
                        "id": "Room4",
                        "type": "Room",
                        "isPattern": false
                },
                "attributes": [{
                        "name": "temperature",
                        "type": "Integer",
                        "value": 155
                }],
                "domainMetadata": [{
                        "name": "location",
                        "type": "point",
                        "value": {
                                "latitude": 49.406393,
                                "longitude": 8.684208
                        }
                }]
        }],
        "updateAction": "UPDATE"
     }'




