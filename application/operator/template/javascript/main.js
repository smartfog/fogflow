'use strict';

const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');

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


function startApp() 
{
    console.log('your application starts listening on port ',  myport);
    
    // to add the initialization part of your application logic here    
}

function stopApp() 
{
	console.log('clean up the app');
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


// global variables to be configured 
var ngsi10client;
var brokerURL;
var myReferenceURL;
var outputs = [];
var isConfigured = false;

// handle the configuration commands issued by FogFlow
function handleAdmin(req, commands, res) 
{	
	handleCmds(commands);
	
	isConfigured = true;
	
	res.status(200).json({});
}

// handle all configuration commands
function handleCmds(commands) 
{
	for(var i = 0; i < commands.length; i++) {
		var cmd = commands[i];
		handleCmd(cmd);
	}	
}

// handle each configuration command accordingly
function handleCmd(commandObj) 
{	
	switch(commandObj.command) {
		case 'CONNECT_BROKER':
			connectBroker(commandObj);
			break;
		case 'SET_OUTPUTS':
			setOutputs(commandObj);
			break;
		case 'SET_YOUR_REFERENCE':
			setReferenceURL(commandObj);
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

function setReferenceURL(cmd) 
{
    myReferenceURL = cmd.referenceURL   
	console.log('your application can subscribe addtional inputs under the reference URL: ', myReferenceURL);
}


// handle all received NGSI notify messages
function handleNotify(req, ctxObjects, res) 
{	
	console.log('handle notify');
	for(var i = 0; i < ctxObjects.length; i++) {
		console.log(ctxObjects[i]);
        fogfunction.handler(ctxObjects[i], publish);
	}
}

// process the input data stream accordingly and generate output stream
function processInputStreamData(data) 
{
	var type = data.entityId.type;
	console.log('type ', type);
	
	// do the internal data processing
	if (type == 'PowerPanel'){
		// to handle this type of input
	} else if (type == 'Rule') {
		// to handle this type of input
	}	
}



