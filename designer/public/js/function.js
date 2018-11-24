$(function(){

// initialization  
var handlers = {};
var CurrentScene = null;

var template = {};


$.ajax({url:'js/function_template.txt', 
    success: function(fileContent){
        template.javascript =  fileContent;      
    }
});

template.python = 'define handler(context): \r\n \t # write your own code to do data processing  \r\n \t # return the generated context entities to be published as outputs';


var myFogFunctionExamples = [
{
    "fogfunction":{"type":"docker","code":"","dockerImage":"privatesite","name":"PrivateSite","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"PrivateSite"}]}],"outputAnnotators":[]},
    "designboard":{"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":24.549998298828086,"y":-148.75000475292967,"type":"FogFunction","module":null,"values":{"name":"PrivateSite","user":"fogflow"}},{"id":2,"x":-197.4500017011719,"y":-146.75000475292967,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-428.4500017011719,"y":-145.08333299999998,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"PrivateSite"}}]}
},
{
    "fogfunction":{"type":"docker","code":"","dockerImage":"publicsite","name":"PublicSite","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"PublicSite"}]}],"outputAnnotators":[]},
    "designboard":{"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":24.549998298828086,"y":-148.75000475292967,"type":"FogFunction","module":null,"values":{"name":"PublicSite","user":"fogflow"}},{"id":2,"x":-197.4500017011719,"y":-146.75000475292967,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-428.4500017011719,"y":-145.08333299999998,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"PublicSite"}}]}
},
{
    "fogfunction":{"type":"docker","code":"","dockerImage":"recommender","name":"Recommender","user":"fogflow","inputTriggers":[{"name":"selector3","selectedAttributeList":["ParkingRequest"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"ConnectedCar"}]}],"outputAnnotators":[]},
    "designboard":{"edges":[{"id":2,"block1":3,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":3,"block1":4,"connector1":["condition","output"],"block2":3,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":104.54999829882809,"y":-135.75000475292967,"type":"FogFunction","module":null,"values":{"name":"Recommender","user":"fogflow"}},{"id":4,"x":-445.4166459882813,"y":-141.75000475292967,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"ConnectedCar"}},{"id":3,"x":-179.4166459882813,"y":-147.75000475292967,"type":"InputTrigger","module":null,"values":{"selectedattributes":["ParkingRequest"],"groupby":["id"]}}]}
},
{
    "fogfunction":{"type":"docker","code":"","dockerImage":"connectedcar","name":"ConnectedCar","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"ConnectedCar"}]}],"outputAnnotators":[]},
    "designboard":{"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":70.4081801170098,"y":-124.33545929838425,"type":"FogFunction","module":null,"values":{"name":"ConnectedCar","user":"fogflow"}},{"id":2,"x":-170.0545471557174,"y":-124.36545929838422,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-407.87272897389914,"y":-123.54727748020238,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"ConnectedCar"}}]}
},
{
    fogfunction: {"type":"docker","code":"","dockerImage":"pushbutton","name":"Pushbutton","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"Pushbutton"}]}],"outputAnnotators":[]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":96,"y":-133,"type":"FogFunction","module":null,"values":{"name":"Pushbutton","user":"fogflow"}},{"id":2,"x":-141,"y":-134,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-373,"y":-136,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"Pushbutton"}}]}
},
{
    fogfunction: {"type":"docker","code":"","dockerImage":"acoustic","name":"Acoustic","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"Microphone"}]}],"outputAnnotators":[]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":96,"y":-133,"type":"FogFunction","module":null,"values":{"name":"Acoustic","user":"fogflow"}},{"id":2,"x":-141,"y":-134,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-373,"y":-136,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"Microphone"}}]}
},
{
    fogfunction: {"type":"docker","code":"","dockerImage":"speaker","name":"Speaker","user":"fogflow","inputTriggers":[{"name":"selector2","selectedAttributeList":["all"],"groupedAttributeList":["id"],"conditionList":[{"type":"EntityType","value":"Speaker"}]}],"outputAnnotators":[]},
    designboard: {"edges":[{"id":1,"block1":2,"connector1":["selector","output"],"block2":1,"connector2":["selectors","input"]},{"id":2,"block1":3,"connector1":["condition","output"],"block2":2,"connector2":["conditions","input"]}],"blocks":[{"id":1,"x":96,"y":-133,"type":"FogFunction","module":null,"values":{"name":"Speaker","user":"fogflow"}},{"id":2,"x":-141,"y":-134,"type":"InputTrigger","module":null,"values":{"selectedattributes":["all"],"groupby":["id"]}},{"id":3,"x":-373,"y":-136,"type":"SelectCondition","module":null,"values":{"type":"EntityType","value":"Speaker"}}]}
}
];


