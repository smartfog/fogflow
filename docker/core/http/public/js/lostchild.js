$(function(){

// initialization  
var handlers = {}

var childphotoURL = 'http://' + config.agentIP + ':' + config.webSrvPort + '/photo/lostchild.png';
var saveLocation = 'http://' + config.agentIP + ':' + config.webSrvPort + '/photo';

var cameraMarkers  = {};

var geoscope = {
    scopeType: "local",
    scopeValue: "local"
};

var curTopology = null;
var curIntent = null;
var checkingTimer = null;
var radiusStepDistance = 1000;
var curMap = null;

var personsFound = [];
var category_dataset = [{key: '#totalbytes', values:[]}];  
var chart;
var MARGIN = { top: 30, right: 20, bottom: 60, left: 80 };

var num_of_task = 0;

addMenuItem('Topology', showTopology);  
addMenuItem('Management', showMgt);      
addMenuItem('Tasks', showTasks);   
addMenuItem('Result', showResult);    

//connect to the socket.io server via the NGSI proxy module
var ngsiproxy = new NGSIProxy();
ngsiproxy.setNotifyHandler(handleNotify);

// client to interact with IoT Broker
var client = new NGSI10Client(config.brokerURL);
subscribeResult();
checkTopology();
checkIntent();

showTopology();
publishChildInfo();


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


//start a timer to send the child information periodically
function startUpdateTimer()
{
    setInterval(publishChildInfo, 2000);
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
    
    if (curIntent != null) {
        personsFound.push(contextObj);
    }
	
	var hash = window.location.hash;
	if (hash == '#Result') {
        updateResult();    
	}
}

function checkTopology() 
{
    var queryReq = {}
    queryReq.entities = [{id:'Topology.child-finder', type: 'Topology', isPattern: false}];           
    
    client.queryContext(queryReq).then( function(resultList) {
        console.log(resultList);
        if(resultList && resultList.length > 0) {
            curTopology = resultList[0];
        }
        
        showTopology();
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });      
}


function checkIntent() 
{
    var queryReq = {};
    queryReq.entities = [{type: 'ServiceIntent', isPattern: true}];               
    queryReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'topology=Topology.child-finder'}]}        
    
    client.queryContext(queryReq).then( function(resultList) {
        console.log(resultList);
        if(resultList && resultList.length > 0) {
            curIntent = resultList[0];
            
            //update the current geoscope as well
            geoscope = curIntent.attributes.intent.value.geoscope;            
        }
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
    });      
}

function showTopology() 
{
    $('#info').html('processing topology of this IoT service');
    
    var html = '';
    html += '<div class="input-prepend">';      
    
    if(curTopology) {
        html += '<button id="submitTopology" type="button" class="btn btn-default" disabled>submit</button>';    
    } else {
        html += '<button id="submitTopology" type="button" class="btn btn-default">submit</button>';    
    }

    html += '<input id="loadTopology" type="file" style="display: none;" accept=".json"></input>';      
    html += '</div> ';     
    
    html += '<div><img src="/img/lostchild.jpg"></img></div>';    

    $('#content').html(html);
    
    // associate functions to clickable buttons
    $('#loadTopology').change(loadTopologyFile);    
    $('#submitTopology').click( function() {
        $('#loadTopology').trigger('click');        
    });
}

function submitTopology(topology)
{
    console.log('submit a topology ', topology);
        
    var topologyCtxObj = {};
    
    topologyCtxObj.entityId = {
        id : 'Topology.' + topology.name, 
        type: topology.name,
        isPattern: false
    };
    
    topologyCtxObj.attributes = {};   
    topologyCtxObj.attributes.status = {type: 'string', value: 'disabled'};
    topologyCtxObj.attributes.template = {type: 'object', value: topology};    
    
    client.updateContext(topologyCtxObj).then( function(data) {
        console.log(data);
        
        // update the current topology
        curTopology = topologyCtxObj;           
    }).catch( function(error) {
        console.log('failed to submit the topology');
    });
    
    //disable the submit button
    $('#submitTopology').prop('disabled', true);
}

function loadTopologyFile(evt)
{
    var files = evt.target.files;

    if(files && files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            try {
                var json = JSON.parse(e.target.result);
                submitTopology(json);                
            }catch(ex) {
                alert('error when trying to load topology file');
            }
        }
        reader.readAsText(files[0]);
    }      
}

