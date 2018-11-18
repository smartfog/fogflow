*****************************************
APIs and examples of their usage 
*****************************************

APIs of FogFlow Discovery
===================================

Look up nearby brokers
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

.. figure:: https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/json/petstore-expanded.json.svg
  :target: https://app.swaggerhub.com/apis/fogflow/broker/1.0.0


Create/update context
-----------------------------------------------

.. note:: It is the same API to create or update a context entity. 
    For a context update, if there is no existing entity, a new entity will be created. 


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
              'http://localhost:8070/ngsi10/updateContext' \
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
                            },{
                                "name": "city",
                                "type": "string",
                                "value": "Heidelberg"                             
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


Query Context via GET
-----------------------------------------------


Fetch a context entity by ID
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity/#eid**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example: 

.. code-block:: console 

   curl http://localhost:8070/ngsi10/entity/Device.temp001

Fetch a specific attribute of a specific context entity
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity/#eid/#attr**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
attr             specify the attribute name to be fetched
==============   ===============

Example: 

.. code-block:: console 

   curl http://localhost:8070/ngsi10/entity/Device.temp001/temp


Check all context entities on a single Broker
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity**

Example: 

.. code-block:: console 

    curl http://localhost:8070/ngsi10/entity



Query context via POST
-----------------------------------------------

**POST /ngsi10/queryContext**

==============   ===============
Param            Description
==============   ===============
entityId         specify the entity filter, which can define a specific entity ID, ID pattern, or type
restriction      a list of scopes and each scope defines a filter based on domain metadata
==============   ===============

query context by the pattern of entity ID
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"id":"Device.*","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{id:'Device.*', isPattern: true}];           
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });          


query context by entity type
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"type":"Temperature","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'Temperature', isPattern: true}];           
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });          


query context by geo-scope (circle)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]};
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


query context by geo-scope (polygon)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
               "entities":[
                  {
                     "id":".*",
                     "isPattern":true
                  }
               ],
               "restriction":{
                  "scopes":[
                     {
                        "scopeType":"polygon",
                        "scopeValue":{
                           "vertices":[
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":135.84594726562503
                              },
                              {
                                 "latitude":37.18657859524883,
                                 "longitude":135.84594726562503
                              },
                              {
                                 "latitude":37.18657859524883,
                                 "longitude":141.51489257812503
                              },
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":141.51489257812503
                              },
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":135.84594726562503
                              }
                           ]
                        }
                    }]
                }
            }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {
               "scopes":[
                  {
                     "scopeType":"polygon",
                     "scopeValue":{
                        "vertices":[
                           {
                              "latitude":34.4069096565206,
                              "longitude":135.84594726562503
                           },
                           {
                              "latitude":37.18657859524883,
                              "longitude":135.84594726562503
                           },
                           {
                              "latitude":37.18657859524883,
                              "longitude":141.51489257812503
                           },
                           {
                              "latitude":34.4069096565206,
                              "longitude":141.51489257812503
                           },
                           {
                              "latitude":34.4069096565206,
                              "longitude":135.84594726562503
                           }
                        ]
                     }
                  }
               ]
            }
                    
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


query context with the filter of domain metadata values
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. note:: the conditional statement can be defined only with the domain matadata of your context entities
    For the time being, it is not supported to filter out entities based on specific attribute values. 

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};        
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


query context with multiple filters
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};          
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


Delete context
-----------------------------------------------

Delete a specific context entity by ID
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi10/entity/#eid**

==============   ===============
Param            Description
==============   ===============
eid              entity ID
==============   ===============

Example: 

.. code-block:: console 

    curl -iX DELETE http://localhost:8070/ngsi10/entity/Device.temp001






Subscribe context
-----------------------------------------------

**POST /ngsi10/subscribeContext**

==============   ===============
Param            Description
==============   ===============
entityId         specify the entity filter, which can define a specific entity ID, ID pattern, or type
restriction      a list of scopes and each scope defines a filter based on domain metadata
reference        the destination to receive notifications
==============   ===============

subscribe context by the pattern of entity ID
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"id":"Device.*","isPattern":true}],
                    "reference": "http://localhost:8066"
                }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        var mySubscriptionId;
        
        var subscribeReq = {}
        subscribeReq.entities = [{id:'Device.*', isPattern: true}];           
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });

subscribe context by entity type
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"type":"Temperature","isPattern":true}]
                    "reference": "http://localhost:8066"                    
                  }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'Temperature', isPattern: true}];           
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });       


subscribe context by geo-scope
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066",                    
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]};
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });   

subscribe context with the filter of domain metadata values
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. note:: the conditional statement can be defined only with the domain matadata of your context entities
    For the time being, it is not supported to filter out entities based on specific attribute values. 

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066",                    
                    "restriction": {
                        "scopes": [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};        
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });      


subscribe context with multiple filters
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066", 
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};          
        
        // use the IP and Port number your receiver is listening
        subscribeReq.reference =  'http://' + agentIP + ':' + agentPort;  
        
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });   

Cancel a subscription by subscription ID
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi10/subscription/#sid**


==============   ===============
Param            Description
==============   ===============
sid              the subscription ID created when the subscription is issued
==============   ===============


curl -iX DELETE http://localhost:8070/ngsi10/subscription/#sid



