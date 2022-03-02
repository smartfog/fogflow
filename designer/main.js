import express from 'express'
import multer from 'multer'
//import http from 'http'
import fetch from 'node-fetch';
//import axios from 'axios'
import bodyParser from 'body-parser'
import {promises as fs} from 'node:fs'
import { Low, JSONFile } from 'lowdb'

import socketio from 'socket.io'

import NGSIAgent from './public/lib/ngsi/ngsiagent.cjs'
import NGSILDAgent from './public/lib/ngsi/LDngsiagent.cjs'

const globalConfigFile = JSON.parse(await fs.readFile('config.json'))

import rabbitmq  from './rabbitmq.cjs';

import {fileURLToPath} from 'node:url';
import path from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));

var jsonParser = bodyParser.json();
fs.mkdir('public/data/meta/', { recursive: true })
const adapter = new JSONFile('public/data/meta/db.json');
const db = new Low(adapter);

await db.read()

db.data ||= {   operators: {}, 
                dockerimages: {}, 
                topologies: {},
                services: {},  
                serviceintents: {}, 
                fogfunctions: {} } 

var app = express();

var config = globalConfigFile.designer;

var masterList = [];

// get the URL of the cloud broker

if ( !('host_ip' in globalConfigFile.broker)) {
    globalConfigFile.broker.host_ip = globalConfigFile.my_hostip    
}
var cloudBrokerURL = "http://" + globalConfigFile.broker.host_ip + ":" + globalConfigFile.broker.http_port + "/ngsi10"

if (config.host_ip) {
    config.agentIP = config.host_ip;        
} else {
    config.agentIP = globalConfigFile.my_hostip;    
}

config.agentPort = globalConfigFile.designer.agentPort;

//Set NGSILD agent port
config.ldAgentPort = globalConfigFile.designer.ldAgentPort;

config.discoveryURL = './ngsi9';
config.brokerURL = './ngsi10';
config.LdbrokerURL = './ngsi-ld';
config.webSrvPort = globalConfigFile.designer.webSrvPort;


if ( !('host_ip' in globalConfigFile.master)) {
    globalConfigFile.master.host_ip = globalConfigFile.my_hostip    
}
const masterURL = "http://" + globalConfigFile.master.host_ip + ":" + globalConfigFile.master.rest_api_port;


if ( !('host_ip' in globalConfigFile.discovery)) {
    globalConfigFile.discovery.host_ip = globalConfigFile.my_hostip    
}
const discoveryURL = "http://" + globalConfigFile.discovery.host_ip + ":" + globalConfigFile.discovery.http_port;


if ( !('host_ip' in globalConfigFile.rabbitmq)) {
    globalConfigFile.rabbitmq.host_ip = globalConfigFile.my_hostip    
}
const rabbitmq_ip = globalConfigFile.rabbitmq.host_ip; 

const rabbitmq_port = globalConfigFile.rabbitmq.port || 5672;
const rabbitmq_user =  globalConfigFile.rabbitmq.username || 'admin';
const rabbitmq_password = globalConfigFile.rabbitmq.password || 'mypass';

const rabbitmq_url = 'amqp://' + rabbitmq_user + ':' + rabbitmq_password + '@' 
                + rabbitmq_ip + ':' + rabbitmq_port.toString();

console.log(config);

console.log(rabbitmq_url);
rabbitmq.Init(rabbitmq_url, handleInternalMessage);

function handleInternalMessage(jsonMsg) {
    console.log(jsonMsg);
    
    var msgType = jsonMsg.Type;
    switch(msgType) {
        case 'MASTER_JOIN':
        
            break;
        case 'MASTER_LEAVE':
        
            break;        
    }        
}

function uuid() {
    var uuid = "",
        i, random;
    for (i = 0; i < 32; i++) {
        random = Math.random() * 16 | 0;
        if (i == 8 || i == 12 || i == 16 || i == 20) {
            uuid += "-"
        }
        uuid += (i == 12 ? 4 : (i == 16 ? (random & 3 | 8) : random)).toString(16);
    }

    return uuid;
}

// all subscriptions that expect data forwarding
var subscriptions = {};


app.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

// to server all static content
app.use(express.static(__dirname + '/public', { cache: false }));

// to receive and save uploaded image content
var storage = multer.diskStorage({
    destination: function(req, file, callback) {
        callback(null, './public/data/photo');
    },
    filename: function(req, file, callback) {
        console.log(file.fieldname);
        callback(null, file.fieldname);
    }
});
var upload = multer({ storage: storage }).any();
app.post('/photo', function(req, res) {
    upload(req, res, function(err) {
        if (err) {
            return res.end("Error uploading file.");
        }
        res.end("File is uploaded");
    });
});


