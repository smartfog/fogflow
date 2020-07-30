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
  To test subscription request 
'''
def test_getSubscription1():
        url=brokerIp+"/ngsi10/subscribeContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata1),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        assert r.status_code == 200


#testCase 2
'''
  To test entity creation with attributes, then susbscribing and get subscription using ID
'''
def test_getSubscription2():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata2),headers=headers)
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
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
	resp=resp['entities']
        sid2=resp[0]["id"]
	if "Result1"==sid2:
		print("\nValidated")
	else:
		print("\nNot Validated")
        assert r.status_code == 200

#testCase 3
'''
 To test entity creation with one  attribute : pressure only followed by subscribing and get using ID
'''
def test_getSubscription3():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata4),headers=headers)
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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result2"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200

#testCase 4
'''
  To test entity creation with one attribute : Temperature only followed by subscription and get using ID
'''
def test_getSubscription4():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata6),headers=headers)
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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result3"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200

#testCase 5
'''
   To test create entity without passing Domain data followed by subscription and get using ID
'''
def test_getSubscription5():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata8),headers=headers)
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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result4"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200

#testCase 6
'''
   To test create entity without attributes followed by subscription and get using Id
'''
def test_getSubscription6():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata10),headers=headers)
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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result5"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200


#testCase 7
'''
   To test create entity without attributes and Metadata and followed by sbscription and get using Id
'''
def test_getSubscription7():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata12),headers=headers)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)
	#print(r.status_code)

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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result6"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200


#testCase 8
'''
   To test create entity without entity type followed by subscription and get using Id
'''
def test_getSubscription8():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r = requests.post(url,data=json.dumps(data_ngsi10.subdata14),headers=headers)
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
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        resp=resp['entities']
        sid2=resp[0]["id"]
        if "Result7"==sid2:
                print("\nValidated")
        else:
                print("\nNot Validated")
        assert r.status_code == 200



#testCase 9
'''
  To test get subscription request by first posting subscription request followed by delete request
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
	print("Subscription with sid-"+sid+" not found")
        assert r.status_code == 404


#testCase 10
'''
  To test the update post request to create entity 
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
  To test  subscription with attributes and using ID to validate it
'''
def test_getSubscription11():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata18),headers=headers)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata19),headers=headers)
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
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #validate via accumulator
        url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : sid})
        print(r.content)
        assert r.status_code == 200

#testCase 12
'''
  To test subscription for its if and else part : 1) for Destination Header
'''
def test_getSubscription12():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata21),headers=headers)
        resp_content=r.content
        resInJson=resp_content.decode('utf8').replace("'",'"')
        resp=json.loads(resInJson)
        #print(resp)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination' : 'orion-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata22),headers=headers)
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
  To test subscription for its if and else part : 2) for User - Agent Header
'''
def test_getSubscription18():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata24),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','User-Agent' : 'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata25),headers=headers)
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
  To test subcription for its if else part : 3) Require-Reliability Header
'''
def test_getSubscription19():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata27),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Require-Reliability' : 'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata28),headers=headers)
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
  To test subscription with two headers simultaneously : 4) Destination and User-Agent
'''
def test_getSubscription20():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata30),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination' : 'orion-broker','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata31),headers=headers)
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
  To test subscription with two headers simultaneously : 4)  User-Agent and Require-Reliability
'''
def test_getSubscription21():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata33),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','User-Agent':'lightweight-iot-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata34),headers=headers)
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
   To test subscription with two headers simultaneously : 4)  Destination  and Require-Reliability headers
'''
def test_getSubscription22():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata36),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination':'orion-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata37),headers=headers)
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
  To test subscription with all  headers simultaneously : 5)  Destination, User-Agent  and Require-Reliability headers
'''
def test_getSubscription23():
        #create an entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata39),headers=headers)
        resp_content1=r.content
        resInJson=resp_content1.decode('utf8').replace("'",'"')
        resp1=json.loads(resInJson)
        #print(resp1)

        #subscribing
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json','Destination':'orion-broker','User-Agent':'lightweight-iot-broker','Require-Reliability':'true'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata40),headers=headers)
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
  To test for get subscripton requests
'''
def test_getsubscription24():
        url=brokerIp+"/ngsi10/subscription"
        r=requests.get(url)
        assert r.status_code == 200

#testCase 20
'''
  To test for get all entities
'''
def test_getallentities():
        url=brokerIp+"/ngsi10/entity"
        r=requests.get(url)
        assert r.status_code == 200

#testCase 21
'''
  To test  for query request using Id
'''
def test_queryrequest1():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata42),headers=headers)
        #print(r.content)
        assert r.status_code == 200

#testCase 22
'''
  To test  for query request using type
'''
def test_queryrequest2():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata43),headers=headers)
        #print(r.content)
        assert r.status_code == 200


#testCase 23
'''
  To test  for query request using geo-scope(polygon)
'''
def test_queryrequest3():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata44),headers=headers)
        #print(r.content)
        assert r.status_code == 200

#testCase 24
'''
  To test  for query request multiple filter
