// simple algorithm to detect anomaly events
function anomalyDetection(msg) {
    var watts = msg.attributes.usage.value;
    var deviceID = msg.entityId.id;
    var shopID = msg.metadata.shop.value;
    var location = msg.metadata.location;

    console.log('============+++++++++++Current usage ', watts, ', current thrshold = ', threshold)

    if (watts > threshold) { // detect an anomaly event
        // publish the detected event
        var anomaly = {};

        var now = new Date();
        anomaly['when'] = now.toISOString();
        anomaly['whichpanel'] = deviceID;
        anomaly['whichshop'] = shopID;
        anomaly['where'] = location;
        anomaly['usage'] = watts;

        updateContext(anomaly)
    }
}


function updateRule(ruleObj) {
    threshold = ruleObj.attributes.threshold.value;
    console.log('update the threshold to ', threshold);
}


var threshold = 30;


//
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = function(contextEntity, publish, query, subscribe) {
    console.log("enter into the user-defined fog function");

    if (contextEntity == null) {
        return;
    }
    if (contextEntity.attributes == null) {
        return;
    }

    var type = contextEntity.entityId.type;
    console.log('type ', type);
    console.log(contextEntity);

    // do the internal data processing
    if (type == 'PowerPanel') {
        var watts = contextEntity.attributes.usage.value;
        var deviceID = contextEntity.entityId.id;
        var shopID = contextEntity.metadata.shop.value;
        var location = contextEntity.metadata.location;

        console.log('============+++++++++++Current usage ', watts, ', current thrshold = ', threshold)

        if (watts > threshold) { // detect an anomaly event
            // publish the detected event
            var anomaly = {};

            var now = new Date();
            anomaly['when'] = now.toISOString();
            anomaly['whichpanel'] = deviceID;
            anomaly['whichshop'] = shopID;
            anomaly['where'] = location;
            anomaly['usage'] = watts;

            var updateEntity = {};

            updateEntity.attributes = {};

            updateEntity.attributes.when = {
                type: 'string',
                value: anomaly['when']
            };
            updateEntity.attributes.whichpanel = {
                type: 'string',
                value: anomaly['whichpanel']
            };

            updateEntity.attributes.shop = {
                type: 'string',
                value: anomaly['whichshop']
            };
            updateEntity.attributes.where = {
                type: 'object',
                value: anomaly['where']
            };
            updateEntity.attributes.usage = {
                type: 'integer',
                value: anomaly['usage']
            };

            updateEntity.metadata = {};
            updateEntity.metadata.shop = {
                type: 'string',
                value: anomaly['whichshop']
            };

            publish(updateEntity)
        }
    } else if (type == 'Rule') {
        threshold = contextEntity.attributes.threshold.value;
        console.log('update the threshold to ', threshold);
    }
};