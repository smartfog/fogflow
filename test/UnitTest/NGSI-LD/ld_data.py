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
      "id": "urn:ngsi-ld:Vehicle:A3000b",
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
                                                "type": "Relationship",
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
                                                                "type": "Relationship",
                                                                "object": "urn:ngsi-ld:camera:c1"
                                                }
                                }
                }
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
                                                "type": "Relationship",
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
                                                                "type": "Relationship",
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

subdata38b=\
{
    "id": "urn:ngsi-ld:Vehicle:A3000b",
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
	     "entities": [{
                             "id": "urn:ngsi-ld:Vehicle:A3000",
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


# to test for single entity creation with upsert API

subdata46=\
{
   "id": "urn:ngsi-ld:Vehicle:A106",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}


# to test for single entity creation with upsert API

subdata47=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A1066",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A1033",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]


# to test for empty payload

subdata48=\
[
{

}
]


# to test multiple entity creation with one empty payload in array

subdata49=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A111",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
}
]

# to test multiple entity creation with missing id in one entity

subdata50=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A1060",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]


# to test multiple entity creation with missing id in one 

subdata51=\
[
{
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A1030",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]


# to test multiple entity creation with id in any entity

subdata52=\
[
{
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test multiple entity creation with missig type in first entity

subdata53=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A500",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
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
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]


# to test for multiple entity creation with missing type in second entity

subdata54=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A700",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A800",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test creation of multiple entity with missing type in all entities

subdata55=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A900",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A10001",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test multiple entity creation with id not as a uri in first entity

subdata56=\
[
{
   "id": "A250",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A350",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test multiple entity creation with id not as a uri for second entity

subdata57=\
[
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
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "A202",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test multiple entity creation with id not as a uri in both entities

subdata58=\
[
{
   "id": "AB",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "CD",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test for multiple entity creation with array in attributes of first entity

subdata59=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A0001",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": [12,13]
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A0002",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test multiple entities creation with array in attributes of those entities

subdata60=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A0101",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": ["Mercedes","Benz"]
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:A9090",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": ["Mercedes","Maruti"]
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test updation using upsert
subdata61=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A0210",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:APC",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test updation through upsert API

subdata62=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A0210",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Maruti"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:APC",
   "type": "Vehicle",
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

# to test subcription with upsert

subdata63=\
{
             "id": "urn:ngsi-ld:Subscription:Upsert",
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


# fire update

subdata64=\
[
{
   "id": "urn:ngsi-ld:Vehicle:A0210",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "MarutiBenz"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
},
{
   "id": "urn:ngsi-ld:Vehicle:APC",
   "type": "Vehicle",
   "brandName": {
                  "type": "Property",
                  "value": "Mercedes"
    },
    "isParked": {
                  "type": "Relationship",
                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
                  "providedBy": {
                                  "type": "Relationship",
                                  "object": "urn:ngsi-ld:Person:Bob"
                   }
     },
     "speed": {
                "type": "Property",
                "value": 120
      },
     "location": {
                    "type": "GeoProperty",
                    "value": {
                              "type": "Point",
                              "coordinates": [-8.5, 41.2]
                    }
     }
}
]

testData74=\
 [{
  	"id": "urn:ngsi-ld:Device:water001",
  	"type": "Device",
  	"on": {
  		"type": "Property",
  		"value": " "
  	}
  }]

testData75=\
 {
        "id": "urn:ngsi-ld:Device:water001",
        "type": "Device",
        "on": {
                "type": "Property",
                "value": " "
        }
  }

subData80 = \
{
                "type": "Subscription",
                "entities": [{
			"id":"urn:ngsi-ld:Device:water080",
                       "type": "Device"
                }],
              "notification": {
                  "format": "normalized",
                  "endpoint": {
                           "uri": "http://127.0.0.1:8888",
                           "accept": "application/ld+json"
                   }
               }
 }

upsertMultipleCommand = \
[
  {
  	"id": "urn:ngsi-ld:Device:water001",
  	"type": "Device",
  	"on1": {
  		"type": "Property",
  		"value": " "
  	}
   },
   {
        "id": "urn:ngsi-ld:Device:water002",
        "type": "Device",
        "on2": {
                "type": "Property",
                "value": " "
        }
   },
   {
        "id": "urn:ngsi-ld:Device:water003",
        "type": "Device",
        "on3": {
                "type": "Property",
                "value": " "
        }
   }

]

upsertCommand = \
[
  {
        "id": "urn:ngsi-ld:Device:water001",
        "type": "Device",
        "on1": {
                "type": "Property",
                "value": " "
        }
   }

]


upsertCommand80 = \
[
  {
        "id": "urn:ngsi-ld:Device:water080",
        "type": "Device",
        "on1": {
                "type": "Property",
                "value": " "
        }
   }

]



DelData= \
 [{
 	"id": "urn:ngsi-ld:Vehicle:A109",
 	"type": "Device",
 	"brandName": {
 		"type": "Property",
 		"value": "xyzeee"
 	},
 	"isParked": {
 		"type": "Relationship",
 		"object": "urn:ngsi-ld:OffStreetParking:Downtown1",
 		"providedBy": {
 			"type": "Relationship",
 			"object": "urn:ngsi-ld:Person:Bob"
 		}
 	},
 	"speed": {
 		"type": "Property",
 		"value": 30
 	},
 	"location": {
 		"type": "GeoProperty",
 		"value": {
 			"type": "Point",
 			"coordinates": [-8.5, 41.2]
 		}
 	}
 }]
test89 = \
[
  {
        "id": "urn:ngsi-ld:Device:test89",
        "type": "Device",
        "on1": {
                "type": "Property",
                "value": " "
        }
   }

]

#NIL in Property
subdata91=\
{
           "id": "urn:ngsi-ld:Vehicle:AQP",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Nil"
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

#null in property

subdata92=\
{
           "id": "urn:ngsi-ld:Vehicle:AQQ",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "null"
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

#Nil in Relationship Object

subdata93=\
{
           "id": "urn:ngsi-ld:Vehicle:AMNM",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercer"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "nil",
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

#Null in Relationship Object

subdata94=\
{
           "id": "urn:ngsi-ld:Vehicle:AXY",
           "type": "Vehicle",
           "brandName": {
                          "type": "Property",
                          "value": "Mercedes"
            },
            "isParked": {
                          "type": "Relationship",
                          "object": "null",
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


#Payload to test heart health predictor use case

subdata95=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor22",
        "type": "HeartSensor",
        "age": {
            "type": "Property",
            "value": 59.5
        },
        "sex": {
            "type": "Property",
            "value": 1
        },
        "cp": {
            "type": "Property",
            "value": 3
        },
        "trestbps": {
            "type": "Property",
            "value": 145
        },
        "chol": {
            "type": "Property",
            "value": 233
        },
        "fbs": {
            "type": "Property",
            "value": 1
        },
        "restecg": {
            "type": "Property",
            "value": 0
        },
        "thalach": {
            "type": "Property",
            "value": 150
        },
        "exang": {
            "type": "Property",
            "value": 0
        },
        "oldpeak": {
            "type": "Property",
            "value": 2.3
        },
        "slope": {
            "type": "Property",
            "value": 0
        },
        "ca": {
            "type": "Property",
            "value": 0
        },
        "thal": {
            "type": "Property",
            "value": 1
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to check Heart Health use case with missing attributes

subdata96=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor23",
        "type": "HeartSensor",
        "trestbps": {
            "type": "Property",
            "value": 145
        },
        "chol": {
            "type": "Property",
            "value": 233
        },
        "ca": {
            "type": "Property",
            "value": 0
        },
        "thal": {
            "type": "Property",
            "value": 1
            },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to update Heart Health use case attribute with wrong payload
subdata97=\
[
    {
        "chol": {
            "type": "Property",
            "value": 233
       },
    }
]


#Payload to check Heart Health use case with empty values of the attributes

subdata98=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor24",
        "type": "HeartSensor",
        "age": {
            "type": " ",
            "value": 59.5
        },
    }
]

#Payload to check NGSI-LD Heart Health use case with different type in attribute

subdata99=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor25",
        "type": "HeartSensor",
        "age": {
            "type": 123,
            "value": 59.5
        },
    }
]

#Payload to check NGSI-LD Heart Health use case with type empty

subdata100=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor26",
        "type":" " ,
        "age": {
            "type": "Property",
            "value": 59.5
        },
        "sex": {
            "type": "Property",
            "value": 1
        },
        "cp": {
            "type": "Property",
            "value": 3
        },
        "trestbps": {
            "type": "Property",
            "value": 145
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]
#Payload to check NGSI-LD Heart Health use case with id empty

subdata101=\
[
    {
        "id": " ",
        "type": "HeartSensor",
        "age": {
            "type": "Property",
            "value": 59.5
        },
        "sex": {
            "type": "Property",
            "value": 1
        },
        "cp": {
            "type": "Property",
            "value": 3
        },
        "thal": {
            "type": "Property",
            "value": 1
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to check NGSI-LD Heart Health use case with multiple values in type of any attribute

subdata102=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor27",
        "type": "HeartSensor",
        "age": {
            "type": "Property,Property1",
            "value": 59.5
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to check Heart Health use case with nested values in value attribute

subdata103=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor28",
        "type": "HeartSensor",
        "age": {
            "type": "Property",
            "value": {
                "value": 54.0
            }
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to check Heart Health use case with multiple values in array format

subdata104=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor29",
        "type": "HeartSensor",
        "age": {
            "type": "Property",
            "value": [55.0,45.0]
        },
        "sex": {
            "type": "Property",
            "value": 1
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Heart Health Prediction use case with empty attribute

subdata105=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor30",
        "type": "HeartSensor",
        "age": {}
    }
]

#Payload to test NGSI-LD Heart Health Prediction use case if wrong information given in type

subdata106=\
[
    {
        "id": "urn:ngsi-ld:Device.HeartSensor31",
        "type": "Example",
        "age": {
            "type": 123 ,
            "value": 59.5
        },
    }
]

#Payload to test NGSI-LD Crop Prediction use case

subdata107=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor40",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property",
            "value": 45
        },
        "airTemp": {
            "type": "Property",
            "value": 20
        },
        "soilmoisture": {
            "type": "Property",
            "value": 23
        },
        "soilpH": {
            "type": "Property",
            "value": 9
        },
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                 "type": "Point",
                 "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with missing attributes
subdata108=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor41",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property",
            "value": 45
        },
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                          35.7,
                          138
                ]
            }
        }
    }
]

#Payload to update NGSI-LD Crop Prediction use case with with wrong payload
subdata109=\
[
    {
        "rainfall": {
            "type": "Property",
            "value": 70
        },
    }
]

#Payload to test NGSI-LD Crop Prediction use case with mismatch of type of values
subdata110=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor42",
        "type": "SoilSensor",
        "airmoisture": {
            "type": 560,
            "value": 45
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with type attribute is empty
subdata111=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor43",
        "type": " ",
        "airmoisture": {
            "type": "Property",
            "value": 45
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with id attribute is empty
subdata112=\
[
    {
        "id": " ",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property",
            "value": 45
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with multiple values in type of any attribute
subdata113=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor44",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property,Property1",
            "value": 45
        },
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with nested values in value attribute
subdata114=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor45",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property,Property1",
            "value": {
                "value": 54.0
                }
        },
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with value taking attribute has empty input

subdata115=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor46",
        "type": "SoilSensor",
        "airmoisture": {},
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]

#Payload to test NGSI-LD Crop Prediction use case with value taking multiple values in array format
subdata116=\
[
    {
        "id": "urn:ngsi-ld:Device.SoilSensor47",
        "type": "SoilSensor",
        "airmoisture": {
            "type": "Property,Property1",
            "value": [55.0,45.0]
        },
        "rainfall": {
            "type": "Property",
            "value": 70
        },
        "location": {
            "type": "GeoProperty",
            "value": {
                "type": "Point",
                "coordinates": [
                    35.7,
                    138
                ]
            }
        }
    }
]
#Payload to test LD visualisation feature
subdata117=\
[
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
                    "geometries": "null"
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

#Payload to delete an entity : LD visualisation feature
subdata118=\
[
        {
            "@context": "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
            "brandName": {
                "type": "Property",
                "value": "Mercedes"
            },
            "id": "urn:ngsi-ld:Vehicle:A848",
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
                    "geometries": "null"
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



#Payload(Basic GeoLocation Payload)

subdata119=\
{
               "id": "urn:ngsi-ld:Subscription:200",
               "type": "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1.02, 28.003 ],
                                "georel": "near;maxDistance==2000",
                                "geoproperty": "loc"
                }
}

#(Payload with minDistance and maxDistance)

subdata120=\
{
               "id": "urn:ngsi-ld:Subscription:A4581",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1.02, 28.003 ],
                                "georel": "near;minDistance==0.12;maxDistance==2000",
                                "geoproperty": "loc"
                }
}


#Payload with different values

subdata121=\
{
               "id": "urn:ngsi-ld:Subscription:A4582",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1.00, 28.000 ],
                                "georel": "near;minDistance==0.12;maxDistance==2000",
                                "geoproperty": "loc"
                }
}


#Payload with large float values

subdata122=\
{
               "id": "urn:ngsi-ld:Subscription:A4583",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1.00000000000, 28.00000000000 ],
                                "georel": "near;minDistance=0;maxDistance==2000",
                                "geoproperty": "loc"
                }
}



#Payload with integral values

subdata123=\
{
               "id": "urn:ngsi-ld:Subscription:A4584",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1200, 1999 ],
                                "georel": "near;minDistance=0;maxDistance==2000",
                                "geoproperty": "loc"
                }
}



#Payload with large integral values

subdata124=\
{
               "id": "urn:ngsi-ld:Subscription:A4585",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1000000000, 280000000 ],
                                "georel": "near;minDistance=0;maxDistance==20000000000000",
                                "geoproperty": "loc"
                }
}

#Payload without near in georel attribute

subdata125=\
{
               "id": "urn:ngsi-ld:Vehicle:A4586",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1000000000, 280000000 ],
                                "georel": "minDistance=0;maxDistance==20000000000000",
                                "geoproperty": "loc"
                }

}


