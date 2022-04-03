$(function() {

    // initialization  
    var handlers = {}
    var geoscope = {
        scopeType: "local",
        scopeValue: "local"
    };

    var curTopology = null;
    var curIntent = null;
    var curResult = [];

    var category_dataset = [{ key: '#health_risk', values: [] }];
    var chart;
    var MARGIN = { top: 30, right: 20, bottom: 60, left: 80 };


    addMenuItem('Topology', showTopology);
    addMenuItem('Management', showMgt);
    addMenuItem('Tasks', showTasks);
    addMenuItem('Alerts', showResult);
    addMenuItem('Prediction', predictionResult);

    //connect to the socket.io server via the NGSI proxy module
    var ngsiproxy = new LDNGSIProxy();
    ngsiproxy.setNotifyHandler(handleNotify);

    //connect to the broker
    var ldclient = new NGSILDclient(config.LdbrokerURL);
    var client = new NGSI10Client(config.brokerURL);
    //var clientDes = new NGSI10Client('./internal');

    subscribeResult();
    checkTopology();
    checkIntent();
    showTopology();

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

    //start a timer to send the rule set periodically
    //function startUpdateTimer() {
    //    setInterval(publishThreshold, 2000);
    //}

    function subscribeResult() {
        var subscribeCtxReq = {};
        subscribeCtxReq.type = "Subscription"

        subscribeCtxReq.entities = [{ type: 'ldStat_health' }];

        subscribeCtxReq.notification = {}
        endPoint = {}
        subscribeCtxReq.notification.format = "normalized"

        endPoint.uri = 'http://' + config.agentIP + ':' + config.ldAgentPort + '/notifyContext';
        endPoint.accept = "application/ld+json"
        subscribeCtxReq.notification.endpoint = endPoint
        console.log("subscribeCtxReq data for topology",subscribeCtxReq)
        ldclient.subscribeContext(subscribeCtxReq).then(function(subscriptionId) {
            console.log("SubscriptionID",subscriptionId);
            ngsiproxy.reportSubID(subscriptionId);
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });

    }

    function handleNotify(contextObj) {
        curResult.push(contextObj);
	
        if (contextObj.hasOwnProperty("timmer") && contextObj.hasOwnProperty("counter")) {
            var time = contextObj.timmer.value;
            var num = contextObj.counter.value;
            var point = [time, num];
            category_dataset[0].values.push(point);
            var hash = window.location.hash;
            console.log(hash);
            if (hash == '#Alerts') {
                updateChart('#chart svg', category_dataset);
            }
        }
    }

    function checkTopology() {
        fetch('/topology/Heart_Health_Predictor').then(res => res.json()).then(topology => {
            if (Object.values(topology).length > 0) { //non-empty       
                curTopology = topology;   
            }
            showTopology();                   
        }).catch(function(error) {
           console.log(error);
           console.log('failed to fetch the required topology');
        });        
    }

    function checkIntent() {
	fetch('/intent/topology/Heart_Health_Predictor').then(res => res.json()).then(intent => {
            if (Object.values(intent).length > 0) { //non-empty 
                curIntent = intent; 
                //update the current geoscope as well
                geoscope = intent[0].geoscope;                 
            }              
        }).catch(function(error) {
           console.log(error);
           console.log('failed to fetch the required intent');
        });    
    }

    function showTopology() {
        $('#info').html('IoT service topology');

        var html = '';
        html += '<div class="input-prepend">';

        if (curTopology) {
            html += '<button id="submitTopology" type="button" class="btn btn-default" disabled>submit</button>';
        } else {
            html += '<button id="submitTopology" type="button" class="btn btn-default">submit</button>';
        }

        html += '<input id="loadTopology" type="file" style="display: none;" accept=".json"></input>';
        html += '</div> ';

        html += '<div><img src="/img/heart_health.png"></img></div>';

        $('#content').html(html);

        // associate functions to clickable buttons
        $('#loadTopology').change(loadTopologyFile);
        $('#submitTopology').click(function() {
            $('#loadTopology').trigger('click');
        });
    }

    function submitTopology(topology) {
        console.log('submit a topology ', topology);

        var topologyCtxObj = {};

        topologyCtxObj.attributes = topology ;
	fetch('/topology/Heart_Health_Predictor').then(res => res.json()).then(topologyCtxObj => {
                curTopology = topologyCtxObj;
        }).catch(function(error) {
           console.log(error);   
        $('#submitTopology').prop('disabled', true);
        });
}

    function loadTopologyFile(evt) {
        var files = evt.target.files;

        if (files && files[0]) {
            var reader = new FileReader();
            reader.onload = function(e) {
                try {
                    json = JSON.parse(e.target.result);
                    submitTopology(json);
                } catch (ex) {
                    alert('error when trying to load topology file');
                }
            }
            reader.readAsText(files[0]);
        }
    }

    function showMgt() {
        $('#info').html('trigger service topology for data sources in a defined geo-scope');

        var html = '';
        html += '<div class="input-prepend">';

        if (curIntent == null) {
            html += '<button id="enableService" type="button" class="btn btn-default">Start</button>';
            html += '<button id="disableService" type="button" class="btn btn-default" disabled>Stop</button>';
        } else {
            html += '<button id="enableService" type="button" class="btn btn-default" disabled>Start</button>';
            html += '<button id="disableService" type="button" class="btn btn-default">Stop</button>';
        }

        html += '</div> ';

        html += '<div id="map"  style="width: 700px; height: 500px"></div>';

        $('#content').html(html);

        // associate functions to clickable buttons
        $('#enableService').click(sendIntent);
        $('#disableService').click(cancelIntent);

        // show up the map
        showMap();
    }

    function sendIntent() {

        console.log('issue an service intent for this service topology ', curTopology);
	var topology = curTopology.name;
	var uid = uuid();
        var sid = 'ServiceIntent.' + uid;
	var serviceintent = {};
	serviceintent.id = sid;
	//serviceintent.stype = "SYN";
	serviceintent.qos = "default";
	serviceintent.geoscope = geoscope;
	serviceintent.priority = {
            'exclusive': false,
            'level': 50
        };
        serviceintent.topology = topology;
	serviceintent.action = 'ADD';
	//console.log(JSON.stringify(serviceintent));
	fetch("/intent", {
            method: "POST",
            headers: {
                Accept: "application/json+ld",
                "Content-Type": "application/json"
            },
	    body: JSON.stringify(serviceintent)
        })
        .then(function(data) {
            console.log(data);
            curIntent = serviceintent;
	    console.log(JSON.stringify(curIntent));
            // change the button status
            $('#enableService').prop('disabled', true);
            $('#disableService').prop('disabled', false);
        }).catch(function(error) {
            console.log('failed to submit the defined intent');
        });
    }

    function cancelIntent() {

        console.log('cancel the issued intent for this service topology ', curTopology.name);
	fetch("/intent/" + curIntent.id, {
            method: "DELETE"
        })	
        .then(function(data) {
	    console.log(data);
            curIntent = null;
            geoscope = {
                scopeType: "local",
                scopeValue: "local"
            };
            $('#enableService').prop('disabled', false);
            $('#disableService').prop('disabled', true);
        }).catch(function(error) {
            console.log('failed to cancel the service intent');
        });
    }


    function showTasks() {
        $('#info').html('list of running data processing tasks');

        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];
        //queryReq.restriction = { scopes: [{ scopeType: 'stringQuery', scopeValue: 'topology=anomaly-detection' }] }

        client.queryContext(queryReq).then(function(taskList) {
            console.log('Task List: ' + taskList);
            displayTaskList(taskList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query task');
        });
    }


    function displayTaskList(tasks) {
        $('#info').html('list of all tasks for this service topology');
	console.log('No of tasks : '+tasks.length)	
        if (tasks.length == 0) {
            $('#content').html('there is no running task for this topology');
            return;
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Service</th>';
        html += '<th>Task</th>';
        html += '<th>worker</th>';
        html += '<th>port</th>';
        html += '<th>status</th>';
        html += '</tr></thead>';

        for (var i = 0; i < tasks.length; i++) {
            var task = tasks[i];
            html += '<tr>';
            html += '<td>' + task.attributes.id.value + '</td>';
            html += '<td>' + task.attributes.service.value + '</td>';
            html += '<td>' + task.attributes.task.value + '</td>';
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


    function predictionResult() {
        $('#info').html('list of prediction data from all Heart Sensors');
        
        var queryReq2 = {};
        queryReq2.type = 'Query';
        queryReq2.entities = [{ type : 'prediction' }];
	console.log(JSON.stringify(queryReq2));
        ldclient.queryContext(queryReq2).then(function(entityList) {
            console.log(entityList);
            displayPredictionList(entityList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query predictions');
        });
    }

    function displayPredictionList(entity) {
	$('#info').html('list of predictions for this service topology');
        console.log('No. of predictions: '+entity[0].length);
	console.log(JSON.stringify(entity));

	if (entity[0].length == 0) {
            $('#content').html('there are no predictions for this topology');
            return;
        }
       
        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Sensor ID</th>';
        html += '<th>Prediction</th>';
        html += '<th>Alert</th>';
        html += '<th>Time</th>';
        html += '</tr></thead>';

        for (var i = 0; i < entity.length; i++) {
            var taskresult = entity[i];
            console.log(taskresult);
            console.log(taskresult.length);
        for (var j = 0; j < taskresult.length; j++){
	    var task = taskresult[j];
            console.log(task);
            html += '<tr>';
            html += '<td>' + task.id + '</td>';
            html += '<td>' + task.Analysis.value + '</td>';
            if (task.alert.value == "1") {
                html += '<td><font color="red">' + task.alert.value + '</font></td>';
            } else {
                html += '<td><font color="green">' + task.alert.value + '</font></td>';
            }

            html += '<td>' + task.Time.value + '</td>';

            html += '</tr>';
        }
       }

        html += '</table>';

        $('#content').html(html);

    }
       
    function showResult() {
        $('#info').html('statistical result');

        var html = '<div id="chart"><svg style="height:500px"></svg></div>';
        $('#content').html(html);

        displayChart(category_dataset, '#chart svg', '', '# of detected heath risk events', [0, ]);
    }


    function displayChart(dataset, divID, model, yLabel, yRange) {
        if (model == 'MULTI-BAR') {
            chart = nv.models.multiBarChart()
                .x(function(d) { return d[0] })
                .y(function(d) { return d[1] })
                .stacked(true)
                .color(d3.scale.category10().range());

            chart.xAxis.tickFormat(function(d) {
                var t = new Date(d);
                var time = t.getHours() + ':' + t.getMinutes() + ':' + t.getSeconds();
                return time;
            });
        } else if (model == 'MULTI-AREA') {
            chart = nv.models.stackedAreaChart()
                .x(function(d) { return d[0] })
                .y(function(d) { return d[1] })
                .clipEdge(true)
                .useInteractiveGuideline(true);

            chart.xAxis.showMaxMin(false).tickFormat(function(d) {
                return d3.time.format('%x')(new Date(d))
            });
        } else if (model == 'PIE-CHART') {
            chart = nv.models.pieChart()
                .x(function(d) { return d.label })
                .y(function(d) { return d.value })
                .showLabels(true);
            console.log('pie chart');
        } else if (model == 'DISCRETE-BAR') {
            chart = nv.models.discreteBarChart()
                .x(function(d) { return d.label })
                .y(function(d) { return d.value })
                .staggerLabels(true)
                .tooltips(false)
                .showValues(true);
        } else {
            chart = nv.models.lineChart()
                .x(function(d) { return d[0] })
                .y(function(d) { return d[1] })
                .color(d3.scale.category10().range());

            chart.xAxis.tickFormat(function(d) {
                return d3.time.format('%X')(new Date(d));
            });
        }

        chart.margin(MARGIN);

        if (model != 'PIE-CHART') {
            chart.yAxis.tickFormat(function(d) { return d });
            chart.yAxis.axisLabel(yLabel);
            chart.forceY(yRange);
        }

        d3.select(divID).datum(dataset)
            .transition().duration(500)
            .call(chart);

        nv.utils.windowResize(chart.update);
    }

    function updateChart(divID, dataset) {
        d3.select(divID).datum(dataset)
            .transition().duration(500)
            .call(chart);
    }
    function showMap() {
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

                geoscope.scopeType = 'polygon';
                geoscope.scopeValue = {
                    vertices: []
                };

                points = geometry.coordinates[0];
                for (i in points) {
                    geoscope.scopeValue.vertices.push({ longitude: points[i][0], latitude: points[i][1] });
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
                for (i in points) {
                    geoscope.scopeValue.vertices.push({ longitude: points[i][0], latitude: points[i][1] });
                }

                console.log(geoscope);
            }

            drawnItems.addLayer(layer);
        });

        // show edge nodes on the map
        displayEdgeNodeOnMap(map);

        // show all devices on the map
        displayDeviceOnMap(map)

        // display the defined scope
        displaySearchScope(map)
    }


    function displaySearchScope(map) {
        console.log(geoscope);
        if (geoscope != null) {
            switch (geoscope.scopeType) {
                case 'circle':
                    L.circle([geoscope.scopeValue.centerLatitude, geoscope.scopeValue.centerLongitude], geoscope.scopeValue.radius).addTo(map);
                    break;
                case 'polygon':
                    var points = [];
                    for (var i = 0; i < geoscope.scopeValue.vertices.length; i++) {
                        points.push(new L.LatLng(geoscope.scopeValue.vertices[i].latitude, geoscope.scopeValue.vertices[i].longitude))
                    }
                    L.polygon(points).addTo(map);
                    break;
            }
        }
    }

    function displayEdgeNodeOnMap(map) {
        var queryReq = {}
        queryReq.entities = [{ type: 'Worker', isPattern: true }];

client.queryContext(queryReq).then(function(edgeNodeList) {
            var edgeIcon = L.icon({
                iconUrl: '/img/gateway.png',
                iconSize: [48, 48]
            });

            for (var i = 0; i < edgeNodeList.length; i++) {
                var worker = edgeNodeList[i];

                latitude = worker.attributes.location.value.latitude;
                longitude = worker.attributes.location.value.longitude;
                edgeNodeId = worker.entityId.id;

                var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
                marker.nodeID = edgeNodeId;
                marker.addTo(map).bindPopup(edgeNodeId);
                marker.on('click', showRunningTasks);
            }

        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function showRunningTasks() {
        var clickMarker = this;

        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];
        queryReq.restriction = { scopes: [{ scopeType: 'stringQuery', scopeValue: 'worker=' + clickMarker.nodeID }] }

        client.queryContext(queryReq).then(function(tasks) {
            console.log(tasks);
            var content = "";
            for (var i = 0; i < tasks.length; i++) {
                var task = tasks[i];

                if (task.attributes.status.value == "paused") {
                    content += '<font color="red">' + task.attributes.id.value + '</font><br>';
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

    function displayDeviceOnMap(map) {
        var queryReq = {}
        queryReq.entities = [{ id: 'Device.*', isPattern: true }];
        client.queryContext(queryReq).then(function(devices) {
            for (var i = 0; i < devices.length; i++) {
                var device = devices[i];
                console.log(device);

                if (device.attributes.iconURL != null) {
                    iconImag = device.attributes.iconURL.value;
                    var edgeIcon = L.icon({
                        iconUrl: iconImag,
                        iconSize: [48, 48]
                    });

                    latitude = device.metadata.location.value.latitude;
                    longitude = device.metadata.location.value.longitude;
                    deviceId = device.entityId.id;

                    var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
                    marker.addTo(map).bindPopup(deviceId);
                }
            }

        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
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

});
