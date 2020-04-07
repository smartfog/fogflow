*****************************************
Integrate an actuator device with Fogflow
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
                        },
                        {
                            "name": "off",
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
Users can have their own customized devices that send the command updates in northbound direction also.

Integration with Non-NGSI supported Devices
-----------------------------------------------

FIWARE-provided IoT Agent will work as an intermediater between a Non-NGSI Device and Fogflow's thin broker in bidirectional manner. For devices based on a specific protocol, separate IoT Agent is there, for example, IoT Agent JSON for MQTT based devices, IoT Agent UL for Ultralight Devices, and so on. Southbound flow for Non-NGSI devices is shown in the figure below. It makes use of a device-protocol specific IoT Agent.

.. figure:: figures/non-ngsi-device-integration.png

   
Using Ultralight devices
===============================================

Integration of an Ultralight actuator device with Fogflow is illustrated in the below example.

To work in Southbound using an Ultralight device, IoT Agent UL and Ultralight devices must be running. Docker-Compose file for this is given `here`_. The "tutorial" service in this file provides the device services. Users need to edit this file based on their environment variables to get started.

.. _`here`: https://github.com/FIWARE/tutorials.IoT-Agent/blob/master/docker-compose.yml

The figure below shows the IoT Device monitor dashboard at http://tutorial_IP:3000/device/monitor

Please note that the "lamp001" is in "off" state. In this integration, we will light the lamp device using Fogflow.
    
.. figure:: figures/device-monitor-1.png
   

**Registering a Device:** Device registeration is done at the IoT Agent to indicate what data the device will be providing. Following is the curl request for creating or registring a device on IoT Agent. Here, a lamp device is registered  with id "lamp001" that is supposed to be the context provider for entity "urn:ngsi-ld:Lamp:001". Corresponding to this, the IoT Agent will register the device in thin broker as well as create the entity for that device in thin broker itself.

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
                     "name": "on",
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

Users can check the status of the Lamp again, it will be in lit-up state as shown in the figure below.

.. figure:: figures/device-monitor-2.png


Using MQTT devices
===============================================

MQTT devices run on MQTT protocol which works on subscribe and publish strategy, where the clients publish and subscribe to an MQTT Broker. All the subscribing clients are notified when another client publishes data on MQTT broker.

Mosquitto Broker is used for MQTT device simulation. Mosquitto broker allows data publishing and subscription on its uniquely identified resources called topics. These topics are defined in the format “/<apikey>/<device_id>/<topicSpecificPart>”. Users can track the updates on these topics by directly subscribing them on the host where Mosquitto is installed.

**Prerequisites for proceding further:**

* Install Mosquitto Broker.
* Start IoT Agent with MQTT Broker location pre-configured. For simplicity, add the following to the environment variables of IoT Agent JSON in the docker-compose file and then run the docker-compose. 

.. code-block:: console

      - IOTA_MQTT_HOST=<MQTT_Broker_Host_IP>
      - IOTA_MQTT_PORT=1883   # Mosquitto Broker runs at port 1883 by default.

In order to let IoT-Agent JSON allow both Northbound as well as Southbound data flow, users need to provide api-key as well for their device registration, so that the IoT-Agent can publish and subscribe to the topics using the api-key. For this, an extra Service-Provisioning request will be sent to IoT Agent. Steps to work with MQTT Devices in Fogflow are given below.


**Create a Service at IoT-Agent** using the following curl request.

.. code-block:: console

      curl -iX POST \
        'http://<IoT_Agent_IP>:4041/iot/services' \
        -H 'Content-Type: application/json' \
        -H 'fiware-service: iot' \
        -H 'fiware-servicepath: /' \
        -d '{
      "services": [
         {
           "apikey":      "FFNN1111",
           "entity_type": "Lamp",
           "resource":    "/iot/json"
         }
      ]
      }'


**Register a Lamp device** using the following curl request.

