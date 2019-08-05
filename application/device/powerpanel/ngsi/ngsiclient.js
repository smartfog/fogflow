(function() {
    
function CtxElement2JSONObject(e) {
    var jsonObj = {};
    jsonObj.entityId = e.entityId;

    jsonObj.attributes = {}    
    for(var i=0; e.attributes && i<e.attributes.length; i++) {
        var attr = e.attributes[i];
        jsonObj.attributes[attr.name] = {
            type: attr.type, 
            value: attr.value
        };
    }
    
    jsonObj.metadata = {}
    for(var i=0; e.domainMetadata && i<e.domainMetadata.length; i++) {
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
    if(ob.attributes) {
        for( key in ob.attributes ) {
            attr = ob.attributes[key];
            contextElement.attributes.push({name: key, type: attr.type, value: attr.value});
        }
    }
    
    contextElement.domainMetadata = [];
    if(ob.metadata) {
        for( key in ob.metadata ) {
            meta = ob.metadata[key];
            contextElement.domainMetadata.push({name: key, type: meta.type, value: meta.value});
        }
    }    

    return contextElement;
}  
    
var NGSI10Client = (function() {
    // initialized with the broker URL
    var NGSI10Client = function(url) {
        this.brokerURL = url;
    };
    
    // update context 
    NGSI10Client.prototype.updateContext = function updateContext(ctxObj) {
        contextElement = JSONObject2CtxElement(ctxObj);
        
        var updateCtxReq = {};
        updateCtxReq.contextElements = [];
        updateCtxReq.contextElements.push(contextElement)
        updateCtxReq.updateAction = 'UPDATE'
         
		console.log(updateCtxReq);
		      
        return axios({
            method: 'post',
            url: this.brokerURL + '/updateContext',
            data: updateCtxReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
    };
    
    // delete context 
    NGSI10Client.prototype.deleteContext = function deleteContext(entityId) {
        var contextElement = {};
        contextElement.entityId = entityId
        
        var updateCtxReq = {};
        updateCtxReq.contextElements = [];
        updateCtxReq.contextElements.push(contextElement)
        updateCtxReq.updateAction = 'DELETE'
        
        return axios({
            method: 'post',
            url: this.brokerURL + '/updateContext',
            data: updateCtxReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
    };    
    
    // query context
    NGSI10Client.prototype.queryContext = function queryContext(queryCtxReq) {        
        return axios({
            method: 'post',
            url: this.brokerURL + '/queryContext',
            data: queryCtxReq
        }).then( function(response){
            if (response.status == 200) {
                var objectList = [];
                var ctxElements = response.data.contextResponses;
                for(var i=0; ctxElements && i<ctxElements.length; i++){                    
                    console.log(ctxElements[i].contextElement);
                    console.log('===========context element=======');
                    console.log(ctxElements[i].contextElement)
                    var obj = CtxElement2JSONObject(ctxElements[i].contextElement);
                    objectList.push(obj);
                }
                return objectList;
            } else {
                return null;
            }
        });
    };    
        
    // subscribe context
    NGSI10Client.prototype.subscribeContext = function subscribeContext(subscribeCtxReq) {        
        return axios({
            method: 'post',
            url: this.brokerURL + '/subscribeContext',
            data: subscribeCtxReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data.subscribeResponse.subscriptionId;
            } else {
                return null;
            }
        });
    };    

    // unsubscribe context    
    NGSI10Client.prototype.unsubscribeContext = function unsubscribeContext(sid) {
        var unsubscribeCtxReq = {};
        unsubscribeCtxReq.subscriptionId = sid;
        
        return axios({
            method: 'post',
            url: this.brokerURL + '/unsubscribeContext',
            data: unsubscribeCtxReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
    };        
    
    return NGSI10Client;
})();

var NGSI9Client = (function() {
    // initialized with the address of IoT Discovery
    var NGSI9Client = function(url) {
        this.discoveryURL = url;
    };
        
    NGSI9Client.prototype.findNearbyIoTBroker = function findNearbyIoTBroker(mylocation, num) 
    {
        var discoveryReq = {};    
        discoveryReq.entities = [{type: 'IoTBroker', isPattern: true}];              
    
        var nearby = {};
        nearby.latitude = mylocation.latitude;
        nearby.longitude = mylocation.longitude;
        nearby.limit = num;
        
        discoveryReq.restriction = {
            scopes: [{
                scopeType: 'nearby',
                scopeValue: nearby
            }]
        };
    
        return this.discoverContextAvailability(discoveryReq).then( function(response) {
            if (response.errorCode.code == 200) {
                var brokers = [];
                for(i in response.contextRegistrationResponses) {
                    contextRegistrationResponse = response.contextRegistrationResponses[i];
                    var providerURL = contextRegistrationResponse.contextRegistration.providingApplication;
                    if (providerURL != '') {
                        brokers.push(providerURL);
                    }
                }
                return brokers;
            } else {
                return nil;
            }            
        });
    }
            
    // discover availability
    NGSI9Client.prototype.discoverContextAvailability = function discoverContextAvailability(discoverReq) {        
        return axios({
            method: 'post',
            url: this.discoveryURL + '/discoverContextAvailability',
            data: discoverReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
    };               
    
    return NGSI9Client;
})();

// initialize the exported object for this module, both for nodejs and browsers
if (typeof module !== 'undefined' && typeof module.exports !== 'undefined'){
    this.axios = require('axios')    
    module.exports.NGSI10Client = NGSI10Client; 
    module.exports.NGSI9Client = NGSI9Client;   
    module.exports.CtxElement2JSONObject = CtxElement2JSONObject;
    module.exports.JSONObject2CtxElement = JSONObject2CtxElement;    
} else {
    window.NGSI10Client = NGSI10Client;  
    window.NGSI9Client = NGSI9Client;
    window.CtxElement2JSONObject = CtxElement2JSONObject;
    window.JSONObject2CtxElement = JSONObject2CtxElement;
}

})();
