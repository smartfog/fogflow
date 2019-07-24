import json
import os

import colorlog
import requests
from flask import (
    Flask,
    Response, request)

from OLD.ngsi.ngsi_parser import context_element_2_json_object
from OLD.request.request import forward_context_information
from settings import cfg

logger = colorlog.getLogger("Main")
# Create the application instance
application = Flask(__name__)
########################################################################################################################
url_orion = "http://" + os.getenv('ORION_IP', cfg['orion']["orion_url"]) + ':' + os.getenv('ORION_PORT',
                                                                                           cfg['orion']["orion_port"])

url_iotbroker = "http://" + os.getenv('IOT_BROKER_IP', cfg['iot_broker']["server_uri"]) + ':' + os.getenv(
    'IOT_BROKER_PORT', cfg['iot_broker']["server_port"])
forward_to_orion = os.getenv("FORWARD_TO_ORION", cfg['orion']["forward_to_orion"]) == 'True'

url_iotdiscovery = "http://" + os.getenv('IOT_DISCOVERY_IP', cfg['iot_discovery']["server_uri"]) + ':' + os.getenv(
    'IOT_DISCOVERY_PORT', cfg['iot_discovery']["server_port"])
########################################################################################################################
logger.info("Forward to Orion: {}".format(forward_to_orion))
if forward_to_orion:
    logger.info("Selected Orion Context Broker: {}".format(url_orion))
logger.info("BackUp IoT Broker: {}".format(url_iotbroker))
logger.info("Selected IoT Discovery: {}".format(url_iotdiscovery))


########################################################################################################################
#                                            FOGFLOW ENTRY POINTS
########################################################################################################################
@application.route('/v1/updateContext', methods=['POST'])
def update_context():
    """
Translate message from IoT Agent Broker to be recognized by FogFlow
:return:"""
    orion_response = requests.Response()
    # -----------------------------------------------------------------------------------------------------------------
    if forward_to_orion:  # if wanna forward to real Orion
        send_message_to_orion('updateContext', request)
    # -----------------------------------------------------------------------------------------------------------------
    forwarded_message = forward_context_information(context_element_2_json_object(json.loads(request.data)))
    if forwarded_message.status_code != 200:
        logger.error("Error forwarding message to FogFlow in updateContext {}".format(forwarded_message.status_code))
        return Response(response=None, content_type="application/json",
                        status=404)

    else:
        logger.info("FogFlow updateContext response: {}".format(json.loads(forwarded_message.content)))
        orion_response = {
            "contextResponses": [
                {
                    "statusCode": {
                        "code": "200",
                        "reasonPhrase": "OK"
                    }
                }
            ]
        }
        return Response(response=json.dumps(orion_response), content_type="application/json", status=200)


@application.route('/NGSI9/registerContext', methods=['POST'])
def register_context():
    """
    When receive registerContext we register the new device and assign an iot broker.
    :return:
    """
    # logger.info("Received registerContext http request")
    if forward_to_orion:  # if wanna forward to real Orion
        send_message_to_orion('registerContext', request)

    import secrets
    registration_id = secrets.token_hex(12)
    resp = \
        {
            "registrationId": registration_id,
            "duration": json.loads(request.data).get('duration')
        }
    resp_status = 200

    return Response(response=resp, content_type="application/json", status=resp_status)


def send_message_to_orion(_type: str, _request: requests):
    session = requests.Session()
    session.headers = {
        "fiware-service": request.headers.get('fiware-service'),
        "fiware-servicepath": request.headers.get('fiware-servicepath'),
        "Cache-Control": "no-cache"
    }
    content_type = {'Content-Type': 'application/json; charset=' + "utf-8"}

    orion_response = requests.Response  # creating a dummy response empty
    if _type == 'registerContext':
        try:
            orion_response = session.post(url_orion + "/v1/registry/registerContext", json="application/json",
                                          data=_request.data, headers=content_type)

        except requests.exceptions.ConnectionError as c:
            logger.error("Error trying to connect to Orion Context Broker. \n{}".format(c))
            return

        if orion_response.status_code != 200:
            logger.error("Orion response error code in registerContext: {}".format(orion_response))

    elif _type == 'updateContext':
        try:
            orion_response = session.post(url_orion + "/v1/updateContext", json="application/json",
                                          data=json.loads(_request.data), headers=content_type)
        except requests.exceptions.ConnectionError as c:
            logger.error("Error trying to connect to Orion Context Broker. \n{}".format(c))
            return

        orion_response = orion_response.content
        if orion_response.get("contextResponses")[0].get('statusCode').get('code') != 200:
            logger.error(
                "Orion response error code in updateContext: \n{}".format(json.dumps(orion_response, indent=4)))

########################################################################################################################
#                                            FIWARE ENTRY POINTS
########################################################################################################################


# If we're running in stand alone mode, run the application
if __name__ == '__main__':
    # print("Starting Server...")
    logger.info("Starting Server...")
    application.run(host='0.0.0.0', port=os.getenv('LISTEN_PORT', 1026), debug=True, use_reloader=False)
