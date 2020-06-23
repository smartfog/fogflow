'use strict';

$(function(){

// initialize the menu bar
var handlers = {}

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
// to interact with designer
var clientDes = new NGSIDesClient(config.designerIP+':'+config.webSrvPort);

var myFogFunctionExamples = [
{
    name: "Convert1",
    topology: {"name":"Convert1","description":"test","tasks":[{"name":"Main","operator":"converter","input_streams":[{"selected_type":"RainSensor","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"RainObservation"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"converter","outputs":["RainObservation"]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"RainSensor","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]},
    intent: {"topology":"Convert1","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Convert2",
    topology: {"name":"Convert2","description":"test","tasks":[{"name":"Main","operator":"geohash","input_streams":[{"selected_type":"SmartAwning","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"geohash","outputs":[]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"SmartAwning","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]},
    intent: {"topology":"Convert2","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Convert3",
    topology: {"name":"Convert3","description":"test","tasks":[{"name":"Main","operator":"converter","input_streams":[{"selected_type":"ConnectedCar","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"RainObservation"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"converter","outputs":["RainObservation"]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"ConnectedCar","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]},
    intent: {"topology":"Convert3","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Prediction",
    topology: {"name":"Prediction","description":"test","tasks":[{"name":"Main","operator":"predictor","input_streams":[{"selected_type":"RainObservation","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"Prediction"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"predictor","outputs":[""]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"RainObservation","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]},
    intent: {"topology":"Prediction","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Prediction",
    topology: {"name":"Prediction","description":"test","tasks":[{"name":"Main","operator":"predictor","input_streams":[{"selected_type":"RainObservation","selected_attributes":[],"groupby":"ALL","scoped":false}],"output_streams":[{"entity_type":"Prediction"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"predictor","outputs":[""]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"RainObservation","selectedattributes":["all"],"groupby":"ALL","scoped":false}}]},
    intent: {"topology":"Prediction","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Controller",
    topology: {"name":"Controller","description":"test","tasks":[{"name":"Main","operator":"controller","input_streams":[{"selected_type":"SmartAwning","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"ControlAction"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"controller","outputs":["ControlAction"]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"SmartAwning","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"Controller","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},{
    name: "Detector",
    topology: {"name":"Detector","description":"test","tasks":[{"name":"Main","operator":"detector","input_streams":[{"selected_type":"Camera","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"detector","outputs":[]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"Camera","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"Detector","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},
{
    name: "Test",
    topology: {"name":"Test","description":"just for a simple test","tasks":[{"name":"Main","operator":"dummy","input_streams":[{"selected_type":"Temperature","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"dummy","outputs":["Out"]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"Temperature","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"Test","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},
/*{
    name: "Agent",
    topology: {"name":"Agent","description":"just for a simple test","tasks":[{"name":"Main","operator":"iotagent","input_streams":[{"selected_type":"Worker","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":123,"y":-99,"type":"Task","module":null,"values":{"name":"Main","operator":"iotagent","outputs":["Out"]}},{"id":2,"x":-194,"y":-97,"type":"EntityStream","module":null,"values":{"selectedtype":"Worker","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"Agent","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
},*/{
    name: "PrivateSiteEstimation",
    topology: {"name":"PrivateSiteEstimation","description":"to estimate the free parking lots from a private parking site","tasks":[{"name":"Estimation","operator":"privatesite","input_streams":[{"selected_type":"PrivateSite","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":26,"y":-47,"type":"Task","module":null,"values":{"name":"Estimation","operator":"privatesite","outputs":["Out"]}},{"id":2,"x":-302,"y":-87,"type":"EntityStream","module":null,"values":{"selectedtype":"PrivateSite","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"PrivateSiteEstimation","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
}, {
    name: "PublicSiteEstimation",
    topology: {"name":"PublicSiteEstimation","description":"to estimate the free parking lot from a public parking site","tasks":[{"name":"PubFreeLotEstimation","operator":"publicsite","input_streams":[{"selected_type":"PublicSite","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]} ,
    designboard:  {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":-37,"y":-108,"type":"Task","module":null,"values":{"name":"PubFreeLotEstimation","operator":"publicsite","outputs":["Out"]}},{"id":2,"x":-340,"y":-128,"type":"EntityStream","module":null,"values":{"selectedtype":"PublicSite","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"PublicSiteEstimation","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
}, {
    name: "ArrivalTimeEstimation",
    topology:  {"name":"ArrivalTimeEstimation","description":"to estimate when the car will arrive at the destination","tasks":[{"name":"CalculateArrivalTime","operator":"connectedcar","input_streams":[{"selected_type":"ConnectedCar","selected_attributes":[],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":-106,"y":-93,"type":"Task","module":null,"values":{"name":"CalculateArrivalTime","operator":"connectedcar","outputs":["Out"]}},{"id":2,"x":-420,"y":-145,"type":"EntityStream","module":null,"values":{"selectedtype":"ConnectedCar","selectedattributes":["all"],"groupby":"EntityID","scoped":false}}]} ,
    intent: {"topology":"ArrivalTimeEstimation","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
}, {
    name: "ParkingLotRecommendation",
    topology: {"name":"ParkingLotRecommendation","description":"to recommend where to park around the destination","tasks":[{"name":"WhereToParking","operator":"recommender","input_streams":[{"selected_type":"ConnectedCar","selected_attributes":["ParkingRequest"],"groupby":"EntityID","scoped":false}],"output_streams":[{"entity_type":"Out"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["streams","input"]}],"blocks":[{"id":1,"x":-14,"y":-46,"type":"Task","module":null,"values":{"name":"WhereToParking","operator":"recommender","outputs":["Out"]}},{"id":2,"x":-379,"y":-110,"type":"EntityStream","module":null,"values":{"selectedtype":"ConnectedCar","selectedattributes":["ParkingRequest"],"groupby":"EntityID","scoped":false}}]},
    intent: {"topology":"ParkingLotRecommendation","priority":{"exclusive":false,"level":0},"qos":"Max Throughput","geoscope":{"scopeType":"global","scopeValue":"global"}}
} 
];


addMenuItem('FogFunction', showFogFunctions);         
addMenuItem('TaskInstance', showTaskInstances);        

showFogFunctions();

queryOperatorList();

queryFogFunctions();


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

function initFogFunctionExamples() 
{
    for(var i=0; i<myFogFunctionExamples.length; i++) {
        var fogfunction = myFogFunctionExamples[i];      
        
        var functionCtxObj = {};    
        functionCtxObj.entityId = {
            id : 'FogFunction.' + fogfunction.name, 
            type: 'FogFunction',
            isPattern: false
        };    
        functionCtxObj.attributes = {};   
        functionCtxObj.attributes.name = {type: 'string', value: fogfunction.name};    
        functionCtxObj.attributes.topology = {type: 'object', value: fogfunction.topology};   
        functionCtxObj.attributes.designboard = {type: 'object', value: fogfunction.designboard};            
        functionCtxObj.attributes.intent = {type: 'object', value: fogfunction.intent};  
        functionCtxObj.attributes.status = {type: 'string', value: 'enabled'};         
          
        submitFogFunction(functionCtxObj);
    }
}

function queryFogFunctions() 
{
    var queryReq = {}
    queryReq.entities = [{type:'FogFunction', isPattern: true}];
    client.queryContext(queryReq).then( function(fogFunctionList) {
        if (fogFunctionList.length == 0) {
			initFogFunctionExamples();
		}
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query fog functions');
    });          
}


function showFogFunctionEditor() 
{
    $('#info').html('to design a fog function');

    var html = '';
    
    html += '<div id="topologySpecification" class="form-horizontal"><fieldset>';            
    
    html += '<div class="control-group"><label class="control-label">name</label>';
    html += '<div class="controls"><input type="text" class="input-large" id="serviceName">';
    html += '</div></div>';
    
    html += '<div class="control-group"><label class="control-label">description</label>';
    html += '<div class="controls"><textarea class="form-control" rows="3" id="serviceDescription"></textarea>';
    html += '</div></div>';      
           
    html += '<div class="control-group"><label class="control-label">topology</label><div class="controls">';
    html += '<span>  </span><button id="cleanBoard" type="button" class="btn btn-default">Clean Board</button>';                            
    html += '<span>  </span><button id="saveBoard" type="button" class="btn btn-default">Save Board</button>';  
    html += '<span>  </span><button id="generateFunction" type="button" class="btn btn-primary">Submit</button>';                                      
    html += '</div></div>';   
       
    html += '</fieldset></div>';   
        
    html += '<div id="blocks" style="width:800px; height:400px"></div>';
       
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
        $('#generateFunction').click(function() {
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

function openFogFunctionEditor(fogfunction)
{
    console.log(fogfunction);
      
    if(fogfunction &&  fogfunction.attributes.designboard){
        CurrentScene = fogfunction.attributes.designboard.value;          
        showFogFunctionEditor(); 
        
        var topology = fogfunction.attributes.topology.value;        
        $('#serviceName').val(topology.name);
        $('#serviceDescription').val(topology.description);
    }
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
    // step 1: construct the service topology object       
    var topologyName = $('#serviceName').val();
    var serviceDescription = $('#serviceDescription').val();

    var topology = {};    
    topology.name = topologyName;
    topology.description = serviceDescription;    
    topology.tasks = generateTaskList(scene);           

    // step 2: construct an intent object
    var intent = {};        
    intent.topology = topologyName;    
    intent.priority = {
        'exclusive': false,
        'level': 0
    };        
    intent.qos = "default";    
    intent.geoscope = {
        "scopeType": "global",
        "scopeValue": "global"
    };   
    
    // step 3: create this fog function            
    var functionCtxObj = {};    
    functionCtxObj.entityId = {
        id : 'FogFunction.' + topologyName, 
        type: 'FogFunction',
        isPattern: false
    };    
    functionCtxObj.attributes = {};   
    functionCtxObj.attributes.name = {type: 'string', value: topologyName};    
    functionCtxObj.attributes.topology = {type: 'object', value: topology};    
    functionCtxObj.attributes.designboard = {type: 'object', value: scene};        
    functionCtxObj.attributes.intent = {type: 'object', value: intent};  
    functionCtxObj.attributes.status = {type: 'string', value: 'enabled'};   
    
    functionCtxObj.metadata = {};      
    var geoScope = {};    
    geoScope.type = "global"
    geoScope.value = "global"
    functionCtxObj.metadata.location = geoScope;     
    
    console.log("=============submit a fog function=============");
    console.log(JSON.stringify(functionCtxObj));    
    
    return clientDes.updateContext(functionCtxObj).then( function(data1) {
        console.log(data1);               
        showFogFunctions();
    }).catch( function(error) {
        console.log('failed to record the created fog function');
    });    
}

function submitFogFunction(functionCtxObj)
{
    console.log("=============submit a fog function=============");
    console.log(JSON.stringify(functionCtxObj));  

    functionCtxObj.metadata = {};      
    var geoScope = {};    
    geoScope.type = "global"
    geoScope.value = "global"
    functionCtxObj.metadata.location = geoScope;   
    
    return clientDes.updateContext(functionCtxObj).then( function(data1) {
        console.log(data1);                 
    }).catch( function(error) {
        console.log('failed to record the created fog function');
    });                  
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

function showFogFunctions() 
{    
    $('#info').html('list of all registered fog functions');
    
    var html = '<div style="margin-bottom: 10px;"><button id="registerFunction" type="button" class="btn btn-primary">register</button></div>';
    html += '<div id="functionList"></div>';

	$('#content').html(html);   
      
    $( "#registerFunction" ).click(function() {
        showFogFunctionEditor();
    });    
                  
    // update the list of submitted fog functions
    updateFogFunctionList();    
}

function updateFogFunctionList() 
{
    var queryReq = {}
    queryReq.entities = [{type:'FogFunction', isPattern: true}];
    client.queryContext(queryReq).then( function(functionList) {
        displayFunctionList(functionList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });       
}

function displayFunctionList(fogFunctions) 
{
    if(fogFunctions == null || fogFunctions.length == 0) {
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>Name</th>';        
    html += '<th>Topology</th>';            
    html += '<th>Intent</th>';                
    html += '</tr></thead>';    
       
    for(var i=0; i<fogFunctions.length; i++){
        var fogfunction = fogFunctions[i];
		
    	html += '<tr>'; 
		html += '<td>' + fogfunction.entityId.id;
		html += '<br><button id="editor-' + fogfunction.entityId.id + '" type="button" class="btn btn-default">editor</button>';
		html += '<br><button id="delete-' + fogfunction.entityId.id + '" type="button" class="btn btn-default">delete</button>';
		html += '</td>';        
                       
		html += '<td>' + fogfunction.attributes.name.value + '</td>';    
        
		html += '<td>' + JSON.stringify(fogfunction.attributes.topology.value) + '</td>'; 
                       
		html += '<td>' + JSON.stringify(fogfunction.attributes.intent.value) + '</td>';                
        
		html += '</tr>';	
	}
       
    html += '</table>';  

	$('#functionList').html(html);  
    
    // associate a click handler to the editor button
    for(var i=0; i<fogFunctions.length; i++){
        var fogfunction = fogFunctions[i];
        
		// association handlers to the buttons
        var editorButton = document.getElementById('editor-' + fogfunction.entityId.id);
        editorButton.onclick = function(myFogFunction) {
            return function(){
                openFogFunctionEditor(myFogFunction);
            };
        }(fogfunction);
		
        var deleteButton = document.getElementById('delete-' + fogfunction.entityId.id);
        deleteButton.onclick = function(myFogFunction) {
            return function(){
                deleteFogFunction(myFogFunction);
            };
        }(fogfunction);		
	}        
}

function deleteFogFunction(fogfunction)
{   
    // delete this fog function
    var functionEntity = {
        id : fogfunction.entityId.id, 
        type: 'FogFunction',
        isPattern: false
    };	    
    
    console.log("delete a fog function");
    console.log(functionEntity);
    
    clientDes.deleteContext(functionEntity).then( function(data) {
        console.log(data);
		showFogFunctions();		
    }).catch( function(error) {
        console.log('failed to delete a fog function');
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
    html += '<th>Service</th>';
    html += '<th>Task</th>';      
    html += '<th>Type</th>';
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