#Payload with float in coordiantes and integral values in minDistance and maxDistance

subdata126=\
{
               "id": "urn:ngsi-ld:Vehicle:A4587",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1000000000.000000000, 280000000.00000000 ],
                                "georel": "near;minDistance=0;maxDistance==20000000000000",
                                "geoproperty": "loc"
                }
}


#Payload with integral coordiantes and float values in minDistance and maxDistance

subdata127=\
{
               "id": "urn:ngsi-ld:Subscription:A4588",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1000000000, 280000000 ],
                                "georel": "near;minDistance=0.000000000;maxDistance==20000000000000.0000000",
                                "geoproperty": "loc"
                }
}


#Payload with x and y with different type

subdata128=\
{
               "id": "urn:ngsi-ld:Subscription:A4589",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 100.223, 100 ],
                                "georel": "near;minDistance=0;maxDistance==2000",
                                "geoproperty": "loc"
                }
}



#Payload to check for a polygon

subdata129=\
{
               "id": "urn:ngsi-ld:Subscription:A4590",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "near;minDistance=0;maxDistance==20000",
                                "geoproperty": "loc"
                }
}


#Payload to check if geoproperty is left empty

subdata130=\
{
               "id": "urn:ngsi-ld:Vehicle:A4591",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Point",
                                "coordinates": [ 1000000000.0000000, 280000000.00000000 ],
                                "georel": "near",
                                "geoproperty": " "
                }
}



