var T1_PM10=0;
var T2_PM10=25;
var T3_PM10=50;
var T4_PM10=90;
var T5_PM10=180;

var T1_PM25=0;
var T2_PM25=15;
var T3_PM25=30;
var T4_PM25=55;
var T5_PM25=100;

//
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = function(contextEntity, publish, query, subscribe)
{
	console.log("enter into the user-defined fog function");	
	if (contextEntity == null) {
		return;
	} 	
	
	// ============================== processing ======================================================
	// processing the received ContextEntity:  

	console.log('ContextEntity.......',contextEntity);
	
	var id = contextEntity.id;
	console.log('id.......',id);	
	var type = contextEntity.type;
	console.log('type.......',type);
	var PM10 = contextEntity.PM10.value;
	console.log('PM10.......',PM10);
	var PM25 = contextEntity.PM25.value;
	console.log('PM25.......',PM25);
	 var healthRiskPM10=0, healthRiskPM25=0, tempVar=0, healthRisk=null;
    //PM10
        if(PM10 >= T1_PM10 && PM10 < T2_PM10)
            healthRiskPM10=1;
        else if(PM10 >= T2_PM10 && PM10 < T3_PM10)
            healthRiskPM10=2;
        else if(PM10 >= T3_PM10 && PM10 < T4_PM10)
            healthRiskPM10=3;
        else if(PM10 >= T4_PM10 && PM10 < T5_PM10)
            healthRiskPM10=4;
        else if(PM10 >= T5_PM10)
            healthRiskPM10=5;
        else
            healthRiskPM10=-1;
    //PM25
        if(PM25 >= T1_PM25 && PM25 < T2_PM25)
            healthRiskPM25=1;
        else if(PM25 >= T2_PM25 && PM25 < T3_PM25)
            healthRiskPM25=2;
        else if(PM25 >= T3_PM25 && PM25 < T4_PM25)
            healthRiskPM25=3;
        else if(PM25 >= T4_PM25 && PM25 < T5_PM25)
            healthRiskPM25=4;
        else if(PM25 >= T5_PM25)
            healthRiskPM25=5;
        else
            healthRiskPM25=-1;
    //HealthRisk
        if(healthRiskPM10>=healthRiskPM25)
            tempVar=healthRiskPM10;
        else if(healthRiskPM10<healthRiskPM25)
            tempVar=healthRiskPM25;
        else if (healthRiskPM10==-1 & healthRiskPM25==-1)
            tempVar=-1

    console.log('Health Risk  PM10......' + healthRiskPM10);
    console.log('Health Risk  PM25......' + healthRiskPM25);

    switch(tempVar)
    {
        case 1:
	    healthRisk='Very Low';
            break;
        case 2:
            healthRisk='Low';
	    break;
        case 3:
            healthRisk='Medium';
            break;
        case 4:
            healthRisk='High';
            break;
        case 5:
            healthRisk='Very High';
            break;
        default:
	    healthRisk='Unknown';
    }

	console.log('Health Risk......' + healthRisk);
	
	// ============================== publish ======================================================
        // if you need to publish the generated result, please refer to the following example    
	if(tempVar>=3)
	{
		console.log('publishing started......' );
		var updateEntity = {};
		updateEntity.id = id
		updateEntity.type =  'result',
		updateEntity.healthRisk = {'type':'property', 'value': healthRisk}
		publish(updateEntity)
		console.log("publish: ", updateEntity);
	}


 	// ============================== subscribe ======================================================   
    	// if you want to subscribe addtional infromation from the assigned nearby broker, please refer to the following example
	
	/*
    	var subscribeCtxReq = {};    
   	subscribeCtxReq.entities = [{type: 'Device', isPattern: true}];
	
        subscribeCtxReq.type = 'Subscription'
	LdSubscription.notification.format = "normalized"
	LdSubscription.notification.endpoint.uri = my_ip + ":" + myport+ "/notifyContext"
        subscribe(subscribeCtxReq);     
    */	

    // For more information about subscription please refer fogflow doc for NGSILD
	
};

