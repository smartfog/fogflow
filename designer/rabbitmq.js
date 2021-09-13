var amqp = require('amqplib/callback_api');
//const amqplib = require('amqplib');
//import {amqp} from 'amqplib/callback_api';

RABBIT_URL = 'amqp://localhost'
EXCHANGE = 'Op'
USER = 'admin'
PASSWORD = 'mypass'

const opt = { credentials: require('amqplib').credentials.plain(USER, PASSWORD) };
const hOpt = {content_type:'application/json'}

msg1 ={
  "entityId": { "id": "KK", "type": "Operator", "isPattern": false },
  "attributes": [
    {
      "name": "designboard",
      "type": "object",
      "value": {"edges":[],"blocks":[{"id":1,"x":-316.85713999999996,"y":-96.71428721679689,"type":"Operator","module":"null","values":{"name":"amqp_test1","description":""}}]}
    },
    {
      "name": "operator",
      "type": 'object',
      "value": {"name":"KK","description":"","parameters":[]}
    }
  ],
  "domainMetadata": [ { "name": "location", "type": "global", "value": "global" } ],
  "dgraph.type": "ContextElement"
}
//send ={'Type': 'Operator', 'RoutingKey': 'Operator.', 'From': 'designer', 'PayLoad': msg1}

function string2json(data){
  console.log("------- stirng 2 json",data);
  for(var i=0;i<data[0].attributes.length;i++){
         console.log("******** data ",data[0].attributes[i].value)
         data[0].attributes[i].value = JSON.parse(data[0].attributes[i].value);
  }
  return data
 
  
}

async function amqpPubTest(msg, exchange_) {

  //payData = string2json(msg.contextElements)
  console.log("inside in amqp ",msg.contextElements[0].attributes)
  send ={'Type': 'Operator', 'RoutingKey': 'Operator.', 'From': 'designer', 'PayLoad': (msg.contextElements[0])}
  amqp.connect(RABBIT_URL, opt,function (error0, connection) {
    if (error0) {
      throw error0;
    }
    connection.createChannel(function (error1, channel) {
      if (error1) {
        throw error1;
      }
      var exchange = exchange_ || EXCHANGE;
      channel.assertQueue(queue='Operator', {durable: true, autoDelete: true});

      channel.assertExchange(exchange, 'topic', {
        durable: true, autoDelete: true
      });
      channel.publish(exchange, 'Operator.', Buffer.from(JSON.stringify(send)),hOpt);
      console.log(" [x] Sent %s", send);
    });

    setTimeout(function () {
      connection.close();
      //process.exit(0);
    }, 500);
  });


}


if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
  this.axios = require('axios')
  module.exports.amqpPubTest = amqpPubTest;
} else {
  window.amqpPubTest = amqpPubTest;
}

amqpPubTest("test")

// var amqp = require('amqplib/callback_api');
// process.env.CLOUDAMQP_URL = 'amqp://localhost';

// var exchange = 'logs';
// var pubChannel = null;

// // if the connection is closed or fails to be established at all, we will reconnect
// var amqpConn = null;
// async function start() {
//   amqp.connect(process.env.CLOUDAMQP_URL + "?heartbeat=60", function(err, conn) {
//     if (err) {
//       console.error("[AMQP]", err.message);
//       return setTimeout(start, 7000);
//     }
//     conn.on("error", function(err) {
//       if (err.message !== "Connection closing") {
//         console.error("[AMQP] conn error", err.message);
//       }
//     });
//     conn.on("close", function() {
//       console.error("[AMQP] reconnecting");
//       return setTimeout(start, 7000);
//     });

//     console.log("[AMQP] connected");
//     amqpConn = conn;
//     pubChannel = conn;
//     //whenConnected();
//   });
// }

// function whenConnected() {
//   startPublisher();
// }

// var offlinePubQueue = [];
// function startPublisher() {
//   amqpConn.createConfirmChannel(function(err, ch) {
//     if (closeOnErr(err)) return;
//     ch.on("error", function(err) {
//       console.error("[AMQP] channel error", err.message);
//     });
//     ch.on("close", function() {
//       console.log("[AMQP] channel closed");
//     });

//     pubChannel = ch;
//     while (true) {
//       var m = offlinePubQueue.shift();
// 			console.log('M = ', m);
//       if (!m) break;
//       publish(m[0], m[1], m[2]);
//     }
//   });
// }

// // method to publish a message, will queue messages internally if the connection is down and resend later
// function publish(exchange, routingKey, content) {
//   try {
//     start()
//     console.log("---------------")
//     pubChannel.publish(exchange, routingKey, content, { persistent: true },
//                        function(err, ok) {
//                          if (err) {
//                            console.error("[AMQP] publish", err);
//                            offlinePubQueue.push([exchange, routingKey, content]);
//                            pubChannel.connection.close();
//                          }
//                        });
//   } catch (e) {
//     console.error("[AMQP] publish", e.message);
//     offlinePubQueue.push([exchange, routingKey, content]);
//   }
// }

// function closeOnErr(err) {
//   if (!err) return false;
//   console.error("[AMQP] error", err);
//   amqpConn.close();
//   return true;
// }

// setTimeout(function() {
//   publish("", "jobs", new Buffer("work work work"));
// }, 3000);

// publish(exchange,'',Buffer.from("hello world"));