function showMgt() 
{
    $('#info').html('trigger service topology for data sources in a defined geo-scope');
    
    var html = '';
    html += '<div class="input-prepend">';      
    
    if(curIntent == null) {
        html += '<button id="enableService" type="button" class="btn btn-default">Start</button>';    
        html += '<button id="disableService" type="button" class="btn btn-default" disabled>Stop</button>';                            
    } else {
        html += '<button id="enableService" type="button" class="btn btn-default" disabled>Start</button>';    
        html += '<button id="disableService" type="button" class="btn btn-default">Stop</button>';                            
    }

    html += '<button id="photoSubmitButton" type="button" class="btn btn-default">Select a photo</button>';  
    html += '<input id="lostChildImg" type="file" style="display: none;" accept="image/gif, image/jpeg, image/png"></input>';  
    html += '<span class="caret"></span>';
    html += '<img id="image_upload_preview" src="' + '/photo/lostchild.png?' +  new Date().getTime() + '" alt="target" style="width: 100px; height: 100px"></img>';      
    html += '</div>'; 

//    html += '<div class="input-prepend"><label class="checkbox"><input type="checkbox" id="ScopeUpdating" value="option1">';
//    html += 'automatically updating the search scope</label>';
//    html += '<select id="checkingInterval"><option value="5">5 seconds</option><option value="30">30 seconds</option><option  value="60">1 minute</option><option  value="120">2 minutes</option><option  value="300">5 minutes</option></select>';        
//    html += '</div>';    

//    html += '<div class="input-prepend">';
//    html += '<label>setting the increase of radius</label>';
//    html += '<select id="radiusInterval"><option  value="10000">10000</option><option  value="20000">20000</option><option  value="50000">50000</option></select>';        
//    html += 'meters</div>';    
    
    html += '<div id="map"  style="width: 700px; height: 500px"></div>';                

    $('#content').html(html);        
    
    // associate functions to clickable buttons
    $('#enableService').click(sendIntent);
    $('#disableService').click(cancelIntent);   
    
    $('#photoSubmitButton').click(openFileDialog);    
    $('#lostChildImg').change(photoSelected);            
    
    // show up the map
    showMap();    
}


function sendIntent() 
{
    if(client == null) {         
        console.log('no nearby broker');
        return;
    }
    
    console.log('issue an service intent for this service topology ', curTopology);
    
    // create the intent object
    var topology = curTopology.attributes.template.value    
    
    var intent = {};        
    intent.topology = topology.name;
    intent.priority = {
        'exclusive': false,
        'level': 50
    };    
    intent.qos = "default";    
    intent.geoscope = geoscope;  
            
    // create the intent entity            
    var intentCtxObj = {};        
    intentCtxObj.entityId = { 
        id: 'ServiceIntent.' + uuid(),           
        type: 'ServiceIntent',
        isPattern: false
    };    
    intentCtxObj.attributes = {};   
    intentCtxObj.attributes.status = {type: 'string', value: 'enabled'};
    intentCtxObj.attributes.intent = {type: 'object', value: intent};  
    
    intentCtxObj.metadata = {};               
    intentCtxObj.metadata.topology = {type: 'string', value: curTopology.entityId.id};    
    
    console.log(JSON.stringify(intentCtxObj));
        
    // check if the dynamic task controlling is selected
    //var scopeUpdating = document.getElementById('ScopeUpdating').checked;
    //var checkingInterval = parseInt($('#checkingInterval option:selected').val(), 10);    
    //var radiusInterval = parseInt($('#radiusInterval option:selected').val(), 10);            
        
    client.updateContext(intentCtxObj).then( function(data) {
        console.log(data);  
        curIntent = intentCtxObj;      
        
        // change the button status
		$('#enableService').prop('disabled', true);
		$('#disableService').prop('disabled', false);     
         
        /*
        if (scopeUpdating == true) {
            console.log('start the timer for checking results and updating the search scope')
            checkingTimer = setInterval(onCheckingTimer, checkingInterval * 1000);
            radiusStepDistance = radiusInterval;
        }*/
                           
    }).catch( function(error) {
        console.log('failed to submit the defined intent');
    });            
}


