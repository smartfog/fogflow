#!/usr/bin/env python
import time
import os
import signal
import sys
import json
import requests 
from flask import Flask, jsonify, abort, request, make_response
from threading import Thread, Lock
import logging
import nxt

app = Flask(__name__, static_url_path = "")

discoveryURL = 'http://192.168.1.80:8070/ngsi9'
brokerURL = ''
profile = {}
subscriptionID = ''

b = nxt.find_one_brick()
mxA = nxt.Motor(b, nxt.PORT_A)
mxB = nxt.Motor(b, nxt.PORT_B)


@app.route('/notifyContext', methods = ['POST'])
def notifyContext():
    if not request.json:
        abort(400)
    
    objs = readContextElements(request.json)
    handleNotify(objs)
    
    return jsonify({ 'responseCode': 200 })


def readContextElements(data):
    # print data

    ctxObjects = []
    
    for response in data['contextResponses']:
        if response['statusCode']['code'] == 200:
            ctxObj = element2Object(response['contextElement'])
            ctxObjects.append(ctxObj)
    
    return ctxObjects

def handleNotify(contextObjs):
    #print("received notification")
    #print(contextObjs)
    
    for ctxObj in contextObjs:
        processInputStreamData(ctxObj)

def processInputStreamData(obj):
    #print '===============receive context entity===================='
    #print obj
    
    if 'attributes' in obj:
        attributes = obj['attributes']
        
        if 'detectedEvent' in attributes:
            event = attributes['detectedEvent']['value']
            handleEvent(event)
        
        # if 'command' in attributes:
        #     command = attributes['command']['value']
        #     handleCommand(command)
   

def signal_handler(signal, frame):
    print('You pressed Ctrl+C!')
    # delete my registration and context entity
    unpublishMySelf()
   
    unsubscribeCmd()

    mxA.brake()
    mxB.brake()
    
    sys.exit(0)  

def findNearbyBroker():    
    global profile, discoveryURL

    nearby = {}
    nearby['latitude'] = profile['location']['latitude']
    nearby['longitude'] = profile['location']['longitude'] 
    nearby['limit'] = 1
 
    discoveryReq = {}
    discoveryReq['entities'] = [{'type': 'IoTBroker', 'isPattern': True}]
    discoveryReq['restriction'] = {'scopes':[{'scopeType': 'nearby', 'scopeValue': nearby}]}
    
    discoveryURL = profile['discoveryURL']
    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}
    response = requests.post(discoveryURL + '/discoverContextAvailability', data=json.dumps(discoveryReq), headers=headers)
    if response.status_code != 200:
        print 'failed to find a nearby IoT Broker'
        return ''
    
    print response.text
    registrations = json.loads(response.text)
    
    for registration in registrations['contextRegistrationResponses']:
        providerURL = registration['contextRegistration']['providingApplication']
        if providerURL != '':
            return providerURL
          
    return '' 
    
def publishMySelf():
    global profile, brokerURL
    
    # for motor1
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + '001'
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['iconURL'] = {'type': 'string', 'value': profile['iconURL']}    
    
    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
        
    updateContext(brokerURL, deviceCtxObj)

    # for motor2
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + '002'
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['iconURL'] = {'type': 'string', 'value': profile['iconURL']}    
    
    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
        
    updateContext(brokerURL, deviceCtxObj)

def unpublishMySelf():
    global profile, brokerURL

    # for motor1
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + '001'
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deleteContext(brokerURL, deviceCtxObj)

    # for motor2
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + '002'
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deleteContext(brokerURL, deviceCtxObj)    

def element2Object(element):
    ctxObj = {}
    
    ctxObj['entityId'] = element['entityId'];
    
    ctxObj['attributes'] = {}  
    if 'attributes' in element:
        for attr in element['attributes']:
            ctxObj['attributes'][attr['name']] = {'type': attr['type'], 'value': attr['value']}   
    
    ctxObj['metadata'] = {}
    if 'domainMetadata' in element:    
        for meta in element['domainMetadata']:
            ctxObj['metadata'][meta['name']] = {'type': meta['type'], 'value': meta['value']}
    
    return ctxObj


def object2Element(ctxObj):
    ctxElement = {}
    
    ctxElement['entityId'] = ctxObj['entityId'];
    
    ctxElement['attributes'] = []  
    if 'attributes' in ctxObj:
        for key in ctxObj['attributes']:
            attr = ctxObj['attributes'][key]
            ctxElement['attributes'].append({'name': key, 'type': attr['type'], 'value': attr['value']})
    
    ctxElement['domainMetadata'] = []
    if 'metadata' in ctxObj:    
        for key in ctxObj['metadata']:
            meta = ctxObj['metadata'][key]
            ctxElement['domainMetadata'].append({'name': key, 'type': meta['type'], 'value': meta['value']})
    
    return ctxElement


