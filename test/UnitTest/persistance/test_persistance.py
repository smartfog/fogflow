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
#brokerIp= "http://localhost:8070"

'''
  testcase 1: To test registration for opearator
'''

def test_persistOPerator():
    #brokerUrl = brokerIp + "/ngsi10/entity/test011"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test0),headers=headers)
    '''
    #print(r.content)
    r = requests.get(brokerUrl, headers=headers)
    #print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "Operator" and resp["entityId"]["id"] == "test011":
        print "\nValidated"
    else:
        print "\nNot Validated"
    '''
    assert r.status_code == 200


'''
  testcase 2: To test registration for opearator with empty payload
'''

def test_persistOPerator1():
    #brokerUrl = brokerIp + "/ngsi10/entity/test011"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test001),headers=headers)

    assert r.status_code == 200

'''
  testcase 3: To test registration for opearator with operator name only.
'''

def test_persistOPerator2():
    #brokerUrl = brokerIp + "/ngsi10/entity/test011"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test002),headers=headers)

    assert r.status_code == 200

'''
  testcase 4: To test registration for opearator with operator name and description.
'''

def test_persistOPerator3():
    #brokerUrl = brokerIp + "/ngsi10/entity/test011"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test003),headers=headers)

    assert r.status_code == 200



'''
  testcase 5: To test registration for fogfunction
'''


def test_persistFogFunction():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test1),headers=headers)

    '''
    #print(r.content)
    r = requests.get(brokerUrl, headers=headers)
    #print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "FogFunction" and resp["entityId"]["id"] == "test2":
        print "\nValidated"
    else:
        print "\nNot Validated"
    '''
    assert r.status_code == 200

'''
  testcase 6: To test registration for fogfunction with empty payload
'''

def test_persistFogFunction1():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test101),headers=headers)
    assert r.status_code == 200

'''
  testcase 7: To test registration for fogfunction with id only.
'''

def test_persistFogFunction2():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test102),headers=headers)
    assert r.status_code == 200

'''
  testcase 8: To test registration for fogfunction with attributes: name and topology.
'''

def test_persistFogFunction3():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test103),headers=headers)
    assert r.status_code == 200

'''
  testcase 9: To test registration for fogfunction with only one attribute: geoscope.
'''

def test_persistFogFunction4():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test2"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test104),headers=headers)
    assert r.status_code == 200


'''
  testcase 10: To test registration for dockerImage
'''

def test_persistDockerImage():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test2),headers=headers)
    '''
    #print(r.content)
    r = requests.get(brokerUrl, headers=headers)
    #print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "DockerImage" and resp["entityId"]["id"] == "test3":
        print "\nValidated"
    else:
        print "\nNot Validated"
    '''
    assert r.status_code == 200

'''
  testcase 11: To test registration for dockerImage with empty payload
'''

def test_persistDockerImage1():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test200),headers=headers)
    assert r.status_code == 200

'''
  testcase 12: To test registration for dockerImage with operator only.
'''

def test_persistDockerImage2():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test201),headers=headers)
    assert r.status_code == 200

'''
  testcase 13: To test registration for dockerImage with operator and name.
'''

def test_persistDockerImage3():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test202),headers=headers)
    assert r.status_code == 200

'''
  testcase 14: To test registration for dockerImage with attributes: hwType and osType
'''

def test_persistDockerImage4():
    #brokerUrl = brokerIp+ "/ngsi10/entity/test3"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test203),headers=headers)
    assert r.status_code == 200


'''
  testcase 15: To test registration for topology
'''

def test_persistopology():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test3),headers=headers)

    '''
    #print(r.content)
    r = requests.get(brokerUrl, headers=headers)
    #print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "Topology" and resp["entityId"]["id"] == "test4":
        print "\nValidated"
    else:
        print "\nNot Validated"
    '''
    assert r.status_code == 200

'''
  testcase 16: To test registration for topology with empty payload.
'''

def test_persistopology1():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test300),headers=headers)
    assert r.status_code == 200

'''
  testcase 17: To test registration for topology with name and description.
'''

def test_persistopology2():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test301),headers=headers)
    assert r.status_code == 200

'''
testcase 18: To test registration for service intent
'''

def test_persistintent():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test4),headers=headers)

    '''
    #print(r.content)
    r = requests.get(brokerUrl, headers=headers)
    #print(r.content)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    if resp["entityId"]["type"] == "Topology" and resp["entityId"]["id"] == "test4":
        print "\nValidated"
    else:
        print "\nNot Validated"
    '''
    assert r.status_code == 200

'''
testcase 19: To test registration for service intent with empty payload
'''

def test_persistintent1():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test400),headers=headers)
    assert r.status_code == 200

'''
testcase 20: To test registration for service intent with attributes: id and topology.
'''

def test_persistintent2():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test401),headers=headers)
    assert r.status_code == 200

'''
testcase 21: To test registration for service intent with attributes: geoscope and topology.
'''

def test_persistintent3():
    #brokerUrl = brokerIp + "/ngsi10/entity/test4"
    #designerUrl = designerIp + "/internal/updateContext"
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test402),headers=headers)
    assert r.status_code == 200


'''
'''
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

