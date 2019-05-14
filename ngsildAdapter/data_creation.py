from flask import Flask, abort, request
import json
import pickle
import datetime
import io
import requests
import constant
class ngsi_data_creation:
    def __init__(self,data):
        self.data=data 
        self.context={}
    def get_ngsi_ld(self):
        data=self.data
        data2=data['contextResponses'][0]
        data3=data2['contextElement']
        #data3=data[data['contextResponses'][0]]
        entity_data=data3['entityId']
        entity_data_type=entity_data['type']
        entity_data_id=entity_data['id']
        relations=[]
        relation_url={}
        relations.append(constant.jsonld_url)
        relation_url[entity_data_type]=constant.brand_url+entity_data_type
        attribute_data=data3['attributes']
        length=len(attribute_data)
        for i in range(length):
            attribute_data=data3['attributes'][i]
            attribute_data_contextvalue=attribute_data['contextValue']
            attribute_data_type=attribute_data['type']
            attribute_data_name=attribute_data['name']
            relation_url[attribute_data_name]=constant.brand_url+attribute_data_name
        relations.append(relation_url)
        self.context['@context']=relations
        #self.context['id']=constant.id_value+entity_data_id
        #self.entity_id=self.context['id']
        #self.context['type']="Vehicle"
        attribute_data=data3['attributes'][0]
        for i in range(length):
            attribute_data=data3['attributes'][i]
            attribute_name=attribute_data['name']
            brand_type={}
            brand_type['type']="Property"
            brand_type['value']=attribute_data['contextValue']
            self.context[attribute_name]=brand_type
        self.context['id']=constant.id_value+entity_data_id
        self.context['type']=entity_data_type
        return self.context     
    def get_entityId(self):
        data=self.data
        data2=data['contextResponses'][0]
        data3=data2['contextElement']
        entity_data=data3['entityId']
        entity_data_id=entity_data['id']
        self.entity_id=constant.id_value+entity_data_id
        return self.entity_id
