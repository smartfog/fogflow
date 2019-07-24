from utils.color_logger import colorlog

from data_models.fogflow_entity import FogFlowEntity
from operations.http_request import HTTPRequest

logger = colorlog.getLogger('FogFlow UpdateContext')


class FogFlowUpdateContext:
    def __init__(self, url: str, request_body: dict):
        self.url = url + "/ngsi10/updateContext"
        self.header = {'Content-Type': "application/json"}
        self.body = {"contextElements": []}
        entities = []
        if request_body.get('contextElements'):
            for entity in request_body['contextElements']:
                _id = entity.get('id')
                _type = entity.get('type')
                is_pattern = entity.get('isPattern')

                attributes = []
                if entity.get('attributes'):
                    attributes = entity.get('attributes')
                else:
                    logger.warning("Received entity without attributes")

                entities.append(FogFlowEntity(_id, _type, is_pattern, attributes).to_json())

        self.body['contextElements'] = entities
        self.body['updateAction'] = 'UPDATE'
        self.request = HTTPRequest(url=self.url, header=self.header, body=self.body)
        self.response = self.request.post()

    def get_response(self):
        return self.response
