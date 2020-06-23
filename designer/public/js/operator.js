$(function(){

// initialization  
var handlers = {};  

//connect to the broker
var client = new NGSI10Client(config.brokerURL);

//DGraph
// to interact with designer
var clientDes = new NGSIDesClient(config.designerIP+':'+config.webSrvPort);

console.log(config.brokerURL);

addMenuItem('Operator', showOperator);  
addMenuItem('DockerImage', showDockerImage);    

initOperatorList();
initDockerImageList();

showOperator();


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


function showOperator()
{
    $('#info').html('list of all registered operators');
    
    var html = '<div style="margin-bottom: 10px;"><button id="registerOperator" type="button" class="btn btn-primary">register</button></div>';
    html += '<div id="operatorList"></div>';

	$('#content').html(html);   
      
    $( "#registerOperator" ).click(function() {
        showOperatorEditor();
    });  
        
    var queryReq = {}
    queryReq.entities = [{type:'Operator', isPattern: true}];           
    
    client.queryContext(queryReq).then( function(operatorList) {
        displayOperatorList(operatorList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
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
            
            var option = document.createElement("option");
            
            option.text = operator.name;
            
            var operatorList = document.getElementById("OperatorList");                           
            operatorList.add(option);
    	} 
        
        // add it into the select list        
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });    
}

function displayOperatorList(operators) 
{
    if(operators == null || operators.length == 0){
        $('#operatorList').html('');           
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>Operator</th>';    
    html += '<th>Name</th>';
    html += '<th>Description</th>';    
    html += '<th>#Parameters</th>';
    html += '<th>#Implementations</th>';
    html += '</tr></thead>';    
           
    for(var i=0; i<operators.length; i++){
        var entity = operators[i];
        
        var operator = entity.attributes.operator.value;
        
        html += '<tr>'; 
		html += '<td>' + entity.entityId.id + '</td>';                        
		html += '<td>' + operator.name + '</td>';                
		html += '<td>' + operator.description + '</td>';                
		html += '<td>' + operator.parameters.length + '</td>';        
		html += '<td>' + 0 + '</td>';                
                              
		html += '</tr>';	                        
	}    
       
    html += '</table>';  
    
	$('#operatorList').html(html);      
}


function showOperatorEditor() 
{
    $('#info').html('to specify an operator');

    var html = '';
    
    html += '<div><button id="generateOperator" type="button" class="btn btn-primary">Submit</button></div>';                  
    html += '<div id="blocks" style="width:1000px; margin-top: 5px; height:400px"></div>';
    
    $('#content').html(html);    

    blocks = new Blocks();
 
    registerAllBlocks(blocks);

    blocks.run('#blocks');
    
    blocks.types.addCompatibility('string', 'choice');
            
    blocks.ready(function() {                
        // associate functions to clickable buttons
        $('#generateOperator').click(function() {
            generateOperator(blocks.export());
        });    
    });    
}

function generateOperator(scene)
{
    // construct the operator based on the design board
    var operator = boardScene2Operator(scene);    
   
    console.log(operator);

    // submit this operator
    submitOperator(operator, scene);
}

function submitOperator(operator, designboard)
{
    var operatorObj = {};
    
    operatorObj.entityId = {
        id : operator.name, 
        type: 'Operator',
        isPattern: false
    };
    
    operatorObj.attributes = {};   
    operatorObj.attributes.designboard = {type: 'object', value: designboard};    	
    operatorObj.attributes.operator = {type: 'object', value: operator};    
    
    operatorObj.metadata = {};      
    var geoScope = {};    
    geoScope.type = "global"
    geoScope.value = "global"
    operatorObj.metadata.location = geoScope;    
        
    
    clientDes.updateContext(operatorObj).then( function(data) {
        showOperator();                       
    }).catch( function(error) {
        console.log('failed to submit the defined operator');
    });           
}

function boardScene2Operator(scene)
{
    console.log(scene);  
    var operator = {};    
        
    for(var i=0; i<scene.blocks.length; i++){
        var block = scene.blocks[i];
                
        if(block.type == "Operator") {
            operator.name = block.values['name'];
            operator.description = block.values['description'];
            
            // construct its controllable parameters
            operator.parameters = findInputParameters(scene, block.id);
            
            break;         
        }
    }        

    console.log(operator);    
    return operator;    
}


function findInputParameters(scene, blockId)
{
    var parameters = [];

    for(var i=0; i<scene.edges.length; i++){
        var edge = scene.edges[i];
        
        if(edge.block2 == blockId) {
            for(var j=0; j<scene.blocks.length; j++) {
                var block = scene.blocks[j];
                
                if(block.id == edge.block1) {
                    var parameter = {};
                    parameter.name = block.values.name
                    parameter.values = block.values.values ;
                    
                    parameters.push(parameter);
                }
            }               
        }
    }
    
    return parameters;
}

function showDockerImage() 
{
    $('#info').html('list of docker images in the docker registry');

    var html = '<div style="margin-bottom: 10px;"><button id="registerDockerImage" type="button" class="btn btn-primary">register</button></div>';
    html += '<div id="dockerImageList"></div>';

	$('#content').html(html);   
      
    updateDockerImageList();       
    
    $( "#registerDockerImage" ).click(function() {
        dockerImageRegistration();
    });                
}


function initOperatorList()
{
    var operatorList = [{
        name: "nodejs",
        description: "",
        parameters:[]
    },{
        name: "python",
        description: "",
        parameters:[]
    },{
        name: "iotagent",
        description: "",
        parameters:[]
    },{
        name: "counter",
        description: "",
        parameters:[]
    },{
        name: "anomaly",
        description: "",
        parameters:[]
    },{
        name: "facefinder",
        description: "",
        parameters:[]
    },{
        name: "connectedcar",
        description: "",
        parameters:[]
    },{
        name: "recommender",
        description: "",
        parameters:[]
    },{
        name: "privatesite",
        description: "",
        parameters:[]
    },{
        name: "publicsite",
        description: "",
        parameters:[]
    },{
        name: "pushbutton",
        description: "",
        parameters:[]
    },{
        name: "acoustic",
        description: "",
        parameters:[]
    },{
        name: "speaker",
        description: "",
        parameters:[]
    },{
        name: "dummy",
        description: "",
        parameters:[]
    },{
        name: "geohash",
        description: "",
        parameters:[]
    },{
        name: "converter",
        description: "",
        parameters:[]
    },{
        name: "predictor",
        description: "",
        parameters:[]
    },{
        name: "controller",
        description: "",
        parameters:[]
    },{
        name: "detector",
        description: "",
        parameters:[]
    }
    ];
    
    var queryReq = {}
    queryReq.entities = [{type:'Operator', isPattern: true}];               
    client.queryContext(queryReq).then( function(existingOperatorList) {
        if (existingOperatorList.length == 0) {
            for(var i=0; i<operatorList.length; i++) {
                submitOperator(operatorList[i], {});
            }          
        }
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query the operator list');
    });     
}

function initDockerImageList()
{
    var imageList = [{
        name: "fogflow/nodejs",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "nodejs",
        prefetched: true
    },{
        name: "fogflow/python",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "python",
        prefetched: false
    },{
        name: "fogflow/counter",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "counter",
        prefetched: false
    },{
        name: "fogflow/anomaly",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "anomaly",
        prefetched: false
    },{
        name: "fogflow/facefinder",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "facefinder",
        prefetched: false
    },{
        name: "fogflow/connectedcar",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "connectedcar",
        prefetched: false
    },{
        name: "fiware/iotagent-json",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "iotagent",
        prefetched: false
    },{
        name: "fogflow/recommender",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "recommender",
        prefetched: false
    },{
        name: "fogflow/privatesite",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "privatesite",
        prefetched: false
    },{
        name: "fogflow/publicsite",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "publicsite",
        prefetched: false
    },{
        name: "pushbutton",
        tag: "latest",
        hwType: "ARM",
        osType: "Linux",
        operatorName: "pushbutton",
        prefetched: false
    },{
        name: "acoustic",
        tag: "latest",
        hwType: "ARM",
        osType: "Linux",
        operatorName: "acoustic",
        prefetched: false
    },{
        name: "speaker",
        tag: "latest",
        hwType: "ARM",
        osType: "Linux",
        operatorName: "speaker",
        prefetched: false
    },{
        name: "pushbutton",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "pushbutton",
        prefetched: false
    },{
        name: "acoustic",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "acoustic",
        prefetched: false
    },{
        name: "speaker",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "speaker",
        prefetched: false
    },{
        name: "fogflow/dummy",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "dummy",
        prefetched: false
    },{
        name: "geohash",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "geohash",
        prefetched: false
    },{
        name: "converter",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "converter",
        prefetched: false
    },{
        name: "predictor",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "predictor",
        prefetched: false
    },{
        name: "controller",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "controller",
        prefetched: false
    },{
        name: "detector",
        tag: "latest",
        hwType: "ARM",
        osType: "Linux",
        operatorName: "detector",
        prefetched: false
    }
    ];

    var queryReq = {}
    queryReq.entities = [{type:'DockerImage', isPattern: true}];               
    client.queryContext(queryReq).then( function(existingImageList) {
        if (existingImageList.length == 0) {
            for(var i=0; i<imageList.length; i++) {
                addDockerImage(imageList[i]);
            }          
        }
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query the image list');
    }); 

}

function addDockerImage(image) 
{    
    //register a new docker image
    var newImageObject = {};

    newImageObject.entityId = {
        id : image.name + '.' + image.tag, 
        type: 'DockerImage',
        isPattern: false
    };

    newImageObject.attributes = {};   
    newImageObject.attributes.image = {type: 'string', value: image.name};        
    newImageObject.attributes.tag = {type: 'string', value: image.tag};    
    newImageObject.attributes.hwType = {type: 'string', value: image.hwType};      
    newImageObject.attributes.osType = {type: 'string', value: image.osType};          
    newImageObject.attributes.operator = {type: 'string', value: image.operatorName};      
    newImageObject.attributes.prefetched = {type: 'boolean', value: image.prefetched};                      
    
    newImageObject.metadata = {};    
    newImageObject.metadata.operator = {
        type: 'string',
        value: image.operatorName
    };         
    
    var geoScope = {};    
    geoScope.type = "global"
    geoScope.value = "global"
    newImageObject.metadata.location = geoScope;            

    clientDes.updateContext(newImageObject).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });      
    
}

function dockerImageRegistration()
{
    $('#info').html('to register a new docker image');
    
    var html = '<div id="dockerRegistration" class="form-horizontal"><fieldset>';                 
    
    html += '<div class="control-group"><label class="control-label" for="input01">Image(*)</label>';
    html += '<div class="controls"><input type="text" class="input-xlarge" id="dockerImageName">';
    html += '</div></div>';
    
    html += '<div class="control-group"><label class="control-label" for="input01">Tag(*)</label>';
    html += '<div class="controls"><input type="text" class="input-xlarge" id="imageTag" placeholder="latest">';
    html += '</div></div>';    
    
    html += '<div class="control-group"><label class="control-label" for="input01">HardwareType(*)</label>';
    html += '<div class="controls"><select id="hwType"><option>X86</option><option>ARM</option></select></div>'
    html += '</div>';    
    
    html += '<div class="control-group"><label class="control-label" for="input01">OSType(*)</label>';
    html += '<div class="controls"><select id="osType"><option>Linux</option><option>Windows</option></select></div>'
    html += '</div>';    

    html += '<div class="control-group"><label class="control-label" for="input01">Operator(*)</label>';
    html += '<div class="controls"><select id="OperatorList"></select>';
    html += '</div></div>';    

    html += '<div class="control-group"><label class="control-label" for="optionsCheckbox">Prefetched</label>';
    html += '<div class="controls"> <label class="checkbox"><input type="checkbox" id="Prefetched" value="option1">';
    html += 'docker image must be fetched by the platform in advance';
    html += '</label></div>';
    html += '</div>';        

    
    html += '<div class="control-group"><label class="control-label" for="input01"></label>';
    html += '<div class="controls"><button id="submitRegistration" type="button" class="btn btn-primary">Register</button>';
    html += '</div></div>';   
       
    html += '</fieldset></div>';

	$('#content').html(html);          
    
    queryOperatorList();
        
    // associate functions to clickable buttons
    $('#submitRegistration').click(registerDockerImage);  
}


function registerDockerImage() 
{    
    console.log('register a new docker image'); 

    // take the inputs    
    var image = $('#dockerImageName').val();
    console.log(image);
    
    var tag = $('#imageTag').val();
    if (tag == '') {
        tag = 'latest';
    }
    
    console.log(tag);    
    
    var hwType = $('#hwType option:selected').val();
    console.log(hwType);
    
    var osType = $('#osType option:selected').val();
    console.log(osType);    
    
    var operatorName = $('#OperatorList option:selected').val();
    console.log(operatorName);        
    
    var prefetched = document.getElementById('Prefetched').checked;
    console.log(prefetched);        
    
               
    if( image == '' || tag == '' || hwType == '' || osType == '' || operatorName == '' ) {
        alert('please provide the required inputs');
        return;
    }    

    //register a new docker image
    var newImageObject = {};

    newImageObject.entityId = {
        id : image + ':' + tag, 
        type: 'DockerImage',
        isPattern: false
    };

    newImageObject.attributes = {};   
    newImageObject.attributes.image = {type: 'string', value: image};        
    newImageObject.attributes.tag = {type: 'string', value: tag};    
    newImageObject.attributes.hwType = {type: 'string', value: hwType};      
    newImageObject.attributes.osType = {type: 'string', value: osType};          
    newImageObject.attributes.operator = {type: 'string', value: operatorName};  
    
    if (prefetched == true) {
        newImageObject.attributes.prefetched = {type: 'boolean', value: true};                      
    } else {
        newImageObject.attributes.prefetched = {type: 'boolean', value: false};                      
    }
            
    newImageObject.metadata = {};    
    newImageObject.metadata.operator = {
        type: 'string',
        value: operatorName
    };               

    clientDes.updateContext(newImageObject).then( function(data) {
        console.log(data);
        
        // show the updated image list
        showDockerImage();
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });      
    
}


function updateDockerImageList()
{
    var queryReq = {}
    queryReq.entities = [{type:'DockerImage', isPattern: true}];           
    
    client.queryContext(queryReq).then( function(imageList) {
        displayDockerImageList(imageList);
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });    
}

function displayDockerImageList(images) 
{
    if(images == null || images.length == 0){
        $('#dockerImageList').html('');           
        return        
    }
    
    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>Operator</th>';    
    html += '<th>Image</th>';
    html += '<th>Tag</th>';
    html += '<th>Hardware Type</th>';
    html += '<th>OS Type</th>';    
    html += '<th>Prefetched</th>';        
    html += '</tr></thead>';    
       
    for(var i=0; i<images.length; i++){
        var dockerImage = images[i];
		
        html += '<tr>'; 
		html += '<td>' + dockerImage.attributes.operator.value + '</td>';                
		html += '<td>' + dockerImage.attributes.image.value + '</td>';                
		html += '<td>' + dockerImage.attributes.tag.value + '</td>';        
		html += '<td>' + dockerImage.attributes.hwType.value + '</td>';                
		html += '<td>' + dockerImage.attributes.osType.value + '</td>';  
        
        if (dockerImage.attributes.prefetched.value == true) {
		    html += '<td><font color="red"><b>' + dockerImage.attributes.prefetched.value + '</b></font></td>';                                            
        } else {
		    html += '<td>' + dockerImage.attributes.prefetched.value + '</td>';                                            
        }
                              
		html += '</tr>';	                        
	}
       
    html += '</table>';  
    
	$('#dockerImageList').html(html);      
}

});


