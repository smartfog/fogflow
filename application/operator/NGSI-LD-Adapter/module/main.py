#!/usr/bin/python
from flask import Flask, abort, request
import json
import pickle
import datetime
import io
import ConfigParser
from common_utilities.config import config_data
from common_utilities.rest_client import Rest_client
from common_utilities.LogerHandler import Handler
from config.check_config import check_ip_port
import requests
from consts import constant
from data_model.ld_generate import ngsi_data_creation
from data_model.orian_ld_genrate import orian_convert_data
from data_model.orian_ld_genrate import orian_convert_data
import copy
import logging
#from config import check_ip_port
from consts import constant

app = Flask(__name__)
File_data = {}
storage = []


class patch_post:
    def __init__(self):
        logger_obj = Handler()
        self.logger = logger_obj.get_logger()

    # storing Entity_id in the file

    def take_backup(self):
        self.logger.info("take_backup function has been started")
        try:
            id_file = open("data_model/storage/data_file.txt", 'r+')
        except FileNotFoundError as fnf_error:
            message = fnf_error
            self.logger.error(message)
        self.logger.debug("storing created Entity from file")
        for x in id_file:
            x = x.rstrip("\n")
            File_data[x] = 1
        self.logger.debug("Clossing the file")
        id_file.close()
        self.logger.info("take_backup function has been end")
    # converting data to ngsild and sending request

    def data_convert_invoker(self, context, dataObj):
        self.logger.info("data_convert_invoker has been started")
        patch_context = copy.deepcopy(context)
        del patch_context["type"]
        del patch_context["id"]
        entity_id = dataObj.get_entityId()
        configobj = config_data()
        entity_url = configobj.get_entity_url()
        url1 = constant.http+entity_url+constant.entity_uri
        url2 = constant.http+entity_url+constant.entity_uri+entity_id+'/attrs'
        if entity_id in File_data.keys():
            self.logger.debug("sending update request")
            payload = json.dumps(patch_context)
            ngb_payload = 'NGB payload:'+str(payload)
            print(ngb_payload)
            self.logger.info(ngb_payload)
            robj = Rest_client(url2, payload)
            r = robj.patch_request()
            ngb_status = 'status code from NGB:'+str(r.status_code)
            print(ngb_status)
            self.logger.info(ngb_status)
            if r.status_code == constant.update_status:
                self.logger.debug("Entity has been updated to NGB")
        else:
            self.logger.debug("Sending create request")
            payload = json.dumps(context)
            ngb_payload = 'NGB payload:'+str(payload)
            print(ngb_payload)
            self.logger.info(ngb_payload)
            robj = Rest_client(url1, payload)
            r = robj.post_request()
            ngb_status = 'status code from NGB:'+str(r.status_code)
            print(ngb_status)
            if r.status_code == constant.create_status:
                self.logger.debug("Entity has been created in NGB")
                id_file = open("data_model/storage/data_file.txt", 'a+')
                id_file.write(entity_id+'\n')
                File_data[entity_id] = 1
                id_file.close()
        self.logger.info("data_convert_invoker has been end")
 # notify app


@app.route('/notifyContext', methods=['POST'])
def noify_server():
    data = request.get_json()
    print('Data from fogflow : '+str(data))
    logger_obj = Handler()
    logger = logger_obj.get_logger()
    logger.debug("notify_server has been started")
    message = 'notify data'+str(data)
    logger.debug(message)
    dataObj = ngsi_data_creation(data)
    context = dataObj.get_ngsi_ld()
    logger.debug("Data is converted to ngsi-ld")
    obj = patch_post()
    obj.data_convert_invoker(context, dataObj)
    logger.info("noify_server has been end")
    return "notify"


@app.route('/notifyContext1', methods=['POST'])
def notify_server_for_orian():
    logger_obj = Handler()
    logger = logger_obj.get_logger()
    data = request.get_json()
    message = 'notify data'+str(data)
    logger.info(message)
    dataObj = orian_convert_data(data)
    context = dataObj.get_data()
    logger.info("Data is converted to ngsi-ld")
    obj = patch_post()
    obj.data_convert_invoker(context, dataObj)
    return "notify"

# subscribe request to the fogflow


@app.route('/subscribeContext', methods=['POST'])
def rest_client():
    logger_obj = Handler()
    logger = logger_obj.get_logger()
    data = request.get_json()
    configobj = config_data()
    fog_url = configobj.get_fogflow_subscription_endpoint()
    url = constant.http+fog_url+constant.subscribe_uri
    payload = json.dumps(data)
    robj = Rest_client(url, payload)
    r = robj.post_request()
    logger.info("Forwarded Subscription to Fogflow.")
    return "subscribe"

# main function


if __name__ == '__main__':
    obj = patch_post()
    obj.take_backup()
    print("Running at 0.0.0.0:8888")
    app.run(host='0.0.0.0', port=8888, debug=True)
