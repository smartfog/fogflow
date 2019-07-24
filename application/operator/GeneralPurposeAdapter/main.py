import json
import os

from flask import (Flask, request, Response)

from operations.fogflow.update_context import FogFlowUpdateContext
from settings import cfg
from utils.color_logger import *

logger = colorlog.getLogger("Main")
# Create the application instance
application = Flask(__name__)

########################################################################################################################
#                                            SELECTED DESTINIES
########################################################################################################################
url_iotbroker = "http://" + os.getenv('IOT_BROKER_IP', cfg['iot_broker']["server_uri"]) + ':' + os.getenv(
    'IOT_BROKER_PORT', cfg['iot_broker']["server_port"])
########################################################################################################################
logger.info("Selected IoT Broker: {}".format(url_iotbroker))


########################################################################################################################
#                                            FOGFLOW ENTRY POINTS
########################################################################################################################
@application.route('/v1/updateContext', methods=['POST'])
def update_context():
    """
    Translate message from IoT Agent Broker to be recognized by FogFlow
    :return:"""
    data = request.json
    if data.get('contextElements'):
        for entity in data['contextElements']:
            logger.info("Received a updateContext from device: {} with attributes: {}".format(
                json.dumps(entity['id']), entity.get('attributes')))

        # create a FogFlow NGSI updateContext request and POST it.
        update_request = FogFlowUpdateContext(url=url_iotbroker, request_body=data)
        response = update_request.get_response()
        if response.status_code != 200:
            logger.error(
                "Error sending updateContext to FogFlow {}{}".format(response.content,
                                                                     response.status_code))

            if response is not None and response.content is not None:
                return Response(response=update_request.response.content, content_type="application/json",
                                status=update_request.response.status_code)
        else:
            response = {
                "contextResponses": [
                    {
                        "statusCode": {
                            "code": "200",
                            "reasonPhrase": "OK"
                        }
                    }
                ]
            }
            return Response(response=json.dumps(response), content_type="application/json", status=200)
    # BAD REQUEST
    return Response(response=json.dumps(None), content_type="application/json", status=400)


@application.route('/NGSI9/registerContext', methods=['POST'])
def register_context():
    """
    When receive registerContext we register the new device and assign an iot broker.
    :return:
    """
    # WE SHOULD NOT RECEIVE A REGISTER CONTEXT IF WE ARE USING DYNAMIC PROVISIONING.
    logger.info("Received register context with\n{}".format(json.dumps(request.json)))
    # BUT IF WE RECEIVE A REGISTER CONTEXT WE RESPOND WITH THE EXPECTED RESPONSE TO IOT AGENT
    import secrets
    registration_id = secrets.token_hex(12)
    resp = \
        {
            "registrationId": registration_id,
            "duration": json.loads(request.data).get('duration')
        }
    resp_status = 200
    return Response(response=resp, content_type="application/json", status=resp_status)


def device_registration(entity):
    """
    ############################################################################################################
    if entity.get('attributes') and not DeviceRegistration.is_registered(entity['id']):
        for attribute in entity['attributes']:
            if attribute.get('name') == 'location':
                # logger.info("Device location: {}".format(json.dumps(attribute)))
                latitude = attribute.get('value').get('coordinates')[0]
                longitude = attribute.get('value').get('coordinates')[1]
                logger.info(
                    "Getting the nearest IoT Broker for location latitude: {} longitude: {}".format(latitude,
                                                                                                    longitude))
                # LOOK FOR A CLOSE IOT BROKER
                resp = FogFlowDiscoverContextAvailability(url_iotdiscovery, latitude, longitude)
                logger.info("Discovery response: {}{}".format(resp.url, json.dumps(resp.body, indent=4)))
                # REGISTER DEVICE INTO THIS SERVICE.
                # logger.info("Registering device in this service")
                # DeviceRegistration.add_device(entity['id'],)
    ############################################################################################################
    """


########################################################################################################################
#                                            FIWARE ENTRY POINTS
########################################################################################################################
# V2 notify from Orion
@application.route('/v2/notify', methods=['POST'])
def notify():
    """
    Create/Update a entity in FogFlow from a Orion V2 Notify Context.
    :return:"""
    request_body = request.json
    if request_body.get('data'):
        entities = []
        for entity_data in request_body['data']:
            attributes = []
            _id = entity_data.pop("id")
            _type = entity_data.pop('type')

            # get the attributes
            for attribute, value in entity_data.items():
                attributes.append({"name": attribute, "type": value.get('type'), "value": value.get('value')})

            entity = {
                "id": _id,
                "type": _type,
                "isPattern": False,
                "attributes": attributes
            }
            # append current entity to the entity list
            entities.append(entity)

        # add the entity list to the request
        entity_request = \
            {
                "contextElements": entities
            }

        # create a FogFlow NGSI updateContext request and POST it.
        update_request = FogFlowUpdateContext(url=url_iotbroker, request_body=entity_request)
        response = update_request.get_response()
        if response.status_code != 200:
            logger.error(
                "Error notifying FogFlow {}{}".format(update_request.response.content,
                                                      update_request.response.status_code))
        return Response(response=json.dumps(""), content_type="application/json", status=200)
    # BAD REQUEST
    return Response(response=json.dumps(None), content_type="application/json", status=400)


# If we're running in stand alone mode, run the application
if __name__ == '__main__':
    logger.info("Starting Server...")
    # Listen port will be on 1026
    application.run(host='0.0.0.0', port=os.getenv('LISTEN_PORT', 1026), debug=True, use_reloader=False)