'''
def test_queryrequest4():
        url=brokerIp+"/ngsi10/queryContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata45),headers=headers)
        #print(r.content)
        assert r.status_code == 200

#testCase 25
'''
  To test if wrong payload is decoded or not
'''
def test_case25():
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata46),headers=headers)
        #print(r.status_code)
        assert r.status_code == 200


#testCase26
'''
  To test the response on passing DELETE in updateAction in payload
'''
def test_case26():
        #create v1 entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata47),headers=headers)
        #print(r.content)

        #get the created entity
        url=brokerIp+"/ngsi10/entity/Result047"
        r=requests.get(url)
	#print(r.content)	

        #passing DELETE in update payload
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata48),headers=headers)
        #print(r.content)

        #get the created entity
        url=brokerIp+"/ngsi10/entity/Result047"
        r=requests.get(url)
        #print(r.content)
        assert r.status_code == 404

#testCase 27
'''
  To test the entity creation with empty payload
'''
def test_case27():
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata48),headers=headers)
        #print(r.content)
        assert r.status_code == 200

#testCase 28
'''
  To test the subscription with empty payload
'''
def test_case28():
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata48),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)
        assert r.status_code == 200

#testCase 29
'''
  To get subscription of empty payload when subscribing
'''
def test_case29():
        url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata48),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

        #get request
        get_url=brokerIp+"/ngsi10/subscription/"
        url=get_url+sid
        r=requests.get(url)
        #print(r.content)
        assert r.status_code == 200

#testCase 30
'''
  To test the action of API on passing an attribute as a command in payload
'''
def test_cases30():
	url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','User-Agent':'lightweight-iot-broker'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata50),headers=headers)
        #print(r.content)
        assert r.status_code == 500

#testCase 31
'''
  To test the fiware header with updateAction equal to UPDATE
'''
def test_case31():
	#create entity
	url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','fiware-service':'iota','fiware-servicepath':'/'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata51),headers=headers)
        #print(r.content)
        
	#get entity
	url=brokerIp +"/ngsi10/entity/Result048"
	r=requests.get(url)
	#print(r.content)
	#print(r.status_code)
        assert r.status_code == 404

#testCase 32
'''
  To test the fiware header with updateAction equal to APPEND
'''
def test_case32():
	url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','fiware-service':'Abc','fiware-servicepath':'pqr'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata52),headers=headers)
        #print(r.content)
        assert r.status_code == 200

#testCase 33
'''
  To test the fiware header with updateAction equal to delete
'''
def test_case33():
        #create v1 entity
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','fiware-service':'Abc','fiware-servicepath':'pqr'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata53),headers=headers)
        #print(r.content)

        #get the created entity
        url=brokerIp+"/ngsi10/entity/Result053"
        r=requests.get(url)
        #print(r.content)

        #passing DELETE in update payload
        url=brokerIp+"/ngsi10/updateContext"
        headers={'Content-Type':'appliction/json','fiware-service':'Abc','fiware-servicepath':'pqr'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata54),headers=headers)
        #print(r.content)

        #get the created entity
        url=brokerIp+"/ngsi10/entity/Result053"
        r=requests.get(url)
        #print(r.content)
        assert r.status_code == 404

#testCase 34
'''
  To test the notifyContext request 
'''
def test_case34():
	url=brokerIp+"/ngsi10/notifyContext"
	headers={'Content-Type':'appliction/json'}
	r=requests.post(url,data=json.dumps(data_ngsi10.subdata55),headers=headers)
	#print(r.content)
	assert r.status_code == 200

#testCase 35
'''
  To test unsubscribing feature
'''
def test_case35():
	#create subscription
	url=brokerIp+"/ngsi10/subscribeContext"
        headers= {'Content-Type': 'application/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata56 ),headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        resp=resp['subscribeResponse']
        sid=resp['subscriptionId']
        #print(sid)

	#unsubscribe Context
	url=brokerIp+"/ngsi10/unsubscribeContext"
	headers={'Content-Type': 'application/json'}
	r=requests.post(url,json={"subscriptionId":sid,"originator":"POMN"},headers=headers)
	#print(r.content)
	assert r.status_code == 200 

#testCase 36
'''
  To test entity creation using other route
'''
def test_case36():
        url=brokerIp+"/v1/updateContext"
        headers={'Content-Type':'appliction/json'}
        r=requests.post(url,data=json.dumps(data_ngsi10.subdata56),headers=headers)
        #print(r.content)
	assert r.status_code == 200

#testCase 37
'''
  To test and fetch unique entity
'''
def test_case37():
	url=brokerIp+"/ngsi10/entity/Result14"
	r=requests.get(url)
	#print(r.content)
	assert r.status_code == 200

#testCase 38
'''
  To test and fetch attribute specific to an entity
'''
def test_case38():
	url=brokerIp+"/ngsi10/entity/Result14/pressure"
	r=requests.get(url)
	resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
	#print(r.content)
        name=resp['name']
        type1=resp['type']
	val=resp['value']
	if name=='pressure' and type1=='float' and val==55:
		print("\nValidated")
	else:
		print("\nNot Validated")
	assert r.status_code == 200


