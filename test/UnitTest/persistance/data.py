# Payload to persist Operator
test001=\
    {
    }

test002=\
    {
    "name":"counter"
    }

test003=\
    {
    "name":"counter",
    "description":"Hi there, all okay!",
    }

test0 =\
    {
        "name":"counter",
        "description":"Hi there, all okay!",
        "parameters":[]
    }

# Payload to persist FogFunction
test101=\
    {
    }

test102=\
    {
        "id": 'FogFunction.ParkingLotRecommendation'
    }        
test103=\
    {
        "id": 'FogFunction.ParkingLotRecommendation',
        "name": 'ParkingLotRecommendation',
        "topology":
        {
            "name": 'ParkingLotRecommendation',
            "description": 'to recommend where to park around the destination',
            "tasks":[[]]
        }
    }
test104=\
    {
        "id": 'FogFunction.ParkingLotRecommendation',
        "name": 'ParkingLotRecommendation',

        "geoscope":
        {
            "scopeType": 'global', "scopeValue": 'global'
        },
        "status": 'enabled',
        "action": 'UPDATE'

    }

test1 =\
    { 
        "id": 'FogFunction.ParkingLotRecommendation',
        "name": 'ParkingLotRecommendation',
        "topology":
        { 
            "name": 'ParkingLotRecommendation',
            "description": 'to recommend where to park around the destination',
            "tasks":[[]]
        },
        "intent":
        { 
            "topology": 'ParkingLotRecommendation',
            "priority": 
            { 
                "exclusive": False, "level": 0 
            },
        "qos": 'Max Throughput',
        "geoscope": 
        {   
            "scopeType": 'global', "scopeValue": 'global' 
        },
        "status": 'enabled',
        "action": 'UPDATE' 
    }
        }
# Payload to persist DockerImage

test200 =\
    {

    }


test201 =\
    {
    "operater": "counter",
    }

test202 =\
    {
    "operater": "counter",
    "name": "fogflow/counter"
    }

test203 =\
    {
    "operater": "counter",
    "name": "fogflow/counter",
    "hwType": "X86",
    "osType": "Linux"
    }

test2 =\
    {
    "operater": "counter",
    "name": "fogflow/counter",
    "tag": "latest",
    "hwType": "X86",
    "osType": "Linux",
    "prefetched": False
}

# payload to persist Topology
test3 =\
    {
   "description": "detect anomaly events in shops",
   "name": "anomaly-detection",
   "tasks": [
      {
         "input_streams": [
            {
               "groupby": "ALL",
               "scoped": True,
               "selected_attributes": [],
               "selected_type": "Anomaly"
            }
         ],
         "name": "Counting",
         "operator": "counter",
         "output_streams": [
            {
               "entity_type": "Stat32_new"
            }
         ]
      },
      {
         "input_streams": [
            {
               "groupby": "EntityID",
               "scoped": True,
               "selected_attributes": [],
               "selected_type": "PowerPanel"
            },
           {
               "groupby": "ALL",
               "scoped": False,
               "selected_attributes": [],
               "selected_type": "Rule"
            }
         ],
         "name": "Detector",
         "operator": "anomaly",
         "output_streams": [
            {
               "entity_type": "Anomaly32_new"
            }
         ]
      }
   ]
}

test300=\
	{
	}

test301=\
	{
	"description": "detect anomaly events in shops",
    "name": "anomaly-detection"
	}

'''        
test302=\
	{
	"description": "detect anomaly events in shops",
        "name": "anomaly-detection",
        "tasks": [
        {
            "input_streams": [
                {
                "groupby": "ALL",
                "scoped": True,
                "selected_attributes": [],
                "selected_type": "Anomaly"
                }   
                ],
        "name": "Counting",
        "operator": "counter",
        "output_streams": [
            {
               "entity_type": "Stat32_new"
            }
            ]
            }
        }   
'''
#payload to persist service intent
test400=\
    {   
    }
test401=\
    {   
    "topology": "anomaly-detection",
    "id": "ServiceIntent.849ecf56-4590-4493-a982-7b1a257053e2"
    }
test402=\
    {   
    "topology": "anomaly-detection",
    "geoscope": { "scopeType": "global", "scopeValue": "global" },
    }

test4=\
    {   
    "topology": "anomaly-detection",
    "priority": { "exclusive": False, "level": 50 },
    "qos": 'NONE',
    "geoscope": { "scopeType": "global", "scopeValue": "global" },
    "id": "ServiceIntent.849ecf56-4590-4493-a982-7b1a257053e2"
    }

