*****************************************
API Walkthrough 
*****************************************

FogFlow Discovery API
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

.. note:: For the Javascript code example, library ngsiclient.js is needed.
    Please refer to the code repository at application/device/powerpanel

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:80/ngsi9/discoverContextAvailability' \
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
                            "scopeType":"nearby",
                            "scopeValue":{
                               "latitude":35.692221,
                               "longitude":139.709059,
                               "limit":1
                            }
                         }
                      ]
                   }
                } '            


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        
        var discoveryURL = "http://localhost:80/ngsi9";
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

  
       

FogFlow Broker API
===============================

.. figure:: https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/json/petstore-expanded.json.svg
  :target: https://app.swaggerhub.com/apis/fogflow/broker/1.0.0

.. note:: Use port 80 for accessing the cloud broker, whereas for edge broker, the default port is 8070.


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
              'http://localhost:80/ngsi10/updateContext' \
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
                              "value": 10
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
        var brokerURL = "http://localhost:80/ngsi10"
    
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

   curl http://localhost:80/ngsi10/entity/Device.temp001

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

   curl http://localhost:80/ngsi10/entity/Device.temp001/temp


Check all context entities on a single Broker
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity**

Example: 

.. code-block:: console 

    curl http://localhost:80/ngsi10/entity



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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"id":"Device.*","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"type":"Temperature","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/queryContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

    curl -iX DELETE http://localhost:80/ngsi10/entity/Device.temp001






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

            curl -X POST 'http://localhost:80/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"id":"Device.*","isPattern":true}],
                    "reference": "http://localhost:8066"
                }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"type":"Temperature","isPattern":true}]
                    "reference": "http://localhost:8066"                    
                  }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/subscribeContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/subscribeContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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

            curl -X POST 'http://localhost:80/ngsi10/subscribeContext' \
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
        var brokerURL = "http://localhost:80/ngsi10"    
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


curl -iX DELETE http://localhost:80/ngsi10/subscription/#sid


FogFlow Designer API
===================================

FogFlow uses its own REST APIs to manage all internal objects, 
including operator, docker image, service topology, service intent, and fog function. 
In addition, FogFlow also provides the extra interface for device registration and 
the management of subscriptions to exchange data with other FIWARE brokers, 
such as Orion/Orion-LD and Scorpio, 
which could be used both as the data sources to fetch the original data 
or as the destination to publish the generated results. 


Operator
-------------------


**a. To create a new Operator**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /operator**

**Example**   
 
.. code-block:: console

	curl -X POST \
	  'http://127.0.0.1:8080/operator' \
	  -H 'Content-Type: application/json' \
	  -d '
	    [{
	        "name": "dummy",
	        "description": "test",
	        "parameters": []
	    }]
	    '


**b. To retrieve all the operators**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /operator**

**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/operator'

**c. To retrieve a specific operator based on operator name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /operator/<name>**

==============   ============================
Param		 Description
==============   ============================
name             Name of existing operator
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/operator/dummy'

DockerImage
-------------------


**a. To create a new DockerImage**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /dockerimage**



**Example**   
 
.. code-block:: console

	curl -X POST \
	  'http://127.0.0.1:8080/dockerimage' \
	  -H 'Content-Type: application/json' \
	  -d '
	    [
	        {
	            "name": "fogflow/dummy",
	            "hwType": "X86",
	            "osType": "Linux",
	            "operatorName": "dummy",
	            "prefetched": false,
	            "tag": "latest"
	        }
	    ]
	    '



**b. To retrieve all the DockerImage**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /dockerimage**

**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/dockerimage'


**c. To retrieve a specific DockerImage based on operator name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /dockerimage/<operator name>**

==============   ============================
Param		 Description
==============   ============================
name             Name of existing operator
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/dockerimage/dummy'


Service Topology
-------------------


**a. To create a new Service**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /service**


**Example**   
 
.. code-block:: console
	
	curl -X POST \
	  'http://127.0.0.1:8080/service' \
	  -H 'Content-Type: application/json' \
	  -d '
	    [
	        {
	            "topology": {
	                "name": "MyTest",
	                "description": "a simple case",
	                "tasks": [
	                    {
	                        "name": "main",
	                        "operator": "dummy",
	                        "input_streams": [
	                            {
	                                "selected_type": "Temperature",
	                                "selected_attributes": [],
	                                "groupby": "EntityID",
	                                "scoped": false
	                            }
	                        ],
	                        "output_streams": [
	                            {
	                                "entity_type": "Out"
	                            }
	                        ]
	                    }
	                ]
	            }
	        }
	    ]
	    '



