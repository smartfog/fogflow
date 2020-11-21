import os
import copy
import json
import requests
import time
import pytest
import data
import sys

# change it by broker ip and port
designerIp="http://localhost:8080"
brokerIp= "http://localhost:8070"

'''
  test registration for opearator
'''


def test_persistOPerator():
    brokerUrl = brokerIp + "/ngsi10/entity/test011"
    designerUrl = designerIp + "/intent/updateContext"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(
        designerUrl,
        data=json.dumps(
            data.test0),
        headers=headers)
    r = requests.get(brokerUrl, headers=headers)
    # print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "Operator" and resp["entityId"]["id"] == "test011":
        print "\nValidated"
    else:
        print "\nNot Validated"
    assert r.status_code == 200


'''
  test registration for fogfunction
'''


def test_persistFogFunction():
    brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    designerUrl = designerIp + "/intent/updateContext"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(
        designerUrl,
        data=json.dumps(
            data.test1),
        headers=headers)
    r = requests.get(brokerUrl, headers=headers)
    # print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "FogFunction" and resp["entityId"]["id"] == "test2":
        print "\nValidated"
    else:
        print "\nNot Validated"
    assert r.status_code == 200


'''
  test registration for dockerImage
'''


def test_persistDockerImage():
    brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    designerUrl = designerIp + "/intent/updateContext"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(
        designerUrl,
        data=json.dumps(
            data.test2),
        headers=headers)
    r = requests.get(brokerUrl, headers=headers)
    # print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "DockerImage" and resp["entityId"]["id"] == "test3":
        print "\nValidated"
    else:
        print "\nNot Validated"
    assert r.status_code == 200


'''
  test registration for topology
'''


def test_persistopology():
    brokerUrl = brokerIp + "/ngsi10/entity/test4"
    designerUrl = designerIp + "/intent/updateContext"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(
        designerUrl,
        data=json.dumps(
            data.test3),
        headers=headers)
    r = requests.get(brokerUrl, headers=headers)
    # print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "Topology" and resp["entityId"]["id"] == "test4":
        print "\nValidated"
    else:
        print "\nNot Validated"
    assert r.status_code == 200


'''
  test if entity does not have domainMetaData
'''
'''def test_DomainMetaDataMissing():
        brokerUrl=brokerIp+"/ngsi10/entity/test4"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test4),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''
'''
   testCase  if entity does not have attribute
'''
'''def test_attributesMissing():
        brokerUrl=brokerIp+"/ngsi10/entity/test5"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test5),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''
'''
   test if type of attributes is string
'''
'''def test_stringAttributes():
        brokerUrl=brokerIp+"/ngsi10/entity/test6"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test6),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''
'''
  test if value and type of attributes is null
'''
'''def test_nullAttributes():
        brokerUrl=brokerIp+"/ngsi10/entity/test7"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test7),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''

'''
  test if type of domainMetaData is point
'''
'''def test_pointDomainMetaData():
        brokerUrl=brokerIp+"/ngsi10/entity/test8"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test8),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''

'''
   test if data have contextElement(test for curl client)
'''

'''def test_forCurlClient():
        brokerUrl=brokerIp+"/ngsi10/entity/test9"
        designerUrl=designerIp+"/ngsi10/updateContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(designerUrl,data=json.dumps(data.test9),headers=headers)
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
'''
