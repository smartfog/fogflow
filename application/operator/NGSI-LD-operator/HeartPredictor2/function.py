import joblib
import pandas as pd
import sys
import json

loaded_rf = joblib.load("./predictor.joblib")

def handleEntity(ctxObj, publish, publishResultOnDesigner):
    print('===============Implement logic====================')
    
    print(ctxObj)
    sys.stdout.flush()
    print(ctxObj["type"])
    sys.stdout.flush()   
    print(ctxObj["age"]["value"])
    sys.stdout.flush() 

    my_data = {
    'age': [ctxObj["age"]["value"]],
    'sex': [ctxObj["sex"]["value"]],
    'cp': [ctxObj["cp"]["value"]],
    'trestbps': [ctxObj["trestbps"]["value"]],
    'chol': [ctxObj["chol"]["value"]],
    'fbs': [ctxObj["fbs"]["value"]],
    'restecg': [ctxObj["restecg"]["value"]],
    'thalach': [ctxObj["thalach"]["value"]],
    'exang': [ctxObj["exang"]["value"]],
    'oldpeak': [ctxObj["oldpeak"]["value"]],
    'slope': [ctxObj["slope"]["value"]],
    'ca': [ctxObj["ca"]["value"]],
    'thal': [ctxObj["thal"]["value"]]
    }

    myvar = pd.DataFrame(my_data)

    prediction = loaded_rf.predict(myvar)
    print("The prediction is ",prediction)
    sys.stdout.flush()
 
    if prediction[0] == 0:
      result = "You are not at risk"
      print(result)
    else:
      result = "You are at Risk"
      print(result)
      
    sys.stdout.flush()

    #generate result to publish 
    updateEntity = \
    {   
        'id' : ctxObj["id"]+".prediction",
        'type' : 'prediction',
        'Analysis' : { 'type' : 'Property', 'value' : result}
    }
    print("Update Entity : ")
    print(json.dumps(updateEntity))
    sys.stdout.flush()

    publishResultOnDesigner(updateEntity)

    if result == "You are at Risk":
      publish(updateEntity)
    sys.stdout.flush()


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
                    print(ctxEle,ctxObjValue['type'],ctxObjValue['value'])
    publish(ctxObj)'''

    
