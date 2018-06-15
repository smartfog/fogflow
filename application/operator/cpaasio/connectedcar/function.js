var speakerID = null;

/**
     * \fn getMetros().
     *
     * \Description: Returns the distance in meters between two coordinate points
     *
     * \param (float) lat1 : Latitude of point  1
     * \param (float) long1 : Longitude of point 1
     * \param (float) lat2 : Latitude of point  2
     * \param (float) long2 : Longitude of point 2
     *
     * \return (float) Distance in meters
     *
     **/

    getMeters = function(lat1,lon1,lat2,lon2)
    {
        rad = function(x) {return x*Math.PI/180;}
        var R = 6378.137*1000; // Radius of the earth in meters
        var dLat = rad( lat2 - lat1 );
        var dLong = rad( lon2 - lon1 );
        var a = Math.sin(dLat/2) * Math.sin(dLat/2) + Math.cos(rad(lat1)) * Math.cos(rad(lat2)) * Math.sin(dLong/2) * Math.sin(dLong/2);
        var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
        var d = R * c;
        return d.toFixed(6); //Retorna metros con tres decimales
     }  

     /**
     * \fn getMinS().
     *
     * \Description: Returns the minutes it takes to reach a coordinate point by car. 
     *  In this case the function can solve it in two ways, the default is a query to 
     *  an OSM API, but it may fail due to loading problems. If this method fails, it 
     *  is performed with the getMetros function and an average speed. This method is 
     *  asynchronous so the promise with await and async functionality has been used.
     *
     * \param (float) lat1 : Latitude of point  1
     * \param (float) long1 : Longitude of point 1
     * \param (float) lat2 : Latitude of point  2
     * \param (float) long2 : Longitude of point 2
     *
     * \return (float) Minutes to destination
     *
     **/
    getMinS = async function(lat1,lon1,lat2,lon2){
      var http = require('http');
      var meanmmin = 330;
      return new Promise(resolve => {
        
        http.get({
          host: 'router.project-osrm.org',
          path: '/route/v1/driving/'+ lat1 +',' + lon1 + ';' + lat2+',' + lon2 +'?overview=false'
        }, function(response) {
          // Continuously update stream with data
          var body = '';
          response.on('data', function(d) {
            body += d;
        });
        response.on('end', function() {
          // Data received, let us parse it using JSON!
          var parsed = JSON.parse(body);
          console.log('Datos: ',parsed);
          var min;
          if(!parsed.routes){
            var m = getMeters(lat1,lon1,lat2,lon2);
            min = m / meanmmin
          }else{
            min = parsed.routes[0].duration/60;
          }
          
          console.log('Min: ',min);
          resolve(min);
          //return await min;
        });
        });
      });
    }


//
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = async function(contextEntity, publish, query, subscribe)
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
    
    // Coordinates are collected and minutes of arrival are determined
    var latitudeDestiny = 37.984737;
    var longitudeDestiny = -1.127266;
    var latitudeCar = contextEntity.metadata.location.value.latitude;
    var longitudeCar = contextEntity.metadata.location.value.longitude;

    // Calculate time to arrived
    var min = await getMinS(latitudeCar,longitudeCar,latitudeDestiny,longitudeDestiny);
    console.log('Time to arrive: ',min);

    //var twentyMinutesLater = new Date();
    //twentyMinutesLater.setMinutes(twentyMinutesLater.getMinutes() + 20);
    
    /* The minutes of arrival are calculated, but the mixture is set on fire. This is because the destination 
    *  should come, in my view, in the entity, I have already tried to change the smartparking.js but it does not 
    *  work for me. I guess it has to be recompiled but I don't know if it's right.
    */
    var parkingreq = {
        arrival_time: min,
        destination: {
            latitude: latitudeDestiny,
            longitude: longitudeDestiny
        }
    };    
    updateEntity.attributes.ParkingRequest = {type: 'object', value: parkingreq};                
       	
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