**b. To retrieve all the Service**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /service**

**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/service'

**c. To retrieve a specific service based on service name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /service/<service name>**

==============   ============================
Param		      Description
==============   ============================
name              Name of existing service
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/service/MyTest'

   

**d. To delete a specific service based on service name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**DELETE /service/<service name>**

==============   ============================
Param		 		Description
==============   ============================
name              Name of existing service
==============   ============================


**Example:**

.. code-block:: console

   curl -X DELETE  'http://localhost:8080/service/MyTest' \
   -H 'Content-Type: application/json'



Intent
-------------------


**a. To create a new Intent**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /intent**

**Example**   
 
.. code-block:: console

	curl -X POST \
	  'http://127.0.0.1:8080/intent' \
	  -H 'Content-Type: application/json' \
	  -d '
	    {
	        "id": "ServiceIntent.594e3d10-59f9-4ee6-97be-fe50b9c99bd8",        
	        "topology": "MyTest",
	        "stype": "ASYN",
	        "priority": {
	            "exclusive": false,
	            "level": 0
	        },
	        "qos": "NONE",
	        "geoscope": {
	            "scopeType": "global",
	            "scopeValue": "global"
	        }
	    }
    '


**b. To retrieve all the Intent**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /intent**

**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/intent'


**c. To retrieve a specific Intent based on Intent ID**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /intent/<intent id>**

==============   ============================
Param		       Description
==============   ============================
id                ID of existing intent
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/intent/ServiceIntent.594e3d10-59f9-4ee6-97be-fe50b9c99bd8'

**d. To delete a specific Intent based on intent id**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**DELETE /intent/<intent id>**

==============   ============================
Param		      Description
==============   ============================
id                ID of the existing intent
==============   ============================


**Example:**

.. code-block:: console

curl -iX DELETE \
  'http://127.0.0.1:8080/intent/ServiceIntent.594e3d10-59f9-4ee6-97be-fe50b9c99bd8'
   

**e. To retrieve list service intents for the given service topology**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /intent/topology/<TopologeName>**

==============   ======================================
Param		       Description
==============   ======================================
TopologeName       name of the given service topology 
==============   ======================================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/intent/topology/MyTest'

Topology
-------------------

**a. To retrieve all the Topology**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /topology**

**Example:**

.. code-block:: console

   curl -X GET 'http://localhost:8080/topology' \
  -H 'Content-Type: application/json'


**b. To retrieve a specific Topology based on topology name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /topology/<topology name>**

==============   ============================
Param		 Description
==============   ============================
Name             name of the existing Topology
==============   ============================	

**Example:** 

.. code-block:: console

   curl -X GET  'http://localhost:8080/topology/MyTest' \
   -H 'Content-Type: application/json' 




Fog Function
-------------------


**a. To create a new Fogfunction**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /fogfunction**



**Example**   
 
.. code-block:: console

	
	curl -X POST \
	  'http://127.0.0.1:8080/fogfunction' \
	  -H 'Content-Type: application/json' \
	  -d '
		[
		    {
		        "name": "ffTest",
		        "topology": {
		            "name": "ffTest",
		            "description": "a fog function",
		            "tasks": [
		                {
		                    "name": "main",
		                    "operator": "mqtt-adapter",
		                    "input_streams": [
		                        {
		                            "selected_type": "HOPU",
		                            "selected_attributes": [],
		                            "groupby": "EntityID",
		                            "scoped": false
		                        }
		                    ],
		                    "output_streams": [
		                        {
		                            "entity_type": "Out"
		                        }
		                    ]
		                }
		            ]
		        },
		        "intent": {
		            "id": "ServiceIntent.1c6396bb-281d-4c14-b61d-f0cc0dcc1006",
		            "topology": "ffTest",
		            "priority": {
		                "exclusive": false,
		                "level": 0
		            },
		            "qos": "default",
		            "geoscope": {
		                "scopeType": "global",
		                "scopeValue": "global"
		            }
		        },
		        "status": "enabled"
		    }
		]
	    '




