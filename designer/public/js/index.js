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
    addMenuItem('uService', showEndPointService);
    addMenuItem('Task', showTasks);    
    addMenuItem('Stream', showStreams);

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
        $('#info').html('list of all IoT devices');

        var html = '<div style="margin-bottom: 10px;"><button id="addNewDevice" type="button" class="btn btn-primary">add</button></div>';

        html += '<div id="deviceList"></div>';

        html += '<div id="lddeviceList"></div>';

        $('#content').html(html);

        updateDeviceList();

        $("#addNewDevice").click(function () {
            deviceRegistration();
        });
    }

    function removeLDDevice(entityID) {
        ldclient.deleteContext(entityID, true).then(function (res) {
            if (res.status == 204) {
                showDevices();
            }
            console.log("ld re aaa ", res);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to upsert API of ld device');
        });
    }

    function showWorkerList() {
        fetch('/info/worker').then(res => res.json()).then(workers => {
            displayWorkerOnMap(workers)
        });        
        
//        var queryReq = {}
//        queryReq.entities = [{ "type": 'Worker', "isPattern": true }];
//        client.queryContext(queryReq).then(function (edgeNodeList) {
//            // show edge nodes on the map
//            displayWorkerOnMap(edgeNodeList);
//        }).catch(function (error) {
//            console.log(error);
//            console.log('failed to query the list of workers');
//        });
    }


    function updateDeviceList() {
        // query and update the list of all NGSI v1 devices
        var queryReq = {}
        queryReq.entities = [{ id: 'Device.*', isPattern: true }];        
        client.queryContext(queryReq).then(function (deviceList) {
            displayDeviceList(deviceList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
        });
                
        // query and update the list of all NGSI-LD devices
        var queryPayload = {
            "type": "Query",
            "entities": [{
                "idPattern": "urn:ngsi-ld:device:.*"
            }]
        }
        ldclient.queryContext(queryPayload).then(function (result) {
            console.log("get device based on ngsi-ld ", result[0]);
            displayLDDeviceList(result[0]);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query the list of ld device');
        });       
    }

    function displayDeviceList(devices) {
        if (devices == null || devices.length == 0) {
            $('#deviceList').html('');
            return
        }
        var html = '<table class="table table-striped table-bordered table-condensed">';
        html += '<div><b>Devices based on NGSI-v1</b></div><br>';
        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Type</th>';
        html += '<th>Attributes</th>';
        html += '<th>DomainMetadata</th>';
        html += '<th>Action</th>';
        html += '</tr></thead>';

        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];
            html += '<tr>';
            html += '<td>' + device.entityId.id + '<br>';
            html += '</td>';
            html += '<td>' + device.entityId.type + '</td>';
            html += '<td>' + JSON.stringify(device.attributes) + '</td>';
            html += '<td>' + JSON.stringify(device.metadata) + '</td>';
            html += '<td>'
            html += '<button id="DOWNLOAD-' + device.entityId.id + '" type="button" class="btn btn-default">Profile</button>';
            html += '<button id="DELETE-' + device.entityId.id + '" type="button" class="btn btn-primary">Delete</button>';
            html += '</td>'
            html += '</tr>';
        }

        html += '</table>';
        $('#deviceList').html(html);

        //associate a click handler to generate device profile on request
        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];

            var profileButton = document.getElementById('DOWNLOAD-' + device.entityId.id);
            profileButton.onclick = function (d) {
                var myProfile = d;
                return function () {
                    downloadDeviceProfile(myProfile);
                };
            }(device);

            var deleteButton = document.getElementById('DELETE-' + device.entityId.id);
            deleteButton.onclick = function (d) {
                var myProfile = d;
                return function () {
                    removeDeviceProfile(myProfile);
                };
            }(device);
        }
    }

    function downloadDeviceProfile(deviceObj) {
        var profile = {};

        profile.id = deviceObj.attributes.id.value;
        profile.type = deviceObj.entityId.type;
        profile.iconURL = deviceObj.attributes.iconURL.value;
        profile.pullbased = deviceObj.attributes.pullbased.value;
        profile.location = deviceObj.metadata.location.value;
        profile.discoveryURL = config.discoveryURL;

        var content = JSON.stringify(profile);
        var dl = document.createElement('a');
        dl.setAttribute('href', 'data:text/json;charset=utf-8,' + encodeURIComponent(content));
        dl.setAttribute('download', 'profile-' + profile.id + '.json');
        dl.click();
    }


    function removeDeviceProfile(deviceObj) {
        var entityid = {
            id: deviceObj.entityId.id,
            isPattern: false
        };

        client.deleteContext(entityid).then(function (data) {
            console.log('remove the device');

            // show the updated device list
            showDevices();
        }).catch(function (error) {
            console.log('failed to cancel a requirement');
        });
    }


    function displayLDDeviceList(devices) {
        if (devices == null || devices.length == 0) {
            $('#lddeviceList').html('');
            return
        }
        var html = '<table class="table table-striped table-bordered table-condensed">';
        html += '<div><b>Devices based on NGSI-v1</b></div><br>';
        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Type</th>';
//        html += '<th>Attributes</th>';
//        html += '<th>DomainMetadata</th>';
        html += '<th>Action</th>';
        html += '</tr></thead>';

        for (var i = 0; i < devices.length; i++) {
            var lddevice = devices[i];
            
            console.log(lddevice);
            
            html += '<tr>';
            html += '<td>' + lddevice.id + '<br>';
            html += '</td>';
            html += '<td>' + lddevice.type + '</td>';
//            html += '<td>' + JSON.stringify(device.attributes) + '</td>';
//            html += '<td>' + JSON.stringify(device.metadata) + '</td>';
            html += '<td>'
            html += '<button id="DOWNLOAD-' + lddevice.id + '" type="button" class="btn btn-default">Profile</button>';
            html += '<button id="DELETE-' + lddevice.id + '" type="button" class="btn btn-primary">Delete</button>';
            html += '</td>'
            html += '</tr>';
        }

        html += '</table>';
        $('#lddeviceList').html(html);

        //associate a click handler to generate device profile on request
        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];

            var profileButton = document.getElementById('DOWNLOAD-' + device.entityId.id);
            profileButton.onclick = function (d) {
                var myProfile = d;
                return function () {
                    downloadDeviceProfile(myProfile);
                };
            }(device);

            var deleteButton = document.getElementById('DELETE-' + device.entityId.id);
            deleteButton.onclick = function (d) {
                var myProfile = d;
                return function () {
                    removeDeviceProfile(myProfile);
                };
            }(device);
        }
    }


    function deviceRegistration() {
        $('#info').html('to register a new IoT device');
        
        var html = ''

        html += '<div class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device Protocol(*)</label>';
        html += '<div class="controls"><select id="deviceProtocol"><option>NGSI-v1</option><option>NGSI-LD</option></select></div>'
        html += '</div>';  

        html += '<div class="control-group"><label class="control-label" for="input01">Device Type(*)</label>';
        html += '<div class="controls"><select id="deviceType"><option>Temperature</option><option>PowerPanel</option><option>Camera</option><option>Alarm</option></select></div>'
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
            var protocolType = $('#deviceProtocol option:selected').val();            
            console.log(protocolType)
            if (protocolType === 'NGSI-LD') {
                registerLDNewDevice()
            } else {
                registerNewDevice()
            }
        });

        $('#autoIDGenerator').click(autoIDGenerator);

        $('#iconImage').change(function () {
            readIconImage(this);
        });
        $('#imageContent').change(function () {
            readContentImage(this);
        });

        showWorkerList();

        $('#tabV1').change(function () {
            document.getElementById('ngsildDeviceRegistration').style.display = 'none';
            document.getElementById('deviceRegistration').style.display = 'block';
        });

        $('#tabLD').change(function () {
            document.getElementById('ngsildDeviceRegistration').style.display = 'block';
            document.getElementById('deviceRegistration').style.display = 'none';
        });

        $('#ldDeviceType').change(function () {
            $('#ldDeviceID').val('');
            var type = $('#ldDeviceType option:selected').val();
            $('#ldDeviceName').val(type.toLowerCase());
        });
    }

    function registerLDNewDevice() {        
        // take the inputs    
        var id = $('#deviceID').val();
        console.log(id);

        var type = $('#deviceType option:selected').val();
        console.log(type);

        if (id == '' || type == '' || locationOfNewDevice == null) {
            alert('please provide the required inputs');
            return;
        }

        console.log(locationOfNewDevice);                

        //createdAt
        var ldPayload = {}
        ldPayload.id = "urn:ngsi-ld:" + "device" + ":" + id;
        ldPayload.type = type;
        
//        ldPayload[deviceName] = {
//            "type": "Property",
//            "value": deviceValue
//        }

        ldPayload.location = {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [locationOfNewDevice.lat, locationOfNewDevice.lng]
            }
        }
        var cTime = new Date().toISOString()
        var createdAt = cTime.split(".")[0]
        ldPayload.createdAt = createdAt;
        console.log("payload ---- ", ldPayload);

        var ldclient = new NGSILDclient(config.LdbrokerURL);

        ldclient.updateContext(ldPayload).then(function (res) {
            showDevices();
        }).catch(function (error) {
            console.log(error);
            console.log('failed to upsert API of ld device');
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
        console.log(id);

        var type = $('#deviceType option:selected').val();
        console.log(type);

        if (id == '' || type == '' || locationOfNewDevice == null) {
            alert('please provide the required inputs');
            return;
        }

        console.log(locationOfNewDevice);

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

        newDeviceObject.entityId = {
            id: 'Device.' + type + '.' + id,
            type: type,
            isPattern: false
        };

        newDeviceObject.attributes = {};
        newDeviceObject.attributes.DeviceID = { type: 'string', value: id };

        var url = 'http://' + config.agentIP + ':' + config.webSrvPort + '/photo/' + contentImageFileName;
        newDeviceObject.attributes.url = { type: 'string', value: url };
        newDeviceObject.attributes.iconURL = { type: 'string', value: '/photo/' + iconImageFileName };

        if (type == "PowerPanel") {
            newDeviceObject.attributes.usage = {
                type: 'integer',
                value: 20
            };
            newDeviceObject.attributes.shop = {
                type: 'string',
                value: id
            };
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

        client.updateContext(newDeviceObject).then(function (data) {
            console.log(data);
            // show the updated device list
            showDevices();
        }).catch(function (error) {
            console.log('failed to register the new device object');
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
        var map = new L.Map('map', { layers: [osm], center: new L.LatLng(35.692221, 139.709059), zoom: 7 });
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

    
//    function displayLDDeviceList(devices) {
        
//        var html = '<table id="ld-device-table" class="table table-striped table-bordered table-condensed">';
        
//        html += '<br><div><b>Device based on NGSI-LD</b></div><br>';
//        html += '<thead><tr>';
//        html += '<th>ID</th>';
//        html += '<th>Type</th>';
//        html += '<th>Property</th>';
//        html += '<th>Action</th>';
//        html += '</tr></thead>';

//        for (var i = 0; i < devices.length; i++) {
//            var device = devices[i];
//            html += '<tr>';
//            html += '<td>' + device.id + '<br>';
//            html += '</td>';
//            html += '<td>' + device.type + '</td>';
//            html += '<td>' + (device.attribute) + '</td>';
//            html += '<td><button id="DELETE-' + device.id + '" type="button" class="btn btn-primary">Delete</button></td>';
//            html += '</tr>';
//        }
//        html += '</table>';
        
//        return html;
//    }
});



