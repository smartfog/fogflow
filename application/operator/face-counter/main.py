from flask import Flask, jsonify, abort, request, make_response
import requests 
import json
import datetime
import threading
import numpy as np
import urllib
import cv2
import openface
import os

app = Flask(__name__, static_url_path = "")

# global variables
brokerURL = ''
outputs = []
timer = None
lock = threading.Lock()
cameraURL = ''
total_size = 0 
align = openface.AlignDlib('/root/openface/models/dlib/shape_predictor_68_face_landmarks.dat')

@app.errorhandler(400)
def not_found(error):
    return make_response(jsonify( { 'error': 'Bad request' } ), 400)

@app.errorhandler(404)
def not_found(error):
    return make_response(jsonify( { 'error': 'Not found' } ), 404)

@app.route('/admin', methods=['POST'])
def admin():    
    if not request.json:
        abort(400)
    
    configObjs = request.json
    handleConfig(configObjs)    
        
    return jsonify({ 'responseCode': 200 })


@app.route('/notifyContext', methods = ['POST'])
def notifyContext():
    if not request.json:
        abort(400)
    
    objs = readContextElements(request.json)
    handleNotify(objs)
    
    return jsonify({ 'responseCode': 200 })


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

def readContextElements(data):
    print data

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
    global cameraURL
    print '===============receive context entity===================='
    print obj

    if 'attributes' in obj:
        attributes = obj['attributes']
        if 'url' in attributes:
            cameraURL = attributes['url']['value'] 
            
            with lock: 
                result = faceCounting(cameraURL)  #fetch the captured image from a web camera
                publishResult(result)  #publish the generated result to the configured broker

def handleConfig(configurations):  
    global brokerURL
    global num_of_outputs  
    for config in configurations:        
        if config['command'] == 'CONNECT_BROKER':
            brokerURL = config['brokerURL']
        if config['command'] == 'SET_OUTPUTS':
            outputs.append({'id': config['id'], 'type': config['type']})
    
def handleTimer():
    global timer

    if cameraURL != '':
        with lock: 
            result = faceCounting(cameraURL)  #fetch the captured image from a web camera
            publishResult(result)  #publish the generated result to the configured broker
        
    timer = threading.Timer(30, handleTimer)  # change to every 30 seconds
    timer.start()


def publishResult(result):
    resultCtxObj = {}
        
    #annotate the context with the configured entity id and type
    if len(outputs) < 1:
        return   
    
    resultCtxObj['entityId'] = {}
    resultCtxObj['entityId']['id'] = outputs[0]['id']
    resultCtxObj['entityId']['type'] = outputs[0]['type']        
    resultCtxObj['entityId']['isPattern'] = False    
    
    resultCtxObj['attributes'] = {}
    resultCtxObj['attributes']['num'] = {'type': 'integer', 'value': result['facenum']}
    resultCtxObj['attributes']['totalbytes'] = {'type': 'integer', 'value': result['totalbytes']}    

    # publish the real time results as context updates    
    updateContext(resultCtxObj)

def updateContext(ctxObj):
    global brokerURL
    if brokerURL == '':
        return
        
    ctxElement = object2Element(ctxObj)
    
    updateCtxReq = {}
    updateCtxReq['updateAction'] = 'UPDATE'
    updateCtxReq['contextElements'] = []
    updateCtxReq['contextElements'].append(ctxElement)

    headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}
    response = requests.post(brokerURL + '/updateContext', data=json.dumps(updateCtxReq), headers=headers)
    if response.status_code != 200:
        print 'failed to update context'
        print response.text


def url2Image(url): 
    global total_size
    resp = urllib.urlopen(url)
    data = resp.read()
    
    total_size += len(data)
    
    image = np.asarray(bytearray(data), dtype=np.uint8)    
    image = cv2.imdecode(image, cv2.CV_LOAD_IMAGE_COLOR)
    rgbImg = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)    
    return rgbImg

def faceCounting(url):     
    image = url2Image(url)
    if image is None:
        raise Exception("Unable to load image: {}".format(url))
        
    faces = align.getAllFaceBoundingBoxes(image)
    if faces is None:
        raise Exception("Unable to find a face: {}".format(url))

    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")      
    result = {"date": now, "facenum": len(faces), "totalbytes": total_size}    

    return result
    

if __name__ == '__main__':
    handleTimer()    
    
    myport = os.environ['myport']
    
    app.run(host='0.0.0.0', port=myport)
    
    timer.cancel()
    
