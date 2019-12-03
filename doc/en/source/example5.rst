*****************************************
Integrate an Actuator Device with Fogflow
*****************************************

The IoT devices can be of two types:

* the sensors which collect environmental data and push it to some application,
* the actuators which perform some action on the environment.

When data flow from a sensor device towards broker, it is called Northbound Flow, while in the other case when data flow from broker towards actuator devices, it is called Southbound Flow.

This tutorial will be focussed on Southbound Flow, i.e., how Fogflow will control the actuator devices to alter the environment.

To get a basic idea of how Southbound actually works in the context of FIWARE, refer `this`_ tutorial .

.. _`this`: https://fiware-tutorials.readthedocs.io/en/latest/iot-agent/index.html


Integration with NGSI supported Devices
-----------------------------------------------

The figure below shows the how southbound flow is accomplished in Fogflow.

.. figure:: figures/ngsi-device-integration.png
   :width: 100 %
   

To use an NGSI device, users can start this simple `Lamp`_ device, which prints the status of the lamp when an on/off command is received.

.. _`Lamp`: https://github.com/smartfog/fogflow/tree/master/application/device/lamp

After starting the lamp device, register the lamp device on Fogflow using the following curl request.

.. code-block:: console

    curl -iX POST \
    'http://<Thin_Broker_IP>:8070/NGSI9/registerContext' \
    -H 'Content-Type: application/json' \
    -H 'fiware-service: openiot' \
    -H 'fiware-servicepath: /' \
    -d '
    {
            "contextRegistrations": [
                {
                    "entities": [
                        {
                            "type": "Lamp",
                            "isPattern": "false",
                            "id": "Lamp.001"
                        }
                    ],
                    "attributes": [
                        {
                            "name": "on",
                            "type": "command"
                        }
                    ],
                    "providingApplication": "http://<Lamp_Host_IP>:8888"
                }
            ],
        "duration": "P1Y"
    }'

Below is the request to run an "on" command on the lamp (the NGSI device) to turn it on. Note that this request will be fired at thin broker. Thin borker will find the provider in the registrations and will send this command update to that provider, i.e. to the device.

.. code-block:: console

    curl -iX POST \
    'http://<Thin_Broker_IP>:8070/ngsi10/updateContext' \
    -H 'Content-Type: application/json' \
    -H 'fiware-service: openiot' \
    -H 'fiware-servicepath: /' \
    -H 'command: true' \
    -d '{	
        "contextElements": [
        {
            "entityId": {
            "id": "Lamp.001",
            "type": "Lamp",
            "isPattern": false
            },
            "attributes": [
                 {
                     "name": "on",
                     "type": "command",
                     "value": ""
                 }
             ]
        }
        ],
        "updateAction": "UPDATE"
    }'

On sending this command update, users can check the status of lamp device that was started in its logs. It will be "Lamp : on". Another supported command is "off" that the users can send to the device.


Integration with Non-NGSI supported Devices
-----------------------------------------------

FIWARE-provided IoT Agent will work as an intermediater between a Non-NGSI Device and Fogflow's thin broker in bidirectional manner. For devices based on a specific protocol, separate IoT Agent is there, for example, IoT Agent JSON for MQTT based devices, IoT Agent UL for Ultralight Devices, and so on. Southbound flow for Non-NGSI devices is shown in the figure below. It makes use of a device-protocol specific IoT Agent.

.. figure:: figures/non-ngsi-device-integration.png
   :width: 100 %
   
Using Ultralight devices
===============================================

Integration of an Ultralight actuator device with Fogflow is illustrated in the below example.

To work in Southbound using an Ultralight device, IoT Agent UL and Ultralight devices must be running. `Docker-Compose`_ file for this is given. The "tutorial" service in this file provides the device services. User need to edit this file based on their environment variables to get started.

.. _`Docker-Compose`: https://github.com/FIWARE/tutorials.IoT-Agent/blob/master/docker-compose.yml

The figure below shows the IoT Device monitor dashboard at http://tutorial_IP:3000/device/monitor

Please note that the "lamp001" is in "off" state. In this integration, we will light the lamp device using Fogflow.
    
.. figure:: figures/device-monitor-1.png
   :width: 100 %
   

**Registering a Device:** Device registeration is done at the IoT Agent to indicate what data the device will be providing. Following is the curl request for creating or registring a device on IoT Agent. Here, we are registering a lamp device with id "lamp001" that is supposed to be the context provider for entity "urn:ngsi-ld:Lamp:001". Corresponding to this, the IoT Agent will register the device in thin broker as well as create the entity for that device in thin broker itself.

.. code-block:: console

    curl -iX POST \
    'http://<IoT_Agent_IP>:4041/iot/devices' \
    -H 'Content-Type: application/json' \
    -H 'fiware-service: openiot' \
    -H 'fiware-servicepath: /' \
    -d '{
      "devices": [
        {
          "device_id": "lamp001",
          "entity_name": "urn:ngsi-ld:Lamp:001",
          "entity_type": "Lamp",
          "protocol": "Ultralight",
          "transport": "HTTP",
          "endpoint": "http://<Device_Host_IP>:3001/iot/lamp001",
          "commands": [
            {"name": "on","type": "command"},
            {"name": "off","type": "command"}
           ],
           "attributes": [
            {"object_id": "s", "name": "state", "type":"Text"},
            {"object_id": "l", "name": "luminosity", "type":"Integer"}
           ],
           "static_attributes": [
             {"name":"refStore", "type": "Relationship","value": "urn:ngsi-ld:Store:001"}
          ]
        }
      ]
    }'

**Sending command to device:** An external application or a Fog Function can control the actuator devices by sending commands like on/off, lock/unlock, open/close, or many others to the devices depending upon the type of device. The commands supported by a device will be known to Thin Broker through the device registration given above.

The below curl request sends an "on" command to the lamp001 device.

.. code-block:: console

    curl -iX POST \
    'http://<Thin_Broker_IP>:8070/ngsi10/updateContext' \
    -H 'Content-Type: application/json' \
    -H 'fiware-service: openiot' \
    -H 'fiware-servicepath: /' \
    -H 'command: true' \
    -d '{
        "contextElements": [
        {
            "entityId": {
            "id": "urn:ngsi-ld:Lamp:001",
            "type": "Lamp",
            "isPattern": false
            },
            "attributes": [
                 {
                     "name": "off",
                     "type": "command",
                     "value": ""
                 }
             ]
        }
        ],
        "updateAction": "UPDATE"
    }'
    
The above request shows Fogflow entity update, which is a bit different from the format suported by other brokers like FIWARE Orion. For that reason, below request is also supported in Fogflow.

.. code-block:: console

    curl -iX POST \
    'http://<Thin_Broker_IP>:8070/v1/updateContext' \
    -H 'Content-Type: application/json' \
    -H 'fiware-service: openiot' \
    -H 'fiware-servicepath: /' \
    -H 'command: true' \
    -d '{
        "contextElements": [
            {
                "type": "Lamp",
                "isPattern": "false",
                "id": "urn:ngsi-ld:Lamp:001",
                "attributes": [
                    {
                        "name": "on",
                        "type": "command",
                        "value": ""
                    }
                ]
            }
        ],
        "updateAction": "UPDATE"
    }'

Users can check the status of the Lamp again, it will in lit-up state as shown in the figure below.

.. figure:: figures/device-monitor-2.png
   :width: 100 %
   
