'use strict';

const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');
const fogfunction = require('./function.js');

var ngsi10client = null;
var brokerURL;
var outputs = [];
var threshold = 30;
var myReferenceURL;
var mySubscriptionId = null;
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
        console.log(cmd);
		handleCmd(cmd);
        console.log("handle next command");
	}	
	
	// send the updates in the buffer
	sendUpdateWithinBuffer();
}

function handleCmd(commandObj) 
{
    if (commandObj.command == 'CONNECT_BROKER') {        
		connectBroker(commandObj);
    } else if (commandObj.command == 'SET_OUTPUTS') {        
		setOutputs(commandObj);
    } else if (commandObj.command == 'SET_REFERENCE'){
        setReferenceURL(commandObj);
	}
}

// connect to the IoT Broker
function connectBroker(cmd) 
{
    brokerURL = cmd.brokerURL;
    ngsi10client = new NGSIClient.NGSI10Client(brokerURL);
	console.log('connected to broker', cmd.brokerURL);
}

function setReferenceURL(cmd) 
{
    myReferenceURL = cmd.url   
	console.log('your application can subscribe addtional inputs under the reference URL: ', myReferenceURL);
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

//
// query results from the assigned nearby IoT broker
//
function query(queryCtxReq, f)
{    
	if (ngsi10client == null) {
        console.log("=== broker is not configured for your query");
		return
	}

    ngsi10client.queryContext(queryCtxReq).then(f).catch( function(error) {
        console.log('failed to subscribe context');
    });
}

//
// send subscriptions to IoT broker
//
function subscribe(subscribeCtxReq)
{
	if (ngsi10client == null) {
        console.log("=== broker is not configured for your subscription");
		return
	}    
    
    subscribeCtxReq.reference =  myReferenceURL;
	
	console.log("================trigger my own subscription===================");
	console.log(subscribeCtxReq);
	console.log("===================");
    
    ngsi10client.subscribeContext(subscribeCtxReq).then( function(subscriptionId) {		
        console.log("subscription id = " + subscriptionId);   
		mySubscriptionId = subscriptionId;
    }).catch(function(error) {
        console.log('failed to subscribe context');
    });
}	
	
//
// publish context entities: 
//
function publish(ctxUpdate)
{
	buffer.push(ctxUpdate)
	
	if (ngsi10client == null) {
        console.log("=== broker is not configured for your update");
		return
	}
	
	for(var i=0; i<buffer.length; i++){
		var update = buffer[i];
				
    	ngsi10client.updateContext(update).then( function(data) {
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
        try {
            fogfunction.handler(ctxObjects[i], publish, query, subscribe);    
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


