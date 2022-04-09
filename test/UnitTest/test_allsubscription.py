import copy
import sys
import data
import pytest
import time
import requests
import json
from datetime import datetime
import os
import sys
# change the path accoring to the test folder in system
sys.path.append('/home/ubuntu/fogflow/test/UnitTest')

# change it by broker ip and port
brokerIp = "http://192.168.100.120:8070"

# testCase 1
'''
  Testing get all subscription
'''


def test_getAllSubscription():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    r = requests.get(url, headers=headers)
    assert r.status_code == 200

# testCase 2


'''
 Testing get subscription by Id
'''


def test_getSubscriptionById():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    resp = resp['subscribeResponse']
    sid = resp['subscriptionId']
    get_url = brokerIp+"/v2/subscription/"
    url = get_url+sid
    r = requests.get(url, headers=headers)
    assert r.status_code == 200
    print("Get subscription by Id testcase passed")


# testCase 2
'''
  Testing get subscription by nil id
'''


def test_getSubscriptionByNilId():
    get_url = brokerIp+"/v2/subscription/nil"
    headers = {'Content-Type': 'application/json'}
    r = requests.get(get_url, headers=headers)
    assert r.status_code == 404

# Test Delete subscription id subscriptionId is persent in the broker
# testCase 2


'''
  Test delete subscription by Id 
'''


def test_deleteSubscriptionById():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    resp = resp['subscribeResponse']
    sid = resp['subscriptionId']
    get_url = brokerIp+"/v2/subscription/"
    url = get_url+sid
    r = requests.delete(url, headers=headers)
    assert r.status_code == 200


# testCase 3
''' 
Test if subscriptionId is not persent in the broker
'''


def test_deleteSubscriptionId():
    delete_url = brokerIp+"/v2/subscription/nil"
    headers = {'Content-Type': 'application/json'}
    r = requests.delete(delete_url, headers=headers)
    assert r.status_code == 200


# testCase 4
'''
  Test with wrong payload
'''


def test_subscriptionWithWrongPayload():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, data=json.dumps(
        data.subscriptionWrongPaylaod), headers=headers)
    assert r.status_code == 500


# testCase 5
'''
  Testing subscription response with the same entity in ngsiv1 and ngsiv2.
'''


def test_v1v2SubscriptionForSameEntity():
    V2url = brokerIp+"/v2/subscriptions"
    ngsi10url = brokerIp+"/ngsi10/subscribeContext"
    headers = {'Content-Type': 'application/json'}
    V2 = requests.post(V2url, data=json.dumps(
        data.subscription_data), headers=headers)
    ngsi10 = requests.post(ngsi10url, data=json.dumps(
        data.v1SubData), headers=headers)
    assert V2.status_code == 201
    assert ngsi10.status_code == 200


# testCase 6
'''
 Test ngsiv2 subscription
'''


def test_Subscription():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    assert response.status_code == 201

# testCase 7
# update request wit create action


'''
  Testing update request with update action
'''


def test_update_request_with_update_action():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    subresponse = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    updateresponse = requests.post(url, data=json.dumps(
        data.updateDataWithupdateaction), headers=headers)
    assert updateresponse.status_code == 201


# testCase 8
'''
  Testing update request with Delete request
'''


def test_upadte_request_with_delete_action():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    subresponse = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    updateresponse = requests.post(url, data=json.dumps(
        data.deleteDataWithupdateaction), headers=headers)
    assert updateresponse.status_code == 201

# testCase 9


'''
  Testing update request with create action
'''


def test_update_request_with_create_action():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    subresponse = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    updateresponse = requests.post(url, data=json.dumps(
        data.createDataWithupdateaction), headers=headers)
    assert updateresponse.status_code == 201


# testCase 10
'''
   Testing notification send by broker 
'''


def test_notifyOneSubscriberv2WithCurrentStatus():
    url = brokerIp+"/v2/subscriptions"
    headers = {'Content-Type': 'application/json'}
    updateresponse = requests.post(url, data=json.dumps(
        data.createDataWithupdateaction), headers=headers)
    subresponse = requests.post(url, data=json.dumps(
        data.subscription_data), headers=headers)
    assert subresponse.status_code == 201
