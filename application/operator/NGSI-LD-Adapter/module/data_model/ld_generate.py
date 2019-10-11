import pickle
import os,sys
sys.path.append('/opt/ngsildAdapter/module')
import datetime
import io
import requests
import sys
import json
from consts import constant
from common_utilities.LogerHandler import Handler
from metadata_jsonToString import metadata_converter
#  class for ngsi data creation

class ngsi_data_creation:
    def __init__(self,data):
        self.data=data 
        self.context={}
        logger_obj=Handler()
        self.logger=logger_obj.get_logger()

    # method for converting the ngsiv1 data ngsi-ld data
    def manage_entity(self,data3):
        if data3.has_key('entityId')==True:
            entity_data=data3['entityId']
            self.entity_data_type=entity_data['type']
            self.entity_data_id=entity_data['id']
        else:
            self.entity_data_type=data3['type']
            self.entity_data_id=data3['id']

    # creating ngsiv1 to ngsild
    def get_ngsi_ld(self):
        self.logger.info("Start converting ngsiv1 data to NGSI-LD")
        data=self.data
        data2=data['contextResponses'][0]
        data3=data2['contextElement']
        self.manage_entity(data3)
        relations=[]
        relation_url={}
        relation_url[self.entity_data_type]=constant.brand_url+self.entity_data_type
        attribute_data=data3['attributes']
        length=len(attribute_data)
        for i in range(length):
            attribute_data=data3['attributes'][i]
            attribute_data_contextvalue=attribute_data['value']
            attribute_data_type=attribute_data['type']
            attribute_data_name=attribute_data['name']
            relation_url[attribute_data_name]=constant.brand_url+attribute_data_name
        relations.append(relation_url)
        self.context['@context']=relations
        attribute_data=data3['attributes'][0]
        for i in range(length):
            attribute_data=data3['attributes'][i]
            attribute_name=attribute_data['name']
            brand_type={}
            brand_type['type']="Property"
            brand_type['value']=attribute_data['value']
            self.context[attribute_name]=brand_type
        self.context['id']=constant.id_value+self.entity_data_id
        self.context['type']=self.entity_data_type
        if data3.has_key('domainMetadata')==True:
            meta_data=data3['domainMetadata'][0]
            meta_data_obj=metadata_converter(self.context,meta_data)
            self.context=meta_data_obj.get_converted_metadata()
            return self.context
        else:
            return self.context

    # creating entity id 

    def get_entityId(self):
        self.logger.info("start creating Entity")
        data=self.data
        data2=data['contextResponses'][0]
        data3=data2['contextElement']
        self.manage_entity(data3)
        self.entity_id=constant.id_value+self.entity_data_id
        self.logger.info("creation of entity has been done")
        return self.entity_id
