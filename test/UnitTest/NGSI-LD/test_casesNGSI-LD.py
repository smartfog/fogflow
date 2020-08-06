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
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 2
'''
  To test create entity with context in payload
'''
def test_case2():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata2),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 3
'''
  To test create entity with context in Link header and request payload is already expanded
'''
def test_case3():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata3),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 4
'''
  To test to append additional attributes to an existing entity
'''
def test_case4():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata4),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 204


#testCase 5
'''
  To test to patch  update specific attributes of an existing entity A100
'''
def test_case5():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata5),headers=headers)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 6
'''
  To test to update the value of a specific attribute of an existing entity with wrong payload
'''
def test_case6():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4580/attrs/brandName"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata6),headers=headers)
        #print(r.content)
	#print(r.status_code)
        assert r.status_code == 204

#testCase 7
'''
  To test create entity  without passing Header
'''
def test_case7():
	url=brokerIp+"/ngsi-ld/v1/entities/"
	r=requests.post(url,data=json.dumps(ld_data.subdata13))
	#print(r.content)
	#print(r.status_code)
	assert r.status_code == 400

#testCase 8
'''
  To test create entity without passing Link Header 
'''
def test_case8():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
	r=requests.post(url,data=json.dumps(ld_data.subdata),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 9
'''
  To test Update entity with two  header namely Content Type and Accept header  and posting duplicate attribute
'''
def test_case9():
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 207

#testCase 10
'''
  To test Update entity with wrong id
'''
def test_case10():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A800/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 404

#testCase 11
'''
  To test Update entity with wrong id and not passing Accept Header
'''
def test_case11():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A800/attrs"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 404

#testCase 12
'''
  To test Update entity with three  header namely Content-Type, Accept and context Link 
'''
def test_case12():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4580/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 13
'''
  To test Update entity with different  headers namely Content-Type, Accept and context Link and passing inappropriate payload
'''
def test_case13():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata15),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 500

#testCase 14
'''
  To test Update entity without header
'''
def test_case14():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        r=requests.post(url,data=json.dumps(ld_data.subdata14))
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 15
'''
  To test Update entity patch request with different header namely Content-Type, Accept and context Link  and passing inappropriate payload
'''
def test_case15():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4580/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata15),headers=headers) 
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 500

#testCase 16
'''
  To test to update entity by first creating entity without corresponding  attribute
'''
def test_case16():
	#create NGSI-LD entity 
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata16),headers=headers)
        #print(r.content)
        #print(r.status_code)

	#update the corresponding entity using patch        
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A700/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 207

#testCase 17
'''
  To test Update entity patch request without header
'''
def test_case17():
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A100/attrs"
        r=requests.patch(url,data=json.dumps(ld_data.subdata14))
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

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

	#update the corresponding entity using patch
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A900/attrs"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata14),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 207

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
        #print(r.content)
        #print(r.status_code)
   
	#to append  the attribute of corresponding entity
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A500/attrs"
        headers={'Content-Type' : 'application/json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata4),headers=headers)
        #print(r.content)
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
        #print(r.content)
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
        print(resp)
        if resp[0]["properties"][1]["id"]=="http://example.org/vehicle/brandName":
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
  To test to retrieve a specific entity by Type, context in Link Header
'''
def test_case29():
        url=brokerIp+"/ngsi-ld/v1/entities?type=Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["type"]=="https://uri.etsi.org/ngsi-ld/default-context/Vehicle":
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
  To test to create a new context source registration, with context in link header
'''
def test_case33():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata7),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 34
'''
  To test to  create a new context source registration, with context in request payload
'''
def test_case34():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata8),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 35
'''
  To test to create a new context source registration, without header
'''
def test_case35():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        r=requests.post(url,data=json.dumps(ld_data.subdata7))
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 36
'''
  To test to create registration  with change in only 2nd  entity in payload
'''
def test_case36():
	url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata21),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 409

#testCase 37
'''
  To test to create registration with only 1 entity in payload
'''
def test_case37():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata22),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testcase 38
'''
  To test get the regiestered source entity on discovery for id = town$
'''
def test_case38():
	url=discoveryIp+"/ngsi9/registration/town$"
	r=requests.get(url)
	resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["ID"]=="town$":
                print("\nValidated")
        else:
                print("\nNot Validated")
	#print(r.status_code)
	assert r.status_code == 200 

#testCase39
'''
  To test Update registration on discovery if it is reflecing or not 
'''
def test_case39():
	#get before update 
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A666"
        r=requests.get(url)
        #print(r.content)
        #print(r.status_code)

	#patch request to update
	url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/urn:ngsi-ld:ContextSourceRegistration:csr1a4001"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata23),headers=headers)
        #print(r.content)
        #print(r.status_code)
        
	#fetching from discovery
	url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A666"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A666":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 40
'''
  To test update request for registration with wrong payload
'''
def test_case40():
	#create registration 
	url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata24),headers=headers)
        #print(r.content)
        #print(r.status_code)

 
        #patch request to update
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/urn:ngsi-ld:ContextSourceRegistration:csr1a4002"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata25),headers=headers)
        #print(r.content)
        #print(r.status_code)
	assert r.status_code == 204

