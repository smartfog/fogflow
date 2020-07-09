*****************************************
NGSI-LD API Walkthrough
*****************************************

This tutorial is focused mainly on the NGSI-LD APIs supported in FogFlow, which include APIs for entities, context registrations and subscriptions. These are discussed in more detail in the following sections. For using NGSI-LD APIs in FogFlow, checkout the latest docker image "fogflow/broker:3.1" from Docker Hub.

FogFlow follows the NGSI-LD data-model, with continuous improvements. For better understanding of NGSI-LD Data-model, refer `this`_.

.. _`this`: https://fiware-datamodels.readthedocs.io/en/latest/ngsi-ld_howto/index.html


Entities
=========================

Entities are the units for representing objects in the environment, each having some properties of their own, may also have some relationships with others. This is how linked data are formed.


Create entities
------------------------------------------

There are several ways of creating an NGSI-LD Entity on FogFlow Broker:

* When context is provided in the Link header: The context for resolving the payload is given through the Link header.
* When context is provided in the payload: Context is in the payload itself, there is no need to attach a Link header in the request.
* When the request payload is already expanded: Some payloads are already expanded using some context.

Curl requests for creating an entity on FogFlow Broker in different ways are given below. All the are the POST requests to FogFlow Broker. Broker returns a response of "201 Created" for a successful creation of a new entity and "409 Conflict" on creating an already existing entity.

**When context is provided in the Link header:**

.. code-block:: console

        curl -iX POST \
        'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
        -H 'Content-Type: application/ld+json' \
        -H 'Accept: application/ld+json' \
        -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
        -d '
        {
                "id": "urn:ngsi-ld:Vehicle:A100",
                "type": "Vehicle",
                "brandName": {
                        "type": "Property",
                        "value": "BMW",
                        "observedAt": "2017-07-29T12:00:04"
                },
                "isParked": {
                        "type": "Relationship",
                        "object": "urn:ngsi-ld:OffStreetParking:Downtown",
                        "observedAt": "2017-07-29T12:00:04",
                        "providedBy": {
                                "type": "Relationship",
                                "object": "urn:ngsi-ld:Person:Bob"
                        }
                },
                "speed": {
                        "type": "Property",
                        "value": 81,
                        "observedAt": "2017-07-29T12:00:04"
                },
                "location": {
                        "type": "GeoProperty",
                        "value": {
                                "type": "Point",
                                "coordinates": [-8.5, 41.2]
                        }
                }
        }'

**When context is provided in the payload:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"@context": [{
			"Vehicle": "https://uri.etsi.org/ngsi-ld/default-context/Vehicle",
			"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName",
			"speed": "https://uri.etsi.org/ngsi-ld/default-context/speed",
			"isParked": {
				"@type": "@id",
				"@id": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
			}
		}],
		"id": "urn:ngsi-ld:Vehicle:A200",
		"type": "Vehicle1",
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
	}'

