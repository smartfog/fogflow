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

    addMenuItem('Architecture', showArch);
    addMenuItem('Discovery', showDiscovery);
    addMenuItem('Broker', showBrokers);
    addMenuItem('Master', showMaster);
    addMenuItem('Worker', showWorkers);
    addMenuItem('Device', showDevices);
    addMenuItem('uService', showEndPointService);
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

        var discoverReq = {}
        discoverReq.entities = [{ type: 'IoTBroker', isPattern: true }];

        var ngsi9client = new NGSI9Client(config.discoveryURL)
        console.log('discorvery url is ',config.discoveryURL);
        ngsi9client.discoverContextAvailability(discoverReq).then(function (response) {
            var brokers = [];
            if (response.errorCode.code == 200 && response.hasOwnProperty('contextRegistrationResponses')) {
                for (var i in response.contextRegistrationResponses) {
                    var contextRegistrationResponse = response.contextRegistrationResponses[i];
                    var brokerID = contextRegistrationResponse.contextRegistration.entities[0].id;
                    var providerURL = contextRegistrationResponse.contextRegistration.providingApplication;
                    if (providerURL != '') {
                        brokers.push({ id: brokerID, brokerURL: providerURL });
                    }
                }
            }
            displayBrokerList(brokers);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
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
            html += '<tr><td>' + broker.id + '</td><td>' + broker.brokerURL + '</td></tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }

    function showMaster() {
        $('#info').html('list of all topology masters');

        var queryReq = {}
        queryReq.entities = [{ type: 'Master', isPattern: true }];
        client.queryContext(queryReq).then(function (masterList) {
            displayMasterList(masterList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
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
        html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < masters.length; i++) {
            var master = masters[i];
            html += '<tr><td>' + master.entityId.id + '</td><td>' + JSON.stringify(master.metadata) + '</td></tr>';
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

        var queryReq = {}
        queryReq.entities = [{ "type": 'Worker', "isPattern": true }];
        client.queryContext(queryReq).then(function (edgeNodeList) {
            // list all edge nodes in the table
            displayWorkerList(edgeNodeList);

            // show edge nodes on the map
            displayWorkerOnMap(edgeNodeList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query the list of workers');
        });
    }


    function displayWorkerList(workerList) {
        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];

        client.queryContext(queryReq).then(function (tasks) {
            var counter = {};
            for (var i = 0; i < tasks.length; i++) {
                var task = tasks[i];
                var workerID = task.metadata.worker.value;

                if (workerID in counter) {
                    counter[workerID]++;
                } else {
                    counter[workerID] = 1;
                }
            }

            var html = '<thead><tr>';
            html += '<th>ID</th>';
            html += '<th>location</th>';
            html += '<th>capacity</th>';
            html += '<th># of tasks</th>';
            html += '</tr></thead>';

            for (var i = 0; i < workerList.length; i++) {
                var worker = workerList[i];

                var wid = worker.entityId.id;
                var num_task_instances = 0;
                if (wid in counter) {
                    num_task_instances = counter[wid];
                }

                html += '<tr>';
                html += '<td>' + wid + '</td>';
                html += '<td>' + JSON.stringify(worker.attributes.location.value) + '</td>';
                html += '<td>' + JSON.stringify(worker.attributes.capacity.value) + '</td>';
                html += '<td>' + num_task_instances + '</td>';
                html += '</tr>';
            }

            $('#workerList').html(html);

        }).catch(function (error) {
            console.log(error);
            console.log('failed to query task');
        });


    }

    function displayWorkerOnMap(workerList) {
        var curMap = showMap();

        var edgeIcon = L.icon({
            iconUrl: '/img/gateway.png',
            iconSize: [48, 48]
        });

        for (var i = 0; i < workerList.length; i++) {
            var worker = workerList[i];

            var latitude = worker.attributes.location.value.latitude;
            var longitude = worker.attributes.location.value.longitude;
            var edgeNodeId = worker.entityId.id;

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



    function displayWorker(workers) {
        if (workers == null || workers.length == 0) {
            $('#content').html('there is no worker running');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Attributes</th>';
        html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < workers.length; i++) {
            var worker = workers[i];

            html += '<tr>';
            html += '<td>' + worker.entityId.id + '</td>';
            html += '<td>' + JSON.stringify(worker.attributes) + '</td>';
            html += '<td>' + JSON.stringify(worker.metadata) + '</td>';
            html += '</tr>';
        }

        html += '</table>';

        $('#content').html(html);
    }


    function showDevices() {
        $('#info').html('list of all IoT devices');

        var html = '<div style="margin-bottom: 10px;"><button id="addNewDevice" type="button" class="btn btn-primary">add</button></div>';
        html += '<div id="deviceList"></div>';

        $('#content').html(html);

        updateDeviceList();

        $("#addNewDevice").click(function () {
            deviceRegistration();
        });
    }

    function showWorkerList() {
        var queryReq = {}
        queryReq.entities = [{ "type": 'Worker', "isPattern": true }];
        client.queryContext(queryReq).then(function (edgeNodeList) {
            // show edge nodes on the map
            displayWorkerOnMap(edgeNodeList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query the list of workers');
        });
    }


    function updateDeviceList() {
        var queryReq = {}
        queryReq.entities = [{ id: 'Device.*', isPattern: true }];

        client.queryContext(queryReq).then(function (deviceList) {
            displayDeviceList(deviceList);
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
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
        html += '<th>Type</th>';
        html += '<th>Attributes</th>';
        html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];

            html += '<tr>';
            html += '<td>' + device.entityId.id + '<br>';
            html += '<button id="DOWNLOAD-' + device.entityId.id + '" type="button" class="btn btn-default">Profile</button>';
            html += '<button id="DELETE-' + device.entityId.id + '" type="button" class="btn btn-primary">Delete</button>';
            html += '</td>';
            html += '<td>' + device.entityId.type + '</td>';
            html += '<td>' + JSON.stringify(device.attributes) + '</td>';
            html += '<td>' + JSON.stringify(device.metadata) + '</td>';
            html += '</tr>';
        }

        html += '</table>';

        $('#deviceList').html(html);

        // associate a click handler to generate device profile on request
        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];
            console.log(device.entityId.id);

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


    function deviceRegistration() {
        $('#info').html('to register a new IoT device');

        var html = '<div id="deviceRegistration" class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device ID(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="deviceID">';
        html += '<span>  </span><button id="autoIDGenerator" type="button" class="btn btn-primary">Autogen</button>';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device Type(*)</label>';
        html += '<div class="controls"><select id="deviceType"><option>Temperature</option><option>PowerPanel</option><option>Camera</option><option>Alarm</option></select></div>'
        html += '</div>';

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

        showWorkerList();

        // associate functions to clickable buttons
        $('#submitRegistration').click(registerNewDevice);
        $('#autoIDGenerator').click(autoIDGenerator);

        $('#iconImage').change(function () {
            readIconImage(this);
        });
        $('#imageContent').change(function () {
            readContentImage(this);
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
            osm = L.tileLayer(osmUrl, { maxZoom: 18, zoom: 7 }),
            map = new L.Map('map', { layers: [osm], center: new L.LatLng(35.692221, 139.709059), zoom: 7 });


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