def updateContext(broker, ctxObj):        
    ctxElement = object2Element(ctxObj)
    
    updateCtxReq = {}
    updateCtxReq['updateAction'] = 'UPDATE'
    updateCtxReq['contextElements'] = []
    updateCtxReq['contextElements'].append(ctxElement)

    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}
    response = requests.post(broker + '/updateContext', data=json.dumps(updateCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to update context'
        print response.text


def deleteContext(broker, ctxObj):        
    ctxElement = object2Element(ctxObj)
    
    updateCtxReq = {}
    updateCtxReq['updateAction'] = 'DELETE'
    updateCtxReq['contextElements'] = []
    updateCtxReq['contextElements'].append(ctxElement)

    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}
    response = requests.post(broker + '/updateContext', data=json.dumps(updateCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to delete context'
        print response.text

def unsubscribeCmd():
    global brokerURL
    global subscriptionID

    print(brokerURL + '/subscription/' + subscriptionID)
    response = requests.delete(brokerURL + '/subscription/' + subscriptionID)
    print(response.text)

def subscribeCmd():                    
    global subscriptionID
    subscribeCtxReq = {}
    subscribeCtxReq['entities'] = []
    
    # subscribe push button on behalf of TPU
    myID = 'Device.Motor.001'
    
    subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})  
    #subscribeCtxReq['attributes'] = ['command']      
    subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':' + str(profile['myPort'])
    
    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json', 'Require-Reliability' : 'true'}
    response = requests.post(brokerURL + '/subscribeContext', data=json.dumps(subscribeCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to subscribe context'
        print response.text  
        return ''          
    else:
        json_data = json.loads(response.text)
        subscriptionID = json_data['subscribeResponse']['subscriptionId']
        print(subscriptionID)
        return subscriptionID
  

    # # subscribe to motor1
    # myID = 'Device.' + profile['type'] + '.' + '001'   
    
    # subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})  
    # subscribeCtxReq['attributes'] = ['command']      
    # subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':' + str(profile['myPort'])
    
    # headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json', 'Require-Reliability' : 'true'}
    # response = requests.post(brokerURL + '/subscribeContext', data=json.dumps(subscribeCtxReq), headers=headers)
    # if response.status_code != 200:
    #     print 'failed to subscribe context'
    #     print response.text          
  
    # # subscribe to motor2
    # myID = 'Device.' + profile['type'] + '.' + '002'   
    
    # subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})  
    # subscribeCtxReq['attributes'] = ['command']      
    # subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':' + str(profile['myPort'])
    
    # headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json', 'Require-Reliability' : 'true'}
    # response = requests.post(brokerURL + '/subscribeContext', data=json.dumps(subscribeCtxReq), headers=headers)
    # if response.status_code != 200:
    #     print 'failed to subscribe context'
    #     print response.text    

    # # subscribe camera on behalf of TPU
    # myID = 'Device.Camera.001'   
    
    # subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})  
    # subscribeCtxReq['attributes'] = ['command']      
    # subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':8008'
    
    # headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json', 'Require-Reliability' : 'true'}
    # response = requests.post(brokerURL + '/subscribeContext', data=json.dumps(subscribeCtxReq), headers=headers)
    # if response.status_code != 200:
    #     print 'failed to subscribe context'
    #     print response.text    

def run():
	# find a nearby broker for data exchange
    global brokerURL
    brokerURL = profile['brokerURL'] #findNearbyBroker()
    if brokerURL == '':
        print 'failed to find a nearby broker'
        sys.exit(0)        

    print "selected broker"
    print brokerURL    
                    
    #announce myself        
    publishMySelf()

    #subscribe to the control commands
    while true:
        sid = subscribeCmd()
        if sid != '':
            break
    
    signal.signal(signal.SIGINT, signal_handler) 
    signal.signal(signal.SIGTERM, signal_handler)

    print('start to handle the incoming control commands')     
    myport = profile['myPort']     
    app.run(host='0.0.0.0', port=myport)

def handleEvent(event):
    print event
    
    eventType = event['type']
    print(eventType)   

    if eventType == 'MOVE_FORWARD': 
        print("MOVE_FORWARD")
        mxA.run(-100)       
        time.sleep(1)
        mxA.brake()   
    
    if eventType == 'MOVE_LEFT': 
        print("MOVE_LEFT")
        mxB.run(100)        
        time.sleep(1)
        mxB.brake()
    
    if eventType == 'MOVE_RIGHT': 
        print("MOVE_RIGHT")
        mxB.run(-100)        
        time.sleep(1)
        mxB.brake()
    
    if eventType == 'MOVE_BACKWARD':
        print("MOVE_BACKWARD")
        mxA.run(100)
        time.sleep(1)
        mxA.brake()


if __name__ == '__main__':
    cfgFileName = 'motor.json' 
    if len(sys.argv) >= 2:
        cfgFileName = sys.argv[1]
    
    try:       
        with open(cfgFileName) as json_file:
            profile = json.load(json_file)
            
        profile['type'] = 'Motor'            
    except Exception as error:
        print 'failed to load the device profile'
        sys.exit(0)  
        
    run()