**When the request payload is already expanded:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"https://uri.etsi.org/ngsi-ld/default-context/brandName": [
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
		"@id": "urn:ngsi-ld:Vehicle:A300",
		"https://uri.etsi.org/ngsi-ld/default-context/isParked": [
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
				"https://uri.etsi.org/ngsi-ld/default-context/providedBy": [
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
		"https://uri.etsi.org/ngsi-ld/default-context/speed": [
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
			"https://uri.etsi.org/ngsi-ld/default-context/Vehicle"
		]
	}'


Update entities
-----------------------------------------------

Entities can be updated by updating their attributes (properties and relationships) and the attributes can be updated in the following ways:

* Add more attributes to the entity: More properties or relationships or both can be added to an existing entity. This is a POST http request to Broker to append more attributes to the entity.
* Update existing attributes of the entity: Existing properties or relationships or both can be updated for an entity. This is a PATCH http request to FogFlow Broker.
* Update specific attribute of the entity: Fields of an existing attribute can be updated for an entity. This update is also called partial update. This is also a PATCH request to the FogFlow Broker.

FogFlow Broker returns "204 NoContent" on a successful attribute update, "404 NotFound" for a non-existing entity. While updating the attributes of an exiting entity, some of the attributes provided in the request payload may not exist. For such cases, FogFlow Broker return a "207 MultiStatus" error.

Here are the curl requests for these Updates.

**Add more attributes to the entity:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs' \
	-H 'Content-Type: application/ld+json' \
	-d '
	{
		"@context": {
			"brandName1": "https://uri.etsi.org/ngsi-ld/default-context/brandName1",
			"isParked1": "https://uri.etsi.org/ngsi-ld/default-context/isParked1"
		},
		"brandName1": {
			"type": "Property",
			"value": "Audi"
		},
		
		"isParked1": {
			"type": "Relationship",
			"object": "Audi"
		}
	}'

**Update existing attributes of the entity:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs' \
	-H 'Content-Type: application/ld+json' \
	-d '
	{
		"@context": {
			"isParked": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
		},
		"brandName": {
			"type": "Property",
			"object": "Audi"
		}
	}'

**Update specific attribute of the entity:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs/<Attribute-Name>' \
	-H 'Content-Type: application/ld+json' \
	-d '
		{
		"@context": {
			"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName"
		},
		"value": "Suzuki"
	}'


Get entities
-----------------------------------------------

This section describes how to retrieve the already created entities from FogFlow Broker. Entities can be retrieved from FogFlow based on different filters, listed below.

* Based on Entity Id: returns an entity whose id is passed in the request URL.
* Based on Attribute Name: returns all those entities which contain the attribute name that is passed in the query parameters of the request URL.
* Based on Entity Id and Entity Type: returns the entity with the entity id same as given in the query parameters along with the type matching.
* Based on Entity Type: returns all the entities that are of the requested type.
* Based on Entity Type with Link header: returns all the entities of requested type, but here the type can be given in a different way in the query parameters of request URL. Refer the request for this in the following sections.
* Based on Entity IdPattern and Entity Type: returns all those entities which lie inside the IdPattern range and the matching type mentioned in the query parameters.

On successful retrieval of at least one entity in the above requests, FogFlow Broker returns a "200 OK" response. For non-existing entities, "404 NotFound" error is returned.

**Based on Entity Id:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Attribute Name:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?attrs=<Expanded-Attribute-Name>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Entity Id and Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?id=<Entity-Id>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Entity Type with Link header:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?type=<Unexpanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'

**Based on Entity IdPattern and Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?idPattern=<Entity-IdPattern>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'


Delete entities
-----------------------------------------------

Either an entity can be deleted, or a specific attribute of that entity can be deleted. Successful deletion returns a "204 NoContent" response, while for non-existing attributes or entities, it returns "404 NotFound" error. 

**Deleting specific attribute of an entity:**

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs/<Attribute-Name>'

**Deleting an entity:**

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>'


Registrations
================================

Registrations or C-Source Registrations are used to indicate which device will be feeding what data to a Broker. Data description like the entity ids and their types, their properties and relationships, endpoint of the provider, location of the data, etc. are given in a C-Source request.


Create registrations
--------------------------------------

A C-Source Registration can be created on FogFlow Broker in the following two ways:

* Using context in Link header: context for resolving the request payload is contained in the Link header of the request.
* Using context in payload: context is given in the payload itself.

On creating a C-Source registration, FogFlow Broker returns "201 Created" response, while in case of at least one already registered entity in the request payload, it will return a "409 Conflict" error.

Curl requests are given in the following sections.

**Using context in Link header:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"id": "urn:ngsi-ld:ContextSourceRegistration:csr1a3400",
		"type": "ContextSourceRegistration",
		"name": "NameExample",
		"description": "DescriptionExample",
		"information": [
			{
				"entities": [
					{
						"id": "urn:ngsi-ld:Vehicle:A500",
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
						"id": "urn:ngsi-ld:Vehicle:A600",
						"type": "OffStreetParking"
					}
				]
			}
		],
		"endpoint": "http://my.csource.org:1026",
		"location": "{ \"type\": \"Polygon\", \"coordinates\": [[[8.686752319335938,49.359122687528746],[8.742027282714844,49.3642654834877],[8.767433166503904,49.398462568451485],[8.768119812011719,49.42750021620163],[8.74305725097656,49.44781634951542],[8.669242858886719,49.43754770762113],[8.63525390625,49.41968407776289],[8.637657165527344,49.3995797187007],[8.663749694824219,49.36851347448498],[8.686752319335938,49.359122687528746]]] }",
		"timestamp": {
			"start": "2017-11-29T14:53:15",
			"end": "2017-12-29T14:53:15"
		},
		"expires": "2030-11-29T14:53:15"
	}'

**Using context in payload:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations/<Registration-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"id": "urn:ngsi-ld:ContextSourceRegistration:csr1a3401",
		"type": "ContextSourceRegistration",
		"name": "NameExample",
		"description": "DescriptionExample",
		"information": [
			{
				"entities": [
					{
						"id": "urn:ngsi-ld:Vehicle:A700",
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
				  "id": "urn:ngsi-ld:Vehicle:A800",
				  "type": "OffStreetParking"
				}
			  ]
			}
		],
		"endpoint": "http://my.csource.org:1026",
		"location": "{ \"type\": \"Polygon\", \"coordinates\": [[[8.686752319335938,49.359122687528746],[8.742027282714844,49.3642654834877],[8.767433166503904,49.398462568451485],[8.768119812011719,49.42750021620163],[8.74305725097656,49.44781634951542],[8.669242858886719,49.43754770762113],[8.63525390625,49.41968407776289],[8.637657165527344,49.3995797187007],[8.663749694824219,49.36851347448498],[8.686752319335938,49.359122687528746]]] }",
		"timestamp": {
			"start": "2017-11-29T14:53:15"
		},
		"expires": "2030-11-29T14:53:15",
		"@context": [
			"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",    
			{
				"Vehicle": "https://uri.etsi.org/ngsi-ld/default-context/Vehicle",
				"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName",
				"brandName1": "https://uri.etsi.org/ngsi-ld/default-context/brandName1",
				"speed": "https://uri.etsi.org/ngsi-ld/default-context/speed",
				"totalSpotNumber": "https://uri.etsi.org/ngsi-ld/default-context/parking/totalSpotNumber",
				"reliability": "https://uri.etsi.org/ngsi-ld/default-context/reliability",
				"OffStreetParking":    "https://uri.etsi.org/ngsi-ld/default-context/parking/OffStreetParking",    
				"availableSpotNumber":    "https://uri.etsi.org/ngsi-ld/default-context/parking/availableSpotNumber",
				 "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
				"isParked": {
					"@type": "@id",
					"@id": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
				},
				"isNextToBuilding":    {    
					"@type":    "@id",    
					"@id":    "https://uri.etsi.org/ngsi-ld/default-context/isNextToBuilding"    
				},    
				"providedBy":    {    
					"@type":    "@id",    
					"@id":    "https://uri.etsi.org/ngsi-ld/default-context/providedBy"    
				},    
				"name":    "https://uri.etsi.org/ngsi-ld/default-context/name"    
			}
		]
	}'