.. code-block:: console

      curl -X POST \
        http://<IoT_Agent_IP>:4041/iot/devices \
        -H 'content-type: application/json' \
        -H 'fiware-service: iot' \
        -H 'fiware-servicepath: /' \
        -d '{
        "devices": [
          {
            "device_id": "lamp001",
            "entity_name": "urn:ngsi-ld:Lamp:001",
            "entity_type": "Lamp",
            "protocol": "IoTA-JSON",
            "transport": "MQTT",
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


**Subscribe to Mosquitto topics:** Once service and device are successfully created, subscribe to the following topics of Mosquitto Broker in separate terminals to track what data are published on these topics:

.. code-block:: console

      mosquitto_sub -h <MQTT_Host_IP> -t "/FFNN1111/lamp001/attrs" 

.. code-block:: console

      mosquitto_sub -h <MQTT_Host_IP> -t "/FFNN1111/lamp001/cmd"
      

**Publish data to Thin Broker:** This section covers the northbound traffic. IoT Agent subscribes to some default topics like ["/+/+/attrs/+","/+/+/attrs","/+/+/configuration/commands","/+/+/cmdexe"]. So, in order to send attribute data to IoT Agent, data need to be published on a topic of Mosquitto Broker using the below command. 

.. code-block:: console

      mosquitto_pub -h <MQTT_Host_IP> -t "/FFNN1111/lamp001/attrs" -m '{"luminosity":78, "state": "ok"}'

Mosquitto broker will notify IoT-Agent for this Update, and consequently, the data will be updated at Thin Broker also.

The updated data can be viewed on the subscribed topic "/FFNN1111/lamp001/attrs" as well , as shown in the figure below.

.. figure:: figures/mqtt-data-update.png


**Run device commands:** This section covers the southbound traffic flow, i.e., how commands are run on the device. For this, send the below command updateContext request to Thin Broker. Thin broker will find the provider for this command update and will forward the UpdateContext request to that provider. In this case, IoT-Agent is the provider. IoT-Agent will publish the command at "/FFNN1111/lamp001/cmd" topic of the Mosquitto broker linked to it.

.. code-block:: console

      curl -iX POST \
      'http://<Thin_Broker_IP>:8070/ngsi10/updateContext' \
      -H 'Content-Type: application/json' \
      -H 'fiware-service: iot' \
      -H 'fiware-servicepath: /' \
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
                       "name": "on",
                       "type": "command",
                       "value": ""
                   }
               ]
          }
          ],
          "updateAction": "UPDATE"
      }'
      
The updated data can be viewed on the subscribed topic "/FFNN1111/lamp001/cmd", as shown in the figure below. This means that "on" command has been run successfully on the MQTT device.

.. figure:: figures/mqtt-cmd-update.png


Users can again have their customized devices to publish the command result on Thin Broker side.

Other APIs for RegisterContext
-----------------------------------------------

**GET a Registration**

Below is the curl request to get a device registration from a thin broker within Fogflow System, it will tell which broker contains the registration information regarding that device.

.. code-block:: console

      curl -iX GET \
      'http://<Thin_Broker_IP>:8070/NGSI9/registration/Lamp001' \
      -H 'fiware-service: openiot' \
      -H 'fiware-servicepath: /'

The device registration id for the above registration would be "Lamp001.openiot.~" within Fogflow. 

Users can also look for the registration at thin broker in the following way, as the Fiware Headers (i.e., "fiware-service" and "fiware-servicepath") are optional in the request. The result is completely dependent on what is being searched for.

.. code-block:: console

      curl -iX GET \
      'http://<Thin_Broker_IP>:8070/NGSI9/registration/Lamp001.openiot.~'


**DELETE a Registration**

Following curl request would delete a device registration in Fogflow.

.. code-block:: console

      curl -iX DELETE \
      'http://<Thin_Broker_IP>:8070/NGSI9/registration/Lamp001' \
      -H 'fiware-service: openiot' \
      -H 'fiware-servicepath: /'

This request would delete the registration "Lamp001.openiot.~". Fiware Headers (i.e., "fiware-service" and "fiware-servicepath") are mandatory.
