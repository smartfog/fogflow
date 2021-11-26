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
import sys


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
    print ('=============notify=============')
    sys.stdout.flush()
    print (request.json)
    sys.stdout.flush()
    if not request.json:
        abort(400)

    objs = readContextElements(request.json)
    print(objs)
    sys.stdout.flush()
    handleNotify(objs)
    return jsonify({'responseCode': 200})

def element2Object(element):
    ctxObj = {}
    for key in element:
        ctxObj[key] = element[key]

    return ctxObj 

def object2Element(ctxObj):
    ctxElement = {}
    ctxElement['id'] = ctxObj['id']
    ctxElement['type'] = ctxObj['type']

    for key in ctxObj:
        if key != 'id' and key != 'type' and key != 'modifiedAt' \
            and key != 'createdAt' and key != 'observationSpace' \
            and key != 'operationSpace' and key != 'location' and key \
            != '@context':
            if ctxObj[key].has_key('createdAt'):
                ctxObj[key].pop('createdAt')
            if ctxObj[key].has_key('modifiedAt'):
                ctxObj[key].pop('modifiedAt')
            ctxElement[key] = ctxObj[key]

    return ctxElement    

def readContextElements(data):
    ctxObjects = []
    print (data['type'])
    if data['type'] == 'Notification':
        for attr in data['data']:
            ctxObj = element2Object(attr)
            ctxObjects.append(ctxObj)
    return ctxObjects

def handleNotify(contextObjs):
    for ctxObj in contextObjs:
        processInputStreamData(ctxObj)

def processInputStreamData(obj):
    print '===============receive context entity===================='
    print obj

    entityId = obj['id']
    print(entityId)
    type1 = obj['type']
    print(type1)
    if obj['type'] == 'lDCamera':
        getCameraURL(obj)
    elif obj['type'] == 'ChildLost':
        getChildInfo(obj)
    
    with lock: 
        faceMatching()

def getCameraURL(entityObj):
    global camera, cameraURL
    
    camera = entityObj    
    if 'url' in entityObj:
        cameraURL = entityObj['url']['value']   
        print("====== camera url ====",cameraURL)

def getChildInfo(entityObj):
    global featuresOfTarget, saveLocation, targetedFeaturesIsSet
    if 'imageURL' in entityObj:
        imageURL = entityObj['imageURL']['value']        
        image =  url2Image(imageURL)
            
        if targetedFeaturesIsSet == False: 
           featuresOfTarget = getRep(image)
           targetedFeaturesIsSet = True
            
        saveLocation = entityObj['saveLocation']['value']  

def handleConfig(configurations):  
    global brokerURL
    global num_of_outputs  
    for config in configurations:
        if config['command'] == 'CONNECT_BROKER':
            brokerURL = config['brokerURL']
	    print("********brkerURL********",brokerURL)
        elif config['command'] == 'SET_OUTPUTS':
            outputs.append({'id': config['id'], 'type': config['type']})
        elif config['command'] == 'SET_INPUTS':
            setInput(config)

def setInput(cmd):
    global input
    print ('cmd')
    print (cmd)
    if 'id' in cmd:
        input['id'] = cmd['id']

    input['type'] = cmd['type']


def handleTimer():
    global timer

    with lock: 
        faceMatching()  
        
    timer = threading.Timer(10, handleTimer)
    timer.start()


def publishResult(ctxObj):
    global brokerURL
    if brokerURL.endswith('/ngsi10') == True:
        brokerURL = brokerURL.rsplit('/', 1)[0]
    if brokerURL == '':
        return

    #ctxElement = object2Element(ctxObj)
    ctxObj['id'] = "urn:ngsi-ld:Device."+outputs[0]['id']
    ctxObj['type'] = outputs[0]['type']

    print("======== published result =============")
    sys.stdout.flush()
    print(ctxObj)
    sys.stdout.flush()

    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/json',
               'Link': '<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld >; rel=http://www.w3.org/ns/json-ld#context"; dtype="application/ld+json"'}
    response = requests.post(brokerURL + '/ngsi-ld/v1/entities/',
                             data=json.dumps(ctxObj),
                             headers=headers)
    if response.status_code != 201:
        print ('failed to update context')
        sys.stdout.flush()
        print (response.text)
        sys.stdout.flush()

