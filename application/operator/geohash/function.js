var geohash = require('ngeohash');

exports.handler = function(contextEntity, publish){
    if (contextEntity == null) {
        return;
    } 
    
    if (contextEntity.attributes == null) {
        return;
    }
    
    if (contextEntity.attributes.location == null) {
        return;
    }
    
    if (contextEntity.attributes.geohash != null) {
        console.log("============this entity has been processed before=============");
        return;
    }
    
    console.log("------------begin-------------");
    console.log("received content entity ", contextEntity)
    console.log("------------end-------------");    
    
    var myposition = contextEntity.attributes.location.value;
    
    var hashstring = geohash.encode(myposition.latitude, myposition.longitude, precision=4);
    
    var updateEntity = {};
    updateEntity.entityId = contextEntity.entityId;    
    
    updateEntity.attributes = {};    
    updateEntity.attributes.geohash = {        
        type: 'string',
        value: hashstring
    };
    
    updateEntity.metadata = {};    
    updateEntity.metadata.geohash = {        
        type: 'string',
        value: hashstring
    };    
    
    publish(updateEntity, -1);
    console.log("publish: ", updateEntity);        
};

