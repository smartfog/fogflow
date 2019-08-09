exports.handler = function(contextEntity, publish, query, subscribe)
{
	console.log("entered into the user-defined fog function");

	if (contextEntity == null) {
		return;
	} 	
	if (contextEntity.attributes == null) {
		return;
	}

    // ============================== processing ======================================================
    // processing the received ContextEntity:    
    

	var BrokerIP = contextEntity.attributes.brokerIP.value;
        BrokerIP = BrokerIP.toString();
	console.log('brokerIP:  ',BrokerIP);

	var BrokerPort = contextEntity.attributes.brokerPort.value;
	BrokerPort = BrokerPort.toString();
	console.log('brokerPort:  ',BrokerPort);

	command = "/opt/iotajson/iota-config.sh " + BrokerIP + " " + BrokerPort;
	
        const shell = require('shelljs');
	shell.exec(command);
};

