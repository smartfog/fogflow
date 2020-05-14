'use strict';

$(function(){

// initialize the menu bar
var handlers = {}

var geoscope = {
    type: 'string',
    value: 'all'
};

var CurrentScene = null;

// icon image for device registration
var iconImage = null;
var iconImageFileName = null;
// content image for camera devices
var contentImage = null;
var contentImageFileName = null;

// the list of all registered operators
var operatorList = [];

// design board
var blocks = null;

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);

//DGraph
//to interact with Designer
var clientDes = new NGSIDesClient(config.designerIP+':'+config.webSrvPort);

var myToplogyExamples = [
{
    topology: {"name":"anomaly-detection","description":"detect anomaly events in shops","tasks":[{"name":"Counting","operator":"counter","input_streams":[{"selected_type":"Anomaly","selected_attributes":[],"groupby":"ALL","scoped":true}],"output_streams":[{"entity_type":"Stat"}]},{"name":"Detector","operator":"anomaly","input_streams":[{"selected_type":"PowerPanel","selected_attributes":[],"groupby":"EntityID","scoped":true},{"selected_type":"Rule","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"Anomaly"}]}]},
    designboard: {"edges":[{"id":2,"block1":3,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]},{"id":3,"block1":2,"connector1":["outputs","output",0],"block2":3,"connector2":["in","input"]},{"id":4,"block1":4,"connector1":["stream","output"],"block2":2,"connector2":["streams","input"]},{"id":5,"block1":5,"connector1":["stream","output"],"block2":2,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":202,"y":-146,"type":"Task","module":null,"values":{"name":"Counting","operator":"counter","outputs":["Stat"]}},{"id":2,"x":-194,"y":-134,"type":"Task","module":null,"values":{"name":"Detector","operator":"anomaly","outputs":["Anomaly"]}},{"id":3,"x":4,"y":-18,"type":"Shuffle","module":null,"values":{"selectedattributes":["all"],"groupby":"ALL"}},{"id":4,"x":-447,"y":-179,"type":"EntityStream","module":null,"values":{"selectedtype":"PowerPanel","selectedattributes":["all"],"groupby":"EntityID","scoped":true}},{"id":5,"x":-438,"y":-5,"type":"EntityStream","module":null,"values":{"selectedtype":"Rule","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]} 
},{
    topology: {"name":"child-finder","description":"search for a lost child based on face recognition","tasks":[{"name":"childfinder","operator":"facefinder","input_streams":[{"selected_type":"Camera","selected_attributes":[],"groupby":"EntityID","scoped":true},{"selected_type":"ChildLost","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"ChildFound"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]},{"id":2,"block1":3,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":7,"y":-107,"type":"Task","module":null,"values":{"name":"childfinder","operator":"facefinder","outputs":["ChildFound"]}},{"id":2,"x":-292,"y":-161,"type":"EntityStream","module":null,"values":{"selectedtype":"Camera","selectedattributes":["all"],"groupby":"EntityID","scoped":true}},{"id":3,"x":-286,"y":-2,"type":"EntityStream","module":null,"values":{"selectedtype":"ChildLost","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]} 
}
];

addMenuItem('Topology', showTopologies);         
addMenuItem('Intent', showIntents);     
addMenuItem('TaskInstance', showTaskInstances);        

showTopologies();

queryOperatorList();

queryTopology();


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
    if(handler != undefined) {
        handler();        
    }
}

function initTopologyExamples() 
{
    for(var i=0; i<myToplogyExamples.length; i++) {
        var example = myToplogyExamples[i];
        submitTopology(example.topology, example.designboard);
    }
}

function queryTopology() 
{
    var queryReq = {}
    queryReq.entities = [{type:'Topology', isPattern: true}];
    client.queryContext(queryReq).then( function(topologyList) {
        if (topologyList.length == 0) {
			initTopologyExamples();
		}
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    });          
}


function showTopologyEditor() 
{
    $('#info').html('to design a service topology');

    var html = '';
    
    html += '<div id="topologySpecification" class="form-horizontal"><fieldset>';            
    
    html += '<div class="control-group"><label class="control-label">name</label>';
    html += '<div class="controls"><input type="text" class="input-large" id="serviceName">';
    html += '</div></div>';
    
    html += '<div class="control-group"><label class="control-label">description</label>';
    html += '<div class="controls"><textarea class="form-control" rows="3" id="serviceDescription"></textarea>';
    html += '</div></div>';      
           
    html += '<div class="control-group"><label class="control-label">task</label><div class="controls">';
    html += '<span>  </span><button id="cleanBoard" type="button" class="btn btn-default">Clean Board</button>';                            
    html += '<span>  </span><button id="saveBoard" type="button" class="btn btn-default">Save Board</button>';  
    html += '<span>  </span><button id="generateTopology" type="button" class="btn btn-primary">Generate Topology</button>';                                      
    html += '</div></div>';   
       
    html += '</fieldset></div>';   
        
    html += '<div id="blocks" style="width:1000px; height:400px"></div>';
    
    $('#content').html(html);    

    blocks = new Blocks();
 
    registerAllBlocks(blocks, operatorList);

    blocks.run('#blocks');
    
    blocks.types.addCompatibility('string', 'choice');
    
    if (CurrentScene != null ) {
        blocks.importData(CurrentScene);
    }
        
    blocks.ready(function() {                
        // associate functions to clickable buttons
        $('#generateTopology').click(function() {
            boardScene2Topology(blocks.export());
        });    
        $('#cleanBoard').click(function() {
            blocks.clear();
        });  
        $('#saveBoard').click(function() {
            CurrentScene = blocks.export();
        });                                              
    });    
}

function openTopologyEditor(topologyEntity)
{		
    if(topologyEntity.attributes.designboard){
        CurrentScene = topologyEntity.attributes.designboard.value;          
        showTopologyEditor(); 
        
        var topology = topologyEntity.attributes.template.value;
        
        $('#serviceName').val(topology.name);
        $('#serviceDescription').val(topology.description);
    }
}

function deleteTopology(topologyEntity)
{
    var entityid = {
        id : topologyEntity.entityId.id, 
        type: 'Topology',
        isPattern: false
    };	    
    
    clientDes.deleteContext(entityid).then( function(data) {
        console.log(data);
		updateTopologyList();		
    }).catch( function(error) {
        console.log('failed to delete a service topology');
    });  	
}



function queryOperatorList()
{
    var queryReq = {}
    queryReq.entities = [{type:'Operator', isPattern: true}];           
    
    client.queryContext(queryReq).then( function(operators) {
        for(var i=0; i<operators.length; i++){
            var entity = operators[i];        
            var operator = entity.attributes.operator.value;                 
            operatorList.push(operator.name);              
    	} 
        
        // add it into the select list        
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });    
}

function boardScene2Topology(scene)
{
    // construct a topology from the provided information
    var topologyName = $('#serviceName').val();
    var serviceDescription = $('#serviceDescription').val();
   
    var topology = {};    
    topology.name = topologyName;
    topology.description = serviceDescription;    
    topology.tasks = generateTaskList(scene);           

    // submit the generated topology
    submitTopology(topology, scene);
}


function generateTaskList(scene)
{    
    var tasklist = [];
    
    for(var i=0; i<scene.blocks.length; i++){
        var block = scene.blocks[i];
        if (block.type == 'Task') {            
            var task = {};
            
            task.name = block.values['name'];
            task.operator = block.values['operator'];

            task.input_streams = [];
            task.output_streams = [];
            
            // look for all input streams associated with this task
            task.input_streams = findInputStream(scene, block.id); 
                        
            // figure out the defined output stream types                        
            for(var j=0; j<block.values['outputs'].length; j++){
                var outputstream = {};
                outputstream.entity_type = block.values['outputs'][j];
                task.output_streams.push(outputstream);
            }
            
            tasklist.push(task);
        }
    }
    
    return tasklist;
}

function findInputStream(scene, blockid)
{
    var inputstreams = [];
    
    for(var i=0; i<scene.edges.length; i++) {
        var edge = scene.edges[i];
        if (edge.block2 == blockid) {
            var inputblockId = edge.block1;
            
            for(var j=0; j<scene.blocks.length; j++){
                var block = scene.blocks[j];
                if (block.id == inputblockId){
                    if (block.type == 'Shuffle') {                        
                        var inputstream = {};
                        
                        inputstream.selected_type = findInputType(scene,  block.id)          
                        
                        if (block.values['selectedattributes'].length == 1 && block.values['selectedattributes'][0].toUpperCase() == 'ALL') {
                            inputstream.selected_attributes = [];
                        } else {
                            inputstream.selected_attributes = block.values['selectedattributes'];                            
                        }
                        
                        inputstream.groupby = block.values['groupby'];                                                                        
                        inputstream.scoped = true;
                        
                        inputstreams.push(inputstream)
                    } else if (block.type == 'EntityStream') {
                        var inputstream = {};
                                                
                        inputstream.selected_type = block.values['selectedtype'];            
                        
                        if (block.values['selectedattributes'].length == 1 && block.values['selectedattributes'][0].toUpperCase() == 'ALL') {
                            inputstream.selected_attributes = [];
                        } else {
                            inputstream.selected_attributes = block.values['selectedattributes'];                            
                        }                                                            
                        
                        inputstream.groupby = block.values['groupby'];                                                
                        inputstream.scoped = block.values['scoped'];
                        
                        inputstreams.push(inputstream)
                    }
                }
            }
        }
    }        
    
    return inputstreams;
}

function findInputType(scene, blockId)
{
    var inputType = "unknown";

    for(var i=0; i<scene.edges.length; i++){
        var edge = scene.edges[i];
        
        if(edge.block2 == blockId) {
            var index = edge.connector1[2];     
            
            for(var j=0; j<scene.blocks.length; j++) {
                var block = scene.blocks[j];                
                if(block.id == edge.block1) {  
                    console.log(block);
                    inputType = block.values.outputs[index];                    
                }
            }               
        }
    }
    
    return inputType;
}

function submitTopology(topology, designboard)
{       
    console.log("==============test========");
    console.log(JSON.stringify(topology));
    console.log(JSON.stringify(designboard));

    var topologyCtxObj = {};
    
    topologyCtxObj.entityId = {
        id : 'Topology.' + topology.name, 
        type: 'Topology',
        isPattern: false
    };
    
    topologyCtxObj.attributes = {};   
    topologyCtxObj.attributes.status = {type: 'string', value: 'enabled'};
    topologyCtxObj.attributes.designboard = {type: 'object', value: designboard};    
    topologyCtxObj.attributes.template = {type: 'object', value: topology};  
    
    topologyCtxObj.metadata = {};    
      
    var geoScope = {};    
    geoScope.type = "global"
    geoScope.value = "global"
    topologyCtxObj.metadata.location = geoScope;    
     
    clientDes.updateContext(topologyCtxObj).then( function(data) {
        console.log(data);                
        // update the list of submitted topologies
        showTopologies();               
    }).catch( function(error) {
        console.log('failed to submit the topology');
    });    
}

function showTopologies() 
{    
    $('#info').html('list of all registered service topologies');
    
    var html = '<div style="margin-bottom: 10px;"><button id="registerTopology" type="button" class="btn btn-primary">register</button></div>';
    html += '<div id="topologyList"></div>';

	$('#content').html(html);   
      
    $( "#registerTopology" ).click(function() {
        showTopologyEditor();
    });    
                  
    // update the list of submitted topologies
    updateTopologyList();    
}

function updateTopologyList() 
{
    var queryReq = {}
    queryReq.entities = [{type:'Topology', isPattern: true}];
    client.queryContext(queryReq).then( function(topologyList) {
        displayTopologyList(topologyList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });       
}

function displayTopologyList(topologies) 
{
    if(topologies == null || topologies.length == 0) {
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>Name</th>';        
    html += '<th>Description</th>';            
    html += '<th>Tasks</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<topologies.length; i++){
        var topology = topologies[i];
		
    	html += '<tr>'; 
		html += '<td>' + topology.entityId.id;
		html += '<br><button id="editor-' + topology.entityId.id + '" type="button" class="btn btn-default">editor</button>';
		html += '<br><button id="delete-' + topology.entityId.id + '" type="button" class="btn btn-default">delete</button>';
		html += '</td>';        
               
        topology = topology.attributes.template.value;
        
		html += '<td>' + topology.name + '</td>';                                 
		html += '<td>' + topology.description + '</td>';                
		html += '<td>' + JSON.stringify(topology.tasks) + '</td>';                
        
		html += '</tr>';	
	}
       
    html += '</table>';  

	$('#topologyList').html(html);  
    
    // associate a click handler to the editor button
    for(var i=0; i<topologies.length; i++){
        var topology = topologies[i];
        
		// association handlers to the buttons
        var editorButton = document.getElementById('editor-' + topology.entityId.id);
        editorButton.onclick = function(mytopology) {
            return function(){
                openTopologyEditor(mytopology);
            };
        }(topology);
		
        var deleteButton = document.getElementById('delete-' + topology.entityId.id);
        deleteButton.onclick = function(mytopology) {
            return function(){
                deleteTopology(mytopology);
            };
        }(topology);		
	}        
}

function showIntents() 
{        
    $('#info').html('list of active intentss');
    
    var html = '<div style="margin-bottom: 10px;"><button id="addIntent" type="button" class="btn btn-primary">add</button></div>';
    html += '<div id="intentList"></div>';

	$('#content').html(html);   
      
    $( "#addIntent" ).click(function() {
        addIntent();
    });    
                  
    var queryReq = {}
    queryReq.entities = [{type:'ServiceIntent', isPattern: true}];    
    
    client.queryContext(queryReq).then( function(intents) {
        displayIntentList(intents);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query intent entities');
    });     
}


function displayIntentList(intents) 
{
    if(intents == null || intents.length == 0){
        $('#intentList').html('');           
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';        
    html += '<th>Topology</th>';    
    html += '<th>Priority</th>';    
    html += '<th>QoS</th>';
    html += '<th>GeoScope</th>';
    //html += '<th>Status</th>';
    //html += '<th>Reason</th>';
    html += '</tr></thead>';    
           
    for(var i=0; i<intents.length; i++){
        var entity = intents[i];
        
        var intent = entity.attributes.intent.value;		
        
        html += '<tr>'; 
		html += '<td>' + entity.entityId.id;
        html += '<br><button id="DELETE-' + entity.entityId.id + '" type="button" class="btn btn-primary">Delete</button>';             
        html + '</td>';                                
		html += '<td>' + intent.topology + '</td>';                        
		html += '<td>' + JSON.stringify(intent.priority)+ '</td>';                
		html += '<td>' + intent.qos + '</td>';        
		html += '<td>' + JSON.stringify(intent.geoscope) + '</td>';                
		//html += '<td>' + JSON.stringify(entity.metadata.status) + '</td>';                
		//html += '<td>' + JSON.stringify(intent.metadata.reason) + '</td>';                
                              
		html += '</tr>';	                        
	}    
       
    html += '</table>';  
    
	$('#intentList').html(html);      

    // associate a click handler to generate device profile on request
    for(var i=0; i<intents.length; i++){
        var entity = intents[i];        
        var deleteButton = document.getElementById('DELETE-' + entity.entityId.id);
        deleteButton.onclick = function(intentID) {
            return function(){
                removeIntent(intentID);
            };                        
        }(entity.entityId.id);        
	}    
}


function addIntent()
{
    $('#info').html('to specify an intent object in order to run your service');         
    
    var html = '<div id="intentRegistration" class="form-horizontal"><fieldset>';   

    html += '<div class="control-group"><label class="control-label" for="input01">Topology</label>';
    html += '<div class="controls"><select id="topologyItems"></select></div>'
    html += '</div>';       
     
    html += '<div class="control-group"><label class="control-label">Priority</label><div class="controls">';
    html += '<select id="priorityLevel"><option>low</option><option>middle</option><option>high</option></select>';    
    html += '</div></div>';    
    
    html += '<div class="control-group"><label class="control-label">Resource usage</label><div class="controls">';
    html += '<select id="resouceUsage"><option>inclusive</option><option>exclusive</option></select>';
    html += '</div></div>';     

    html += '<div class="control-group"><label class="control-label">Objective</label><div class="controls">';
    html += '<select id="QoS">';
    html += '<option value="NONE">None</option>';
    html += '<option value="MAX_THROUGHPUT">Max Throughput</option>';
    html += '<option value="MIN_LATENCY">Min Latency</option>';
    html += '<option value="MIN_COST">Min Cost</option></select>';    
    html += '</div></div>';  

    html += '<div class="control-group"><label class="control-label">Geoscope</label><div class="controls">';
    html += '<select id="geoscope"><option value="global">global</option><option value="custom">custom</option></select>';    
    html += '</div></div>';  
    
    html += '<div class="control-group" id="mapDiv">'
    html += '</div>';    
        
    html += '<div class="control-group"><label class="control-label" for="input01"></label>';
    html += '<div class="controls"><button id="submitIntent" type="button" class="btn btn-primary">Send</button>';
    html += '</div></div>';   
       
    html += '</fieldset></div>';
    
    $('#content').html(html);        
       
    // add all service topologies into the selection list
    listAllTopologies();
    
    // associate functions to clickable buttons
    $('#submitIntent').click(submitIntent);  
    
    $('#geoscope').change( function() {
        var scope = $(this).val();
        
        if (scope == "custom") {
            // show the map to set locations
            showMap();               
        } else {
            removeMap();
        }
    });
}


function removeIntent(entityID)
{
    var entityid = {
        id : entityID, 
        isPattern: false
    };	    
    
    clientDes.deleteContext(entityid).then( function(data) {
        console.log('remove the service intent');		
        // show the updated intent list
        showIntents();				
    }).catch( function(error) {
        console.log('failed to delete this service intent');
    }); 
}


function listAllTopologies() 
{
    var queryReq = {}
    queryReq.entities = [{type:'Topology', isPattern: true}];
    client.queryContext(queryReq).then( function(topologyList) {
        var topologySelect = document.getElementById('topologyItems');
        for(var i=0; i<topologyList.length; i++) {
            var topology = topologyList[i].attributes.template.value;
            topologySelect.options[topologySelect.options.length] = new Option(topology.name, topology.name);
        }
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query topology');
    });          
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
  

function submitIntent()
{       
    var intent = {};    
    
    var topology = $('#topologyItems option:selected').val();    
    intent.topology = topology;
        
    var temp1 = $('#PriorityLevel option:selected').val();
    var priorityLevel = 0;
    switch(temp1) {
        case 'low':
            priorityLevel = 0;
            break;
        case 'middle': 
            priorityLevel = 50;
            break;
        case 'high':            
            priorityLevel = 100;
            break;        
    }    

    var temp2 = $('#ResouceUsage option:selected').val();
    var exclusiveResourceUsage = false;    
    if(temp2 == 'exclusive'){
        exclusiveResourceUsage = true; 
    }    
    
    intent.priority = {
        'exclusive': exclusiveResourceUsage,
        'level': priorityLevel
    };
    
    var qos = $('#QoS option:selected').val();    
    intent.qos = qos;    
        
    var scope = $('#geoscope option:selected').val(); 

    var operationScope = {};    
    if (scope == 'custom') {        
        operationScope.scopeType = geoscope.type
        operationScope.scopeValue = geoscope.value    
        intent.geoscope = operationScope;   
    } 
    
    var intentCtxObj = {};
    
    var uid = uuid();
    
    intentCtxObj.entityId = { 
        id: 'ServiceIntent.' + uid,           
        type: 'ServiceIntent',
        isPattern: false
    };
    
    intentCtxObj.attributes = {};   
    intentCtxObj.attributes.status = {type: 'string', value: 'enabled'};
    intentCtxObj.attributes.intent = {type: 'object', value: intent};  
    
    intentCtxObj.metadata = {};         
    var geoScope = {};    
    if (scope == 'custom') {        
        geoScope.type = geoscope.type
        geoScope.value = geoscope.value    
    } else {
        geoScope.type = scope
        geoScope.value = scope
    }    
    intentCtxObj.metadata.location = geoScope;         
    
    console.log(JSON.stringify(intentCtxObj));
        
    clientDes.updateContext(intentCtxObj).then( function(data) {
        console.log(data);  
        // update the list of submitted intents
        showIntents();               
    }).catch( function(error) {
        console.log('failed to submit the defined intent');
    });    
}


function showMap() 
{   
    var htmlContent = '<label class="control-label" for="input01">Polygon</label><div class="controls"><div id="map"  style="width: 500px; height: 400px"></div></div>';
    $('#mapDiv').html(htmlContent);
    
    var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        osm = L.tileLayer(osmUrl, {maxZoom: 7, zoom: 7}),
        map = new L.Map('map', {zoomControl:false, layers: [osm], center: new L.LatLng(35.692221, 138.709059), zoom: 7});

	//disable zoom in/out
	map.dragging.disable();
	map.touchZoom.disable();
	map.doubleClickZoom.disable();
	map.scrollWheelZoom.disable();
	map.boxZoom.disable();
	map.keyboard.disable();

    var drawnItems = new L.FeatureGroup();
    map.addLayer(drawnItems);

    var drawControl = new L.Control.Draw({
        draw: {
            position: 'topleft',
            polyline: false,            
			polygon: {
                showArea: false
            },
            rectangle: {
                shapeOptions: {
                    color: '#E3225C',
                    weight: 2,
                    clickable: false
                }
            },
            circle: {
                shapeOptions: {
                    color: '#E3225C',
                    weight: 2,
                    clickable: false
                }
            },          
            marker: false
        },
        edit: {
            featureGroup: drawnItems
        }
    });
    map.addControl(drawControl);

    map.on('draw:created', function (e) {
		var type = e.layerType;
		var layer = e.layer;

		if (type === 'rectangle') {
            var geometry = layer.toGeoJSON()['geometry'];
            console.log(geometry);
            
            geoscope.type = 'polygon';
            geoscope.value = {
                vertices: []
            };
            
            var points = geometry.coordinates[0];
            for(var i in points){
                geoscope.value.vertices.push({longitude: points[i][0], latitude: points[i][1]});
            }
            
			console.log(geoscope);            
		} 
		if (type === 'circle') {
            var geometry = layer.toGeoJSON()['geometry'];
            console.log(geometry);
            var radius = layer.getRadius();
            
            geoscope.type = 'circle';
            geoscope.value = {
                centerLatitude: geometry.coordinates[1],
                centerLongitude: geometry.coordinates[0],
                radius: radius
            }
            
			console.log(geoscope);            
		} 
		if (type === 'polygon') {
            var geometry = layer.toGeoJSON()['geometry'];
            console.log(geometry);
            
            geoscope.type = 'polygon';
            geoscope.value = {
                vertices: []
            };
            
            var points = geometry.coordinates[0];
            for(var i in points){
                geoscope.value.vertices.push({longitude: points[i][0], latitude: points[i][1]});
            }
            
			console.log(geoscope);            
		}
        
        drawnItems.addLayer(layer);
    });  
}

function removeMap() {
    $('#mapDiv').html('');    
}


function showTaskInstances() 
{
    $('#info').html('list of running data processing tasks');

    var queryReq = {}
    queryReq.entities = [{type:'Task', isPattern: true}];    
    
    client.queryContext(queryReq).then( function(taskList) {
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
    html += '<th>Service</th>';
    html += '<th>Task</th>';    
    html += '<th>Worker</th>';    
    html += '<th>port</th>';	
    html += '<th>status</th>';		    
    html += '</tr></thead>';    
       
    for(var i=0; i<tasks.length; i++){
        var task = tasks[i];
		
        html += '<tr>'; 
		html += '<td>' + task.entityId.id + '</td>';
		html += '<td>' + task.entityId.type + '</td>'; 
        html += '<td>' + task.attributes.service.value + '</td>';		
        html += '<td>' + task.attributes.task.value + '</td>';	        
	    html += '<td>' + task.metadata.worker.value + '</td>';

		html += '<td>' + task.attributes.port.value + '</td>';
				
		if (task.attributes.status.value == "paused") {
			html += '<td><font color="red">' + task.attributes.status.value + '</font></td>';			
		} else {
			html += '<td><font color="green">' + task.attributes.status.value + '</font></td>';
		}
        
		html += '</tr>';	
	}
       
    html += '</table>'; 

	$('#content').html(html);   
}


});



