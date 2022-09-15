const amqplib = require('amqplib');
const NGSIClient = require('./ngsiclient');

// list of active tasks launched by FogFlow workers
var taskMap = {};

//const fogflow_host_ip = 'localhost';
const fogflow_host_ip = '10.11.11.166';
const fogflow_port = '8080';

//const locations = [{
//    latitude: 38.01506335734526,
//    longitude: -1.1715629779027177
//}];


const locations = [{
    latitude: 38.018837048090326,
    longitude: -1.1693835099259555
},{
    latitude: 38.01506335734526,
    longitude: -1.1715629779027177
},{
    latitude: 38.016654446651685,
    longitude: -1.1680698394775393
}];

const rabbitmq_ip = fogflow_host_ip;
const rabbitmq_port = 5672;
const rabbitmq_user = 'admin';
const rabbitmq_password = 'mypass';

const rabbitmq_url = 'amqp://' + rabbitmq_user + ':' + rabbitmq_password + '@'
    + rabbitmq_ip + ':' + rabbitmq_port.toString();

console.log(rabbitmq_url);

var ngsi10client = new NGSIClient.NGSI10Client("http://" + fogflow_host_ip + ":" + fogflow_port +"/ngsi10");


const myArgs = process.argv.slice(2);
console.log('myArgs: ', myArgs);
var num_input_entities = parseInt(myArgs[0]);

// for rabbitmq client
var amqp_url = null;
var amqpConn = null;
var amqpChannel = null;
var msgHandler = null;

const exchange_name = 'fogflow';
const exchange_type = 'topic';
const queue_name = 'fogflow-test-client';
const subscribed_keys = ['designer.*', 'task.'];

const TIME_INTERVAL_RECONNECT = 5000;

var cb_after_connected = null;

var start_time;
var stop_time;
var migrate_time;

var test_status = "INIT";

Init(rabbitmq_url, handleInternalMessage, startTest);

process.on('SIGINT', function() {
    if (test_status == "ENABLED") {
        migrate_time = Date.now();       
        update_temperature_entities(num_input_entities);         
        test_status = "MIGRATE";        
    }else if (test_status == "MIGRATE") {        
        stop_time = Date.now();    
        stopTest();        
        test_status = "DISABLED";                
    }else if (test_status == "DISABLED") {
        printTestResult();      
        test_status = "DELETE";          
    }else if (test_status == "DELETE") {
        delete_temperature_entities(num_input_entities);        
        test_status = "EXIT";        
    }else if (test_status == "EXIT") {
        process.exit();
    }
});

// this is only triggered when starting the task designer
async function startTest() {
    console.log("[RabbitMQ] is already connected")

    await create_temperature_entities(num_input_entities);
    console.log("create the input entities...");    

    await enable_test_fog_function();
    console.log("enable the fog function...");     
    start_time = Date.now();

    test_status = "ENABLED";
}


function maxMinAvg(arr) {
    var max = arr[0];
    var min = arr[0];
    var sum = arr[0]; //changed from original post
    var num = 1;
    for (var i = 1; i < arr.length; i++) {
        if (arr[i] == null) {
            continue;
        }
        
        if (arr[i] > max) {
            max = arr[i];
        }
        if (arr[i] < min) {
            min = arr[i];
        }
        sum = sum + arr[i];
        num = num + 1
    }
    
    return [max, min, sum/num]; //changed from original post
}

function printTestResult() {
    console.log("task_id", ",", "assignment_time", ",", "orchestration_time", ",", "migration_time", ",", "termination_time");

    var array_assignment_time = [];
    var array_orchestration_time = [];
    var array_migration_time = [];    
    var array_termination_time = [];

    for(var tid in taskMap) {
        var task = taskMap[tid];
        console.log(tid, ",", task["assignment_time"], ",", task["orchestration_time"], ",", task["migration_time"], ",", task["termination_time"]);
        array_assignment_time.push(task["assignment_time"]);
        array_orchestration_time.push(task["orchestration_time"]);
        array_migration_time.push(task["migration_time"]);        
        array_termination_time.push(task["termination_time"]);
    }

    console.log("assignment_time: ", maxMinAvg(array_assignment_time));
    console.log("orchestration_time: ", maxMinAvg(array_orchestration_time));
    console.log("migration_time: ", maxMinAvg(array_migration_time));    
    console.log("termination_time: ", maxMinAvg(array_termination_time));
}

async function stopTest() {
    await disable_test_fog_function();
    console.log("disable the fog function...");     
}

async function create_temperature_entities(num) {
    for(var i=1; i<=num; i++) {
        var deviceObj = {};

        deviceObj.entityId = {
            id: "Device." + i,
            type: "Temperature",
            isPattern: false		
        };

        deviceObj.attributes = {}
        deviceObj.attributes["protocol"] = { type:"string", value: "NGSI-v1" }; 
    
        deviceObj.metadata = {}

        num_locations = locations.length;
        var idx = i%num_locations
        mylocation = locations[idx];

        deviceObj.metadata.location = {
            type: 'point',
            value: { 'latitude': mylocation.latitude, 'longitude': mylocation.longitude }
        };
    
        // create the device entity in FogFlow Cloud Broker
        ngsi10client.updateContext(deviceObj).then(function (data) {
            console.log(data);
        }).catch(function (error) {
            console.log('failed to register the new device object');
        });         
    }
}