def fetchInputByQuery():
    ctxQueryReq = {}

    ctxQueryReq['entities'] = []
    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/ld+json'}
    response = requests.get(brokerURL + '/ngsi-ld/v1/entities/'
                            + input['id'], headers=headers)

    if response.status_code != 200:
        print ('failed to query the input data')
        return {}
    else:
        jsonResult = response.json()

        ctxObj = element2Object(jsonResult)
        ctxElments = []
        ctxElments.append(ctxObj)
        return ctxElments


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
            fileName = 'childfound-' + camera['cameraID']['value'] + '-' + str(int(time.time())) + '.png'
            bgrImg = cv2.cvtColor(image, cv2.COLOR_RGB2BGR)    
            cv2.imwrite(fileName, bgrImg)
            files = {fileName: open(fileName, 'rb')}
            requests.post(saveLocation, files=files)
           
	    print("camera***********",camera) 
            result = {  "date": {"type":"Property", "value":now}, 
                        "cameraID": {"type":"Property", "value":camera['cameraID']['value']}, 
                        "where": {"type":"Property", "value":[35.0878,140.578]}, 
                        "delta": {"type":"Property", "value":distance}, 
                        "image": {"type":"Property", "value":saveLocation + '/' + fileName},
                        "totalbytes": {"type":"Property", "value":total_size}
                    }
            
            print("**** Result *******",result)
            #update context
            publishResult(result)  
            

def requestInputBySubscription():
    ctxSubReq = {}

    ctxSubReq['entities'] = []

    if id in input:
        ctxSubReq['entities'].append({'id': input['id'],
                'isPattern': False})
    else:
        ctxSubReq['entities'].append({'type': input['type'],
                'isPattern': True})

    subrequestUri['uri'] = 'http://host.docker.internal:' \
        + os.getenv('myport')
    subrequestEndPoint['endpoint'] = subrequestUri
    ctxSubReq['notification'] = subrequestEndPoint

    ctxSubReq['reference'] = 'http://host.docker.internal:' \
        + os.getenv('myport')

    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/ld+json',
               'Link': '<https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    if brokerURL.endswith('/ngsi10') == True:
        brokerURL = brokerURL.rsplit('/', 1)[0]

    response = requests.post(brokerURL + '/ngsi-ld/v1/subscriptions/',
                             data=json.dumps(ctxSubReq),
                             headers=headers)

    if response.status_code != 200:
        print ('failed to query the input data')
    else:
        print ('subscribed to the input data')


# continuous execution to handle received notifications

def notify2execution():
    myport = int(os.getenv('myport'))
    print ('listening on port ' + os.getenv('myport'))

    app.run(host='0.0.0.0', port=myport)
    timer.cancel()


def runInOperationMode():
    print ('===== OPERATION MODEL========')

    # apply the configuration received from the environmental varible

    myCfg = os.getenv('adminCfg')

    print (myCfg)

    adminCfg = json.loads(myCfg)
    handleConfig(adminCfg)

    syncMode = os.getenv('sync')
    if syncMode != None and syncMode == 'yes':
        query2execution()
    else:
        notify2execution()


# one time execution triggered by query

def query2execution():
    ctxObjects = fetchInputByQuery()
    handleNotify(ctxObjects)


def runInTestMode():
    print ('===== TEST NGSILD MODEL========')

    # load the configuration

    with open('config.json') as json_file:
        config = json.load(json_file)
        print (config)

        handleConfig(config)

        # trigger the data processing

        query2execution()


if __name__ == '__main__':
    parameters = sys.argv

    if len(parameters) == 2 and parameters[1] == '-o':
        runInOperationMode()
    else:
        runInTestMode()    
    
