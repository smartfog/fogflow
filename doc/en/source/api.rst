*****************************************
APIs and examples of their usage 
*****************************************

APIs of FogFlow Discovery
===================================

lookup of nearby brokers
-----------------------------------------------

For any external application or IoT devices, the only interface they need from FogFlow Discovery is to find out a nearby 
Broker based on its own location. After that, they only need to interact with the assigned nearby Broker. 

**POST /ngsi9/discoverContextAvailability**

==============   ===============
Param            Description
==============   ===============
latitude         latitude of your location
longitude        latitude of your location
limit            number of expected brokers
==============   ===============


Please check the following examples. 

.. note:: For the Javascript code example, you need the library ngsiclient.js. 
    Please refer to the code repository at application/device/powerpanel

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                   "entities":[
                      {
                         "type":"IoTBroker",
                         "isPattern":true
                      }
                   ],
                   "restriction":{
                      "scopes":[
                         {
                            "type":"nearby",
                            "value":{
                               "latitude":35.692221,
                               "longitude":139.709059,
                               "limit":1
                            }
                         }
                      ]
                   }
                }             


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        
        var discoveryURL = "http://localhost:8071/ngsi9";
        var myLocation = {
                "latitude": 35.692221,
                "longitude": 139.709059
            };
        
        // find out the nearby IoT Broker according to my location
        var discovery = new NGSI.NGSI9Client(discoveryURL)
        discovery.findNearbyIoTBroker(myLocation, 1).then( function(brokers) {
            console.log('-------nearbybroker----------');    
            console.log(brokers);    
            console.log('------------end-----------');    
        }).catch(function(error) {
            console.log(error);
        });

  
       

APIs of FogFlow Broker
===============================

Create and update a context entity
-----------------------------------------------

**POST /ngsi10/updateContext**

==============   ===============
Param            Description
==============   ===============
latitude         latitude of your location
longitude        latitude of your location
limit            number of expected brokers
==============   ===============

Example: 

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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


Fetch a context entity by ID
-----------------------------------------------

**GET /ngsi10/entity/#eid**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example: 

.. code-block:: console 

   curl http://localhost:8070/ngsi10/entity/test001

Check all context entities on a single Broker
-------------------------------------------------

**GET /ngsi10/entity**

Example: 

.. code-block:: console 

    curl http://localhost:8070/ngsi10/entity


Delete a specific context entity by ID
-----------------------------------------------

**DELETE /ngsi10/entity/#eid**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example: 

.. code-block:: console 

    curl -iX DELETE http://localhost:8070/ngsi10/entity/test001


Check all received subscriptions
-------------------------------------------------

**GET /ngsi10/subscription**

Example: 

.. code-block:: console 

    curl http://localhost:8070/ngsi10/subscription


Query context entities
-----------------------------------------------

**GET /ngsi10/queryContext**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example1: **query context by entity type**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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


Example2: **query context by geo-scope**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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


Example3: **query context by combination of several constraints**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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



Subscribe context entities
-----------------------------------------------

**GET /ngsi10/queryContext**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example1: **subscribe context by entity type**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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


Example2: **subscribe context by geo-scope**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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


Example3: **subscribe context by combination of several constraints**

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "contextValue": 10
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


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
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







