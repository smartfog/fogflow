import colorlog

# from utils.color_logger import logger

logger = colorlog.getLogger("NGSI-Parser")

"""
ngsi_parser.py: esta clase es la encargada de hacer la transformación entre NGSIs, se compone de dos funciones, 
context_element_2_json_object() y json_object_2_context_element(), la primera está encargada de convertir el mensaje 
para que FogFlow pueda entenderlo y aceptarlo, y la segunda se encarga de convertir el mensaje ya transformado en un 
objeto contextElement, objeto que se envía en el cuerpo del mensaje en formato JSON.

ngsi_parser.py: this class is responsible for making the transformation between NGSIs, it consists of two functions, 
context_element_2_json_object() and json_object_2_context_element(), the former is in charge of converting the message 
so that FogFlow can understand and accept it, and the second is responsible for converting the message already 
transformed into an object contextElement, object that is sent in the body of the message in JSON format.
"""


def json_object_2_context_element(context_object):
    """
    this method is responsible for converting the message already transformed into an object contextElement,
    object that is sent in the body of the message in JSON format.
    :param context_object:
    :return:
    """

    # logger.info('convert json object to context element\n {}'.format(json.dumps(context_object, indent=4)))
    context_element = \
        {
            'entityId': context_object['entityId'],
            'attributes': []
        }

    if 'attributes' in context_object:
        for key in context_object['attributes']:
            context_element['attributes'].append(
                {'name': key['name'], 'type': key['type'], 'contextValue': key['value']})

    context_element['domainMetadata'] = []
    if 'metadata' in context_object:
        for key in context_object['metadata']:
            meta = context_object['metadata'][key]
            context_element['domainMetadata'].append({'name': key, 'type': meta['type'], 'value': meta['value']})

    return context_element


def context_element_2_json_object(entity):
    """
    this method is for converting the message to json object so that FogFlow can understand and accept it
    :param entity:
    :return:
    """
    json_object = {
        "contextElements": [{
            "entityId":
                {
                    "id": entity["contextElements"][0]["id"],
                    "type": entity["contextElements"][0]["type"],
                    "isPattern": False
                },
            "attributes": {}
        }]
    }

    """
    json_object = {}
    json_object['contextElements']['entityId']['id'] = element["contextElements"][0]["id"]
    json_object['contextElements']['entityId']['type'] = element["contextElements"][0]["type"]
    """
    json_object["contextElements"][0]["attributes"] = []
    if "attributes" in entity["contextElements"][0]:
        for attr in entity["contextElements"][0]["attributes"]:
            json_object["contextElements"][0]["attributes"].append(
                {
                    "name": attr["name"],
                    "type": attr["type"],
                    "contextValue": attr["value"]
                })
    # json_object = {'contextElements': json_object, "metadata": {}}
    if "domainMetadata" in entity:
        for meta in entity["domainMetadata"]:
            json_object["metadata"][meta["name"]] = \
                {
                    "type": meta["type"],
                    "value": meta["value"]
                }
    json_object["updateAction"] = "UPDATE"

    """
    json_object["domainMetadata"] = []
    json_object["domainMetadata"].append(
        {
            "name": "location",
            "type": "point",
            "value": {
                "latitude": 49.406393,
                "longitude": 8.684208
            }
        }
    )"""
    # logger.info("json object parsed to FogFlow format like: \n{}".format(json.dumps(json_object, indent=4)))
    return json_object