async function update_temperature_entities(num) {
    console.log("update the location of all device entities");            
    
    for(var i=1; i<=num; i++) {
        var deviceObj = {};

        deviceObj.entityId = {
            id: "Device." + i,
            type: "Temperature",
            isPattern: false		
        };

        deviceObj.attributes = {}
        deviceObj.attributes["protocol"] = { type:"string", value: "NGSI-v1" }; 
    
        deviceObj.metadata = {}

        num_locations = locations.length;
        var idx = (i+1)%num_locations
        mylocation = locations[idx];

        deviceObj.metadata.location = {
            type: 'point',
            value: { 'latitude': mylocation.latitude, 'longitude': mylocation.longitude }
        };
    
        // create the device entity in FogFlow Cloud Broker
        ngsi10client.updateContext(deviceObj).then(function (data) {
            console.log(data);
        }).catch(function (error) {
            console.log('failed to update the device object');
        });         
    }
}

async function delete_temperature_entities(num) {
    for(var i=1; i<=num; i++) {

        entityId = {
            id: "Device." + i,
            type: "Temperature",
            isPattern: false		
        };
    
        console.log("delete the entity: ", "Device." + i);
    
        // delete the device entity in FogFlow Cloud Broker
        await ngsi10client.deleteContext(entityId).then(function (data) {
            console.log(data);
        }).catch(function (error) {
            console.log('failed to delete the entiti');
        });         
    }
}

async function enable_test_fog_function() {
    var url = "http://" + fogflow_host_ip + ":8080/fogfunction/mytest/enable";
    await fetch(url);    
}

async function disable_test_fog_function() {
    var url = "http://" + fogflow_host_ip + ":8080/fogfunction/mytest/disable";
    await fetch(url);    
}


function handleInternalMessage(jsonMsg) {
    var msgType = jsonMsg.Type;
    switch (msgType) {
        case 'MASTER_JOIN':

            break;
        case 'MASTER_LEAVE':

            break;
        case 'TASK_UPDATE':
            onTaskUpdate(jsonMsg)
            break;
        case 'TASK_INFO':
            onTaskInfo(jsonMsg)
            break;            
    }
}


function onTaskUpdate(msg) {
    console.log("========task_update=========");

    var payload = msg.PayLoad;
    var workerID = msg.From;

    if (payload.TopologyName != 'mytest') {
        console.log("non-relevant task");
        return
    }

    var taskID = payload.TaskID
    if ( !(taskID in taskMap)) {
        taskMap[taskID] = {}
    }   

    if (payload.Status == 'running') {
        if (test_status == "MIGRATE") {
            taskMap[taskID].migration_time = Date.now() - migrate_time;                        
        } else {        
            taskMap[taskID].assignment_time = Date.now() - start_time;            
        }
    } else if (payload.Status == 'removed') {
        taskMap[taskID].termination_time = Date.now() - stop_time;
    }

    console.log(payload);
}

function onTaskInfo(msg) {
    //console.log("========task_info=========");
  
    var payload = msg.PayLoad;
    var workerID = msg.From;
    var taskID = payload.TaskID;

    if (payload.TopologyName != 'mytest') {
        console.log("non-relevant task");
        return
    }
    
    if (taskID in taskMap) {
        taskMap[taskID].orchestration_time = Date.now() - start_time;
        console.log(taskID, ": started");        
    } else {
        console.log("non-relevant task");
    }

    //console.log(payload);
}


function updateTaskList(updateMsg, fromWorker) {
    console.log("========task_update=========");

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


function Init(rabbitmqURL, fnConsumer, afterConnected) 
{
    console.log("[RabbitMQ] connecting to ", rabbitmqURL);    
    amqp_url = rabbitmqURL;
    msgHandler = fnConsumer
	cb_after_connected = afterConnected;
    
    amqplib.connect(amqp_url).then(function(conn) {       
        console.log("[RabbitMQ] connected");
        amqpConn = conn;
        
        whenConnected();
    }).catch( function(err) {
        console.error("[RabbitMQ]", err.message);
        return setTimeout(reConnect, TIME_INTERVAL_RECONNECT);
    });
}

function reConnect() {
    console.log("[RabbitMQ] reconnecting to ", amqp_url);    
	
    amqplib.connect(amqp_url).then(function(conn) {       
        console.log("[RabbitMQ] connected");
        amqpConn = conn;
        
        whenConnected();
    }).catch( function(err) {
        console.error("[RabbitMQ]", err.message);
        return setTimeout(reConnect, TIME_INTERVAL_RECONNECT);
    });
}

async function whenConnected() {
    amqpChannel = await amqpConn.createChannel()

    //create the exchange 
    await amqpChannel.assertExchange(exchange_name, exchange_type, {durable: true, autoDelete: true}).catch(console.error);       
    
    //start the consumer
    await amqpChannel.assertQueue(queue_name, {durable: true});
    
    for(var i=0; i<subscribed_keys.length; i++){
        var key = subscribed_keys[i];
        console.log("[RabbitMQ] subscribed to ", key);
        await amqpChannel.bindQueue(queue_name, exchange_name, key);           
    }
    
    await amqpChannel.consume(queue_name, processMsg, { noAck: true });        
	
	cb_after_connected();
}

function processMsg(msg) {
    var jsonMsg = JSON.parse(msg.content)
    msgHandler(jsonMsg);
}

async function Publish(msg){
    var msgContent = JSON.stringify(msg);
    await amqpChannel.publish(exchange_name, msg.RoutingKey, Buffer.from(msgContent));
}

