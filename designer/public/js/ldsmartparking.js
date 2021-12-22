$(function() {

    // initialization  
    var handlers = {}


    // location of new device
    var locationOfNewDevice = null;
    // icon image for device registration
    var iconImage = null;
    var iconImageFileName = null;

    var curMap = null;
    var myCenter = new L.LatLng(37.990905, -1.131133);
    var timerID = null;

    var recommendedParkingSite = null;
    var carMarker = null;

    var privateParkingSites = [{
        id: "001",
        location: {
            latitude: 37.997849,
            longitude: -1.124129
        },
        iconURL: "/img/parkingsite.png",
        datasource: "http://fiware-dev.inf.um.es:1026/v2/"
    }, {
        id: "002",
        location: {
            latitude: 37.991524,
            longitude: -1.115503
        },
        iconURL: "/img/parkingsite.png",
        datasource: "http://fiware-dev.inf.um.es:1026/v2/"
    }, {
        id: "003",
        location: {
            latitude: 37.986957,
            longitude: -1.128292
        },
        iconURL: "/img/parkingsite.png",
        datasource: "http://fiware-dev.inf.um.es:1026/v2/"
    }];

    var publicParkingSites = [{
        id: "006",
        location: {
            latitude: 38.000011,
            longitude: -1.128464
        },
        iconURL: "/img/parkingsite.png",
        datasource: "http://fiware-dev.inf.um.es:1026/v2/"
    }, {
        id: "007",
        location: {
            latitude: 38.009850,
            longitude: -1.142369
        },
        iconURL: "/img/parkingsite.png",
        datasource: "http://fiware-dev.inf.um.es:1026/v2/"
    }];


    addMenuItem('ProcessingFlow', 'Processing Flow', showProcessingFlows);
    addMenuItem('DigitalTwin', 'Digital Twin', showTwins);
    addMenuItem('RunningTask', 'Running Task', showTasks);
    addMenuItem('SmartParking', 'Smart Parking', showParking);

    //connect to the socket.io server via the NGSI proxy module
    var ngsiproxy = new NGSIProxy();
    ngsiproxy.setNotifyHandler(handleNotify);

    // client to interact with IoT Broker
    var ldclient = new NGSILDclient(config.LdbrokerURL);
    var client = new NGSI10Client(config.brokerURL);
    subscribeResult();

    initParkingSite();

    showProcessingFlows();

    $(window).on('hashchange', function() {
        var hash = window.location.hash;

        if (hash != '#SmartParking' && timerID != null) {
            console.log('terminate the current timer ' + timerID + ' when switch the menu items');
            clearInterval(timerID);
        }

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
        handler();
    }

    function subscribeResult() {
        var subscribeCtxReq = {};
        /*subscribeCtxReq.entities = [{
            id: 'Twin.ConnectedCar.01',
            type: 'ConnectedCar',
            isPattern: false
        }];
        subscribeCtxReq.attributes = ['RecommendedParkingSite'];*/
        //subscribeCtxReq.reference = 'http://' + config.agentIP + ':' + config.agentPort;

	subscribeCtxReq.type = "Subscription"

        subscribeCtxReq.entities = [{ id:'urn:ngsi-ld:Twin.ConnectedCar.01' ,type: 'LDconnectedcar' }];

        subscribeCtxReq.notification = {}
        endPoint = {}
        subscribeCtxReq.notification.format = "normalized"

        endPoint.uri = 'http://' + config.agentIP + ':' + config.ldAgentPort + '/notifyContext';
        endPoint.accept = "application/ld+json"
        subscribeCtxReq.notification.endpoint = endPoint
        console.log("subscribeCtxReq data for topology",subscribeCtxReq)

        ldclient.subscribeContext(subscribeCtxReq).then(function(subscriptionId) {
            console.log(subscriptionId);
            ngsiproxy.reportSubID(subscriptionId);
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });
    }

    function handleNotify(contextObj) {
        console.log(contextObj);

        if (contextObj.RecommendedParkingSite == null) {
            return
        }

        recommendedParkingSite = contextObj.RecommendedParkingSite.value;

        var hash = window.location.hash;
        if (hash == '#SmartParking') {
            console.log("recommend parking site " + recommendedParkingSite);
            updateRecommendationResult();
        }
    }

    function updateRecommendationResult() {
        var message = '<b>you can park at <font color="red">' + recommendedParkingSite + '</font></b>';
        carMarker.bindPopup(message, { closeOnClick: false });
        carMarker.openPopup();
    }

    function initParkingSite() {
        // for private parking sites
        for (var i = 0; i < privateParkingSites.length; i++) {
            var privatesite = privateParkingSites[i];
            createParkingSiteEntity(privatesite, "LDPrivateSite");
        }

        // for public parking sites    
        for (var i = 0; i < publicParkingSites.length; i++) {
            var publicsite = publicParkingSites[i];
            createParkingSiteEntity(publicsite, "LDPublicSite");
        }
    }

    function createParkingSiteEntity(site, siteType) {
        var siteEntity = {};

        /*siteEntity.entityId = {
            id: 'Twin.ParkingSite.' + site.id,
            type: siteType,
            isPattern: false
        };*/
	siteEntity.id = 'urn:ngsi-ld:Twin.ParkingSite.' + site.id;
	siteEntity.type = siteType;

        //siteEntity.attributes = {};
        siteEntity.iconURL = { type: 'Property', value: site.iconURL };
        siteEntity.datasource = { type: 'Property', value: site.datasource };

        //siteEntity.metadata = {};
        siteEntity.location = {
            type: 'GeoProperty',
            value: {
	    	type: 'Point',
		coordinates: [site.location.latitude, site.location.longitutde] 
            }
        };

        ldclient.updateContext(siteEntity).then(function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to create a parking site entity');
        });
    }


    function showProcessingFlows() {
        $('#info').html('to show the logical data processing flows behind this service');

        var html = '';
        html += '<div><img width="50%" src="/img/smart-parking.png"></img></div>';

        $('#content').html(html);
    }

    function showTwins() {
        $('#info').html('list of all digital twins and each of them is a virtual entity');

        var html = '<div id="twinList"></div>';
        $('#content').html(html);
        updateTwinList();
    }

    function updateTwinList() {
        var queryReq = {}
        //queryReq.entities = [{ id: 'Twin.*', isPattern: true }];
	queryReq.type = 'Query';
	queryReq.entities = [{ idPattern : 'urn:ngsi-ld:Twin.*' }];
        ldclient.queryContext(queryReq).then(function(twinList) {
            console.log(twinList);
            displayTwinList(twinList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function displayTwinList(twins) {
        if (twins == null || twins.length == 0) {
            $('#twinList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Type</th>';
        //html += '<th>Attributes</th>';
        //html += '<th>DomainMetadata</th>';
        html += '</tr></thead>';

        for (var i = 0; i < twins.length; i++) {
            var twin = twins[i];
            console.log(twin);
        for (var j = 0; j < twin.length; j++) {
	    var displayTwin = twin[j];
            html += '<tr>';
            html += '<td>' + displayTwin.id + '<br>';
            html += '<button id="remove-' + displayTwin.id + '" type="button" class="btn btn-default">remove</button>';
            html += '</td>';
            html += '<td>' + displayTwin.type + '</td>';
            //html += '<td>' + JSON.stringify(twin.attributes) + '</td>';
            //html += '<td>' + JSON.stringify(twin.metadata) + '</td>';
            html += '</tr>';
        }
      }

        html += '</table>';

        $('#twinList').html(html);

        // associate a click handler to generate twin profile on request
        for (var i = 0; i < twins.length; i++) {
            var twin = twins[i];
            console.log(twin);
        for (var j = 0; j < twin.length; j++) {
	     var display = twin[j];
            var removeButton = document.getElementById('remove-' + display.id);
            removeButton.onclick = function(d) {
                var myProfile = d;
                return function() {
                    removeDigitalTwin(myProfile);
                };
            }(twin);
        }
      }
    }


    function removeDigitalTwin(deviceObj) {
        var entityid = {
            id: deviceObj.entityId.id,
            isPattern: false
        };

        client.deleteContext(entityid).then(function(data) {
            console.log('remove the digital twin');

            // show the updated digital twin list
            showTwins();
        }).catch(function(error) {
            console.log('failed to cancel a requirement');
        });
    }


    function showTasks() {
        $('#info').html('list of all triggerred function tasks');

        var queryReq = {}
        queryReq.entities = [{ type: 'Task', isPattern: true }];

        client.queryContext(queryReq).then(function(taskList) {
            console.log(taskList);
            displayTaskList(taskList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query task');
        });
    }


    function displayTaskList(tasks) {
        $('#info').html('list of all function tasks that have been triggerred');

        if (tasks.length == 0) {
            $('#content').html('');
            return;
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>ID</th>';
        html += '<th>Type</th>';
        html += '<th>Service</th>';
        html += '<th>Task</th>';
        html += '<th>Worker</th>';
        html += '<th>port</th>';
        html += '<th>status</th>';
        html += '</tr></thead>';

        for (var i = 0; i < tasks.length; i++) {
            var task = tasks[i];
            html += '<tr>';
            html += '<td>' + task.entityId.id + '</td>';
            html += '<td>' + task.entityId.type + '</td>';
            html += '<td>' + task.attributes.service.value + '</td>';
            html += '<td>' + task.attributes.task.value + '</td>';
            html += '<td>' + task.metadata.worker.value + '</td>';

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


    function showParking() {
        $('#info').html('to illustrate the smart parking use case for Murcia');

        var html = '';

        html += '<div id="map"  style="width: 800px; height: 600px"></div>';

        $('#content').html(html);

        var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
        var osm = L.tileLayer(osmUrl, { maxZoom: 15, zoom: 13 });
        var map = new L.Map('map', {
            layers: [osm],
            center: myCenter,
            zoom: 13,
            zoomControl: true
        });

        var drawnItems = new L.FeatureGroup();
        map.addLayer(drawnItems);

        // show edge nodes on the map
        displayEdgeNodeOnMap(map);

        // show moving car
        drawConnectedCar(map);

        // display parking sites
        displayParkingSites(map);

        // remember the created map
        curMap = map;
    }

    function displayEdgeNodeOnMap(map) {
        var queryReq = {}
        queryReq.entities = [{ type: 'Worker', isPattern: true }];
        queryReq.restriction = { scopes: [{ scopeType: 'stringQuery', scopeValue: 'role=EdgeNode' }] }
        client.queryContext(queryReq).then(function(edgeNodeList) {
            console.log(edgeNodeList);

            var edgeIcon = L.icon({
                iconUrl: '/img/gateway.png',
                iconSize: [48, 48]
            });

            for (var i = 0; i < edgeNodeList.length; i++) {
                var worker = edgeNodeList[i];

                console.log(worker);

                latitude = worker.attributes.location.value.latitude;
                longitude = worker.attributes.location.value.longitude;
                edgeNodeId = worker.entityId.id;

                console.log(latitude, longitude, edgeNodeId);

                var marker = L.marker(new L.LatLng(latitude, longitude), { icon: edgeIcon });
                marker.nodeID = edgeNodeId;
                marker.addTo(map).bindPopup(edgeNodeId);
                marker.on('click', showRunningTasks);

                console.log('=======draw edge on the map=========');
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

    function drawConnectedCar(map) {
        var taxiIcon = L.icon({
            iconUrl: '/img/taxi.png',
            iconSize: [80, 80]
        });

        var path = [
            [37.996655, -1.150094],
            [37.984174, -1.141039]
        ];
        carMarker = L.Marker.movingMarker(path, [10000], { autostart: false, loop: true });
        carMarker.options.icon = taxiIcon;

        map.addLayer(carMarker);

        carMarker.on('click', function() {
            if (carMarker.isRunning()) {
                console.log('timerID = ', timerID);
                if (timerID != null) {
                    clearInterval(timerID);
                }
                carMarker.pause();
                carMarker.bindPopup('<b>Click me to start !</b>').openPopup();
            } else {
                carMarker.start();
                carMarker.bindPopup('<b>Click me to pause !</b>').openPopup();
                timerID = setInterval(function() {
                    var mylocation = carMarker.getLatLng();
                    updateMobileObject(mylocation);
                }, 1000);
                console.log('timerID = ', timerID);
            }
        });

        carMarker.bindPopup('<b>Click me to start !</b>', { closeOnClick: false });
        carMarker.openPopup();
    }

    function displayParkingSites(map) {
        var queryReq = {}
	queryReq.type = 'Query';
        queryReq.entities = [{ idPattern: 'urn:ngsi-ld:Twin.ParkingSite.*'}];
        ldclient.queryContext(queryReq).then(function(sites) {
            console.log(sites);

            for (var i = 0; i < sites.length; i++) {
                var site = sites[i];
	    for (var j = 0; j < site.length; j++) {
		 var displaySite = site[j];
                 console.log(" display sites ",displaySite)
                if (displaySite.iconURL != null) {
                    var iconImag = displaySite.iconURL.value;
                    var icon = L.icon({
                        iconUrl: iconImag,
                        iconSize: [48, 48]
                    });

                    latitude = displaySite.location.value.coordinates[0];
                    longitude = displaySite.location.value.coordinates[1];
                    siteId = displaySite.id;

                    var marker = L.marker(new L.LatLng(latitude, longitude), { icon: icon });
                    marker.addTo(map).bindPopup(siteId);
                }
            }
       }

        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function updateMobileObject(location) {
        //register a new device
        var movingCarObject = {};

        movingCarObject.id = 'urn:ngsi-ld:Twin.ConnectedCar.01';
        
        movingCarObject = {};
        movingCarObject.iconURL = { type: 'Property', value: '/img/taxi.png' };
        movingCarObject.location = {
            type: 'GeoProperty',
            value: { 
			type: 'Point',
			coordinates: [location.lat, location.lng ]}
        };


        ldclient.updateContext(movingCarObject).then(function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to update car object');
        });

    }

});
