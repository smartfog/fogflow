# Payload to create entity without passing Link Header
subdata=\
{
           "id": "urn:ngsi-ld:Vehicle:A020",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload to create entity with context in link header  
subdata1=\
{
           "id": "urn:ngsi-ld:Vehicle:A100",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }


# Payload to create entity with context inside payload
subdata2=\
 {
              "@context": [{
              "Vehicle": "http://example.org/vehicle/Vehicle",
              "brandName": "http://example.org/vehicle/brandName",
              "speed": "http://example.org/vehicle/speed",
              "isParked": {
                             "@type": "@id",
                             "@id": "http://example.org/common/isParked"
              },
              "providedBy": {
                               "@type": "@id",
                               "@id": "http://example.org/common/providedBy"
               }
           }],
           "id": "urn:ngsi-ld:Vehicle:A4580",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                         "type": "Property",
                         "value": 80
              },
              "createdAt": "2017-07-29T12:00:04",
              "location": {
                            "type": "GeoProperty",
                            "value": {
                                        "type": "Point",
                                        "coordinates": [-8.5, 41.2]
                             }
                }
       }

# create entity with context in Link header and request payload is already expanded
subdata3=\
 {
             "https://example.org/vehicle/brandName": [
             {
                  "@type": [
                              "https://uri.etsi.org/ngsi-ld/Property"
                   ],
                   "https://uri.etsi.org/ngsi-ld/hasValue": [
                            {
                                 "@value": "Mercedes"
                            }
                      ]
                }
          ],
            "https://uri.etsi.org/ngsi-ld/createdAt": [
             {
                  "@type": "https://uri.etsi.org/ngsi-ld/DateTime",
                  "@value": "2017-07-29T12:00:04"
              }
          ],
            "@id": "urn:ngsi-ld:Vehicle:A8866",
            "https://example.org/common/isParked": [
             {
                        "https://uri.etsi.org/ngsi-ld/hasObject": [
                         {
                                 "@id": "urn:ngsi-ld:OffStreetParking:Downtown1"
                          }
                       ],
                        "https://uri.etsi.org/ngsi-ld/observedAt": [
                         {
                                "@type": "https://uri.etsi.org/ngsi-ld/DateTime",
                                "@value": "2017-07-29T12:00:04"
                          }
                       ],
                         "https://example.org/common/providedBy": [
                          {
                                   "https://uri.etsi.org/ngsi-ld/hasObject": [
                                   {
                                           "@id": "urn:ngsi-ld:Person:Bob"
                                   }
                                ],
                                "@type": [
                                            "https://uri.etsi.org/ngsi-ld/Relationship"
                                  ]
                           }
                        ],
                          "@type": [
                                      "https://uri.etsi.org/ngsi-ld/Relationship"
                            ]
                  }
            ],
             "https://uri.etsi.org/ngsi-ld/location": [
              {
                        "@type": [
                                    "https://uri.etsi.org/ngsi-ld/GeoProperty"
                          ],
                        "https://uri.etsi.org/ngsi-ld/hasValue": [
                         {
                               "@value": "{ \"type\":\"Point\", \"coordinates\":[ -8.5, 41.2 ] }"
                          }
                        ]
               }
            ],
             "https://example.org/vehicle/speed": [
              {
                       "@type": [
                                   "https://uri.etsi.org/ngsi-ld/Property"
                        ],
                        "https://uri.etsi.org/ngsi-ld/hasValue": [
                         {
                               "@value": 80
                          }
                        ]
                }
            ],
             "@type": [
                        "https://example.org/vehicle/Vehicle"
            ]

   }


# Payload  to append additional attributes to an existing entity 
subdata4=\
{
      "id": "urn:ngsi-ld:Vehicle:A100",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "BMW"
       }
 }

subdata4b=\
{
      "id": "urn:ngsi-ld:Vehicle:A500",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "BMW"
       }
 }


subdata4c=\
{
      "id": "urn:ngsi-ld:Vehicle:A3000",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "BMW"
       }
 }


# Payload  to patch  update specific attributes of an existing entity A100
subdata5=\
{
      "id": "urn:ngsi-ld:Vehicle:A100",
      "type": "Vehicle",
      "brandName1": {
                           "type": "Property",
                           "value": "AUDI"
         }
  }

