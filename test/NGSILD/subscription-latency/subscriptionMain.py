from flask import Flask, jsonify, abort, request, make_response
import requests
import json
import time
import datetime
import threading
import os
import sys

import ngsildClient as fogflow

app = Flask(__name__, static_url_path='')

start = False
startTime = ''
BrokerURL = ''
myport = 8085 
myIp = ''
noOfRequest = 2 

'''@app.route('/notifyContext', methods=['POST'])
def notify():
    //handle notification
    if start = True :
        startTime = ()
'''
# find The throughtput

def runApp():
    myport = int(8085)
    app.run(host='0.0.0.0', port=myport)


def handleConfig(config):
    global myIp
    global BrokerURL
    BrokerURL = config['subscribe_broker_url']
    myIp = config['my_ip']
    print(BrokerURL)
    print(myIp)

def setConfig():
    # load the configuration
    with open('config.json') as json_file:
        config = json.load(json_file)
        handleConfig(config)


def updateRequest(requestNo):
    global BrokerURL
    updateCtxRequest = {}
    updateCtxRequest['id'] = 'urn:ngsi-ld:Car:A0' + requestNo
    updateCtxRequest['type'] = 'Vehicle'
    brand = {}
    brand['type'] = 'property'
    brand['value'] = 'BMW'
    
    CarRelation = {}
    CarRelation['type'] = 'relationship'
    CarRelation['object'] = 'urn:ngsi-ld:Car:A111'
    updateCtxRequest['brand'] = brand 
    updateCtxRequest['CarRelation'] = CarRelation
    
    responseStartus = fogflow.updateRequest(updateCtxRequest,BrokerURL)
    print(responseStartus)

def update(noOfRequest):
    for requestNo in range(int(noOfRequest)):
        updateRequest(requestNo)

def subscribe():
    global myIp
    global BrokerURL
    subscriptionRequest = {}
    subscriptionRequest['type']  = 'Subscription'

    entities = []

    entity = {}
    entity['id'] = 'urn:ngsi-ld:Car:A0' + '.*'
    entity['type'] = 'Vehicle'
    entities.append(entity)

    subscriptionRequest['entities'] = entities

    notification = {}
    notification['format'] = 'keyValues'
    endPoint = {}

    endPoint['uri'] = str(myIp) + ':' + str(myport) + '/notifyContext'
    endPoint['accept'] = 'application/ld+json'
    notification['endpoint'] = endPoint

    subscriptionRequest['notification'] = notification
    print(subscriptionRequest)
    statusCode = fogflow.subscribeContext(subscriptionRequest,BrokerURL)
    if statusCode == '':
        print('check BrokerURL in config.json')
    else :
        print(statusCode)
    
if __name__ == '__main__':
     #run app
    setConfig()
    
    #sid = subscribe()
    
    #update(noOfRequest)

    runApp()
    
     

