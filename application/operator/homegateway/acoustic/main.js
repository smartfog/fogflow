'use strict';

const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');
const fogfunction = require('./function.js');

var ngsi10client = null;
var brokerURL;
var outputs = [];
var threshold = 30;
var isConfigured = false;

var buffer = [];

function startApp() 
{
    console.log('start to receive input data streams via a listening port');
}

function stopApp() 
{
	console.log('clean up the app');
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
	
	// send the updates in the buffer
	sendUpdateWithinBuffer();
}

function handleCmd(commandObj) 
{	
	switch(commandObj.command) {
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

function sendUpdateWithinBuffer()
{
	for(var i=0; i<buffer.length; i++){
		var tmp = buffer[i];
		
		if (tmp.outputIdx > 0) {
			tmp.ctxObj.entityId.id = outputs[i].id;
			tmp.ctxObj.entityId.type = outputs[i].type;
		}
		
    		ngsi10client.updateContext(tmp.ctxObj).then( function(data) {
			console.log('======send update======');
    		    console.log(data);
    		}).catch(function(error) {
			console.log(error);
    		    console.log('failed to update context');
    		});  		
	}
	
	buffer= [];
}

function publish(ctxUpdate, index)
{
    console.log("publish an update: ", ctxUpdate, " at outputIndex ", index);    
	
	buffer.push({ctxObj: ctxUpdate, outputIdx: index})
	
	if (ngsi10client == null) {
		return
	}
	
	for(var i=0; i<buffer.length; i++){
		var tmp = buffer[i];
		
		if (tmp.outputIdx >=0) {
			tmp.ctxObj.entityId.id = outputs[i].id;
			tmp.ctxObj.entityId.type = outputs[i].type;
			tmp.ctxObj.entityId.isPattern = false;            
		}
		
        ngsi10client.updateContext(tmp.ctxObj).then( function(data) {
		console.log('======send update======');
    		console.log(data);
        }).catch(function(error) {
			console.log(error);
    		console.log('failed to update context');
        });  		
	}
	
	buffer= [];
}

// handle the received results
function handleNotify(req, ctxObjects, res) 
{	
	console.log('handle notify');
	for(var i = 0; i < ctxObjects.length; i++) {
		console.log(ctxObjects[i]);
        fogfunction.handler(ctxObjects[i], publish);
	}
}

// get the listening port number from the environment variables given by the FogFlow edge worker
var myport = process.env.myport;

// set up the NGSI agent to listen on 
NGSIAgent.setNotifyHandler(handleNotify);
NGSIAgent.setAdminHandler(handleAdmin);
NGSIAgent.start(myport, startApp);

fogfunction.handler(null, publish);

process.on('SIGINT', function() {
	NGSIAgent.stop();	
	stopApp();
	
	process.exit(0);
});


