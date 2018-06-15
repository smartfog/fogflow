var OrionNGSI = require('ngsijs');

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

    // query the data source to know the available parking lots
    if (contextEntity.attributes.datasource != null) {
        var dsURL = contextEntity.attributes.datasource.value;
        console.log("connected to the orion broker to fetch the parking information");
        console.log(dsURL);
        //var connection = new NGSI.Connection("http://orion.example.com:1026");       

        var updateEntity = {};

        // The content of the entity's identification is copied
        // The type is changed so that the entities of the private and public car parks 
        // are equal, so that the recommender does not make distinctions between them.
        console.log("=========ID copy========="); 
        updateEntity.entityId = {
            id: "SmartParking: " + contextEntity.entityId.id,
            type: "SmartParking",
            isPattern: false
        };

        console.log("=========Attributes copy========="); 
        updateEntity.attributes = {};
        //The attribute name is copied
        //In this case the O.R.A. does not have a name attribute, so it puts its id in
        updateEntity.attributes.NameParking = {type: 'string', value: contextEntity.entityId.id};          
        //The free places attribute is copied 
        updateEntity.attributes.FreeParkingSpots = {type: 'integer', value: contextEntity.attributes.libres.value};          

        //The coordinates are copied as metadata
        console.log("Latitude: ", contextEntity.metadata.location.value.latitude); 
        console.log("Longitude: ",contextEntity.metadata.location.value.longitude); 
        var location = {
            latitude: contextEntity.metadata.location.value.latitude,
            longitude: contextEntity.metadata.location.value.longitude
        };    
        updateEntity.attributes.LocationParking = {type: 'object', value: location};                
            
        publish(updateEntity);
        console.log("publish: ", updateEntity);	        
    }    
    
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