#Payload to check if georel is within

subdata131=\
{
               "id": "urn:ngsi-ld:Vehicle:A4592",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "within",
                                "geoproperty": "loc"
                }
}



#Payload to check if georel is contains

subdata132=\
{
               "id": "urn:ngsi-ld:Vehicle:A4593",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "contains",
                                "geoproperty": "loc"
                }
}


#Payload to check if georel is overlaps

subdata133=\
{
               "id": "urn:ngsi-ld:Vehicle:A4594",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "overlaps",
                                "geoproperty": "loc"
                }
}


#Payload to check if georel is intersects

subdata134=\
{
               "id": "urn:ngsi-ld:Vehicle:A4595",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "intersects",
                                "geoproperty": "loc"
                }
}



#Payload to check if georel is disjoint

subdata135=\
{
               "id": "urn:ngsi-ld:Vehicle:A4596",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "disjoint",
                                "geoproperty": "loc"
                }
}


#Payload to check if georel is equals

subdata136=\
{
               "id": "urn:ngsi-ld:Vehicle:A4597",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "equals",
                                "geoproperty": "loc"
                }
}


#Payload to check if georel is empty

subdata137=\
{
               "id": "urn:ngsi-ld:Vehicle:A4598",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": " ",
                                "geoproperty": "loc"
                }
}


