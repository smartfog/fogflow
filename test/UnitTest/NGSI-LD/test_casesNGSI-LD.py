#!/usr/bin/python
# -*- coding: utf-8 -*-
import os
import sys
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import ld_data
import sys

# change it by broker ip and port

brokerIp = 'http://127.0.0.1:8070'
discoveryIp = 'http://127.0.0.1:8090'
accumulatorURl = 'http://127.0.0.1:8888'


# testCase 1

def test_case1():
    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata1),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 2

def test_case2():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata2),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 3

def test_case3():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata3),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 4

def test_case4():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata4),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 5

def test_case5():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata5),
                      headers=headers)
    print r.status_code
    assert r.status_code == 201


# testCase 6

def test_case6():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata6),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 7

def test_case7():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    r = requests.post(url, data=json.dumps(ld_data.subdata13))
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 8

def test_case8():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 9

def test_case9():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 10

def test_case10():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14b),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 11

def test_case11():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14b),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 12

def test_case12():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 13

def test_case13():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata15),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 14

def test_case14():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    r = requests.post(url, data=json.dumps(ld_data.subdata14))
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 15

def test_case15():

    # create NGSI-LD entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata16),
                      headers=headers)
    print r.content
    print r.status_code

    # update the corresponding entity using post
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14c),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 16

def test_case16():

    # create NGSI-LD entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata17),
                      headers=headers)
    print r.content
    print r.status_code

    # update the corresponding entity using post
        # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata14d),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 17
# def test_case17():
# ....#create NGSI-LD  entity
# ....time.sleep(3)
#        url=brokerIp+"/ngsi-ld/v1/entities/"
#        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
 #       r=requests.post(url,data=json.dumps(ld_data.subdata32),headers=headers)
 #       print(r.content)
 #       print(r.status_code)
#
#   ....#to delete corresponding entity
# ....time.sleep(3)
# ....url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A999"
#        headers={'Content-Type':'application/json','Accept':'application/ld+json'}
# ....r=requests.delete(url,headers=headers)
#        print(r.status_code)
#        assert r.status_code == 204

# testCase 18

# testCase 19

def test_case19():

    # create NGSI-LD entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata18),
                      headers=headers)
    print r.content
    print r.status_code

    # to delete attribute of corresponding entity
    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A501/attrs/brandName1'
    r = requests.delete(url)
    print r.content
    print r.status_code
    assert r.status_code == 404


# testCase 20

def test_case20():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A999'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 404


# testCase 21

def test_case21():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A4580'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Vehicle:A4580':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 22
# def test_case22():
# ....time.sleep(3)
#        url=brokerIp+"/ngsi-ld/v1/entities?attrs=https://uri.etsi.org/ngsi-ld/default-context/brandName"
#        headers={'Content-Type':'application/json','Accept':'application/ld+json'}
#        r=requests.get(url,headers=headers)
        # resp_content=r.content
        # resInJson= resp_content.decode('utf8').replace("'", '"')
        # resp=json.loads(resInJson)
#        print(r.content)
        # if resp[0]["brandName"]["value"]=="Mercedes":
         #       print("\nValidated")
        # else:
        #        print("\nNot Validated")
#        print(r.status_code)
#        assert r.status_code == 200

# testCase 23

def test_case23():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities?attrs'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 24

# testCase 25

def test_case25():

    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/entities?type=https://uri.etsi.org/ngsi-ld/default-context/Vehicle'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp[0]['type'] == 'http://example.org/vehicle/Vehicle':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 26

def test_case26():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities?type=http://example.org'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 404


# testCase 27

def test_case27():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities?type=Vehicle'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp[0]['type'] == 'Vehicle':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 28

def test_case28():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities?type'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 29
# def test_case29():
 #       url=brokerIp+"/ngsi-ld/v1/entities?idPattern=urn:ngsi-ld:Vehicle:A.*&type=http://example.org/vehicle/Vehicle"
  #      headers={'Content-Type' : 'application/json','Accept':'application/ld+json'}
   #     r=requests.get(url,headers=headers)
    #    resp_content=r.content
   #     resInJson= resp_content.decode('utf8').replace("'", '"')
   #     resp=json.loads(resInJson)
   #     print(resp)
   #     if resp[0]["type"]=="http://example.org/vehicle/Vehicle" and resp[0]["id"].find("A")!=-1:
   #             print("\nValidated")
   #     else:
   #             print("\nNot Validated")
   #     print(r.status_code)
   #     assert r.status_code == 200

# testCase 30

def test_case30():
    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A4580'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A4580':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 31

def test_case31():

        # create subscription
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata10),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 32

def test_case32():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 200


# testCase 33

