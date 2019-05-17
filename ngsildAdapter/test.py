#Tests to be add
#from module.data_model.ld_generate import ngsi_data_creation
#from test_data_models import test_data_model
#import module
#import data_module
#import ld_generate
from module.data_model import ld_generate
import json
import unittest
import requests
import sys 
#sys.path.insert(0, '../module/data_model')
#from ld_generate import ngsi_data_creation

#from root.TRANSFORMER.Next_transform.fogflow.ngsildAdapter.module.data_mode.ld_generate import ngsi_data_creation
#import o
ngsi_data= \
 {
       'originator':u'',
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

convert_data_output= \
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
                'type': 'Vehicle',
                'id': 'urn:ngsi-ld:Car31'
 }


id_value="urn:ngsi-ld:Car31"

class TestStringMethods(unittest.TestCase):
    def setUp(self):
        pass
    def test_get_ngsi_ld(self):
        obj=ngsi_data_creation(ngsi_data)
        result_data=obj.get_ngsi_ld()
        self.assertEqual(json.dumps(result_data),json.dumps(convert_data_output))
    def test_get_entityId(self):
        obj=ngsi_data_creation(ngsi_data)
        entity_id=obj.get_entityId()
        self.assertEqual(entity_id,id_value)
if __name__ == '__main__':
    unittest.main()


