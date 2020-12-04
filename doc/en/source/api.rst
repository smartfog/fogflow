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




FogFlow Service Orchestrator API
=========================================


The overall development process of an IoT Service in FogFlow is shown in the following figure. 
For the development of a fog function, the steps 4 and 5 are combined, which means a default requirement 
is issued by the FogFlow editor when a fog function is submmited. 


.. figure:: figures/development_process.png
   :width: 100 %



Implement an operator
-----------------------------------------------

Before defining the designed service topology, 
all operators used in your service topology must be provided by you or the other provider in the FogFlow system. 


* `nodejs-based`_ 

* `python-based`_ 


.. _`nodejs-based`: https://github.com/smartfog/fogflow/tree/master/application/template/javascript
.. _`python-based`: https://github.com/smartfog/fogflow/tree/master/application/template/python


.. note:: currently two templates are provided: one for nodejs based Implement and the other for python-based implementation



Publish the operator
-----------------------------------------------

The image of operator can be published to the public docker registry or on private docker registery. 
If you do not want to use any docker registry, you have to make sure that 
the docker image of an operator is built on all edge nodes. 
Currently, when the FogFlow worker receives a command to launch a task instance, 
it will first search the required docker image from the local storage. If it does not find it, 
it will start to fetch the required docker image for the docker registry (the public one or any private one, which is up to the 
configuration of the FogFlow worker). 

If anyone would like to publish the image, then following docker command can be used. 


.. code-block:: console   
	
	docker push  [the name of your image]


.. note:: this step is done with only docker commands


Define and register operator
-----------------------------------------------

An operator docker image can also be registered by sending a constructed NGSI update message to the IoT Broker deployed in the cloud. 

Here is a Javascript-based code example to register an operator docker image. 
Within this code example, the Javascript-based library is being used to interact with FogFlow IoT Broker. 
The library can be found from the github code repository (designer/public/lib/ngsi), ngsiclient.js shall be included into the web page. 


.. code-block:: javascript

    var image = {
        name: "counter",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "counter",
        prefetched: false
    };

    //register a new docker image
    var newImageObject = {};

    newImageObject.entityId = {
        id : image.name + ':' + image.tag, 
        type: 'DockerImage',
        isPattern: false
    };

    newImageObject.attributes = {};   
    newImageObject.attributes.image = {type: 'string', value: image.name};        
    newImageObject.attributes.tag = {type: 'string', value: image.tag};    
    newImageObject.attributes.hwType = {type: 'string', value: image.hwType};      
    newImageObject.attributes.osType = {type: 'string', value: image.osType};          
    newImageObject.attributes.operator = {type: 'string', value: image.operatorName};      
    newImageObject.attributes.prefetched = {type: 'boolean', value: image.prefetched};                      
    
    newImageObject.metadata = {};    
    newImageObject.metadata.operator = {
        type: 'string',
        value: image.operatorName
    };               
    
    // assume the config.brokerURL is the IP of cloud IoT Broker
    var client = new NGSI10Client(config.brokerURL);    
    client.updateContext(newImageObject).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });        



Define and register your service topology
-----------------------------------------------

Usually,service topology can be defined and registered via the FogFlow topology editor. 
However, it can also be defined and registered with own code. 

To register a service topology, A constructed NGSI update message is needed to be send by the code to the IoT Broker deployed in the cloud. 

Here is a Javascript-based code example to register an operator docker image. 
Within this code example, the Javascript-based library is used to interact with FogFlow IoT Broker. 
The library can be found from the github code repository (designer/public/lib/ngsi). An ngsiclient.js must be included into the web page. 

.. code-block:: javascript

    // the json object that represent the structure of your service topology
    // when using the FogFlow topology editor, this is generated by the editor
    var topology = {  
       "description":"detect anomaly events from time series data points",
       "name":"anomaly-detection",
       "priority": {
            "exclusive": false,
            "level": 100
       },
       "trigger": "on-demand",   
       "tasks":[  
          {  
             "name":"AnomalyDetector",
             "operator":"anomaly",
             "groupBy":"shop",
             "input_streams":[  
                {  
                      "type": "PowerPanel",
                    "shuffling": "unicast",
                      "scoped": true
                },
                {  
                      "type": "Rule",
                    "shuffling": "broadcast",
                      "scoped": false               
                }                       
             ],
             "output_streams":[  
                {  
                   "type":"Anomaly"
                }
             ]
          },
          {  
             "name":"Counter",
             "operator":"counter",
             "groupBy":"*",
             "input_streams":[  
                {  
                   "type":"Anomaly",
                   "shuffling": "unicast",
                   "scoped": true               
                }           
             ],
             "output_streams":[  
                {  
                   "type":"Stat"
                }
             ]
          }          
       ]
    }
    
    //submit it to FogFlow via NGSI Update
    var topologyCtxObj = {};
    
    topologyCtxObj.entityId = {
        id : 'Topology.' + topology.name, 
        type: topology.name,
        isPattern: false
    };
    
    topologyCtxObj.attributes = {};   
    topologyCtxObj.attributes.status = {type: 'string', value: 'enabled'};
    topologyCtxObj.attributes.template = {type: 'object', value: topology};    
    
    // assume the config.brokerURL is the IP of cloud IoT Broker
    var client = new NGSI10Client(config.brokerURL);    

    // send NGSI10 update    
    client.updateContext(topologyCtxObj).then( function(data) {
        console.log(data);                
    }).catch( function(error) {
        console.log('failed to submit the topology');
    });    



Create a requirement entity to trigger the service topology
--------------------------------------------------------------


Here is the Javascript-based code example to trigger a service topology by sending a customized requirement entity to FogFlow. 


