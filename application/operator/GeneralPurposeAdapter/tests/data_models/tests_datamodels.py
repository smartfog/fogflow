import json
import unittest

import colorlog

from data_models.fiware_entity import FiwareEntity
from data_models.fogflow_entity import FogFlowEntity

logger = colorlog.getLogger('testsDataModelsFogFlow')
fiware_entity = \
    {
        "id": "smartspot:HOP240ac4066f62",
        "type": "smartspot",
        "isPattern": "false",
        "attributes": [
            {
                "name": "temperature",
                "type": "Number",
                "value": 21.233592694220384
            },
            {
                "name": "humidity",
                "type": "Number",
                "value": 58.148624256789
            }
        ]
    }

fogflow_entity = \
    {
        "entityId":
            {
                "id": "smartspot:HOP240ac4066f62",
                "type": "smartspot",
                "isPattern": "false",
            },
        "attributes":
            [
                {
                    "name": "temperature",
                    "type": "Number",
                    "contextValue": 21.233592694220384
                },
                {
                    "name": "humidity",
                    "type": "Number",
                    "contextValue": 58.148624256789
                }
            ],
        "domainMetadata": []
    }


class TestDataModels(unittest.TestCase):
    def setUp(self):
        print("Setting all up")
        self.fogflowEntity = FogFlowEntity(fiware_entity['id'],
                                           fiware_entity['type'],
                                           fiware_entity['isPattern'],
                                           fiware_entity['attributes'])
        self.fiwareEntity = FiwareEntity(fogflow_entity['entityId'], fogflow_entity['attributes'])

    def test_data_model_fiware_to_fogflow(self):
        self.assertEqual(self.fogflowEntity.to_json(), json.dumps(fogflow_entity), "Entities should be equals")

    def test_data_model_fogflow_to_fiware(self):
        self.assertEqual(self.fiwareEntity.to_json(), json.dumps(fiware_entity), "Entities should be equals")


if __name__ == '__main__':
    unittest.main()
