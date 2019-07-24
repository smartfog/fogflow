import json
import unittest

from data_models.fogflow_entity import FogFlowEntity
from utils.color_logger import colorlog

logger = colorlog.getLogger('FogFlow entity creation')
logger.info("Testing FogFlow Entity Creation")


# TODO Tests is not comparing dicts as expected
class FogFlowEntityTest(unittest.TestCase):
    entity = FogFlowEntity('HOP30aea4ca71fe:smartspot', 'smartspot', 'false',
                           [
                               {
                                   "name": 'temperature',
                                   "type": 'Number',
                                   "value": '32.3'
                               },
                               {
                                   "name": 'altitude',
                                   "type": 'Number',
                                   "value": '78.796997'
                               },
                               {
                                   "name": 'latitude',
                                   "type": 'Number',
                                   "value": '-1.32'
                               },
                               {
                                   "name": 'longitude',
                                   "type": 'Number',
                                   "value": '30.3'
                               },
                           ]
                           )
    entity_part1 = FogFlowEntity('HOP30aea4ca71fe:smartspot', 'smartspot', 'false',
                                 [
                                     {
                                         "name": 'temperature',
                                         "type": 'Number',
                                         "value": '32.3'
                                     },
                                     {
                                         "name": 'altitude',
                                         "type": 'Number',
                                         "value": '78.796997'
                                     }
                                 ]
                                 )
    entity_part2 = FogFlowEntity('HOP30aea4ca71fe:smartspot', 'smartspot', 'false',
                                 [
                                     {
                                         "name": 'latitude',
                                         "type": 'Number',
                                         "value": '-1.32'
                                     }
                                 ]
                                 )
    entity_part3 = FogFlowEntity('HOP30aea4ca71fe:smartspot', 'smartspot', 'false',
                                 [
                                     {
                                         "name": 'longitude',
                                         "type": 'Number',
                                         "value": '30.3'
                                     },
                                 ]
                                 )

    desired_entity = {
        "entityId": {
            "id": "HOP30aea4ca71fe:smartspot",
            "type": "smartspot",
            "isPattern": False
        },
        "attributes": [
            {
                "name": "temperature",
                "type": "Number",
                "contextValue": "32.3"
            },
            {
                "name": "altitude",
                "type": "Number",
                "contextValue": "78.796997"
            }
        ],
        "domainMetadata": [
            {
                "name": "location",
                "type": "point",
                "value": {
                    "latitude": -1.32,
                    "longitude": 30.3
                }
            }]
    }
    logger.info("created entity\n{}".format(json.dumps(entity.to_json(), indent=4)))
    logger.info("desired entity\n{}".format(json.dumps(desired_entity, indent=4)))

    def setUp(self):
        pass

    def fogflowEntityCreation(self):
        self.assertDictEqual(FogFlowEntityTest.entity.to_json(), FogFlowEntityTest.desired_entity)

    def fogflowEntityCreationParts(self):
        dicti = FogFlowEntityTest.entity_part1.to_json()
        dicti.update(FogFlowEntityTest.entity_part2.to_json())
        dicti.update(FogFlowEntityTest.entity_part3.to_json())
        self.assertDictEqual(dicti, FogFlowEntityTest.desired_entity)

        self.assertEqual(json.dumps(dicti), json.dumps(FogFlowEntityTest.desired_entity))


if __name__ == '__main__':
    unittest.main()
