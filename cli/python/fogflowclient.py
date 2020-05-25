import requests
import json


class ContextEntity:
    def __init__(self):
        self.id = ""
        self.type = ""        
        self.attributes = {}
        self.metadata = {}
    
    def toJSON(self):
        ctxElement = {}
    
        ctxElement['entityId'] = {}
        ctxElement['entityId']['id'] = self.id
        ctxElement['entityId']['type'] = self.type        
        ctxElement['entityId']['isPattern'] = False        
        
        ctxElement['attributes'] = self.attributes
        ctxElement['metadata'] = self.metadata

        return json.dumps(ctxElement)
        
class FogFlowClient:
    def __init__(self, url):
        self.fogflowURL = url
    
    def sendIntent(self, topology):
        intent = {}
        intent['topology'] = topology
        intent['priority'] = {
            'exclusive': False,
            'level': 50
        };        
        intent['geoscope'] = {
            'scopeType': "global",
            'scopeValue': "global"
        };    
    
        intentEntity = ContextEntity()
                
        intentEntity.id = "ServiceIntent." + topology
        intentEntity.type = "ServiceIntent"          
        
        intentEntity.attributes["status"] =  {'type': 'string', 'value': 'enabled'}                
        intentEntity.attributes["intent"] =  {'type': 'object', 'value': intent}
        intentEntity.metadata["location"]  = {
            'type': 'global',
            'value': 'global'
        };    
                     
        self.put(intentEntity)
       
    def get(self, entityID):
        print(entityID)
        return True        
    
    def put(self, ctxEntity):
        headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}
                
        response = requests.post(self.fogflowURL + '/updateContext', data=ctxEntity.toJSON(), headers=headers)
        if response.status_code != 200:
            print('failed to update context')
            print(response.text)
            return False    
        else:
            print('update context')        
            return True

        
    def delete(self, entityID):
        print(entityID)
        return True      