# Payload  to update the value of a specific attribute of an existing entity with wrong payload
subdata6=\
{
         "value": "MARUTI"
  }

# Payload to create a new Subscription to with context in Link header
subdata10=\
{
	     "id": "urn:ngsi-ld:Subscription:7",
             "type": "Subscription",
             "entities": [{
                             "idPattern": ".*",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName"],
             "notification": {
                             "attributes": ["brandName"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }

# Payload to create entity which is to be tested for delete attribute request
subdata11=\
{
           "id": "urn:ngsi-ld:Vehicle:A500",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload to Update the subscription
subdata12=\
{
             "id": "urn:ngsi-ld:Subscription:7",
             "type": "Subscription",
             "entities": [{
                             "type": "Vehicle"
               }],
             "watchedAttributes": ["http://example.org/vehicle/brandName2"],
             "q":"http://example.org/vehicle/brandName2!=Mercedes",
             "notification": {
             "attributes": ["http://example.org/vehicle/brandName2"],
             "format": "keyValues",
             "endpoint": {
                             "uri": "http://127.0.0.1:8888/ld-notify",
                             "accept": "application/json"
              }
           }
       }

# Payload to  create entity  without passing Header
subdata13=\
{
           "id": "urn:ngsi-ld:Vehicle:A600",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload to update entity with different header and posting duplicate attribute
subdata14=\
{
      "id": "urn:ngsi-ld:Vehicle:A800",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "MARUTI"
       }
 }



subdata14b=\
{
      "id": "Vehicle:A800",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "MARUTI"
       }
 }



subdata14c=\
{
      "id": "urn:ngsi-ld:Vehicle:A700",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "MARUTI"
       }
 }


subdata14d=\
{
      "id": "urn:ngsi-ld:Vehicle:A900",
      "type": "Vehicle",
      "brandName1": {
                          "type": "Property",
                          "value": "MARUTI"
       }
 }


# Payload to  Update entity with different headers and passing inappropriate payload
subdata15=\
{
        "id": "urn:ngsi-ld:Vehicle:A100",
	"type": "Vehicle",
        "brandName1": {
                          
         }
  }


# Payload to create entity without attribute
subdata16=\
{
           "id": "urn:ngsi-ld:Vehicle:A700",
           "type": "Vehicle",
           "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload to create entity without any attributes
subdata17=\
{
           "id": "urn:ngsi-ld:Vehicle:A900",
           "type": "Vehicle",
            "createdAt": "2017-07-29T12:00:04"
     }

# Payload to create entity without any attribute to be tested for delete attribute request
subdata18=\
{
           "id": "urn:ngsi-ld:Vehicle:A501",
           "type": "Vehicle",
            "createdAt": "2017-07-29T12:00:04"
     }


#Payload to update a specific subscription based on subscription id, with context in Link header and different payload
subdata20=\
{
             "id": "urn:ngsi-ld:Subscription:7",
             "type": "Subscription",
             "entities": [{
                             "type": "Vehicle"
               }],
             
             "notification": {
             "format": "keyValues",
             "endpoint": {
                             "uri": "http://127.0.0.1:8888/ld-notify",
                             "accept": "application/json"
              }
           }
       }

# Inappropriate payload to perform patch update 
subdata25=\
{
	"type":"ContextSourceRegistration"
}

# Empty Payload
subdata26=\
{
}

# Payload to create entity to check for CreatedAt and ModifiedAt
subdata27=\
{
           "id": "urn:ngsi-ld:Vehicle:A6000",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload to create csource registration with idPattern
subdata28=\
{
             "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a7000",
             "type": "ContextSourceRegistration",
             "name": "NameExample",
             "description": "DescriptionExample",
             "information": [
             {
                     "entities": [
                      {
                             "id": "urn:ngsi-ld:Vehicle:A.*",
                             "type": "Vehicle"
                      }
                   ],
                     "properties": [
                     "brandName",
                     "speed"
                   ],
                     "relationships": [
                     "isParked"
                     ]
             },
             {
                     "entities": [
                     {
                             "idPattern": "pqr$",
                             "type": "OffStreetParking"
                     }
                   ]
             }
            ],
             "endpoint": "http://127.0.0.1:8888/csource",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

# Payload to create Subscription to check for Modified At and Created At in susbcription
subdata29=\
{
	     "id": "urn:ngsi-ld:Subscription:8",
             "type": "Subscription",
             "entities": [{
                             "idPattern": ".*",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName"],
             "notification": {
                             "attributes": ["brandName"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }

# Payload to create a Subscription with id as urn:ngsi-ld:Subscription:10
subdata30=\
{
	     "id": "urn:ngsi-ld:Subscription:10",
             "type": "Subscription",
             "entities": [{
                             "idPattern": ".*",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName"],
             "notification": {
                             "attributes": ["brandName"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }

# Payload to update the corresponding subscription
subdata31=\
{
             "id": "urn:ngsi-ld:Subscription:10",
             "type": "Subscription",
             "entities": [{
                             "type": "Vehicle"
               }],
             "watchedAttributes": ["http://example.org/vehicle/brandName2"],
             "q":"http://example.org/vehicle/brandName2!=Mercedes",
             "notification": {
             "attributes": ["http://example.org/vehicle/brandName2"],
             "format": "keyValues",
             "endpoint": {
                             "uri": "http://127.0.0.1:8888/ld-notify",
                             "accept": "application/json"
              }
           }
       }
# Payload to create an entity which is to be checked for delete request
subdata32=\
{
           "id": "urn:ngsi-ld:Vehicle:A999",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                          "observedAt": "2017-07-29T12:00:04",
                          "providedBy": {
                                          "type": "Relationship",
                                          "object": "urn:ngsi-ld:Person:Bob"
                           }
             },
             "speed": {
                        "type": "Property",
                        "value": 80
              },
             "createdAt": "2017-07-29T12:00:04",
             "location": {
                            "type": "GeoProperty",
                            "value": {
                                      "type": "Point",
                                      "coordinates": [-8.5, 41.2]
                            }
             }
     }

# Payload for entity creation with nested property with context in payload
subdata33=\
{
                "id": "urn:ngsi-ld:Vehicle:B990",
                "type": "Vehicle",
                "brandName": {
                                "type": "Property",
                                "value": "Mercedes"
                },
                "isParked": {
                                "type": "Relationship",
                                "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                                "providedBy": {
                                                "type": "person",
                                                "object": "urn:ngsi-ld:Person:Bob"
                                },
                                "parkingDate": {
                                                "type": "Property",
                                                "value": "2017-07-29T12:00:04"
                                },
                                "availableSpotNumber": {
                                                "type": "Property",
                                                "value": "121",
                                                "reliability": {
                                                                "type": "Property",
                                                                "value": "0.7"
                                                },
                                                "providedBy": {
                                                                "type": "relationship",
                                                                "object": "urn:ngsi-ld:camera:c1"
                                                }
                                }
                },
                "@context": [
                                {
                                                "Vehicle": "http://example.org/vehicle/Vehicle",
                                                "brandName": "http://example.org/vehicle/brandName",
                                                "speed": "http://example.org/vehicle/speed",
                                                "isParked": {
                                                                "@type": "@id",
                                                                "@id": "http://example.org/common/isParked"
                                                },
                                                "providedBy": {
                                                                "@type": "@id",
                                                                "@id": "http://example.org/common/providedBy"
                                                }
                                }
                ]
}

# Payload to create nested entity with context in link header
subdata34=\
{
                "id": "urn:ngsi-ld:Vehicle:C001",
                "type": "Vehicle",
                "brandName": {
                                "type": "Property",
                                "value": "Mercedes"
                },
                "isParked": {
                                "type": "Relationship",
                                "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                                "providedBy": {
                                                "type": "person",
                                                "object": "urn:ngsi-ld:Person:Bob"
                                },
                                "parkingDate": {
                                                "type": "Property",
                                                "value": "2017-07-29T12:00:04"
                                },
                                "availableSpotNumber": {
                                                "type": "Property",
                                                "value": "121",
                                                "reliability": {
                                                                "type": "Property",
                                                                "value": "0.7"
                                                },
                                                "providedBy": {
                                                                "type": "relationship",
                                                                "object": "urn:ngsi-ld:camera:c1"
                                                }
                                }
                }
}


#Payload to create new entity with different context link
subdata35=\
{
	"id": "urn:ngsi-ld:Vehicle:A909",
	"type": "Vehicle",
	"brandName": {
		"type": "Property",
		"value": "Mercedes"
	},
	"isParked": {
		"type": "Relationship",
		"object": "urn:ngsi-ld:OffStreetParking:Downtown1",
		"observedAt": "2017-07-29T12:00:04",
		"providedBy": {
			"type": "Relationship",
			"object": "urn:ngsi-ld:Person:Bob"
		}
	},
	"speed": {
		"type": "Property",
		"value": 80
	},
	"createdAt": "2017-07-29T12:00:04",
	"location": {
		"type": "GeoProperty",
		"value": {
			"type": "Point",
			"coordinates": [-8.5, 41.2]
		}
	}
}

subdata36=\
{
             "id": "urn:ngsi-ld:Subscription:19",
             "type": "Subscription",
             "entities": [{
                             "idPattern": ".*",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName","speed"],
             "notification": {
                             "attributes": ["brandName","speed"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }


subdata37=\
{
             "id": "urn:ngsi-ld:Subscription:20",
             "type": "Subscription",
             "entities": [{
                             "id": "urn:ngsi-ld:Vehicle:A3000",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName","speed"],
             "notification": {
                             "attributes": ["brandName","speed"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             },
	     "context": ["https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"]
      }


subdata38=\
{
    "id": "urn:ngsi-ld:Vehicle:A3000",
    "type": "Vehicle",
    "brandName": {
        "type": "Property",
        "value": "Mercedes"
    },
    "isParked": {
        "type": "Relationship",
        "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
        "observedAt": "2017-07-29T12:00:04",
        "providedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Person:Bob"
        }
    },
    "speed": {
        "type": "Property",
        "value": 80
    },
    "createdAt": "2017-07-29T12:00:04",
    "location": {
        "type": "GeoProperty",
        "value": {
            "type": "Point",
            "coordinates": [-8.5, 41.2]
        }
    }
}

subdata39=\
{
             "id": "urn:ngsi-ld:Subscription:020",
             "type": "Subscription",
             "entities": [{
                             "idPattern": ".*",
                             "type": "Vehicle"
             }],
             "watchedAttributes": ["brandName","speed"],
             "notification": {
                             "attributes": ["brandName","speed"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
     }

subdata40=\
{
  "id": "urn:ngsi-ld:Vehicle:A4000",
 "type": "Vehicle",
 "brandName": {
  "type": "Property",
  "value": "Mercedes"
 },
 "isParked": {
  "type": "Relationship",
  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
  "observedAt": "2017-07-29T12:00:04",
  "providedBy": {
   "type": "Relationship",
   "object": "urn:ngsi-ld:Person:Bob"
  }
 },
 "speed": {
  "type": "Property",
  "value": 80
 },
 "createdAt": "2017-07-29T12:00:04",
 "location": {
  "type": "GeoProperty",
  "value": {
   "type": "Point",
   "coordinates": [-8.5, 41.2]
  }
 }
}

# for update request
subdata41=\
{
	"id": "urn:ngsi-ld:Vehicle:A4000",
        "type": "Vehicle",
	"brandName1": {
		"type" : "Property",
		"value": "BMW1"
	}
}

# to test if instance id is fetched or not
subdata42=\
{
        "id": "urn:ngsi-ld:Vehicle:C001",
        "type": "Vehicle",
        "brandName1": {
                "type" : "Property",
                "value": "BMW1",
                "instanceId": "instance1"
        }
}

# to test if datasetId is fetched or not
subdata43=\
{
        "id": "urn:ngsi-ld:Vehicle:C002",
        "type": "Vehicle",
        "brandName1": {
                "type" : "Property",
                "value": "BMW1",
                "datasetId": "dataset1"
        }
}

# to test if error is thrown if entities are missing
subdata44=\
{
             "id": "urn:ngsi-ld:Subscription:79",
             "type": "Subscription",
             "watchedAttributes": ["brandName"],
             "notification": {
                             "attributes": ["brandName"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }

# to test that no other type is taken for NGSI-LD subscription other than Subscription
subdata45=\
{
             "id": "urn:ngsi-ld:Subscription:80",
             "type": "SubscriptionXY",
             "watchedAttributes": ["brandName"],
             "notification": {
                             "attributes": ["brandName"],
                             "format": "keyValues",
                             "endpoint": {
                                             "uri": "http://127.0.0.1:8888/ld-notify",
                                             "accept": "application/json"
                               }
             }
      }
