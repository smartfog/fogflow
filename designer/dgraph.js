const dgraph = require("dgraph-js");

var config_fs_name = './config.json';
var axios = require('axios')
var fs = require('fs');
const bodyParser = require('body-parser');


var globalConfigFile = require(config_fs_name)
var config = globalConfigFile.designer;
config.grpcPort = globalConfigFile.persistent_storage.port;


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

// Drop All - discard all data and start from a clean slate.
async function dropAll(dgraphClient) {
    const op = new dgraph.Operation();
    op.setDropAll(true);
    await dgraphClient.alter(op);
}


/*
   convert object domainmetadata data into string to store entity as single node  
*/

async function resolveDomainMetaData(data) {
    if ('domainMetadata' in data) {
        var len = data.domainMetadata.length
        for (var i = 0; i < len; i++) {
            if ('value' in data.domainMetadata[i]) {
                if (data.domainMetadata[i].type != 'global' && data.domainMetadata[i].type != 'stringQuery') {
                    data.domainMetadata[i].value = JSON.stringify(data.domainMetadata[i].value)
                }
            }
        }
    }
}

/*
   convert object attributes data into string to store entity as single node 
*/

async function resolveAttributes(data) {
    if ('attributes' in data) {
        var length = data.attributes.length
        for (var i = 0; i < length; i++) {
            if ('type' in data.attributes[i]) {
                if (data.attributes[i].type == 'object')
                    data.attributes[i].value = JSON.stringify(data.attributes[i].value)
                else {
                    data.attributes[i].value = data.attributes[i].value.toString()
                }
            }
        }
    }
}


function CtxElement2JSONObject(e) {
    var jsonObj = {};
    jsonObj.entityId = e.entityId;

    jsonObj.attributes = {}
    for (var i = 0; e.attributes && i < e.attributes.length; i++) {
        var attr = e.attributes[i];
        jsonObj.attributes[attr.name] = {
            type: attr.type,
            value: attr.value
        };
    }

    jsonObj.metadata = {}
    for (var i = 0; e.domainMetadata && i < e.domainMetadata.length; i++) {
        var meta = e.domainMetadata[i];
        jsonObj.metadata[meta.name] = {
            type: meta.type,
            value: meta.value
        };
    }

    return jsonObj;
}

function JSONObject2CtxElement(ob) {
    console.log('convert json object to context element')
    var contextElement = {};

    contextElement.entityId = ob.entityId;

    contextElement.attributes = [];
    if (ob.attributes) {
        for (key in ob.attributes) {
            attr = ob.attributes[key];
            contextElement.attributes.push({ name: key, type: attr.type, value: attr.value });
        }
    }

    contextElement.domainMetadata = [];
    if (ob.metadata) {
        for (key in ob.metadata) {
            meta = ob.metadata[key];
            contextElement.domainMetadata.push({ name: key, type: meta.type, value: meta.value });
        }
    }

    return contextElement;
}

/*
   insert data into database
*/

async function createData(dgraphClient, ctx) {
    const txn = dgraphClient.newTxn();
    try {
        const mu = new dgraph.Mutation();
        mu.setSetJson(ctx);
        const response = await txn.mutate(mu);
        await txn.commit();
    } finally {
        await txn.discard();
    }
}

