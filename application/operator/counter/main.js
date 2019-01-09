'use strict';

const NGSIClient = require('./ngsi/ngsiclient.js');
const NGSIAgent = require('./ngsi/ngsiagent.js');

var ngsi10client;
var brokerURL;
var outputs = [];
var tBegin = Date.now();
console.log(tBegin);

var isConfigured = false;
var tInterval = 5000; // 5 seconds
var totalNumAnomaly = 0; // total number of detected anomaly events

function startApp() 
{
    console.log('start to receive input data streams via a listening port');
    setInterval(onTimer, tInterval);
}

function stopApp() 
{                         
    console.log('clean up the app');
}

// process the input data stream and generate output stream
function processInputStreamData(data) 
{
    // do the internal data processing    
    counter(data);
}

function onTimer()
{
    var now = Date.now();
    
    if( (now - tBegin) >= tInterval ) {
        // publish the total number of anomaly events in the current time window
        var stat = {};
        
        stat['time'] = tBegin; // time.toISOString(); var time = new Date(now);
        stat['counter'] = totalNumAnomaly;    
        
        updateContext(stat)                
        
        tBegin = tBegin + tInterval;
        totalNumAnomaly = 0;
    } 
}

function counter(msg) 
{    
    var now = Date.now();
    
    if( (now - tBegin) >= tInterval ) {
        // publish the total number of anomaly events in the current time window
        var stat = {};
        
        stat['time'] = tBegin; //time.toISOString(); var time = new Date(now);
        stat['counter'] = totalNumAnomaly;    
        
        updateContext(stat)                
        
        tBegin = tBegin + tInterval;
        totalNumAnomaly = 1;
    } else {
        totalNumAnomaly = totalNumAnomaly + 1;
    }
}

// update context for streams
function updateContext(stat) 
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
            
    ctxObj.attributes.time = {        
        type: 'integer',
        value: stat.time
    };
    ctxObj.attributes.counter = {
        type: 'integer',
        value: stat.counter
    };  
    
    console.log(ctxObj);    
    ngsi10client.updateContext(ctxObj).then( function(data) {
        console.log('======send update======, ');
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


// handle the received results
function handleNotify(req, ctxObjects, res) 
{    
    for(var i = 0; i < ctxObjects.length; i++) {
        //console.log(ctxObjects[i]);
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



