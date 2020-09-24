const dgraph = require("dgraph-js");
const grpc = require("grpc");
var request = require('request');
var config_fs_name = './config.json';
var axios = require('axios')
var fs = require('fs');
const bodyParser = require('body-parser');
var globalConfigFile = require(config_fs_name)
var config = globalConfigFile.designer;
config.grpcPort = globalConfigFile.persistent_storage.port;
config.HostIp = globalConfigFile.external_hostip;
config.brokerIp=globalConfigFile.coreservice_ip
config.brokerPort=globalConfigFile.broker.http_port

/*
   creating grpc client for making connection with dgraph
*/

async function newClientStub() {
    return new dgraph.DgraphClientStub(config.HostIp+":"+config.grpcPort, grpc.credentials.createInsecure());
}

// Create a client.

async function newClient(clientStub) {
    return new dgraph.DgraphClient(clientStub);
}

// Drop All - discard all data and start from a clean slate.
/*async function dropAll(dgraphClient) {
    const op = new dgraph.Operation();
    op.setDropAll(true);
    await dgraphClient.alter(op);
}*/

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
	    value: string . 
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
		if('value' in data.domainMetadata[i]) {
			if(data.domainMetadata[i].type != 'global' && data.domainMetadata[i].type != 'stringQuery'){
                 		data.domainMetadata[i].value=JSON.stringify(data.domainMetadata[i].value)
                        }
		}
	  }
      }
}

/*
   convert object attributes data into string to store entity as single node 
*/

async function resolveAttributes(data) {
     if ('attributes' in data){
     	var length=data.attributes.length
     	for(var i=0;i < length; i++) {
        	if('type' in data.attributes[i]) {           
        		if (data.attributes[i].type=='object')
			data.attributes[i].value=JSON.stringify(data.attributes[i].value)
			else {
                        	data.attributes[i].value=data.attributes[i].value.toString()
			}
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
        const response = await txn.mutate(mu);
        await txn.commit();
    }
	 finally {
        await txn.discard();
    }
}

/*
   send data to cloud broker
*/

async function sendData(contextEle) {
	var  updateCtxReq = {};
	updateCtxReq.contextElements = [];
	updateCtxReq.updateAction = 'UPDATE'	
	updateCtxReq.contextElements.push(contextEle)
        await axios({
            method: 'post',
            url: 'http://'+config.brokerIp+':'+config.brokerPort+'/ngsi10/updateContext',
            data: updateCtxReq
            }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
}

/*
   convert string object into structure to register data into cloud broker
*/

async function changeFromDataToobject(contextElement) {
    contextEle=contextElement['contextElements']
    for (var ctxEle=0; ctxEle < contextEle.length; ctxEle=ctxEle+1) {
	ctxEleReq=contextEle[ctxEle]
        if('attributes' in ctxEleReq) {
		for(var ctxAttr=0; ctxAttr<ctxEleReq.attributes.length ;ctxAttr=ctxAttr+1) {
			if (ctxEleReq.attributes[ctxAttr].type=='object') {
				const value=ctxEleReq.attributes[ctxAttr].value
				ctxEleReq.attributes[ctxAttr].value=JSON.parse(value)
	       		}
			if (ctxEleReq.attributes[ctxAttr].type=='integer') {
                        	const value=ctxEleReq.attributes[ctxAttr].value
                        	ctxEleReq.attributes[ctxAttr].value=parseInt(value)
                	}
			if (ctxEleReq.attributes[ctxAttr].type=='float') {
                        	const value=ctxEleReq.attributes[ctxAttr].value
                        	ctxEleReq.attributes[ctxAttr].value=parseFloat(value)
                	}
			if (ctxEleReq.attributes[ctxAttr].type=='boolean') {
                        	const value=ctxEleReq.attributes[ctxAttr].value
				if(value=='false')
                        	ctxEleReq.attributes[ctxAttr].value=false
				else 
				ctxEleReq.attributes[ctxAttr].value=true
                	}

        	}
        }
        if ('domainMetadata' in ctxEleReq){
		for(ctxdomain=0; ctxdomain<ctxEleReq.domainMetadata.length; ctxdomain=ctxdomain+1) {
			if ('value' in ctxEleReq.domainMetadata[ctxdomain]) {
				if(ctxEleReq.domainMetadata[ctxdomain].type!='global'&& ctxEleReq.domainMetadata[ctxdomain].type!='stringQuery'){
					const value=ctxEleReq.domainMetadata[ctxdomain].value
					ctxEleReq.domainMetadata[ctxdomain].value=JSON.parse(value)
				}
			}
		}
          }
    await sendData(ctxEleReq)
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
    //const responseObject=JSON.stringify(responsData)
    const responseObject=responsData
    sendPostRequestToBroker(responsData)
}

/*
   main handler 
*/

async function db(contextData) {
     try{
     const dgraphClientStub = await newClientStub();
     const dgraphClient = await  newClient(dgraphClientStub);
    //await dropAll(dgraphClient);
    if ('contextElements' in contextData) {	
	contextData=contextData['contextElements']
	contextData=contextData[0]
    }
    await resolveAttributes(contextData)
    await resolveDomainMetaData(contextData)
    await createData(dgraphClient,contextData);
    await dgraphClientStub.close();
     }catch(err){
	console.log('DB ERROR::',err);
    }
}

process.on('unhandledRejection', (reason, promise) => {
 //bug-113 fix
 console.log('Retrying.. There is an Unhandled Rejection at:', promise, 'reason:', reason);
  queryForEntity();
});


async function queryForEntity() {
	const dgraphClientStub = await newClientStub();
	const dgraphClient = await newClient(dgraphClientStub);
	await setSchema(dgraphClient);
	await queryData(dgraphClient);
	await dgraphClientStub.close();
}

module.exports={db,queryForEntity}
