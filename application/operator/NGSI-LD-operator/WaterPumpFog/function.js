var current_observation = ""
var previous_observation = ""
var pon_observation = ""
var onHasCreated = "false"
var publishStatus = "false"
//var startTime = new Date().getTime();
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker

function sleep(milliseconds) {
  const date = Date.now();
  let currentDate = null;
  do {
    currentDate = Date.now();
  } while (currentDate - date < milliseconds);
}


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

	var con_observation = ""

	if("on_status" in contextEntity) {
		var con_observation = contextEntity.on_status.observedAt
	}
	current_observation = con_observation
	if (current_observation != previous_observation) { 
		previous_observation = current_observation
		if (con_observation != "" && onHasCreated == "false") {
			console.log("=====initial timer has been created======")
			value = "off"
			onHasCreated = "true"
			createDate = con_observation.split("T")
			publishStatus = "true"
		}

		if (publishStatus == "true") {
			for(var i = 0 ; i<60 ;i++) {
				console.log("Wait for publish")
				sleep(1000);	
			}
			updateEntity["command"] = {"type":"Property", "value": value}
			publish(updateEntity)
			publishStatus = "false"
			onHasCreated = "false"
		}
	} else {
		console.log("Status is already off")
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

