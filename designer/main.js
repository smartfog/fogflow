var request = require('request');
var express =   require('express');
var multer  =   require('multer');
var https = require('https');

var config_fs_name = './config.json';


var args = process.argv.slice(2);
if (args.length > 0) {
    config_fs_name = args[0];
} 
console.log(config_fs_name);

var fs = require('fs');
var globalConfigFile = require(config_fs_name)
var app     =   express();
var NGSIAgent = require('./public/lib/ngsi/ngsiagent.js');
var NGSIClient = require('./public/lib/ngsi/ngsiclient.js');

var config = globalConfigFile.designer;

// set the agent IP address from the environment variable
config.agentIP = 'designer';
config.agentPort = globalConfigFile.designer.agentPort; 

config.discoveryURL = './ngsi9';
config.brokerURL = './ngsi10';    

config.webSrvPort = globalConfigFile.designer.webSrvPort

console.log(config);


// all subscriptions that expect data forwarding
var subscriptions = {};

// to server all static content
app.use(express.static(__dirname + '/public', {cache: false}));

app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

// to receive and save uploaded image content
var storage = multer.diskStorage({
  destination: function (req, file, callback) {
    callback(null, './public/photo');
  },
  filename: function (req, file, callback) {
    console.log(file.fieldname);    
    callback(null, file.fieldname);
  }
});
var upload = multer({ storage : storage}).any();
app.post('/photo',function(req, res){
    upload(req, res, function(err) {
        if(err) {
            return res.end("Error uploading file.");
        }
        res.end("File is uploaded");
    });
});

app.get('/config.js', function(req, res){
	res.setHeader('Content-Type', 'application/json');		
	var data = 'var config = ' + JSON.stringify(config) + '; '	
	res.end(data);
});

// fetch the requested URL from the edge node within the internal network
app.get('/proxy', function(req, res){
    console.log(req.query.url);    
    
    if(req.query.url) {
        request(req.query.url).pipe(res);
    }
});


// handle the received results
function handleNotify(req, ctxObjects, res) {	
	console.log('handle notify');
    var sid = req.body.subscriptionId;
    console.log(sid);
    if(sid in subscriptions) {
        for(var i = 0; i < ctxObjects.length; i++) {
            console.log(ctxObjects[i]);
            var client = subscriptions[sid];
            client.emit('notify', ctxObjects[i]);
        }
    }
}

NGSIAgent.setNotifyHandler(handleNotify);
NGSIAgent.start(config.agentPort);

var webServer;
webServer = app.listen(config.webSrvPort, function(){
    console.log("HTTP-based web server is listening on port ", config.webSrvPort);
});
    

var io = require('socket.io').listen(webServer);

io.on('connection', function (client) {
    console.log('a client is connecting');       
    client.on('subscriptions', function (subList) {
        console.log(subList);
        for(var i=0; subList && i<subList.length; i++){
            sid = subList[i];
            subscriptions[sid] = client;
        }
    });
    client.on('disconnect', function () {
        console.log('disconnected');
        
        //remove the subscriptions associated with this socket
        for(sid in subscriptions) {
            if(subscriptions[sid] == client) {
                delete subscriptions[sid];
            }
        }
    });
});
