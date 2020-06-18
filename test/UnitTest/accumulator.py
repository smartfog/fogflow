
from flask import Flask, abort, request
import requests
import json
app = Flask(__name__)

myStatus = 'off'

subId=[]


@app.route('/accumulate',methods=['POST'])
def getUpdateNotification():
        data=request.get_json()
        print(data)
        pload=data["subscriptionId"]
        subId.append(data["subscriptionId"])
        print(pload)
        return "Done"

@app.route('/v2/notifyContext',methods=['POST'])
def getUpdateNotificatio1n():
        data=request.get_json()
        print(data)
        print(data["subscriptionId"])
        subId.append(data["subscriptionId"])
        print(subId)
        return "Done."

@app.route('/validateNotification',methods=['POST'])
def getValidationNotification():
        data=request.get_json(force=True)
        print(data)
        pload=data["subscriptionId"]
        print(pload)
        print(subId)
        if pload in subId:
                print("validated the method")
                return "Validated"
        else:
                print("Not validated")
                return "Not validated"


# main file for starting application
if __name__ == '__main__':
    app.run(host= '0.0.0.0', port=8888, debug=True)
