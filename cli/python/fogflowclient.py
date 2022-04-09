import requests
import json
import socketio
from threading import Lock, Thread


class WebSocketClient(Thread):
    def __init__(self, url):
        Thread.__init__(self)
        self.url = url
        self.lock = Lock()
        self.connected = False
        self.subscriptions = []
        self.cbList = {}

        self.sio = socketio.Client()
        #self.sio = socketio.AsyncClient()

    def setCallback(self, subscriptionId, callback):
        self.lock.acquire()
        self.cbList[subscriptionId] = callback
        self.subscriptions.append(subscriptionId)
        if self.connected == True:
            self.sio.emit('subscriptions', self.subscriptions)

        self.subscriptions = []
        self.lock.release()

    def onConnect(self):
        self.lock.acquire()
        self.connected = True
        print('connection established')
        if len(self.subscriptions) > 0:
            sio.emit('subscriptions', self.subscriptions)
        self.subscriptions = []
        self.lock.release()

    def onNotify(self, msg):
        sid = msg['subscriptionID']
        entities = msg['entities']

        self.lock.acquire()
        cb = self.cbList[sid]
        self.lock.release()

        cb(entities)

    def onDisconnect(self):
        self.lock.acquire()
        self.connected = False
        print('disconnected from server')
        self.lock.release()

    def run(self):
        self.sio.on('connect', self.onConnect)
        self.sio.on('disconnect', self.onDisconnect)

        self.sio.on('notify', self.onNotify)

        self.sio.connect(self.url)
        self.sio.wait()

    def stop(self):
        self.sio.disconnect()


