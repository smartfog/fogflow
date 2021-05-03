var pof_observation = ""
var pon_observation = ""
var activateTimer = "off"
//var startTime = new Date().getTime();
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
	console.log(updateEntity)
	var con_observation = ""
	var cof_observation = ""
	if("on_status" in contextEntity) {
		con_observation = contextEntity.on_status.observedAt
	}
	if("off_status" in contextEntity) {
		cof_observation = contextEntity.off_status.observedAt
	}
	if(cof_observation == undefined) {
		console.log("not a off status")
	}
	console.log("onstatus")
	console.log(con_observation)
	console.log(cof_observation )
	if (con_observation != "" && cof_observation == "") {
		value = "off"
		activateTimer = "on"
	} else if (con_observation != "" && cof_observation != "") {
		if (con_observation != pon_observation) {
			value = "off"
			activateTimer = "on"
			pon_observation = con_observation
		}
	}
	var id = contextEntity.id;
	console.log('id.......',id);	
	var type = contextEntity.type;
	console.log('type.....',type)
	var diff = 0
	if (activateTimer == "on") {
		var startTime = new Date().getTime();
		activateTimer == "off"
		while (diff < 1) {
			var endTime = new Date().getTime();
			var diff = Math.abs(endTime - startTime) / 60000
			console.log("waitting for update")
		}
		updateEntity["command"] = {'type':'Property', 'value': value}
		publish(updateEntity) //pulish result on broker
		console.log("publish: ", updateEntity);
		
	} else {
		console.log("The water pump is off")
	}
	// ============================== publish ======================================================
        // if you need to publish the generated result, please refer to the following example   
	/*
		This example will publish the result is speed is under risk
	*/ 
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

