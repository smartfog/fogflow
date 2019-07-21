'use strict';

const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');

var ngsi10client;
var brokerURL;
var outputs = [];
var threshold = 30;
var isConfigured = false;

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
    var type = data.entityId.type;
    console.log('type ', type);
    
    // do the internal data processing
    if (type == 'PowerPanel'){
        anomalyDetection(data);        
    } else if (type == 'Rule') {
        updateRule(data);
    }    
}

// simple algorithm to detect anomaly events
function anomalyDetection(msg) 
{	
	var watts = msg.attributes.usage.value;
	var deviceID = msg.entityId.id;
	var shopID =  msg.metadata.shop.value;	
	var location = msg.metadata.location;
	
	console.log('============+++++++++++Current usage ', watts, ', current thrshold = ', threshold)
	
	if(watts > threshold) { // detect an anomaly event
		// publish the detected event
        var anomaly = {};
        
        var now = new Date();
        anomaly['when'] = now.toISOString();
        anomaly['whichpanel'] = deviceID;    
        anomaly['whichshop'] = shopID;            
        anomaly['where'] = location;
        anomaly['usage'] = watts;            
        
        updateContext(anomaly)                    
    } 
}

function updateRule(ruleObj)
{
    threshold = ruleObj.attributes.threshold.value;
    console.log('update the threshold to ', threshold);
}

// update context for streams
function updateContext(anomaly) 
{
    if (isConfigured == false) {
        console.log('the task is not configured yet!!!');
        return;
    }
        
    var ctxObj = {};
    
    ctxObj.entityId = {};
    
    var outputStream = outputs[0];
    
    ctxObj.entityId.id = outputStream.id;
    ctxObj.entityId.type = outputStream.type;
    ctxObj.entityId.isPattern = false;
    
    ctxObj.attributes = {};
       
    ctxObj.attributes.when = {        
        type: 'string',
        value: anomaly['when']
    };
    ctxObj.attributes.whichpanel = {
        type: 'string',
        value: anomaly['whichpanel']
    };  
        
    ctxObj.attributes.shop = {
        type: 'string',
        value: anomaly['whichshop']
    };  
    ctxObj.attributes.where = {
        type: 'object',
        value: anomaly['where']
    };  
    ctxObj.attributes.usage = {
        type: 'integer',
        value: anomaly['usage']
    };
        
    ctxObj.metadata = {};        
    ctxObj.metadata.shop = {
        type: 'string',
        value: anomaly['whichshop']
    };  
              
    ngsi10client.updateContext(ctxObj).then( function(data) {
        console.log('======send update======');
        console.log(ctxObj);
        console.log(data);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to update context');
    });    
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
    switch(commandObj.command){
        case 'CONNECT_BROKER':
            connectBroker(commandObj);
            break;
        case 'SET_OUTPUTS':
            setOutputs(commandObj);
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

function setOutputs(cmd) 
{
    var outputStream = {};
    outputStream.id = cmd.id;
    outputStream.type = cmd.type;

    outputs.push(outputStream);

    console.log('output has been set: ', cmd);
}

// handle the received results
function handleNotify(req, ctxObjects, res) 
{    
    console.log('handle notify');

	for(var i = 0; i < ctxObjects.length; i++) {
		console.log(ctxObjects[i]);
        
        try {
            processInputStreamData(ctxObjects[i]);        
        } catch (error) {
            console.log(error)
        }        
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


