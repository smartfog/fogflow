'use strict';
const AWS = require('aws-sdk');
const mqtt = require('./mqtt-lib.js');
const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');
var Promise = require('promise');

var ngsi10client;
var brokerURL;
var myReferenceURL = '';
var isConfigured = false;

var mqttclient = null;
var bufferedTopics = [];
var deviceList = {};
var entity2topic = {};
var topic2entity = {};


function startApp() 
{
    console.log('start to receive input data streams via a listening port');
}

function stopApp() 
{
	console.log('clean up the app');
}

// process the input data stream and generate output stream
function processInputStreamData(data) 
{
    if(data.entityId.type == 'DeviceProfile.awsiot') {
        doMediation(data);		
    } else {
        onReceivedNGSIUpdate(data);
    }	    
}

function doMediation(deviceProfileObj)
{
    console.log('receiving a device profile');
    var id = deviceProfileObj.attributes.id.value;    
    var profile = deviceProfileObj.attributes.profile.value;        
    
    if (id in deviceList) {
        return;
    }    

    // extract the searchable attributes from AWS IoT
    getThingObjectFromAWS(profile).then( function(response) {        
        var awsThingObj = response;

        // create a ngsi context entity for this device
        publishNewDevice(profile, awsThingObj);    

        // establish the attribute association and data communcation between entity attributes and topics
        if( mqttclient == null ) {
            establishMQTTClient(profile);
        } else {
            bridgeDevice(profile);    
        }            
    });    
    
    
    deviceList[id] = true;
}

function publishNewDevice(profile, thing)
{    
    var ctxObj = {};
	
	ctxObj.entityId = {};
		
    ctxObj.entityId.id = 'Stream.ngsi.' + profile.type + '.' + profile.id;
	ctxObj.entityId.type = 'ngsi.' + profile.type;
	ctxObj.entityId.isPattern = false;

    // convert all searchable attributes of the thing into the metadata of the device entity
	ctxObj.metadata = {};		    
    for(var key in thing.attributes) {
        var attrValue = thing.attributes[key];
        ctxObj.metadata[key] = {type: 'string', value: attrValue};        
    }
	
    // convert pull or push based attributes in the device profile into the attributes of the device entity
    ctxObj.attributes = {};        
    if(profile.publish) {
        for(var i=0; i<profile.publish.length; i++) {
            var name = profile.publish[i];
            ctxObj.attributes[name] = {type: 'object', value: null};
        }
    }    
    if(profile.subscribe) {
        for(var i=0; i<profile.subscribe.length; i++) {
            var name = profile.subscribe[i];
            ctxObj.attributes[name] = {type: 'object', value: null};
        }        
    }
    	   		      
    ngsi10client.updateContext(ctxObj).then( function(data) {
		console.log('======send update to create device entity======');
		console.log(ctxObj);
        console.log(data);
    }).catch(function(error) {
		console.log(error);
        console.log('failed to update context');
    });       
}

function bridgeDevice(profile)
{
    if(profile.publish) {
        for(var i=0; i<profile.publish.length; i++) {
            var name = profile.publish[i];
            var topic = '/' + profile.type + '/' + profile.id + '/' + name;
            
            subscribeMQTT(topic);
            
            var entityid = 'Stream.ngsi.' + profile.type + '.' + profile.id;
            var entitytype = 'ngsi.' + profile.type;
            topic2entity[topic] = {'id': entityid, 'type': entitytype, 'isPattern': false};
        }        
    }
    
    if(profile.subscribe) {
        for(var i=0; i<profile.subscribe.length; i++) {
            var attribute = profile.subscribe[i];        
            var entityid = 'Stream.ngsi.' + profile.type + '.' + profile.id;      
            
            subscribeNGSI(entityid, attribute);
            
            var topic = '/' + profile.type + '/' + profile.id + '/' + attribute;
            entity2topic[entityid] = {'topic': topic, 'attribute': attribute};
        }
    }
}

function subscribeMQTT(topic)
{
    if( mqttclient.connected ) {
        mqttclient.subscribe(topic);
    } else {
        bufferedTopics.push(topic);
    }                    
}

function subscribeNGSI(entityid, attribute)
{
    var subscribeCtxReq = {};    
    subscribeCtxReq.entities = [{id: entityid, isPattern: false}];
    subscribeCtxReq.reference = myReferenceURL;
    
    ngsi10client.subscribeContext(subscribeCtxReq).then( function(subscriptionId) {
		console.log('======send NGSI subscription======');        
        console.log(subscriptionId);   
    }).catch(function(error) {
        console.log('failed to subscribe context');
    });    
}

