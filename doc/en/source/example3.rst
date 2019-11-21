*****************************************
Connect an IoT device to FogFlow
*****************************************

If your device can communicate with FogFlow via NGSI, connecting your device to FogFlow
can be very easy. It requires some small application to be running on your device,
for example, a raspberry Pi with several connected sensors or actuators. 

In the following example, we show how a simulated PowerPanel device can be connected to FogFlow via NGSI. 
This example code is also accessible from `FogFlow code repository`_ in the application folder. 

You need Node.js to run this example code. Please install Node.js and npm.

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


You need to modify discoveryURL in profile1.json.

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


You need to install the packages as follows:

.. code-block:: console

    npm install


Run this example code as follows:

.. code-block:: console

    node powerpanel.js profile1.json


To connect Non-NGSI IoT Devices, FIWARE provides IoT Agents that work with IoT devices based on various protocols like MQTT, Ultralight, 
LoRaWAN, etc. IoT Agents can communicate over both, either NGSIv1 or NGSIv2, however, currently Fogflow supports only NGSIv1.
For IoT Agents, following two scenarios can be there:

- When IoT Agent uses NGSIv1, Fogflow can directly understand IoT Agent requests.
- When IoT Agent uses NGSIv2, `General Purpose Adapter`_ is required as an intermediater between IoT Agent and Fogflow.

Users can run IoT Agent on cloud node by directly running `docker-compose`_ file used to start the cloud node. By default, IoT Agent is 
already allowed. Users can opt out if they do not require it.

For running IoT Agent on edge node, users can uncomment the related command in `Start Edge`_ file.

.. _`General Purpose Adapter`: https://github.com/smartfog/fogflow/tree/master/application/operator/GeneralPurposeAdapter
   
.. _`docker-compose`: https://github.com/smartfog/fogflow/blob/master/docker/core/http/docker-compose.yml

.. _`Start Edge`: https://github.com/smartfog/fogflow/blob/master/docker/edge/http/start.sh