function onCheckingTimer()
{
    console.log("updating the scope if no result is found")
    console.log(personsFound);
    console.log(curIntent);
    console.log(geoscope);        
    
    if (personsFound.length == 0 && curIntent != null && geoscope.scopeType == 'circle') {
        console.log('current radius = ', geoscope.scopeValue.radius)
	        
        // increase the search scope by updating the requirement
        geoscope.scopeValue.radius += radiusStepDistance;
        curIntent.attributes.intent.value.geoscope = geoscope
        
        console.log(curIntent.attributes.intent.value);
        
        client.updateContext(curIntent).then( function(data) {
            console.log('already updated the current requirement with the increased scope');            
            displaySearchScope();
        }).catch( function(error) {
            console.log('failed to update the current requirement');
        }); 	        
    } else {
		// make the camera image to pop up, if the target is founded in the camera
		ctxObj = personsFound[0];
		var cameraDeviceID = "Device.Camera." + ctxObj.attributes.cameraID.value
		var marker = cameraMarkers[cameraDeviceID];				
		if (marker){
			marker.openPopup();
		}		
	}						
}


function cancelIntent() 
{
    if(client == null) {         
        console.log('no nearby broker');
        return;
    }    
    
    console.log('cancel the issued intent for this service topology ', curTopology.entityId.id);
        
    //stop the timer for result checking
    if (checkingTimer != null) {
        console.log('stop the timer for checking results and updating the search scope')        
        clearInterval(checkingTimer);
    }        
    
    var entityid = {
        id : curIntent.entityId.id, 
        type: 'ServiceIntent',
        isPattern: false
    };	    
    
    client.deleteContext(entityid).then( function(data) {
        console.log(data);
        
        curIntent = null;		        
        geoscope = {
            scopeType: "local",
            scopeValue: "local"
        };
        
        personsFound = [];
        
		$('#enableService').prop('disabled', false);
		$('#disableService').prop('disabled', true);		
    }).catch( function(error) {
        console.log('failed to cancel the service intent');
    });        
}


function openFileDialog() {
    $('#lostChildImg').trigger('click');
}

function photoSelected(evt) {
    var files = evt.target.files;

    if(files && files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            var dataURI = e.target.result;
            $('#image_upload_preview').attr('src', dataURI);
            Webcam.params.upload_name = 'lostchild.png';
            Webcam.upload(dataURI,  '/photo', function(code, text) {
                console.log(code);
                console.log(text);
            });
        }
        reader.readAsDataURL(files[0]);
    }    
}

function publishChildInfo()
{
    console.log('publish the information of the lost child, image at ', childphotoURL);
    
    var lostChildCtxObj = {};
    
    lostChildCtxObj.entityId = {
        id : 'Stream.ChildLost.01', 
        type: 'ChildLost',
        isPattern: false
    };
    
    lostChildCtxObj.attributes = {};   
    lostChildCtxObj.attributes.imageURL = {type: 'string', value: childphotoURL};
    lostChildCtxObj.attributes.saveLocation = {type: 'string', value: saveLocation};    
    
    client.updateContext(lostChildCtxObj).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to update the threshold for anomaly detection');
    });    
}


function showTasks() 
{
    $('#info').html('list of running data processing tasks');
    
    if(curTopology == null) {
        $('#content').html('please load the topology first');
        return;        
    }        
        
    var queryReq = {}
    queryReq.entities = [{type:'Task', isPattern: true}];
    queryReq.restriction = {scopes: [{scopeType: 'stringQuery', scopeValue: 'topology=child-finder'}]}    
    
    client.queryContext(queryReq).then( function(taskList) {
        console.log(taskList);
        displayTaskList(taskList);
	num_of_task = taskList.length;
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query task');
    });        
}


