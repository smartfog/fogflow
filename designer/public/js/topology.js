'use strict';

$(function(){

// initialize the menu bar
var handlers = {}
var CurrentScene = null;

// location of new device
var locationOfNewDevice = null;
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

var myToplogyExamples = [
{
    topology: {"name":"anomaly-detection","description":"detect anomaly events from time series data points","priority":{"exclusive":false,"level":50},"trigger":"on-demand","tasks":[{"name":"AnomalyDetector","operator":"anomaly","groupBy":"shop","input_streams":[{"type":"PowerPanel","scoped":true,"shuffling":"unicast"},{"type":"Rule","scoped":false,"shuffling":"broadcast"}],"output_streams":[{"type":"Anomaly"}]},{"name":"Counter","operator":"counter","groupBy":"all","input_streams":[{"type":"Anomaly","scoped":true,"shuffling":"unicast"}],"output_streams":[{"type":"Stat"}]}]},
    designboard: {"edges":[{"id":1,"block1":1,"connector1":["outputs","output",0],"block2":2,"connector2":["inputs","input",0]},{"id":2,"block1":4,"connector1":["stream","output"],"block2":1,"connector2":["inputs","input",1]},{"id":3,"block1":3,"connector1":["stream","output"],"block2":1,"connector2":["inputs","input",0]}],"blocks":[{"id":1,"x":-21,"y":-95,"type":"Task","module":null,"values":{"name":"AnomalyDetector","operator":"anomaly","groupby":"shop","inputs":["unicast","broadcast"],"outputs":["Anomaly"]}},{"id":2,"x":194,"y":-97,"type":"Task","module":null,"values":{"name":"Counter","operator":"counter","groupby":"all","inputs":["unicast"],"outputs":["Stat"]}},{"id":3,"x":-280,"y":-138,"type":"InputStream","module":null,"values":{"entitytype":"PowerPanel","scoped":true}},{"id":4,"x":-279,"y":24,"type":"InputStream","module":null,"values":{"entitytype":"Rule","scoped":false}}]}
}, {
    topology: {"name":"crowd-detection","description":"detect the number of faces from IP cameras","priority":{"exclusive":false,"level":50},"trigger":"on-demand","tasks":[{"name":"facesum","operator":"sum","groupBy":"all","input_streams":[{"type":"FaceNum","scoped":true,"shuffling":"unicast"}],"output_streams":[{"type":"FaceSum"}]},{"name":"facenum","operator":"facecounter","groupBy":"cameraID","input_streams":[{"type":"Camera","scoped":true,"shuffling":"unicast"}],"output_streams":[{"type":"FaceNum"}]}]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["outputs","output",0],"block2":1,"connector2":["inputs","input",0]},{"id":2,"block1":3,"connector1":["stream","output"],"block2":2,"connector2":["inputs","input",0]}],"blocks":[{"id":1,"x":241,"y":-113,"type":"Task","module":null,"values":{"name":"facesum","operator":"sum","groupby":"all","inputs":["unicast"],"outputs":["FaceSum"]}},{"id":2,"x":-70,"y":-164,"type":"Task","module":null,"values":{"name":"facenum","operator":"facecounter","groupby":"cameraID","inputs":["unicast"],"outputs":["FaceNum"]}},{"id":3,"x":-367,"y":-100,"type":"InputStream","module":null,"values":{"entitytype":"Camera","scoped":true}}]}
}, {
    topology: {"name":"child-finder","description":"search for a lost child based on face recognition","priority":{"exclusive":true,"level":100},"trigger":"on-demand","tasks":[{"name":"childfinder","operator":"facefinder","groupBy":"cameraID","input_streams":[{"type":"Camera","scoped":true,"shuffling":"unicast"},{"type":"ChildLost","scoped":false,"shuffling":"broadcast"}],"output_streams":[{"type":"ChildFound"}]}]},
    designboard: {"edges":[{"id":1,"block1":3,"connector1":["stream","output"],"block2":1,"connector2":["inputs","input",1]},{"id":2,"block1":2,"connector1":["stream","output"],"block2":1,"connector2":["inputs","input",0]}],"blocks":[{"id":1,"x":-48,"y":-113,"type":"Task","module":null,"values":{"name":"childfinder","operator":"facefinder","groupby":"cameraID","inputs":["unicast","broadcast"],"outputs":["ChildFound"]}},{"id":2,"x":-344,"y":-159,"type":"InputStream","module":null,"values":{"entitytype":"Camera","scoped":true}},{"id":3,"x":-336,"y":3,"type":"InputStream","module":null,"values":{"entitytype":"ChildLost","scoped":false}}]}
}
];


