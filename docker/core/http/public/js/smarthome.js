$(function(){

// initialization  
var handlers = {}

   
addMenuItem('Home', showHomeGateways);  
addMenuItem('Device', showDevices);      
addMenuItem('Task', showTasks);      
addMenuItem('CityConsole', showCityConsole);    

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);

showHomeGateways();

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


function showCityConsole() 
{
    $('#info').html('to broadcast an audio message to all homes');
    
    var html = '<div id="dockerRegistration" class="form-horizontal"><fieldset>';                 
        
    html += '<div class="control-group"><label class="control-label" for="input01">Audio Message</label>';
    html += '<div class="controls"><select id="audioMessage"><option>Help</option><option>Fire Alarm</option><option>Emergency</option></select></div>'
    html += '</div>';    

    html += '<div class="control-group"><label class="control-label" for="input01"></label>';
    html += '<div class="controls"><button id="broadcastInfo" type="button" class="btn btn-primary">Broadcast</button>';
    html += '</div></div>';   
       
    html += '</fieldset></div>';

	$('#content').html(html);          
        
    // associate functions to clickable buttons
    $('#broadcastInfo').click(broadcast);      
}

function broadcast() 
{       
    var annoucementUpdate = {};

    annoucementUpdate.entityId = {
        id : 'Stream.CityConcole.Tokyo', 
        type: 'Announcement',
        isPattern: false
    };

    annoucementUpdate.attributes = {};  
    
    var selectedType = $('#audioMessage option:selected').val();
    console.log(selectedType);
        
    annoucement = {
        type: selectedType,
        audiofile: 'test.wav'
    }     
    
    annoucementUpdate.attributes.annoucement = {type: 'object', value: annoucement};    
	
    client.updateContext(annoucementUpdate).then( function(data) {
        console.log(data);       
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });         
}


function showHomeGateways() {    
    $('#info').html('to show the map with all devices and edges');    
    
    var html = '<div id="deviceList"></div><div id="map"  style="width: 800px; height: 600px"></div>';                
    $('#content').html(html);     
	
	// display the map
	var backgroundMap = showMap();              
    
    // show edge nodes on the map
    displayHomeGateway(backgroundMap);
}

function showMap()
{
    var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var osm = L.tileLayer(osmUrl, {maxZoom: 13, zoom: 10});    
    var myCenter = new L.LatLng(35.84, 139.48);    
    var map = new L.Map('map', {layers: [osm], 
                                center: myCenter, 
                                zoom: 10,
                                zoomControl:false});
                                
    var drawnItems = new L.FeatureGroup();
    map.addLayer(drawnItems);
	
	return map;                               
}

function displayHomeGateway(map)
{
    var queryReq = {}
    queryReq.entities = [{type:'Home', isPattern: true}];
    client.queryContext(queryReq).then( function(gatewayList) {
        showOnMap(map, gatewayList);
        showAsTable(gatewayList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });  
}


function showOnMap(map, homeGateways)
{
    if(homeGateways == null || homeGateways.length == 0){
        return        
    }        
    
    var homeIcon = L.icon({
    iconUrl: '/img/home.png',
        iconSize: [48, 48]
    });      
    
    for(var i=0; i<homeGateways.length; i++){
        var gateway = homeGateways[i];    
        
        console.log(gateway);
        
        latitude = gateway.metadata.location.value.latitude;
        longitude = gateway.metadata.location.value.longitude;
        homeID = gateway.entityId.id;
                    
        var marker = L.marker(new L.LatLng(latitude, longitude), {icon: homeIcon});
		marker.nodeID = homeID;
        marker.addTo(map).bindPopup(homeID);
    }                                
}

function showAsTable(homeGateways) 
{
    if(homeGateways == null || homeGateways.length == 0){
        $('#deviceList').html('');           
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>Attributes</th>';
    html += '<th>DomainMetadata</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<homeGateways.length; i++){
        var gateway = homeGateways[i];
		
        html += '<tr>'; 
		html += '<td>' + gateway.entityId.id + '</td>';
		html += '<td>' + JSON.stringify(gateway.attributes) + '</td>';        
		html += '<td>' + JSON.stringify(gateway.metadata) + '</td>';
		html += '</tr>';	                        
	}
       
    html += '</table>';  
    
	$('#deviceList').html(html);        
}


function showTasks() 
{
    $('#info').html('list of running data processing tasks');

    var queryReq = {}
    queryReq.entities = [{type:'Task', isPattern: true}];    
    
    client.queryContext(queryReq).then( function(taskList) {
        console.log(taskList);
        displayTaskList(taskList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });     
}

function displayTaskList(tasks) 
{
    if(tasks == null || tasks.length ==0){
        $('#content').html('');                   
        return
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>Type</th>';
    html += '<th>Attributes</th>';
    html += '<th>DomainMetadata</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<tasks.length; i++){
        var task = tasks[i];
		
        html += '<tr>'; 
		html += '<td>' + task.entityId.id + '</td>';
		html += '<td>' + task.entityId.type + '</td>'; 
		html += '<td>' + JSON.stringify(task.attributes) + '</td>';        
		html += '<td>' + JSON.stringify(task.metadata) + '</td>';
		html += '</tr>';	
	}
       
    html += '</table>'; 

	$('#content').html(html);   
}


function showDevices() 
{
    $('#info').html('list of all IoT devices');        

    var html = '<div id="deviceList"></div>';
	$('#content').html(html);        
    
    updateDeviceList();   
}

function updateDeviceList()
{
    var queryReq = {}
    queryReq.entities = [{id:'Device.*', isPattern: true}];           
    
    client.queryContext(queryReq).then( function(deviceList) {
        console.log(deviceList);
        displayDeviceList(deviceList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });    
}

function displayDeviceList(devices) 
{
    if(devices == null || devices.length == 0){
        $('#deviceList').html('');           
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>Type</th>';
    html += '<th>HomeID</th>';        
    html += '<th>Attributes</th>';
    html += '<th>DomainMetadata</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<devices.length; i++){
        var device = devices[i];
		
        html += '<tr>'; 
		html += '<td>' + device.entityId.id + '</td>';
		html += '<td>' + device.entityId.type + '</td>'; 
		html += '<td>' + device.metadata.homeID.value + '</td>';        
		html += '<td>' + JSON.stringify(device.attributes) + '</td>';        
		html += '<td>' + JSON.stringify(device.metadata) + '</td>';
		html += '</tr>';	                        
	}
       
    html += '</table>';  
    
	$('#deviceList').html(html);        
}



});