function displayTaskList(tasks) 
{
    $('#info').html('list of all tasks for this service topology');

    if(tasks.length == 0) {
        $('#content').html('there is no running task for this topology');
        return;
    }          

    var html = '<table class="table table-striped table-bordered table-condensed">';
   
    html += '<thead><tr>';
    html += '<th>ID</th>';
    html += '<th>worker</th>';
    html += '<th>port</th>';	
    html += '<th>status</th>';		
    html += '</tr></thead>';    

    for(var i=0; i<tasks.length; i++){
        var task = tasks[i];
        html += '<tr>';
        html += '<td>' + task.attributes.id.value + '</td>';		
		html += '<td>' + task.attributes.worker.value + '</td>';
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


function showResult() 
{
    $('#info').html('updated result');
	    
    // table to show the search result
    var html = '';
    
    html += '<table id="searchResult" class="table table-striped table-bordered table-condensed"></table>';            
    html += '<div id="chart"><svg style="height:500px"></svg></div>';    
    
	$('#content').html(html);      

    updateResult();
}


function displayChart(dataset, divID, model, yLabel, yRange)
{
    if(model == 'MULTI-BAR'){
        chart = nv.models.multiBarChart()
                        .x( function(d) { return d[0]} )
                        .y( function(d) { return d[1]} )
                        .stacked(true)
                        .color( d3.scale.category10().range());  
                        
        chart.xAxis.tickFormat( function(d) {
            var t = new Date(d);            
            var time = t.getHours() + ':' + t.getMinutes() + ':' + t.getSeconds();            
            return time;
        });          
    }else if(model == 'MULTI-AREA') {
        chart = nv.models.stackedAreaChart()
                        .x( function(d) { return d[0]} )
                        .y( function(d) { return d[1]} )
                        .clipEdge(true)
						.useInteractiveGuideline(true);

 		chart.xAxis.showMaxMin(false).tickFormat(function(d) { 
			return d3.time.format('%x')(new Date(d)) });                   
    }else if(model == 'PIE-CHART') {
        chart = nv.models.pieChart()
                        .x( function(d) { return d.label} )
                        .y( function(d) { return d.value} )
                        .showLabels(true);      
        console.log('pie chart');                        
    }else if(model == 'DISCRETE-BAR') {
        chart = nv.models.discreteBarChart()
                        .x( function(d) { return d.label  } )
                        .y( function(d) { return d.value  })
                        .staggerLabels(true)
                        .tooltips(false)
                        .showValues(true);        
    } else {        
        chart = nv.models.lineChart()
                        .x( function(d) { return d[0] } )
                        .y( function(d) { return d[1] })
                        .color( d3.scale.category10().range());    
        
        chart.xAxis.tickFormat(function(d) { 
            return d3.time.format('%X')(new Date(d)); 
        });                       
    }
    
    chart.margin(MARGIN);

    if(model != 'PIE-CHART'){    
        chart.yAxis.tickFormat( function(d) { return d } );
        chart.yAxis.axisLabel(yLabel);
        chart.forceY(yRange);
    }
        
    d3.select(divID).datum(dataset)
                           .transition().duration(500)
                           .call(chart);
    
    nv.utils.windowResize(chart.update);                    
}

function updateResult()
{
    var html = '';
    html += '<thead><tr>';
    html += '<th>Time</th>';
    html += '<th>Image</th>';        
    html += '<th>Which Camera</th>';
    html += '<th>Location</th>';
    html += '</tr></thead>';        
    
    for(var i = 0; i < personsFound.length; i++) {
        ctxObj = personsFound[i];
        
        var url = document.createElement('a');
        url.href = ctxObj.attributes.image.value;               
        
        html += '<tr>'; 
        html += '<td>' + ctxObj.attributes.when.value + '</td>';        
        html += '<td><img src="' + url.pathname + '" width="200px"></img></td>';		                
        html += '<td>' + ctxObj.attributes.cameraID.value + '</td>';
        html += '<td>' + JSON.stringify(ctxObj.attributes.where.value) + '</td>';		
        html += '</tr>';	
    }
    
    //update the table content with the received result
    $('#searchResult').empty();        
    $('#searchResult').append(html);        
}

function updateChart(divID, dataset)
{
    d3.select(divID).datum(dataset)
                           .transition().duration(500)
                           .call(chart);
}

function showMap() 
{
    var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var osm = L.tileLayer(osmUrl, {maxZoom: 7, zoom: 7});    
    var map = new L.Map('map', {zoomControl:false, layers: [osm], center: new L.LatLng(35.692221, 138.709059), zoom: 7 });

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
            polygon: true,
            rectangle: true,
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
            
            geoscope.scopeType = 'polygon';
            geoscope.scopeValue = {
                vertices: []
            };
            
            points = geometry.coordinates[0];
            for(i in points){
                geoscope.scopeValue.vertices.push({longitude: points[i][0], latitude: points[i][1]});
            }
            
			console.log(geoscope);            
		} 
		if (type === 'circle') {
            var geometry = layer.toGeoJSON()['geometry'];
            console.log(geometry);
            var radius = layer.getRadius();
            
            geoscope.scopeType = 'circle';
            geoscope.scopeValue = {
                centerLatitude: geometry.coordinates[1],
                centerLongitude: geometry.coordinates[0],
                radius: radius
            }
            
			console.log(geoscope);            
		} 
		if (type === 'polygon') {
            var geometry = layer.toGeoJSON()['geometry'];
            console.log(geometry);
            
            geoscope.scopeType = 'polygon';
            geoscope.scopeValue = {
                vertices: []
            };
            
            points = geometry.coordinates[0];
            for(i in points){
                geoscope.scopeValue.vertices.push({longitude: points[i][0], latitude: points[i][1]});
            }
            
			console.log(geoscope);            
		}
        
        drawnItems.addLayer(layer);
    });  


    // show edge nodes on the map
    displayEdgeNodeOnMap(map);
    
    // show all devices on the map
    displayDeviceOnMap(map);  
        
    // remember the created map
    curMap = map;
    
    // display the current search scope
    displaySearchScope();    
}

