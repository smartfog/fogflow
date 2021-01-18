import os,sys
# change the path accoring to the test folder in system
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import ld_data
import sys

# change it by broker ip and port
brokerIp="http://localhost:8070"
discoveryIp="http://localhost:8090"

print("Testing of NGSI-LD")
# testCase 1
'''
  To test create entity with context in Link Header
'''
def test_case1():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata1),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201


#testCase 2
'''
  To test create entity with context in payload
'''
def test_case2():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata2),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201


#testCase 3
'''
  To test create entity with context in Link header and request payload is already expanded
'''
def test_case3():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata3),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 4
'''
  To test to append additional attributes to an existing entity
'''
def test_case4():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json', 'Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata4),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201


#testCase 5
'''
  To test to update specific attributes of an existing entity A100
'''
def test_case5():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata5),headers=headers)
        print(r.status_code)
        assert r.status_code == 201

#testCase 6
'''
  To test to update the value of a specific attribute of an existing entity with wrong payload
'''
def test_case6():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata6),headers=headers)
        print(r.content)
	print(r.status_code)
        assert r.status_code == 400

#testCase 7
'''
  To test create entity  without passing Header
'''
def test_case7():
	url=brokerIp+"/ngsi-ld/v1/entities/"
	r=requests.post(url,data=json.dumps(ld_data.subdata13))
	print(r.content)
	print(r.status_code)
	assert r.status_code == 400

#testCase 8
'''
  To test create entity without passing Link Header 
'''
def test_case8():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
	r=requests.post(url,data=json.dumps(ld_data.subdata),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 9
'''
  To test Update entity with two  header namely Content Type and Accept header  and posting duplicate attribute
'''
def test_case9():
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 10
'''
  To test Update entity with wrong id format 
'''
def test_case10():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14b),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 400 

#testCase 11
'''
  To test Update entity with wrong id format and  not passing Accept Header
'''
def test_case11():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14b),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 400

#testCase 12
'''
  To test Update entity with three  header namely Content-Type, Accept and context Link 
'''
def test_case12():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 13
'''
  To test Update entity with different  headers namely Content-Type, Accept and context Link and passing inappropriate payload
'''
def test_case13():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata15),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 500

#testCase 14
'''
  To test Update entity without header
'''
def test_case14():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        r=requests.post(url,data=json.dumps(ld_data.subdata14))
        print(r.content)
        print(r.status_code)
        assert r.status_code == 400


#testCase 16
'''
  To test to update entity by first creating entity without corresponding  attribute
'''
def test_case16():
	#create NGSI-LD entity 
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata16),headers=headers)
        print(r.content)
        print(r.status_code)

	#update the corresponding entity using post        
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14c),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 18
'''
  To test to update entity by first creating entity with all  attributes missing 
'''
def test_case18():
	#create NGSI-LD entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata17),headers=headers)
        #print(r.content)
        #print(r.status_code)

	#update the corresponding entity using post
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14d),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 19
'''
  To test to delete NGSI-LD context entity
'''
def test_19():
	#create NGSI-LD  entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata32),headers=headers)
        #print(r.content)
        #print(r.status_code)

   	#to delete corresponding entity
 	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A999"
        headers={'Content-Type':'application/json','Accept':'application/ld+json'}
	r=requests.delete(url,headers=headers)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 20
'''
  To test to delete an attribute of an NGSI-LD context entity
'''
def test_case20():
	#create NGSI-LD entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata11),headers=headers)
        print(r.content)
        #print(r.status_code)
   
	#to append  the attribute of corresponding entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata4b),headers=headers)
        print(r.content)
        #print(r.status_code)
	
	#to delete the attribute of corresponding entity
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A500/attrs/brandName1"
        r=requests.delete(url)
	#print(r.status_code)
        assert r.status_code == 204

#testCase 21
'''
  To test to delete an attribute of an NGSI-LD context entity which does not have any attribute
'''
def test_case21():
	#create NGSI-LD entity 
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata18),headers=headers) 
        #print(r.content)
        #print(r.status_code)

	#to delete attribute of corresponding entity 
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A501/attrs/brandName1"
        r=requests.delete(url)
        print(r.content)
	#print(r.status_code)
        assert r.status_code == 404


#testCase 22
'''
  To test to retrieve a specific entity which is deleted 
'''
def test_case22():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A999"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
	#print(r.status_code)
        assert r.status_code == 404

#testCase 23
'''
  To test to retrieve a specific entity which is existing
'''
def test_case23():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4580"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
	resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["id"]=="urn:ngsi-ld:Vehicle:A4580":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 24