//============= FogFlow API =================================

app.get('/info/master', async function(req, res) {    
    try {
        var url = masterURL + "/status";        
        const response = await fetch(url);
        const master = await response.json();   
        res.json([master]);    
    } catch(error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});


app.get('/info/worker', async function(req, res) {   
    try {
        var url = masterURL + "/workers";        
        const response = await fetch(url);
        const workers = await response.json(); 
        var workerList = Array.from(Object.values(workers));          
        res.json(workerList);
    } catch(error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});

app.get('/info/task', async function(req, res) {    
    try {
        var url = masterURL + "/tasks";        
        const response = await fetch(url);
        const taskList = await response.json();   
        res.json(taskList);
    } catch(error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});

app.get('/info/type', async function(req, res) {   
    try {
        var url = discoveryURL + "/etype";        
        const response = await fetch(url);
        const typeList = await response.json();   
        res.json(typeList);
    } catch(error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});


app.get('/operator', async function(req, res) {    
    var operators = db.data.operators;         
    res.json(operators);           
});
app.get('/operator/:name', async function(req, res) {    
    var name = req.params.name;
    var operator = db.data.operators[name];
    res.json(operator);
});

app.post('/operator', jsonParser, async function (req, res) {
    var operators = req.body;    
    for(var i=0; i<operators.length; i++){
        var operator = operators[i];
        db.data.operators[operator.name] = operator
    }
    
    await db.write();    
    
    res.sendStatus(200)    
});

app.get('/dockerimage', async function(req, res) {    
    var dockerimages = db.data.dockerimages;    
    res.json(dockerimages);             
});
app.get('/dockerimage/:operator', async function(req, res) {    
    var operator = req.params.operator;
    var imageList = [];
 
    var dockerimages = db.data.dockerimages;   
    Object.values(dockerimages).forEach(dockerimage => {
        if (dockerimage.operatorName == operator) {
            imageList.push(dockerimage)    
        }        
    })
        
    res.json(imageList);
});

app.post('/dockerimage', jsonParser, async function (req, res) {
    var dockerimages = req.body;    
    console.log(dockerimages);
    
    for(var i=0; i<dockerimages.length; i++){
        var dockerimage = dockerimages[i];
        console.log(dockerimage);
        db.data.dockerimages[dockerimage.name] = dockerimage;     
    }
    
    await db.write();
    
    res.sendStatus(200)    
});

app.get('/topology', async function(req, res) {    
    var topologies = [];    
    Object.values(db.data.topologies).forEach(topology => {
        topologies.push(topology)
    })
    
    res.json(topologies);    
});
app.get('/topology/:name', async function(req, res) {    
    var name = req.params.name;
    var topology = db.data.topologies[name];
    
    topology.dockerimages = [];
    topology.operators = [];
        
    // include the related docker images for the operators used by this topology
    for(var i=0; i<topology.tasks.length; i++){        
        var task = topology.tasks[i];
        var name = task.operator;
        
        var dockerimages = getDockerImages(name);        
        var operator = getOperator(name);                
        if (operator != null) {
            operator.dockerimages = dockerimages;            
            topology.operators.push(operator);            
        }
    }    
    
    res.json(topology);
});


app.get('/service', async function(req, res) {    
    var services = db.data.services;    
    res.json(services);    
});
app.get('/service/:name', async function(req, res) {    
    var name = req.params.name;
    var service = db.data.services[name];
    res.json(service);
});
app.post('/service', jsonParser, async function (req, res) {
    var services = req.body;    
    
    for(var i=0; i<services.length; i++){
        var service = services[i];
        console.log(service);   
        var name = service.topology.name     
        db.data.services[name] = service;
        db.data.topologies[name] = service.topology;
    }
    
    await db.write();    
    
    res.sendStatus(200)     
});
app.delete('/service/:name', async function (req, res) {
    var name = req.params.name;

    delete db.data.services[name];
    delete db.data.topologies[name];    
    await db.write();  
        
    res.sendStatus(200)      
});


app.get('/intent', async function(req, res) {    
    var serviceintents = db.data.serviceintents;    
    res.json(serviceintents);  
});
app.get('/intent/:id', async function(req, res) {    
    var id = req.params.id;
    var serviceintent = db.data.serviceintents[id];
    res.json(serviceintent);
});
app.post('/intent', jsonParser, async function (req, res) {
    var serviceintent = req.body;    
    console.log(serviceintent)
    db.data.serviceintents[serviceintent.id] = serviceintent;    
    await db.write();    
    
    serviceintent.action = 'ADD';    
    publishMetadata("ServiceIntent", serviceintent);                            
    
    res.sendStatus(200)       
});
app.delete('/intent/:id', async function (req, res) {
    var id = req.params.id;
    
    var serviceintent = db.data.serviceintents[id];
    serviceintent.action = 'DELETE';
    publishMetadata("ServiceIntent", serviceintent);       
    
    delete db.data.serviceintents[id];    
    await db.write();    
    
    res.sendStatus(200)        
});


app.get('/fogfunction', async function(req, res) {    
    var fogfunctions = db.data.fogfunctions;    
    res.json(fogfunctions);  
});
app.get('/fogfunction/:name', async function(req, res) {    
    var name = req.params.name;
    var fogfunction = db.data.fogfunctions[name];
    res.json(fogfunction);
});
app.post('/fogfunction', jsonParser, async function (req, res) {
    var fogfunctions = req.body;    
    
    for(var i=0; i<fogfunctions.length; i++){
        var fogfunction = fogfunctions[i];
        console.log(fogfunction);      
        db.data.fogfunctions[fogfunction.name] = fogfunction;
        db.data.topologies[fogfunction.name] = fogfunction.topology;
        
        if (fogfunction.intent.hasOwnProperty('id') == false) {
            var uid = uuid();
            var sid = 'ServiceIntent.' + uid;            
            fogfunction.intent.id = sid;
        }
        
        var serviceintent = fogfunction.intent;
        serviceintent.action = 'ADD';
        publishMetadata("ServiceIntent", serviceintent);                             
    }
    
    await db.write();    
    
    res.sendStatus(200)  
});
app.delete('/fogfunction/:name', async function (req, res) {
    var name = req.params.name;

    var fogfunction = db.data.fogfunctions[name]
    
    var serviceintent = fogfunction.intent
    serviceintent.action = 'DELETE';    
    publishMetadata("ServiceIntent", serviceintent);                                
    
    delete db.data.topologies[name];        
    delete db.data.fogfunctions[name];    
    await db.write();        
    
    res.sendStatus(200)        
});

app.get('/config.js', function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    var data = 'var config = ' + JSON.stringify(config) + '; '
    res.end(data);
});


// return the docker images for a given operator name
function getDockerImages(name)
{
    var dockerimages = [];
    Object.values(db.data.dockerimages).forEach(dockerimage => {
        if (dockerimage.operatorName == name) {
            dockerimages.push(dockerimage);    
        }        
    });
    
    return dockerimages;
}


// return the operator for a given operator name
function getOperator(name)
{
    if (name in db.data.operators) {
        return db.data.operators[name];    
    } else {
        return null
    };
}

// publish the created metadata related to service orchestration
function publishMetadata(dType, dObject)
{
    var jsonMsg = { Type: dType, 
                    RoutingKey: "orchestration.", 
                    From: "designer", 
                    PayLoad: dObject 
                };
        
    rabbitmq.Publish(jsonMsg);    
}


// handle the received results
function handleNotify(req, ctxObjects, res) {
    console.log('handle notify');
    var sid = req.body.subscriptionId;
    console.log(sid);
    if (sid in subscriptions) {
        for (var i = 0; i < ctxObjects.length; i++) {
            console.log(ctxObjects[i]);
            var client = subscriptions[sid];
            client.emit('notify', { 'subscriptionID': sid, 'entities': ctxObjects[i] });
        }
    }
}


NGSIAgent.setNotifyHandler(handleNotify);
NGSIAgent.start(config.agentPort);

NGSILDAgent.setNotifyHandler(handleNotify);
NGSILDAgent.start(config.ldAgentPort);

var webServer;
webServer = app.listen(config.webSrvPort, function() {
    console.log("HTTP-based web server is listening on port ", config.webSrvPort);
});

var io = socketio.listen(webServer);

io.on('connection', function(client) {
    console.log('a client is connecting');
    client.on('subscriptions', function(subList) {
	console.log(subList);
        for (var i = 0; subList && i < subList.length; i++) {
            sid = subList[i];
            subscriptions[sid] = client;
        }
    });
    client.on('disconnect', function() {
        console.log('disconnected');

        //remove the subscriptions associated with this socket
        for (sid in subscriptions) {
            if (subscriptions[sid] == client) {
                delete subscriptions[sid];
            }
        }
    });
});
