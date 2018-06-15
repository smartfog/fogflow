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
     *  Web OSM API: http://project-osrm.org/
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

    /**
     * \fn enoughPlaces().
     *
     * \Function that checks if it is possible to park in this parking. At the moment 
     *  it is a simple function that checks if there are places, taking into consideration 
     *  that it takes one per minute. This function is a very simple version and should 
     *  be improved, at the moment there are no resources available as histories to see how 
     *  long it takes to fill a place.
     *
     * \param (float) freePlaces : Free places of parking
     * \param (float) minToPark : Minutes of arrival at the parking
     * \param (float) minToArrived : Minutes of arrival at the destination
     *
     * \return (boolean) If parking is available
     *
     **/
    enoughPlaces = function(freePlaces,minToPark,minToArrived)
    {
        return freePlaces > (minToPark+minToArrived);
    } 


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
	if (contextEntity.attributes.RecommendedParkingSite != null) {
		return;
	}

    // to inform the driver where to park for the given parking request
    if (contextEntity.attributes.ParkingRequest != null) {
        var parkingReq = contextEntity.attributes.ParkingRequest.value;
        
        console.log("~~~~~~~~~~~~~~~~~~~~~~receive the following parking request~~~~~~~~~~~~~~~~~");
        console.log(parkingReq);

        var parkingSaved = null;
        var minSaved;

        var queryReq = {}
        queryReq.entities = [{type:'SmartParking', isPattern: true}];
        var entities;

        var handleQueryResult = async function(entityList) {
            for(var i=0; i<entityList.length; i++) {
                var entity = entityList[i];
                console.log('===============' + i + '===================');
                console.log(entity);

                // Coordinates are collected and minutes of arrival are determined
                var latitudeEntity = entity.attributes.LocationParking.value.latitude;
                var longitudeEntity = entity.attributes.LocationParking.value.longitude;
                var latitudeCar = contextEntity.attributes.ParkingRequest.value.destination.latitude;
                var longitudeCar = contextEntity.attributes.ParkingRequest.value.destination.longitude;

                console.log('Latitude entity: ',latitudeEntity);
                console.log('Longitude entity: ',longitudeEntity);
                console.log('Latitude car: ',latitudeCar);
                console.log('Longitude car: ',longitudeCar);

                var minToPark = await getMinS(latitudeCar,longitudeCar,latitudeEntity,longitudeEntity);
                console.log('This parking min: ',minToPark);

                var places = enoughPlaces(entity.attributes.FreeParkingSpots.value,minToPark,contextEntity.attributes.ParkingRequest.value.arrival_time);
                console.log('This parking places: ',entity.attributes.FreeParkingSpots.value);
                console.log('Result of enoughPlaces: ',places);

                //Compared to the previous one
                if(places && (parkingSaved == null || minToPark < minSaved)){
                    parkingSaved = entity;
                    minSaved = minToPark;
                    console.log("New best parking: ",parkingSaved);
                    console.log('Min: ',minSaved);
                }
                
            }

            console.log("Parking designation: ",parkingSaved);

            // send an update to tell which parking site the driver should go        
            var updateEntity = {};
            updateEntity.entityId = {
                id: contextEntity.entityId.id,
                type: contextEntity.entityId.type,
                isPattern: false
            };

            updateEntity.attributes = {};

            // If there is no parking space, the following is assigned 'No parking available'
            if(parkingSaved == null){
                updateEntity.attributes.RecommendedParkingSite = {type: 'string', value: 'No parking available'}; 
            }else{
                updateEntity.attributes.RecommendedParkingSite = {type: 'string', value: parkingSaved.attributes.NameParking.value}; 
            } 
            
                
            publish(updateEntity);
            console.log("publish: ", updateEntity); 
        }

        query(queryReq, handleQueryResult);        
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

