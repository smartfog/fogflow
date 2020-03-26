subscription_data=\
{
  "description": "A subscription to get info about Room1",
  "subject": {
    "entities": [
      {
        "id": "Room1",
        "type": "Room",
      }
    ],
    "condition": {
      "attrs": [
	  "p3"
      ]
    }
  },
  "notification": {
    "http": {
      "url": "http://192.168.100.162:8888"
    },
    "attrs": [	
	"p1",
	"p2",
	"p3"
    ]
  },
  "expires": "2040-01-01T14:00:00.00Z",
  "throttling": 5
}
	

#data to test the following code for broker.thinBroker.go:946
'''
	 subReqv2 := SubscriptionRequest{}

        err := r.DecodeJsonPayload(&subReqv2)
        if err != nil {
                rest.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

'''

subscriptionWrongPaylaod=\
{
  "description": "A subscription to get info about Room1",
  "subject": {
    "entities": [
      {
        "id": "Room1",
        "type": "Room",
	"ispattern":"false"
      }
    ],
    "condition": {
      "attrs": [
	  "p3"
      ]
    }
  },
  "notification": {
    "http": {
      "url": "http://192.168.100.162:8888"
    },
    "attrs": [	
	"p1",
	"p2",
	"p3"
    ]
  },
  "expires": "2040-01-01T14:00:00.00Z",
  "throttling": 5
}

v1SubData=\
{
  "entities": [
    {
      "id": "Room1",
      "type": "Room",
    }
  ],
  "reference": "http://192.168.100.162:8668/ngsi10/updateContext"
}

updateDataWithupdateaction=\
{
"contextElements": [
{
"entityId": {
"id": "Room1",
"type": "Room"
},
"attributes": [
{
"name": "p1",
"type": "float",
"value": 60
},
{
"name": "p3",
"type": "float",
"value": 69
},
{
"name": "p2",
"type": "float",
"value": 32
}
],
"domainMetadata": [
{
"name": "location",
"type": "point",
"value": {
"latitude": 49.406393,
"longitude": 8.684208
}
}
]
}
],
"updateAction": "UPDATE"
}

createDataWithupdateaction=\
{
"contextElements": [
{
"entityId": {
"id": "Room1",
"type": "Room"
},
"attributes": [
{
"name": "p1",
"type": "float",
"value": 90
},
{
"name": "p3",
"type": "float",
"value": 70
},
{
"name": "p2",
"type": "float",
"value": 12
}
],
"domainMetadata": [
{
"name": "location",
"type": "point",
"value": {
"latitude": 49.406393,
"longitude": 8.684208
}
}
]
}
],
"updateAction": "CRETAE"
}

deleteDataWithupdateaction=\
{
"contextElements": [
{
"entityId": {
"id": "Room1",
"type": "Room"
},
"attributes": [
{
"name": "p1",
"type": "float",
"value": 12
},
{
"name": "p3",
"type": "float",
"value": 13
},
{
"name": "p2",
"type": "float",
"value": 14
}
],
"domainMetadata": [
{
"name": "location",
"type": "point",
"value": {
"latitude": 49.406393,
"longitude": 8.684208
}
}
]
}
],
"updateAction": "DELETE"
}

