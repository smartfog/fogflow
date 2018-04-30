$(function(){

// initialization  
var handlers = {}


// location of new device
var locationOfNewDevice = null;
// icon image for device registration
var iconImage = null;
var iconImageFileName = null;

var curMap = null;
var myCenter = new L.LatLng(37.990905, -1.131133);
var timerID = null;

   
addMenuItem('ProcessingFlow', showProcessingFlows);  
addMenuItem('SmartParking', showParking);      

//connect to the socket.io server via the NGSI proxy module
var ngsiproxy = new NGSIProxy();
ngsiproxy.setNotifyHandler(handleNotify);

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);
subscribeResult();

showProcessingFlows();

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

function subscribeResult()
{
    var subscribeCtxReq = {};    
    subscribeCtxReq.entities = [{type: 'ChildFound', isPattern: true}];
    subscribeCtxReq.reference =  'http://' + config.agentIP + ':' + config.agentPort;
    
    client.subscribeContext(subscribeCtxReq).then( function(subscriptionId) {
        console.log(subscriptionId);   
        ngsiproxy.reportSubID(subscriptionId);		
    }).catch(function(error) {
        console.log('failed to subscribe context');
    });
}

function handleNotify(contextObj)
{
    console.log(contextObj);
    
    if (curRequirement != null) {
        personsFound.push(contextObj);
    }
	
	var hash = window.location.hash;
	if (hash == '#Result') {
        updateResult();    
	}
}


function showProcessingFlows() 
{
    $('#info').html('to show the logical data processing flows behind this service');
       
    var html = '';
    html += '<div><img src="/img/smart-parking.png"></img></div>';    
    
    $('#content').html(html);	
}


function showParking() 
{
    $('#info').html('to show the map with all devices and edges');    

    var html = '';
    
    html += '<div style="margin-bottom: 10px;">';
    html += '<button id="ResetCar" type="button" class="btn btn-default">Reset</button>';        
    html += '</div>';
    
    html += '<div id="map"  style="width: 800px; height: 600px"></div>';                
   
    $('#content').html(html);       
    
    var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var osm = L.tileLayer(osmUrl, {maxZoom: 13, zoom: 13});    
    var map = new L.Map('map', {layers: [osm], 
                                center: myCenter, 
                                zoom: 13,
                                zoomControl:false});
                                
    var drawnItems = new L.FeatureGroup();
    map.addLayer(drawnItems);

    // show edge nodes on the map
    displayEdgeNodeOnMap(map);
       
    // show moving car
    drawConnectedCar(map);

    // remember the created map
    curMap = map;
}

function displayEdgeNodeOnMap(map)
{
    var queryReq = {}
    queryReq.entities = [{type:'Worker', isPattern: true}];
    queryReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'role=EdgeNode'}]}    
    client.queryContext(queryReq).then( function(edgeNodeList) {
        console.log(edgeNodeList);

        var edgeIcon = L.icon({
            iconUrl: '/img/gateway.png',
            iconSize: [48, 48]
        });      
        
        for(var i=0; i<edgeNodeList.length; i++){
            var worker = edgeNodeList[i];    
            
            console.log(worker);
            
            latitude = worker.attributes.physical_location.value.latitude;
            longitude = worker.attributes.physical_location.value.longitude;
            edgeNodeId = worker.entityId.id;
            
            console.log(latitude, longitude, edgeNodeId);
            
            var marker = L.marker(new L.LatLng(latitude, longitude), {icon: edgeIcon});
			marker.nodeID = edgeNodeId;
            marker.addTo(map).bindPopup(edgeNodeId);
		    marker.on('click', showRunningTasks);                        
            
            console.log('=======draw edge on the map=========');
        }                            
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });     
}

function showRunningTasks()
{
	var clickMarker = this;
	
    var queryReq = {}
    queryReq.entities = [{type:'Task', isPattern: true}];
    queryReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'worker=' + clickMarker.nodeID}]}    
    
    client.queryContext(queryReq).then( function(tasks) {
        console.log(tasks);		
		var content = "";		
        for(var i=0; i<tasks.length; i++){
        	var task = tasks[i];		
		
			if (task.attributes.status.value == "paused") {
				content += '<font color="red">' + task.attributes.id.value +'</font><br>';				
			} else {
				content += '<font color="green"><b>' + task.attributes.id.value + '</b></font><br>';												
			}		
		}
		
		clickMarker._popup.setContent(content);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    }); 	
}