#Payload to check if geometry is empty

subdata138=\
{
               "id": "urn:ngsi-ld:Vehicle:A4600",
               "type" : "Subscription",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": " ",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "within",
                                "geoproperty": "loc"
                }
}


#Payload to check if type empty

subdata139=\
{
               "id": "urn:ngsi-ld:Vehicle:A4601",
               "type" : " ",
               "entities": [
                               {
                                               "type": "Vehicle"

                                }
               ] ,
               "geoQ": {
                                "geometry": "Polygon",
                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]],
                                "georel": "equals",
                                "geoproperty": "loc"
                }
}

#Payload to test Geolocation Feature

subdata140=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [1.01, 4.003]
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to test Geolocation Feature with type null


subdata141=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": " ",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [1.01, 4.003]
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to test Geolocation Feature if in type property value is large float value


subdata142=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": 121.000000000000009
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [1.01, 4.003]
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to test Geolocation Feature if in type property value is integer

subdata143=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": 121
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [1.01, 4.003]
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to test Geolocation Feature if in type geoproperty the coordinates are string

subdata144=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": 121
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [1.01, 4.003]
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to test Geolocation Feature to check if datapart is empty

subdata145=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle"
}]






#Payload to test Geolocation Feature to check if coordinates are left empty


subdata146=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": []
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Payload to check if id is empty in id