function establishMQTTClient(profile)
{
    var myConfig = new AWS.Config({
        accessKeyId: profile.accessKeyId, 
        secretAccessKey: profile.secretAccessKey, 
        region: profile.region
    });    

    var iot = new AWS.Iot(myConfig);
    iot.describeEndpoint({}, function(err, data){
        if (err) {
            console.log(err);
            return;
        }       
            
        var connectOpts = {
            accessKey: profile.accessKeyId,
            clientId: profile.id,    
            endpoint: data.endpointAddress,
            secretKey: profile.secretAccessKey,
            regionName: profile.region
        };                                
                
        mqttclient = new mqtt.MQTTClient(connectOpts);     
                                
        function onConnect() {
            // handle the delayed subscriptions
            handleBufferedTopics();
        }   
        
        function onMessage(message) {
            console.log(message);            
            onReceivedMQTTUpdate(message);
        }  
                     
        function onConnectionLost() {
            console.log('connection lost');            
        }                
                
        mqttclient.on('onConnect', onConnect);
        mqttclient.on('onMessageArrived', onMessage);
        mqttclient.on('onConnectionLost', onConnectionLost);        
                      
        bridgeDevice(profile);
    });
}

function getThingObjectFromAWS(profile)
{
    var myConfig = new AWS.Config({
        accessKeyId: profile.accessKeyId, 
        secretAccessKey: profile.secretAccessKey, 
        region: profile.region
    });

    var iot = new AWS.Iot(myConfig);

    var params = {
        thingTypeName: profile.type
    };

    var targetedThingName = profile.type + profile.id;

    return new Promise( function(resolve, reject)     {
        iot.listThings(params, function(err, data) {
            if (err) {
                reject(err);
            } 
    
            for (var i=0; i<data.things.length; i++) {
                var thing = data.things[i];    
                if(thing.thingName == targetedThingName) {
                    resolve(thing);
                }
            }  
            
            reject(new Error('not found'));
        });                  
    });             
}

function handleBufferedTopics()
{
    for(var i=0; i<bufferedTopics.length; i++){
        var topic = bufferedTopics[i];         
        mqttclient.subscribe(topic);
        console.log('subscribe to topic: ', topic);
    }
    
    bufferedTopics = [];
}

function onReceivedMQTTUpdate(msg)
{    
    var topic = msg.topic;
    var jsondata = JSON.parse(msg.message);
    
    var entityid = topic2entity[topic];
    
    var ctxObj = {};	
	ctxObj.entityId = entityid;    

    var items = topic.split('/');
    var len = items.length - 1;
    var attribute = items[len];
    	
    ctxObj.attributes = {};        
    ctxObj.attributes[attribute] = {type: 'object', value: jsondata};
            	   		      
    ngsi10client.updateContext(ctxObj).then( function(data) {
		console.log('======send update to device entity======');
		console.log(ctxObj);
        console.log(data);
    }).catch(function(error) {
		console.log(error);
        console.log('failed to update context');
    });       
}


function onReceivedNGSIUpdate(entity) 
{
    console.log('---------------------------ngsi------------------')
    console.log(entity);
    
    var eid = entity.entityId.id;
    var topic = entity2topic[eid]
    
    console.log('topic : ', topic);
    
    var attributeObj = entity.attributes[topic.attribute];    
    if( mqttclient.connected ) {        
        mqttclient.publish(topic.topic, attributeObj.value);
    }    
}

// handle the commands received from the engine
function handleAdmin(req, commands, res) 
{	
    console.log('=============commands=============');
    console.log(commands);
	handleCmds(commands);
	
	isConfigured = true;
	
	res.status(200).json({});
}

function handleCmds(commands) 
{
	for(var i = 0; i < commands.length; i++) {
		var cmd = commands[i];
		handleCmd(cmd);
	}	
}

function handleCmd(commandObj) 
{	
	switch(commandObj.command) {
		case 'CONNECT_BROKER':
			connectBroker(commandObj);
			break;
        case 'SET_REFERENCE':
            setReference(commandObj);
            break;
	}	
}

// connect to the IoT Broker
function connectBroker(cmd) 
{
    brokerURL = cmd.brokerURL;
    ngsi10client = new NGSIClient.NGSI10Client(brokerURL);
	console.log('connected to broker', cmd.brokerURL);
}

function setReference(cmd) 
{
    myReferenceURL = cmd.url;
	console.log('reference: ', myReferenceURL);
}


// handle the received results
function handleNotify(req, ctxObjects, res) 
{	
	console.log('handle notify');

	for(var i = 0; i < ctxObjects.length; i++) {
		console.log(ctxObjects[i]);
        processInputStreamData(ctxObjects[i]);
	}
}

// get the listening port number from the environment variables given by the FogFlow edge worker
var myport = process.env.myport;

// set up the NGSI agent to listen on 
NGSIAgent.setNotifyHandler(handleNotify);
NGSIAgent.setAdminHandler(handleAdmin);
NGSIAgent.start(myport, startApp);
process.on('SIGINT', function() {
	NGSIAgent.stop();	
	stopApp();
	
	process.exit(0);
});




