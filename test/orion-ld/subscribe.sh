curl -iX POST \
          'http://localhost:1026/ngsi-ld/v1/subscriptions' \
          -H 'Content-Type: application/json' \
          -H 'Accept: application/ld+json' \
          -H 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
          -d ' {
        "type": "Subscription",
        "entities": [{
           "type": "Vehicle"
        }],
        "notification": {
           "format": "keyValues",
           "endpoint": {
               "uri": "http://192.168.0.59:8070/ngsi-ld/v1/notifyContext/",
               "accept": "application/ld+json"
            }
        }
                }'

