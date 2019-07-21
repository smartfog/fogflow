var parkingRequest = null;

//
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = function(contextEntity, publish, query, subscribe)
{
    if (contextEntity == null) {
        return;
    } 	
    if (contextEntity.attributes == null) {
        return;
    }    
    
    if (contextEntity.attributes.ParkingRequest != null) {
        return;
    }  
    
    if (contextEntity.attributes.RecommendedParkingSite != null) {
        return;
    }  
        
    
    // to calculate how long it will take to arrive the planned destination accordingly
    // todo: @Javier, please add your implementation here
    
    
    // create and update the ParkingRequest attribute of this connected car
    var updateEntity = {};
    updateEntity.entityId = {
           id: contextEntity.entityId.id,
           type: contextEntity.entityId.type,
           isPattern: false
    };	    	
    updateEntity.attributes = {};	 
        
    var twentyMinutesLater = new Date();
    twentyMinutesLater.setMinutes(twentyMinutesLater.getMinutes() + 20);

    parkingRequest = {
        arrival_time: twentyMinutesLater,
        destination: {
            latitude: 37.984737,
            longitude: -1.127266
        }
    };    
    updateEntity.attributes.ParkingRequest = {type: 'object', value: parkingRequest};                
       	
    publish(updateEntity);
    console.log("publish: ", updateEntity);	    
    
    

    // ============================== publish ======================================================
    
    // if you need to publish the generated result, please refer to the following example    
    
    /*
    var updateEntity = {};
    updateEntity.entityId = {
           id: "Twin.Home.0001",
           type: 'Home',
           isPattern: false
    };	    	
    updateEntity.attributes = {};	 
    updateEntity.attributes.city = {type: 'string', value: 'Heidelberg'};                
    
    updateEntity.metadata = {};    
    updateEntity.metadata.location = {
        type: 'point',
        value: {'latitude': 33.0, 'longitude': -1.0}
    };        
   	
    publish(updateEntity);
    console.log("publish: ", updateEntity);		                
    */

    
    // ============================== query ======================================================
    
    // if you want to query addtional information from the assigned nearby broker, please refer to the following example
    
    /*
    var queryReq = {}
    queryReq.entities = [{type:'PublicSite', isPattern: true}];    
    var handleQueryResult = function(entityList) {
        for(var i=0; i<entityList.length; i++) {
            var entity = entityList[i];
            console.log('===============' + i + '===================');
            console.log(entity);   
        }
    }  
    
    query(queryReq, handleQueryResult);
    */
    
    
    // ============================== subscribe ======================================================
    
    // if you want to subscribe addtional infromation from the assigned nearby broker, please refer to the following example

    /*
    var subscribeCtxReq = {};    
    subscribeCtxReq.entities = [{type: 'Home', isPattern: true}];
    subscribeCtxReq.attributes = ['alert'];        
    //subscribeCtxReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'geohash='+contextEntity.attributes.geohash.value}]}        	
    
    subscribe(subscribeCtxReq);     
    */
	
};

