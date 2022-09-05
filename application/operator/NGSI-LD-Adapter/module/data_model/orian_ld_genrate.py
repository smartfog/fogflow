from common_utilities.LogerHandler import Handler
from consts import constant
import logging
import json
import sys
import requests
import io
import datetime
import pickle
import sys
import os
sys.path.append('/opt/ngsildAdapter/module')


class orian_convert_data:
    def __init__(self, data):
        self.data = data
        self.context = {}

    def url_conversion(self, keys, relation_url):
        for attr in keys:
            relation_url[attr] = constant.brand_url+attr
        return relation_url

    def add_entity(self, data_keys, data2):
        data = data2
        for keys in data_keys:
            data3 = data[keys]
            if data3.has_key('metadata'):
                del data3['metadata']
            self.context[keys] = data[keys]

    def get_keys(self, data):
        keys = data.keys()
        keys.remove('id')
        keys.remove('type')
        return keys

    def get_data(self):
        logger_obj = Handler()
        logger = logger_obj.get_logger()
        logger.info("Start converting ngsiv2 data to NGSI-LD")
        data = self.data
        data1 = data['data']
        data2 = data1[0]
        relations = []
        relation_url = {}
        relations.append(constant.context_url)
        entity_data_type = data2['type']
        relation_url[entity_data_type] = constant.brand_url+entity_data_type
        data_keys = self.get_keys(data2)
        relation_url = self.url_conversion(data_keys, relation_url)
        relations.append(relation_url)
        self.context['@context'] = relations
        self.context['id'] = constant.id_value+data2['id']
        self.context['type'] = data2['type']
        self.add_entity(data_keys, data2)
        logger.info("Data has been converted to NGSI-LD")
        return self.context

    def get_entityId(self):
        logger_obj = Handler()
        logger = logger_obj.get_logger()
        logger.info("Creating Entity")
        data = self.data
        data1 = data['data']
        data2 = data1[0]
        entity_id = constant.id_value+data2['id']
        logger.info("Entity has been created")
        return entity_id
