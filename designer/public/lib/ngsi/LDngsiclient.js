(function () {

    function CtxElement2JSONObject(e) {
        var jsonObj = {};
        for (var ctxElement in e) {
            jsonObj[ctxElement] = e[ctxElement]
        }
        return jsonObj;
    }

    function JSONObject2CtxElement(ctxObj) {
        console.log('convert json object to context element')
        var ctxElement = {};

        ctxElement['id'] = ctxObj['id']
        ctxElement['type'] = ctxObj['type']

        for (key in ctxObj) {
            if (key != 'id' && key != 'type' && key != 'modifiedAt' && key != 'createdAt' && key != 'observationSpace' && key != 'operationSpace' && key != 'location' && key != '@context') {
                ctxElement[key] = ctxObj[key]
            }
        }

        return ctxElement

    }


    var NGSILDclient = (function () {
        // initialized with the broker URL
        var NGSILDclient = function (url) {
            if (url.includes('ngsi10')) {
                url = url.substring(0, url.lastIndexOf("/"));
            }
            this.brokerURL = url;
        };

        // update context 
        NGSILDclient.prototype.updateContext = function updateContext(ctxObj) {
            updateCtxReq = JSONObject2CtxElement(ctxObj)
            const updateCtxElements = []
            updateCtxElements.push(updateCtxReq)
            return axios({
                method: 'post',
                url: this.brokerURL + '/v1/entityOperations/upsert',
                headers: {
                    'content-type': 'application/json',
                    'Accept': 'application/ld+json',
                    'Link': '<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld>; rel="http://www.w3.org/ns/json-ld#context"; type="application/ld+json"'
                },
                data: updateCtxElements
            }).then(function (response) {
                console.log("response", response)
                if (response.status == 204) {
                    console.log("Successfully updated")
                    return response.data;
                } else if (response.status == 207) {
                    console.log("Failed to update some entities")
                    return response.data;
                } else {
                    return null;
                }
            });
        };

        // delete context 
        NGSILDclient.prototype.deleteContext = function deleteContext(entityId, isFromDeviceTab = false) {
            ldURL = this.brokerURL + '/ngsi-ld/v1/entities/' + entityId
            if (isFromDeviceTab) ldURL = this.brokerURL + '/v1/entities/' + entityId
            return axios({
                method: 'delete',
                headers: {
                    'content-type': 'application/json',
                    'Accept': 'application/ld+json',
                    'Link': '<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld>; rel="http://www.w3.org/ns/json-ld#context"; type="application/ld+json"'
                },
                url: ldURL,
                data: { entities: [{ type: "IoTBroker", isPattern: true }] }
            }).then(function (response) {
                if (response.status == 201) {
                    return response.data;
                } else if (response.status == 204) {
                    return response;
                } else {
                    return null;
                }
            });
        };

        NGSILDclient.prototype.GetEntitiesContext = function GetEntitiesContext(type_) {
            console.log("inside in get entities context", this.brokerURL);

            return axios({
                method: "get",
                url: this.brokerURL + '/v1/entities',
                headers: {
                    'Accept': 'application/ld+json',
                    'content-type': 'application/json'
                },
                params: { type: type_ },
                data: { entities: [{ type: "IoTBroker", isPattern: true }] }
            }).then(function (response) {
                if (response.status == 200) {
                    return response.data;
                } else {
                    return null;
                }
            });
        };

        // query context
        NGSILDclient.prototype.queryContext = function queryContext(queryCtxReq) {
            return axios({
                method: 'post',
                url: this.brokerURL + '/v1/entityOperations/query',
                headers: {
                    'content-type': 'application/json',
                    'Accept': 'application/ld+json',
                    'Link': '<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld>; rel="http://www.w3.org/ns/json-ld#context"; type="application/ld+json"'
                },
                data: queryCtxReq
            }).then(function (response) {
                //console.log(response);
                if (response.status == 200) {
                    var objectList = [];
                    var ctxElements = response.data;
                    objectList.push(ctxElements)
                    return objectList;
                } else {
                    return null;
                }
            });
        };

        // subscribe context
        NGSILDclient.prototype.subscribeContext = function subscribeContext(subscribeCtxReq) {
            return axios({
                method: 'post',
                url: this.brokerURL + '/v1/subscriptions/',
                headers: {
                    'content-type': 'application/json',
                    'Accept': 'application/ld+json',
                    'Link': '<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld>; rel="http://www.w3.org/ns/json-ld#context"; type="application/ld+json"'
                },
                data: subscribeCtxReq
            }).then(function (response) {
                console.log("response", response)
                if (response.status == 201) {
                    //console.log("response.subscribeResponse.subscriptionId",response.subscribeResponse.subscriptionId)
                    return response.data.subscribeResponse.subscriptionId;
                } else {
                    console.log("response")
                    return null;
                }
            });
        };

        // unsubscribe context    
        NGSILDclient.prototype.unsubscribeContext = function unsubscribeContext(sid) {
            var unsubscribeCtxReq = {};
            unsubscribeCtxReq.subscriptionId = sid;

            return axios({
                method: 'post',
                url: this.brokerURL + '/unsubscribeContext',
                data: unsubscribeCtxReq
            }).then(function (response) {
                if (response.status == 200) {
                    return response.data;
                } else {
                    return null;
                }
            });
        };

        return NGSILDclient;
    })();

    var NGSI9Client = (function () {
        // initialized with the address of IoT Discovery
        var NGSI9Client = function (url) {
            this.discoveryURL = url;
        };

        NGSI9Client.prototype.findNearbyIoTBroker = function findNearbyIoTBroker(mylocation, num) {
            var discoveryReq = {};
            discoveryReq.entities = [{ type: 'IoTBroker', isPattern: true }];

            var nearby = {};
            nearby.latitude = mylocation.latitude;
            nearby.longitude = mylocation.longitude;
            nearby.limit = num;

            discoveryReq.restriction = {
                scopes: [{
                    type: 'nearby',
                    value: nearby
                }]
            };

            return this.discoverContextAvailability(discoveryReq).then(function (response) {
                if (response.errorCode.code == 200) {
                    var brokers = [];
                    for (i in response.contextRegistrationResponses) {
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
            }).then(function (response) {
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
    if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
        this.axios = require('axios')
        module.exports.NGSILDclient = NGSILDclient;
        module.exports.NGSI9Client = NGSI9Client;
        module.exports.CtxElement2JSONObject = CtxElement2JSONObject;
        module.exports.JSONObject2CtxElement = JSONObject2CtxElement;
    } else {
        window.NGSILDclient = NGSILDclient
        window.NGSI9Client = NGSI9Client;
        window.CtxElement2JSONObject = CtxElement2JSONObject;
        window.JSONObject2CtxElement = JSONObject2CtxElement;
    }

})();
