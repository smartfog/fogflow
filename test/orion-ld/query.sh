curl -iX GET \
          'http://localhost:8070/ngsi-ld/v1/entities?type=Vehicle' \
          -H 'Content-Type: application/ld+json' \
          -H 'Accept: application/ld+json' \
          -H 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'

