// simple algorithm to detect anomaly events
function anomalyDetection(msg) {
    var watts = msg.usage.value;
    var deviceID = msg.id;
    //var shopID = msg.metadata.shop.value;
    //var location = msg.metadata.location;

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
    threshold = threshold.value;
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
    console.log("enter into the user-defined fog function",contextEntity);
    if (contextEntity == null) {
        return;
    }
    
    var type = contextEntity.type;
    console.log('type ', type);
    console.log("contextEntity",contextEntity);

    // do the internal data processing
    if (type == "PowerPanelNew") {
	console.log("Enter into powerpanel")
	console.log("contextEntity.usage",contextEntity.usage)
	
        var watts = contextEntity.usage.value;
        var deviceID = contextEntity.id;
        var shopID = contextEntity.shop.value;
        //var location = contextEntity.location;
	console.log('==========watts',watts, ', ======deviceID=====',deviceID)

        if (watts > threshold) { // detect an anomaly event
            // publish the detected event
            var anomaly = {};

            var now = new Date();
            anomaly['when'] = now.toISOString();
            anomaly['whichpanel'] = deviceID;
            anomaly['whichshop'] = shopID;
            //anomaly['where'] = location;
            anomaly['usage'] = watts;

            var updateEntity = {};


            updateEntity.when = {
                type: 'Property',
                value: anomaly['when']
            };
            updateEntity.whichpanel = {
                type: 'Property',
                value: anomaly['whichpanel']
            };

            updateEntity.shop = {
                type: 'Property',
                value: anomaly['whichshop']
            };
            /*updateEntity.where = {
                type: 'Property',
                value: anomaly['where']
            };*/
            updateEntity.usage = {
                type: 'Property',
                value: anomaly['usage']
            };

            /*updateEntity.metadata = {};
            updateEntity.metadata.shop = {
                type: 'string',
                value: anomaly['whichshop']
            };*/
	    console.log("publish")
            console.log("updateEntity",updateEntity)
            publish(updateEntity)
        }
    } else if (type == 'Rule') {
	console.log(contextEntity.threshold.value)
        threshold = contextEntity.threshold.value;
        console.log('update the threshold to ', threshold);
    }
};
