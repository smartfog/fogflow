import ConfigParser
class config_data:
    def __init__(self):
        self.config = ConfigParser.ConfigParser()
        self.config.readfp(open(r'config/config.ini'))
    def get_entity_url(self):
        entity_url= self.config.get('My Section', 'ngsi-ld-broker')
        return entity_url
    def get_fogflow_subscription_endpoint(self):
        fogflow_subscription_url=self.config.get('My Section', 'fogflow_subscription_endpoint')
        return fogflow_subscription_url
        
