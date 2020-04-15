*****************************************
Integration
*****************************************

Connect an IoT device to FogFlow
====================================

With NGSI supported Devices
--------------------------------

If the device can communicate with FogFlow via NGSI, connecting device to FogFlow
can be very easy. It requires some small application to be running on the device,
for example, a raspberry Pi with several connected sensors or actuators. 

In the following example, it is shown how a simulated PowerPanel device can be connected to FogFlow via NGSI. 
This example code is also accessible from `FogFlow code repository`_ in the application folder. 

Node.js need to be run this example code. Please install Node.js and npm.

.. _`FogFlow code repository`: https://github.com/smartfog/fogflow/blob/master/application/device/powerpanel/powerpanel.js

.. code-block:: javascript

    'use strict';
    
    const NGSI = require('./ngsi/ngsiclient.js');
    const fs = require('fs');
    
    // read device profile from the configuration file
    var args = process.argv.slice(2);
    if(args.length != 1){
        console.log('please specify the device profile');
        return;
    }
    
    var cfgFile = args[0];
    var profile = JSON.parse(
        fs.readFileSync(cfgFile)
    );
    
    var ngsi10client;
    var timer;
    
    // find out the nearby IoT Broker according to my location
    var discovery = new NGSI.NGSI9Client(profile.discoveryURL)
    discovery.findNearbyIoTBroker(profile.location, 1).then( function(brokers) {
        console.log('-------nearbybroker----------');    
        console.log(brokers);    
        console.log('------------end-----------');    
        if(brokers && brokers.length > 0) {
            ngsi10client = new NGSI.NGSI10Client(brokers[0]);
    
            // generating data observations periodically
            timer = setInterval(function(){ 
                updateContext();
            }, 1000);    
    
            // register my device profile by sending a device update
            registerDevice();
        }
    }).catch(function(error) {
        console.log(error);
    });
    
    // register device with its device profile
    function registerDevice() 
    {
        var ctxObj = {};
        ctxObj.entityId = {
            id: 'Device.' + profile.type + '.' + profile.id,
            type: profile.type,
            isPattern: false
        };
        
        ctxObj.attributes = {};
        
        var degree = Math.floor((Math.random() * 100) + 1);        
        ctxObj.attributes.usage = {
            type: 'integer',
            value: degree
        };   
        ctxObj.attributes.shop = {
            type: 'string',
            value: profile.id
        };       
        ctxObj.attributes.iconURL = {
            type: 'string',
            value: profile.iconURL
        };                   
        
        ctxObj.metadata = {};
        
        ctxObj.metadata.location = {
            type: 'point',
            value: profile.location
        };    
       
        ngsi10client.updateContext(ctxObj).then( function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to update context');
        });  
    }
    
    // update context for streams
    function updateContext() 
    {
        var ctxObj = {};
        ctxObj.entityId = {
            id: 'Stream.' + profile.type + '.' + profile.id,
            type: profile.type,
            isPattern: false
        };
        
        ctxObj.attributes = {};
        
        var degree = Math.floor((Math.random() * 100) + 1);        
        ctxObj.attributes.usage = {
            type: 'integer',
            value: degree
        };
        ctxObj.attributes.deviceID = {
            type: 'string',
            value: profile.type + '.' + profile.id
        };   	     
        
        ctxObj.metadata = {};
        
        ctxObj.metadata.location = {
            type: 'point',
            value: profile.location
        }; 
        ctxObj.metadata.shop = {
            type: 'string',
            value: profile.id
        };	          
        
        ngsi10client.updateContext(ctxObj).then( function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to update context');
        });    
    }
    
    process.on('SIGINT', function() 
    {    
        if(ngsi10client) {
            clearInterval(timer);
            
            // to delete the device
            var entity = {
                id: 'Device.' + profile.type + '.' + profile.id,
                type: 'Device',
                isPattern: false
            };
            ngsi10client.deleteContext(entity).then( function(data) {
                console.log(data);
            }).catch(function(error) {
                console.log('failed to delete context');
            });        
    
            // to delete the stream    
            var entity = {
                id: 'Stream.' + profile.type + '.' + profile.id,
                type: 'Stream',
                isPattern: false
            };
            ngsi10client.deleteContext(entity).then( function(data) {
                console.log(data);
            }).catch(function(error) {
                console.log('failed to delete context');
            });        
        }
    });


discoveryURL is need to modify in profile1.json.

.. code-block:: json

    {
        "discoveryURL":"http://35.198.104.115:443/ngsi9",
        "location": {
            "latitude": 35.692221,
            "longitude": 139.709059
        },
        "iconURL": "/img/shop.png",
        "type": "PowerPanel",
        "id": "01"
    }


 Packages that need to be installed as follows:

.. code-block:: console

    npm install


Run this example code as follows:

.. code-block:: console

    node powerpanel.js profile1.json

With Non-NGSI supported Devices
----------------------------------

To connect Non-NGSI IoT Devices, FIWARE provides IoT Agents that work with IoT devices based on various protocols like MQTT, Ultralight,
etc. IoT Agents can communicate over both, either NGSIv1 or NGSIv2, however, currently Fogflow supports only NGSIv1. So, users need to configure IoT Agent to use NGSIv1 format.

Users can run IoT Agent on Fogflow cloud node by directly running `docker-compose`_ file used to start the cloud node. By default, IoT Agent is already allowed. Users can opt out if they do not require it.

For running IoT Agent on edge node, users can uncomment the related command in `Start Edge`_ file.
   
.. _`docker-compose`: https://github.com/smartfog/fogflow/blob/master/docker/core/http/docker-compose.yml

.. _`Start Edge`: https://github.com/smartfog/fogflow/blob/master/docker/edge/http/start.sh


An example usage of Fiware IoT-Agent JSON sending location-based temerature data to thin broker is given below. Iot Agent requires following three requests for sending NGSI Data to broker.

- **Service Provisioning:** Service provisioning or group provisioning is used by IoT Agent to set some default commands or attributes like authentication key, optional context broker endpoint, etc. for anonymous devices.

Following is the curl request for creating or registring a service on IoT Agent.

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
         "entity_type": "Thing",
         "resource":    "/iot/json"
       }
    ]
    }'

- **Device Provisioning:** Device provisioning is used to specify what data and data attributes a device will be sending to the IoT Agent.

The below curl request is used to register a device having Device ID "Device1111" which would be sending the data of entity "Thing1111" to IoT Agent.

.. code-block:: console

    curl -X POST \
      http://<IoT_Agent_IP>:4041/iot/devices \
      -H 'content-type: application/json' \
      -H 'fiware-service: iot' \
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
                            "type": "integer"
                    }
                    ]}]
    }'

- **Sensor Data Updation:** IoT Agent maps the received data with its device registration and creates an NGSI update corresponding to the same. Note that IoT Agent receives data from Device in Non-NGSI format.

Curl request that actually sends the "Thing1111" entity update to IoT Agent on behalf of "Device1111" is given below.

.. code-block:: console

    curl -X POST \
      'http://<IoT_Agent_IP>:7896/iot/json?i=Device1111&k=FFNN1111' \
      -H 'content-type: application/json' \
      -H 'fiware-service: iot' \
      -H 'fiware-servicepath: /' \
      -d '{ 
        "locationName":"Heidelberg",
        "locationId":"0011",
        "Temperature":20
    }'

As soon as the IoT Agent recieves update from device, it requests thin broker to update the entity data in the form of an NGSIv1 UpdateContext request.
