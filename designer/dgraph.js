const dgraph = require("dgraph-js");

var axios = require('axios')
var fs = require('fs');
const bodyParser = require('body-parser');

var config_fs_name = './config.json';
var globalConfigFile = require(config_fs_name)
var config = globalConfigFile.designer;
config.grpcPort = globalConfigFile.persistent_storage.port;


// cache all FogFlow metadata
var operatorList = {};
var dockerimageList = {};
var servicetopologyList = {};
var serviceIntentList = {};
var fogfunctionList = {};


if ('host_ip' in globalConfigFile.persistent_storage){
    config.HostIp = globalConfigFile.persistent_storage.host_ip;
} else {
    config.HostIp = globalConfigFile.my_hostip;
}


if ('host_ip' in globalConfigFile.broker){
    config.brokerIp = globalConfigFile.broker.host_ip;    
} else {
    config.brokerIp = globalConfigFile.my_hostip    
}


config.brokerPort = globalConfigFile.broker.http_port


/*
   creating grpc client for making connection with dgraph
*/

async function newClientStub() {
    var clientStub = new dgraph.DgraphClientStub(config.HostIp + ":" + config.grpcPort);
    return clientStub;
}

// Create a client.

async function newClient(clientStub) {
    return new dgraph.DgraphClient(clientStub);
}

/*
   retrieve an object with the specified data type and name
*/
function GetObject(dtype, name) {
    var resultObj = {};
    
    switch(dtype){
        case 'Operator':
            if (name in operatorList){
                resultObj = operatorList[name];
            }
            break;
        case 'DockerImage':
            if (name in dockerimageList){
                resultObj = dockerimageList[name];
            }        
            break;
        case 'ServiceTopology':
            if (name in servicetopologyList){
                resultObj = servicetopologyList[name];
            }        
            break;
        case 'ServiceIntent':
            if (name in serviceIntentList){
                resultObj = serviceIntentList[name];
            }        
            break;            
        case 'FogFunction':
            if (name in fogfunctionList){
                resultObj = fogfunctionList[name];
            }        
            break;
    }
    
    return resultObj;
}



/*
   retrieve the list of objects with the specified data type
*/
function GetObjectList(dtype) {    
    if (dtype == 'Operator') {
        return operatorList;
    } else if (dtype == 'DockerImage') {
        return dockerimageList;        
    } else if (dtype == 'ServiceTopology') {
        return servicetopologyList;        
    } else if (dtype == 'ServiceIntent') {
        return serviceIntentList;    
    } else if (dtype == 'FogFunction') {
        return fogfunctionList;        
    }            
}

/*
   write an object with the specified data type and name
*/
async function PutObject(dtype, name) {
    if (dtype == 'Operator') {
        operatorList[name];
    } else if (dtype == 'DockerImage') {
        return dockerimageList;        
    } else if (dtype == 'ServiceTopology') {
        return servicetopologyList;        
    } else if (dtype == 'ServiceIntent') {
        return serviceIntentList;    
    } else if (dtype == 'FogFunction') {
        return fogfunctionList;        
    }    
}


/*
   load all FogFlow internal entities from dgraph and then publish them into FogFlow broker
*/
async function LoadEntity() {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    var fogfunctions = dgraph.GetObjectList('FogFunction');        



    await dgraphClientStub.close();
}

/*
   query entity by uid
*/
async function QueryNodeByEntityId(eid) {
    try {
        const dgraphClientStub = await newClientStub();
        const dgraphClient = await newClient(dgraphClientStub);

        console.log("query context elements by entity ID: " + eid)

        const query = `query all($eid: string) {
            contextElements(func: has(entityId)) {
               {
                    uid
                    entityId @filter(ge(id, $eid)) {
                         id
                    }                     
               }
            }
        }`;
        const vars = { $eid: eid };
        const responseBody = await dgraphClient.newTxn().queryWithVars(query, vars);

        console.log(responseBody.getJson())

        ctxElements = responseBody.getJson().contextElements;

        await dgraphClientStub.close();

        return ctxElements;
    } catch (err) {
        console.log('DB ERROR::', err);
    }
}



/*
   write model profiles into dgraph
*/
async function WriteJsonWithType(jsonData, dtype) {
    try {
        const dgraphClientStub = await newClientStub();
        const dgraphClient = await newClient(dgraphClientStub);

        if (dtype != "") {
            jsonData["dgraph.type"] = dtype
        }
        
        console.log(jsonData)

        const txn = dgraphClient.newTxn();
        try {
            const mu = new dgraph.Mutation();
            mu.setSetJson(jsonData);
            const response = await txn.mutate(mu);
            await txn.commit();
        } finally {
            await txn.discard();
        }

        await dgraphClientStub.close();
    } catch (err) {
        console.log('DB ERROR::', err);
    }
}

/*
   retrieve all json objects with the specified data type
*/
async function QueryJsonWithType(dtype) {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    const query = `{
        result(func: type(${dtype})) {
            {
             uid
             expand(_all_)
            }
        }            
     } `;

    const responseBody = await dgraphClient.newTxn().queryWithVars(query);

    await dgraphClientStub.close();

    return responseBody.getJson()['result'];
}



/*
   delete entity by uid
*/
async function DeleteNodeById(id) {
    try {
        const dgraphClientStub = await newClientStub();
        const dgraphClient = await newClient(dgraphClientStub);

        const txn = dgraphClient.newTxn();
        try {
            const deleteJsonObj = {
                uid: id,
            };

            const mu = new dgraph.Mutation();
            mu.setDeleteJson(deleteJsonObj);
            const response = await txn.mutate(mu);
            await txn.commit();
        } finally {
            await txn.discard();
        }

        await dgraphClientStub.close();
    } catch (err) {
        console.log('DB ERROR::', err);
    }
}


/*
    clean up the entire graph database
*/
async function DropAll() {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    const op = new dgraph.Operation();
    op.setDropAll(true);

    await dgraphClient.alter(op);
}

async function InitSchema(fogflow_schema) {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    const op = new dgraph.Operation();
    op.setSchema(fogflow_schema);
    await dgraphClient.alter(op);

    await dgraphClientStub.close();
}


async function Init(myGraphQLSchema){
    try {
        await InitSchema(myGraphQLSchema);
        console.log("init schema")
        await LoadEntity();
        console.log("load entity")        
    }catch(e) { 
        console.log('[Dgraph] Retrying to connect to dgraph');
        setTimeout(Init, 2000);
    }
}


module.exports = { Init, GetObject, PutObject, DeleteNodeById, GetObjectList, WriteJsonWithType, QueryJsonWithType, DropAll}