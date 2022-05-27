const mqtt = require('mqtt');

var topicMap = {};


function prepare_entity_update(msg, attribute_mappings, device_id, device_location) {
  var contextUpdates = [];

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
    
      updateEntity.attributes["refDevice"] = { type: "relationship", value: device_id };
      var now = new Date();
      updateEntity.attributes["dateObserved"] = { type: "datetime", value: now.toISOString() };

      switch(mapping.type.toLowerCase()) {
        case "number":
          updateEntity.attributes[mapping.name] = { type: "float", value: parseFloat(msg[attr]) };
          break;        
        case "float":
          updateEntity.attributes[mapping.name] = { type: "float", value: parseFloat(msg[attr]) };
          break;
        case "integer":
          updateEntity.attributes[mapping.name] = { type: "float", value: parseInt(msg[attr]) };
          break;        
        case "string":
          updateEntity.attributes[mapping.name] = { type: "string", value: msg[attr] };
          break;
        default:
          console.log("mapping_type ", mapping.type, " is not supported yet");
      }
    
      updateEntity.metadata = {};
      updateEntity.metadata.location = device_location;

      contextUpdates.push(updateEntity);
    } else {
      console.log("there is no defined mapping for this attribute ", attr);
    }
  }

  return contextUpdates;
}

exports.handler = function (contextEntity, publish, query, subscribe) {
  console.log("enter into the user-defined fog function");

  if (contextEntity == null) {
    return;
  }
  if (contextEntity.attributes == null) {
    return;
  }

  console.log(contextEntity);

  var deviceID = contextEntity.entityId.id;
  var deviceLocation = contextEntity.metadata.location;

  var mqttbrokerURL = contextEntity.attributes.mqttbroker.value;
  var client = mqtt.connect(mqttbrokerURL);

  var topic = contextEntity.attributes.topic.value;

  var attribute_mappings = contextEntity.attributes.mappings.value;
  console.log(attribute_mappings);

  // register the link from topic to attribute-mappings
  if (topicMap.hasOwnProperty(topic) == false) {
    topicMap[topic] = attribute_mappings;
  }

  client.on('connect', () => {
    console.log('Connected')
    client.subscribe([topic], () => {
      console.log("Subscribe to topic ", topic);
    })
  })

  client.on('message', (topic, payload) => {
    // get the corresponding attribute_mappings for this topic
    var mappings = topicMap[topic];
    
    try {
        var msg = JSON.parse(payload.toString());        
    } catch (error) {
        console.log('Error happened here!');
        console.error(error);
        return;        
    }

    // create an entity update with the transformed attributes
    update_messages = prepare_entity_update(msg, mappings, deviceID, deviceLocation);
    
    // publish the created entity update
    for(var i=0; i<update_messages.length; i++) {
      publish(update_messages[i]);
      console.log("publish: ", update_messages[i]);
    }
  })

};

