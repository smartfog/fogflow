import express from 'express'
import multer from 'multer'
import httpProxy from 'http-proxy';
import fetch from 'node-fetch';
import bodyParser from 'body-parser'
import { promises as fs } from 'node:fs'
import { Low, JSONFile } from 'lowdb'

import socketio from 'socket.io'

import NGSIClient from './public/lib/ngsi/ngsiclient.cjs'
import NGSIAgent from './public/lib/ngsi/ngsiagent.cjs'
import NGSILDAgent from './public/lib/ngsi/LDngsiagent.cjs'

const globalConfigFile = JSON.parse(await fs.readFile('config.json'))

import rabbitmq from './rabbitmq.cjs';

import { fileURLToPath } from 'node:url';
import path from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));

var jsonParser = bodyParser.json();


var metadata_folder = './public/data/meta';
fs.mkdir(metadata_folder, { recursive: true })

var photo_folder = './public/data/photo';
fs.mkdir(photo_folder, { recursive: true })



const adapter = new JSONFile(metadata_folder + '/db.json');
const db = new Low(adapter);

await db.read()

db.data ||= {
    devices: {},
    subscriptions: {},
    operators: {},
    dockerimages: {},
    topologies: {},
    services: {},
    serviceintents: {},
    fogfunctions: {}
}

var send_loaded_intents = false

var app = express();

var ngsiProxy = httpProxy.createProxyServer();

var config = globalConfigFile.designer;

var taskMap = {};

// get the URL of the cloud broker

if (!('host_ip' in globalConfigFile.broker)) {
    globalConfigFile.broker.host_ip = globalConfigFile.my_hostip
}
var cloudBrokerURL = "http://" + globalConfigFile.broker.host_ip + ":" + globalConfigFile.broker.http_port
var ngsi10client = new NGSIClient.NGSI10Client(cloudBrokerURL + "/ngsi10");

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


if (!('host_ip' in globalConfigFile.master)) {
    globalConfigFile.master.host_ip = globalConfigFile.my_hostip
}
const masterURL = "http://" + globalConfigFile.master.host_ip + ":" + globalConfigFile.master.rest_api_port;


if (!('host_ip' in globalConfigFile.discovery)) {
    globalConfigFile.discovery.host_ip = globalConfigFile.my_hostip
}
const discoveryURL = "http://" + globalConfigFile.discovery.host_ip + ":" + globalConfigFile.discovery.http_port;


if (!('host_ip' in globalConfigFile.rabbitmq)) {
    globalConfigFile.rabbitmq.host_ip = globalConfigFile.my_hostip
}
const rabbitmq_ip = globalConfigFile.rabbitmq.host_ip;

const rabbitmq_port = globalConfigFile.rabbitmq.port || 5672;
const rabbitmq_user = globalConfigFile.rabbitmq.username || 'admin';
const rabbitmq_password = globalConfigFile.rabbitmq.password || 'mypass';

const rabbitmq_url = 'amqp://' + rabbitmq_user + ':' + rabbitmq_password + '@'
    + rabbitmq_ip + ':' + rabbitmq_port.toString();

console.log(config);

console.log(rabbitmq_url);
rabbitmq.Init(rabbitmq_url, handleInternalMessage, issueLoadedIntents);

function handleInternalMessage(jsonMsg) {
    console.log(jsonMsg.Type);

    var msgType = jsonMsg.Type;
    switch (msgType) {
        case 'MASTER_JOIN':

            break;
        case 'MASTER_LEAVE':

            break;
        case 'TASK_UPDATE':
            onTaskUpdate(jsonMsg)
            break;
    }
}

function onTaskUpdate(msg) {
    var payload = msg.PayLoad;
    var workerID = msg.From;
    if (payload.Status == 'removed') {
        removeTask(payload.TaskID)
    } else {
        updateTaskList(payload, workerID);
    }
}

