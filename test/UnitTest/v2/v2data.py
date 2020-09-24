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
      "url": "http://0.0.0.0:8888/accumulate"
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
      "url": "http://0.0.0.0:8888/accumulate"
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
  "reference": "http://0.0.0.0:8888/accumulate"
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

subdata1=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial10",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 69
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 75
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}

subdata2=\
{
  "description": "A subscription to get info about RoomTrial10",
  "subject": {
    "entities": [
      {
        "id": "RoomTrial10",
        "type": "Room"
      }
    ],
    "condition": {
      "attrs": [
        "pressure"
      ]
    }
  },
  "notification": {
    "http": {
      "url": "http://0.0.0.0:8888/accumulate"
    },
    "attrs": [
      "temperature"
    ]
  },
  "expires": "2040-01-01T14:00:00.00Z",
  "throttling": 5
}

subdata3=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial10",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 50
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 80
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}

subdata4=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial20",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 69
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 75
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}



subdata5=\
{
  "description": "A subscription to get info about RoomTrial20",
  "subject": {
    "entities": [
      {
        "id": "RoomTrial20",
        "type": "Room"
      }
    ],
    "condition": {
      "attrs": [
        "pressure"
      ]
    }
  },
  "notification": {
    "http": {
      "url": "http://0.0.0.0:8888/accumulate"
    },
    "attrs": [
      "temperature"
    ]
  },
  "expires": "2040-01-01T14:00:00.00Z",
  "throttling": 5
}



subdata6=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial20",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 40
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 85
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}



subdata7=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial30",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 69
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 75
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}



subdata8=\
{
  "description": "A subscription to get info about RoomTrial30",
  "subject": {
    "entities": [
      {
        "id": "RoomTrial30",
        "type": "Room"
      }
    ],
    "condition": {
      "attrs": [
        "pressure"
      ]
    }
  },
  "notification": {
    "http": {
      "url": "http://0.0.0.0:8888/accumulate"
    },
    "attrs": [
      "temperature"
    ]
  },
  "expires": "2040-01-01T14:00:00.00Z",
  "throttling": 5
}





subdata9=\
{
            "contextElements": [
                {
                    "entityId": {
                        "id": "RoomTrial30",
                        "type": "Room"
                        },
                    "attributes": [
                            {
                            "name": "temperature",
                            "type": "float",
                            "value": 44
                            },
                            {
                            "name": "pressure",
                            "type": "float",
                            "value": 60
                            }
                        ],
                    "domainMetadata": [
                            {
                            "name": "location",
                            "type": "point",
                            "value": {
                            "latitude": -33.1,
                            "longitude": -1.1
                            }}
                        ]
                }
            ],
            "updateAction": "UPDATE"
}

