function makeid(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}

var myID = makeid(5);


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
	
    if (contextEntity.id == null) {
        return;
    }
    
    if (contextEntity.type == null) {
        return;
    }  

    if (contextEntity.Updater != null && contextEntity.Updater == myID) {
        // this is the update from myself
        return
    }
      
   // query the data source to know the available parking lots
    if (contextEntity.datasource != null) {
        var dsURL = contextEntity.datasource.value;
        console.log("connected to the orion broker to fetch the parking information");
        console.log(dsURL);
        //var connection = new NGSI.Connection("http://orion.example.com:1026");

        var updateEntity = {};
        /*updateEntity.entityId = {
            id: contextEntity.entityId.id,
            type: contextEntity.entityId.type,
            isPattern: false
        };*/
	updateEntity.id = contextEntity.id;
	updateEntity.type = contextEntity.type;
        #updateEntity.attributes = {};
        updateEntity.FreeParkingSpots = {type: 'Property', value: 10};
        updateEntity.Updater = {type: 'Property', value: myID};

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