addMenuItem('Topology', showTopologies);         
addMenuItem('Requirement', showRequirements);         
addMenuItem('Editor', showEditor);         

showEditor();

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


function showEditor() 
{
    $('#info').html('to design a service topology');

    var html = '';
    
    html += '<div id="topologySpecification" class="form-horizontal"><fieldset>';            
    
    html += '<div class="control-group"><label class="control-label">topology name</label>';
    html += '<div class="controls"><input type="text" class="input-large" id="topologyName">';
    html += '</div></div>';
    
    html += '<div class="control-group"><label class="control-label">service description</label>';
    html += '<div class="controls"><textarea class="form-control" rows="3" id="serviceDescription"></textarea>';
    html += '</div></div>';    
    
    html += '<div class="control-group"><label class="control-label">priority</label><div class="controls">';
    html += '<select id="priorityLevel"><option>low</option><option>middle</option><option>high</option></select>';    
    html += '</div></div>';    
    
    html += '<div class="control-group"><label class="control-label">resource usage</label><div class="controls">';
    html += '<select id="resouceUsage"><option>inclusive</option><option>exclusive</option></select>';
    html += '</div></div>';        
           
    html += '<div class="control-group"><label class="control-label">data processing graph</label><div class="controls">';
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

function openEditor(topologyEntity)
{
	//selectMenuItem('Editor');
	//window.location.hash = '#Editor';
		
    if(topologyEntity.attributes.designboard){
        CurrentScene = topologyEntity.attributes.designboard.value;          
        showEditor(); 
        
        var topology = topologyEntity.attributes.template.value;
        
        $('#topologyName').val(topology.name);
        $('#serviceDescription').val(topology.description);
		
		var priority = topology.priority;
		
		if (priority.level == 0) {
			$('#priorityLevel').val('low');
		} else if (priority.level == 50) {
			$('#priorityLevel').val('middle');
		} else if (priority.level == 100) {
			$('#priorityLevel').val('high');
		}
		
		if (priority.exclusive == true) {
			$('#resouceUsage').val('exclusive');
		} else {
			$('#resouceUsage').val('inclusive');
		}
    }
}

function deleteTopology(topologyEntity)
{
    var entityid = {
        id : topologyEntity.entityId.id, 
        type: 'Topology',
        isPattern: false
    };	    
    
    client.deleteContext(entityid).then( function(data) {
        console.log(data);
		updateTopologyList();		
    }).catch( function(error) {
        console.log('failed to delete a service topology');
    });  	
}

function queryOperatorList()
{
    var queryReq = {}
    queryReq.entities = [{type:'DockerImage', isPattern: true}];           
    
    client.queryContext(queryReq).then( function(imageList) {
        console.log(imageList);

        for(var i=0; i<imageList.length; i++){
            var dockerImage = imageList[i];            
            var operatorName = dockerImage.attributes.operator.value;
            
            var exist = false;
            for(var j=0; j<operatorList.length; j++){
                if(operatorList[j] == operatorName){
                    exist = true;
                    break;
                }
            }
            
            if(exist == false){
                operatorList.push(operatorName);                
            }            
        }
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query the operator list');
    });     
}

function boardScene2Topology(scene)
{
    // construct a topology from the provided information
    var topologyName = $('#topologyName').val();
    var serviceDescription = $('#serviceDescription').val();

    var temp1 = $('#priorityLevel option:selected').val();
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
    
    var temp2 = $('#resouceUsage option:selected').val();
    var exclusiveResourceUsage = false;    
    if(temp2 == 'exclusive'){
        exclusiveResourceUsage = true; 
    }    
    
    var topology = {};    
    topology.name = topologyName;
    topology.description = serviceDescription;
    topology.priority = {
        'exclusive': exclusiveResourceUsage,
        'level': priorityLevel
    };
    
    topology.trigger = 'on-demand';
    
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
            task.groupBy = block.values['groupby'];
            task.input_streams = [];
            task.output_streams = [];
            
            for(var j=0; j<block.values['inputs'].length; j++){
                var inputstream = findInputStream(scene, block.id, j);  
                
                if( inputstream != null ) {
                    inputstream.shuffling = block.values['inputs'][j];                
                    task.input_streams.push(inputstream);                    
                }                              
            }
                        
            for(var j=0; j<block.values['outputs'].length; j++){
                var outputstream = {};
                outputstream.type = block.values['outputs'][j];
                task.output_streams.push(outputstream);
            }
            
            tasklist.push(task);
        }
    }
    
    return tasklist;
}