Update registrations
--------------------------------------

An existing C-Source Registration can be updated by its id. Context for resolving the payload is given in the request payload. In case the context object is not given in the request payload, FogFlow Broker will resolve the payload using the default context. "204 NoContent" response is returned on a successful registration update on FogFlow Broker.

Curl request for C-Source Registration update is given below.

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations/<Registration-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"type": "ContextSourceRegistration",
		"name": "NameExample",
		"description": "DescriptionExample",
		"information": [
			{
				"entities": [
					{
						"id": "urn:ngsi-ld:Vehicle:A500",
						"type": "Vehicle"
					}
				],
				"properties": [
					"brandName",
					"speed",
					"brandName1"
				],
				"relationships": [
					"isParked",
					"isParked1"
				]
			},
			{
				"entities": [
					{
						"id": "urn:ngsi-ld:Vehicle:A600",
						"type": "Vehicle"
					}
				],
				"properties": [
					"brandName"
				],
				"relationships": [
					"isParked"
				]
			}
		],
		"endpoint": "http://my.csource.org:1026",
		"location": "{ \"type\": \"Polygon\", \"coordinates\": [[[8.686752319335938,49.359122687528746],[8.742027282714844,49.3642654834877],[8.767433166503904,49.398462568451485],[8.768119812011719,49.42750021620163],[8.74305725097656,49.44781634951542],[8.669242858886719,49.43754770762113],[8.63525390625,49.41968407776289],[8.637657165527344,49.3995797187007],[8.663749694824219,49.36851347448498],[8.686752319335938,49.359122687528746]]] }",
		"timestamp": {
			"start": "2017-11-29T14:53:15"
		},
		"expires": "2030-11-29T14:53:15",
		"@context": [
            "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",    
			{
				"Vehicle": "https://uri.etsi.org/ngsi-ld/default-context/Vehicle",
				"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName",
				"brandName1": "https://uri.etsi.org/ngsi-ld/default-context/brandName1",
				"speed": "https://uri.etsi.org/ngsi-ld/default-context/speed",
				"totalSpotNumber": "https://uri.etsi.org/ngsi-ld/default-context/parking/totalSpotNumber",
				"reliability": "https://uri.etsi.org/ngsi-ld/default-context/reliability",
				"OffStreetParking":    "https://uri.etsi.org/ngsi-ld/default-context/parking/OffStreetParking",    
				"availableSpotNumber":    "https://uri.etsi.org/ngsi-ld/default-context/parking/availableSpotNumber",
				 "timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
				"isParked": {
					"@type": "@id",
					"@id": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
				},
				"isNextToBuilding":    {    
					"@type":    "@id",    
					"@id":    "https://uri.etsi.org/ngsi-ld/default-context/isNextToBuilding"    
				},    
				"providedBy":    {    
					"@type":    "@id",    
					"@id":    "https://uri.etsi.org/ngsi-ld/default-context/providedBy"    
				},    
				"name":    "https://uri.etsi.org/ngsi-ld/default-context/name",
				"timestamp": "http://uri.etsi.org/ngsi-ld/timestamp",
				"expires":"http://uri.etsi.org/ngsi-ld/expires"
			}
		]
	}'


