$(function(){

// initialization  
var handlers = {}

   
addMenuItem('Home', showHome);  
addMenuItem('Device', showDevices);      
addMenuItem('Task', showTasks);      

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);

showHome();

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


function showHome() {    
    $('#info').html('use case for smart industry');    
    
    var html = '';
    html += '<div><img src="/img/factory.png"></img></div>';    
    
    $('#content').html(html);	
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
    $('#info').html('list of all IoT devices in Factor A');        

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
    html += '<th>Action</th>';        
    //html += '<th>Attributes</th>';
    //html += '<th>DomainMetadata</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<devices.length; i++){
        var device = devices[i];
		
        html += '<tr>'; 
		html += '<td>' + device.entityId.id + '</td>';
		html += '<td>' + device.entityId.type + '</td>'; 
		html += '<td>' + " " + '</td>';         
		//html += '<td>' + JSON.stringify(device.attributes) + '</td>';        
		//html += '<td>' + JSON.stringify(device.metadata) + '</td>';
		html += '</tr>';	                        
	}
       
    html += '</table>';  
    
	$('#deviceList').html(html);        
}



});



