var amqp = require('amqplib/callback_api');


RABBIT_URL = 'amqp://localhost'
EXCHANGE = 'Op'
USER = 'admin'
PASSWORD = 'mypass'
var exchange = undefined;
var channel_ = undefined;
const opt = { credentials: require('amqplib').credentials.plain(USER, PASSWORD) };
const hOpt = { content_type: 'application/json' }

/** Hardcode value for testing */
msg1 = {
  "entityId": { "id": "KK", "type": "Operator", "isPattern": false },
  "attributes": [
    {
      "name": "designboard",
      "type": "object",
      "value": { "edges": [], "blocks": [{ "id": 1, "x": -316.85713999999996, "y": -96.71428721679689, "type": "Operator", "module": "null", "values": { "name": "amqp_test1", "description": "" } }] }
    },
    {
      "name": "operator",
      "type": 'object',
      "value": { "name": "KK", "description": "", "parameters": [] }
    }
  ],
  "domainMetadata": [{ "name": "location", "type": "global", "value": "global" }],
  "dgraph.type": "ContextElement"
}

function string2json(data) {
  for (var i = 0; i < data[0].attributes.length; i++) {
    console.log("******** data ", data[0].attributes[i].value)
    data[0].attributes[i].value = JSON.parse(data[0].attributes[i].value);
  }
  return data


}

function reconnect() {
  amqpConnection()
}

var refreshIntervalId = undefined;
var isAmqpUp = false;

function amqpConnection(exchange_) {
  amqp.connect(RABBIT_URL, opt, function (error0, connection) {
    if (error0) {
      console.log("Connection Error : Retrying to connect RabbitMQ in 10 seconds ");
      if (refreshIntervalId === undefined){
        refreshIntervalId = setInterval(reconnect, 10000);
      }
      //throw error0;
    }
   if(connection){
    console.log("connected with RabbitMQ");
    clearInterval(refreshIntervalId);
    refreshIntervalId = undefined;
    global.isAmqpUp = true;

    connection.createChannel(function (error1, channel) {
      if (error1) {
        throw error1;
      }
      exchange = exchange_ || EXCHANGE;
      channel_ = channel;
      channel.assertQueue(queue = 'Operator', { durable: true, autoDelete: true });

      channel.assertExchange(exchange, 'topic', {
        durable: true, autoDelete: true
      });
     
    });
   } 
  });
 
}



async function amqpPub(data, exchange_) {
  if (channel_ === undefined) amqpConnection();
  let msg = JSON.parse(JSON.stringify(data));
  if (msg.attribute == undefined) {
    return
  }
  if (msg.attribute.hasOwnProperty("designboard")) delete msg.attribute.designboard;
  if (msg.attribute.hasOwnProperty("uid")) delete msg.attribute.uid;
  console.log("final amqp msg ** ", msg.attribute)
  send = { 'Type': msg.internalType, 'RoutingKey': 'Operator.', 'From': 'designer', 'PayLoad': msg.attribute }
  try{
    channel_.publish(exchange, 'Operator.', Buffer.from(JSON.stringify(send)), hOpt);
  }catch(err){
    amqpConnection();
  }
  console.log(" [x] Sent %s", send.PayLoad);
  //return
}


if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
  this.axios = require('axios')
  module.exports.amqpPub = amqpPub;
  module.exports.amqpConnection =amqpConnection;
} else {
  window.amqpPub = amqpPub;
  window.amqpConnection = amqpConnection;
}
