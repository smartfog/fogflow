from flask import Flask, abort, request
import json
import pickle
import datetime
import io
import requests
from consts import constant
class Rest_client:
    def __init__(self,url,payload):
        self.url=url
        self.payload=payload
        self.headers=constant.header
    def post_request(self):
        response = requests.post(self.url, data=self.payload, headers=self.headers)
        print(dir(response))
        if response.ok: 
            return response
        else:
            return None 
    def patch_request(self):
        response = requests.patch(self.url, data=self.payload, headers=self.headers)
        if response.ok:
           return response
        else:
            return None
            
         
