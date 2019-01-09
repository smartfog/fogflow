var geohash = require('ngeohash');

exports.handler = function(contextEntity, publish){
    if (contextEntity == null) {
        return;
    } 
    
    if (contextEntity.attributes == null) {
        return;
    }
    
    publish(contextEntity, 0);
};