def test_case33():

    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    headers = {'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Subscription:7':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 34

def test_case34():

    # time.sleep(3)
        # get subscription before update

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    headers = {'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    print r.status_code

    # Update the subscription
    # time.sleep(3)....

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.patch(url, data=json.dumps(ld_data.subdata12),
                       headers=headers)
    print r.content
    print r.status_code

    # get subscription after update
    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    headers = {'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Subscription:7':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 35

def test_case35():

    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    r = requests.patch(url, data=json.dumps(ld_data.subdata12))
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 36

def test_case36():

    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:7'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.patch(url, data=json.dumps(ld_data.subdata20),
                       headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 204


# testCase 37

# testCase 38

def test_case38():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata26),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 39

def test_case39():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata26),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 40

def test_case40():

    # time.sleep(3)
        # Entity Creation

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata27),
                      headers=headers)
    print r.content
    print r.status_code

        # Fetching Entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A6000'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Vehicle:A6000':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 41

def test_case41():

    # create subscription
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata29),
                      headers=headers)
    print r.content
    print r.status_code

    # making a get request
    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:8'
    headers = {'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    r = requests.get(url, headers=headers)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Subscription:8':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 42

def test_case42():

    # create a subscription
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata30),
                      headers=headers)
    print r.content
    print r.status_code

    # get subscription over discovery before update
    # time.sleep(3)

    url = discoveryIp + '/ngsi9/subscription'
    r = requests.get(url)
    print r.content
    print r.status_code

    # update the subscription
    # time.sleep(3)

    url = brokerIp \
        + '/ngsi-ld/v1/subscriptions/urn:ngsi-ld:Subscription:10'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.patch(url, data=json.dumps(ld_data.subdata31),
                       headers=headers)
    print r.content
    print r.status_code

    # get subscription after update
    # time.sleep(3)

    url = discoveryIp + '/ngsi9/subscription'
    r = requests.get(url)
    print r.content
    print r.status_code
    assert r.status_code == 200


# testCase 43

def test_case43():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata33),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 44

def test_case44():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata34),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 45

def test_case45():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:B990'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['id'] == 'urn:ngsi-ld:Vehicle:B990':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 46

def test_case46():
    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:C001'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp)

    if resp['ID'] == 'urn:ngsi-ld:Vehicle:C001':
        print '\nValidated'
    else:
        print '\nNot Validated'

        # print(r.status_code)

    assert r.status_code == 200


# testCase 47

def test_case47():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata36),
                      headers=headers)
    print r.content
    print r.status_code

        # get subscription over discovery before update
    # time.sleep(3)

    url = discoveryIp + '/ngsi9/subscription'
    r = requests.get(url)
    print r.content
    print r.status_code

    # to create same subscription again
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata36),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 201


# testCase 48

def test_case48():

    # to create an entity

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata38),
                      headers=headers)

        # print(r.content)
        # print(r.status_code)

    # to fetch the registration from discovery

    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A3000'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"])

    print '\nchecking if brandName attribute is present in discovery before deletion'
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A3000':
        if resp['AttributesList'
                ]['https://uri.etsi.org/ngsi-ld/default-context/brandName'
                  ]['type'] == 'Property':
            print '\n-----> brandName is existing...!!'
        else:
            print '\n-----> brandName does not exist..!'
    else:
        print '\nNot Validated'

        # print(r.status_code)

    # to delete brandName attribute

    url = brokerIp \
        + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A3000/attrs/brandName'
    r = requests.delete(url)

        # print(r.content)
        # print(r.status_code)

    # To fetch registration again from discovery

    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A3000'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"])

    print '\nchecking if brandName attribute is present in discovery after deletion'
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A3000':
        if 'https://uri.etsi.org/ngsi-ld/default-context/brandName' \
            in resp['AttributesList']:
            print '\n-----> brandName is existing...!!'
        else:
            print '\n-----> brandName does not exist because deleted...!'
    else:
        print '\nNot Validated'
    assert r.status_code == 200


# testCase 49

def test_case49():

    # to fetch registration of entity from discovery before appending

    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A3000'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"])

    print '\nchecking if brandName1 attribute is present in discovery before appending'
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A3000':
        if 'https://uri.etsi.org/ngsi-ld/default-context/brandName1' \
            in resp['AttributesList']:
            print '\n-----> brandName1 is existing...!!'
        else:
            print '\n-----> brandName1 does not exist yet...!'
    else:
        print '\nNot Validated'

    # to append an entity with id as urn:ngsi-ld:Vehicle:A3000

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata4c),
                      headers=headers)

        # print(r.content)
        # print(r.status_code)

    assert r.status_code == 201

    # to fetch registration of entity from discovery after appending

    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A3000'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"])

    print '\nchecking if brandName1 attribute is present in discovery after appending'
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A3000':
        if 'https://uri.etsi.org/ngsi-ld/default-context/brandName1' \
            in resp['AttributesList']:
            print '\n-----> brandName1 is existing after appending...!!'
        else:
            print '\n-----> brandName1 does not exist yet...!'
    else:
        print '\nNot Validated'
    assert r.status_code == 200


