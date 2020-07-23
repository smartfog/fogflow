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

subdata4=\
{
      "@context": {
                        "brandName1": "http://example.org/vehicle/brandName1"
       },
      "brandName1": {
                          "type": "Property",
                          "value": "BMW"
       }
 }

subdata5=\
{
        "@context": {
                        "brandName1": "http://example.org/vehicle/brandName1"
         },
        "brandName1": {
                           "type": "Property",
                           "value": "AUDI"
         }
  }

subdata6=\
{
        "@context": {
                         "brandName": "http://example.org/vehicle/brandName"
         },
         "value": "MARUTI"
  }

subdata7=\
{
             "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a3459",
             "type": "ContextSourceRegistration",
             "name": "NameExample",
             "description": "DescriptionExample",
             "information": [
             {
                     "entities": [
                      {
                             "id": "urn:ngsi-ld:Vehicle:A456",
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
                             "id": "downtown",
                             "type": "OffStreetParking"
                     }
                   ]
             }
            ],
             "endpoint": "http://http://my.csource.org:1026",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

subdata8=\
{
         "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a4000",
         "type": "ContextSourceRegistration",
         "name": "NameExample",
         "description": "DescriptionExample",
         "information": [
         {
                 "entities": [
                  {
                         "id": "urn:ngsi-ld:Vehicle:A555",
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
                         "id": "town$",
                         "type": "OffStreetParking"
                  }
                 ]
         }
         ],
         "endpoint": "http://my.csource.org:1026",
         "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
         "timestamp": {
         "start": "2017-11-29T14:53:15"
         },
         "expires": "2030-11-29T14:53:15",
         "@context": [

         "https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld",
         {
                 "Vehicle": "http://example.org/vehicle/Vehicle",
                 "brandName": "http://example.org/vehicle/brandName",
                 "brandName1": "http://example.org/vehicle/brandName1",
                 "speed": "http://example.org/vehicle/speed",
                 "totalSpotNumber": "http://example.org/parking/totalSpotNumber",
                 "reliability": "http://example.org/common/reliability",
                 "OffStreetParking":  "http://example.org/parking/OffStreetParking",
                 "availableSpotNumber":  "http://example.org/parking/availableSpotNumber",
                 "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
                 "isParked": {
                                 "@type": "@id",
                                 "@id": "http://example.org/common/isParked"
                 },
                 "isNextToBuilding":    {
                                 "@type":  "@id",
                                 "@id":  "http://example.org/common/isNextToBuilding"
                  },
                 "providedBy":    {
                                 "@type":  "@id",
                                 "@id":  "http://example.org/common/providedBy"
                 },
                  "name":    "http://example.org/common/name"
         }
       ]
 }

subdata9=\
{
         "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a3459",
         "type": "ContextSourceRegistration",
         "name": "NameExample",
         "description": "DescriptionExample",
         "information": [
         {
                 "entities": [
                 {
                         "id": "urn:ngsi-ld:Vehicle:A456",
                         "type": "Vehicle"
                 }
               ],
                 "properties": [
                                 "brandName",
                                 "speed",
                                 "brandName1"
               ],
                 "relationships": [
                                 "isParked"
               ]
         }
       ],
         "endpoint": "http://my.csource.org:1026",
         "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
         "timestamp": {
         "start": "2017-11-29T14:53:15"
         },
         "expires": "2030-11-29T14:53:15",
         "@context": [

                         "https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld",
                         {
                                 "Vehicle": "http://example.org/vehicle/Vehicle",
                                 "brandName": "http://example.org/vehicle/brandName",
                                 "brandName1": "http://example.org/vehicle/brandName1",
                                 "speed": "http://example.org/vehicle/speed",
                                 "totalSpotNumber": "http://example.org/parking/totalSpotNumber",
                                 "reliability": "http://example.org/common/reliability",
                                 "OffStreetParking":  "http://example.org/parking/OffStreetParking",
                                 "availableSpotNumber":  "http://example.org/parking/availableSpotNumber",
                                 "isParked": {
                                                 "@type": "@id",
                                                 "@id": "http://example.org/common/isParked"
                                  },
                                 "isNextToBuilding":  {
                                                 "@type": "@id",
                                                 "@id": "http://example.org/common/isNextToBuilding"
                                  },
                                 "providedBy": {
                                                 "@type":  "@id",
                                                 "@id":  "http://example.org/common/providedBy"
                                  },
                          "name": "http://example.org/common/name",
                          "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
                          "expires":"http://uri.etsi.org/ngsi-ld/expires"
                         }
         ]
  }


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
                                             "uri": "http://my.endpoint.org/notify",
                                             "accept": "application/json"
                               }
             }
      }

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
                             "uri": "http://my.endpoint.org/notify",
                             "accept": "application/json"
              }
           }
       }

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


subdata14=\
{
      "@context": {
                        "brandName1": "http://example.org/vehicle/brandName1"
       },
      "brandName1": {
                          "type": "Property",
                          "value": "MARUTI"
       }
 }

subdata15=\
{
        "@context": {
                        "brandName1": "http://example.org/vehicle/brandName1"
         },
        "brandName1": {
                          
         }
  }

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


subdata17=\
{
           "id": "urn:ngsi-ld:Vehicle:A900",
           "type": "Vehicle",
            "createdAt": "2017-07-29T12:00:04"
     }

subdata18=\
{
           "id": "urn:ngsi-ld:Vehicle:A501",
           "type": "Vehicle",
            "createdAt": "2017-07-29T12:00:04"
     }

