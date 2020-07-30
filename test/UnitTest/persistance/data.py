# Payload to persist Operator 
test0=\
{
	"contextElements": [{
		"entityId": {
			"type": "Operator",
			"id": "test011"
		},
		"attributes": [{
			"name": "designboard",
			"type": "object",
			"value": {}
		}, {
			"name": "operator",
			"type": "object",
			"value": {
				"description": "",
				"name": "recommender",
				"parameters": []
			}
		}],
		"domainMetadata": [{
			"name": "location",
			"type": "global",
			"value": "global"
		}]
	}],
	"updateAction": "UPDATE"
}

# Payload to persist FogFunction
test1=\
{
  	"contextElements": [{
  			"entityId": {
  				"type": "FogFunction",
  				"id": "test2"
  			},
  			"attributes": [{
  				"name": "name",
  				"type": "string",
  				"value": "Test"
  			}, {
  				"name": "topology",
  				"type": "object",
  				"value": {
  					"description": "just for a simple test",
  					"name": "Test",
  					"tasks": [{
  						"input_streams": [{
  							"groupby": "EntityID",
  							"scoped": "false",
  							"selected_attributes": [],
  							"selected_type": "Temperature"
  						}],
  						"name": "Main",
  						"operator": "dummy",
  						"output_streams": [{
  							"entity_type": "Out"
  						}]
  					}]
  				}
  			}, {
  				"name": "designboard",
  				"type": "object",
  				"value": {
  					"blocks": [{
  						"id": 1,
  						"module": "null",
  						"type": "Task",
  						"values": {
  							"name": "Main",
  							"operator": "dummy",
  							"outputs": ["Out"]
  						},
  						"x": 123,
  						"y": -99
  					}, {
  						"id": 2,
  						"module": "null",
  						"type": "EntityStream",
  						"values": {
  							"groupby": "EntityID",
  							"scoped": "false",
  							"selectedattributes": ["all"],
  							"selectedtype": "Temperature"
  						},
  						"x": -194,
  						"y": -97
  					}],
  					"edges": [{
  						"block1": 2,
  						"block2": 1,
  						"connector1": ["stream", "output"],
  						"connector2": ["streams", "input"],
  						"id": 1
  					}]
  				}
  			}, {
  				"name": "intent",
  				"type": "object",
  				"value": {
  					"geoscope": {
  						"scopeType": "global",
  						"scopeValue": "global"
  					},
  					"priority": {
  						"exclusive": "false",
  						"level": 0
  					},
  					"qos": "Max Throughput",
  					"topology": "Test"
  				}
  			}, {
  				"name": "status",
  				"type": "string",
  				"value": "enabled"
  			}],
  			"domainMetadata": [{
  				"name": "location",
  				"type": "global",
  				"value": "global"
  			}]
  		}

  	],
  	"updateAction": "UPDATE"
  } 
    
# Payload to persist DockerImage
test2=\
{
	"contextElements": [{
		"entityId": {
			"type": "DockerImage",
			"id": "test3"
		},
		"attributes": [{
			"name": "image",
			"type": "string",
			"value": "fogflow/counter"
		}, {
			"name": "tag",
			"type": "string",
			"value": "latest"
		}, {
			"name": "hwType",
			"type": "string",
			"value": "X86"
		}, {
			"name": "osType",
			"type": "string",
			"value": "Linux"
		}, {
			"name": "operator",
			"type": "string",
			"value": "counter"
		}, {
			"name": "prefetched",
			"type": "boolean",
			"value": "false"
		}],
		"domainMetadata": [{
			"name": "operator",
			"type": "string",
			"value": "counter"
		}, {
			"name": "location",
			"type": "global",
			"value": "global"
		}]
	}],
	"updateAction": "UPDATE"
}

# payload to persist Topology
test3=\
 {
  	"contextElements": [{
  			"entityId": {
  				"type": "Topology",
  				"id": "test4"
  			},
  			"attributes": [{
  				"name": "status",
  				"type": "string",
  				"value": "enabled"
  			}, {
  				"name": "designboard",
  				"type": "object",
  				"value": {
  					"blocks": [{
  						"id": 1,
  						"module": "null",
  						"type": "Task",
  						"values": {
  							"name": "Counting",
  							"operator": "counter",
  							"outputs": ["Stat"]
  						},
  						"x": 202,
  						"y": -146
  					}, {
  						"id": 2,
  						"module": "null",
  						"type": "Task",
  						"values": {
  							"name": "Detector",
  							"operator": "anomaly",
  							"outputs": ["Anomaly"]
  						},
  						"x": -194,
  						"y": -134
  					}, {
  						"id": 3,
  						"module": "null",
  						"type": "Shuffle",
  						"values": {
  							"groupby": "ALL",
  							"selectedattributes": ["all"]
  						},
  						"x": 4,
  						"y": -18
  					}, {
  						"id": 4,
  						"module": "null",
  						"type": "EntityStream",
  						"values": {
  							"groupby": "EntityID",
  							"scoped": "true",
  							"selectedattributes": ["all"],
  							"selectedtype": "PowerPanel"
  						},
  						"x": -447,
  						"y": -179
  					}, {
  						"id": 5,
  						"module": "null",
  						"type": "EntityStream",
  						"values": {
  							"groupby": "ALL",
  							"scoped": "false",
  							"selectedattributes": ["all"],
  							"selectedtype": "Rule"
  						},
  						"x": -438,
  						"y": -5
  					}],
  					"edges": [{
  						"block1": 3,
  						"block2": 1,
  						"connector1": ["stream", "output"],
  						"connector2": ["streams", "input"],
  						"id": 2
  					}, {
  						"block1": 2,
  						"block2": 3,
  						"connector1": ["outputs", "output", 0],
  						"connector2": ["in", "input"],
  						"id": 3
  					}, {
  						"block1": 4,
  						"block2": 2,
  						"connector1": ["stream", "output"],
  						"connector2": ["streams", "input"],
  						"id": 4
  					}, {
  						"block1": 5,
  						"block2": 2,
  						"connector1": ["stream", "output"],
  						"connector2": ["streams", "input"],
  						"id": 5
  					}]
  				}
  			}, {
  				"name": "template",
  				"type": "object",
  				"value": {
  					"description": "detect anomaly events in shops",
  					"name": "anomaly-detection",
  					"tasks": [{
  						"input_streams": [{
  							"groupby": "ALL",
  							"scoped": "true",
  							"selected_attributes": [],
  							"selected_type": "Anomaly"
  						}],
  						"name": "Counting",
  						"operator": "counter",
  						"output_streams": [{
  							"entity_type": "Stat"
  						}]
  					}, {
  						"input_streams": [{
  							"groupby": "EntityID",
  							"scoped": "true",
  							"selected_attributes": [],
  							"selected_type": "PowerPanel"
  						}, {
  							"groupby": "ALL",
  							"scoped": "false",
  							"selected_attributes": [],
  							"selected_type": "Rule"
  						}],
  						"name": "Detector",
  						"operator": "anomaly",
  						"output_streams": [{
  							"entity_type": "Anomaly"
  						}]
  					}]
  				}
  			}],
  			"domainMetadata": [{
  				"name": "location",
  				"type": "global",
  				"value": "global"
  			}]
  		}

  	],
  	"updateAction": "UPDATE"
  }

