from flask import Flask, abort, request
import json
import pickle
import datetime
import io
import requests
import constant
class Rest_client:
    def __init__(self,url,payload):
        self.url=url
        self.payload=payload
        self.headers=constant.header
    def post_request(self):
	print("POST URL: ", self.url)
        print("Requested Payload: ", self.payload)
        r = requests.post(self.url, data=self.payload, headers=self.headers)
        print(r.content)
        print(r.status_code)
        return r
    def patch_request(self):
	print("PATCH URL: ", self.url)
	print("Requested Payload: ", self.payload)
        r = requests.patch(self.url, data=self.payload, headers=self.headers)
        print(r.content)
        print(r.status_code)
        return r
         
