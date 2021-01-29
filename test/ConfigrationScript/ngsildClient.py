import requests
import json


# request that send the subscribe Context request

def subscribeContext(subscribeCtxEle, BrokerURL):
    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/json',
               'Link': '<https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    response = requests.post(BrokerURL + '/ngsi-ld/v1/subscriptions/',
                             data=json.dumps(subscribeCtxEle),
                             headers=headers)
    if response.status_code == 201:
        return response.status_code
    else:
        return ''


# request that send the update Context request to Broker
def updateRequest(updateCtxEle, FogFlowBrokerURL):
    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/json',
               'Link': '<https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    response = requests.post(FogFlowBrokerURL + '/ngsi-ld/v1/entities/',
                             data=json.dumps(updateCtxEle),
                             headers=headers)

    if response.status_code == 201:
        return response.status_code
    else:
        return ''