# testCase 50

def test_case50():

    # to create entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata40),
                      headers=headers)
    print r.content
    print r.status_code

    # to create subscription
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata39),
                      headers=headers)
    print r.content
    print r.status_code

    # Update entity to fire notification
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata41),
                      headers=headers)
    print r.status_code

    # to validate
    # time.sleep(3)

    url = 'http://0.0.0.0:8888/validateNotification'
    r = requests.post(url,
                      json={'subscriptionId': 'urn:ngsi-ld:Subscription:020'
                      })
    print r.content
    assert r.status_code == 200


# testCase 51

def test_case51():

    # to create entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata42),
                      headers=headers)
    print r.content
    print r.status_code

    # to fetch and verify instanceId
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:C001'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['brandName1']['instanceId'] == 'instance1':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 52

def test_case52():

        # to create entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata43),
                      headers=headers)
    print r.content
    print r.status_code

        # to fetch and verify instanceId
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:C002'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['brandName1']['datasetId'] == 'dataset1':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase 53

def test_case53():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata44),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 54

def test_case54():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata45),
                      headers=headers)
    print r.content
    print r.status_code
    assert r.status_code == 400


# testCase 55

def test_case55():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/ld + json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata46),
                      headers=headers)
    print r.content
    assert r.status_code == 400


# testCase 56

def test_case56():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata47),
                      headers=headers)
    print r.content
    assert r.status_code == 204


# testCase 57

def test_case57():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata48),
                      headers=headers)
    print r.content
    assert r.status_code == 404


# testCase 58

def test_case58():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/abc+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata47),
                      headers=headers)
    print r.content
    assert r.status_code == 400


# testCase 59

def test_case59():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata49),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 60

def test_case60():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata50),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 61

def test_case61():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata51),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 62

def test_case62():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata52),
                      headers=headers)
    print r.content
    assert r.status_code == 404


# testCase 63

def test_case63():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata53),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 64

def test_case64():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata54),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 65

def test_case65():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata55),
                      headers=headers)
    print r.content
    assert r.status_code == 404


# testCase 66

def test_case66():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata56),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 67

def test_case67():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata57),
                      headers=headers)
    print r.content
    assert r.status_code == 207


# testCase 68

def test_case68():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata58),
                      headers=headers)
    print r.content
    assert r.status_code == 404


# testCase 69

def test_case69():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata59),
                      headers=headers)
    print r.content

    # to fetch the corresponding entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A0001'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    assert r.status_code == 200


# testCase 70

def test_case70():

    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata60),
                      headers=headers)
    print r.content

        # to fetch the corresponding first entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A0101'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content

        # to fetch the corresponding second entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A9090'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content

    assert r.status_code == 200


# testCase71

def test_case71():

    # create entities
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata61),
                      headers=headers)
    print r.content

        # to fetch the corresponding first entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A0210'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"])

    print '\nchecking the value of  brandName attribute on entity creation'
    if resp['id'] == 'urn:ngsi-ld:Vehicle:A0210':
        print resp['brandName']['value']

        # create entities
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata62),
                      headers=headers)
    print r.content

        # to fetch the corresponding first entity
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A0210'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.get(url, headers=headers)
    print r.content
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)

        # print(resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"])

    print '\nchecking the value of  brandName attribute after entity update'
    if resp['id'] == 'urn:ngsi-ld:Vehicle:A0210':
        print resp['brandName']['value']

    assert r.status_code == 200


# testCase72

def test_case72():
    url = discoveryIp \
        + '/ngsi9/ngsi-ld/registration/urn:ngsi-ld:Vehicle:A0210'
    r = requests.get(url)
    resp_content = r.content
    resInJson = resp_content.decode('utf8').replace("'", '"')
    resp = json.loads(resInJson)
    print resp
    if resp['ID'] == 'urn:ngsi-ld:Vehicle:A0210':
        print '\nValidated'
    else:
        print '\nNot Validated'
    print r.status_code
    assert r.status_code == 200


# testCase73

def test_case73():

    # time.sleep(3)
        # to create subscription

    url = brokerIp + '/ngsi-ld/v1/subscriptions/'
    headers = {'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.subdata63),
                      headers=headers)
    print r.content
    print r.status_code

        # Update entity to fire notification
    # time.sleep(3)

    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json'}
    r = requests.post(url, data=json.dumps(ld_data.subdata64),
                      headers=headers)
    print r.status_code

        # to validate
    # time.sleep(3)

    url = 'http://0.0.0.0:8888/validateNotification'
    r = requests.post(url,
                      json={'subscriptionId': 'urn:ngsi-ld:Subscription:Upsert'
                      })
    print r.content
    assert r.status_code == 200


