import json
from fogflowclient import FogFlowClient, ContextEntity


def onReceiveContextEntity(ctxEntity):
    print(ctxEntity)


def main():
    ffclient = FogFlowClient("http://localhost")

    #create the input data
    deviceID = "Device.Car.0" 
    
    tempSensor = ContextEntity()                
    tempSensor.id = deviceID
    tempSensor.type = "Car"              
    tempSensor.attributes["temperature"] =  {'type': 'integer', 'value': 30}
    tempSensor.metadata["location"] = {
        "type":"point",
        "value":{
            "latitude":35.97800618085566,
            "longitude":139.41650390625003
        }
    }
        
    ffclient.put(tempSensor)

    intentId = ffclient.sendIntent("myTest")
    
    ffclient.get(deviceID)

    ffclient.removeIntent(intentId)


if __name__ == "__main__":
    main()