**b. To retrieve all the Fogfunction**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /fogfunction**

**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/fogfunction'


**c. To retrieve a specific Fogfunction based on fogfunction name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /fogfunction/<name>**

==============   ============================
Param		 Description
==============   ============================
name             Name of existing fogfunction
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/fogfunction/ffTest'

**d. To delete a specific fogfunction based on fogfunction name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**DELETE /fogfunction/<fogfunction name>**

==============   ============================
Param		 Description
==============   ============================
name              Name of existing fogfunction
==============   ============================


**Example:**

.. code-block:: console

	curl -iX DELETE \
	  'http://127.0.0.1:8080/fogfunction/ffTest'


**e. To enable a specific fogfunction based on fogfunction name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**GET /fogfunction/<fogfunction name>/enable**

==============   ============================
Param		 Description
==============   ============================
name              Name of existing fogfunction
==============   ============================


**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/fogfunction/ffTest/enable'


**F. To disable a specific fogfunction based on fogfunction name**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**GET /fogfunction/<fogfunction name>/disable**

==============   ============================
Param		 Description
==============   ============================
name              Name of existing fogfunction
==============   ============================


**Example:**

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/fogfunction/ffTest/disable'



Device
-------------------

**a. To create a new device**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /device**



**Example**   
 
.. code-block:: console

	curl -iX POST \
	  'http://127.0.0.1:8080/device' \
	  -H 'Content-Type: application/json' \
	  -d '
	    {
	        "id": "urn:Device.12345",
	        "type": "HOPU",
	        "attributes": {
	            "protocol": {
	                "type": "string",
	                "value": "MQTT"
	                },
	            "mqttbroker":  {
	                "type": "string",
	                "value": "mqtt://mqtt.cdtidev.nec-ccoc.com:1883"
	                },
	            "topic": {
	                "type": "string",
	                "value": "/api/12345/attrs"
	                },
	            "mappings": {
	                "type": "object",
	                "value": {
	                    "temp8": {
	                        "name": "temperature",
	                        "type": "float",
	                        "entity_type": "AirQualityObserved"
	                    },
	                    "hum8": {
	                        "name": "humidity",
	                        "type": "float",
	                        "entity_type": "AirQualityObserved"
	                    }
	                }
	            }          
	        },
	        "metadata": {
	            "location": {
	                    "type": "point",
	                    "value": {
	                        "latitude": 49.406393,
	                        "longitude": 8.684208
	                    }
	            }
	        }
	    }'


**b. To get the list of all registered devices**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /device**

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/device'


**c. To delete a specific device**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /device/<device_id>** 

==============   ============================
Param		      Description
==============   ============================
device_id         entity ID of this device
==============   ============================	

**Example:** 

.. code-block:: console

	curl -iX DELETE \
	  'http://127.0.0.1:8080/device/urn:Device.12345'



Subscription
-------------------

**a. To create a subscription for a given destination**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /subscription**


**Example**   
 
.. code-block:: console

	curl -iX POST \
	  'http://127.0.0.1:8080/subscription' \
	  -H 'Content-Type: application/json' \
	  -d '{
	            "entity_type": "AirQualityObserved",
	            "destination_broker": "NGSI-LD",
	            "reference_url": "http://127.0.0.1:9090",
	            "tenant": "ccoc"
	      }'


**b. To get the list of all registered subscription**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /subscription**

**Example:** 

.. code-block:: console

	curl -iX GET \
	  'http://127.0.0.1:8080/subscription'


**c. To delete a specific subscription**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /subscription/<subscription_id>** 

=================   ============================
Param                 Description
=================   ============================
subscription_id      ID of this subscription
=================   ============================


**Example:** 

.. code-block:: console


	curl -iX DELETE \
	  'http://127.0.0.1:8080/subscription/88bba05c-dda2-11ec-ba1d-acde48001122'


NGSI-LD Supported API's
============================

The following figure shows a brief overview of how the APIs in current scope will be used to achieve the goal of NGSI-LD API support in FogFlow. The API support includes Entity creation, registration, subscription and notification.



.. figure:: figures/ngsild_architecture.png

Entities API
------------
For the purpose of interaction with Fogflow, IOT devices approaches broker with entity creation request where it is resolved as per given context. Broker further forwards the registration request to Fogflow Discovery in correspondence to the created entity.

