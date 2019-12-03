from flask import Flask, abort, request
import json
app = Flask(__name__)

myStatus = 'off'

# Receive commands from Broker
@app.route('/ngsi10/updateContext',methods=['POST'])
def getUpdateNotification():
        #dir(request)
        data=request.get_json()
        if data.has_key('contextElements')==True:
                dataContext=data['contextElements']
                dataAttribute=dataContext[0]
                if dataAttribute.has_key('attributes')==True:
                        attribute=dataAttribute['attributes']
                        if attribute[0]['name']=='off':
                                my_status='off'
                        elif attribute[0]['name']=='on':
                                my_status='on'
                        else:
                                print("Command not found!!")
                                return ""
        print('Lamp : {}'.format(my_status))
        return ""
# main file for starting application
if __name__ == '__main__':
    app.run(host= '0.0.0.0', port=8888, debug=True)
