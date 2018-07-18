*****************************************
APIs and examples of their usage 
*****************************************

APIs of FogFlow Discovery
===================================

lookup of nearby brokers
-----------------------------------------------

For any external application or IoT devices, the only interface they need from FogFlow Discovery is to find out a nearby 
Broker based on its own location. After that, they only need to interact with the assigned nearby Broker. 
Please check the following examples. 

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8080/ngsi10/updateContext' \
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
        var ngsi10client = new NGSI.NGSI10Client(config.brokerURL);	
    
    	// send NGSI10 update	
        ngsi10client.updateContext(topologyCtxObj).then( function(data) {
            console.log(data);                
        }).catch( function(error) {
            console.log('failed to submit the topology');
        });
   
       

APIs of FogFlow Broker
===============================

Create and update a context entity
-----------------------------------------------


Fetch a context entity by ID
-----------------------------------------------


Check all context entities on a single Broker
-------------------------------------------------








