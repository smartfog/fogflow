import os,sys
sys.path.append('/root/GO/src/persistance/fogflow/test/UnitTest/Persistance')
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import data
import sys

# change it by broker ip and port
designerIp="http://180.179.214.208:8080"
brokerIp="http://180.179.214.208:8070"

'''
  test registration for opearator
'''
def test_persistOPerator():
	brokerUrl=brokerIp+"/ngsi10/entity/test0"
 	designerUrl=designerIp+"/ngsi10/updateContext"
	headers= {'Content-Type': 'application/json'}
	r=requests.post(designerUrl,data=json.dumps(data.test0),headers=headers)
	r=requests.get(url,headers=headers)
	assert r.status_code == 200
	
'''
  test registration for fogfunction
'''
def test_persistOPerator():
        brokerUrl=brokerIp+"/ngsi10/entity/test1"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test1),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200

'''
  test registration for dockerImage
'''
def test_persistOPerator():
        brokerUrl=brokerIp+"/ngsi10/entity/test2"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test2),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200

'''
  test registration for topology
'''
def test_persistOPerator():
        brokerUrl=brokerIp+"/ngsi10/entity/test3"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test3),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200