#testCase 41
'''
  To test the get entity from discovery for enitity Id = A662
'''
def test_case41():
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A662"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A662":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200


#testCase 42
'''
  To test to update an existing context source registration, with context in request payload
'''
def test_case42():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/urn:ngsi-ld:ContextSourceRegistration:csr1a4002"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata9),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 204


#testCase 43
'''
  To test to update an existing context source registration with idPattern , with context in request payload regarding one entity
'''
def test_case43():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/urn:ngsi-ld:ContextSourceRegistration:csr1a3459"
        headers={'Content-Type' : 'application/json'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata19),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 44
'''
  To test to get a registration by Type
'''
def test_case44():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["information"][0]["entities"][0]["type"]=="http://example.org/vehicle/Vehicle":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200


#testCase 45
'''
  To test to get a registration by Type, context in Link Header
'''
def test_case45():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?type=Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["information"][0]["entities"][0]["type"]=="https://uri.etsi.org/ngsi-ld/default-context/Vehicle":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 46
'''
  To test to get a registration by ID and Type
'''
def test_case46():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?id=urn:ngsi-ld:Vehicle:A456&type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp[0]["information"][0]["entities"][0]["type"]=="http://example.org/vehicle/Vehicle" and resp[0]["information"][0]["entities"][0]["id"]=="urn:ngsi-ld:Vehicle:A456":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 47
'''
  To test to get a registration by IdPattern and Type
'''
def test_case47():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?idPattern=urn:ngsi-ld:Vehicle:A*&type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        assert r.status_code == 404

#testCase 48
'''
  To test to delete an existing context source registration based on registration id
'''
def test_case48():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/urn:ngsi-ld:ContextSourceRegistration:csr1a3459"
        r=requests.delete(url)
        #print(r.status_code)
        assert r.status_code == 204


#testCase 49
'''
  To test to get registration by Type with wrong query
'''
def test_case49():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?type="
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 50
'''
  To test to get registration by idPattern and Type with wrong query
'''
def test_case50():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?idPattern=&type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 51
'''
  To test to get a registration by ID and Type with werong query
'''
def test_case51():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?id&type=http://example.org/vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 52
'''
  To test to create a new Subscription to with context in Link header
'''
def test_case52():
        #create subscription
	url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata10),headers=headers)
	#print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 53
'''
  To test to retrieve all the subscriptions
'''
def test_case53():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Accept' : 'application/ld+json'}
        r=requests.get(url,headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 200

#testCase 54
'''
  To test to retrieve a specific subscription based on subscription id
'''
def test_case54():
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

#testCase 55
'''
  To test to update a specific subscription based on subscription id, with context in Link header
'''
def test_case55():
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

#testCase 56
'''
  To test to update a specific subscription based on subscription id, without header
'''
def test_case56():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        r=requests.patch(url,data=json.dumps(ld_data.subdata12))
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 57
'''
  To test to update a specific subscription based on subscription id, with context in Link header and different payload
'''
def test_case57():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        headers={'Content-Type':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.patch(url,data=json.dumps(ld_data.subdata20),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 58
'''
  To test to delete a specific subscription based on subscription id
'''
def test_case58():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7"
        r=requests.delete(url)
        #print(r.status_code)
        assert r.status_code == 204

#testCase 59
'''
  To test for empty payload in entity creation
'''
def test_case59():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 60
'''
  To test for empty payload in csource registration
'''
def test_case60():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 61
'''
  To test for empty payload in subscription
'''
def test_case61():
        url=brokerIp+"/ngsi-ld/v1/subscriptions/"
        headers={'Content-Type' : 'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata26),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 62
'''
  To test for ModifiedAt and CreatedAt in entity creation
'''
def test_case62():
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

#testCase 63
'''
  To test for ModifiedAt and CreatedAt in Csource Registration
'''
def test_case63():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations?type=http://example.org/vehicle/Vehicle"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if  resp[0]["information"][0]["entities"][0]["type"]=="http://example.org/vehicle/Vehicle":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 64
'''
  To test for ModifiedAt and CreatedAt in susbcription
'''
def test_case64():
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

#testCase 65
'''
  To test for csource registartion with Id pattern
'''
def test_case65():
        url=brokerIp+"/ngsi-ld/v1/csourceRegistrations/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata28),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 400

#testCase 66
'''
  To test for update subscription over discovery
'''
def test_case66():
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

#testCase 67
'''
  To test entity creation with nested property with context in payload
'''
def test_case67():
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
        r=requests.post(url,data=json.dumps(ld_data.subdata33),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201


#testCase 68
'''
  To test create entity with nested property with  context in Link
'''
def test_case68():
	url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata34),headers=headers)
        #print(r.content)
        #print(r.status_code)
        assert r.status_code == 201

#testCase 69
'''
  To test to retrieve entity with id as urn:ngsi-ld:B990
'''
def test_case69():
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:B990"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json'}
        r=requests.get(url,headers=headers)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp)
        if resp["id"]=="urn:ngsi-ld:Vehicle:B990":
                print("\nValidated")
        else:
                print("\nNot Validated")
        #print(r.status_code)
        assert r.status_code == 200

#testCase 70
'''
  To test and retrieve the entity from discovery
'''
def test_case70():
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

