.. _trigger-topology:

How to trigger a service topology
========================================

Once a service topology has been submitted to FogFlow, we can use a NGSI application to 
trigger the full or partial data processing of this service topology, 
by sending a specific NGSI entity update, which is called "requirement" in FogFlow. 

Here is the Javascript-based code example to send a customized requirement object via NGSI10 update. 
This requirement entity is identified by its entity ID. 
Later on, if you would like to issue any change to the same requirement, 
you can send a NGSI10 update to this requirement entity based on its entity ID. 

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