//connect to the broker
var client = new NGSI10Client(config.brokerURL);

addMenuItem('Function', showFunction);  
addMenuItem('Task', showTask);  
addMenuItem('Editor', showEditor);  


showFunction();

// the list of all registered operators
var operatorList = [];
queryOperatorList();

queryFunctionList();


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
    if(handler != null) {
        handler();        
    }
}

function initFogFunctionExamples() 
{
    for(var i=0; i<myFogFunctionExamples.length; i++) {
        var example = myFogFunctionExamples[i];
        submitFunction(example.fogfunction, example.designboard);
    }
}

function queryFunctionList() 
{
    var queryReq = {}
    queryReq.entities = [{type:'FogFunction', isPattern: true}];
    client.queryContext(queryReq).then( function(fogFunctionList) {
        if (fogFunctionList.length == 0) {
			initFogFunctionExamples();
		}
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    });          
}


function showDesignBoard()
{
    var html = '';
    
    html += '<div class="input-prepend">';         
    html += '<button id="cleanBoard" type="button" class="btn btn-default">Clean Board</button>';                            
    html += '<button id="saveBoard" type="button" class="btn btn-default">Save Board</button>';                                
    html += '<button id="generateFunction" type="button" class="btn btn-primary">Create a Fog Function</button>';        
    html += '<button id="displayFogFunctionObject" type="button" class="btn btn-default">Display as JSON</button>';         
    html += '</div>'; 
        
    html += '<div id="blocks" style="width:1000px; height:400px"></div>';
    
    html += '<div style="margin-top: 10px;"><h4 class="text-left">Function code</h4>';
    html += '<select id="codeType"><option value="javascript">javascript</option><option value="python"">python</option><option value="docker"">dockerimage</option></select>';    
    html += '<div id="codeBoard"></div>';            
    html += '</div>'    
    
    $('#content').html(html);  
    
    var boardHTML = '<textarea id="codeText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
    $('#codeBoard').html(boardHTML);
    $('#codeText').val(template.javascript);
   
	
    var blocks = new Blocks();
 
    // prepare the configuration
    var config = {};

    // prepare the design board
    registerAllBlocks(blocks);
  
    blocks.run('#blocks');
    
    if (CurrentScene != null ) {
		console.log(CurrentScene);
        blocks.importData(CurrentScene);
    }  		
	
    blocks.ready(function() {                
        $('#generateFunction').click(function() {
            generateFogfunction(blocks.export());
        });    
        $('#cleanBoard').click(function() {
            blocks.clear();
        });                      
        $('#saveBoard').click(function() {
            CurrentScene = blocks.export();
        });                              
        $('#displayFogFunctionObject').click(function() {
            var board = blocks.export();
            var fogfunction = boardScene2fogfunction(board);    
            var ffObj = {
                fogfunction: fogfunction,
                designboard: board
            };
            alert(JSON.stringify(ffObj));
        });                                      
    }); 	  	
}


