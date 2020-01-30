ngsi_data=\
 {
        "originator": "",
        "subscriptionId": "d0c08d50-6296-4ef3-9b0f-ff48b3cd5528",
        "contextResponses": [{
                "contextElement": {
                        "attributes": [{
                                "value": 34,
                                "type": "float",
                                "name": "tempo"
                        }],
                        "entityId": {
                                "type": "Tempr",
                                "id": "Temprature702"
                        },
                        "domainMetadata": [{
                                "type": "point",
                                "name": "location",
                                "value": {
                                        "latitude": 49.406393,
                                        "longitude": 8.684208
                                }
                        }]
                },
                "statusCode": {
                        "code": 200,
                        "reasonPhrase": "OK"
                }
        }]
}

# Expected output for withDomainMetadata
convert_data_output=\
{
        "@context": [ {
                "Tempr": "http://example.org/Tempr",
                "tempo": "http://example.org/tempo"
        }],
        "location": {
                "type": "GeoProperty",
                "value": "{\"type\": \"point\", \"coordinates\": [49.406393, 8.684208]}"
        },
        "tempo": {
                "type": "Property",
                "value": 34
        },
        "id": "urn:ngsi-ld:Temprature702",
        "type": "Tempr"
}

patch_data_output=\
{
        "@context": ["https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld", {
                "Tempr": "http://example.org/Tempr",
                "tempo": "http://example.org/tempo"
        }],
        "tempo": {
                "type": "Property",
                "value": 34
        },
        "location": {
                "type": "GeoProperty",
                "value": "{\"type\": \"Point\", \"coordinates\": [49.406393, 8.684208]}"
        }
}

orian_notify_data=\
{
        "subscriptionId": "5d09cdeb0016e878cf94e070",
        "data": [{
                "type": "roomie",
                "id": "Room5",
                "temp": {
                        "type": "Float",
                        "value": 50,
                        "metadata": {}
                }
        }]
}

orian_notify_output_data=\
{
        "@context": ["https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld", {
                "roomie": "http://example.org/roomie",
                "temp": "http://example.org/temp"
        }],
        "type": "roomie",
        "id": "urn:ngsi-ld:Room5",
        "temp": {
                "type": "Float",
                "value": 50
        }
}
id_value="urn:ngsi-ld:Temprature702"

#Input for ld_generate without DomainMetadata

wdmdata=\
 {
        "originator": "",
        "sbscriptionId": "2cdc4d34-0216-459a-b582-1c13a6d1bec1",
        "contextResponses": [{
                "contextElement": {
                        "attributes": [{
                                "type": "string",
                                "name": "livestockType",
                                "value": "Catte321"
                        }, {
                                "type": "string",
                                "name": "breed",
                                "value": "Japan"
                        }, {
                                "type": "string",
                                "name": "gender",
                                "value": "Male"
                        }],
                        "entityId": {
                                "type": "Livestock",
                                "id": "Livestock006"
                        }
                },
                "statsCode": {
                        "code": 200,
                        "reasonPhrase": "OK"
                }
        }]
}

# Expected output for without DomainMetadata

Eodata=\
{
        'gender': {
                'type': 'Property',
                'value': 'Male'
        },
        'breed': {
                'type': 'Property',
                'value': 'Japan'
        },
        'livestockType': {
                'type': 'Property',
                'value': 'Catte321'
        },
        '@context': [{
                'livestockType': 'http://example.org/livestockType',
                'breed': 'http://example.org/breed',
                'Livestock': 'http://example.org/Livestock',
                'gender': 'http://example.org/gender'
        }],
        'type': 'Livestock',
        'id': 'urn:ngsi-ld:Livestock006'
}

# point metadata

point_metaData=\
{
'type': 'point', 'name': 'location', 'value': {'latitude': 49.406393, 'longitude': 8.684208}
}

ngb_pointMetaDataOutput=\
{
      'type': 'GeoProperty', 'value': '{"type": "point", "coordinates": [49.406393, 8.684208]}'
}

#input to test metadata

get_coordinate_input=\
{'type': 'point', 'name': 'location', 'value': {'latitude': 49.406393, 'longitude': 8.684208}}

get_coordinate_output=\
{'latitude': 49.406393, 'longitude': 8.684208}

# input data to test metadata without type
get_coordinate_input_forWithoutType=\
{'name': 'location', 'value': {'latitude': 49.406393, 'longitude': 8.684208}}

get_coordinate_output_forWithoutType=-1

# input data to test metadata without value

get_coordinate_input_forWithoutValue=\
{'type': 'point', 'name': 'location'}

get_coordinate_output_forWithoutValue=-1

get_coordinate_input_forWithoutTypeValue=\
{'name': 'location'}

get_coordinate_output_forWithoutTypeValue=-1

domainMetaData_with_polygon_input=\
{
	"type": "polygon",
	"name": "location",
	"value": {
		"vertices": [{
			"latitude": 49.406393,
			"longitude": 822.684208
		}, {
			"latitude": 429.406393,
			"longitude": 811.684208
		}, {
			"latitude": 491.406393,
			"longitude": 822.684208
		}]
	}
}

domainMetaData_with_withpolygon_output=\
{
	'type': 'GeoProperty',
	'value': '{"type": "polygon", "coordinates": [[49.406393, 822.684208], [429.406393, 811.684208], [491.406393, 822.684208]]}'
}

domainMetaData_with_polygon_withoutVertices_input=\
{
	'type': 'polygon',
	'name': 'location',
	'value': {
		'vertices': None
	}
}

domainMetaData_with_polygon_withoutVertices_output=\
{
	'type': 'GeoProperty',
	'value': '{"type": "polygon", "coordinates": []}'
}