subdata19=\
{
  "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a3459",
  "type": "ContextSourceRegistration",
  "name": "NameExample",
  "description": "DescriptionExample",
  "information": [
    {
      "entities": [
        {
          "id": "urn:ngsi-ld:Vehicle:A456",
          "type": "Vehicle"
        }
      ],
      
      "relationships": [
        "isParked"
      ]
    },
    {
      "entities": [
        {
          "idPattern": "downtown$",
          "type": "OffStreetParking"
        }
      ]
    }
  ],
  "endpoint": "http://my.csource.org:1026",
  "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
  "timestamp": {
    "start": "2017-11-29T14:53:15"
  },
  "expires": "2030-11-29T14:53:15",
"@context": [

                "https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld",    
    {
    "Vehicle": "http://example.org/vehicle/Vehicle",
    "brandName": "http://example.org/vehicle/brandName",
    "brandName1": "http://example.org/vehicle/brandName1",
    "speed": "http://example.org/vehicle/speed",
    "totalSpotNumber": "http://example.org/parking/totalSpotNumber",
    "reliability": "http://example.org/common/reliability",
    "OffStreetParking":    "http://example.org/parking/OffStreetParking",    
    "availableSpotNumber":    "http://example.org/parking/availableSpotNumber",    
    "isParked": {
        "@type": "@id",
        "@id": "http://example.org/common/isParked"
    },
    "isNextToBuilding":  {    
        "@type":  "@id",    
        "@id":  "http://example.org/common/isNextToBuilding"    
    },    
    "providedBy": {    
        "@type":  "@id",    
        "@id":  "http://example.org/common/providedBy"    
    },    
    "name": "http://example.org/common/name",
    "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
    "expires":"http://uri.etsi.org/ngsi-ld/expires"
}
]
}

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
                             "uri": "http://my.endpoint.org/notify",
                             "accept": "application/json"
              }
           }
       }

subdata21=\
{
             "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a4001",
             "type": "ContextSourceRegistration",
             "name": "NameExample",
             "description": "DescriptionExample",
             "information": [
             {
                     "entities": [
                      {
                             "id": "urn:ngsi-ld:Vehicle:A456",
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
                             "id": "uptown$",
                             "type": "OffStreetParking"
                     }
                   ]
             }
            ],
             "endpoint": "http://http://my.csource.org:1026",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

subdata22=\
{
             "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a4001",
             "type": "ContextSourceRegistration",
             "name": "NameExample",
             "description": "DescriptionExample",
             "information": [
             {
                     "entities": [
                      {
                             "id": "urn:ngsi-ld:Vehicle:A666",
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
             }
            ],
             "endpoint": "http://http://my.csource.org:1026",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

subdata23=\
{
  "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a4001",
  "type": "ContextSourceRegistration",
  "name": "NameExample",
  "description": "DescriptionExample",
  "information": [
    {
      "entities": [
        {
          "id": "urn:ngsi-ld:Vehicle:A666",
          "type": "Vehicle"
        }
      ],
      "properties": [
        "brandName",
        "speed",
        "brandName1"
      ],
      "relationships": [
        "isParked"
      ]
    }
  ],
  "endpoint": "http://my.csource.org:1026",
  "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
  "timestamp": {
    "start": "2017-11-29T14:53:15"
  },
  "expires": "2030-11-29T14:53:15",
"@context": [

                "https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld",    
    {
    "Vehicle": "http://example.org/vehicle/Vehicle",
    "brandName": "http://example.org/vehicle/brandName",
    "brandName1": "http://example.org/vehicle/brandName1",
    "speed": "http://example.org/vehicle/speed",
    "totalSpotNumber": "http://example.org/parking/totalSpotNumber",
    "reliability": "http://example.org/common/reliability",
    "OffStreetParking":    "http://example.org/parking/OffStreetParking",    
    "availableSpotNumber":    "http://example.org/parking/availableSpotNumber",    
    "isParked": {
        "@type": "@id",
        "@id": "http://example.org/common/isParked"
    },
    "isNextToBuilding":  {    
        "@type":  "@id",    
        "@id":  "http://example.org/common/isNextToBuilding"    
    },    
    "providedBy": {    
        "@type":  "@id",    
        "@id":  "http://example.org/common/providedBy"    
    },    
    "name": "http://example.org/common/name",
    "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
    "expires":"http://uri.etsi.org/ngsi-ld/expires"
}
]
}

subdata24=\
{
             "id": "urn:ngsi-ld:ContextSourceRegistration:csr1a4002",
             "type": "ContextSourceRegistration",
             "name": "NameExample",
             "description": "DescriptionExample",
             "information": [
             {
                     "entities": [
                      {
                             "id": "urn:ngsi-ld:Vehicle:A662",
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
             }
            ],
             "endpoint": "http://http://my.csource.org:1026",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

subdata25=\
{
	"type":"ContextSourceRegistration"
}

subdata26=\
{
}

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
             "endpoint": "http://http://my.csource.org:1026",
             "location": "{ \"type\": \"Point\", \"coordinates\": [-8.5, 41.2] }",
             "timestamp": {
                             "start": "2017-11-29T14:53:15"
                             },
             "expires": "2030-11-29T14:53:15"
     }

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
                                             "uri": "http://my.endpoint.org/notify",
                                             "accept": "application/json"
                               }
             }
      }

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
                                             "uri": "http://my.endpoint.org/notify",
                                             "accept": "application/json"
                               }
             }
      }


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
                             "uri": "http://my.endpoint.org/notify",
                             "accept": "application/json"
              }
           }
       }

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


