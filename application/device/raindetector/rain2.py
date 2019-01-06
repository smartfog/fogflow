#!/usr/bin/env python
import time
import os
import signal
import sys
import json
import requests 
import RPi.GPIO as GPIO 


discoveryURL = 'http://192.168.1.80:8070/ngsi9'
brokerURL = ''
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
    
    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    
    updateContext(brokerURL, deviceCtxObj)

    # stream entity
    streamCtxObj = {}
    streamCtxObj['entityId'] = {}
    streamCtxObj['entityId']['id'] = 'Stream.' + profile['type'] + '.' + profile['id']
    streamCtxObj['entityId']['type'] = profile['type']        
    streamCtxObj['entityId']['isPattern'] = False
    
    streamCtxObj['attributes'] = {}
    streamCtxObj['attributes']['raining'] = {'type': 'bool', 'value': False}
    
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
    streamCtxObj['entityId']['type'] = profile['type']        
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
    streamCtxObj['entityId']['type'] = profile['type']        
    streamCtxObj['entityId']['isPattern'] = False
    
    streamCtxObj['attributes'] = {}
    streamCtxObj['attributes']['raining'] = {'type': 'bool', 'value': raining}
    
    streamCtxObj['metadata'] = {}
    streamCtxObj['metadata']['location'] = {'type': 'point', 'value': {'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude'] }}
    
    updateContext(brokerURL, streamCtxObj)        
        
def callback(channel):  
    if GPIO.input(channel):
        print "No Rain"
        updateMyStatus(False)
    else:
        print "Rain"  
        updateMyStatus(True)
  
def setRainDetector(pinNo):
    # Set our GPIO numbering to BCM
    GPIO.setmode(GPIO.BCM)

    # Define the GPIO pin that we have our digital output from our sensor connected to
    channel = pinNo
    # Set the GPIO pin to an input
    GPIO.setup(channel, GPIO.IN)

    # This line tells our script to keep an eye on our gpio pin and let us know when the pin goes HIGH or LOW
    GPIO.add_event_detect(channel, GPIO.BOTH, bouncetime=300)
    # This line asigns a function to the GPIO pin so that when the above line tells us there is a change on the pin, run this function
    GPIO.add_event_callback(channel, callback)

    
def run():
    global brokerURL
    brokerURL = findNearbyBroker()
    if brokerURL == '':
        print 'failed to find a nearby broker'
        sys.exit(0)
        
    #announce myself        
    publishMySelf()

    #set up the detector
    setRainDetector(17)    
    
    signal.signal(signal.SIGINT, signal_handler) 

    print('start to detect rain drops')      
    while True:
        time.sleep(5)
  
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

