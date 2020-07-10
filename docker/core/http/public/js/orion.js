$(function(){

// initialization  
var handlers = {}


addMenuItem('Test', testIntegration);  

//connect to the socket.io server via the NGSI proxy module
var ngsiproxy = new NGSIProxy();
ngsiproxy.setNotifyHandler(handleNotify);

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);

testIntegration();

$(window).on('hashchange', function() {
    var hash = window.location.hash;
    selectMenuItem(location.hash.substring(1));
});

function addMenuItem(name, func) {
    handlers[name] = func; 
    $('#menu').append('<li id="' + name + '"><a href="' + '#' + name + '">' + name + '</a></li>');
}

function selectMenuItem(name) {
    $('#menu li').removeClass('active');
    var element = $('#' + name);
    element.addClass('active');
    
    var handler = handlers[name];
    handler();
}

function handleNotify(contextObj)
{
    console.log(contextObj);
    
	var curText = $('#fogflowText').val();
	curText += JSON.stringify(contextObj);
	
	$('#fogflowText').val(curText);        	
}


function testIntegration() 
{
    $('#info').html('to showcase how FogFlow can be integrated with FIWARE Orion Broker');
       
    var html = '';
	
    html += '<div class="input-prepend">'; 
	html += '<button id="sendUpdates" type="button" class="btn btn-primary">send update</button>';	
    html += '</div>';  	
	
    html += '<div class="input-prepend">'; 
    html += '<label class="control-label">Orion Broker at</label>';
	html += '<input type="text" class="input-xlarge" id="orionBroker">';	   
    html += '</div>';  	
	
    html += '<div class="input-prepend">'; 
	html += '<button id="fogflowSubscribe" type="button" class="btn btn-default">subscribe to FogFlow</button>';	
	html += '<textarea id="fogflowText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
    html += '</div>';    
	
    html += '<div class="input-prepend">';    	
	html += '<button id="orionSubscribe" type="button" class="btn btn-default">subscribe to Orion</button>';
    html += '<textarea id="orionText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';    
	html += '</div>';    

    $('#content').html(html);	
	
	$('#orionBroker').val("localhost:1026");
	
    $('#sendUpdates').click(sendUpdate);			
	
    $('#fogflowSubscribe').click( function() {
        var orionBroker = $('#orionBroker').val();		
		$('#fogflowText').val("subscribe to FogFlow broker on behalf of Orion Broker at " + orionBroker + "<br>");        
		subscribeFogFlow('PowerPanel', orionBroker);
    });	
	
	
    $('#orionSubscribe').click( function() {
        var orionBroker = $('#orionBroker').val();		
		$('#orionText').val("subscribe to Orion broker to check the received updates");        
    });		
}

function sendUpdate()
{
	var profile = {
	    "location": {
        		"latitude": 35.508822,	
        		"longitude": 139.696867
    		},
    		"iconURL": "/img/shop.png",    
		"type": "PowerPanel",
		"id": "05"	
	};
	
    var ctxObj = {};
    ctxObj.entityId = {
        id: 'Stream.' + profile.type + '.' + profile.id,
        type: profile.type,
        isPattern: false
    };
    
    ctxObj.attributes = {};
    
    var degree = Math.floor((Math.random() * 100) + 1);        
    ctxObj.attributes.usage = {
        type: 'integer',
        value: degree
    };
    ctxObj.attributes.deviceID = {
        type: 'string',
        value: profile.type + '.' + profile.id
    };   	     
    
    ctxObj.metadata = {};
    
    ctxObj.metadata.location = {
        type: 'point',
        value: profile.location
    }; 
    ctxObj.metadata.shop = {
        type: 'string',
        value: profile.id
    };	          
    
    client.updateContext(ctxObj).then( function(data) {
        console.log(data);
    }).catch(function(error) {
        console.log('failed to update context');
    });    	
}

function subscribeFogFlow(entityType, orionBroker)
{
    var subscribeCtxReq = {};    
    subscribeCtxReq.entities = [{type: entityType, isPattern: true}];
    subscribeCtxReq.reference =  'http://' + orionBroker + '/v2';
    
    client.subscribeContext4Orion(subscribeCtxReq).then( function(subscriptionId) {
        console.log(subscriptionId);   
        ngsiproxy.reportSubID(subscriptionId);		
    }).catch(function(error) {
        console.log('failed to subscribe context');
    });	
}

});



