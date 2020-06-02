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

    # synchronized remote call
    def remoteCall(self, serviceTopology):
        intentId = sendIntend(serviceTopology)
        
        result = self.fetchResult
        
        removeIntent(intentId)
        
    # asynchronize 
    def execute(self, serviceTopology, callback):
        sendIntend(serviceTopology)
    
    def sendIntent(self, serviceTopology):
        intent = {}
        intent['topology'] = serviceTopology
        intent['priority'] = {
            'exclusive': False,
            'level': 50
        };        
        intent['geoscope'] = {
            'scopeType': "global",
            'scopeValue': "global"
        };            
        
        headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}                
        response = requests.post(self.fogflowURL + '/intent', data=json.dumps(intent), headers=headers)
        if response.status_code != 200:
            print('failed to update context')
            print(response.text)
            return False    
        else:
            print('send intent')        
            return True    
       
    def removeIntent(self, intentEntityId):
        paramter = {}
        paramter['id'] = intentEntityId
    
        headers = {'Accept' : 'application/json', 'Content-Type' : 'application/json'}                
        response = requests.delete(self.fogflowURL + '/intent', data=json.dumps(paramter), headers=headers)
        if response.status_code != 200:
            print('failed to remove intent')
            print(response.text)
            return False    
        else:
            print('remove intent')        
            return True    

    def fetchResult(self, serviceTopology)
        
        return True
        
    def subscribeResult(self, serviceTopology, callback)
        
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

    def get(self, entityID):
    
        print(entityID)
        
        return True  
        
                
    def delete(self, entityID):
        print(entityID)
        return True      

    def query(self, eid, eType):
        print(entityID)
        return True      

