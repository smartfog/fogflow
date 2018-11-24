exports.handler = function(contextEntity, publish){
	if (contextEntity == null) {
		return;
	} 
	
	console.log("------------begin-------------");
    console.log("received detected event", contextEntity)
	console.log("------------end-------------");		
	
	if (contextEntity.attributes == null) {
		return;
	}
	    
	if (contextEntity.metadata.homeID == null) {
		return;
	}    
	
    var homeID = contextEntity.metadata.homeID.value;
		
	var updateEntity = {};
	updateEntity.entityId = {
        id: "Home." + homeID,
        type: 'Home',
        isPattern: false
    };	    	
	updateEntity.attributes = {};	

    if (contextEntity.attributes.detectedEvent != null) {        
	    var detectedEvent = contextEntity.attributes.detectedEvent.value;        
        
        var alert = {
            type: "ASK_FOR_HELP",
            reason: detectedEvent,
            device: contextEntity.entityId.id        
        }    
        
	    updateEntity.attributes.alert = {		
            type: 'object',
            value: alert
        };        
    }

    updateEntity.metadata = {};	

	if (contextEntity.metadata.location != null) {
        updateEntity.metadata.location = contextEntity.metadata.location;
	}        
        
	updateEntity.attributes.pushbutton = {type: 'string', value: contextEntity.entityId.id};

	console.log("publish: ", updateEntity);		    	
	publish(updateEntity, -1);
};

