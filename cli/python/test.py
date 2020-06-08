import json
import time
import os

from fogflowclient import FogFlowClient, ContextEntity


def onResult(ctxEntity):
    print(ctxEntity)


def main():
    ffclient = FogFlowClient("http://localhost")    
        
    #push entities to the FogFlow system    
    for i in range(5):
        #create the input data
        deviceID = "Device.Car." + str(i) 
        
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

    time.sleep(2)

    #trigger a service topology and then subscribe to the generated result
    sessionId = ffclient.start("test", onResult)    
    
    while True:
        try:
            time.sleep(3)
        except:
            print("exit")
            
            entities = ffclient.getById(tempSensor.id)            
            for entity in entities:
                print(entity.toJSON())
            
            entities = ffclient.getByType(tempSensor.type)            
            for entity in entities:
                print(entity.toJSON())
                                                        
            for i in range(5):
                deviceID = "Device.Car." + str(i)                             
                ffclient.delete(deviceID)            
                        
            #stop the service topology
            ffclient.stop(sessionId)            
                        
            os._exit(1)



if __name__ == "__main__":
    main()