subdata147=\
[{
                                "id": "urn:ngsi-ld:value:",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": []
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]




#Paylaod to check if Report:id is removed in id


subdata148=\
[{
                                "id": "urn:ngsi-ld:",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": []
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]



#Payload to test Geolocation Feature to check if id is empty


subdata149=\
[{
                                "id": " ",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "deep shikhar"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": []
                                                }
                                },
                                                                "hotel": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": ["Deep", "Shikhar"]
                                                }
                                }
                }]


#Payload to test Geolocation Feature if value is null imside type property

subdata150=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": []
                                                }
                                }
                }]



#Payload to test Geolocation Feature if coordinates are given string value


subdata151=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": "hello"
                                                }
                                }
                }]



#Payload to test Geolocation Feature if coordinates are given string as empty string value


subdata152=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": " "
                                                }
                                }
                }]



#Payload to test Geolocation Feature if entity is given empty value

subdata153=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "121"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": "null"
                }]



#Payload to test Geolocation Feature if entity type is given null

subdata154=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": " "
                                                }
                                }
                }]



#Payload to test Geolocation Feature if type Polygon has 1-D coordinates

subdata155=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [1.3,2.4]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": " "
                                                }
                                }
                }]




#Payload to test Geolocation Feature if type Point has 2-D coordinates


subdata156=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [1.3,2.4]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                }
                }]



#Payload to check if type given any arbitrary value

subdata157=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "XYZ",
                                                "value": "null"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [1.3,2.4]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                }
                }]





#Payload to check if value is given boolean

subdata158=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "False"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [1.3,2.4]
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                }
                }]



#Payload to check for nested values

subdata159=\
[{
                                "id": "urn:ngsi-ld:Vehicle:id2",
                                "type": "Vehicle",
                                "BrandName": {
                                                "type": "Property",
                                                "value": "False"
                                },
                                "location": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                                                             "type": "GeoProperty",
                                                             "value": {
                                                                "type": "Polygon",
                                                                "coordinates": [1.3,2.4]
                                                }
                                                }
                                },
                                "room": {
                                                "type": "GeoProperty",
                                                "value": {
                                                                "type": "Point",
                                                                "coordinates": [[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]
                                                }
                                }
                }]