function drawConnectedCar(map)
{    
    var taxiIcon = L.icon({
        iconUrl: '/img/taxi.png',
        iconSize: [80, 80]
    });    
                   
    var path = [[35.722266, 139.725322], [35.722266, 139.801368]];    
    var carMarker = L.Marker.movingMarker(path, [10000], {autostart: false, loop: true} );    
    carMarker.options.icon = taxiIcon;    
    
    map.addLayer(carMarker);

    carMarker.on('click', function() {
        if (carMarker.isRunning()) {
            console.log('timerID = ', timerID);
            if (timerID != null) {
                clearInterval(timerID);
            }
            carMarker.pause();
            carMarker.bindPopup('<b>Click me to start !</b>').openPopup();                
        } else {                
            carMarker.start();
            carMarker.bindPopup('<b>Click me to pause !</b>').openPopup();
            timerID = setInterval(function() {
            var mylocation = carMarker.getLatLng();
                updateMobileObject(mylocation);
            }, 1000);                
            console.log('timerID = ', timerID);                
        }
    });        

    carMarker.bindPopup('<b>Click me to start !</b>', {closeOnClick: false});
    carMarker.openPopup();        
}

function updateMobileObject(location)
{   
    //register a new device
    var movingCarObject = {};

    movingCarObject.entityId = {
        id : 'Device.ConnectedCar.01', 
        type: 'ConnectedCar',
        isPattern: false
    };

    movingCarObject.attributes = {};   
    movingCarObject.attributes.iconURL = {type: 'string', value: '/img/taxi.png'};
    movingCarObject.attributes.pullbased = {type: 'boolean', value: false};   	
	movingCarObject.attributes.location = {
        type: 'point',
        value: {'latitude': location.lat, 'longitude': location.lng}
    };    		 
    
    if (location.lng > 139.748122 && location.lng < 139.765496) {
        movingCarObject.attributes.raining = {type: 'boolean', value: true};   	            
    } else {
        movingCarObject.attributes.raining = {type: 'boolean', value: false};   	                    
    }    
            
    movingCarObject.metadata = {};    
    movingCarObject.metadata.location = {
        type: 'point',
        value: {'latitude': location.lat, 'longitude': location.lng}
    };               

    client.updateContext(movingCarObject).then( function(data) {
        console.log(data);        
    }).catch( function(error) {
        console.log('failed to update car object');
    });   
    
    console.log('observation of mobile sensor, raining = ', movingCarObject.attributes.raining);    
}

function readIconImage(input) 
{
    console.log('read icon image');
    if( input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            //var filename = $('#image_file').val();
            iconImage = e.target.result;
        }
        reader.readAsDataURL(input.files[0]);
        iconImageFileName = input.files[0].name;
    }    
}

function readContentImage(input)
{
    console.log('read content image'); 
    if( input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            contentImage = e.target.result;
        }
        reader.readAsDataURL(input.files[0]);
        contentImageFileName = input.files[0].name;
    }      
}

function registerNewDevice() 
{    
    console.log('register a new device'); 

    // take the inputs    
    var id = $('#deviceID').val();
    console.log(id);
    
    var type = $('#deviceType option:selected').val();
    console.log(type);        
    
    if( id == '' || type == '' || locationOfNewDevice == null) {
        alert('please provide the required inputs');
        return;
    }    
    
    console.log(locationOfNewDevice);
    
    //set up the the icon image according to the device type
    if (type == 'RainSensor') {
        iconImageFileName = 'rainsensor.png';
    } else if (type == 'SmartAwning') {
        iconImageFileName = 'awning.png';
    }
        
    //register a new device
    var newDeviceObject = {};

    newDeviceObject.entityId = {
        id : 'Device.' + type + '.' + id, 
        type: type,
        isPattern: false
    };

    newDeviceObject.attributes = {};   
    newDeviceObject.attributes.id = {type: 'string', value: id};    
    newDeviceObject.attributes.iconURL = {type: 'string', value: '/img/' + iconImageFileName};
    newDeviceObject.attributes.pullbased = {type: 'boolean', value: false};   
	
	newDeviceObject.attributes.location = {
        type: 'point',
        value: {'latitude': locationOfNewDevice.lat, 'longitude': locationOfNewDevice.lng}
    };    		 
            
    newDeviceObject.metadata = {};    
    newDeviceObject.metadata.location = {
        type: 'point',
        value: {'latitude': locationOfNewDevice.lat, 'longitude': locationOfNewDevice.lng}
    };               

    client.updateContext(newDeviceObject).then( function(data) {
        console.log(data);        
        showDevices();
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });          
}
       
    
function autoIDGenerator()
{
    var id = uuid();
    $('#deviceID').val(id);                   
}        


function uuid() {
    var uuid = "", i, random;
    for (i = 0; i < 32; i++) {
        random = Math.random() * 16 | 0;
        if (i == 8 || i == 12 || i == 16 || i == 20) {
            uuid += "-"
        }
        uuid += (i == 12 ? 4 : (i == 16 ? (random & 3 | 8) : random)).toString(16);
    }
    
    return uuid;
} 


});