function updateTaskList(updateMsg, fromWorker) {
    var taskID = updateMsg.TaskID;

    if (taskID in taskMap) {
        if ("Status" in updateMsg) {
            taskMap[taskID].Status = updateMsg.Status
        }
        if ("Info" in updateMsg) {
            taskMap[taskID].Info = updateMsg.Info
        }
    } else {
        taskMap[taskID] = updateMsg;
        taskMap[taskID]['Worker'] = fromWorker;
    }
}

function removeTask(taskID) {
    if (taskID in taskMap) {
        delete taskMap[taskID];
    }
}


function isEmpty(obj) {
    for (var prop in obj) {
        if (obj.hasOwnProperty(prop))
            return false;
    }

    return true;
}

// this is only triggered when starting the task designer
function issueLoadedIntents() {
    console.log("[RabbitMQ] is already connected")

    if (send_loaded_intents == false) {
        console.log("issue the loaded intents to Master");

        // existing service intents
        Object.keys(db.data.serviceintents).forEach(function (key) {
            var intent = db.data.serviceintents[key];
            intent.action = 'ADD';
            publishMetadata("ServiceIntent", intent);
        });

        // existing fog functions
        Object.keys(db.data.fogfunctions).forEach(function (key) {
            var fogfunction = db.data.fogfunctions[key];

            fogfunction.status = 'enabled';

            var intent = fogfunction.intent;
            intent.action = 'ADD';
            publishMetadata("ServiceIntent", intent);
        });

        send_loaded_intents = true;
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


app.use(function (req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    res.header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE");

    next();
});

// to server all static content
app.use(express.static(__dirname + '/public', { cache: false }));

// to receive and save uploaded image content
var storage = multer.diskStorage({
    destination: function (req, file, callback) {
        callback(null, './public/data/photo');
    },
    filename: function (req, file, callback) {
        console.log(file.fieldname);
        callback(null, file.fieldname);
    }
});
var upload = multer({ storage: storage }).any();
app.post('/data/photo', function (req, res) {
    upload(req, res, function (err) {
        if (err) {
            return res.end("Error uploading file.");
        }
        res.end("File is uploaded");
    });
});

//============= FogFlow API =================================

app.all("/ngsi10/*", function (req, res) {
    console.log('redirecting to ngsi-v1 broker', cloudBrokerURL);
    ngsiProxy.web(req, res, { target: cloudBrokerURL });
});

app.all("/ngsi-ld/*", function (req, res) {
    //console.log('redirecting to ngsi-ld broker');
    ngsiProxy.web(req, res, { target: cloudBrokerURL });
});

app.all("/ngsi9/*", function (req, res) {
    //console.log('redirecting to ngsi-v1 discovery');
    ngsiProxy.web(req, res, { target: discoveryURL });
});


//============= FogFlow API =================================

app.get('/info/master', async function (req, res) {
    try {
        var url = masterURL + "/status";
        const response = await fetch(url);
        const master = await response.json();
        res.json([master]);
    } catch (error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});

app.get('/info/broker', async function (req, res) {
    try {
        var url = discoveryURL + "/ngsi9/broker";
        const response = await fetch(url);
        const workers = await response.json();
        var workerList = Array.from(Object.values(workers));
        res.json(workerList);
    } catch (error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});

app.get('/info/worker', async function (req, res) {
    try {
        var url = masterURL + "/workers";
        const response = await fetch(url);
        const workers = await response.json();
        var workerList = Array.from(Object.values(workers));
        res.json(workerList);
    } catch (error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});

app.get('/info/task', async function (req, res) {
    var taskList = [];
    console.log(taskMap)
    for (const taskID of Object.keys(taskMap)) {
        var task = taskMap[taskID];
        task['ID'] = taskID;
        taskList.push(task);
    }

    res.json(taskList);
});

app.get('/info/task/:intent', async function (req, res) {
    var intentID = req.params.intent;

    var taskList = [];

    for (const taskID of Object.keys(taskMap)) {
        var task = taskMap[taskID];

        if (task['ServiceIntentID'] == intentID) {
            task['ID'] = taskID;
            taskList.push(task);
        }
    }

    res.json(taskList);
});

app.get('/info/type', async function (req, res) {
    try {
        var url = discoveryURL + "/etype";
        const response = await fetch(url);
        const typeList = await response.json();
        res.json(typeList);
    } catch (error) {
        console.log("failed to connect the master at ", url, '[ERROR CODE]', error.code);
        res.json([]);
    };
});


app.get('/operator', async function (req, res) {
    var operators = db.data.operators;
    res.json(operators);
});
app.get('/operator/:name', async function (req, res) {
    var name = req.params.name;
    var operator = db.data.operators[name];
    res.json(operator);
});

app.post('/operator', jsonParser, async function (req, res) {
    var operators = req.body;
    for (var i = 0; i < operators.length; i++) {
        var operator = operators[i];
        db.data.operators[operator.name] = operator
    }

    await db.write();

    res.sendStatus(200)
});

app.get('/dockerimage', async function (req, res) {
    var dockerimages = db.data.dockerimages;
    res.json(dockerimages);
});
app.get('/dockerimage/:operator', async function (req, res) {
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

    for (var i = 0; i < dockerimages.length; i++) {
        var dockerimage = dockerimages[i];
        console.log(dockerimage);
        db.data.dockerimages[dockerimage.name] = dockerimage;
    }

    await db.write();

    res.sendStatus(200)
});

app.get('/topology', async function (req, res) {
    var topologies = [];
    Object.values(db.data.topologies).forEach(topology => {
        topologies.push(topology)
    })

    res.json(topologies);
});
app.get('/topology/:name', async function (req, res) {
    var name = req.params.name;
    if (db.data.topologies.hasOwnProperty(name) == false) {
        res.json({});
        return
    }

    var topology = db.data.topologies[name];

    topology.dockerimages = [];
    topology.operators = [];

    // include the related docker images for the operators used by this topology
    for (var i = 0; i < topology.tasks.length; i++) {
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


app.get('/service', async function (req, res) {
    var services = db.data.services;
    res.json(services);
});
app.get('/service/:name', async function (req, res) {
    var name = req.params.name;
    var service = db.data.services[name];
    res.json(service);
});
app.post('/service', jsonParser, async function (req, res) {
    var services = req.body;

    for (var i = 0; i < services.length; i++) {
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


app.get('/intent', async function (req, res) {
    var serviceintents = db.data.serviceintents;
    res.json(serviceintents);
});
app.get('/intent/:id', async function (req, res) {
    var id = req.params.id;
    var serviceintent = db.data.serviceintents[id];
    res.json(serviceintent);
});
app.get('/intent/topology/:topology', async function (req, res) {
    var topology = req.params.topology;

    var intents = [];
    for (var i = 0; i < db.data.serviceintents.length; i++) {
        var intent = db.data.serviceintents[i];
        if (intent.topology == topology) {
            intents.push(intent)
        }
    }

    res.json(intents);
});
app.post('/intent', jsonParser, async function (req, res) {
    var serviceintent = req.body;

    if (isEmpty(serviceintent) == true) {
        res.sendStatus(200)
        return
    }

    console.log(serviceintent)
    db.data.serviceintents[serviceintent.id] = serviceintent;
    await db.write();

    serviceintent.action = 'ADD';
    publishMetadata("ServiceIntent", serviceintent);

    res.sendStatus(200)
});
app.delete('/intent/:id', async function (req, res) {
    var id = req.params.id;

    try {
        var serviceintent = db.data.serviceintents[id];
        if (serviceintent === undefined) {
            res.sendStatus(404)
            return;
        }
        serviceintent.action = 'DELETE';
        publishMetadata("ServiceIntent", serviceintent);

        delete db.data.serviceintents[id];
        await db.write();
        res.sendStatus(200)
    } catch (error) {
        console.log("Delete Intent API failed for [ID] ", id, ', [ERROR]', error.message);
        res.sendStatus(404)
    }
});


app.get('/fogfunction', async function (req, res) {
    var fogfunctions = db.data.fogfunctions;
    res.json(fogfunctions);
});
app.get('/fogfunction/:name', async function (req, res) {
    var name = req.params.name;
    var fogfunction = db.data.fogfunctions[name];
    res.json(fogfunction);
});

app.get('/fogfunction/:name/enable', async function (req, res) {
    var name = req.params.name;
    var fogfunction = db.data.fogfunctions[name];
    fogfunction.status = 'enabled';

    console.log("ENABLE fog function", name);

    var serviceintent = fogfunction.intent;
    serviceintent.action = 'ADD';
    publishMetadata("ServiceIntent", serviceintent);

    res.json(fogfunction);
});

app.get('/fogfunction/:name/disable', async function (req, res) {
    var name = req.params.name;
    var fogfunction = db.data.fogfunctions[name];
    fogfunction.status = 'disabled';

    console.log("DISABLE fog function", name);

    var serviceintent = fogfunction.intent;
    serviceintent.action = 'DELETE';
    publishMetadata("ServiceIntent", serviceintent);

    res.json(fogfunction);
});


app.post('/fogfunction', jsonParser, async function (req, res) {
    var fogfunctions = req.body;

    for (var i = 0; i < fogfunctions.length; i++) {
        var fogfunction = fogfunctions[i];

        fogfunction.status = 'enabled';

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
    try {
        var name = req.params.name;

        var fogfunction = db.data.fogfunctions[name]
        if (fogfunction === undefined) {
            res.sendStatus(404);
            return;
        }
        var serviceintent = fogfunction.intent
        serviceintent.action = 'DELETE';
        publishMetadata("ServiceIntent", serviceintent);

        delete db.data.topologies[name];
        delete db.data.fogfunctions[name];
        await db.write();

        res.sendStatus(200)
    } catch (error) {
        console.log("Delete Fogfunction API failed for [Name] ", name, ', [ERROR]', error.message);
        res.sendStatus(404)
    }
});

app.get('/task', async function (req, res) {
    var tasks = db.data.tasks;
    res.json(tasks);
});
app.get('/task/:intentID', async function (req, res) {
    var intentID = req.params.intentID;

    var tasks = [];
    //    for(var i=0; i<db.data.tasks.length; i++) {
    //        if db.data.tasks[i]
    //    }

    res.json(tasks);
});


app.get('/subscription', async function (req, res) {
    var subscriptions = db.data.subscriptions;
    res.json(subscriptions);
});
app.post('/subscription', jsonParser, async function (req, res) {
    var subscription = req.body;

    if (isEmpty(subscription) == true) {
        res.sendStatus(200)
        return
    }

    console.log(subscription)

    var headers = {};

    if (subscription.destination_broker == 'NGSI-LD') {
        headers["Content-Type"] = "application/json";
        headers["Destination"] = "NGSI-LD";
        headers["NGSILD-Tenant"] = subscription.tenant;
    } else if (subscription.destination_broker == 'NGSIv2') {
        headers["Content-Type"] = "application/json";
        headers["Destination"] = "NGSIv2";        
    }

     // issue the subscription to FogFlow Cloud Broker
    var subscribeCtxReq = {};
    subscribeCtxReq.entities = [{ type: subscription['entity_type'], isPattern: true }];
    subscribeCtxReq.reference = subscription['reference_url'];
    ngsi10client.subscribeContextWithHeaders(subscribeCtxReq, headers).then(function (subscriptionId) {
        console.log("subscription id = " + subscriptionId);
        db.data.subscriptions[subscriptionId] = subscription;
        db.write();

        res.sendStatus(200)
    }).catch(function (error) {
        console.log('failed to subscribe context, ', error);
        res.sendStatus(500)
    });
});
app.delete('/subscription/:id', async function (req, res) {
    var sid = req.params.id;
    
    if (db.data.subscriptions.hasOwnProperty(sid) == false) {
        res.sendStatus(404)
        return;
    }    

    // unsubscribe 
    ngsi10client.unsubscribeContext(sid).then(function (subscriptionId) {
        console.log("remove the subscription sid = " + subscriptionId);
        delete db.data.subscriptions[sid];
        db.write();
        res.sendStatus(200)
    }).catch(function (error) {
        console.log('failed to unsubscribe the specified subscription');
        res.sendStatus(500)
    });
});
 

app.get('/device', async function (req, res) {
    var devices = db.data.devices;
    res.json(devices);
});
app.post('/device', jsonParser, async function (req, res) {
    var deviceObj = req.body;
    
    if (deviceObj.hasOwnProperty('id')) {
        var did = deviceObj.id;        
    } else {
        var did = deviceObj.type + uuid();        
    }
        
    if (isEmpty(deviceObj) == true) {
        res.sendStatus(200)
        return
    }

    console.log(deviceObj)

    deviceObj.entityId = {
           id: deviceObj.id,
           type: deviceObj.type,
           isPattern: false		
	};

    // create the device entity in FogFlow Cloud Broker
    ngsi10client.updateContext(deviceObj).then(function (data) {
        console.log(data);
        db.data.devices[did] = deviceObj;
        db.write();
        res.sendStatus(200)        
    }).catch(function (error) {
        console.log('failed to register the new device object');
        res.sendStatus(500)        
    });    
});
app.delete('/device/:id', async function (req, res) {
    var did = req.params.id;
    if (db.data.devices.hasOwnProperty(did) == false) {
        res.sendStatus(404)
        return;
    }

    // delete the associated ngsi entity for this device
    var entityid = {
        id: did,
        isPattern: false
    };
     
    ngsi10client.deleteContext(entityid).then(function (subscriptionId) {
        console.log('remove the device ', did);
        delete db.data.devices[did];
        db.write();
        res.sendStatus(200)
    }).catch(function (error) {
        console.log('failed to delete the corresponding entity');
        res.sendStatus(500)
    });
});



app.get('/config.js', function (req, res) {
    res.setHeader('Content-Type', 'application/json');
    var data = 'var config = ' + JSON.stringify(config) + '; '
    res.end(data);
});


// return the docker images for a given operator name
function getDockerImages(name) {
    var dockerimages = [];
    Object.values(db.data.dockerimages).forEach(dockerimage => {
        if (dockerimage.operatorName == name) {
            dockerimages.push(dockerimage);
        }
    });

    return dockerimages;
}


// return the operator for a given operator name
function getOperator(name) {
    if (name in db.data.operators) {
        return db.data.operators[name];
    } else {
        return null
    };
}

// publish the created metadata related to service orchestration
function publishMetadata(dType, dObject) {
    var jsonMsg = {
        Type: dType,
        RoutingKey: "orchestration.",
        From: "designer",
        PayLoad: dObject
    };

    rabbitmq.Publish(jsonMsg);
    //console.log("Published: ", JSON.stringify(jsonMsg));
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
webServer = app.listen(config.webSrvPort, function () {
    console.log("HTTP-based web server is listening on port ", config.webSrvPort);
});

var io = socketio.listen(webServer);

io.on('connection', function (client) {
    console.log('a client is connecting');
    client.on('subscriptions', function (subList) {
        console.log(subList);
        for (var i = 0; subList && i < subList.length; i++) {
            var sid = subList[i];
            subscriptions[sid] = client;
        }
    });

    client.on('disconnect', function () {
        console.log('disconnected');
        //remove the subscriptions associated with this socket
        for (var sid in subscriptions) {
            if (subscriptions[sid] == client) {
                delete subscriptions[sid];
            }
        }
    });
});
