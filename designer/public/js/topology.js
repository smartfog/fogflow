'use strict';

$(function() {

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
    var eTypeList = [];
    var selectedServiceIntent = null;
    // design board
    var blocks = null;

    // client to interact with IoT Broker
    var client = new NGSI10Client(config.brokerURL);

    addMenuItem('Topology', 'Service Topology', showTopologies);
    addMenuItem('Intent', 'Service Intent', showIntents);
    
    initTopologyExamples();
    
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

    function initTopologyExamples() {
        fetch('/service').then(res => res.json()).then(tologies => {
            if (Object.keys(tologies).length === 0) {
                fetch("/service", {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(myToplogyExamples)
                })
                .then(response => {
                    console.log("create a list of service ontologies", response.status)
                    showTopologies();                                                       
                })
                .catch(err => console.log(err));                                                                 
            } else {
                showTopologies();
            }               
        })
    }


    function showTopologyEditor() {
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

        registerAllBlocks(blocks, operatorList, eTypeList);

        blocks.run('#blocks');

        blocks.types.addCompatibility('string', 'choice');

        if (CurrentScene != null) {
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

    function openTopologyEditor(topologyEntity) {
        if (topologyEntity.designboard) {
            CurrentScene = topologyEntity.designboard;
            showTopologyEditor();

            var topology = topologyEntity.topology;
            $('#serviceName').val(topology.name);
            $('#serviceDescription').val(topology.description);
        }
    }

    function deleteTopology(topologyEntity) {
        console.log("DELETE");
        var msg = {name: topologyEntity.topology.name}
        
        fetch("/service/" + topologyEntity.topology.name, {
            method: "DELETE"
        })
        .then(response => {
            console.log("delete a service tology: ", response.status)
            showTopologies();
        })
        .catch(err => console.log(err));           
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

    function boardScene2Topology(scene) {
        // construct a topology from the provided information
        var uidTopology =  $('#uidTopology').val();

        var topologyName = $('#serviceName').val();
        var serviceDescription = $('#serviceDescription').val();

        var serviceTopology = {};

        var topology = {};
        topology.name = topologyName;
        topology.description = serviceDescription;
        topology.tasks = generateTaskList(scene);
        if (uidTopology){
            topology.uid = uidTopology;
        }
        
        serviceTopology.topology = topology;
        serviceTopology.designboard = scene;
       
        // submit the generated topology
        submitTopology(serviceTopology);
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

    function submitTopology(topology) {
        console.log(topology)
        
        fetch("/service", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify([topology])
        })
        .then(response => {
            console.log("submit a new service topology: ", response.status)
            showTopologies();
        })
        .catch(err => console.log(err));
    }

    function showTopologies() {
        $('#info').html('list of all registered service topologies');

        var html = '<div style="margin-bottom: 10px;"><button id="registerTopology" type="button" class="btn btn-primary">register</button></div>';
        html += '<div id="topologyList"></div>';

        $('#content').html(html);

        $("#registerTopology").click(function() {
            showTopologyEditor();
        });

        // update the list of submitted topologies
        updateTopologyList();
    }

    function updateTopologyList() {
        fetch('/service').then(res => res.json()).then(topologies => {
            var topologyList = Object.values(topologies);
            console.log("get topology list ", topologyList)
            displayTopologyList(topologyList);             
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of service ontologies');
        });
    }

    function displayTopologyList(topologies) {
        if (topologies == null || topologies.length == 0) {
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Name</th>';
        html += '<th>Description</th>';
        html += '<th>#Tasks</th>';
        html += '<th>Actions</th>';        
        html += '</tr></thead>';

        for (var i = 0; i < topologies.length; i++) {
            var topologyEntity = topologies[i];
            
            console.log(topologyEntity)

            html += '<td>' + topologyEntity.topology.name + '</td>';

            html += '<td>' + topologyEntity.topology.description + '</td>';
            html += '<td>' + topologyEntity.topology.tasks.length + '</td>';

            html += '<td class="singlecolumn">';
            html += '<button id="editor-' + topologyEntity.topology.name  + '" type="button" class="btn btn-primary btn-separator">view</button>';
            html += '<button id="delete-' + topologyEntity.topology.name + '" type="button" class="btn btn-primary btn-separator">delete</button>';
            html += '</td>';
            
            html += '</tr>';
        }

        html += '</table>';

        $('#topologyList').html(html);

        // associate a click handler to the editor button
        for (var i = 0; i < topologies.length; i++) {
            var topologyEntity = topologies[i];
            // association handlers to the buttons
            var editorButton = document.getElementById('editor-' + topologyEntity.topology.name );
            editorButton.onclick = function(mytopology) {
                return function() {
                    openTopologyEditor(mytopology);
                };
            }(topologyEntity);

            var deleteButton = document.getElementById('delete-' + topologyEntity.topology.name );
            deleteButton.onclick = function(mytopology) {
                return function() {
                    deleteTopology(mytopology);
                };
            }(topologyEntity);
        }
    }

    function showIntents() {
        $('#info').html('list of active intentss');

        var html = '<div style="margin-bottom: 10px;"><button id="addIntent" type="button" class="btn btn-primary">add</button></div>';
        html += '<div id="intentList"></div>';

        $('#content').html(html);

        $("#addIntent").click(function() {
            addIntent();
        });

        fetch('/intent').then(res => res.json()).then(intents => {
            var intentList = Object.values(intents);
            displayIntentList(intentList);            
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of service intents');
        });
    }


    function displayIntentList(intents) {
        if (intents == null || intents.length == 0) {
            $('#intentList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Action</th>';
        html += '<th>Topology</th>';
        html += '<th>Service Type</th>';
        html += '<th>Priority</th>';
        html += '<th>Sevice Level Objective</th>';
        html += '<th>GeoScope</th>';
        html += '</tr></thead>';

        for (var i = 0; i < intents.length; i++) {
            var intent = intents[i];

            html += '<tr>';
            html += '<td>' + intent.id;
            html += '</td>';
            html += '<td class="singlecolumn">';

            html += '<button id="TASKS-' + intent.id + '" type="button" class="btn btn-primary">Tasks</button>  ';            
            html += '<button id="UPDATE-' + intent.id + '" type="button" class="btn btn-primary">Update</button>  ';
            html += '<button id="DELETE-' + intent.id + '" type="button" class="btn btn-primary">Delete</button>';

            html += '</td>';
            html += '<td>' + intent.topology + '</td>';
            html += '<td>' + intent.stype + '</td>';
            html += '<td>' + JSON.stringify(intent.priority) + '</td>';
            html += '<td>' + intent.qos + '</td>';
            html += '<td>' + JSON.stringify(intent.geoscope) + '</td>';            

            html += '</tr>';
        }

        html += '</table>';

        $('#intentList').html(html);

        // associate a click handler to generate device profile on request
        for (var i = 0; i < intents.length; i++) {
            var intent = intents[i];
            
            var taskButton = document.getElementById('TASKS-' + intent.id);
            taskButton.onclick = function(intent) {
                return function() {
                    queryTasks(intent);
                };
            }(intent);            
            
            var deleteButton = document.getElementById('DELETE-' + intent.id);
            deleteButton.onclick = function(intent) {
                return function() {
                    removeIntent(intent);
                };
            }(intent);

            var updateButton = document.getElementById('UPDATE-' + intent.id);
            updateButton.onclick = function(intent) {
                return function() {
                    $('#intentDgraphUID').val(intent.uid);
                    selectedServiceIntent = intent;
                    showIntent(intent);
                };
            }(intent);
        }
    }

    function showIntent(intentEntity) {
        var html = '<div id="intentRegistration" class="form-horizontal"><fieldset>';
        html +=    '<div id="intentDgraphUID" style="display: none;"> </div>';
        html += '<div class="control-group hidediv"><label class="control-label hidediv" for="input01">ID</label>';
        html += '<div class="controls hidediv"><lable class="hidediv" id="SID">sid</label></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Topology</label>';
        html += '<div class="controls"><select id="topologyItems"></select></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label">Type</label><div class="controls">';
        html += '<select id="SType">';
        html += '<option value="SYN">Synchronous</option>';
        html += '<option value="ASYN">Asynchronous</option></select>';
        html += '</div></div>';
                
        html += '<div class="control-group"><label class="control-label">Priority</label><div class="controls">';
        html += '<select id="priorityLevel"><option>low</option><option>middle</option><option>high</option></select>';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label">Resource usage</label><div class="controls">';
        html += '<select id="resouceUsage"><option>inclusive</option><option>exclusive</option></select>';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label">Objective</label><div class="controls">';
        html += '<select id="SLO">';
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
        html += '<div class="controls"><button id="submitIntent" type="button" class="btn btn-primary">Apply</button>';
        html += '</div></div>';

        html += '</fieldset></div>';

        $('#content').html(html);

        // add all service topologies into the selection list
        listAllServiceTopologies();

        // set the value accordingly
        if (intentEntity == undefined) {            
            var uid = uuid();
            var sid = 'ServiceIntent.' + uid;

            $("#SID").text(sid);
        } else {
            $('#intentDgraphUID').val(intentEntity.uid);
            var sid = intentEntity.id;
            $("#SID").text(sid);

            console.log(intentEntity);
            var intent = intentEntity;//intentEntity.attributes["intent"].value;
            console.log(intent);

            $('#topologyItems').val(intent.topology);
            //$('#SType').val(intent.stype);
            $('#SLO').val(intent.qos);

            if (intent.priority.exclusive) {
                $('#resouceUsage').val("exclusive");
            } else {
                $('#resouceUsage').val("inclusive");
            }

            if (intent.priority.level == 0) {
                $('#priorityLevel').val("low");
            } else if (intent.priority.level == 50) {
                $('#priorityLevel').val("middle");
            } else if (intent.priority.level == 100) {
                $('#priorityLevel').val("high");
            }

            if (intent.geoscope.scopeType == "global") {
                $('#geoscope').val("global");
            } else {
                $('#geoscope').val("custom");
                geoscope = intent.geoscope;
                showMap();
            }
        }

        // associate functions to clickable buttons
        $('#submitIntent').click(submitIntent);

        $('#geoscope').change(function() {
            var scope = $(this).val();

            if (scope == "custom") {
                // show the map to set locations
                showMap();
            } else {
                removeMap();
            }
        });
    }


    function addIntent() {
        selectedServiceIntent = null;
        $('#info').html('to specify an intent object in order to run your service');
        showIntent();
    }

    function removeIntent(intentObj) {
        fetch("/intent/" + intentObj.id, {
            method: "DELETE"
        })
        .then(data => {
            showIntents();
        })
        .catch(err => console.log(err));    
    }

    function listAllServiceTopologies() {
        fetch('/service').then(res => res.json()).then(servicetopologies => {
            Object.values(servicetopologies).forEach(servicetopology => {
                var name = servicetopology.topology.name;
                var topologySelect = document.getElementById('topologyItems');                    
                topologySelect.options[topologySelect.options.length] = new Option(name, name);                
            })            
        }); 
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


    function submitIntent() {
        var intent = {};

        var topology = $('#topologyItems option:selected').val();
        intent.topology = topology;
        var sType = $('#SType option:selected').val();
        intent.stype = sType;
        var temp1 = $('#priorityLevel option:selected').val();
        var priorityLevel = 0;
        switch (temp1) {
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
        if (temp2 == 'exclusive') {
            exclusiveResourceUsage = true;
        }

        intent.priority = {
            'exclusive': exclusiveResourceUsage,
            'level': priorityLevel
        };

        var slo = $('#SLO option:selected').val();
        intent.qos = slo;

        var scope = $('#geoscope option:selected').val();

        var operationScope = {};
        if (scope == 'custom') {
            operationScope.scopeType = geoscope.type;
            operationScope.scopeValue = geoscope.value;
            intent.geoscope = operationScope;
        } else {
            operationScope.scopeType = scope;
            operationScope.scopeValue = scope;
            intent.geoscope = operationScope;
        }

        var sid = $("#SID").text();
        intent.id = sid;
        
        fetch("/intent", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify(intent)
        })
        .then(response => {
            console.log("issue a new intent: ", response.status)
            showIntents();
        })
        .catch(err => console.log(err));        
    }


    function showMap() {
        var htmlContent = '<label class="control-label" for="input01">Polygon</label><div class="controls"><div id="map"  style="width: 500px; height: 400px"></div></div>';
        $('#mapDiv').html(htmlContent);

        var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
            osm = L.tileLayer(osmUrl, { maxZoom: 7, zoom: 7 }),
            map = new L.Map('map', { zoomControl: false, layers: [osm], center: new L.LatLng(35.692221, 138.709059), zoom: 7 });

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

        map.on('draw:created', function(e) {
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
                for (var i in points) {
                    geoscope.value.vertices.push({ longitude: points[i][0], latitude: points[i][1] });
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
                for (var i in points) {
                    geoscope.value.vertices.push({ longitude: points[i][0], latitude: points[i][1] });
                }

                console.log(geoscope);
            }

            drawnItems.addLayer(layer);
        });

        // if geoscope is already set, display it on the map
        if (geoscope) {
            if (geoscope.type == 'circle') {
                var circle = geoscope.value;
                L.circle([circle.centerLatitude, circle.centerLongitude], circle.radius).addTo(map);
            } else if (geoscope.type == 'rectangle') {
                var circle = geoscope.value.vertices;
                L.polygon().addTo(map);
            }
        }

    }

    function removeMap() {
        $('#mapDiv').html('');
    }


    function showTaskInstances() {
        $('#info').html('list of running data processing tasks');

        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];

        client.queryContext(queryReq).then(function(taskList) {
            displayTaskList(taskList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }
    
    function queryTasks(intent) {
        fetch('/info/task/' + intent.id).then(res => res.json()).then(tasks => {
            console.log("the list of tasks ");
            console.log(tasks);
            var taskList = Object.values(tasks);
            displayTaskList(taskList);            
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of create tasks');
        });        
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
