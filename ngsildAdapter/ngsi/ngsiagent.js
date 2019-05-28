'use strict';

var http = require('http'),
    express = require('express'),
    logger = require('logops'),
    bodyParser = require('body-parser'),
    northboundServer,    
	notifyHandler,
	adminHandler;


function CtxElement2JSONObject(e) {
    var jsonObj = {};
    jsonObj.entityId = e.entityId;

    jsonObj.attributes = {}    
    for(var i=0; e.attributes && i<e.attributes.length; i++) {
        var attr = e.attributes[i];
        jsonObj.attributes[attr.name] = {
            type: attr.type, 
            value: attr.contextValue
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


console.log('CtxElement2JSONObject jsonObject',jsonObj);    
    return jsonObj;
}    

function ensureType(req, res, next) {
    console.log("inside ensureType....");
    if (req.is('json')) {
        next();
    } else {
        next(new errors.UnsupportedContentType(req.headers['content-type']));
    }
	
}

function loadContextRoutes(router) {
    var notifyHandlers = [
            ensureType,
            handleNotify
        ],
		adminHandlers = [
            ensureType,
            handleAdmin
        ];

    router.post('/notifyContext', notifyHandlers);
    console.log('router.post/notifyContext');
    router.post('/admin', adminHandlers);		
    console.log('router.post/admin');
}
		
function handleError(error, req, res, next) {
    var code = 500;

    logger.debug('Error [%s] handing request: %s', error.name, error.message);

    if (error.code && String(error.code).match(/^[2345]\d\d$/)) {
        code = error.code;
    }

    res.status(code).json({
        name: error.name,
        message: error.message
    });
}

function traceRequest(req, res, next) {
    logger.debug('Request for path [%s] from [%s]', req.path, req.get('host'));

    next();
}

function setNotifyHandler(newHandler) {
    notifyHandler = newHandler;
 console.log('notifyHandler set!');
}

function setAdminHandler(newHandler) {
    adminHandler = newHandler;
 console.log('adminHandler set!');

}

function readContextElements(body) {
 console.log('inside from readContextElements', ctxObjects);
	var ctxObjects = [];
	
	for(var i = 0; i < body.contextResponses.length; i++){
		var response = body.contextResponses[i];
		if(response.statusCode.code == '200'){		       
			var flexObj = CtxElement2JSONObject(response.contextElement);
			ctxObjects.push(flexObj);
		}
	}
 console.log('ctxObjects from readContextElements', ctxObjects);

	return ctxObjects;
}

function handleNotify(req, res, next) {
	if (notifyHandler) {
        console.log('req........ in handleNotify...',req);
        console.log('res........ in handleNotify...',res);
        logger.debug('Handling notification from [%s]', req.get('host'));		

		var ctxs = readContextElements(req.body);
		notifyHandler(req, ctxs, res);		
        console.log('ctxs........ in handleNotify...',ctxs);
        next();
    } else {
        var errorNotFound = new Error({
            message: 'Notification handler not found'
        });
        logger.error('Tried to handle a notification before notification handler was established.');
        next(errorNotFound);
    }
}

function handleAdmin(req, res, next) {
	if (adminHandler) {
        //console.log('req........ in handleAdmin...',req);
        //console.log('res........ in handleAdmin...',res);
        console.log('Handling admin from [%s]', req.get('host'));
        console.log('reqBody array--------', req.body);		

//	req.body.command.value='SET_OUTPUTS';
//        req.body.id.value=1;
//        req.body.type.value='healthRisk';


//        req.body.entityId={};
//        req.body.entityId.id=1;
//        req.body.entityId.type='healthRisk';
	
	if( Array.isArray(req.body) ) {                   
			adminHandler(req, req.body, res);			
		}
		
        next();
    } else {
        var errorNotFound = new Error({
            message: 'admin handler not found'
        });
        logger.error('Tried to handle an admin command before admin handler was established.');
        next(errorNotFound);
    }
}


function start(port, callback) {
    var baseRoot = '/';

    northboundServer = {
        server: null,
        app: express(),
        router: express.Router()
    };

    logger.info('Starting NGSI Agent listening on port [%s]', port);

    northboundServer.app.set('port', port);
    northboundServer.app.set('host', '0.0.0.0');
    northboundServer.app.use(bodyParser.json());
	
    northboundServer.app.use(baseRoot, northboundServer.router);

    loadContextRoutes(northboundServer.router);
    console.log('done with loadContextRoutes');
    northboundServer.app.use(handleError);


    northboundServer.server = http.createServer(northboundServer.app);
    console.log('done with http.createServer');
    northboundServer.server.listen(northboundServer.app.get('port'), northboundServer.app.get('host'), callback);
    console.log('done with northboundServer.server.listen');
}

function stop() {
    logger.info('Stopping NGSI Agent');

    if (northboundServer) {
        northboundServer.server.close();
    } 
}

exports.start = start;
exports.stop = stop;
exports.setNotifyHandler  = setNotifyHandler;
exports.setAdminHandler  = setAdminHandler;