# circle metaData

circle_metaData_input=\
{
	'type': 'polygon',
	'name': 'location',
	'value': {
		'centerLongitude': 0,
		'radius': 80.8332,
		'centerLatitude': 0
	}
}

circle_metaData_output=\
{
	'type': 'GeoProperty',
	'value': '{"type": "polygon", "coordinates": [0, 0, 80.8332]}'
}

# test data for ld_generate for polygon 

polygonDataInput=\
{
	'originator': '',
	'subscriptionId': 'f3b423fe-aa0c-484f-8d4b-4629399f0d12',
	'contextResponses': [{
		'contextElement': {
			'attributes': [{
				'type': 'string',
				'name': 'livestockType',
				'value': 'Cattle'
			}],
			'entityId': {
				'type': 'Livestock',
				'id': 'Livestock9999'
			},
			'domainMetadata': [{
				'type': 'polygon',
				'name': 'location',
				'value': {
					'vertices': [{
						'latitude': 49.406393,
						'longitude': 822.684208
					}, {
						'latitude': 429.406393,
						'longitude': 811.684208
					}, {
						'latitude': 491.406393,
						'longitude': 822.684208
					}]
				}
			}]
		},
		'statusCode': {
			'code': 200,
			'reasonPhrase': 'OK'
		}
	}]
}

polygonDataOutput=\
{
	'@context': [{
		'livestockType': 'http://example.org/livestockType',
		'Livestock': 'http://example.org/Livestock'
	}],
	'location': {
		'type': 'GeoProperty',
		'value': '{"type": "polygon", "coordinates": [[49.406393, 822.684208], [429.406393, 811.684208], [491.406393, 822.684208]]}'
	},
	'type': 'Livestock',
	'id': 'urn:ngsi-ld:Livestock9999',
	'livestockType': {
		'type': 'Property',
		'value': 'Cattle'
	}
}

#test data for ld_generate for circle
 
circleDataInput=\
{
	'originator': '',
	'subscriptionId': 'f3b423fe-aa0c-484f-8d4b-4629399f0d12',
	'contextResponses': [{
		'contextElement': {
			'attributes': [{
				'type': 'string',
				'name': 'livestockType',
				'value': 'Cattle'
			}],
			'entityId': {
				'type': 'Livestock',
				'id': 'Livestock9999'
			},
			'domainMetadata': [{
				'type': 'polygon',
				'name': 'location',
				'value': {
					'centerLongitude': 0,
					'radius': 30.9443,
					'centerLatitude': 0
				}
			}]
		},
		'statusCode': {
			'code': 200,
			'reasonPhrase': 'OK'
		}
	}]
}

circleDataOutput=\
{
	'@context': [{
		'livestockType': 'http://example.org/livestockType',
		'Livestock': 'http://example.org/Livestock'
	}],
	'location': {
		'type': 'GeoProperty',
		'value': '{"type": "polygon", "coordinates": [0, 0, 30.9443]}'
	},
	'type': 'Livestock',
	'id': 'urn:ngsi-ld:Livestock9999',
	'livestockType': {
		'type': 'Property',
		'value': 'Cattle'
	}
}

polygonTestDataWV_input=\
{
	'originator': '',
	'subscriptionId': 'f3b423fe-aa0c-484f-8d4b-4629399f0d12',
	'contextResponses': [{
		'contextElement': {
			'attributes': [{
				'type': 'string',
				'name': 'livestockType',
				'value': 'Cattle'
			}],
			'entityId': {
				'type': 'Livestock',
				'id': 'Livestock9999'
			},
			'domainMetadata': [{
				'type': 'polygon',
				'name': 'location',
				'value': {
					'vertices': None
				}
			}]
		},
		'statusCode': {
			'code': 200,
			'reasonPhrase': 'OK'
		}
	}]
}

polygonTestDataWV_output=\
{
	'@context': [{
		'livestockType': 'http://example.org/livestockType',
		'Livestock': 'http://example.org/Livestock'
	}],
	'location': {
		'type': 'GeoProperty',
		'value': '{"type": "polygon", "coordinates": []}'
	},
	'type': 'Livestock',
	'id': 'urn:ngsi-ld:Livestock9999',
	'livestockType': {
		'type': 'Property',
		'value': 'Cattle'
	}
}

# Test data for ld_generate for object value 
TestDataForObject_input=\
{
	'originator': '',
	'subscriptionId': 'aa58629f-5638-4dd9-ac95-1253d6648654',
	'contextResponses': [{
		'contextElement': {
			'attributes': [{
				'type': 'Work_results',
				'name': 'livestockType1',
				'value': {
					'name2': {
						'processName': 'process',
						'startTime': '26-08-2019'
					},
					'name1': {
						'processName': 'process1',
						'startTime': '26-08-2019'
					}
				}
			}],
			'entityId': {
				'type': 'Work_results',
				'id': 'Livestock2000'
			}
		},
		'statusCode': {
			'code': 200,
			'reasonPhrase': 'OK'
		}
	}]
}

TestDataForObject_output=\
{
	'@context': [{
		 'Work_results':  'http://example.org/Work_results',
		 'livestockType1':  'http://example.org/livestockType1'
	}],
	 'livestockType1': {
		'type': 'Property',
		'value': {
			 'name2': {
				 'processName':  'process',
				 'startTime':  '26-08-2019'
			},
			 'name1': {
				 'processName':  'process1',
				 'startTime':  '26-08-2019'
			}
		}
	},
	'type':  'Work_results',
	'id':  'urn:ngsi-ld:Livestock2000'
}
