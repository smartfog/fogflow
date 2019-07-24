exports.handler = function(contextEntity, publish, query, subscribe) {
    console.log("entered into the user-defined fog function");

    var entityID = contextEntity.entityId.id;

    if (contextEntity == null) {
        return;
    }
    if (contextEntity.attributes == null) {
        return;
    }

// For broker-config.json

    var BrokerIP = contextEntity.attributes.brokerIP.value;
    BrokerIP = BrokerIP.toString();
    console.log('BrokerIP:  ',BrokerIP);

    var BrokerPort = contextEntity.attributes.brokerPort.value;
    BrokerPort = BrokerPort.toString();
    console.log('BrokerPort:  ',BrokerPort);

    var command = "./gpadapter-config.sh " + BrokerIP + " " + BrokerPort;
    console.log(command);
    const shell = require('shelljs');
    shell.exec(command);

};
