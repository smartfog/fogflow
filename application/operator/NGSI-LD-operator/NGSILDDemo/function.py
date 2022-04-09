def handleEntity(ctxObj, publish):
    print('===============Implement losic====================')

    print(ctxObj)

    ctxObjKeys = ctxObj.keys()

    for ctxEle in ctxObjKeys:
        if ctxEle != 'id' and ctxEle != 'type' and ctxEle != 'modifiedAt' \
                and ctxEle != 'createdAt' and ctxEle != 'observationSpace' \
                and ctxEle != 'operationSpace' and ctxEle\
                != '@context':
            ctxObjValue = ctxObj[ctxEle]
            if ctxObjValue.has_key('type') == True:
                if ctxObjValue['type'] == 'Relationship':
                    print(ctxEle, ctxObjValue['type'], ctxObjValue['object'])
                else:
                    print(ctxEle, ctxObjValue['type'], ctxObjValue['value'])
    publish(ctxObj)
