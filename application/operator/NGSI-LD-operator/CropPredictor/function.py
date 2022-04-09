import joblib
import pandas as pd
import sys
import json
from datetime import datetime

crops = ['wheat', 'mungbean', 'Tea', 'millet', 'maize', 'lentil', 'jute', 'cofee', 'cotton', 'ground nut', 'peas', 'rubber', 'sugarcane', 'tobacco', 'kidney beans', 'moth beans',
         'coconut', 'blackgram', 'adzuki beans', 'pigeon peas', 'chick peas', 'banana', 'grapes', 'apple', 'mango', 'muskmelon', 'orange', 'papaya', 'watermelon', 'pomegranate']
cr = 'rice'


loaded_rf = joblib.load("./croppredictor.joblib")


def handleEntity(ctxObj, publish):
    print('===============Implement logic====================')

    print(ctxObj)
    sys.stdout.flush()
    print(ctxObj["type"])
    sys.stdout.flush()
    print(ctxObj["airmoisture"]["value"])
    sys.stdout.flush()

    atemp = ctxObj["airTemp"]["value"]
    shum = ctxObj["soilmoisture"]["value"]
    pH = ctxObj["soilpH"]["value"]
    rain = ctxObj["rainfall"]["value"]
    ah = ctxObj["airmoisture"]["value"]

    l = []
    l.append(ah)
    l.append(atemp)
    l.append(pH)
    l.append(rain)
    predictcrop = [l]

    predictions = loaded_rf.predict(predictcrop)
    count = 0
    for i in range(0, 30):
        if(predictions[0][i] == 1):
            c = crops[i]
            count = count+1
            break
        i = i+1
    if(count == 0):
        print('The predicted crop is %s' % cr)
        result = "rice"
    else:
        print('The predicted crop is %s' % c)
        result = c

    sys.stdout.flush()

    '''if prediction[0] == 0:
      result = "You are not at risk"
      print(result)
    else:
      result = "You are at Risk"
      print(result)
      
    sys.stdout.flush()'''

    # generate result to publish
    updateEntity = \
        {
            'id': ctxObj["id"]+".prediction",
            'type': 'CropPrediction',
            'soilmoisture': {'type': 'Property', 'value': shum},
            'soilph': {'type': 'Property', 'value': pH},
            'rainfall': {'type': 'Property', 'value': rain},
            'airmoisture': {'type': 'Property', 'value': ah},
            'cropprediction': {'type': 'Property', 'value': str(result)}
        }
    print("Update Entity : ")
    print(json.dumps(updateEntity))
    sys.stdout.flush()

    '''publishResultOnDesigner(updateEntity)

    if result == "You are at Risk":
      publish(updateEntity)
    sys.stdout.flush()'''

    '''ctxObjKeys = ctxObj.keys()
    
    for ctxEle in ctxObjKeys:
        if ctxEle != 'id' and ctxEle != 'type' and ctxEle != 'modifiedAt' \
            and ctxEle != 'createdAt' and ctxEle != 'observationSpace' \
            and ctxEle != 'operationSpace' and ctxEle\
            != '@context':
            ctxObjValue = ctxObj[ctxEle]
            if ctxObjValue.has_key('type') == True:
                if ctxObjValue['type'] == 'Relationship':
                    print(ctxEle,ctxObjValue['type'],ctxObjValue['object'])
                else:
                    print(ctxEle,ctxObjValue['type'],ctxObjValue['value'])'''
    publish(updateEntity)
    sys.stdout.flush()
