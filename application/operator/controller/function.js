var awningID = null;

exports.handler = function(contextEntity, publish, subscribe){
    
    if (contextEntity == null) {
        return;
    } 
    
    if (contextEntity.attributes == null) {
        return;
    }
    
    if (contextEntity.entityId.type == 'SmartAwning') {
        awningID = contextEntity.entityId.id;
        
        // trigger a subsciption to fetch the prediction result
        var subscribeCtxReq = {};    
        subscribeCtxReq.entities = [{type: 'Prediction', isPattern: true}];
        //subscribeCtxReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'geohash='+contextEntity.attributes.geohash.value}]}           
        
        subscribe(subscribeCtxReq);     
           
    } else if (contextEntity.entityId.type == 'Prediction') {
        // issue control commands to awning devices according to the received prediction results
        // var geohash = contextEntity.attributes.geohash.value;
                
        var ctxObj = {};
    
        ctxObj.entityId = {};
    
        ctxObj.attributes = {};         
                
        ctxObj.metadata = {};       
        ctxObj.metadata.awningID = {
            type: 'string',
            value: awningID
        };                     
        
        if (contextEntity.attributes.raining.value == true) {
            ctxObj.attributes.command = {       
                type: 'string',
                value: 'open'
            };              
        } else {
            ctxObj.attributes.command = {       
                type: 'string',
                value: 'close'
            };                          
        }
        
        publish(ctxObj, 0);
    }
    
};

