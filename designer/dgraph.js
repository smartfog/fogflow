const dgraph = require("dgraph-js");
const grpc = require("grpc");
var request = require('request');
var config_fs_name = './config.json';
//var axios = require('axios')
var fs = require('fs');
const bodyParser = require('body-parser');
var globalConfigFile = require(config_fs_name)
var config = globalConfigFile.designer;
config.grpcPort = globalConfigFile.persistent_storage.port;
config.HostIp = globalConfigFile.external_hostip;
config.brokerIp=globalConfigFile.coreservice_ip
config.brokerPort=globalConfigFile.broker.http_port
console.log(config.grpcPort)
console.log(config.HostIp)

function newClientStub() {
    return new dgraph.DgraphClientStub(config.HostIp+":"+config.grpcPort, grpc.credentials.createInsecure());
}

// Create a client.
function newClient(clientStub) {
    return new dgraph.DgraphClient(clientStub);
}

// Drop All - discard all data and start from a clean slate.
async function dropAll(dgraphClient) {
    const op = new dgraph.Operation();
    op.setDropAll(true);
    await dgraphClient.alter(op);
}

/*
   create scheema for node
*/
async function setSchema(dgraphClient) {
    const schema = `
            attributes: [uid] .
            domainMetadata: [uid] .
            entityId: uid .
            updateAction: string .
            id: string .
            isPattern: bool .
            latitude: float .
            longitude: float .
            name: string .
            type: string . 
    `;
    const op = new dgraph.Operation();
    op.setSchema(schema);
    await dgraphClient.alter(op);
}

/*
   convert object domainmetadata data into string to store entity as single node  
*/

async function resolveDomainMetaData(data) {
     if ('domainMetadata' in data) {
     	var len=data.domainMetadata.length
     	for(var i=0;i < len; i++) {
		if('value' in data.domainMetadata[i])
                 data.domainMetadata[i].value=JSON.stringify(data.domainMetadata[i].value)
                if('location' in data.domainMetadata[i])
	        data.domainMetadata[i].location=JSON.stringify(data.domainMetadata[i].location)
	  }
      }
}

/*
   convert object attributes data into string to store entity as single node 
*/

async function resolveAttributes(data) {
     if ('attributes' in data){
     	var length=data.attributes.length
     	console.log(length)
     	for(var i=0;i < length; i++) {
        	if('value' in data.attributes[i]) {           
        		if (data.attributes[i].type=='object')
			data.attributes[i].value=JSON.stringify(data.attributes[i].value)
         	}
       	  }
    }
}

/*
   insert data into database
*/

async function createData(dgraphClient,ctx) {
    const txn = dgraphClient.newTxn();
    try {
        const mu = new dgraph.Mutation();
        mu.setSetJson(ctx);
	console.log(mu)
        const response = await txn.mutate(mu);
	console.log(response)
	console.log(txn)
        await txn.commit();
    }
	 finally {
        await txn.discard();
    }
}

/*
   send data to cloud broker
*/

async function sendData(contextElement) {
	var  updateCtxReq = {};
	updateCtxReq.contextElements = [];
	updateCtxReq.updateAction = 'UPDATE'	
	updateCtxReq.contextElements.push(contextElement)
	request({
            method: 'post',
            url:'http://'+config.brokerIp+':'+config.brokerPort+'/ngsi10/updateContext',
            body:JSON.stringify(updateCtxReq)
        })
	console.log(updateCtxReq)

}

/*
   convert string object into structure to register data into cloud broker
*/

async function changeFromDataToobject(contextElement) {
    contextEle=contextElement['contextElements']
    for (var ctxEle=0; ctxEle < contextEle.length; ctxEle=ctxEle+1) {
	for(var ctxAttr=0; ctxAttr<contextEle[ctxEle].attributes.length ;ctxAttr=ctxAttr+1) {
		if (contextEle[ctxEle].attributes[ctxAttr].type=='object') {
			const value=contextEle[ctxEle].attributes[ctxAttr].value
			contextEle[ctxEle].attributes[ctxAttr].value=JSON. parse(value)
	       	}
        }
	for(ctxdomain=0; ctxdomain<contextEle[ctxEle].domainMetadata.length; ctxdomain=ctxdomain+1) {
		if ('value' in contextEle[ctxEle].domainMetadata[ctxdomain]) {
			const value=contextEle[ctxEle].domainMetadata[ctxdomain].value
			contextEle[ctxEle].domainMetadata[ctxdomain].value=JSON.parse(value)
		}
		if ('location' in contextEle[ctxEle].domainMetadata[ctxdomain]) {
			const loc=contextEle[ctxEle].domainMetadata[ctxdomain].location
			contextEle[ctxEle].domainMetadata[ctxdomain].location=JSON.parse(loc)
		}
	}
     sendData(contextEle[ctxEle])
    }
}

async function sendPostRequestToBroker(contextElement) {
    await changeFromDataToobject(contextElement)	
}

/*
   Query for getting the registered node
*/

async function queryData(dgraphClient) {
   const query = `{

        contextElements(func: has(entityId)) {
           {
                  entityId{
                       id
                       type
                       isPattern
                  }
                  attributes{
                       name
                       type
		       value
                  }
                  domainMetadata{
                       name
                       type
		       value
                       location
                  }
              }
           }

    }`;

    responseBody = await dgraphClient.newTxn().queryWithVars(query);
    const responsData= responseBody.getJson();
    const responseObject=JSON.stringify(responsData)
    sendPostRequestToBroker(responsData)
}

/*
   main handler 
*/

async function db(contextData) {
    const dgraphClientStub = newClientStub();
    const dgraphClient = newClient(dgraphClientStub);
    await dropAll(dgraphClient);
    if ('contextElements' in contextData) {	
	contextData=contextData['contextElements']
	contextData=contextData[0]
	console.log(contextData)
    }
    console.log(contextData)
    await resolveAttributes(contextData)
    console.log(contextData)
    await resolveDomainMetaData(contextData)
    await setSchema(dgraphClient);
    await createData(dgraphClient,contextData);
    await queryData(dgraphClient);
    dgraphClientStub.close();
}
process.on('unhandledRejection', (reason, promise) => {
  console.log('Unhandled Rejection at:', promise, 'reason:', reason);
});

module.exports=db;
