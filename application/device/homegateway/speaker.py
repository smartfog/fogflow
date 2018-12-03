#!/usr/bin/env python
import time
import os
import signal
import sys
import json
import requests 
import wiringpi
from flask import Flask, jsonify, abort, request, make_response
from threading import Thread, Lock
import logging

app = Flask(__name__, static_url_path = "")

discoveryURL = 'http://192.168.1.80:8070/ngsi9'
brokerURL = ''
profile = {}

from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer


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
    for ctxObj in contextObjs:
        processInputStreamData(ctxObj)

def processInputStreamData(obj):
    print '===============receive context entity===================='
    print obj
    
    if 'attributes' in obj:
        attributes = obj['attributes']
        if 'command' in attributes:
            alert = attributes['command']['value']
            handleCommand(alert)


def handleCommand(alert):
    print alert
    
    alertType = alert['type']
        
    if alertType == 'ASK_FOR_HELP': 
        #os.system('mpg123 sos.mp3 &')
		print "play sos.mp3 when someone asks for helps"
    elif alertType == 'EMERGENCY': 
        #os.system('mpg123 alarm.mp3 &')    
		print "play alarm.mp3 when there is an emergency"
    elif alertType == 'BROADCAST': 
        #os.system('========broadcast======== &')       
		print "play out an audio file broadcasted by city console" 

def signal_handler(signal, frame):
    print('You pressed Ctrl+C!')
    # delete my registration and context entity
    unpublishMySelf()
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
    
    # device entity
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + profile['homeID']
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['iconURL'] = {'type': 'string', 'value': profile['iconURL']}    
    
    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    deviceCtxObj['metadata']['homeID'] = {'type': 'string', 'value': profile['homeID']}        
        
    updateContext(brokerURL, deviceCtxObj)

def unpublishMySelf():
    global profile, brokerURL

    # device entity
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + profile['homeID']
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deleteContext(brokerURL, deviceCtxObj)
    

def element2Object(element):
    ctxObj = {}
    
    ctxObj['entityId'] = element['entityId'];
    
    ctxObj['attributes'] = {}  
    if 'attributes' in element:
        for attr in element['attributes']:
            ctxObj['attributes'][attr['name']] = {'type': attr['type'], 'value': attr['contextValue']}   
    
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
            ctxElement['attributes'].append({'name': key, 'type': attr['type'], 'contextValue': attr['value']})
    
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

def subscribeCmd():                    
    subscribeCtxReq = {}
    subscribeCtxReq['entities'] = []
    
    myID = 'Device.' + profile['type'] + '.' + profile['homeID']    
    
    subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})  
    subscribeCtxReq['attributes'] = ['command']      
    subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':' + str(profile['myPort'])
    
    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json', 'Require-Reliability' : 'true'}
    response = requests.post(brokerURL + '/subscribeContext', data=json.dumps(subscribeCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to subscribe context'
        print response.text      
        
    
def run():
	# find a nearby broker for data exchange
    global brokerURL
    brokerURL = findNearbyBroker()
    if brokerURL == '':
        print 'failed to find a nearby broker'
        sys.exit(0)        

    print "selected broker"
    print brokerURL    
                    
    #announce myself        
    publishMySelf()

    #subscribe to the control commands
    subscribeCmd()
    
    signal.signal(signal.SIGINT, signal_handler) 

    print('start to handle the incoming control commands')     
    myport = profile['myPort']     
    app.run(host='0.0.0.0', port=myport)

  
if __name__ == '__main__':
    cfgFileName = 'profile.json' 
    if len(sys.argv) >= 2:
        cfgFileName = sys.argv[1]
    
    try:       
        with open(cfgFileName) as json_file:
            profile = json.load(json_file)
            
        profile['type'] = 'Speaker'            
    except Exception as error:
        print 'failed to load the device profile'
        sys.exit(0)  
        
    run()

