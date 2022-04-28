const mqtt = require('mqtt');

exports.handler = function (contextEntity, publish, query, subscribe) {
  console.log("enter into the user-defined fog function");

  if (contextEntity == null) {
    return;
  }
  if (contextEntity.attributes == null) {
    return;
  }

  console.log(contextEntity);

  var device_id = contextEntity.entityId.id;

  var mqttbrokerURL = contextEntity.attributes.mqttbroker.value;
  var client = mqtt.connect(mqttbrokerURL);

  var topic = contextEntity.attributes.topic.value;

  var attribute_mappings = contextEntity.attributes.mappings.value;
  console.log(attribute_mappings);

  client.on('connect', () => {
    console.log('Connected')
    client.subscribe([topic], () => {
      console.log("Subscribe to topic ", topic);
    })
  })

  client.on('message', (topic, payload) => {
    var msg = JSON.parse(payload.toString());
    for (const attr in msg) {
      console.log(attr, ":", msg[attr]);

      if (attr in attribute_mappings) {
        var mapping = attribute_mappings[attr];

        var updateEntity = {};
        updateEntity.entityId = {
          id: mapping.entity_type + ":" + device_id,
          type: mapping.entity_type,
          isPattern: false
        };
        updateEntity.attributes = {};

        updateEntity.attributes["device"] = { type: "string", value: device_id };
        var now = new Date();
        updateEntity.attributes["dateObserved"] = { type: "DateTime", value: now.toISOString() };

        if (mapping.type == "Number") {
          updateEntity.attributes[mapping.name] = { type: "float", value: parseFloat(msg[attr]) };
        } else if (mapping.type == "String") {
          updateEntity.attributes[mapping.name] = { type: "string", value: msg[attr] };
        }

        publish(updateEntity);
        console.log("publish: ", updateEntity);
      }
    }
  })

};

