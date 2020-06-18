import os,sys
# change the path accoring to the test folder in system
#sys.path.append('/home/ubuntu/setup/src/fogflow/test/UnitTest/v1')
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import data_ngsi10
import sys

# change it by broker ip and port
brokerIp="http://localhost:8070"

print("Testing of v1 API")
# testCase 1
'''
  Testing post subscription
'''
def test_getSubscription1():
	url=brokerIp+"/ngsi10/subscribeContext"
	headers={'Content-Type' : 'application/json'}
	r=requests.post(url,data=json.dumps(data_ngsi10.subdata1),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
	#print(r.status_code)
	assert r.status_code == 200

#testCase 2
'''
  Testing entity creation with attributes, then susbscribing and get subscription using ID
'''
def test_getSubscription2():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
	r = requests.post(url,data=json.dumps(data_ngsi10.subdata2),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata3),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
        
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        #print(url)
        r=requests.get(url)
        assert r.status_code == 200
       
#testCase 3
'''
 Testing entity creation with one  attribute : pressure only followed by subscribing and get using ID
'''
def test_getSubscription3():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
	r = requests.post(url,data=json.dumps(data_ngsi10.subdata4),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata5),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
        
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        
#testCase 4
'''
  Testing entity creation with one attribute : Temperature only followed by subscription and get using ID
'''
def test_getSubscription4():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'} 
	r = requests.post(url,data=json.dumps(data_ngsi10.subdata6),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata7),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
	#print(sid)
        
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        
#testCase 5
'''
   Testing create entity without passing Domain data followed by subscription and get using ID
'''
def test_getSubscription5():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'} 
	r = requests.post(url,data=json.dumps(data_ngsi10.subdata8),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata9),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
        
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        
#testCase 6
'''
   Testing create entity without attributes followed by subscription and get using Id
'''
def test_getSubscription6():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}	
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata10),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata11),headers=headers)
        resp_content=r.content
	resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
       
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        
#testCase 7
'''
   Testing create entity without attributes and Metadata and followed by sbscription and get using Id
'''
def test_getSubscription7():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
	r = requests.post(url,data=json.dumps(data_ngsi10.subdata12),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	
	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata13),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
        
	#get request to fetch subscription
	get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        
#testCase 8
'''
   Testing create entity without entity type followed by subscription and get using Id
'''
def test_getSubscription8():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}	
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata14),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata15),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#get request to fetch subscription
        get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 200
        

#testCase 9
'''
  Testing get subscription request by first posting subscription request followed by delete request
'''
def test_getSubscription9():
	#create an entity
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata16),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #subscribing
	url_del=brokerIp+"/ngsi10/subscription/"
        url=url_del+sid
	r = requests.delete(url,headers=headers)
        #print(r.status_code)
	
	#get request to fetch subscription
        get_url=brokerIp+"/ngsi10/subscription"
        url=get_url+sid
        r=requests.get(url)
        assert r.status_code == 404
        
       
#testCase 10 
'''
  Testing the update post request 
'''
def test_getSubscription10():
	url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
	r=requests.post(url,data=json.dumps(data_ngsi10.subdata17),headers=headers)
	resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
	#print(r.status_code)
	assert r.status_code == 200
        
#testCase 11
'''
  Testing  subscription with attributes and using ID to validate it
'''
def test_getSubscription11():
	#create an entity
	url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata18),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
	#print(resp)
	
	#subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata19),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update the created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata20),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	# validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 12
'''
  Testing subscription for its if and else part : 1) for Destination Header
'''
def test_getSubscription12():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
	r=requests.post(url,data=json.dumps(data_ngsi10.subdata21),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
	#print(resp)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination' : 'orion-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata22),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update the created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata23),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 13
'''
  Tseting subscription for its if and else part : 2) for User - Agent Header
'''
def test_getSubscription18():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata24),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','User-Agent' : 'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata25),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata26),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 14
'''
  Testing subcription for its if else part : 3) Require-Reliability Header
'''
def test_getSubscription19():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata27),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Require-Reliability' : 'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata28),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
	resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #update the created entity
	url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata29),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 15
'''
  Testing subscription with two headers simultaneously : 4) Destination and User-Agent
'''
def test_getSubscription20():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata30),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination' : 'orion-broker','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata31),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata32),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)
	
	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 16
'''
  Testing subscription with two headers simultaneously : 4)  User-Agent and Require-Reliability
'''
def test_getSubscription21():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata33),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','User-Agent':'lightweight-iot-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata34),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata35),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        

#testCase 17
'''
   Testing subscription with two headers simultaneously : 4)  Destination  and Require-Reliability headers
'''
def test_getSubscription22():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata36),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination':'orion-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata37),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata38),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 18
'''
  Testing subscription with all  headers simultaneously : 5)  Destination, User-Agent  and Require-Reliability headers
'''
def test_getSubscription23():
	#create an entity
        url=brokerIp+"/ngsi10/updateContext"
	headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata39),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
	#print(resp1)

	#subscribing
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination':'orion-broker','User-Agent':'lightweight-iot-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata40),headers=headers)
        #print(r.status_code)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#update created entity
        url=brokerIp+"/ngsi10/updateContext"
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata41),headers=headers)
        #print(r.status_code)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

	#validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200
        
#testCase 19
'''
  Testing for get subscripton requests
'''
def test_getsubscription24():
	url=brokerIp+"/ngsi10/subscription"
	r=requests.get(url)
	assert r.status_code == 200

#testCase 20
'''
  Testing for get all entities
'''
def test_getallentities():
	url=brokerIp+"/ngsi10/entity"
	r=requests.get(url)
	assert r.status_code == 200
	
#testCase 21
'''
  Testing  for query request using Id
'''
def test_queryrequest1():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata42),headers=headers)
        #print(r.status_code)
        #print(r.content)
        assert r.status_code == 200

#testCase 22
'''
  Testing  for query request using type
'''
def test_queryrequest2():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata43),headers=headers)
        #print(r.status_code)
        #print(r.content)
        assert r.status_code == 200

#testCase 23
'''
  Testing  for query request using geo-scope(polygon)
'''
def test_queryrequest3():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata44),headers=headers)
        #print(r.status_code)
        #print(r.content)
        assert r.status_code == 200
	
#testCase 24
'''
  Testing  for query request multiple filter
'''
def test_queryrequest4():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata45),headers=headers)
        #print(r.status_code)
        #print(r.content)
        assert r.status_code == 200
