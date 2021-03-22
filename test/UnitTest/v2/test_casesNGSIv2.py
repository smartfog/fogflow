import os,sys
# change the path accoring to the test folder in system
#sys.path.append('/home/ubuntu/setup/src/fogflow/test/UnitTest/v2')
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import v2data
import sys


# change it by broker ip and port
brokerIp="http://127.0.0.1:8070"

print(" The Validation test begins ")

# testCase 1
'''
  Testing get all subscription
'''
def test_getAllSubscription():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        r=requests.get(url,headers=headers)
        assert r.status_code == 200


#testCase 2

'''
 Testing get subscription by Id
'''

def test_getSubscriptionById():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        get_url=brokerIp+"/v2/subscription/"
        url=get_url+sid
        r=requests.get(url,headers=headers)
        assert r.status_code == 200
        print("Get subscription by Id testcase passed")


#testCase 3
'''
  Testing get subscription by nil id
'''

def test_getSubscriptionByNilId():
        get_url=brokerIp+"/v2/subscription/nil"
        headers= {'Content-Type': 'application/json'}
        r=requests.get(get_url,headers=headers)
        assert r.status_code == 404

# Test Delete subscription id subscriptionId is persent in the broker
#testCase 4

'''
  Test delete subscription by Id
'''

def test_deleteSubscriptionById():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        get_url=brokerIp+"/v2/subscription/"
        url=get_url+sid
        r=requests.delete(url,headers=headers)
        assert r.status_code == 200



#testCase 5
'''
Test if subscriptionId is not persent in the broker
'''

def test_deleteSubscriptionId():
        delete_url=brokerIp+"/v2/subscription/nil"
        headers= {'Content-Type': 'application/json'}
        r=requests.delete(delete_url,headers=headers)
        assert r.status_code == 200

#testCase 6
'''
  Test with wrong payload
'''

def test_subscriptionWithWrongPayload():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subscriptionWrongPaylaod),headers=headers)
        assert r.status_code == 500

#testCase 7
'''
  Testing subscription response with the same entity in ngsiv1 and ngsiv2.
'''

def test_v1v2SubscriptionForSameEntity():
        V2url=brokerIp+"/v2/subscriptions"
        ngsi10url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        V2=requests.post(V2url,data=json.dumps(v2data.subscription_data),headers=headers)
        ngsi10=requests.post(ngsi10url,data=json.dumps(v2data.v1SubData),headers=headers)
        assert V2.status_code == 201
        assert ngsi10.status_code == 200

#testCase 8
'''
 Test ngsiv2 subscription
'''

def test_Subscription():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        response=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        assert response.status_code==201

#testCase 9
#update request wit create action

'''
  Testing update request with update action
'''
def test_update_request_with_update_action():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        subresponse=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        updateresponse=requests.post(url,data=json.dumps(v2data.updateDataWithupdateaction),headers=headers)
        assert updateresponse.status_code == 201

#testCase 10
'''
  Testing update request with Delete request
'''

def test_upadte_request_with_delete_action():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        subresponse=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        updateresponse=requests.post(url,data=json.dumps(v2data.deleteDataWithupdateaction),headers=headers)
        assert updateresponse.status_code==201


#testCase 11

'''
  Testing update request with create action
'''

def test_update_request_with_create_action():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        subresponse=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        updateresponse=requests.post(url,data=json.dumps(v2data.createDataWithupdateaction),headers=headers)
        assert updateresponse.status_code==201

#testCase 12
'''
   Testing notification send by broker
'''

def test_notifyOneSubscriberv2WithCurrentStatus():
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        updateresponse=requests.post(url,data=json.dumps(v2data.createDataWithupdateaction),headers=headers)
        subresponse=requests.post(url,data=json.dumps(v2data.subscription_data),headers=headers)
        assert subresponse.status_code==201

# testCase 13
'''
  Testing  subscription with attributes and using ID
'''
def test_getSubscription1validation():
        #update request to create entity at broker
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subdata1),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #subsciption request in v2 format
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subdata2),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #update to trigger notification
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(v2data.subdata3),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #validation based on subscriptionId
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200


# testCase 14
'''
  Testing subscription with attributes and with Header : User-Agent
'''
def test_getsubscription2validate():
        #update request to create entity at broker
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subdata4),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #subsciption request in v2 format
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(v2data.subdata5),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #update to trigger notification
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(v2data.subdata6),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #vaidation based on subscriptionId
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200


# testCase 15
'''
  Testing with subscribing , updating and deleting and validating
'''
def test_getsubscription3():
        #update request to create entity at broker
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'application/json'}
        r=requests.post(url,data=json.dumps(v2data.subdata7),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #subsciption request in v2 format
        url=brokerIp+"/v2/subscriptions"
        headers= {'Content-Type': 'application/json','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(v2data.subdata8),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #update to trigger Notification
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(v2data.subdata9),headers=headers)
        #print(r.status_code)
        #print(r.content)

        #delete the subscription
        url=brokerIp+"/v2/subscription/"+sid
        r=requests.delete(url)
        #print(r.status_code)
        #print(r.content)
        print("The subscriptionId "+sid+" is deleted successfully")

        #validate based on get subscriptionId
        url=brokerIp+"/v2/subscription/"+sid
        r=requests.get(url)
        #print(r.status_code)
        #print(r.content)
        assert r.status_code == 404
        print("The subscriptionId "+sid+" coud not be fetched via get since deleted")