Get registrations
--------------------------------------

C-Source Registrations can be retrieved from FogFlow Broker using the following filters, which are passed in the request through query parameters.

* Based on Entity Type: returns all the registrations with the matching entity type.
* Based on Entity Type with Link header: returns all the registrations with matching entity type, but here, entity type is passed differently.
* Based on Entity Id and Entity Type: returns the registration which contains the requested entity id and type.
* Based on Entity IdPattern and Entity Type: returns all those registrations which lie within the range of requested entity id pattern and also matching the entity type.

Successful retrieval returns "200 OK" response while in case on not-existing registrations, Broker returns "404 NotFound" error. Send the following curl requests to Broker to view how it works.

**Based on Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations?type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Entity Type with Link header:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations?type=<Unexpanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'

**Based on Entity Id and Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations?id=<Entity-Id>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Based on Entity IdPattern and Entity Type:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations?idPattern=<Entity-IdPattern>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'


Delete registrations
--------------------------------------

C-Source registration can be deleted using the following request.

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/csourceRegistrations/<Registration-Id>'


Subscriptions
================================

Subscribers can subscribe for entities using a subscription request to the FogFlow Broker.


Create subscriptions
--------------------------------------

Subscriptions can be created either for an Entity Id or an Entity Id Pattern. Whenever entity update is there for that subscription, FogFlow Broker will automatically notify the updated entity to the subscribers. "201 Created" response is returned on a successful subscription on Broker, along with the Subscription Id, which can later be used to retrieve, update or delete the subscription.

Refer the following curl requests, but before running the subscriptions, make sure some notify receiver is running, that can simply view the contents of the notification. For already subscribed entities, when entity creation or update takes place, a notification will be received by the subscriber. Notification is also received by a subscriber in case of subscription to an already existing entity.

**Subscribing for an Entity Id**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"id" : "urn:ngsi-ld:Vehicle:A100",
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
	}'

**Subscribing for an IdPattern:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"idPattern" : ".*",
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
	}'


Update subscriptions
--------------------------------------

An existing subscription on FogFlow Broker can be updated by id using the curl request given below.

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"type": "Vehicle1"
		}],
		"watchedAttributes": ["https://uri.etsi.org/ngsi-ld/default-context/brandName11"],
		"notification": {
			"attributes": ["https://uri.etsi.org/ngsi-ld/default-context/brandName223"],
			"format": "keyValues",
			"endpoint": {
				"uri": "http://my.endpoint.org/notify",		
				"accept": "application/json"
			}
		}
	}'
	

Get subscriptions
--------------------------------------

All the subscriptions or a subscription with specific id, both can be retrieved from FogFlow Broker with a response of "200 OK". Curl requests are given below.

**All Subscriptions:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Accept: application/ld+json'

**Specific subscription:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>' \
	-H 'Accept: application/ld+json'


Delete subscriptions
--------------------------------------

A subscription can be deleted by sending the following request to FogFlow Broker, with a response of 204 "NoContent".

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>'



**The NGSI-LD support in FogFlow also carries some limitations with it. Improvements are continued.**