/*
   send data to cloud broker
*/
async function sendData(contextEle) {
    var updateCtxReq = {};
    updateCtxReq.contextElements = [];
    updateCtxReq.updateAction = 'UPDATE'
    updateCtxReq.contextElements.push(contextEle)
    await axios({
        method: 'post',
        url: 'http://' + config.brokerIp + ':' + config.brokerPort + '/ngsi10/updateContext',
        data: updateCtxReq
    }).then(function (response) {
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
async function sendPostRequestToBroker(contextElement) {
    contextEle = contextElement['contextElements']

    for (var ctxEle = 0; ctxEle < contextEle.length; ctxEle = ctxEle + 1) {
        ctxEleReq = contextEle[ctxEle]
        if ('attributes' in ctxEleReq) {
            for (var ctxAttr = 0; ctxAttr < ctxEleReq.attributes.length; ctxAttr = ctxAttr + 1) {
                if (ctxEleReq.attributes[ctxAttr].type == 'object') {
                    const value = ctxEleReq.attributes[ctxAttr].value
                    ctxEleReq.attributes[ctxAttr].value = JSON.parse(value)
                }
                if (ctxEleReq.attributes[ctxAttr].type == 'integer') {
                    const value = ctxEleReq.attributes[ctxAttr].value
                    ctxEleReq.attributes[ctxAttr].value = parseInt(value)
                }
                if (ctxEleReq.attributes[ctxAttr].type == 'float') {
                    const value = ctxEleReq.attributes[ctxAttr].value
                    ctxEleReq.attributes[ctxAttr].value = parseFloat(value)
                }
                if (ctxEleReq.attributes[ctxAttr].type == 'boolean') {
                    const value = ctxEleReq.attributes[ctxAttr].value
                    if (value == 'false')
                        ctxEleReq.attributes[ctxAttr].value = false
                    else
                        ctxEleReq.attributes[ctxAttr].value = true
                }
            }
        }
        if ('domainMetadata' in ctxEleReq) {
            for (ctxdomain = 0; ctxdomain < ctxEleReq.domainMetadata.length; ctxdomain = ctxdomain + 1) {
                if ('value' in ctxEleReq.domainMetadata[ctxdomain]) {
                    if (ctxEleReq.domainMetadata[ctxdomain].type != 'global' && ctxEleReq.domainMetadata[ctxdomain].type != 'stringQuery') {
                        const value = ctxEleReq.domainMetadata[ctxdomain].value
                        ctxEleReq.domainMetadata[ctxdomain].value = JSON.parse(value)
                    }
                }
            }
        }

        await sendData(ctxEleReq)
    }
}

/*
   load all context elemented that have been saved into the dgraph databasefor getting the registered node
*/

async function loadContextElements(dgraphClient) {
    const query = `{
        contextElements(func: type(ContextData)) {
           {
            expand(_all_)
              }
           }
    }`;

    responseBody = await dgraphClient.newTxn().queryWithVars(query);
    const responsData = responseBody.getJson();
    console.log("all data---- ",responsData)
    sendPostRequestToBroker(responsData)
}


/*
   write entity data into dgraph
*/

async function WriteEntity(contextData1) {
    var contextData = contextData1;
    contextData.attribute = JSON.stringify(contextData.attribute)
    console.log("write entity *** ",contextData)
    try {
       // console.log("inside in write entity-----",contextData)
        const dgraphClientStub = await newClientStub();
        const dgraphClient = await newClient(dgraphClientStub);

        // if ('contextElements' in contextData) {
        //     contextData = contextData['contextElements']
        //     contextData = contextData[0]
        // }

        // await resolveAttributes(contextData)
        // await resolveDomainMetaData(contextData)

        contextData["dgraph.type"] = "ContextData"

        // console.log(contextData);

        await createData(dgraphClient, contextData);
        await dgraphClientStub.close();
    } catch (err) {
        console.log('DB ERROR::', err);
    }
}


/*
   delete entity
*/
async function DeleteEntity(updateCtxReq) {
    ctxElement = null

    console.log(updateCtxReq);

    if ('contextElements' in updateCtxReq) {
        elements = updateCtxReq['contextElements']
        ctxElement = elements[0]
    }

    if (ctxElement == null) {
        console.log("there is no context element in the request");
        return
    }

    eid = ctxElement.entityId.id;

    elements = await QueryNodeByEntityId(eid)

    for (var i = 0; i < elements.length; i++) {
        var element = elements[i];

        DeleteNodeById(element.uid)
    }
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
   query entity by uid
*/
async function QueryNodeByActionType(internalType) {
    try {
        const dgraphClientStub = await newClientStub();
        const dgraphClient = await newClient(dgraphClientStub);

        console.log("query context elements by internal type: " + internalType)

        const query = `query all($internalType: string) {
            contextElements(func: has(ContextData)) {
               {
                    uid
                    ContextData @filter(ge(internalType, $internalType)) {
                        internalType
                    }                     
               }
            }
        }`;
        const vars = { $internalType: internalType };
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

async function UpdateByUID(id_,attribute) {
    WriteEntity(attribute)
    //DeleteNodeById(id_);
}
/* */
// async function UpdateByUID(id_,attribute) {
//     try {
//         uid_ = id_; 
//         attr = JSON.stringify(attribute);
//         attr = attr.replace(/"/g, '\'');

//         const dgraphClientStub = await newClientStub();
//         const dgraphClient = await newClient(dgraphClientStub);

//         const txn = dgraphClient.newTxn();
//         try {
//             const query = `
//             query {
//                 user as var(func: eq(uid, "${uid_}"))
//             }`
          
//           const mu = new dgraph.Mutation();
//           mu.setSetNquads(`uid(user) <attribute> "${attr}"  .`);
//         //  mu.setCond(`@if(eq(len(user), 1))`);
          
//           const req = new dgraph.Request();
//           req.setQuery(query);
//           req.addMutations(mu);
//           req.setCommitNow(true);
          
//           await dgraphClient.newTxn().doRequest(req);

//           await txn.commit();
//         } catch (e) {
//           if (e === dgraph.ERR_ABORTED) {
//             // Retry or handle exception.
//           } else {
//             throw e;
//           }
//         } finally {
//           // Clean up. Calling this after txn.commit() is a no-op
//           // and hence safe.
//           await txn.discard();
//         }



//         // const txn = dgraphClient.newTxn();
//         // attr = JSON.stringify(attribute);
//         // attr = attr.replace(/"/g, '\'');
//         // console.log("**************** ddd ",attr);
//         // uid_ = id_; 
//         // console.log("aaaa ",uid_);
//         // try {
//         //     // const updateObj = {
//         //     //     uid: id,
//         //     //     attribute:JSON.stringify(attribute)
//         //     // };
//         //     const query = `
//         //         query {
//         //             user as var(func: eq(uid, "${uid_}"))
//         //         }`
//         //     const mu = new dgraph.Mutation();
//         //     //mu.setSetNquads(` uid(user) <attribute> "${attr}" .`);
//         //     // txn.setMutationsList([mu]);
//         //     // txn.setCommitNow(true);
            
//         //     //const mu = new dgraph.Mutation();
//         //     mu.setSetNquads(`uid(user) <attribute> "${attr}" .`);

//         //     const req = new dgraph.Request();
//         //     req.setQuery(query);
//         //     req.setMutationsList([mu]);
//         //     req.setCommitNow(true);

//         //     // Upsert: If wrong_email found, update the existing data
//         //     // or else perform a new mutation.
//         //     await dgraphClient.newTxn().doRequest(req);

//         //     //const response = await txn.mutate(mu);
//         //     await txn.commit();
//         // } finally {
//         //     await txn.discard();
//         // }

//         // await dgraphClientStub.close();
//     } catch (err) {
//         console.log('DB ERROR::', err);
//     }
// }

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



process.on('unhandledRejection', (reason, promise) => {
    console.log('Unhandled Rejection at:', promise, 'reason:', reason);
});



/*
   load all FogFlow internal entities from dgraph and then publish them into FogFlow broker
*/
async function LoadEntity() {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    await loadContextElements(dgraphClient);

    await dgraphClientStub.close();
}


/*
   retrieve all json objects with the specified data type
*/


/** {
    result(func: type(${dtype})) {
        {
             uid
             expand(_all_)
        }
     }
 } */
async function QueryJsonWithType(internalType) {
    console.log("&&&&&&&&&&&&&&&&&&&&&& in query ---");
    internalType = 'Operator'
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    const query = `{
        contextElements(func: type(ContextData)) {
           {
               uid
            expand(_all_)
              }
           }
    }`;

    const responseBody = await dgraphClient.newTxn().queryWithVars(query);
    //console.log("inside in query --- ",responseBody.getJson())

    await dgraphClientStub.close();

    return responseBody.getJson();
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


/*
   set the schema used by FogFlow
*/
const fogflow_schema = ` 

name: string @index(term) .
formattype: string @index(term) .
description: string @index(term) .
filepath: string @index(term) .
url: string @index(term) .
flavor: string @index(term) .
inputdata: string @index(term) .
version: string @index(term) .
attribute: string @index(term) .
internalType: string @index(term) .
updateAction: string @index(term) .




type ContextData {
    attribute
    internalType
    updateAction
}

`;

async function InitSchema() {
    const dgraphClientStub = await newClientStub();
    const dgraphClient = await newClient(dgraphClientStub);

    const op = new dgraph.Operation();
    op.setSchema(fogflow_schema);
    await dgraphClient.alter(op);

    await dgraphClientStub.close();
}


async function Init(){
    try {
        await InitSchema();
        console.log("init schema")
        await LoadEntity();
        console.log("load entity")        
    }catch(e) { 
        console.error('==========' + e.details + '============='); 
        console.log('Retrying to connect to dgraph');
        setTimeout(Init, 2000);
    }
}


process.on('unhandledRejection', (reason, promise) => {
    console.log('Retrying to connect to dgraph');
    //Init();
});


module.exports = { Init, WriteEntity, DeleteEntity, DeleteNodeById, WriteJsonWithType, QueryJsonWithType, DropAll, QueryNodeByActionType,UpdateByUID }
