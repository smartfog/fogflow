import json

import colorlog
import requests
from flask import Response

logger = colorlog.getLogger('HTTPRequest')


class HTTPRequest(object):
    def __init__(self, url: str, header: dict, body: dict):
        self.url = url
        self.header = header
        self.body = body

    def get_url(self):
        return self.url

    def get_header(self):
        return self.url

    def post(self):
        try:
            logger.info("Sending the next post request: \nurl: {}\nheader: {}\nbody: {}".format(self.url, self.header,
                                                                                                json.dumps(self.body,
                                                                                                           indent=4)))
            return requests.post(url=self.url, data=json.dumps(self.body), headers=self.header)
        except requests.exceptions.ConnectionError as ce:
            logger.error("Error trying to connect to {}\n{}".format(self.url, ce))
            return Response(status=404)
