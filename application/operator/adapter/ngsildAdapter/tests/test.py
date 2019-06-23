import sys, os
sys.path.append('/root/TRANSFORMER/Next_transform/fogflow/ngsildAdapter/module')
from common_utilities.rest_client import Rest_client
from common_utilities import rest_client
from data_model.ld_generate  import ngsi_data_creation
import json
import unittest
import requests
import mock 
from mock import patch
from consts import constant
from common_utilities.config import config_data
ngsi_data=\
 {
       'originator':'',
           'subscriptionId': '195bc4c6-882e-40ce-a98f-e9b72f87bdfd',
           'contextResponses':
              [
                    {
                      'contextElement': {'attributes':
                            [
                                   {
                                          'contextValue': 'ford5',
                                           'type': 'string',
                                           'name': 'brand40'
                                   },
                                   {
                                          'contextValue': 'ford6',
                                          'type': 'string',
                                          'name': 'brand50'
                                   }
                        ],
                        'entityId':
                             {
                                   'type': 'Car',
                                    'id': 'Car31',
                                    'isPattern':True
                                 },
                      'domainMetadata':
                       [
                             {
                                  'type': 'point',
                                  'name': 'location',
                                  'value':
                                     {
                                               'latitude': 49.406393,
                                                'longitude': 8.684208
                                        }
                                }
                          ]
                  },
                  'statusCode':
                   {
                           'code': 200,
                           'reasonPhrase': 'OK'
                   }
         }
   ]
}
convert_data_output=\
 {
   '@context':
        [
                   'https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld',
                        {
                                     'Car':     'http://example.org/Car',
                                     'brand40': 'http://example.org/brand40',
                                     'brand50': 'http://example.org/brand50'
                            }
                ],
                    'brand50':
                        {
                            'type': 'Property',
                            'value': 'ford6'
                        },
                     'brand40':
                         {
                             'type': 'Property',
                             'value': 'ford5'
                         },
                'type': 'Car',
                'id': 'urn:ngsi-ld:Car31'
 }
patch_data_output=\
 {
   '@context':
        [
                   'https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld',
                        {
                                     'Car':     'http://example.org/Car',
                                     'brand40': 'http://example.org/brand40',
                                     'brand50': 'http://example.org/brand50'
                            }
                ],
                    'brand50':
                        {
                            'type': 'Property',
                            'value': 'ford6'
                        },
                     'brand40':
                         {
                             'type': 'Property',
                             'value': 'ford5'
                         },
 }
id_value="urn:ngsi-ld:Car31"

class TestStringMethods(unittest.TestCase):
    '''def setUp(self):
        pass'''
    def test_get_ngsi_ld(self):
        obj=ngsi_data_creation(ngsi_data)
        result_data=obj.get_ngsi_ld()
        self.assertEqual(json.dumps(result_data),json.dumps(convert_data_output))
    def test_get_entityId(self):
        obj=ngsi_data_creation(ngsi_data)
        entity_id=obj.get_entityId()
        self.assertEqual(entity_id,id_value)
    def test_mock_post(self):
        with patch('common_utilities.rest_client.requests.post') as mock_get:
            mock_get.return_value.status_code = 201
            configobj=config_data()
            entity_url=configobj.get_entity_url()
            url1 =constant.http+entity_url+constant.entity_uri
            payload=convert_data_output
            obj=Rest_client(url1,payload)
            response=obj.post_request()
        self.assertEqual(response.status_code, 201)
    def test_mock_patch(self):
        with patch('common_utilities.rest_client.requests.patch') as mock_get:
            mock_get.return_value.status_code = 204
            obj=ngsi_data_creation(ngsi_data)
            entity_id=obj.get_entityId()
            configobj=config_data()
            entity_url=configobj.get_entity_url()
            url=constant.http+entity_url+constant.entity_uri+entity_id+'/attrs'
            payload=patch_data_output
            obj=Rest_client(url,payload)
            response=obj.patch_request()
        self.assertEqual(response.status_code, 204)
if __name__ == '__main__':
    unittest.main()


