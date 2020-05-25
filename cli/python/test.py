import json
from fogflowclient import FogFlowClient, ContextEntity

def main():
    ffclient = FogFlowClient("http://localhost")

    #create the input data
    tempSensor = ContextEntity()                
    tempSensor.id = "Device.Car." + str(x)
    tempSensor.type = "Car"      
        
    tempSensor.attributes["temperature"] =  {'type': 'integer', 'value': 30}
        
    tempSensor.metadata["location"] = {
        "type":"point",
        "value":{
            "latitude":35.97800618085566,
            "longitude":139.41650390625003
        }
    }
        
    deviceID = client.put(tempSensor)

    intentId = ffclient.sendIntent("mytest")

    result = ffclient.get(deviceID)
    
    ffclient.get(deviceID)


    ffclient.removeIntent(intentId)


if __name__ == "__main__":
    main()
