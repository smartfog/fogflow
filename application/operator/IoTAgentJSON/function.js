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
    
//	console.log('ContextEntity.......',contextEntity);

	var MongoIP = contextEntity.attributes.mongoIP.value;
	MongoIP = MongoIP.toString();
	console.log('mongoIP:  ',MongoIP);

	var MongoPort = contextEntity.attributes.mongoPort.value;
	MongoPort = MongoPort.toString();
	console.log('mongoPort:  ',MongoPort);

	var BrokerIP = contextEntity.attributes.brokerIP.value;
        BrokerIP = BrokerIP.toString();
	console.log('brokerIP:  ',BrokerIP);

	var BrokerPort = contextEntity.attributes.brokerPort.value;
	BrokerPort = BrokerPort.toString();
	console.log('brokerPort:  ',BrokerPort);

	command = "/opt/iotajson/iota-config.sh " + MongoIP  + " " + MongoPort + " " + BrokerIP + " " + BrokerPort;
	
        const shell = require('shelljs');
	shell.exec(command);
};

