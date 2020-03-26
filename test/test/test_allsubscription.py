import os,sys
sys.path.append('/home/ubuntu/fogflow/test/test')
from datetime import datetime
import copy 
import json
import requests
import time
import pytest
import data
import sys
# testCase 1
def test_getAllSubscription():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	r=requests.get(url,headers=headers)
	assert r.status_code == 200

#testCase 2
def test_getSubscriptionById():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	r=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	resp_content=r.content
	resInJson= resp_content.decode('utf8').replace("'", '"')
	resp=json.loads(resInJson)
	resp=resp['subscribeResponse']
	sid=resp['subscriptionId']
	get_url="http://192.168.100.120:8070/v2/subscription/"
	url=get_url+sid
	r=requests.get(url,headers=headers)
	assert r.status_code == 200
	print("Get subscription by Id testcase passed")

#testCase 2
def test_getSubscriptionByNilId():	
	get_url="http://192.168.100.120:8070/v2/subscription/nil"
	headers= {'Content-Type': 'application/json'}
	r=requests.get(get_url,headers=headers)
	assert r.status_code == 404	

# Test Delete subscription id subscriptionId is persent in the broker
#testCase 2
def test_deleteSubscriptionById():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	r=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	resp_content=r.content
	resInJson= resp_content.decode('utf8').replace("'", '"')
	resp=json.loads(resInJson)
	resp=resp['subscribeResponse']
	sid=resp['subscriptionId']
	get_url="http://192.168.100.120:8070/v2/subscription/"
	url=get_url+sid
	r=requests.delete(url,headers=headers)
	assert r.status_code == 200

#testCase 3
# Test if subscriptionId is not persent in the broker
def test_deleteSubscriptionId():	
	delete_url="http://192.168.100.120:8070/v2/subscription/nil"
	headers= {'Content-Type': 'application/json'}
	r=requests.delete(delete_url,headers=headers)
	assert r.status_code == 200

#testCase 4
def test_subscriptionWithWrongPayload():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	r=requests.post(url,data=json.dumps(data.subscriptionWrongPaylaod),headers=headers)
	assert r.status_code == 500

#testCase 5
def test_v1v2SubscriptionForSameEntity():
	V2url="http://192.168.100.120:8070/v2/subscriptions"
	ngsi10url="http://192.168.100.120:8070/ngsi10/subscribeContext"
	headers= {'Content-Type': 'application/json'}
	V2=requests.post(V2url,data=json.dumps(data.subscription_data),headers=headers)
	ngsi10=requests.post(ngsi10url,data=json.dumps(data.v1SubData),headers=headers)
	assert V2.status_code == 201
	assert ngsi10.status_code == 200
	
#testCase 6
def test_Subscription():	
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	response=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	assert response.status_code==201

#testCase 7
#update request wit create action 
def test_update_request_with_create_action():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	subresponse=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	updateresponse=requests.post(url,data=json.dumps(data.updateDataWithupdateaction),headers=headers)	
	assert updateresponse.status_code == 201
#testCase 8
def test_upadte_request_with_delete_action():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	subresponse=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	updateresponse=requests.post(url,data=json.dumps(data.deleteDataWithupdateaction),headers=headers)
	assert updateresponse.status_code==201
	
#testCase 9
def test_update_request_with_create_action():	
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	subresponse=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	updateresponse=requests.post(url,data=json.dumps(data.createDataWithupdateaction),headers=headers)
	assert updateresponse.status_code==201

#testCase 10
def test_notifyOneSubscriberv2WithCurrentStatus():
	url="http://192.168.100.120:8070/v2/subscriptions"
	headers= {'Content-Type': 'application/json'}
	updateresponse=requests.post(url,data=json.dumps(data.createDataWithupdateaction),headers=headers)
	subresponse=requests.post(url,data=json.dumps(data.subscription_data),headers=headers)
	assert subresponse.status_code==201
