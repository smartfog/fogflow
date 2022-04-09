import json
import time
import datetime
import threading
import os
import sys

import ngsildClient as fogflow

NGSILDBrokerURL = ''
FogFlowBrokerURL = ''
myIp = ''
id = ''
type = ''


def handleConfig(config):
    global myIp, NGSILDBrokerURL, FogFlowBrokerURL, myIp, id,  type
    global myIp
    NGSILDBrokerURL = config['subscribe_broker_url']
    myIp = config['my_ip']
    FogFlowBrokerURL = config['update_broker_url']
    id = config['id']
    type = config['type']


def setConfig():
    with open('config.json') as json_file:
        config = json.load(json_file)
        handleConfig(config)


def updateRequest():
    global myIp, NGSILDBrokerURL, FogFlowBrokerURL, myIp, id, type
    updateCtxRequest = {}
    updateCtxRequest['id'] = id
    updateCtxRequest['type'] = type
    brand = {}
    brand['type'] = 'property'
    brand['value'] = 'BMW'

    updateCtxRequest['brand'] = brand

    responseStartus = fogflow.updateRequest(updateCtxRequest, FogFlowBrokerURL)
    print(responseStartus)


def subscribe():
    global myIp, NGSILDBrokerURL, FogFlowBrokerURL, myIp, id, type
    subscriptionRequest = {}
    subscriptionRequest['type'] = 'Subscription'

    entities = []

    entity = {}
    entity['id'] = id
    entity['type'] = type
    entities.append(entity)

    subscriptionRequest['entities'] = entities

    notification = {}
    notification['format'] = 'keyValues'
    endPoint = {}

    endPoint['uri'] = str(NGSILDBrokerURL) + '/ngsi-ld/v1/notifyContext/'
    endPoint['accept'] = 'application/ld+json'
    notification['endpoint'] = endPoint

    subscriptionRequest['notification'] = notification
    print(subscriptionRequest)
    statusCode = fogflow.subscribeContext(subscriptionRequest, NGSILDBrokerURL)
    if statusCode == '':
        print('check BrokerURL in config.json')
    else:
        print(statusCode)


if __name__ == '__main__':
    # run app
    setConfig()
    subscribe()
    updateRequest()
