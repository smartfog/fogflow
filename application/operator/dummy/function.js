var entityID = "";
exports.handler = function(contextEntity, publish, query, subscribe)
{
    console.log("enter into the user-defined fog function");
    
    entityID = contextEntity.entityId.id
    
    var updateEntity = {};
    updateEntity.entityId = {
           id: "Result" + entityID,
           type: 'Result',
           isPattern: false
    };	    	
    updateEntity.attributes = contextEntity.attributes;	      
   	
    publish(updateEntity);	
};

