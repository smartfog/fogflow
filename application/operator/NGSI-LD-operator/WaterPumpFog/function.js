var pof_observation = ""
var pon_observation = ""
var onHasCreated = "false"
var initiateUpdate = "false"
var gotUpdate = "false"
//var startTime = new Date().getTime();
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker

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
		var con_observation = contextEntity.on_status.observedAt
	}
	if (con_observation != "" && onHasCreated == "false") {
		value = "off"
		onHasCreated = "true"
		createDate = con_observation.split("T")
		console.log("=====initial timer has been created======")
	} else if (con_observation != "" && onHasCreated == "true") {
		date = con_observation.split("T")
		initiateUpdate = "true"
		value  = "off"
	}
	if (initiateUpdate == "true") {
		//console.log(date)
		//console.log(createDate)
		console.log(createDate[0])
		console.log(date[0])
		if (createDate[0] != date[0]) {
			console.log("date is not equal")
			updateEntity["command"] = {"type":"Property", "value": value}
			onHasCreated = "false"
			initiateUpdate = "false"
			gotUpdate = "true"
			publish(updateEntity)
		 } else {
			createH = createDate[1].split(":")
			updateH = date[1].split(":")
			console.log(parseInt(updateH[0]))
			console.log(parseInt(createH[0]))
			if (parseInt(updateH[0]) - parseInt(createH[0]) >= 1) {
				updateEntity["command"] = {"type":"Property", "value": value}
				publish(updateEntity)
				onHasCreated = "false"
				initiateUpdate = "false"
				gotUpdate = "true"
			}
		}
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