# test if header content-Type application/json is allowed or not

def test_case74():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.testData74),
                      headers=headers)
    assert r.status_code == 204


# test if header content-Type is application/ld+json then the link header should not be persent in request

def test_case75():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/ld+json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.testData74),
                      headers=headers)
    assert r.status_code == 404


# test if Allowd Content-Type are only appliation/json and application/ld+json

def test_case76():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application1/ld1+json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.testData74),
                      headers=headers)
    print r.status_code
    assert r.status_code == 400


# test create and get the entity in openiot FiwareService

def test_case77():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.testData74),
                      headers=headers)

        # print(r.status_code)

    url = brokerIp + '/ngsi-ld/v1/entities/' \
        + 'urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


def test_case78():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application1/ld1+json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.testData74),
                      headers=headers)

        # print(r.status_code)

    headers = {
        'Content-Type': 'application1/ld1+json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiott',
        'fiware-servicepath': 'test',
        }

    url = brokerIp + '/ngsi-ld/v1/entities/' \
        + 'urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 404


# To test upsert Api support only array of entities

def test_case79():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.testData75),
                      headers=headers)

        # print(r.status_code)

    assert r.status_code == 500


def test_case80():
    upsertURL = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Integration': 'NGSILDBroker',
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    rUpsert = requests.post(upsertURL,
                            data=json.dumps(ld_data.upsertCommand80),
                            headers=headers)
    subscribeURL = brokerIp + '/ngsi-ld/v1/subscriptions/'
    rSubscribe = requests.post(subscribeURL,
                               data=json.dumps(ld_data.subData80),
                               headers=headers)
    print rSubscribe.status_code
    time.sleep(5)
    getURL = accumulatorURl + '/validateupsert'
    rget = requests.get(getURL)
    print rget.content
    assert rget.content == '200'


# To test get Entity by Eid from broker if FiwareService is provided

def test_case81():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# To test get Entity form broker By Eid if FiwareService is not provided

def test_case82():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# To test get all Entity from broker if FiwareService is provided

def test_case83():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url,
                      data=json.dumps(ld_data.upsertMultipleCommand),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities?type=Device'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# To test get Entity form broker if FiwareService is not provided

def test_case84():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url,
                      data=json.dumps(ld_data.upsertMultipleCommand),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities?type=Device'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# Test if registration of entity is available in discovery or not if fiwareService is provided

def test_case85():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = discoveryIp \
        + '/ngsi9/registration/urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# test if registration of entity is available in discovery or not if fiwareService is not provided

def test_case86():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = discoveryIp \
        + '/ngsi9/registration/urn:ngsi-ld:Device:water001'
    r = requests.get(url, headers=headers)
    assert r.status_code == 200


# test response of discovery if Entity does not exist in the disovery with fiwareService

def test_case87():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = discoveryIp \
        + '/ngsi9/registration/urn:ngsi-ld:Device:water0010'
    r = requests.get(url, headers=headers)
    assert r.content == 'null'


# test response of discovery if Entity does not exist in the disovery without fiwareService

def test_case88():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.upsertCommand),
                      headers=headers)
    url = discoveryIp \
        + '/ngsi9/registration/urn:ngsi-ld:Device:water0010'
    r = requests.get(url, headers=headers)
    assert r.content == 'null'


# test Delete Entity from thinbroker if FiwareService is provided

def test_case89():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    r = requests.post(url, data=json.dumps(ld_data.DelData),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A109'
    rget = requests.get(url, headers=headers)
    assert rget.status_code == 200
    delURL = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A109'
    r = requests.delete(delURL, headers=headers)
    rget = requests.get(url, headers=headers)
    assert rget.status_code == 404


# test creation of Entity in one fiwareService and delete the Entity with same id in different fiwareService

def test_case90():
    url = brokerIp + '/ngsi-ld/v1/entityOperations/upsert'
    headers = {'Content-Type': 'application/json',
               'Accept': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
    r = requests.post(url, data=json.dumps(ld_data.test89),
                      headers=headers)
    url = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Device:test89'
    rget = requests.get(url, headers=headers)
    assert rget.status_code == 200
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/ld+json',
        'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"',
        'fiware-service': 'openiot',
        'fiware-servicepath': 'test',
        }
    delURL = brokerIp + '/ngsi-ld/v1/entities/urn:ngsi-ld:Device:test89'
    rDel = requests.delete(delURL, headers=headers)
    assert rDel.status_code == 404
    assert rget.status_code == 200

