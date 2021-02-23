var nolmalspeed = 30;
var averageSpeed1 = 30;
var averageSpeed2 = 50;
var highSpeed = 50;
//
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
	
	var id = contextEntity.id;
	console.log('id.......',id);	
	var type = contextEntity.type;
	console.log('type.......',type);
	var speed = contextEntity.speed.value;
	console.log('speed.......',speed);
	var brand = contextEntity.brandName.value;
	console.log('brand.......',brand);
	var resultSpeed = null
    //PM10
        if(speed <= 30)
            resultSpeed =  "speed is slow"
	else if(speed > averageSpeed1 && speed <= averageSpeed2)
	    resultSpeed =  "Speed is normal"
	else 
	    resultSpeed = "speed is under risk"
        
	console.log('Speed......' + resultSpeed);
	
	// ============================== publish ======================================================
        // if you need to publish the generated result, please refer to the following example   
	/*
		This example will publish the result is speed is under risk
	*/ 
	if(speed > 50)
	{
		console.log('publishing started......' );
		var updateEntity = {};
		updateEntity.id = id+"daresult"
		updateEntity.type =  'daresult',
		updateEntity.speed = {'type':'Property', 'value': resultSpeed}
		//console.log(updateEntity)
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

