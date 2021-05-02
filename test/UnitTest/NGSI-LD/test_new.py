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
brokerIp="http://127.0.0.1:8070"
accumulatorURl ="http://127.0.0.1:8888"
discoveryIp="http://127.0.0.1:8090"

# test if header content-Type application/json is allowed or not 
def test_case74():
	url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
	headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
	r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
	assert r.status_code == 204


# test if header content-Type is application/ld+json then the link header should not be persent in request

def test_case75():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        assert r.status_code == 404

#test if Allowd Content-Type are only appliation/json and application/ld+json

def test_case76():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
	print(r.status_code)
        assert r.status_code == 400


# test create and get the entity in openiot FiwareService

def test_case77():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        #print(r.status_code)
	url=brokerIp+'/ngsi-ld/v1/entities/'+'urn:ngsi-ld:Device:water001'
	r = requests.get(url,headers=headers)
	assert r.status_code == 200


def test_case78():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        #print(r.status_code)
	headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiott','fiware-servicepath' :'test'}

        url=brokerIp+'/ngsi-ld/v1/entities/'+'urn:ngsi-ld:Device:water001'
        r = requests.get(url,headers=headers)
        assert r.status_code == 404

# To test upsert Api support only array of entities
def test_case79():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData75),headers=headers)
        #print(r.status_code)
	assert r.status_code == 500

def test_case80():
        upsertURL=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Integration': 'NGSILDBroker','Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        rUpsert=requests.post(upsertURL,data=json.dumps(ld_data.upsertCommand80),headers=headers)
	subscribeURL=brokerIp+"/ngsi-ld/v1/subscriptions/"
	rSubscribe=requests.post(subscribeURL,data=json.dumps(ld_data.subData80),headers=headers)
        print(rSubscribe.status_code)
	time.sleep(5)
	getURL=accumulatorURl+"/validateupsert"
	rget = requests.get(getURL)
	print(rget.content)
        assert rget.content == "200"

#To test get Entity by Eid from broker if FiwareService is provided
def test_case81():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Device:water001"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200


# To test get Entity form broker By Eid if FiwareService is not provided

def test_case82():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Device:water001"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200


# To test get all Entity from broker if FiwareService is provided
def test_case83():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.upsertMultipleCommand),headers=headers)
        url=brokerIp+"/ngsi-ld/v1/entities?type=Device"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200

# To test get Entity form broker if FiwareService is not provided
def test_case84():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.upsertMultipleCommand),headers=headers)
        url=brokerIp+"/ngsi-ld/v1/entities?type=Device"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200


# Test if registration of entity is available in discovery or not if fiwareService is provided
def test_case85():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' ,'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Device:water001"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200


#test if registration of entity is available in discovery or not if fiwareService is not provided
def test_case86():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Device:water001"
        r = requests.get(url,headers=headers)
        assert r.status_code == 200

#test response of discovery if Entity does not exist in the disovery with fiwareService 
def test_case87():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"','fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Device:water0010"
        r = requests.get(url,headers=headers)
        assert r.content == "null"

#test response of discovery if Entity does not exist in the disovery without fiwareService
def test_case88():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.upsertCommand),headers=headers)
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Device:water0010"
        r = requests.get(url,headers=headers)
        assert r.content == "null"

# test Delete Entity from thinbroker if FiwareService is provided

def test_case89():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"','fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.DelData),headers=headers)
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A109"
        rget = requests.get(url,headers=headers)
        assert rget.status_code == 200
	delURL=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A109"
	r = requests.delete(delURL,headers=headers)
	rget = requests.get(url,headers=headers)
	assert rget.status_code == 404

# test creation of Entity in one fiwareService and delete the Entity with same id in different fiwareService

def test_case90():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.test89),headers=headers)
	url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Device:test89"
	rget = requests.get(url,headers=headers)
	assert rget.status_code == 200
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"','fiware-service' : 'openiot','fiware-servicepath' :'test'}
        delURL=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Device:test89"
	rDel = requests.delete(delURL,headers=headers)
	assert rDel.status_code == 404
	assert rget.status_code == 200



