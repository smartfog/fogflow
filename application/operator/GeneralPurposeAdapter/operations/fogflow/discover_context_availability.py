from operations.http_request import HTTPRequest


# SINCE THIS IS DEPLOYED WITH DOCKER SWARM THIS IS NOT NEEDED, maybe with Adapter as Operator.
class FogFlowDiscoverContextAvailability:
    def __init__(self, url: str, latitude: float, longitude: float, limit=1):
        self.url = url + "/ngsi9/discoverContextAvailability"
        self.header = {'Content-Type': "application/json"}
        self.body = {
            "entities": [
                {"type": "IoTBroker", "isPattern": True}
            ],
            'restriction': {
                "scopes": [{"type": "nearby", "value": {"latitude": latitude, "longitude": longitude, "limit": limit}}]
            }
        }
        self.request = HTTPRequest(url=self.url, header=self.header, body=self.body)
        self.response = self.request.post()


""" EXAMPLE:
curl -iX POST \
  'http://localhost:8071/ngsi9/discoverContextAvailability' \
  -H 'Content-Type: application/json' \
  -d '
    {
       "entities":[
          {
             "type":"IoTBroker",
             "isPattern":true
          }
       ],
       "restriction":{
          "scopes":[
             {
                "type":"nearby",
                "value":{
                   "latitude":35.692221,
                   "longitude":139.709059,
                   "limit":1
                }
             }
          ]
       }
    }
"""
