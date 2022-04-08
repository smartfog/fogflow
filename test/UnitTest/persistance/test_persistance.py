import os
import copy
import json
import requests
import time
import pytest
import data
import sys

designerIp="http://localhost:8080"

'''
  testcase 1: To test registration for opearator
'''

def test_persistOPerator():
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test0), headers=headers)
    assert r.status_code == 200

'''
  testcase 2: To test registration for opearator with empty payload
'''

def test_persistOPerator1():
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test001),headers=headers)

    assert r.status_code == 200

'''
  testcase 3: To test registration for opearator with operator name only.
'''

def test_persistOPerator2():
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test002),headers=headers)

    assert r.status_code == 200

'''
  testcase 4: To test registration for opearator with operator name and description.
'''

def test_persistOPerator3():
    designerUrl = designerIp + "/operator"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test003),headers=headers)

    assert r.status_code == 200



'''
  testcase 5: To test registration for fogfunction
'''


def test_persistFogFunction():
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test1),headers=headers)

    assert r.status_code == 200

'''
  testcase 6: To test registration for fogfunction with empty payload
'''

def test_persistFogFunction1():
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test101),headers=headers)
    assert r.status_code == 200

'''
  testcase 7: To test registration for fogfunction with id only.
'''

def test_persistFogFunction2():
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test102),headers=headers)
    assert r.status_code == 200

'''
  testcase 8: To test registration for fogfunction with attributes: name and topology.
'''

def test_persistFogFunction3():
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test103),headers=headers)
    assert r.status_code == 200

'''
  testcase 9: To test registration for fogfunction with only one attribute: geoscope.
'''

def test_persistFogFunction4():
    designerUrl = designerIp + "/fogfunction"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test104),headers=headers)
    assert r.status_code == 200


'''
  testcase 10: To test registration for dockerImage
'''

def test_persistDockerImage():
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test2),headers=headers)
    assert r.status_code == 200

'''
  testcase 11: To test registration for dockerImage with empty payload
'''

def test_persistDockerImage1():
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test200),headers=headers)
    assert r.status_code == 200

'''
  testcase 12: To test registration for dockerImage with operator only.
'''

def test_persistDockerImage2():
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test201),headers=headers)
    assert r.status_code == 200

'''
  testcase 13: To test registration for dockerImage with operator and name.
'''

def test_persistDockerImage3():
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test202),headers=headers)
    assert r.status_code == 200

'''
  testcase 14: To test registration for dockerImage with attributes: hwType and osType
'''

def test_persistDockerImage4():
    designerUrl = designerIp + "/dockerimage"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test203),headers=headers)
    assert r.status_code == 200


'''
  testcase 15: To test registration for topology
'''

def test_persistopology():
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test3),headers=headers)
    
    assert r.status_code == 200

'''
  testcase 16: To test registration for topology with empty payload.
'''

def test_persistopology1():
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test300),headers=headers)
    assert r.status_code == 200

'''
  testcase 17: To test registration for topology with name and description.
'''

def test_persistopology2():
    designerUrl = designerIp + "/service"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test301),headers=headers)
    assert r.status_code == 200

'''
testcase 18: To test registration for service intent
'''

def test_persistintent():
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test4),headers=headers)

    assert r.status_code == 200

'''
testcase 19: To test registration for service intent with empty payload
'''

def test_persistintent1():
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test400),headers=headers)
    assert r.status_code == 200

'''
testcase 20: To test registration for service intent with attributes: id and topology.
'''

def test_persistintent2():
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test401),headers=headers)
    assert r.status_code == 200

'''
testcase 21: To test registration for service intent with attributes: geoscope and topology.
'''

def test_persistintent3():
    designerUrl = designerIp + "/intent"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(designerUrl,data=json.dumps(data.test402),headers=headers)
    assert r.status_code == 200

