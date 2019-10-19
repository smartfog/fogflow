var cameraID = null;

exports.handler = function(contextEntity, publish, subscribe){
    
    if (contextEntity == null) {
        return;
    } 
    
    if (contextEntity.attributes == null) {
        return;
    }
    
    console.log("========RECEIVE============");
    console.log(contextEntity);
    
    return 
    
//    if (contextEntity.entityId.type == 'Camera') {               
//        var ctxObj = {};
    
//        ctxObj.entityId = {};
    
//        ctxObj.attributes = {};         
                
//        ctxObj.metadata = {};       
//        ctxObj.metadata.awningID = {
//            type: 'string',
//            value: awningID
//        };                     
        
//        if (contextEntity.attributes.raining.value == true) {
//            ctxObj.attributes.command = {       
//                type: 'string',
//                value: 'open'
//            };              
//        } else {
//            ctxObj.attributes.command = {       
//                type: 'string',
//                value: 'close'
//            };                          
//        }
        
//        publish(ctxObj, 0);
//    }
    
};

