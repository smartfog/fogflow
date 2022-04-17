'use strict';

$(function() {

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

    // the list of all available entity type registered by brokers
    var eTypeList = [];

    // design board
    var blocks = null;

    // client to interact with IoT Broker
    var client = new NGSI10Client(config.brokerURL);

    var selectedFogFunction = null;


    addMenuItem('FogFunction', 'Fog Function', showFogFunctions);
    //addMenuItem('TaskInstance', 'Task Instance', showTaskInstances);
   
    initFogFunctionExamples();
    
    queryOperatorList();
    queryEntityTypeList();


    $(window).on('hashchange', function() {
        var hash = window.location.hash;

        selectMenuItem(location.hash.substring(1));
    });

    function addMenuItem(id, name, func) {
        handlers[id] = func;
        $('#menu').append('<li id="' + id + '"><a href="' + '#' + id + '">' + name + '</a></li>');
    }

    function selectMenuItem(name) {
        $('#menu li').removeClass('active');
        var element = $('#' + name);
        element.addClass('active');

        var handler = handlers[name];
        if (handler != undefined) {
            handler();
        }
    }

    function initFogFunctionExamples() {
        fetch('/fogfunction').then(res => res.json()).then(fogfunctions => {
            if (Object.keys(fogfunctions).length === 0) {
                fetch("/fogfunction", {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(myFogFunctionExamples)
                })
                .then(response => {
                    console.log("create the initial set of fog functions: ", response.status)
                    showFogFunctions();
                })
                .catch(err => console.log(err));                                                                
            } else {
                showFogFunctions();
            }               
        })        
    }

    function showFogFunctionEditor() {
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

        registerAllBlocks(blocks, operatorList, eTypeList);

        blocks.run('#blocks');

        blocks.types.addCompatibility('string', 'choice');

        if (CurrentScene != null) {
            blocks.importData(CurrentScene);
        }

        blocks.ready(function() {
            // associate functions to clickable buttons
            $('#generateFunction').click(function() {
                boardScene2FogFunction(blocks.export());
            });
            $('#cleanBoard').click(function() {
                blocks.clear();
            });
            $('#saveBoard').click(function() {
                CurrentScene = blocks.export();
            });
        });

    }

    function openFogFunctionEditor(fogfunction) {
        console.log(fogfunction);

        if (fogfunction && fogfunction.designboard) {
            CurrentScene = fogfunction.designboard;
            showFogFunctionEditor();

            var topology = fogfunction.topology;
            $('#serviceName').val(topology.name);
            $('#serviceDescription').val(topology.description);
        }
    }


    function queryOperatorList() {
        fetch('/operator').then(res => res.json()).then(operators => {
            Object.values(operators).forEach(operator => {
                operatorList.push(operator.name);
            })         
        });         
    }
    
    function queryEntityTypeList() {
        fetch('/info/type').then(res => res.json()).then(dtypes => {
            for (var i = 0; i < dtypes.length; i++) {
                eTypeList.push(dtypes[i]);
            }            
        }); 
    }    
    

    function boardScene2FogFunction(scene) {
        // step 1: construct the service topology object   
        var attribute = {}    
        var topologyName = $('#serviceName').val();
        var serviceDescription = $('#serviceDescription').val();

        var topology = {};
        topology.name = topologyName;
        topology.description = serviceDescription;
        topology.tasks = generateTaskList(scene);

        // step 2: construct an intent object
        var intent = {};

        var uid = uuid();
        var sid = 'ServiceIntent.' + uid;
                    
        intent.id = sid;        
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
       
        var fogfunction = {};        
        fogfunction.name = topologyName;
        fogfunction.topology = topology;
        fogfunction.intent = intent;
        fogfunction.designboard = scene;
        fogfunction.status = 'enabled';
        
               
        if (topologyName == '' || topology.name == '' || topology.tasks.length==0 || topology.tasks[0].operator == 'null' || 
        topology.tasks[0].operator == '' || topology.tasks[0].input_streams.length==0){
            alert('please provide the required inputs');
            return;
        }
        
        submitFogFunction(fogfunction);
    }

    function submitFogFunction(functionCtxObj) {
        fetch("/fogfunction", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify([functionCtxObj])
        })
        .then(response => {
            console.log("submit a new fog function: ", response.status)
            showFogFunctions();
        })
        .catch(err => console.log(err));  
    }


    function generateTaskList(scene) {
        var tasklist = [];

        for (var i = 0; i < scene.blocks.length; i++) {
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
                for (var j = 0; j < block.values['outputs'].length; j++) {
                    var outputstream = {};
                    outputstream.entity_type = block.values['outputs'][j];
                    task.output_streams.push(outputstream);
                }

                tasklist.push(task);
            }
        }

        return tasklist;
    }

    function findInputStream(scene, blockid) {
        var inputstreams = [];

        for (var i = 0; i < scene.edges.length; i++) {
            var edge = scene.edges[i];
            if (edge.block2 == blockid) {
                var inputblockId = edge.block1;

                for (var j = 0; j < scene.blocks.length; j++) {
                    var block = scene.blocks[j];
                    if (block.id == inputblockId) {
                        if (block.type == 'Shuffle') {
                            var inputstream = {};

                            inputstream.selected_type = findInputType(scene, block.id)

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

    function findInputType(scene, blockId) {
        var inputType = "unknown";

        for (var i = 0; i < scene.edges.length; i++) {
            var edge = scene.edges[i];

            if (edge.block2 == blockId) {
                var index = edge.connector1[2];

                for (var j = 0; j < scene.blocks.length; j++) {
                    var block = scene.blocks[j];
                    if (block.id == edge.block1) {
                        console.log(block);
                        inputType = block.values.outputs[index];
                    }
                }
            }
        }

        return inputType;
    }

    function showFogFunctions() {
        $('#info').html('list of all registered fog functions');

        var html = '<div style="margin-bottom: 10px;"><button id="registerFunction" type="button" class="btn btn-primary">register</button></div>';
        html += '<div id="functionList"></div>';

        $('#content').html(html);

        $("#registerFunction").click(function() {
            selectedFogFunction = null;
            CurrentScene = null
            showFogFunctionEditor();
        });

        // update the list of submitted fog functions
        updateFogFunctionList();
    }

    function updateFogFunctionList() {
        fetch('/fogfunction').then(res => res.json()).then(fogfunctions => {
            var fogfunctionList = Object.values(fogfunctions);
            displayFunctionList(fogfunctionList);             
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of fog functions');
        });
    }

    function displayFunctionList(fogFunctions) {
        if (fogFunctions == null || fogFunctions.length == 0) {
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Name</th>';
        html += '<th>Action</th>';
        html += '<th>Topology</th>';
        html += '<th>Intent</th>';
        html += '</tr></thead>';

        for (var i = 0; i < fogFunctions.length; i++) {
            var fogfunction = fogFunctions[i];
            
            console.log(fogfunction)

            html += '<td>' + fogfunction.name + '</td>';

            html += '<td>';
            
            html += '<button id="task-' + fogfunction.name + '" type="button" class="btn btn-primary btn-separator">tasks</button>';
            
            html += '<button id="editor-' + fogfunction.name + '" type="button" class="btn btn-primary btn-separator">view</button>';
            html += '<button id="delete-' + fogfunction.name + '" type="button" class="btn btn-primary btn-separator">delete</button>';
            
            if (fogfunction.status == 'enabled') {
                html += '<button id="status-' + fogfunction.name + '" type="button" class="btn btn-secondary btn-separator">disable</button>';                
            } else {
                html += '<button id="status-' + fogfunction.name + '" type="button" class="btn btn-success btn-separator">enable</button>';                
            }
            
            html += '</td>';

            html += '<td>' + JSON.stringify(fogfunction.topology.name) + '</td>';

            html += '<td>' + JSON.stringify(fogfunction.intent) + '</td>';

            html += '</tr>';
        }

        html += '</table>';

        $('#functionList').html(html);

        // associate a click handler to the editor button
        for (var i = 0; i < fogFunctions.length; i++) {
            var fogfunction = fogFunctions[i];
            console.log(fogfunction)
            
            // association handlers to the buttons
            var editorButton = document.getElementById('editor-' + fogfunction.name);
            editorButton.onclick = function(myFogFunction) {
                return function() {
                    console.log("editor buttion ",myFogFunction);
                    selectedFogFunction = myFogFunction;
                    openFogFunctionEditor(myFogFunction);
                };
            }(fogfunction);

            var deleteButton = document.getElementById('delete-' + fogfunction.name);
            deleteButton.onclick = function(myFogFunction) {
                return function() {
                    console.log("delete buttion ", myFogFunction);
                    deleteFogFunction(myFogFunction);
                };
            }(fogfunction);
            
            var statusButton = document.getElementById('status-' + fogfunction.name);
            statusButton.onclick = function(myFogFunction) {
                return function() {
                    if (statusButton.innerHTML == "enable") {
                        enableFogFunction(myFogFunction);                        
                    } else {
                        disableFogFunction(myFogFunction);                                                
                    }
                };
            }(fogfunction);
            
            var taskButton = document.getElementById('task-' + fogfunction.name);
            taskButton.onclick = function(myFogFunction) {
                return function() {
                    console.log("taskButton buttion ",myFogFunction);
                    selectedFogFunction = myFogFunction;
                    getTaskByFogFunction(myFogFunction);
                };
            }(fogfunction);                        
        }
    }

    function deleteFogFunction(fogfunction) {                
        fetch("/fogfunction/" + fogfunction.name, {
            method: "DELETE"
        })
        .then(response => {
            console.log("delete a fog function: ", response.status)
            showFogFunctions();
        })
        .catch(err => console.log(err));   
    }
    
    function enableFogFunction(fogfunction) {                
        fetch("/fogfunction/" + fogfunction.name + "/enable")
        .then(response => {
            console.log("enable a fog function: ", response.status)
            showFogFunctions();
        })
        .catch(err => console.log(err));   
    }    
    
    function disableFogFunction(fogfunction) {                
        fetch("/fogfunction/" + fogfunction.name + "/disable")
        .then(response => {
            console.log("disable a fog function: ", response.status)
            showFogFunctions();
        })
        .catch(err => console.log(err));   
    }      

    function getTaskByFogFunction(fogfunction) {              
        console.log(fogfunction);  
        fetch("/info/task/" + fogfunction.intent.id).then(res => res.json()).then(tasks => {
            displayTaskList(tasks);
        })
        .catch(err => console.log(err));   
    } 

    function uuid() {
        var uuid = "",
            i, random;
        for (i = 0; i < 32; i++) {
            random = Math.random() * 16 | 0;
            if (i == 8 || i == 12 || i == 16 || i == 20) {
                uuid += "-"
            }
            uuid += (i == 12 ? 4 : (i == 16 ? (random & 3 | 8) : random)).toString(16);
        }

        return uuid;
    }

    function displayTaskList(tasks) {
        if (tasks == null || tasks.length == 0) {
            $('#content').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Service</th>';
        html += '<th>Task</th>';
        html += '<th>Worker</th>';
        html += '<th>Status</th>';
        html += '</tr></thead>';

        for (var i = 0; i < tasks.length; i++) {
            var task = tasks[i];

            console.log(task);

            html += '<tr>';
            html += '<td>' + task.TaskID + '</td>';
            html += '<td>' + task.TopologyName + '</td>';
            html += '<td>' + task.TaskName + '</td>';
            html += '<td>' + task.Worker + '</td>';            

            if (task.Status == "paused") {
                html += '<td><font color="red">' + task.Status + '</font></td>';
            } else {
                html += '<td><font color="green">' + task.Status + '</font></td>';
            }

            html += '</tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }


});