import json

import colorlog
import requests

from settings import cfg

"""
request.py: en esta clase se produce el envío del mensaje transformado al IoT Broker mediante una llamada al método 
request() especificando la IP del IoT Broker y el puerto donde escucha.

request.py: in this class the transformed message is sent to the IoT Broker by means of a call to the method 
request() specifying the IP of the IoT Broker and the port where it listens.
"""

logger = colorlog.getLogger("request")


def forward_context_information(msg):
    """
    send to the IoT Broker the transformed message from ngsi_parser.context_element_2_json_object() method
    :param msg:
    :return:
    """

    send_request = {
        "url": "http://" + cfg['iot_broker']["server_uri"] + ':' + cfg['iot_broker'][
            "server_port"] + '/ngsi10/updateContext',
        "method": 'POST',
        "payload": msg,
        "headers": {
            'Content-Type': 'application/json',
            'cache-control': 'no-cache'
        }
    }
    resp = requests.Response()
    resp.status_code = 404
    # logger.info("sending to FogFlow: \n{}".format(json.dumps(options['payload'], indent=4)))
    try:
        resp = requests.request(send_request['method'], url=send_request["url"],
                                data=json.dumps(send_request['payload']), headers=send_request['headers'])
    except requests.exceptions.ConnectionError as c:
        logger.error("Connection Error trying to connect to {}:\n{}".format(send_request['url'], c))
    # logger.info("FogFlow response: {}".format(resp.json()))

    return resp
