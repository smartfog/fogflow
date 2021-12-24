var express =   require('express');
var multer  =   require('multer');
var https = require('https');
const bodyParser = require('body-parser');
var axios = require('axios')
var dgraph = require('./dgraph.js');
var amqp = require('./rabbitmq.js');
var jsonParser = bodyParser.json();
var config_fs_name = './config.json';


var args = process.argv.slice(2);
if (args.length > 0) {
    config_fs_name = args[0];
}
console.log(config_fs_name);

var fs = require('fs');
var globalConfigFile = require(config_fs_name)
var app = express();
var NGSIAgent = require('./public/lib/ngsi/ngsiagent.js');
var NGSIClient = require('./public/lib/ngsi/ngsiclient.js');
var NGSILDAgent = require('./public/lib/ngsi/LDngsiagent.js');

var config = globalConfigFile.designer;

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
config.webSrvPort = globalConfigFile.designer.webSrvPort


console.log(config);

dgraph.Init();
amqp.amqpConnection();
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
        callback(null, './public/photo');
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


// create a new intent to trigger the corresponding service topology
app.post('/intent', jsonParser, function(req, res) {
    intent = req.body

    var intentCtxObj = {};

    var uid = uuid();

    intentCtxObj.entityId = {
        id: 'ServiceIntent.' + uid,
        type: 'ServiceIntent',
        isPattern: false
    };

    intentCtxObj.attributes = {};
    intentCtxObj.attributes.status = { type: 'string', value: 'enabled' };
    intentCtxObj.attributes.intent = { type: 'object', value: intent };

    intentCtxObj.metadata = {};
    intentCtxObj.metadata.location = intent.geoscope;

    ngsi10client = new NGSIClient.NGSI10Client(cloudBrokerURL);
    ngsi10client.updateContext(intentCtxObj).then(function(data) {
        console.log('======create intent======');
        console.log(data);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to create intent');
    });

    // prepare the response
    reply = { 'id': intentCtxObj.entityId.id, 'outputType': 'Result' }
    res.json(reply);
});

/*
  api to create fogFlow internal entities
*/

function getResult(filterBy, objList) {
    console.log("objlist is ---- ",objList);
    var arr = [];
    for(var i in objList){
            if (objList[i].hasOwnProperty("attribute") && Object.keys(objList[i].attribute).length === 0) {
                console.log("aaaa",objList[i]);
                continue;}
            var fType = objList[i]['internalType'];
            if (fType != undefined && filterBy == fType){
                    attrObj = JSON.parse(objList[i].attribute) 
                    attrObj.uid= objList[i].uid
                    arr.push(attrObj);
                    continue;
            }
    }
    return arr;
  }

app.post('/internal/updateContext', jsonParser, async function (req, res) {
    //var updateContextReq = await req.body
    let updateContextReq = Object.assign({}, await req.body);
    console.log("****************** update",updateContextReq);


    if (updateContextReq.updateAction == "DELETE") {
        console.log("delete entity is ",updateContextReq);
        amqp.amqpPub(updateContextReq)
        let tmpVar = JSON.parse(JSON.stringify(updateContextReq));
        await dgraph.DeleteNodeById(tmpVar.uid);

        
    } else if (updateContextReq.updateAction == "UPDATE") {
        //console.log("main js obj  ++++ ",updateContextReq)
        let tmpVar1 = JSON.parse(JSON.stringify(updateContextReq));
        await amqp.amqpPub(tmpVar1)
        let tmpVar = JSON.parse(JSON.stringify(updateContextReq));
        await dgraph.WriteEntity(tmpVar)
    }
    res.send("")
});

app.post('/internal/getContext', jsonParser, async function (req, res) {
    var queryContext = await req.body
    var dgraphOp = await dgraph.QueryJsonWithType(queryContext.internalType)
    res.send({data:getResult(queryContext.internalType,dgraphOp.contextElements)})
});

app.post('/masterNotify', jsonParser, async function (req, res) {
    //dgraphSendToMaster();
    //res.json({"status":"DONE"})
    if(global.isAmqpUp){
        dgraphSendToMaster();
        res.json({"status":"DONE"})
    }
    else {
        
        res.status(404).send({"msg":"rabbitMQ is not reachable"});
    }
});

// to remove an existing intent
app.delete('/intent', jsonParser, function(req, res) {
    eid = req.body.id

    var entityid = {
        id: eid,
        isPattern: false
    };

    ngsi10client = new NGSIClient.NGSI10Client(cloudBrokerURL);
    ngsi10client.deleteContext(entityid).then(function(data) {
        console.log('======delete intent======');
        console.log(data);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to delete intent');
    });

    res.end("OK");
});


app.get('/config.js', function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    var data = 'var config = ' + JSON.stringify(config) + '; '
    res.end(data);
});


/*
// fetch the requested URL from the edge node within the internal network
app.get('/proxy', function(req, res) {
    console.log(req.query.url);

    if (req.query.url) {
        request(req.query.url).pipe(res);
    }
});
*/

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

async function dgraphSendToMaster(){
    var dgraphResult = await dgraph.QueryJsonWithType('all')
    if (dgraphResult.hasOwnProperty("contextElements")){
        dgraphData = dgraphResult.contextElements
        for (var i in dgraphData) {
            if (dgraphData[i].attribute == undefined) {
                continue;
              }
            dgraphData[i].attribute = JSON.parse(dgraphData[i].attribute)
            amqp.amqpPub(dgraphData[i])
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


//var io = require('socket.io')();
var io = require('socket.io').listen(webServer);

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