function findInputStream(scene, blockid, inputIdx)
{
    for(var i=0; i<scene.edges.length; i++) {
        var edge = scene.edges[i];
        if (edge.block2 == blockid && edge.connector2[2] == inputIdx) {
            var inputblockId = edge.block1;
            
            for(var j=0; j<scene.blocks.length; j++){
                var block = scene.blocks[j];
                if (block.id == inputblockId){
                    if (block.type == 'Task') {
                        var inputstream = {};
                        
                        var attributeName = edge.connector1[0];
                        var valueIndx = edge.connector1[2];
                        inputstream.type = block.values[attributeName][valueIndx];
                        inputstream.scoped = true;
                        
                        return inputstream;                        
                    } else if (block.type == 'InputStream') {
                        var inputstream = {};
                        
                        inputstream.type = block.values['entitytype'];
                        inputstream.scoped = block.values['scoped'];
                        
                        return inputstream;                                                
                    }
                }
            }
        }
    }        
    
    return null;
}

function submitTopology(topology, designboard)
{       
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
        
    client.updateContext(topologyCtxObj).then( function(data) {
        console.log(data);  
              
        // update the list of submitted topologies
        showTopologies();               
    }).catch( function(error) {
        console.log('failed to submit the topology');
    });    
}

function showTopologies() 
{    
    $('#info').html('list of submitted topologies');

    var html = '<div id="itemList"></div>';         
    
	$('#content').html(html);       
                  
    // update the list of submitted topologies
    updateTopologyList();    
}

function updateTopologyList() 
{
    var queryReq = {}
    queryReq.entities = [{type:'Topology', isPattern: true}];
    client.queryContext(queryReq).then( function(topologyList) {
        console.log(topologyList);
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
    html += '<th>Type</th>';
    html += '<th>Status</th>';    
    html += '<th>Template</th>';    
    html += '</tr></thead>';    
       
    for(var i=0; i<topologies.length; i++){
        var topology = topologies[i];
		
    		html += '<tr>'; 
		html += '<td>' + topology.entityId.id;
		html += '<br><button id="editor-' + topology.entityId.id + '" type="button" class="btn btn-default">editor</button>';
		html += '<br><button id="delete-' + topology.entityId.id + '" type="button" class="btn btn-default">delete</button>';
		html += '</td>';        
		html += '<td>' + topology.entityId.type + '</td>'; 
                        
		html += '<td>' + topology.attributes.status.value + '</td>';        
		html += '<td>' + JSON.stringify(topology.attributes.template.value) + '</td>';
		html += '</tr>';	
	}
       
    html += '</table>';  

	$('#itemList').html(html);  
    
    // associate a click handler to the editor button
    for(var i=0; i<topologies.length; i++){
        var topology = topologies[i];
        
		// association handlers to the buttons
        var editorButton = document.getElementById('editor-' + topology.entityId.id);
        editorButton.onclick = function(mytopology) {
            return function(){
                openEditor(mytopology);
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

function showRequirements() 
{        
    $('#info').html('list of scoped requirements to trigger processing tasks');

    var queryReq = {}
    queryReq.entities = [{type:'Requirement', isPattern: true}];    
    
    client.queryContext(queryReq).then( function(requirements) {
        console.log(requirements);
        displayRequirementList(requirements);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });     
}

function displayRequirementList(requirements) 
{
    if(requirements == null || requirements.length ==0){
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
       
    for(var i=0; i<requirements.length; i++){
        var requirement = requirements[i];
		
        html += '<tr>'; 
		html += '<td>' + requirement.entityId.id + '</td>';
		html += '<td>' + requirement.entityId.type + '</td>'; 
		html += '<td>' + JSON.stringify(requirement.attributes) + '</td>';        
		html += '<td>' + JSON.stringify(requirement.metadata) + '</td>';
		html += '</tr>';	
	}
       
    html += '</table>'; 

	$('#content').html(html);   
}

});



