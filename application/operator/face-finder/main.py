from flask import Flask, jsonify, abort, request, make_response
import requests 
import json
import time
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
camera = None
cameraURL = ''
total_size = 0 

threshold = 1.0
imageDemension = 96
align = openface.AlignDlib('/root/openface/models/dlib/shape_predictor_68_face_landmarks.dat')
net = openface.TorchNeuralNet('/root/openface/models/openface/nn4.small2.v1.t7', imageDemension)

saveLocation = ''
featuresOfTarget = None
targetedFeaturesIsSet = False

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
    print "notify"

    if not request.json:
        abort(400)
    
    print(request.json)
    
    objs = readContextElements(request.json)
    
    print(objs)

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
    print '===============receive context entity===================='
    print obj

    entityId = obj['entityId']
    if entityId['type'] == 'Camera':
        getCameraURL(obj)
    elif entityId['type'] == 'ChildLost':
        getChildInfo(obj)
    
    with lock: 
        faceMatching()

def getCameraURL(entityObj):
    global camera, cameraURL
    
    camera = entityObj    
    if 'attributes' in entityObj:
        attributes = entityObj['attributes']
        if 'url' in attributes:
            cameraURL = attributes['url']['value']    

def getChildInfo(entityObj):
    global featuresOfTarget, saveLocation, targetedFeaturesIsSet
    if 'attributes' in entityObj:
        attributes = entityObj['attributes']
        if 'imageURL' in attributes:
            imageURL = attributes['imageURL']['value']        
            image =  url2Image(imageURL)
            
            if targetedFeaturesIsSet == False: 
                featuresOfTarget = getRep(image)
                targetedFeaturesIsSet = True
            
            saveLocation = attributes['saveLocation']['value']  

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

    with lock: 
        faceMatching()  
        
    timer = threading.Timer(10, handleTimer)
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
    resultCtxObj['attributes']['cameraID'] = {'type': 'string', 'value': result['cameraID']}
    resultCtxObj['attributes']['where'] = {'type': 'object', 'value': result['location']}    
    resultCtxObj['attributes']['when'] = {'type': 'string', 'value': result['date']}        
    resultCtxObj['attributes']['image'] = {'type': 'string', 'value': result['image']}            

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

    
def getRep(img):
    bb = align.getLargestFaceBoundingBox(img)
    if bb is None:
        raise Exception("Unable to find a face: {}")
        
    alignedFace = align.align(imageDemension, img, bb, landmarkIndices=openface.AlignDlib.OUTER_EYES_AND_NOSE)
    if alignedFace is None:
        raise Exception("Unable to align image: {}")

    rep = net.forward(alignedFace)
    return rep
    

def faceMatching():     
    global camera, cameraURL, featuresOfTarget, targetedFeaturesIsSet
    
    if camera == None or cameraURL == '' or targetedFeaturesIsSet == False:        
        print 'parameters are not yet set', camera, cameraURL, featuresOfTarget
        return     

    image = url2Image(cameraURL)
    if image is None:
        raise Exception("Unable to load image: {}".format(camera))
        
    faces = align.getAllFaceBoundingBoxes(image)
    if faces is None:
        raise Exception("Unable to find a face: {}".format(camera))
    
    for face in faces:
        alignedFace = align.align(96, image, face, landmarkIndices=openface.AlignDlib.OUTER_EYES_AND_NOSE)
        if alignedFace is None:
            raise Exception("Unable to align image: {}")

        rep = net.forward(alignedFace)
        d = rep - featuresOfTarget
        print("  + Squared l2 distance between representations: {:0.3f}".format(np.dot(d, d)))
            
        distance = np.dot(d, d);
        if  distance < threshold:
            distance = "{:0.3f}".format(distance);        
            now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")              
        
            # save the image and post it to the remote image server
            fileName = 'childfound-' + camera['metadata']['cameraID']['value'] + '-' + str(int(time.time())) + '.png'
            bgrImg = cv2.cvtColor(image, cv2.COLOR_RGB2BGR)    
            cv2.imwrite(fileName, bgrImg)
            files = {fileName: open(fileName, 'rb')}
            requests.post(saveLocation, files=files)
            
            result = {  "date": now, 
                        "cameraID": camera['metadata']['cameraID']['value'], 
                        "location": camera['metadata']['location']['value'], 
                        "delta": distance, 
                        "image": saveLocation + '/' + fileName,
                        "totalbytes": total_size
                    }
            
            
            #update context
            publishResult(result)  
            
                                

if __name__ == '__main__':
    #handleTimer()    
    
    myport = os.environ['myport']
    app.run(host='0.0.0.0', port=myport)
    
    timer.cancel()
    
