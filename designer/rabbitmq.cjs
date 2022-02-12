const amqplib = require('amqplib');

var amqp_url = null;
var amqpConn = null;
var amqpChannel = null;
var msgHandler = null;

const exchange_name = 'fogflow';
const exchange_type = 'topic';
const queue_name = 'fogflow-designer';
const subscribed_keys = ['designer.*'];

const TIME_INTERVAL_RECONNECT = 5000;

function Init(rabbitmqURL, fnConsumer) 
{    
    amqp_url = rabbitmqURL;
    msgHandler = fnConsumer
    
    amqplib.connect(amqp_url).then(function(conn) {       
        console.log("[RabbitMQ] connected");
        amqpConn = conn;
        
        whenConnected();
    }).catch( function(err) {
        console.error("[RabbitMQ]", err.message);
        return setTimeout(Init, TIME_INTERVAL_RECONNECT);
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
}

function processMsg(msg) {
    var jsonMsg = JSON.parse(msg.content)
    msgHandler(jsonMsg);
}

async function Publish(msg){
    var msgContent = JSON.stringify(msg);
    await amqpChannel.publish(exchange_name, msg.RoutingKey, Buffer.from(msgContent));
    console.log("[RabbitMQ] published ", msg);
}

module.exports = { Init, Publish }
