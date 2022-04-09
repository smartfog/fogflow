#!/usr/bin/env python
import time
import os
import signal
import sys
import json
import requests
from datetime import datetime
import serial
import binascii

discoveryURL = 'http://192.168.1.100/ngsi9'
brokerURL = ''
profile = {}
subscriptionID = ''


def handle_exit(sig, frame):
    unpublishMySelf()
    unsubscribe()
    raise(SystemExit)


signal.signal(signal.SIGINT, handle_exit)
signal.signal(signal.SIGTERM, handle_exit)


def subscribe():
    global brokerURL
    global subscriptionID

    subscribeCtxReq = {}
    subscribeCtxReq['entities'] = []

    # subscribe push button on behalf of TPU
    myID = 'Device.Pushbutton.0001'

    subscribeCtxReq['entities'].append({'id': myID, 'isPattern': False})
    subscribeCtxReq['reference'] = 'http://' + profile['myIP'] + ':8008'

    headers = {'Accept': 'application/json',
               'Content-Type': 'application/json', 'Require-Reliability': 'true'}
    response = requests.post(brokerURL + '/subscribeContext',
                             data=json.dumps(subscribeCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to subscribe context'
        print response.text
        return ''
    else:
        json_data = json.loads(response.text)
        subscriptionID = json_data['subscribeResponse']['subscriptionId']
        print(subscriptionID)
        return subscriptionID


def unsubscribe():
    print(brokerURL + '/subscription/' + subscriptionID)

    response = requests.delete(brokerURL + '/subscription/' + subscriptionID)

    print(response.text)


def findNearbyBroker():
    global profile, discoveryURL

    nearby = {}
    nearby['latitude'] = profile['location']['latitude']
    nearby['longitude'] = profile['location']['longitude']
    nearby['limit'] = 1

    discoveryReq = {}
    discoveryReq['entities'] = [{'type': 'IoTBroker', 'isPattern': True}]
    discoveryReq['restriction'] = {'scopes': [
        {'scopeType': 'nearby', 'scopeValue': nearby}]}

    discoveryURL = profile['discoveryURL']
    headers = {'Accept': 'application/json',
               'Content-Type': 'application/json'}
    response = requests.post(discoveryURL + '/discoverContextAvailability',
                             data=json.dumps(discoveryReq), headers=headers)
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
    deviceCtxObj['entityId']['id'] = 'Device.' + \
        profile['type'] + '.' + profile['homeID']
    deviceCtxObj['entityId']['type'] = profile['type']
    deviceCtxObj['entityId']['isPattern'] = False

    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['iconURL'] = {
        'type': 'string', 'value': profile['iconURL']}

    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['location'] = {'type': 'point', 'value': {
        'latitude': profile['location']['latitude'], 'longitude': profile['location']['longitude']}}
    deviceCtxObj['metadata']['homeID'] = {
        'type': 'string', 'value': profile['homeID']}

    return updateContext(brokerURL, deviceCtxObj)


def unpublishMySelf():
    global profile, brokerURL

    # device entity
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + \
        profile['type'] + '.' + profile['homeID']
    deviceCtxObj['entityId']['type'] = profile['type']
    deviceCtxObj['entityId']['isPattern'] = False

    deleteContext(brokerURL, deviceCtxObj)


def object2Element(ctxObj):
    ctxElement = {}

    ctxElement['entityId'] = ctxObj['entityId']

    ctxElement['attributes'] = []
    if 'attributes' in ctxObj:
        for key in ctxObj['attributes']:
            attr = ctxObj['attributes'][key]
            ctxElement['attributes'].append(
                {'name': key, 'type': attr['type'], 'value': attr['value']})

    ctxElement['domainMetadata'] = []
    if 'metadata' in ctxObj:
        for key in ctxObj['metadata']:
            meta = ctxObj['metadata'][key]
            ctxElement['domainMetadata'].append(
                {'name': key, 'type': meta['type'], 'value': meta['value']})

    return ctxElement


def updateContext(broker, ctxObj):
    ctxElement = object2Element(ctxObj)

    updateCtxReq = {}
    updateCtxReq['updateAction'] = 'UPDATE'
    updateCtxReq['contextElements'] = []
    updateCtxReq['contextElements'].append(ctxElement)

    headers = {'Accept': 'application/json',
               'Content-Type': 'application/json'}
    response = requests.post(broker + '/updateContext',
                             data=json.dumps(updateCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to update context'
        print response.text
        return False
    else:
        return True


def deleteContext(broker, ctxObj):
    ctxElement = object2Element(ctxObj)

    updateCtxReq = {}
    updateCtxReq['updateAction'] = 'DELETE'
    updateCtxReq['contextElements'] = []
    updateCtxReq['contextElements'].append(ctxElement)

    headers = {'Accept': 'application/json',
               'Content-Type': 'application/json'}
    response = requests.post(broker + '/updateContext',
                             data=json.dumps(updateCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to delete context'
        print response.text


def reportEvent(eType):
    print eType

    # update my device profile with the latest observation
    deviceCtxObj = {}
    deviceCtxObj['entityId'] = {}
    deviceCtxObj['entityId']['id'] = 'Device.' + \
        profile['type'] + '.' + profile['homeID']
    deviceCtxObj['entityId']['type'] = profile['type']
    deviceCtxObj['entityId']['isPattern'] = False

    detectedEvent = {}
    detectedEvent['type'] = eType
    detectedEvent['time'] = str(datetime.now())

    deviceCtxObj['attributes'] = {}
    deviceCtxObj['attributes']['detectedEvent'] = {
        'type': 'object', 'value': detectedEvent}

    deviceCtxObj['metadata'] = {}
    deviceCtxObj['metadata']['homeID'] = {
        'type': 'string', 'value': profile['homeID']}

    updateContext(brokerURL, deviceCtxObj)


def run():
    # find a nearby broker for data exchange
    global brokerURL
    brokerURL = profile['brokerURL']  # findNearbyBroker()
    if brokerURL == '':
        print 'failed to find a nearby broker'
        sys.exit(0)

    print(brokerURL)

    # announce myself to the nearby broker
    while True:
        ok = publishMySelf()
        if ok == True:
            break
        else:
            time.sleep(1)

    print("publish myself")

    while True:
        sid = subscribe()
        if sid != '':
            break
        else:
            time.sleep(1)

    # detect the button-push event
    ser = serial.Serial('/dev/ttyUSB0', 57600, timeout=1)
    state = "off"
    press_time = datetime.now()
    while True:
        try:
            data = ser.read()
            number = binascii.hexlify(data)

            if number == '55':
                print("STATE = %s" % (state))

                if state == "off":
                    state = "on"
                    press_time = datetime.now()
                    print("BUTTON_PRESS")
                else:
                    print("BUTTON_RELEASE")
                    release_time = datetime.now()
                    delta = release_time - press_time

                    # print(delta)
                    print(delta.seconds)

                    if delta.seconds > 5:
                        reportEvent("RESET")
                    else:
                        reportEvent("CLICK")

                    state = "off"

        except KeyboardInterrupt:
            print('You pressed Ctrl+C!')
            # delete my registration and context entity
            unpublishMySelf()
            break
        except Exception as e:
            continue


if __name__ == '__main__':
    cfgFileName = 'button.json'
    if len(sys.argv) >= 2:
        cfgFileName = sys.argv[1]

    try:
        with open(cfgFileName) as json_file:
            profile = json.load(json_file)

        profile['type'] = 'Pushbutton'

    except Exception as error:
        print 'failed to load the device profile'
        sys.exit(0)

    run()
