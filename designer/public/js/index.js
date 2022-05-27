'use strict';

$(function () {

    // initialize the menu bar
    var handlers = {}

    // location of new device
    var locationOfNewDevice = null;
    // icon image for device registration
    var iconImage = null;
    var iconImageFileName = null;
    // content image for camera devices
    var contentImage = null;
    var contentImageFileName = null;

    // client to interact with IoT Broker
    var client = new NGSI10Client(config.brokerURL);
    var ldclient = new NGSILDclient(config.LdbrokerURL);        

    var ldDeviceObj = []

    addMenuItem('Architecture', showArch);
    addMenuItem('Discovery', showDiscovery);
    addMenuItem('Broker', showBrokers);
    addMenuItem('Master', showMaster);
    addMenuItem('Worker', showWorkers);
    addMenuItem('Device', showDevices);
    addMenuItem('Subscriptions', showSubscriptions);
    addMenuItem('uService', showEndPointService);
    addMenuItem('Task', showTasks);    
    addMenuItem('Entity', showStreams);

    showArch();

    $(window).on('hashchange', function () {
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

    function showArch() {
        $('#info').html('Elastic IoT Platform with Standard-based Edge Computing');
        $('#content').html('<img width="50%" height="50%" src="/img/arch.jpg"></img>');
    }


    function showDiscovery() {
        $('#info').html('information of IoT Discovery');

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Type</th>';
        html += '<th>URL</th>';
        html += '</tr></thead>';

        html += '<tr><td>IoT Discovery</td><td>' + config.discoveryURL + '</td></tr>';

        html += '</table>';

        $('#content').html(html);
    }

    function showBrokers() {
        $('#info').html('list of all IoT Brokers');

        fetch('/info/broker').then(res => res.json()).then(brokers => {
            displayBrokerList(brokers);
        });
    }


    function displayBrokerList(brokers) {
        if (brokers == null || brokers.length == 0) {
            $('#content').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>BrokerURL</th>';
        html += '</tr></thead>';

        for (var i = 0; i < brokers.length; i++) {
            var broker = brokers[i];
            html += '<tr><td>' + broker.id + '</td><td>' + broker.myURL + '</td></tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }

    function showMaster() {
        $('#info').html('list of all topology masters');

        fetch('/info/master').then(res => res.json()).then(masters => {
            displayMasterList(masters)
        });
    }

    function displayMasterList(masters) {

        if (masters == null || masters.length == 0) {
            $('#content').html('there is no topology master running');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Location</th>';
        html += '<th>Agent</th>';
        html += '</tr></thead>';

        for (var i = 0; i < masters.length; i++) {
            var master = masters[i];
            console.log(master);
            html += '<tr><td>' + master.id + '</td><td>' + JSON.stringify(master.location) + '</td><td>' + master.agent + '</td></tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }


    function showEndPointService() {
        $('#info').html('list of all available end-point services');

        var queryReq = {}
        queryReq.entities = [{ type: 'uService', isPattern: true }];
        client.queryContext(queryReq).then(function (serviceList) {
            displayServiceList(serviceList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function displayServiceList(serviceList) {
        if (serviceList == null || serviceList.length == 0) {
            $('#content').html('there is no endpoint service running');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Access Point</th>';
        html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < serviceList.length; i++) {
            var service = serviceList[i];
            html += '<tr><td>' + service.entityId.id + '</td>';
            html += '<td>' + service.metadata.IP.value + ":" + service.metadata.port.value + '</td>';
            html += '<td>' + JSON.stringify(service.metadata) + '</td></tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }

    function showWorkers() {
        $('#info').html('show all edge nodes on the map');

        var html = '<table class="table table-striped table-bordered table-condensed" id="workerList"></table>';
        html += '<div id="map"  style="width: 700px; height: 500px"></div>';

        $('#content').html(html);

        fetch('/info/worker').then(res => res.json()).then(workers => {
            displayWorkerList(workers)
            displayWorkerOnMap(workers)
        });
    }

    function displayWorkerList(workerList) {
        var html = '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>location</th>';
        html += '<th>capacity</th>';
        html += '<th># of tasks</th>';
        html += '</tr></thead>';

        for (var i = 0; i < workerList.length; i++) {
            var worker = workerList[i];

            var wid = worker.id;

            html += '<tr>';
            html += '<td>' + wid + '</td>';
            html += '<td>' + JSON.stringify(worker.location) + '</td>';
            html += '<td>' + worker.capacity + '</td>';
            html += '<td>' + worker.workload + '</td>';
            html += '</tr>';
        }

        $('#workerList').html(html);
    }


    function displayWorkerOnMap(workerList) {
        var curMap = showMap();

        var edgeIcon = L.icon({
            iconUrl: '/img/gateway.png',
            iconSize: [48, 48]
        });

        for (var i = 0; i < workerList.length; i++) {
            var worker = workerList[i];
            console.log(worker)

            var latitude = worker.location.latitude;
            var longitude = worker.location.longitude;
            var edgeNodeId = worker.id;

            var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
            marker.nodeID = edgeNodeId;
            marker.addTo(curMap).bindPopup(edgeNodeId);
            marker.on('click', showRunningTasks);
        }
    }

    function showRunningTasks() {
        var clickMarker = this;

        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];
        queryReq.restriction = { scopes: [{ scopeType: 'stringQuery', scopeValue: 'worker=' + clickMarker.nodeID }] }

        client.queryContext(queryReq).then(function (tasks) {
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
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query task');
        });
    }


    function showDevices() {
        $('#info').html('list of all registered devices');

        var html = '<div style="margin-bottom: 10px;"><button id="addNewDevice" type="button" class="btn btn-primary">add</button></div>';

        html += '<div id="deviceList"></div>';

        $('#content').html(html);

        updateDeviceList();

        $("#addNewDevice").click(function () {
            deviceRegistration();
        });
    }

    function updateDeviceList() {
        fetch('/device').then(res => res.json()).then(devices => {
            console.log("the list of all registered devices");
            console.log(devices);
            var deviceList = Object.values(devices);
            displayDeviceList(deviceList);            
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of permanent subscriptions');
        });       
    }

    function displayDeviceList(devices) {
        if (devices == null || devices.length == 0) {
            $('#deviceList').html('');
            return
        }
        
        var html = '<table class="table table-striped table-bordered table-condensed">';
        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Device Type</th>';
        html += '<th>Attributes</th>';
        html += '<th>DomainMetadata</th>';
        html += '<th>Action</th>';
        html += '</tr></thead>';

        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];
            html += '<tr>';
            html += '<td>' + device.id + '</td>';
            html += '<td>' + device.type + '</td>';
            html += '<td>' + JSON.stringify(device.attributes) + '</td>';
            html += '<td>' + JSON.stringify(device.metadata) + '</td>';
            html += '<td>'
            html += '<button id="DELETE-' + device.id + '" type="button" class="btn btn-primary">Delete</button>';
            html += '</td>'
            html += '</tr>';
        }

        html += '</table>';
        $('#deviceList').html(html);

        //associate a click handler to generate device profile on request
        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];

            var deleteButton = document.getElementById('DELETE-' + device.id);
            deleteButton.onclick = function (deviceID) {
                return function() {
                    console.log("delete the device ", deviceID);
                    removeDevice(deviceID);
                };                
            }(device.id);
        }
    }

    function removeDevice(deviceID) {
        fetch("/device/" + deviceID, {
            method: "DELETE"
        })
        .then(response => {
            console.log("delete a registered device: ", response.status)
            updateDeviceList();
        })
        .catch(err => console.log(err));           
    }


    function deviceRegistration() {
        $('#info').html('to register a new IoT device');
        
        var html = ''

        html += '<div class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device Protocol(*)</label>';
        html += '<div class="controls"><select id="deviceProtocol"><option>NGSI-v1</option><option>NGSI-LD</option><option>MQTT</option></select></div>'
        html += '</div>';  
        
        html += '<div class="control-group" style="display:none;" id="mqttbrokerInfo">';        
        html += '<label class="control-label" for="input01">MQTTBroker</label>';
        html += '<div class="controls"><input class="input-xlarge" id="mqttbrokerURL"></div>'        
        html += '</div>';                       

        html += '<div class="control-group"><label class="control-label" for="input01">Device Type(*)</label>';
        html += '<div class="controls"><select id="deviceType"><option>Temperature</option><option>PowerPanel</option><option>Camera</option><option>Alarm</option><option>HOPU</option></select></div>'
        html += '</div>';

        html += '<div class="control-group" style="display:none;" id="topicInfo"><label class="control-label" for="input01">Topic</label>';
        html += '<div class="controls"><textarea id="topic"></textarea></div>'        
        html += '</div>';

        html += '<div class="control-group" style="display:none;" id="mappingInfo"><label class="control-label" for="input01">Attribute-Mappings</label>';
        html += '<div class="controls"><textarea id="mappingRules"></textarea></div>'        
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device ID(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="deviceID">';
        html += '<span>  </span><button id="autoIDGenerator" type="button" class="btn btn-primary">Autogen</button>';
        html += '</div></div>';


        html += '<div class="control-group"><label class="control-label" for="input01">Icon Image</label>';
        html += '<div class="controls"><input class="input-file" id="iconImage" type="file" accept="image/png"></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Camera Image</label>';
        html += '<div class="controls"><input class="input-file" id="imageContent" type="file" accept="image/png"></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Location(*)</label>';
        html += '<div class="controls"><div id="map"  style="width: 500px; height: 400px"></div></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01"></label>';
        html += '<div class="controls"><button id="submitRegistration" type="button" class="btn btn-primary">Register</button>';
        html += '</div></div>';

        html += '</fieldset></div>';

        $('#content').html(html);

        // associate functions to clickable buttons
        $('#submitRegistration').click(function () {
            registerNewDevice()
        });

        $('#autoIDGenerator').click(autoIDGenerator);

        $('#deviceProtocol').change(function () {
            var dType = $('#deviceProtocol option:selected').val();
            if (dType == "MQTT") {                
                $('#mqttbrokerInfo').show();
            } else {
                $('#mqttbrokerInfo').hide();
            }
        });

        $('#deviceType').change(function () {
            var dType = $('#deviceType option:selected').val();
            if (dType == "HOPU") {                
                $('#mappingInfo').show();
                $('#topicInfo').show();                
            } else {
                $('#mappingInfo').hide();
                $('#topicInfo').hide();                
            }
        });        

        $('#iconImage').change(function () {
            readIconImage(this);
        });
        $('#imageContent').change(function () {
            readContentImage(this);
        });

        showWorkerList();
    }

    function showWorkerList() {
        fetch('/info/worker').then(res => res.json()).then(workers => {
            displayWorkerOnMap(workers)
        });        
    }

    function readIconImage(input) {
        console.log('read icon image');
        if (input.files && input.files[0]) {
            var reader = new FileReader();
            reader.onload = function (e) {
                //var filename = $('#image_file').val();
                iconImage = e.target.result;
            }
            reader.readAsDataURL(input.files[0]);
            iconImageFileName = input.files[0].name;
        }
    }

    function readContentImage(input) {
        console.log('read content image');
        if (input.files && input.files[0]) {
            var reader = new FileReader();
            reader.onload = function (e) {
                contentImage = e.target.result;
            }
            reader.readAsDataURL(input.files[0]);
            contentImageFileName = input.files[0].name;
        }
    }

    function registerNewDevice() {
        console.log('register a new device');

        // take the inputs    
        var id = $('#deviceID').val();
        var protocol = $('#deviceProtocol').val();
        var type = $('#deviceType option:selected').val();

        if (id == '' || type == '' || locationOfNewDevice == null) {
            alert('please provide the required inputs');
            return;
        }
        
        var mappingRules = $('#mappingRules').val();        
        if (mappingRules != null && mappingRules != '') {
            try {  
                const json = JSON.parse(mappingRules);  
            } catch (e) {  
                alert('please provide the required inputs for the mapping rules');
                return;
            }            
        }           
        
        //upload the icon image
        if (iconImage != null) {
            Webcam.params.upload_name = iconImageFileName;
            Webcam.upload(iconImage, '/photo', function (code, text) {
                console.log(code);
                console.log(text);
            });
        } else {
            switch (type) {
                case "PowerPanel":
                    iconImageFileName = 'shop.png';
                    break;
                case "Camera":
                    iconImageFileName = 'camera.png';
                    break;
                default:
                    iconImageFileName = 'defaultIcon.png';
                    break;
            }
        }

        // if the device is pull-based, publish a stream entity with its provided URL as well        
        if (contentImage != null) {
            Webcam.params.upload_name = contentImageFileName;
            Webcam.upload(contentImage, '/photo', function (code, text) {
                console.log(code);
                console.log(text);
            });
        }

        //register a new device
        var newDeviceObject = {};
        
        newDeviceObject.id = id;
        newDeviceObject.type = type;
        
        newDeviceObject.attributes = {};
                
        newDeviceObject.attributes.protocol = { type: 'string', value: protocol };         
        
        var mqttbrokerURL = $('#mqttbrokerURL').val();
        if (mqttbrokerURL != null && mqttbrokerURL != '') {
            newDeviceObject.attributes.mqttbroker = { type: 'string', value: mqttbrokerURL };          
        }
        
        if (protocol == 'MQTT' && type == 'HOPU') {
            //var topic = '/api/' + id + '/attrs';            
            var topic = $('#topic').val();                    
            newDeviceObject.attributes.topic = { type: 'string', value: topic };          
        }        
        
        if (mappingRules != null && mappingRules != '') {
            newDeviceObject.attributes.mappings =  { type: 'object', value: JSON.parse(mappingRules) };
        }        

        var url = 'http://' + config.agentIP + ':' + config.webSrvPort + '/photo/' + contentImageFileName;
        newDeviceObject.attributes.url = { type: 'string', value: url };
        newDeviceObject.attributes.iconURL = { type: 'string', value: '/photo/' + iconImageFileName };
                
        if (type == "PowerPanel") {
            newDeviceObject.attributes.usage = { type: 'integer', value: 20 }; 
            newDeviceObject.attributes.shop =  { type: 'string', value: id }; 
        }                

        newDeviceObject.metadata = {};
        newDeviceObject.metadata.location = {
            type: 'point',
            value: { 'latitude': locationOfNewDevice.lat, 'longitude': locationOfNewDevice.lng }
        };

        if (type == "PowerPanel") {
            newDeviceObject.metadata.shop = {
                type: 'string',
                value: id
            };
        } else if (type == "Camera") {
            newDeviceObject.metadata.cameraID = {
                type: 'string',
                value: id
            };
        }
        
        console.log(newDeviceObject);
        
        fetch("/device", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify(newDeviceObject)
        })
        .then(response => {
            console.log("submit a new device: ", response.status)
            showDevices();
        })
        .catch(err => console.log(err));         
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


    function autoIDGenerator() {
        var id = uuid();
        $('#deviceID').val(id);
    }

    function showMap() {
        var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
            osm = L.tileLayer(osmUrl, { maxZoom: 18, zoom: 7 });
        try {
            map.remove();
        } catch (err) {

        }
        var map = new L.Map('map', { layers: [osm], center: new L.LatLng(38.018837048090326, -1.1715629779027177), zoom: 15});
        var drawnItems = new L.FeatureGroup();
        map.addLayer(drawnItems);

        var cameraIcon = L.icon({
            iconUrl: '/img/location.png',
            iconSize: [48, 48]
        });

        var drawControl = new L.Control.Draw({
            draw: {
                position: 'topleft',
                polyline: false,
                polygon: false,
                rectangle: false,
                circle: false,
                marker: {
                    zIndexOffset: 2000,
                    repeatMode: true,
                    icon: cameraIcon
                }
            },
            edit: false
        });
        map.addControl(drawControl);

        map.on('draw:created', function (e) {
            var type = e.layerType, layer = e.layer;

            if (type === 'marker') {
                console.log(layer.getLatLng());
                locationOfNewDevice = layer.getLatLng();
            }

            drawnItems.clearLayers();
            drawnItems.addLayer(layer);
        });

        return map;
    }

    function showSubscriptions() {
        $('#info').html('list of all permanent subscriptions');

        var html = '<div style="margin-bottom: 10px;"><button id="addSubscription" type="button" class="btn btn-primary">add</button></div>';

        html += '<div id="subscriptionList"></div>';

        $('#content').html(html);

        $("#addSubscription").click(function () {
            subscriptionRegistration();
        });

        fetch('/subscription').then(res => res.json()).then(subscriptions => {
            console.log("the list of permanent subscriptions");
            console.log(subscriptions);
            var subscriptionList = Object.values(subscriptions);
            displaySubscriptionList(subscriptions);            
        }).catch(function(error) {
            console.log(error);
            console.log('failed to fetch the list of permanent subscriptions');
        });  
    }
    
    function displaySubscriptionList(subscriptions) {
        if (subscriptions == null || Object.values(subscriptions).length == 0) {
            $('#subscriptionList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Subscription ID</th>';
        html += '<th>Entity Type</th>';
        html += '<th>Destination Broker</th>';
        html += '<th>Reference URL</th>';
        html += '<th>Tenant</th>';        
        html += '<th>Action</th>';
        html += '</tr></thead>';

        for(var sid in subscriptions) {
            var subscription = subscriptions[sid];

            html += '<tr>';
            html += '<td>' + sid + '</td>';            
            html += '<td>' + subscription.entity_type + '</td>';
            html += '<td>' + subscription.destination_broker + '</td>';
            html += '<td>' + subscription.reference_url + '</td>';
            html += '<td>' + subscription.tenant + '</td>';

            html += '<td><button id="delete-' + sid + '" type="button" class="btn btn-primary btn-separator">delete</button></td>';

            html += '</tr>';
        }

        html += '</table>';

        $('#subscriptionList').html(html);

        // associate a click handler to the editor button
        for(var sid in subscriptions) {
            var deleteButton = document.getElementById('delete-' + sid);
            deleteButton.onclick = function(mySubscriptionID) {
                return function() {
                    console.log("delete buttion ", mySubscriptionID);
                    deleteSubscription(mySubscriptionID);
                };
            }(sid);                   
        }        
    }
    
    function deleteSubscription(subscriptionID) {                
        fetch("/subscription/" + subscriptionID, {
            method: "DELETE"
        })
        .then(response => {
            console.log("delete a registered subscription: ", response.status)
            showSubscriptions();
        })
        .catch(err => console.log(err));   
    }

    function subscriptionRegistration() {
        $('#info').html('to register a new permanent subscription');
        
        var html = ''

        html += '<div class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Entity Type(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="entityType"></div>';
        html += '</div>';  

        html += '<div class="control-group"><label class="control-label" for="input01">Destination Broker (*)</label>';
        html += '<div class="controls"><select id="destinationBroker"><option>NGSIv1</option><option>NGSIv2</option><option>NGSI-LD</option></select></div>'        
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Reference URL (*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="referenceURL"></div>';
        html += '</div>';        

        html += '<div class="control-group"><label class="control-label" for="input01">Tenant (*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="tenant"></div>';
        html += '</div>';              

        html += '<div class="control-group"><label class="control-label" for="input01"></label>';
        html += '<div class="controls"><button id="submitRegistration" type="button" class="btn btn-primary">Register</button>';
        html += '</div></div>';

        html += '</fieldset></div>';


        $('#content').html(html);

        // associate functions to clickable buttons
        $('#submitRegistration').click(function () {
            registerNewSubscription();
        });
    }    

    function registerNewSubscription() {
        console.log('register a new permanent subscription');

        // take the inputs    
        var eType = $('#entityType').val();
        console.log(eType);
        var destinationBroker = $('#destinationBroker').val();
        console.log(destinationBroker);
        var referenceURL = $('#referenceURL').val();
        console.log(referenceURL);        
        var tenant = $('#tenant').val();
        console.log(tenant);            

        var subscription = {
            "entity_type": eType,
            "destination_broker": destinationBroker,
            "reference_url": referenceURL,
            "tenant": tenant
        };

        fetch("/subscription", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify(subscription)
        })
        .then(response => {
            console.log("submit a new permanent subscription: ", response.status)
            showSubscriptions();
        })
        .catch(err => console.log(err));               
    }


    function showTasks() {
        $('#info').html('list of running data processing tasks');

        fetch('/info/task').then(res => res.json()).then(tasks => {
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

    function showStreams() {
        $('#info').html('list of all entities');

        var queryReq = {}
        queryReq.entities = [{ id: '.*', isPattern: true }];

        client.queryContext(queryReq).then(function (entityList) {
            displayEntityList(entityList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function displayEntityList(entities) {
        if (entities == null || entities.length == 0) {
            $('#content').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Entity Type</th>';
        html += '<th>Attributes</th>';
        html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < entities.length; i++) {
            var entity = entities[i];

            html += '<tr>';
            html += '<td>' + entity.entityId.id + '</td>';
            html += '<td>' + entity.entityId.type + '</td>';
            html += '<td>' + JSON.stringify(entity.attributes) + '</td>';
            html += '<td>' + JSON.stringify(entity.metadata) + '</td>';
            html += '</tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }

});