'''
  To test to retrieve entities by attributes  
'''
def test_case24():
        url=brokerIp+"/ngsi-ld/v1/entities?attrs=http://example.org/vehicle/brandName"
        headers={'Content-Type':'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["http://example.org/vehicle/brandName"]["value"]=="MARUTI":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 25
'''
  To test to retrieve entities by attributes with wrong query 
'''
def test_case25():
        url=brokerIp+"/ngsi-ld/v1/entities?attrs"
        headers={'Content-Type':'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 26
'''
  To test to retrieve a specific entity by ID and Type
'''
def test_case26():
        url=brokerIp+"/ngsi-ld/v1/entities?id=urn:ngsi-ld:Vehicle:A4580&type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["type"]=="http://example.org/vehicle/Vehicle"  and resp[0]["id"]=="urn:ngsi-ld:Vehicle:A4580":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 27
'''
  To test to retrieve a specific entity by Type
'''
def test_case27():
        url=brokerIp+"/ngsi-ld/v1/entities?type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
       	resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["type"]=="http://example.org/vehicle/Vehicle":
                print("\nValidated")
        else:
                print("\nNot Validated")
	#print(r.status_code)
        assert r.status_code == 200

#testCase 28
'''
  To test to retrieve a specific entity by Type with wrong query
'''
def test_case28():
        url=brokerIp+"/ngsi-ld/v1/entities?type=http://example.org"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 404

#testCase 29
'''
  To test to retrieve  entities by Type, with context in Link Header
'''
def test_case29():
        url=brokerIp+"/ngsi-ld/v1/entities?type=Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["type"]=="Vehicle":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 30
'''
  To test to retrieve a specific entity by Type, context in Link Header and wrong query
'''
def test_case30():
        url=brokerIp+"/ngsi-ld/v1/entities?type"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 31
'''
  To test to To retrieve a specific entity by IdPattern and Type
'''
def test_case31():
        url=brokerIp+"/ngsi-ld/v1/entities?idPattern=urn:ngsi-ld:Vehicle:A.*&type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["type"]=="http://example.org/vehicle/Vehicle" and resp[0]["id"].find("A")!=-1:
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 32
'''
  To test to retrieve an entity registered over Discovery
'''
def test_case32():
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A4580"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A4580":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 33
'''
  To test to create a new Subscription to with context in Link header
'''
def test_case33():
        #create subscription
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata10),headers=headers)
	#print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 34
'''
  To test to retrieve all the subscriptions
'''
def test_case34():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Accept' : 'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 200

#testCase 35
'''
  To test to retrieve a specific subscription based on subscription id
'''
def test_case35():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if  resp["id"]=="urn:ngsi-ld:Subscription:7":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 36
'''
  To test to update a specific subscription based on subscription id, with context in Link header
'''
def test_case36():
        #get subscription before update 
	url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)

	#Update the subscription	
	url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Content-Type':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata12),headers=headers)
        #print(r.content)
        #print(r.status_code)
        
	#get subscription after update
	url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if  resp["id"]=="urn:ngsi-ld:Subscription:7":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
	assert r.status_code == 200

#testCase 37
'''
  To test to update a specific subscription based on subscription id, without header
'''
def test_case37():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        r=requests.patch(url,data=json.dumps(ld_data.subdata12))
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 38
'''
  To test to update a specific subscription based on subscription id, with context in Link header and different payload
'''
def test_case38():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Content-Type':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata20),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 39
'''
  To test to delete a specific subscription based on subscription id
'''
def test_case39():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        r=requests.delete(url)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 40
'''
  To test for empty payload in entity creation
'''
def test_case40():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 41
'''
  To test for empty payload in csource registration
'''
def test_case41():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 42
'''
  To test for empty payload in subscription
'''
def test_case42():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 43
'''
  To test for ModifiedAt and CreatedAt in entity creation
'''
def test_case43():
        #Entity Creation
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata27),headers=headers)
        #print(r.content)
        #print(r.status_code)

        #Fetching Entity
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A6000"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if  resp["id"]=="urn:ngsi-ld:Vehicle:A6000":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 44
'''
  To test for ModifiedAt and CreatedAt in susbcription
'''
def test_case44():
	#create subscription 
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata29),headers=headers)
        #print(r.content)
        #print(r.status_code)
        
	#making a get request 
	url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:8"
        headers={'Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if  resp["id"]=="urn:ngsi-ld:Subscription:8":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 45
'''
  To test for csource registartion with Id pattern
'''
def test_case45():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata28),headers=headers)
        print(r.content)
        #print(r.status_code)
        assert r.status_code == 400  

#testCase 46
'''
  To test for update subscription over discovery
'''
def test_case46():
	#create a subscription
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata30),headers=headers)
        #print(r.content)
        #print(r.status_code)

	#get subscription over discovery before update
	url=discoveryIp+"/ngsi9/subscription"
	r=requests.get(url)
	#print(r.content)
	#print(r.status_code)

	#update the subscription
	url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:10"
        headers={'Content-Type':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata31),headers=headers)
        #print(r.content)
        #print(r.status_code)
	
	#get subscription after update 
	url=discoveryIp+"/ngsi9/subscription"
        r=requests.get(url)
        #print(r.content)
        #print(r.status_code)
	assert r.status_code == 200

#testCase 47
'''
  To test entity creation with nested property with context in payload
'''
def test_case47():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata33),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 48
'''
  To test create entity with nested property with  context in Link
'''
def test_case48():
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata34),headers=headers)
        print(r.content)
        print(r.status_code)
        assert r.status_code == 201

