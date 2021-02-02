curl -i --location --request POST 'http://localhost:1026/ngsi-ld/v1/entityOperations/upsert?options=update' \
--header 'Content-Type: application/json' \
--header 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
--data-raw '[
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
]'


