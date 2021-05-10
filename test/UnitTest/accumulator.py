
from flask import Flask, abort, request
import requests
import json
app = Flask(__name__)

myStatus = 'off'

subId = []
entityIdDict = {}
# Getting notification for Quantumleap and sending response 200 to the
# test module

@app.route('/ngsi-ld/v1/entityOperations/upsert', methods=['POST'])
def upsertNotification():
    global entityIdDict
    print(dir(request))
    entities = request.get_json()
    entity = entities[0]
    id = entity["id"]
    print("id")
    print(id)
    entityIdDict[id] = 1
    return "Done"


@app.route('/ngsi-ld/v1/entities/urn:ngsi-ld:Device:water001/attrs/on', methods=['PATCH'])
def upsertNotificationNew():
    entities = request.get_json()
    print(dir(request))
    print(entities)
    return "Done"

@app.route('/validateupsert', methods=['POST'])
def upsertNotificationvalidator():
    global entityIdDict
    print(entityIdDict)
    if entityIdDict['urn:ngsi-ld:Device:water080'] == 1:
         statusCode = "200"
    else :	
	 statusCode = "404"
    return statusCode

@app.route('/accumulate', methods=['POST'])
def getUpdateNotification():
    print(dir(request))
    data = request.get_json()
    print(data)
    pload = data["subscriptionId"]
    subId.append(data["subscriptionId"])
    print(pload)
    return "Done"


@app.route('/csource', methods=['POST'])
def getNotifiedLD_csource():
    data = request.get_json()
    print(data)
    return "Done"


@app.route('/ld-notify', methods=['POST'])
def getNotifiedLD_subscription():
    data = request.get_json()
    print(data)
    pload = data["subscriptionId"]
    subId.append(data["subscriptionId"])
    print(pload)
    print(subId)
    return "Done"


@app.route('/v2/notifyContext', methods=['POST'])
def getUpdateNotificatio1n():
    # dir(request)
    data = request.get_json()
    print(data)
    print(data["subscriptionId"])
    subId.append(data["subscriptionId"])
    print(subId)
    return "Done."


@app.route('/validateNotification', methods=['POST'])
def getValidationNotification():
    data = request.get_json(force=True)
    print(data)
    pload = data["subscriptionId"]
    print(pload)
    print(subId)
    if pload in subId:
        print("validated the method")
        return "Validated"
    else:
        print("Not validated")
        return "Not validated"

@app.route('/ngsi10/updateContext',methods=['POST'])
def getNotified_southbound():
    data = request.get_json()
    print(data)
    if data["contextElements"][0]["attributes"][0]["type"]=="command" and data["contextElements"][0]["attributes"][0]["name"]=="on":
    	print("validated the method")
    else:
	print("Not validated")
    #pload = data["subscriptionId"]
    #subId.append(data["subscriptionId"])
    #print(pload)
    #print(subId)
    return "Notification of command recieved"



# main file for starting application
if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8888, debug=True)
