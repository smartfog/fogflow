var amqp = require('amqplib/callback_api');
//const amqplib = require('amqplib');
//import {amqp} from 'amqplib/callback_api';

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
//send ={'Type': 'Operator', 'RoutingKey': 'Operator.', 'From': 'designer', 'PayLoad': msg1}

function string2json(data) {
  console.log("------- stirng 2 json", data);
  for (var i = 0; i < data[0].attributes.length; i++) {
    console.log("******** data ", data[0].attributes[i].value)
    data[0].attributes[i].value = JSON.parse(data[0].attributes[i].value);
  }
  return data


}

function amqpConnection(exchange_) {
  amqp.connect(RABBIT_URL, opt, function (error0, connection) {
    if (error0) {
      throw error0;
    }
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
    
  });
}



async function amqpPub(data, exchange_) {
  if (channel_ === undefined) amqpConnection();
  let msg = JSON.parse(JSON.stringify(data));
  //console.log("inside in amqp ",msg.attribute)
  //payData = string2json(msg.contextElements)
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
    // reconnect
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
