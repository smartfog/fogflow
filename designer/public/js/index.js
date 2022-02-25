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
    addMenuItem('Edge', showEdge);
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


    // visualization : start
    var map = undefined;

    var curMap = undefined;
        
    var LDTestData = [
        {
            "@context": "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
            "brandName": {
                "type": "Property",
                "value": "Mercedes"
            },
            "createdAt": "2022-01-06 10:35:18.506423711 +0530 IST m=+41.091378705",
            "id": "urn:ngsi-ld:Vehicle:A100",
            "isParked": {
                "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                "observedAt": "2017-07-29T12:00:04",
                "providedBy": {
                    "type": "Relationship",
                    "object": "urn:ngsi-ld:Person:Bob"
                },
                "type": "Relationship"
            },
            "location": {
                "type": "GeoProperty",
                "value": {
                    "type": "Point",
                    "coordinates": [
                        -8.5,
                        41.2
                    ],
                    "geometries": null
                }
            },
            "modifiedAt": "2022-01-06 10:35:18.506415761 +0530 IST m=+41.091370789",
            "speed": {
                "type": "Property",
                "value": 80
            },
            "type": "Vehicle"
        }
    ]

    var edgeNodesList = undefined;
    var brokers = [];

    function getBrokerList() {
        brokers =[]
        var discoverReq = {}
        
        discoverReq.entities = [{ type: 'IoTBroker', isPattern: true }];
        // v1 API for get broker list
        var ngsi9client = new NGSI9Client(config.discoveryURL)
        console.log('discorvery url is ',config.discoveryURL);
        displayEdgeNode();
        // v1 API for get broker list      
        ngsi9client.discoverContextAvailability(discoverReq).then(function (response) {
            console.log("broker response: ",response);
            
            if (response.errorCode.code == 200 && response.hasOwnProperty('contextRegistrationResponses')) {
                for (var i in response.contextRegistrationResponses) {
                    var tmpIn = parseInt(i)+1;
                    var contextRegistrationResponse = response.contextRegistrationResponses[i];
                    console.log("broker details is *** ",contextRegistrationResponse);
                    var brokerID = contextRegistrationResponse.contextRegistration.entities[0].id;
                    var providerURL = contextRegistrationResponse.contextRegistration.providingApplication;
                    if (providerURL != '') {
                        brokers.push({ id: brokerID, brokerURL: providerURL });
                        console.log("inside broker ",brokers);
                        if (i==0) {
                            var option = document.createElement("option");
                            option.text ='All'; 
                            var allEdgeBrokerList = document.getElementById("allEdgeBrokerList");
                            allEdgeBrokerList.add(option);
                        }
                        var option = document.createElement("option");
                        var edgeName = brokerID.replace('Broker','Edge');
                        if (edgeName.includes(".001")) edgeName = edgeName.replace('Edge','Cloud');
                        option.text = edgeName
                        console.log("get broker ID: " + brokerID);
                        var allEdgeBrokerList = document.getElementById("allEdgeBrokerList");
                        console.log("all edge obje ",allEdgeBrokerList);
                        allEdgeBrokerList.add(option);

                        if (tmpIn === response.contextRegistrationResponses.length){
                            getWorkerList(displayEdgeOnMap,brokers)
                        }
                    }
                }
                
            }

        }).catch(function (error) {
            console.log(error);
            console.log('failed to query context');
        });
       // return brokers;
    }
    var ldDeviceObj = []

    function resetEdges(){
        try{
        curMap = showMap();
        }catch(err){console.log(err);}
        ldDeviceObj=[]
        $('#deviceList').html('');
    }
    
    // for LD payload
    function getEntityProperties(ldDataList){
        for(var i=0;i<ldDataList.length;i++){
            var objKeyList = Object.keys(ldDataList[i])
            var objValue ={} 
            var finalLDEntityDetails = {}
            objKeyList.forEach(ldKey => {
                if (ldDataList[i][ldKey].hasOwnProperty("type") && ldDataList[i][ldKey]["type"] == "Property"){
                    objValue[ldKey] = ldDataList[i][ldKey]["value"]
                    console.log(objValue)
                }
            });
            if(ldDataList[i].hasOwnProperty("location") && ldDataList[i].location.hasOwnProperty("value")){
                if(ldDataList[i].location.value.type == "Point")
                {
                    finalLDEntityDetails['isLocation']= true;
                }else finalLDEntityDetails['isLocation']= false;
            }
            finalLDEntityDetails['id']=ldDataList[i].id;
            finalLDEntityDetails['type']=ldDataList[i].type;
            finalLDEntityDetails['attribute']=JSON.stringify(objValue);
            ldDeviceObj.push(finalLDEntityDetails)

        }
    }

    // API call for LD devices
    function ldEntities(ldBrokerURL){
        var sourcePath=ldBrokerURL.replace("ngsi10","ngsi-ld");
        var ldclient = new NGSILDclient(sourcePath);
        ldclient.GetEntitiesContext('All').then(function (edgeNodeList) {
        getEntityProperties(edgeNodeList)
        for (var i = 0; i < edgeNodeList.length; i++) {
            var edgeEntity = edgeNodeList[i];
            if(i+1 == edgeNodeList.length)
            {
                $('#deviceList').html(displayLDDeviceList4Edge(ldDeviceObj));
            }
            try{
                var latitude = edgeEntity.location.value.coordinates[0];
                var longitude = edgeEntity.location.value.coordinates[1];
                var edgeNodeId = edgeEntity.id;
                var edgeIcon = edgeDeviceIcon(edgeEntity.type,true);
                var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
                marker.nodeID = edgeNodeId;
                marker.addTo(curMap).bindPopup(edgeNodeId);
            }catch (e) {
                console.log(e);
            }
        }
        }).catch(function (error) {
            console.log(error);
            console.log('failed to query the list of ld device');
        });
        
    }

    var workerWithDeviceList = {};
    function getWorkerList(callback,brokerObj) {
        var edgeBrokerEntityList = [];
        workerWithDeviceList = {};
        for (var i = 0; i < brokerObj.length; i++){
            var tmpI = i;
            var tmpClient = new NGSI10Client(brokerObj[i].brokerURL);
            var queryReq = {}
            /**
             * get only worker list from the broker
             */
            queryReq.entities = [{ "type": 'Worker', "isPattern": true }];
            console.log("get broker id is --- ",brokerObj[i].brokerURL);
            ldEntities(brokerObj[i].brokerURL)
            // call v1 API for get worker list and v1 devices
            tmpClient.queryContext(queryReq).then(function (edgeNodeList) {
                console.log("inside in worker list ",edgeNodeList);
                if (edgeNodeList){
                    edgeBrokerEntityList = edgeBrokerEntityList.concat(edgeNodeList);
                    let workerObj = edgeNodeList.find(o => o.entityId.type === 'Worker');
                    if (workerObj){
                    var workerId = workerObj.entityId.id;
                    workerWithDeviceList[workerId] = edgeNodeList;
                    }
                }
                if (tmpI+1==brokerObj.length){
                    callback(edgeBrokerEntityList)
                }
            }).catch(function (error) {
                console.log(error);
                console.log('failed to query the list of workers');
            });
        }
    }

    function showEdge(){
        resetEdges();
        $('#info').html('list of all IoT Edge');
        var html = '<div id="dockerRegistration" class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Edge(*)</label>';
        html += '<div class="controls"><select id="allEdgeBrokerList" ></select>';
        html += '</div></div>';
        html += '</fieldset></div>';

        html += '<div id="deviceList"></div>';

        html += '<div id="map"  style="width: 900px; height: 500px"></div>';
        $('#content').html(html);
        curMap = showMap();
        getBrokerList();
    }

    function displayEdgeNode() {
        $('#allEdgeBrokerList').change(function() {
            var selectedEdgeBroker = $('#allEdgeBrokerList option:selected').val();
            console.log('selected edge broker is ',selectedEdgeBroker)
            if (selectedEdgeBroker === 'All'){
                resetEdges();
                getWorkerList(displayEdgeOnMap,brokers)
            }
            else
            {
                selectedEdgeBroker = selectedEdgeBroker.replace(/Edge|Cloud/gi,'Broker')
                let obj = brokers.find(o => o.id === selectedEdgeBroker);
                console.log('selected edge broker is obj ',obj)
                resetEdges()
                getWorkerList(displayEdgeOnMap,[obj])
            }
          });
    }
    
    function edgeDeviceIcon(type_,isLDDevice){
        var edgeIcon = undefined;
        if (type_ === 'Worker'){
             edgeIcon = L.icon({
                iconUrl: '/img/edge7.png',
                iconSize: [48, 48]
            });
        }else {
             var iconURL = '/img/device.png';
             if (isLDDevice) iconURL = '/img/lddevice.png'
             edgeIcon = L.icon({
                iconUrl: iconURL,
                iconSize: [40, 40]
            });
        }
        return edgeIcon;
    }
    
    // selected broker validate with worker 
    function selectedWorkerValidate(workerID) {
     var selectedEdge =  $('#allEdgeBrokerList option:selected').val();
     if (selectedEdge === 'All') return true;
     if (workerID.includes(".") && selectedEdge !== undefined) {
        let wId = workerID.substr(workerID.indexOf('.'));
        let selectedId = selectedEdge.substr(selectedEdge.indexOf('.'));
        if (wId === selectedId) {
            return true;
        }
     }
     return false;
    }
    
    //for  v1 devices and worker edge 
    function displayEdgeOnMap(workerList) {
        $('#ld-device-table tr:last').after(displayDeviceList4Edge(workerList,false));
        for (var i = 0; i < workerList.length; i++) {
            var edgeEntity = workerList[i];
            try{
                var latitude = edgeEntity.metadata.location.value.latitude;
                var longitude = edgeEntity.metadata.location.value.longitude;
                var edgeNodeId = edgeEntity.entityId.id;
                var edgeIcon = edgeDeviceIcon(edgeEntity.entityId.type,false);
                var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
                marker.nodeID = edgeNodeId;
                  
                if (edgeEntity.entityId.type === 'Worker' && selectedWorkerValidate(edgeNodeId)){
                    console.log("selected edge -- ",$('#allEdgeBrokerList option:selected').val());
                    var id = edgeEntity.entityId.id;
                    console.log("get all worker data ",workerWithDeviceList[id]);
                    var container = $('<div />');
                    container.html(displayDeviceList4Edge(workerWithDeviceList[id],true));
                    marker.addTo(curMap).bindPopup(id);
                }
                else {
                    //marker.addTo(curMap).bindPopup(edgeNodeId);
                }
            }catch (e) {
                console.log(e);
            }
        }
    }

    // for ld devices list
    function displayLDDeviceList4Edge(devices,isFromDevice=false){
        if (devices == null || devices.length == 0) {
            $('#deviceList').html('');
            return
        }
        var html = '<table id="ld-device-table" class="table table-striped table-bordered table-condensed">';
        if (isFromDevice)   html += '<br><div><b>NGSI-LD Table</b></div><br>';
        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Type</th>';
        html += '<th>Attributes</th>';
        if (isFromDevice)   html += '<th>Action</th>';
        html += '</tr></thead>';

        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];
            if (!(device.isLocation)) html += '<tr style="color:red">';
            else html += '<tr>';
            
            html += '<td>' + device.id + '<br>';
            html += '</td>';
            html += '<td>' + device.type + '</td>';
            html += '<td>' + (device.attribute) + '</td>';
            if (isFromDevice) html += '<td><button id="DELETE-' + device.id + '" type="button" class="btn btn-primary">Delete</button></td>';

            html += '</tr>';
        }
        html += '</table>';
        return html;
    }

    // for v1 device list table
    function displayDeviceList4Edge(devices,isFromMap) {
        if (devices == null || devices.length == 0) {
            $('#deviceList').html('');
            return
        }
        
        var html ='';
        for (var i = 0; i < devices.length; i++) {
            var device = devices[i];
            if(device.entityId.type === 'Worker') continue;

            if (!(device.metadata.hasOwnProperty('location'))) html = '<tr style="color:red">';
            else html = '<tr>';
            
            html += '<td>' + device.entityId.id + '<br>';
            html += '</td>';
            html += '<td>' + device.entityId.type + '</td>';
            if (!isFromMap){
            html += '<td>' + JSON.stringify(device.attributes) + '</td>';
            }
            html += '</tr>';
        }
        html += '</table>';
        return html;
    }

    // visualization:end

    function showDevices() {
        $('#info').html('list of all IoT devices');

        var html = '<div style="margin-bottom: 10px;"><button id="addNewDevice" type="button" class="btn btn-primary">add</button></div>';
       

        html += '<div id="deviceList"></div>';

        
        html += '<div id="lddeviceList"></div>';

        $('#content').html(html);
        $('#lddeviceList').html('');
        ldDeviceObj = []
        updateDeviceList();
        ldEntities4DeviceTab();

        $("#addNewDevice").click(function () {
            deviceRegistration();
        });
    }

    function ldEntities4DeviceTab(){
        
        var ldclient = new NGSILDclient(config.LdbrokerURL);
        ldclient.GetEntitiesContext('All').then(function (edgeNodeList) {

        getEntityProperties(edgeNodeList)
        $('#lddeviceList').html(displayLDDeviceList4Edge(ldDeviceObj,true));
        // associate a click handler to generate device profile on request
        for (var i = 0; i < edgeNodeList.length; i++) {
            var device = edgeNodeList[i];
            console.log(device.id);
            var deleteButton = document.getElementById('DELETE-' + device.id);
            deleteButton.onclick = function (d) {
                var myProfile = d;
                return function () {
                    console.log("delete obj ",myProfile)
                    if (confirm('Do you want to delete the entity '+d.id)) {
                        removeLDDevice(d.id);
                      }
                };
            }(device);
        }

        }).catch(function (error) {
            console.log(error);
            console.log('failed to query the list of ld device');
        });
        
    }

    function removeLDDevice(entityID){
        var ldclient = new NGSILDclient(config.LdbrokerURL);
        
        ldclient.deleteContext(entityID,true).then(function (res) {
            if (res.status == 204){
                showDevices();
            }
            console.log("ld re aaa ",res);
            }).catch(function (error) {
                console.log(error);
                console.log('failed to upsert API of ld device');
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
        html += '<div><b>NGSI-V1 Table</b></div><br>';
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
  
    function deviceRegistration() {
        $('#info').html('to register a new IoT device');
        var html = '<p>'
        html += '<input type="radio" checked="checked" name="r1" id="e1" value="NGSIV1"/> NGSI-V1&emsp;&emsp;'   
        html += '<input type="radio"  name="r1" id="e2" value="NGSILD"/> NGSI-LD'
        html += '</p>'
        
        html += '<div class="form-horizontal"><fieldset>';
        html += '<div id="ngsildDeviceRegistration" style="display: none">';

        html += '<div class="control-group"><label class="control-label" for="input01">Device Type(*)</label>';
        html += '<div class="controls"><select id="ldDeviceType"><option>Temperature</option><option>PowerPanel</option></select></div>'
        html += '</div>';

  
        html += '<div class="control-group"><label class="control-label" for="input01">LD Device ID(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="ldDeviceID">';
        html += '<span>  </span><button id="ldAutoIDGenerator" type="button" class="btn btn-primary">Autogen</button>';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Device Name</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="ldDeviceName">';
        html += '</div></div>';


        html += '<div class="control-group"><label class="control-label" for="input01">Device Value</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="ldDeviceValue" placeholder="only number value">';
        html += '</div></div>';
    
        html += '</div>';
       
        html += '<div id="deviceRegistration">';

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

        html += '</div>'; // for v1 tag

        html += '<div class="control-group"><label class="control-label" for="input01">Location(*)</label>';
        html += '<div class="controls"><div id="map"  style="width: 500px; height: 400px"></div></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01"></label>';
        html += '<div class="controls"><button id="submitRegistration" type="button" class="btn btn-primary">Register</button>';
        html += '</div></div>';

        html += '</fieldset></div>';

       
        $('#content').html(html);

        var type = $('#ldDeviceType option:selected').val();
        $('#ldDeviceName').val(type.toLowerCase());
        // associate functions to clickable buttons
        $('#submitRegistration').click(function () {
            console.log("aaaaaaaa")
            var radioVal = $('input[name="r1"]:checked').val();
            console.log("radio val 000000 ",radioVal);
            if( radioVal === 'NGSILD'){
                registerLDNewDevice()
            }else  registerNewDevice()

            });

        $('#autoIDGenerator').click(autoIDGenerator);

        $('#ldAutoIDGenerator').click(ldAutoIDGenerator);


        $('#iconImage').change(function () {
            readIconImage(this);
        });
        $('#imageContent').change(function () {
            readContentImage(this);
        });

       showWorkerList();
       // confirm("Press a button!");
        $('#e1').change(function () {
            var a = $('input[name="r1"]:checked').val();
            console.log('radio button val ',a)
            document.getElementById('ngsildDeviceRegistration').style.display = 'none';
            document.getElementById('deviceRegistration').style.display = 'block';
            showWorkerList();
        });
        $('#e2').change(function () {
            var a = $('input[name="r1"]:checked').val();
            console.log('radio button val ',a)
            document.getElementById('ngsildDeviceRegistration').style.display = 'block';
            document.getElementById('deviceRegistration').style.display = 'none';
            showWorkerList();
        });

        $('#ldDeviceType').change(function () {
            console.log("value of ldDeviceType ",$('#ldDeviceType option:selected').val());
            $('#ldDeviceID').val('');
            var type = $('#ldDeviceType option:selected').val();
            $('#ldDeviceName').val(type.toLowerCase());
        });

    }

    function ldAutoIDGenerator(){
        

        var type = $('#ldDeviceType option:selected').val();
        console.log(type);

        if(type == '') return;
        var date = new Date();
        var epochForId = Math.floor(date.getTime()/1000)
        var uid = "urn:ngsi-ld:"+type.toLowerCase()+":"+epochForId.toString();
        //var uid = "urn:ngsi-ld:"+"device"+":"+epochForId.toString();
        console.log("ngsi ld auto ID",uid);
        $('#ldDeviceID').val(uid);
        
    }


    function registerLDNewDevice(){
        // take the inputs    
        var id = $('#ldDeviceID').val();
        console.log(id);

        var type = $('#ldDeviceType option:selected').val();
        console.log(type);
        

        var deviceName = $('#ldDeviceName').val();
        console.log(deviceName);
        
        var deviceValue = $('#ldDeviceValue').val();
        console.log(deviceValue);

        console.log(locationOfNewDevice);
        

        if (id == '' || type == '' || locationOfNewDevice == '' || locationOfNewDevice == null) {
            alert('please provide the required inputs');
            return;
        }
        if (type === 'Temperature') type = 'Device'
        //createdAt
        var ldPayload = {}
        ldPayload.id = id;
        ldPayload.type = type;
        ldPayload[deviceName] = {
            "type": "Property",
            "value": deviceValue
        }

        ldPayload.location = {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [locationOfNewDevice.lat,locationOfNewDevice.lng]
            }
        }
        var cTime = new Date().toISOString()
        var createdAt = cTime.split(".")[0]
        ldPayload.createdAt = createdAt;
        console.log("payload ---- ",ldPayload);

        var ldclient = new NGSILDclient(config.LdbrokerURL);
        
        ldclient.ldupdateContext(ldPayload).then(function (res) {
            console.log("ld re aaa ",res);
            if (confirm('Do you want to add more attributes?')) {
                $('#ldDeviceName').val('');
                $('#ldDeviceValue').val('');
                // Save it!
                console.log('Thing was saved to the database.');
              } else {
                showDevices();
                // Do nothing!
                console.log('Thing was not saved to the database.');
              }
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

        var tmpAttr = {}
        newDeviceObject.attributes = [];
        tmpAttr = { name:'DeviceID', type: 'string', value: id };
        newDeviceObject.attributes.push(tmpAttr);
        var url = 'http://' + config.agentIP + ':' + config.webSrvPort + '/photo/' + contentImageFileName;
        tmpAttr = { name:'url', type: 'string', value: url };
        newDeviceObject.attributes.push(tmpAttr);
        tmpAttr = { name:'iconURL', type: 'string', value: '/photo/' + iconImageFileName };
        newDeviceObject.attributes.push(tmpAttr);

        if (type == "PowerPanel") {
            tmpAttr = {
                name:'usage',
                type: 'integer',
                value: 20
            };
            newDeviceObject.attributes.push(tmpAttr);
            tmpAttr = {
                name:'shop',
                type: 'string',
                value: id
            };
            newDeviceObject.attributes.push(tmpAttr);
        }

        newDeviceObject.domainMetadata = [];
        var metadata = {};
        metadata = {
            name:'location',
            type: 'point',
            value: { 'latitude': locationOfNewDevice.lat, 'longitude': locationOfNewDevice.lng }
        };
        newDeviceObject.domainMetadata.push(metadata);
        if (type == "PowerPanel") {
            metadata = {
                type: 'string',
                value: id
            };
            newDeviceObject.domainMetadata.push(metadata);
        } else if (type == "Camera") {
            metadata = {
                type: 'string',
                value: id
            };
            newDeviceObject.domainMetadata.push(metadata);
        }
        var finalDeviceObject = {"contextElements":[newDeviceObject],"updateAction": "UPDATE"}
        client.updateContext(finalDeviceObject).then(function (data) {

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
            try{
                map.remove();
            }catch(err){
             
            }
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