function displayEdgeNodeOnMap(map)
{
    var queryReq = {}
    queryReq.entities = [{type:'Worker', isPattern: true}];
    client.queryContext(queryReq).then( function(edgeNodeList) {
        console.log(edgeNodeList);

        var edgeIcon = L.icon({
            iconUrl: '/img/gateway.png',
            iconSize: [48, 48]
        });      
        
        for(var i=0; i<edgeNodeList.length; i++){
            var worker = edgeNodeList[i];    
            
            latitude = worker.attributes.location.value.latitude;
            longitude = worker.attributes.location.value.longitude;
            edgeNodeId = worker.entityId.id;
            
            var marker = L.marker(new L.LatLng(latitude, longitude), {icon: edgeIcon});
			marker.nodeID = edgeNodeId;
            marker.addTo(map).bindPopup(edgeNodeId);
		    marker.on('click', showRunningTasks);
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


function displaySearchScope()
{
    console.log(geoscope);            
    if(geoscope != null) {
        switch (geoscope.scopeType) {
            case 'circle': 
                L.circle([geoscope.scopeValue.centerLatitude, geoscope.scopeValue.centerLongitude], geoscope.scopeValue.radius).addTo(curMap);                                
                break;
            case 'polygon':
                var points = [];
                for(var i=0; i<geoscope.scopeValue.vertices.length; i++){
                    points.push(new L.LatLng(geoscope.scopeValue.vertices[i].latitude, geoscope.scopeValue.vertices[i].longitude))
                }            
                L.polygon(points).addTo(curMap);
                break;
        }                
    }   
}

function displayDeviceOnMap(map)
{
    var queryReq = {}
    queryReq.entities = [{id:'Device.*', isPattern: true}];
    client.queryContext(queryReq).then( function(devices) {
        console.log(devices);
        
        for(var i=0; i<devices.length; i++){
            var device = devices[i];    
            
           if (device.attributes.iconURL != null){
                iconImag = device.attributes.iconURL.value;            
                var edgeIcon = L.icon({
                    iconUrl: iconImag,
                    iconSize: [48, 48]
                });                  
                
                latitude = device.metadata.location.value.latitude;
                longitude = device.metadata.location.value.longitude;
                deviceId = device.entityId.id;
                
                var marker = L.marker(new L.LatLng(latitude, longitude), {icon: edgeIcon});
                if(device.entityId.type == 'Camera') { 
    				var imageURL = device.attributes.url.value;
    				marker.addTo(map).bindPopup("<img src=" + window.location.origin + "/proxy?url="  + imageURL + "></img>");                             				
    			    cameraMarkers[deviceId] = marker;				
                } else {
                    marker.addTo(map).bindPopup(deviceId);                 
                }               
            }           
        }                            
    }).catch(function(error) {
        console.log(error);
        console.log('failed to query context');
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


});



