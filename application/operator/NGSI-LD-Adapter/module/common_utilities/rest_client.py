import logging
from LogerHandler import Handler
from consts import constant
import requests
import io
import datetime
import pickle
import json
from flask import Flask, abort, request
import sys
import os
sys.path.append('/opt/ngsildAdapter/module')


class Rest_client:
    def __init__(self, url, payload):
        self.url = url
        self.payload = payload
        self.headers = constant.header
        logger_obj = Handler()
        self.logger = logger_obj.get_logger()

# sending post request
    def post_request(self):
        self.logger.info("Sending post request")
        response = requests.post(
            self.url, data=self.payload, headers=self.headers)
        if response.ok:
            self.logger.debug("post response is ok")
            return response
        else:
            self.logger.debug("Response is None Entity may already exits")
            return None

    # sending patch request
    def patch_request(self):
        self.logger.info("Patch request is sending")
        response = requests.patch(
            self.url, data=self.payload, headers=self.headers)
        if response.ok:
            self.logger.debug("Patch response is ok")
            return response
        else:
            self.logger.info(
                "Patch response is None Entity there may some problem in entity")
            return None
