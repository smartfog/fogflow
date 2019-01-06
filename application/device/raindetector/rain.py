#!/usr/bin/env python
import time
import os
import signal
import sys
import json
import requests 
import RPi.GPIO as io 

brokerURL = 'http://192.168.1.102:8070/ngsi10'
profile = {}

from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer

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
    discoveryReq['restriction'] = {'scopes':[{'type': 'nearby', 'value': nearby}]}
    
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
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + profile['id']
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['iconURL'] = {'type': 'string', 'value': profile['iconURL']};                           	
    deviceCtxObj['attributes']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }};
    
    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    
    updateContext(brokerURL, deviceCtxObj)

    # stream entity
    streamCtxObj = {}
    streamCtxObj['entityId'] = {}
    streamCtxObj['entityId']['id'] = 'Stream.' + profile['type'] + '.' + profile['id']
    streamCtxObj['entityId']['type'] = 'RainObservation'        
    streamCtxObj['entityId']['isPattern'] = False
    
    streamCtxObj['attributes'] = {}
    streamCtxObj['attributes']['raining'] = {'type': 'bool', 'value': False}
    streamCtxObj['attributes']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}    

    streamCtxObj['metadata'] = {}
    streamCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    
    updateContext(brokerURL, streamCtxObj)

def unpublishMySelf():
    global profile, brokerURL

    # device entity
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + profile['type'] + '.' + profile['id']
    deviceCtxObj['entityId']['type'] = profile['type']        
    deviceCtxObj['entityId']['isPattern'] = False
    
    deleteContext(brokerURL, deviceCtxObj)
    
    # stream entity
    streamCtxObj = {}
    streamCtxObj['entityId'] = {}
    streamCtxObj['entityId']['id'] = 'Stream.' + profile['type'] + '.' + profile['id']
    streamCtxObj['entityId']['type'] = 'RainObservation'        
    streamCtxObj['entityId']['isPattern'] = False
    
    deleteContext(brokerURL, streamCtxObj)    

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

def updateMyStatus(raining):                
    # stream entity
    streamCtxObj = {}
    streamCtxObj['entityId'] = {}
    streamCtxObj['entityId']['id'] = 'Stream.' + profile['type'] + '.' + profile['id']
    streamCtxObj['entityId']['type'] = 'RainObservation'      
    streamCtxObj['entityId']['isPattern'] = False
    
    streamCtxObj['attributes'] = {}
    streamCtxObj['attributes']['raining'] = {'type': 'bool', 'value': raining}
    
    streamCtxObj['metadata'] = {}
    streamCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    
    updateContext(brokerURL, streamCtxObj)        
        
    
def run():
    #global brokerURL
    #brokerURL = findNearbyBroker()
    #if brokerURL == '':
    #    print 'failed to find a nearby broker'
    #    sys.exit(0)
        
    #announce myself        
    publishMySelf()  
    
    signal.signal(signal.SIGINT, signal_handler) 

    print('start to detect rain drops')      

    water_sensor = 23
    io.setmode(io.BCM)
    io.setup(water_sensor, io.IN)
    
    # read the current state
    rain_detected = False
    
    if io.input(water_sensor):    
        print("No Rain Detected")
        rain_detected = False
        updateMyStatus(False)
    else:
        print("Rain Detected")    
        rain_detected = True
        updateMyStatus(True)

    while True:
        time.sleep(1.0)    
            
        if io.input(water_sensor):  
            if rain_detected == True: 
                print("No Rain Detected")
                rain_detected = False                
                updateMyStatus(False)
        else:
            if rain_detected == False:
                print("Rain Detected")
                rain_detected = True                
                updateMyStatus(True)

  
if __name__ == '__main__':
    cfgFileName = 'profile.json' 
    if len(sys.argv) >= 2:
        cfgFileName = sys.argv[1]
    
    try:       
        with open(cfgFileName) as json_file:
            profile = json.load(json_file)
    except Exception as error:
        print 'failed to load the device profile'
        sys.exit(0)

    run()

