function defaultOperatorList(){
        var operatorList = [{
            name: "nodejs",
            description: "",
            parameters: []
        }, {
            name: "python",
            description: "",
            parameters: []
        }, {
            name: "iotagent",
            description: "",
            parameters: []
        }, {
            name: "counter",
            description: "",
            parameters: []
        }, {
            name: "anomaly",
            description: "",
            parameters: []
        }, {
            name: "facefinder",
            description: "",
            parameters: []
        }, {
            name: "connectedcar",
            description: "",
            parameters: []
        }, {
            name: "recommender",
            description: "",
            parameters: []
        }, {
            name: "privatesite",
            description: "",
            parameters: []
        }, {
            name: "publicsite",
            description: "",
            parameters: []
        }, {
            name: "dummy",
            description: "",
            parameters: []
        }];
	
	return operatorList;
}
  

function defaultDockerImageList() {
        var imageList = [{
            name: "fogflow/nodejs",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "nodejs",
            prefetched: true
        }, {
            name: "fogflow/python",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "python",
            prefetched: false
        }, {
            name: "fogflow/counter",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "counter",
            prefetched: false
        }, {
            name: "fogflow/anomaly",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "anomaly",
            prefetched: false
        }, {
            name: "fogflow/facefinder",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "facefinder",
            prefetched: false
        }, {
            name: "fogflow/connectedcar",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "connectedcar",
            prefetched: false
        }, {
            name: "fiware/iotagent-json",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "iotagent",
            prefetched: false
        }, {
            name: "fogflow/recommender",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "recommender",
            prefetched: false
        }, {
            name: "fogflow/privatesite",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "privatesite",
            prefetched: false
        }, {
            name: "fogflow/publicsite",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "publicsite",
            prefetched: false
        }, {
            name: "fogflow/dummy",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "dummy",
            prefetched: false
        }, {
            name: "geohash",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "geohash",
            prefetched: false
        }, {
            name: "converter",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "converter",
            prefetched: false
        }, {
            name: "predictor",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "predictor",
            prefetched: false
        }, {
            name: "controller",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "controller",
            prefetched: false
        }, {
            name: "detector",
            tag: "latest",
            hwType: "ARM",
            osType: "Linux",
            operatorName: "detector",
            prefetched: false
        }];

    return imageList;
}


	var myToplogyExamples = [{
            topology: { "name": "anomaly-detection", "description": "detect anomaly events in shops", "tasks": [{ "name": "Counting", "operator": "counter", "input_streams": [{ "selected_type": "Anomaly", "selected_attributes": [], "groupby": "ALL", "scoped": true }], "output_streams": [{ "entity_type": "Stat" }] }, { "name": "Detector", "operator": "anomaly", "input_streams": [{ "selected_type": "PowerPanel", "selected_attributes": [], "groupby": "EntityID", "scoped": true }, { "selected_type": "Rule", "selected_attributes": [], "groupby": "ALL", "scoped": false }], "output_streams": [{ "entity_type": "Anomaly" }] }] },
            designboard: { "edges": [{ "id": 2, "block1": 3, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }, { "id": 3, "block1": 2, "connector1": ["outputs", "output", 0], "block2": 3, "connector2": ["in", "input"] }, { "id": 4, "block1": 4, "connector1": ["stream", "output"], "block2": 2, "connector2": ["streams", "input"] }, { "id": 5, "block1": 5, "connector1": ["stream", "output"], "block2": 2, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": 202, "y": -146, "type": "Task", "module": null, "values": { "name": "Counting", "operator": "counter", "outputs": ["Stat"] } }, { "id": 2, "x": -194, "y": -134, "type": "Task", "module": null, "values": { "name": "Detector", "operator": "anomaly", "outputs": ["Anomaly"] } }, { "id": 3, "x": 4, "y": -18, "type": "Shuffle", "module": null, "values": { "selectedattributes": ["all"], "groupby": "ALL" } }, { "id": 4, "x": -447, "y": -179, "type": "EntityStream", "module": null, "values": { "selectedtype": "PowerPanel", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": true } }, { "id": 5, "x": -438, "y": -5, "type": "EntityStream", "module": null, "values": { "selectedtype": "Rule", "selectedattributes": ["all"], "groupby": "ALL", "scoped": false } }] }
        }, 
        {
            topology: { "name": "child-finder", "description": "search for a lost child based on face recognition", "tasks": [{ "name": "childfinder", "operator": "facefinder", "input_streams": [{ "selected_type": "Camera", "selected_attributes": [], "groupby": "EntityID", "scoped": true }, { "selected_type": "ChildLost", "selected_attributes": [], "groupby": "ALL", "scoped": false }], "output_streams": [{ "entity_type": "ChildFound" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }, { "id": 2, "block1": 3, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": 7, "y": -107, "type": "Task", "module": null, "values": { "name": "childfinder", "operator": "facefinder", "outputs": ["ChildFound"] } }, { "id": 2, "x": -292, "y": -161, "type": "EntityStream", "module": null, "values": { "selectedtype": "Camera", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": true } }, { "id": 3, "x": -286, "y": -2, "type": "EntityStream", "module": null, "values": { "selectedtype": "ChildLost", "selectedattributes": ["all"], "groupby": "ALL", "scoped": false } }] }
        }];
    
    var myFogFunctionExamples = [{
            name: "Detector",
            topology: { "name": "Detector", "description": "test", "tasks": [{ "name": "Main", "operator": "detector", "input_streams": [{ "selected_type": "Camera", "selected_attributes": [], "groupby": "EntityID", "scoped": false }], "output_streams": [] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": 123, "y": -99, "type": "Task", "module": null, "values": { "name": "Main", "operator": "detector", "outputs": [] } }, { "id": 2, "x": -194, "y": -97, "type": "EntityStream", "module": null, "values": { "selectedtype": "Camera", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "Detector", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        },
        {
            name: "Test",
            topology: { "name": "Test", "description": "just for a simple test", "tasks": [{ "name": "Main", "operator": "dummy", "input_streams": [{ "selected_type": "Temperature", "selected_attributes": [], "groupby": "EntityID", "scoped": false }], "output_streams": [{ "entity_type": "Out" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": 123, "y": -99, "type": "Task", "module": null, "values": { "name": "Main", "operator": "dummy", "outputs": ["Out"] } }, { "id": 2, "x": -194, "y": -97, "type": "EntityStream", "module": null, "values": { "selectedtype": "Temperature", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "Test", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        },
        {
            name: "PrivateSiteEstimation",
            topology: { "name": "PrivateSiteEstimation", "description": "to estimate the free parking lots from a private parking site", "tasks": [{ "name": "Estimation", "operator": "privatesite", "input_streams": [{ "selected_type": "PrivateSite", "selected_attributes": [], "groupby": "EntityID", "scoped": false }], "output_streams": [{ "entity_type": "Out" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": 26, "y": -47, "type": "Task", "module": null, "values": { "name": "Estimation", "operator": "privatesite", "outputs": ["Out"] } }, { "id": 2, "x": -302, "y": -87, "type": "EntityStream", "module": null, "values": { "selectedtype": "PrivateSite", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "PrivateSiteEstimation", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        }, 
        {
            name: "PublicSiteEstimation",
            topology: { "name": "PublicSiteEstimation", "description": "to estimate the free parking lot from a public parking site", "tasks": [{ "name": "PubFreeLotEstimation", "operator": "publicsite", "input_streams": [{ "selected_type": "PublicSite", "selected_attributes": [], "groupby": "EntityID", "scoped": false }], "output_streams": [{ "entity_type": "Out" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": -37, "y": -108, "type": "Task", "module": null, "values": { "name": "PubFreeLotEstimation", "operator": "publicsite", "outputs": ["Out"] } }, { "id": 2, "x": -340, "y": -128, "type": "EntityStream", "module": null, "values": { "selectedtype": "PublicSite", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "PublicSiteEstimation", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        }, 
        {
            name: "ArrivalTimeEstimation",
            topology: { "name": "ArrivalTimeEstimation", "description": "to estimate when the car will arrive at the destination", "tasks": [{ "name": "CalculateArrivalTime", "operator": "connectedcar", "input_streams": [{ "selected_type": "ConnectedCar", "selected_attributes": [], "groupby": "EntityID", "scoped": false }], "output_streams": [{ "entity_type": "Out" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": -106, "y": -93, "type": "Task", "module": null, "values": { "name": "CalculateArrivalTime", "operator": "connectedcar", "outputs": ["Out"] } }, { "id": 2, "x": -420, "y": -145, "type": "EntityStream", "module": null, "values": { "selectedtype": "ConnectedCar", "selectedattributes": ["all"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "ArrivalTimeEstimation", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        }, 
        {
            name: "ParkingLotRecommendation",
            topology: { "name": "ParkingLotRecommendation", "description": "to recommend where to park around the destination", "tasks": [{ "name": "WhereToParking", "operator": "recommender", "input_streams": [{ "selected_type": "ConnectedCar", "selected_attributes": ["ParkingRequest"], "groupby": "EntityID", "scoped": false }], "output_streams": [{ "entity_type": "Out" }] }] },
            designboard: { "edges": [{ "id": 1, "block1": 2, "connector1": ["stream", "output"], "block2": 1, "connector2": ["streams", "input"] }], "blocks": [{ "id": 1, "x": -14, "y": -46, "type": "Task", "module": null, "values": { "name": "WhereToParking", "operator": "recommender", "outputs": ["Out"] } }, { "id": 2, "x": -379, "y": -110, "type": "EntityStream", "module": null, "values": { "selectedtype": "ConnectedCar", "selectedattributes": ["ParkingRequest"], "groupby": "EntityID", "scoped": false } }] },
            intent: { "topology": "ParkingLotRecommendation", "priority": { "exclusive": false, "level": 0 }, "qos": "Max Throughput", "geoscope": { "scopeType": "global", "scopeValue": "global" } }
        }
    ];
    