#testCase 49
'''
  To test to retrieve entity with id as urn:ngsi-ld:B990
'''
def test_case49():
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:B990"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
	#print(r.content)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["id"]=="urn:ngsi-ld:Vehicle:B990":
                print("\nValidated")
        else:
                print("\nNot Validated")
        print(r.status_code)
        assert r.status_code == 200

#testCase 50
'''
  To test and retrieve the entity from discovery
'''
def test_case50():
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:C001"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["ID"]=="urn:ngsi-ld:Vehicle:C001":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 51
'''
  To test if multiple subscription can be created with same subscription Id
'''
def test_case51():
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
	headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
	r=requests.post(url,data=json.dumps(ld_data.subdata36),headers=headers)
        #print(r.content)
        #print(r.status_code)

        #get subscription over discovery before update
        url=discoveryIp+"/ngsi9/subscription"
        r=requests.get(url)
        #print(r.content)
        #print(r.status_code)
	
	#to create same subscription again
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata36),headers=headers)
        print(r.content)
        #print(r.status_code)
	assert r.status_code == 201

#testCase 52
'''
  To test if delete attribute is reflected over discovery
'''
def test_case52():
	#to create an entity
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata38),headers=headers)
        #print(r.content)
        #print(r.status_code)

	#to fetch the registration from discovery
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"])
	print("\nchecking if brandName attribute is present in discovery before deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"]["type"] == "Property":
			print("\n-----> brandName is existing...!!")
		else:
			print("\n-----> brandName does not exist..!")
        else:
                print("\nNot Validated")
        #print(r.status_code)


	#to delete brandName attribute
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A3000/attrs/brandName"
        r=requests.delete(url)
        #print(r.content)
        #print(r.status_code)

	#To fetch registration again from discovery
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"])
	print("\nchecking if brandName attribute is present in discovery after deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if "https://uri.etsi.org/ngsi-ld/default-context/brandName" in resp["AttributesList"]:
                        print("\n-----> brandName is existing...!!")
                else:
                        print("\n-----> brandName does not exist because deleted...!")
        else:
                print("\nNot Validated")
	assert r.status_code == 200

#testCase 53
'''
  To test if appended attribute is reflected on discovery 
'''
def test_case53():
	#to fetch registration of entity from discovery before appending
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"])
        print("\nchecking if brandName1 attribute is present in discovery after deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if "https://uri.etsi.org/ngsi-ld/default-context/brandName1" in resp["AttributesList"]:
                        print("\n-----> brandName1 is existing...!!")
                else:
                        print("\n-----> brandName1 does not exist yet...!")
        else:
                print("\nNot Validated")

	#to append an entity with id as urn:ngsi-ld:Vehicle:A3000
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata4c),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

	#to fetch registration of entity from discovery after appending
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"])
        print("\nchecking if brandName1 attribute is present in discovery after deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if "https://uri.etsi.org/ngsi-ld/default-context/brandName1" in resp["AttributesList"]:
                        print("\n-----> brandName1 is existing after appending...!!")
                else:
                        print("\n-----> brandName1 does not exist yet...!")
        else:
                print("\nNot Validated")
	assert r.status_code == 200

#testCase 54
'''
  To test if discovery's context availablity is updated on updating 
'''
def test_case54():
	#to create entity
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata40),headers=headers)
        #print(r.content)
        #print(r.status_code)

	#to create subscription
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata39),headers=headers)
	#print(r.content)
        #print(r.status_code)

	#Update entity to fire notification
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata41),headers=headers)
        #print(r.status_code)
        
	#to validate
	url="http://0.0.0.0:8888/validateNotification"
        r=requests.post(url,json={"subscriptionId" : "urn:ngsi-ld:Subscription:020"})
        print(r.content)
        assert r.status_code == 200

# testCase 55
'''
  To test if instanceId is fetched while creating entity 
'''
def test_case55():
	# to create entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata42),headers=headers)
        #print(r.content)
        #print(r.status_code)
        
	# to fetch and verify instanceId
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:C001"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["brandName1"]["instanceId"]=="instance1":
                print("\nValidated")
        else:
                print("\nNot Validated")
        print(r.status_code)
        assert r.status_code == 200


# testCase 56
'''
  To test if datasetId is fetched while creating entity
'''
def test_case56():
        # to create entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata43),headers=headers)
        #print(r.content)
        #print(r.status_code)

        # to fetch and verify instanceId
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:C002"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["brandName1"]["datasetId"]=="dataset1":
                print("\nValidated")
        else:
                print("\nNot Validated")
        print(r.status_code)
        assert r.status_code == 200


#testCase 57
'''
  To test for subscription without entities in Payload
'''
def test_case57():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata44),headers=headers)
        print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 58
'''
  To test for subscription with different type in payload 
'''
def test_case58():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata45),headers=headers)
        print(r.content)
        #print(r.status_code)
        assert r.status_code == 400 
	
