#!/usr/bin/python
from flask import Flask, abort, request
import json
import pickle
import datetime
import io
import requests
import copy
import constant
from data_creation import ngsi_data_creation
from request_file import Rest_client
app = Flask(__name__)
File_data={}
storage=[]


def take_backup():
    id_file=open("data_file.txt",'r+')
    for x in id_file:
        x=x.rstrip("\n")
        File_data[x]=1
    print(File_data)
    id_file.close()
#notify context
@app.route('/notifyContext',methods=['POST'])
def noify_server():
    data=request.get_json()
    print(data)
    dataObj=ngsi_data_creation(data)
    #print(data)
    context=dataObj.get_ngsi_ld()
    #print(context)
    patch_context= copy.deepcopy(context)
    del patch_context['type']
    del patch_context['id']
    entity_id=dataObj.get_entityId()
    url1 = constant.entity_url
    url2=constant.entity_url+entity_id+'/attrs'
    if entity_id in File_data.keys():
        payload=json.dumps(patch_context)
        robj=Rest_client(url2,payload)
        r=robj.patch_request()
    else:
        payload=json.dumps(context)
        #print(payload)
        robj=Rest_client(url1,payload)
        r=robj.post_request()
        if r.status_code==201:
            id_file=open("data_file.txt",'a+')
            id_file.write(entity_id+'\n')
            File_data[entity_id]=1
            id_file.close()
    return "notify"
@app.route('/subscribeContext',methods=['POST'])
def rest_client():
    data=request.get_json()
    url=constant.s_url
    payload = json.dumps(data)
    robj=Rest_client(url,payload)
    r=robj.post_request()
    return "subscribe"
if __name__ == '__main__':
    take_backup()
    app.run(host= '0.0.0.0', port=8888, debug=True)