function showEditor() 
{
    $('#info').html('editor to design a fog function');
    
	showDesignBoard();
    
    $('#codeType').change(function() {
        var fType = $(this).val();
        switch(fType) {
            case 'javascript':
                var boardHTML = '<textarea id="codeText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
                $('#codeBoard').html(boardHTML);
                $('#codeText').val(template.javascript);
                break;
            case 'python':
                var boardHTML = '<textarea id="codeText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
                $('#codeBoard').html(boardHTML);
                $('#codeText').val(template.python);
                break;
            case 'docker':
                var boardHTML = '<select id="codeImage"></select>';
                $('#codeBoard').html(boardHTML);                
                for(var i=0; i<operatorList.length; i++){
                    var operatorName = operatorList[i];
                    $('#codeImage').append($("<option></option>").attr("value", operatorName).text(operatorName));                                                    
                }                
                break;
        }        
    });    
    
    //initialize the content in the code textarea
    $('#codeText').val(template.javascript);              
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

function generateFogfunction(scene)
{
    // construct the fog function object based on the design board
    var fogfunction = boardScene2fogfunction(scene);    
   
    // submit this fog function
    submitFunction(fogfunction, scene);
}



function boardScene2fogfunction(scene)
{
    console.log(scene);  
    var fogfunction = {};    
    
    // check the function type and the provided function code
    var fType = $('#codeType option:selected').val();    
    fogfunction.type = fType;
    
    switch(fType) {
        case 'javascript':
            var fCode = $('#codeText').val();
            fogfunction.code = fCode;    
            fogfunction.dockerImage = 'nodejs';           
            break;
        case 'python':
            var fCode = $('#codeText').val();
            fogfunction.code = fCode;    
            fogfunction.dockerImage = 'pythonbase';           
            break;
        case 'docker':
            var dockerImage = $('#codeImage option:selected').val();            
            fogfunction.code = '';    
            fogfunction.dockerImage = dockerImage;           
            break;        
    }     
    
    // check the defined inputs and outputs of this function
    for(var i=0; i<scene.blocks.length; i++){
        var block = scene.blocks[i];
        
        console.log(block.name);
        
        if(block.type == "FogFunction") {
            fogfunction.name = block.values['name'];
            fogfunction.user = block.values['user'];
            
            // construct its input streams
            fogfunction.inputTriggers = findInputTriggers(scene, block.id);
            
            // construct its output streams
            fogfunction.outputAnnotators = findOutputAnnotators(scene, block.id);   
            
            break;         
        }
    }        
    
    return fogfunction;    
}

function findInputTriggers(scene, blockId)
{
    var selectors = [];

    for(var i=0; i<scene.edges.length; i++){
        var edge = scene.edges[i];
        
        if(edge.block2 == blockId) {
            for(var j=0; j<scene.blocks.length; j++) {
                var block = scene.blocks[j];
                
                if(block.id == edge.block1) {
                    var selector = {};
                    selector.name = "selector" + block.id
                    selector.selectedAttributeList = block.values.selectedattributes;
                    selector.groupedAttributeList = block.values.groupby;
                    selector.conditionList = findConditions(scene, block.id);
                    
                    selectors.push(selector);
                }
            }               
        }
    }
    
    return selectors;
}


function findConditions(scene, blockId)
{
    var conditions = [];
    
    for(var i=0; i<scene.edges.length; i++){
        var edge = scene.edges[i];    
        
        if(edge.block2 == blockId) {        
            for(var j=0; j<scene.blocks.length; j++) {
                var block = scene.blocks[j];
                
                if(block.id == edge.block1) {        
                    var condition = {};                    
                    condition.type = block.values.type;
                    condition.value = block.values.value;                                    
                    conditions.push(condition);
                }
            }
        }
    }
            
    return conditions
}

function findOutputAnnotators(scene, blockId)
{
    var annotators = [];
    
    for(var i=0; i<scene.edges.length; i++){
        var edge = scene.edges[i];    
        
        if(edge.block1 == blockId) {                    
            for(var j=0; j<scene.blocks.length; j++) {
                var block = scene.blocks[j];
                
                if(block.id == edge.block2) {        
                    var annotator = {};    
                    
                    annotator.entityType = block.values.entitytype;
                    annotator.groupInherited = block.values.herited;                
                    
                    annotators.push(annotator);
                }
            }
        }
    }            
    
    return annotators;    
}


function submitFunction(fogfunction, designboard)
{
	console.log("==============================")
    console.log(JSON.stringify(fogfunction));  
	console.log(JSON.stringify(designboard));
	console.log("============end===============")

        
    var functionCtxObj = {};
    
    functionCtxObj.entityId = {
        id : 'FogFunction.' + fogfunction.name, 
        type: 'FogFunction',
        isPattern: false
    };
    
    functionCtxObj.attributes = {};   
    functionCtxObj.attributes.status = {type: 'string', value: 'enabled'};
    functionCtxObj.attributes.designboard = {type: 'object', value: designboard};    	
    functionCtxObj.attributes.fogfunction = {type: 'object', value: fogfunction};    
    
    client.updateContext(functionCtxObj).then( function(data) {
        console.log(data);  
              
        // update the list of submitted topologies
        updateFogFunctionList();                       
    }).catch( function(error) {
        console.log('failed to submit the fog function');
    });           
}

function showFunction() 
{
    var queryReq = {}
    queryReq.entities = [{type:'FogFunction', isPattern: true}];
    client.queryContext(queryReq).then( function(fogFunctionList) {
        console.log(fogFunctionList);
        displayFunctionList(fogFunctionList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    });          
}

function updateFogFunctionList() 
{
    var queryReq = {}
    queryReq.entities = [{type:'FogFunction', isPattern: true}];
    client.queryContext(queryReq).then( function(functionList) {
        console.log(functionList);
        displayFunctionList(functionList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });       
}

function displayFunctionList(fogfunctions) 
{
    $('#info').html('list of all submitted fog functions');

    if(fogfunctions.length == 0) {
        $('#content').html('');
        return;
    }          

    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>FogFunction</th>';
    html += '<th>Status</th>';    
    html += '<th>Control</th>';        
    html += '</tr></thead>';    

    for(var i=0; i<fogfunctions.length; i++){
        var fogfunction = fogfunctions[i];
        
        html += '<tr>';		
		html += '<td>' + fogfunction.entityId.id + '<br><button id="editor-' + fogfunction.entityId.id + '" type="button" class="btn btn-default">editor</button>';        
		html += '<br><button id="delete-' + fogfunction.entityId.id + '" type="button" class="btn btn-default">delete</button></td>';        
		html += '<td>' + JSON.stringify(fogfunction.attributes['fogfunction'].value) + '</td>'; 
        
        var status = fogfunction.attributes['status'].value;        
        
        if ( status == 'disabled' ) {
		    html += '<td>inactive</td>';               
            html += '<td><button id="control-' + fogfunction.entityId.id + '" type="button" class="btn btn-primary">enable</button></td>';
        } else {
		    html += '<td>active</td>';                           
            html += '<td><button id="control-' + fogfunction.entityId.id + '" type="button" class="btn btn-primary">disable</button></td>';
        }
             
		html += '</tr>';			
	}
       
    html += '</table>';            
	
	$('#content').html(html);     
    
    // associate a click handler to the control button
    for(var i=0; i<fogfunctions.length; i++){
        var fogfunction = fogfunctions[i];
                
        var ctrlButton = document.getElementById("control-" + fogfunction.entityId.id);
        ctrlButton.onclick = function(f) {
            var myFunction = f;
            return function(){
                switchFogFunctionStatus(myFunction);
            };
        }(fogfunction);	
		
        var editorButton = document.getElementById("editor-" + fogfunction.entityId.id);
        editorButton.onclick = function(f) {
            var myFunction = f;
            return function(){
                openEditor(myFunction);
            };
        }(fogfunction);		
		
        var deleteButton = document.getElementById("delete-" + fogfunction.entityId.id);
        deleteButton.onclick = function(f) {
            var myFunction = f;
            return function(){
                deleteFunction(myFunction);
            };
        }(fogfunction);			
	}      	   	
}


function switchFogFunctionStatus(fogFunc)
{
    var functionCtxObj = {};    
    
    // switch the status
    functionCtxObj.entityId = fogFunc.entityId
    
    functionCtxObj.attributes = {};   
    
    if (fogFunc.attributes.status.value == "enabled") {
        functionCtxObj.attributes.status = {type: 'string', value: 'disabled'};        
    } else {
        functionCtxObj.attributes.status = {type: 'string', value: 'enabled'};        
    }
    
    client.updateContext(functionCtxObj).then( function(data) {
        console.log(data);                
        // update the list of submitted topologies
        updateFogFunctionList();                       
    }).catch( function(error) {
        console.log('failed to submit the topology');
    });      
}

function deleteFunction(fogFunc)
{
    var entityid = {
        id : fogFunc.entityId.id, 
        type: 'FogFunction',
        isPattern: false
    };	    
    
    client.deleteContext(entityid).then( function(data) {
        console.log(data);
		updateFogFunctionList();		
    }).catch( function(error) {
        console.log('failed to delete a fog function');
    });  	
}

function showTask() 
{
    $('#info').html('list of all triggerred function tasks');
            
    var queryReq = {}
    queryReq.entities = [{type:'Task', isPattern: true}];
    queryReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'topology=system'}]}    

    client.queryContext(queryReq).then( function(taskList) {
        console.log(taskList);
        displayTaskList(taskList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    });        
}


function displayTaskList(tasks) 
{
    $('#info').html('list of all function tasks that have been triggerred');

    if(tasks.length == 0) {
        $('#content').html('');
        return;
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


function openEditor(fogfunctionEntity)
{
    if(fogfunctionEntity.attributes.designboard){
        CurrentScene = fogfunctionEntity.attributes.designboard.value;   	
    }
		
	//selectMenuItem('Editor');
	//window.location.hash = '#Editor';			
	showEditor();
       
    var fogfunction = fogfunctionEntity.attributes.fogfunction.value;
			    		
	// check the function type and the provided function code
	$('#codeType').val(fogfunction.type);
	
    switch(fogfunction.type) {
        case 'javascript':		
            var boardHTML = '<textarea id="codeText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
            $('#codeBoard').html(boardHTML);		
       		$('#codeText').val(fogfunction.code);          
       		break;
   		case 'python':			
            var boardHTML = '<textarea id="codeText" class="form-control" style="min-width: 800px; min-height: 200px"></textarea>';
            $('#codeBoard').html(boardHTML);		
       		$('#codeText').val(fogfunction.code);          
       		break;
   		case 'docker':				
            var boardHTML = '<select id="codeImage"></select>';
            $('#codeBoard').html(boardHTML); 
            for(var i=0; i<operatorList.length; i++){
                var operatorName = operatorList[i];
                $('#codeImage').append($("<option></option>").attr("value", operatorName).text(operatorName));                                                    
            } 			               
       		$('#codeImage').val(fogfunction.dockerImage);                     
       		break;        
	} 

}

});