.. note:: Use port 80 for accessing the cloud broker, whereas for edge broker, the default port is 8070. The localhost is the coreservice IP for the system hosting fogflow. 

**POST /ngsi-ld/v1/entities/**

**a. To create NGSI-LD context entity, with context in Link in Header**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

=============     ===========================================================
key               Value
=============     ===========================================================
Content-Type      application/json
Accept            application/ld+json
Link              <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; 
                  type="application/ld+json"
=============     ===========================================================

**Request**

.. code-block:: console

   curl -iX POST \
     'http://localhost:80/ngsi-ld/v1/entities/' \
      -H 'Content-Type: application/json' \
      -H 'Accept: application/ld+json' \
      -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
        {
	      "id": "urn:ngsi-ld:Vehicle:A100",
	      "type": "Vehicle",
	      "brandName": {
		             "type": "Property",
		             "value": "Mercedes"
	       },
	       "isParked": {
		             "type": "Relationship",
		             "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
		             "observedAt": "2017-07-29T12:00:04",
		             "providedBy": {
			                     "type": "Relationship",
			                     "object": "urn:ngsi-ld:Person:Bob"
		              }
	        },
	        "speed": {
		           "type": "Property",
		           "value": 80
	         },
	        "createdAt": "2017-07-29T12:00:04",
	        "location": {
		               "type": "GeoProperty",
		               "value": {
			                 "type": "Point",
			                 "coordinates": [-8.5, 41.2]
		               }
	        }   
        }'
	
**b.  To create a new NGSI-LD context entity, with context in Link header and request payload is already expanded**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

=============     ======================================
key               Value
=============     ======================================
Content-Type      application/json
Accept            application/ld+json
=============     ======================================

**Request**

.. code-block:: console

      curl -iX POST \
     'http://localhost:80/ngsi-ld/v1/entities/' \
      -H 'Content-Type: application/json' \
      -H 'Accept: application/ld+json' \
      -d'
       {
                  "http://example.org/vehicle/brandName": [
                  {
                       "@type": [
                                   "http://uri.etsi.org/ngsi-ld/Property"
                        ],
                        "http://uri.etsi.org/ngsi-ld/hasValue": [
                                 {
                                      "@value": "Mercedes"
                                 }
                           ]
                     }
               ],
                 "http://uri.etsi.org/ngsi-ld/createdAt": [
                  {
                       "@type": "http://uri.etsi.org/ngsi-ld/DateTime",
                       "@value": "2017-07-29T12:00:04"
                   }
               ],
                 "@id": "urn:ngsi-ld:Vehicle:A8866",
                 "http://example.org/common/isParked": [
                  {
                             "http://uri.etsi.org/ngsi-ld/hasObject": [
                              {
                                      "@id": "urn:ngsi-ld:OffStreetParking:Downtown1"
                               }
                            ],
                             "http://uri.etsi.org/ngsi-ld/observedAt": [
                              {
                                     "@type": "http://uri.etsi.org/ngsi-ld/DateTime",
                                     "@value": "2017-07-29T12:00:04"
                               }
                            ],
                              "http://example.org/common/providedBy": [
                               {
                                        "http://uri.etsi.org/ngsi-ld/hasObject": [
                                        {
                                                "@id": "urn:ngsi-ld:Person:Bob"
                                        }
                                     ],
                                     "@type": [
                                                 "http://uri.etsi.org/ngsi-ld/Relationship"
                                       ]
                                }
                             ],
                               "@type": [
                                           "http://uri.etsi.org/ngsi-ld/Relationship"
                                 ]
                       }
                 ],
                  "http://uri.etsi.org/ngsi-ld/location": [
                   {
                             "@type": [
                                         "http://uri.etsi.org/ngsi-ld/GeoProperty"
                               ],
                             "http://uri.etsi.org/ngsi-ld/hasValue": [
                              {
                                    "@value": "{ \"type\":\"Point\", \"coordinates\":[ -8.5, 41.2 ] }"
                               }
                             ]
                    }
                 ],
                  "http://example.org/vehicle/speed": [
                   {
                            "@type": [
                                        "http://uri.etsi.org/ngsi-ld/Property"
                             ],
                             "http://uri.etsi.org/ngsi-ld/hasValue": [
                              {
                                    "@value": 80
                               } 
                             ]
                     }
                 ],
                  "@type": [
                             "http://example.org/vehicle/Vehicle"
                 ]

        }'

**c. To append additional attributes to an existing entity**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /ngsi-ld/v1/entities/**

=============     ======================================
key               Value
=============     ======================================
Content-Type      application/json
Accept            application/ld+json
=============     ======================================

**Request**

.. code-block:: console

       curl -iX POST \
       'http://localhost:80/ngsi-ld/v1/entities/' \
       -H 'Content-Type: application/json' \
       -H 'Accept: application/ld+json' \
       -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \       
       -d'
        {
	      "id": ""urn:ngsi-ld:Vehicle:A100",
              "type": "Vehicle",

	     ""brandName1"": {
		                 "type": "Property",
		                 "value": "BMW"
	      }
        }'

**d. To update specific attributes of an existing entity**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /ngsi-ld/v1/entities/**

=============     ======================================
key               Value
=============     ======================================
Content-Type      application/json
Accept            application/ld+json
=============     ======================================

**Request**

.. code-block:: console

        curl -iX POST \
       'http://localhost:80/ngsi-ld/v1/entities/' \
       -H 'Content-Type: application/json' \
       -H 'Accept: application/ld+json' \
       -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
       -d'
        {
		"id": ""urn:ngsi-ld:Vehicle:A100",
	        "type": "Vehicle",

	       "brandName": {
		                  "type": "Property",
		                  "value": "AUDI"
	        }
	}'


**e. To delete an NGSI-LD context entity**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi-ld/v1/entities/#eid**

==============   ============================
Param		 Description
==============   ============================
eid              Entity Id
==============   ============================

**Example:**

.. code-block:: console

   curl -iX DELETE http://localhost:80/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100  -H 'Content-Type: application/json' -H 'Accept: application/ld+json'


**f. To delete an attribute of an NGSI-LD context entity**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi-ld/v1/entities/#eid/attrs/#attrName**

==============   ============================
Param		 Description
==============   ============================
eid              Entity Id
attrName         Attribute Name
==============   ============================

**Example:**

.. code-block:: console

   curl -iX DELETE http://localhost:80/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs/brandName1

**g. To retrieve a specific entity**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/entities/#eid**

==============   ============================
Param		 Description
==============   ============================
eid              Entity Id
==============   ============================

**Example:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4569
   
Subscription API
-------------------

A new subscription is issued by the subscriber which is enrouted to broker where the details of subscriber is stored for notification purpose. The broker initiate a request to Fogflow Discovery, where this is registered as new subscription and looks for availabltiy of corresponding data. On receiving data is passes the information back to subscribing broker.

**a. To create a new Subscription to with context in Link header**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**POST /ngsi-ld/v1/subscriptions**

**Header Format**

=============     ===========================================================
key               Value
=============     ===========================================================
Content-Type      application/ld+json
Link              <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; 
                  type="application/ld+json"
=============     ===========================================================

**Request**   

.. code-block:: console

   curl -iX POST\
     'http://localhost:80/ngsi-ld/v1/subscriptions/' \
      -H 'Content-Type: application/ld+json' \
      -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
       {
		"type": "Subscription",
		"id"  : "urn:ngsi-ld:Subscription:71",
		"entities": [{
				"id": "urn:ngsi-ld:Vehicle:71",
				"type": "Vehicle"
		}],
		"watchedAttributes": ["brandName"],
		"notification": {
				"attributes": ["brandName"],
				"format": "keyValues",
				"endpoint": {
						"uri": "http://my.endpoint.org/notify",
						"accept": "application/json"
				  }
		}
	 }'



**b. To retrieve all the subscriptions**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/subscriptions**

**Example:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/subscriptions/ -H 'Accept: application/ld+json'


**c. To retrieve a specific subscription based on subscription id**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/subscriptions/#sid**

==============   ============================
Param		 Description
==============   ============================
sid              subscription Id
==============   ============================	

**Example:** 

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:71
   
**d. To delete a specific subscription based on subscription id**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


**DELETE /ngsi-ld/v1/subscriptions/#sid**

==============   ============================
Param		 Description
==============   ============================
sid              subscription Id
==============   ============================


**Example:**

.. code-block:: console

   curl -iX DELETE http://localhost:80/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:71