.. code-block:: javascript

    var rid = 'Requirement.' + uuid();    
   
    var requirementCtxObj = {};    
    requirementCtxObj.entityId = {
        id : rid, 
        type: 'Requirement',
        isPattern: false
    };
    
    var restriction = { scopes:[{scopeType: geoscope.type, scopeValue: geoscope.value}]};
                
    requirementCtxObj.attributes = {};   
    requirementCtxObj.attributes.output = {type: 'string', value: 'Stat'};
    requirementCtxObj.attributes.scheduler = {type: 'string', value: 'closest_first'};    
    requirementCtxObj.attributes.restriction = {type: 'object', value: restriction};    
                        
    requirementCtxObj.metadata = {};               
    requirementCtxObj.metadata.topology = {type: 'string', value: curTopology.entityId.id};
    
    console.log(requirementCtxObj);
            
    // assume the config.brokerURL is the IP of cloud IoT Broker
    var client = new NGSI10Client(config.brokerURL);                
    client.updateContext(requirementCtxObj).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to send a requirement');
    });    




Remove a requirement entity to terminate the service topology
---------------------------------------------------------------


Here is the Javascript-based code example to terminate a service topology by deleting the requirement entity. 


.. code-block:: javascript

    var rid = [the id of your created requirement entity];    
            
    // 
    var client = new NGSI10Client(config.brokerURL);                
    client.deleteContext(rid).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to send a requirement');
    });    


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



**c.  To create a new NGSI-LD context entity, with context in Link header and request payload is already expanded**
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


**d. To append additional attributes to an existing entity**
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
	      "id": "urn:ngsi-ld:Vehicle:A4580",
              "type": "Vehicle",

	     ""brandName1"": {
		                 "type": "Property",
		                 "value": "BMW"
	      }
        }'


**e. To update specific attributes of an existing entity**
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
		"id": "urn:ngsi-ld:Vehicle:A4580",
	        "type": "Vehicle",

	       "brandName": {
		                  "type": "Property",
		                  "value": "AUDI"
	        }
	}'
  
**g. To delete an NGSI-LD context entity**
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


**h. To delete an attribute of an NGSI-LD context entity**
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


**i. To retrieve a specific entity**
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


**j. To retrieve entities by attributes**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/entities?attrs=(Value 1)**

==============   ============================
Param		 Description
==============   ============================
Value 1          Attriute Value
==============   ============================

**Example:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/entities?attrs=http://example.org/vehicle/brandName -H 'Content-Type: application/ld+json' -H 'Accept: application/ld+json'



**k. To retrieve a specific entity by ID and Type**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/entities?id=(value 1)&type=(value 2)**

==============   ============================
Param		 Description
==============   ============================
value 1          Attribute Value of Entity
Value 2          Type Value of Entity
==============   ============================

**Example:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/entities?id=urn:ngsi-ld:Vehicle:A4569&type=http://example.org/vehicle/Vehicle


**l. To retrieve a specific entity by Type**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi-ld/v1/entities?type=(Value 1)**

==============   ============================
Param		 Description
==============   ============================
Value 1          Type Value
==============   ============================

**Example:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/entities?type=http://example.org/vehicle/Vehicle


**m. To retrieve a specific entity by Type, context in Link Header**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET ngsi-ld/v1/entities?type=(Value 1)**

==============   ============================
Param		 Description
==============   ============================
Value 1          Type Value
==============   ============================

**Header Format**

=============     ===========================================================
key               Value
=============     ===========================================================
Content-Type      application/json
Accept            application/ld+json
Link              <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; 
                  type="application/ld+json"
=============     ===========================================================

**Example:**

.. code-block:: console

   curl -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'  http://localhost:80//ngsi-ld/v1/entities?type=Vehicle -H 'Content-Type: application/ld+json' -H 'Accept: application/ld+json'


**n. To retrieve a specific entity by IdPattern and Type**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET : /ngsi-ld/v1/entities?idPattern=(Value 1)&type=(Value 2)**

==============   ============================
Param		 Description
==============   ============================
value 1          idPattern Value of Entity
Value 2          Type Value of Entity
==============   ============================ample:**

.. code-block:: console

   curl http://localhost:80/ngsi-ld/v1/entities?idPattern=urn:ngsi-ld:Vehicle:A.*&type=http://example.org/vehicle/Vehicle

       
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
		"entities": [{
				"idPattern": ".*",
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

   curl http://localhost:80/ngsi-ld/subscriptions/ -H 'Accept: application/ld+json'


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

   curl http://localhost:80/ngsi-ld/subscriptions/urn:ngsi-ld:Subscription:71


**d. To update a specific subscription based on subscription id, with context in Link header**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**PATCH  /ngsi-ld/v1/subscriptions/#sid**

==============   ============================
Param		 Description
==============   ============================
sid              subscription Id
==============   ============================


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
     'http://localhost:80/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:71' \
      -H 'Content-Type: application/ld+json' \
      -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
       {
	 	"id": "urn:ngsi-ld:Subscription:7",
	 	"type": "Subscription",
	 	"entities": [{
	  			"type": "Vehicle"
	 	  }],
	 	"watchedAttributes": ["http://example.org/vehicle/brandName2"],
	        "q":"http://example.org/vehicle/brandName2!=Mercedes",
	 	"notification": {
	  	"attributes": ["http://example.org/vehicle/brandName2"],
	  	"format": "keyValues",
	  	"endpoint": {
	   			"uri": "http://my.endpoint.org/notify",
				"accept": "application/json"
	  	 }
	      }
	  }'


**e. To delete a specific subscription based on subscription id**
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

