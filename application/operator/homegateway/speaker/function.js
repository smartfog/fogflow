var speakerID = null;
var myHomeID = null;

exports.handler = function(contextEntity, publish, subscribe){	
	if (contextEntity == null) {
		return;
	} 
	
	console.log("------------begin-------------");
    console.log("received detected event", contextEntity)
	console.log("------------end-------------");		
	
	if (contextEntity.attributes == null) {
		return;
	}
    
    if (contextEntity.entityId.type == 'Speaker') {
        speakerID = contextEntity.entityId.id;
        myHomeID = "Home." + contextEntity.metadata.homeID.value;
		        
        //  update the home entity with this speaker information
	    var updateEntity = {};
	    updateEntity.entityId = {
            id: myHomeID,
            type: 'Home',
            isPattern: false
        };	    	
	    updateEntity.attributes = {};	 
	    updateEntity.attributes.speaker = {type: 'string', value: speakerID};      
		
    		updateEntity.metadata = {};	
		if (contextEntity.metadata.location != null) {
        		updateEntity.metadata.location = contextEntity.metadata.location;
		}    
    	
	    console.log("publish: ", updateEntity);		                
	    publish(updateEntity, -1);        
        
        // trigger a subsciption to fetch the prediction result
        var subscribeCtxReq = {};    
        subscribeCtxReq.entities = [{type: 'Home', isPattern: true}];
        subscribeCtxReq.attributes = ['alert'];        
        //subscribeCtxReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'geohash='+contextEntity.attributes.geohash.value}]}        	
        
        subscribe(subscribeCtxReq);     

        // trigger a subsciption to fetch the prediction result
        var subscribeCtxReq = {};    
        subscribeCtxReq.entities = [{type: 'Announcement', isPattern: true}];
        subscribeCtxReq.attributes = ['annoucement'];                
        subscribe(subscribeCtxReq);   
           
    } else if (contextEntity.entityId.type == 'Home') {
        // issue control commands to the speaker device according to the detected home events        
        
        var attachedSpeakerID = null;
        
        if (contextEntity.attributes.speaker == null) {
            attachedSpeakerID =  speakerID;
        } else {
            attachedSpeakerID = contextEntity.attributes.speaker.value;            
        }
               
        if (attachedSpeakerID == null) {
            return;
        }
		     
        var alert =  contextEntity.attributes.alert.value;
		var fromWhichHome = contextEntity.entityId.id;
		if (fromWhichHome == myHomeID) {
			console.log("the events are from the current home");
			return;	
		} 
                
        var ctxObj = {};
	
	    ctxObj.entityId = {};
        ctxObj.entityId.id = attachedSpeakerID;
        ctxObj.entityId.isPattern = false;        
	
        ctxObj.attributes = {};	        
        ctxObj.attributes.command = {		
            type: 'object',
            value: alert
        };              
        
	    console.log("publish: ", ctxObj);		                        
        publish(ctxObj, -1);        
    } else if (contextEntity.entityId.type == 'Announcement') {
        // issue control commands to the speaker device according to the received announcement       
        
        var attachedSpeakerID = null;
        
        if (contextEntity.attributes.speaker == null) {
            attachedSpeakerID =  speakerID;
        } else {
            attachedSpeakerID = contextEntity.attributes.speaker.value;            
        }
        
        console.log("====speaker ID=========", attachedSpeakerID);
               
        if (attachedSpeakerID == null) {
            return;
        }
		
		var annoucement = contextEntity.attributes.annoucement.value;		
		
		var alert = {};
		alert.type = "BROADCAST";
		alert.value = annoucement;
		                    
        var ctxObj = {};
	
	    ctxObj.entityId = {};
        ctxObj.entityId.id = attachedSpeakerID;
        ctxObj.entityId.isPattern = false;        
	
        ctxObj.attributes = {};	        
        ctxObj.attributes.command = {		
            type: 'object',
            value: alert
        };              
        
	    console.log("publish: ", ctxObj);		                        
        publish(ctxObj, -1);        
    }
	
};