class ContextEntity:
    def __init__(self):
        self.id = ""
        self.type = ""
        self.attributes = {}
        self.metadata = {}

    def toContextElement(self):
        ctxElement = {}

        ctxElement['entityId'] = {}
        ctxElement['entityId']['id'] = self.id
        ctxElement['entityId']['type'] = self.type
        ctxElement['entityId']['isPattern'] = False

        ctxElement['attributes'] = []
        for key in self.attributes:
            attr = self.attributes[key]
            ctxElement['attributes'].append(
                {'name': key, 'type': attr['type'], 'value': attr['value']})

        ctxElement['domainMetadata'] = []
        for key in self.metadata:
            meta = self.metadata[key]
            ctxElement['domainMetadata'].append(
                {'name': key, 'type': meta['type'], 'value': meta['value']})

        return ctxElement

    def fromContextElement(self, ctxElement):
        entityId = ctxElement['entityId']

        self.id = entityId['id']
        self.type = entityId['type']

        if 'attributes' in ctxElement:
            for attr in ctxElement['attributes']:
                self.attributes[attr['name']] = {
                    'type': attr['type'], 'value': attr['value']}

        if 'domainMetadata' in ctxElement:
            for meta in ctxElement['domainMetadata']:
                self.attributes[meta['name']] = {
                    'type': meta['type'], 'value': meta['value']}

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
        self.wsclient = WebSocketClient(url)
        self.wsclient.start()

    def __del__(self):
        self.wsclient.stop()

    # synchronized remote call
    def remoteCall(self, serviceTopology):
        response = requests.get(
            self.fogflowURL + '/remoteCall?serviceTopology=' + serviceTopology)
        print(response.text)

    # asynchronize way to trigger a service topology
    def start(self, serviceTopology, callback):
        # issue an intent to trigger the service topology
        response = self.sendIntent('Test')

        print(response)

        # set up the callback function to handle the subscribed results
        refURL = 'http://host.docker.internal:1030'
        sid = self.subscribe(response['outputType'], refURL)
        if sid != '':
            self.wsclient.setCallback(sid, callback)
        else:
            print('failed to create the subscription')

        return response['id']

    # stop the triggered service topology
    def stop(self, sid):
        self.removeIntent(sid)

    def sendIntent(self, serviceTopology):
        intent = {}
        intent['topology'] = serviceTopology
        intent['priority'] = {
            'exclusive': False,
            'level': 50
        }
        intent['geoscope'] = {
            'scopeType': "global",
            'scopeValue': "global"
        }

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.post(
            self.fogflowURL + '/intent', data=json.dumps(intent), headers=headers)
        if response.status_code != 200:
            print('failed to update context')
            print(response.text)

        return response.json()

    def removeIntent(self, intentEntityId):
        paramter = {}
        paramter['id'] = intentEntityId

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.delete(
            self.fogflowURL + '/intent', data=json.dumps(paramter), headers=headers)
        if response.status_code != 200:
            print('failed to remove intent')
            print(response.text)
            return False
        else:
            print('remove intent')
            return True

    def put(self, ctxEntity):
        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}

        updateCtxReq = {}
        updateCtxReq['contextElements'] = []
        updateCtxReq['contextElements'].append(ctxEntity.toContextElement())
        updateCtxReq['updateAction'] = 'UPDATE'

        response = requests.post(self.fogflowURL + '/ngsi10/updateContext',
                                 data=json.dumps(updateCtxReq), headers=headers)
        if response.status_code != 200:
            print('failed to update context')
            print(response.text)
            return False
        else:
            return True

    def getById(self, entityId):
        queryReq = {}
        queryReq['entities'] = []

        idObj = {}
        idObj['id'] = entityId
        idObj['isPattern'] = False
        queryReq['entities'].append(idObj)

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.post(
            self.fogflowURL + '/ngsi10/queryContext', data=json.dumps(queryReq), headers=headers)

        entityList = []

        if response.status_code != 200:
            print('failed to query context')
            print(response.text)
        else:
            jsonData = response.json()
            for element in jsonData['contextResponses']:
                entity = ContextEntity()
                entity.fromContextElement(element['contextElement'])
                entityList.append(entity)

        return entityList

    def getByType(self, entityType):
        queryReq = {}
        queryReq['entities'] = []

        idObj = {}
        idObj['type'] = entityType
        idObj['isPattern'] = True
        queryReq['entities'].append(idObj)

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.post(
            self.fogflowURL + '/ngsi10/queryContext', data=json.dumps(queryReq), headers=headers)

        entityList = []

        if response.status_code != 200:
            print('failed to query context')
            print(response.text)
        else:
            jsonData = response.json()
            for element in jsonData['contextResponses']:
                entity = ContextEntity()
                entity.fromContextElement(element['contextElement'])
                entityList.append(entity)

        return entityList

    def delete(self, entityId):
        contextElement = {}

        idObj = {}
        idObj['id'] = entityId
        idObj['isPattern'] = False
        contextElement['entityId'] = idObj

        updateCtxReq = {}
        updateCtxReq['contextElements'] = []
        updateCtxReq['contextElements'].append(contextElement)
        updateCtxReq['updateAction'] = 'DELETE'

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.post(self.fogflowURL + '/ngsi10/updateContext',
                                 data=json.dumps(updateCtxReq), headers=headers)
        if response.status_code != 200:
            print('failed to delete context entity ' + entityId)
            print(response.text)
            return False
        else:
            print('delete context entity ' + entityId)
            return True

    def subscribe(self, entityType, refURL):
        subscribeCtxReq = {}
        subscribeCtxReq['entities'] = []

        idObj = {}
        idObj['type'] = entityType
        idObj['isPattern'] = True
        subscribeCtxReq['entities'].append(idObj)

        subscribeCtxReq['reference'] = refURL

        headers = {'Accept': 'application/json',
                   'Content-Type': 'application/json'}
        response = requests.post(self.fogflowURL + '/ngsi10/subscribeContext',
                                 data=json.dumps(subscribeCtxReq), headers=headers)

        if response.status_code != 200:
            print('failed to subscribe')
            print(response.text)
            return ''
        else:
            return response.json()['subscribeResponse']['subscriptionId']
