var pof_observation = ""
var pon_observation = ""

//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = function(contextEntity, publish, query, subscribe)
{
	console.log("enter into the user-defined fog function");	
	if (contextEntity == null) {
		return;
	} 	
	
	// ============================== processing ======================================================
	// processing the received ContextEntity:  
	

	console.log('ContextEntity.......',contextEntity);
	var updateEntity = {};
	for (var key in contextEntity) {
		updateEntity[key] = contextEntity[key]
	}
	var on_status = contextEntity["on_status"]
	var off_status = contextEntity["off_status"]
	var con_observation = ""
	var cof_observation = ""
	if ("observedAt" in on_status) {
		con_observation = on_status["observedAt"]
	}
	
        if ("observedAt" in of_status) {
                cof_observation = of_status_["observedAt"]
        }
	var value = ""
	if (con_observation == "" && cof_observation != "") {
		value = "off"
	} else if (con_observation != "" && cof_observation == "") {
		value = "on"
	} else if (con_observation != "" && cof_observation != "") {
		if (con_observation != pon_observation) {
			value = "on"
			pon_observation = con_observation
		}
	} else if (con_observation != "" && cof_observation != "") { 
		if (cof_observation != pof_observation) {
                        value = "of"
                        pon_observation = con_observation
                }
	}

	var id = contextEntity.id;
	console.log('id.......',id);	
	var type = contextEntity.type;
	console.log('type.......',type);
	// ============================== publish ======================================================
        // if you need to publish the generated result, please refer to the following example   
	/*
		This example will publish the result is speed is under risk
	*/ 
	if (value != "") {
		console.log('publishing started......' );
		updateEntity.id = id
		updateEntity.type =  type
		updateEntity["command"] = {'type':'Property', 'value': value}
		console.log(updateEntity)
		publish(updateEntity)
		console.log("publish: ", updateEntity);
	}
	
 	// ============================== subscribe ======================================================   
    	// if you want to subscribe addtional infromation from the assigned nearby broker, please refer to the following example
	
	/*
    	var subscribeCtxReq = {};    
   	subscribeCtxReq.entities = [{type: 'Device', isPattern: true}];
	
        subscribeCtxReq.type = 'Subscription'
	LdSubscription.notification.format = "normalized"
	LdSubscription.notification.endpoint.uri = my_ip + ":" + myport+ "/notifyContext"
        subscribe(subscribeCtxReq);     
    */	

    // For more information about subscription please refer fogflow doc for NGSILD
	
};

