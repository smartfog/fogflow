package main

import (
	"encoding/json"
	"bytes"
	"net/http"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	. "fogflow/common/communicator"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"

	. "fogflow/common/config"
)

type Master struct {
	cfg *Config

	BrokerURL string

	id           string
	myURL        string
	messageBus   string
	discoveryURL string
	designerURL  string
        SecurityCfg     *HTTPS

	communicator *Communicator
	communicator2 *Communicator
	ticker       *time.Ticker
	agent        *NGSIAgent

	//list of all workers
	workers         map[string]*WorkerProfile
	workerList_lock sync.RWMutex

	//list of all operators
	operatorList      map[string]Operator
	operatorList_lock sync.RWMutex

	//list of all docker images
	dockerImageList      map[string][]DockerImage
	dockerImageList_lock sync.RWMutex

	//list of all submitted topologies
	topologyList      map[string]*Topology
	topologyList_lock sync.RWMutex

	//list of all submitted topologies
	fogfunctionList      map[string]*FogFunction
	fogfunctionList_lock sync.RWMutex

	//to manage the orchestration of service topology
	serviceMgr *ServiceMgr

	//to manage the orchestration of tasks
	taskMgr *TaskMgr

	//number of deployed task
	curNumOfTasks int
	prevNumOfTask int
	counter_lock  sync.RWMutex

	//type of subscribed entities
	subID2Type map[string]string
}

func (master *Master) Start(configuration *Config) {
	master.cfg = configuration
	master.SecurityCfg = &configuration.HTTPS
	
	master.messageBus = configuration.GetMessageBus()
	master.discoveryURL = configuration.GetDiscoveryURL()
	master.designerURL = configuration.GetDesignerURL()
	
	master.workers = make(map[string]*WorkerProfile)

	master.operatorList = make(map[string]Operator)
	master.dockerImageList = make(map[string][]DockerImage)
	master.topologyList = make(map[string]*Topology)
	master.fogfunctionList = make(map[string]*FogFunction)

	master.subID2Type = make(map[string]string)

	// communicate with the cloud_broker
	master.BrokerURL = configuration.GetBrokerURL()
	INFO.Println("communicate with the cloud broker via ", master.BrokerURL)

	// initialize the manager for both fog function and service topology
	master.taskMgr = NewTaskMgr(master)
	master.taskMgr.Init()

	master.serviceMgr = NewServiceMgr(master)
	master.serviceMgr.Init()

	// announce myself to the nearby IoT Broker
	for {
		// announce myself to the nearby IoT Broker
		err := master.registerMyself()
		if err != nil {
			INFO.Println("wait for the assigned broker to be ready")
			time.Sleep(5 * time.Second)
		} else {
			INFO.Println("annouce myself to the nearby broker")
			break
		}
	}

	master.myURL = "http://" + configuration.GetMasterIP() + ":" + strconv.Itoa(configuration.Master.AgentPort)

	// start the NGSI agent
	master.agent = &NGSIAgent{Port: configuration.Master.AgentPort, SecurityCfg: master.cfg.HTTPS}
	master.agent.Start()
	//master.agent.SetContextNotifyHandler(master.onReceiveContextNotify)
	master.agent.SetContextAvailabilityNotifyHandler(master.onReceiveContextAvailability)

	 go func() {
	      
              body, err := json.Marshal(map[string]string{
                       "status" : "Master is Up"})
                if err != nil {
                        fmt.Println(err)
                }
               //fmt.Println(master.cfg.HTTPS)
               master.cfg.HTTPS.LoadConfig()
               client := master.cfg.HTTPS.GetHTTPClient()
               fmt.Println("==== client =====",client)
               req2, err := http.NewRequest("POST", master.designerURL+"/masterNotify", bytes.NewBuffer(body))
               fmt.Println("++++++ req ++++++ and err",req2,err)
               for {
                        //resp, err := client.Post(url,"application/json" , bytes.NewBuffer(body))
                       time.Sleep(5 * time.Second)
                       resp, err := client.Do(req2)
                       fmt.Println(err)
                        if(resp != nil) {
                                defer resp.Body.Close()
                                break
                        }
                        
                }
        }()
	
	// start the message consumer
	go func() {
		cfg := MessageBusConfig{}
		cfg.Broker = configuration.GetMessageBus()
		cfg.Exchange = "fogflow"
		cfg.ExchangeType = "topic"
		cfg.DefaultQueue = master.id
		cfg.BindingKeys = []string{master.id + ".", "heartbeat.*"}

		// create the communicator with the broker info and topics
		master.communicator = NewCommunicator(&cfg)
		for {
			retry, err := master.communicator.StartConsuming(master.id, master)
			if retry {
				INFO.Printf("Going to retry launching the rabbitmq. Error: %v", err)
			} else {
				INFO.Printf("stop retrying")
				break
			}
		}
	}()

	go func() {
                cfg1 := MessageBusConfig{}
                cfg1.Broker = configuration.GetMessageBus()
                cfg1.Exchange = "Op"
                cfg1.ExchangeType = "topic"
                cfg1.DefaultQueue = "Operator"
                cfg1.BindingKeys = []string{"Operator.", "heartbeat.*"}

                // create the communicator with the broker info and topics
                master.communicator2 = NewCommunicator(&cfg1)
                for {
                        retry, err := master.communicator2.StartConsuming("Operator", master)
                        if retry {
                                INFO.Printf("Going to retry launching the rabbitmq. Error: %v", err)
                        } else {
                                INFO.Printf("stop retrying")
                                break
                        }
                }
        }()


	master.prevNumOfTask = 0
	master.curNumOfTasks = 0

	// start a timer to do something periodically
	master.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			<-master.ticker.C
			//master.onTimer()
		}
	}()

	// subscribe to the update of required context information
	//master.triggerInitialSubscriptions()
}

func (master *Master) onTimer() {
	master.counter_lock.Lock()
	delta := master.curNumOfTasks - master.prevNumOfTask
	fmt.Printf("# of orchestrated tasks = %d, throughput = %d/s\r\n", master.curNumOfTasks, delta)
	master.prevNumOfTask = master.curNumOfTasks
	master.counter_lock.Unlock()
}

func (master *Master) Quit() {
	INFO.Println("to stop the master")
	master.unregisterMyself()
	INFO.Println("unregister myself")
	master.ticker.Stop()
	INFO.Println("stop the timer")
	master.communicator.StopConsuming()
	INFO.Println("stop consuming the message")
}

func (master *Master) registerMyself() error {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = master.id
	ctxObj.Entity.Type = "Master"
	ctxObj.Entity.IsPattern = false

	ctxObj.Metadata = make(map[string]ValueObject)

	mylocation := Point{}
	mylocation.Latitude = master.cfg.Location.Latitude
	mylocation.Longitude = master.cfg.Location.Longitude
	ctxObj.Metadata["location"] = ValueObject{Type: "point", Value: mylocation}

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	return err
}

func (master *Master) unregisterMyself() {
	entity := EntityId{}
	entity.ID = master.id
	entity.Type = "Master"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	err := client.DeleteContext(&entity)
	if err != nil {
		ERROR.Println(err)
	}
}

/*func (master *Master) triggerInitialSubscriptions() {
	master.subscribeContextEntity("Operator")
	master.subscribeContextEntity("DockerImage")
	master.subscribeContextEntity("Topology")
	master.subscribeContextEntity("FogFunction")
	master.subscribeContextEntity("ServiceIntent")
	master.subscribeContextEntity("TaskIntent")
}*/

func (master *Master) subscribeContextEntity(entityType string) {
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}
	newEntity.Type = entityType
	newEntity.IsPattern = true
	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)
	subscription.Reference = master.myURL

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		ERROR.Println(err)
	}
	INFO.Println(sid)

	master.subID2Type[sid] = entityType
}

/*func (master *Master) onReceiveContextNotify(notifyCtxReq *NotifyContextRequest) {
	sid := notifyCtxReq.SubscriptionId
	stype := master.subID2Type[sid]

	if len(notifyCtxReq.ContextResponses) == 0 {
		return
	}
	fmt.Println(" %%%%%%%%%%%%%%%%% notify %%%%%%%%%%%%%%",notifyCtxReq.ContextResponses[0].ContextElement)
	contextObj := CtxElement2Object(&notifyCtxReq.ContextResponses[0].ContextElement)
	switch stype {
	// registry of an operator
	case "Operator":
		fmt.Println(" ***** Operastor registry ********",contextObj)
		//master.handleOperatorRegistration(contextObj)

	// registry of a docker image
	case "DockerImage":
		//master.handleDockerImageRegistration(contextObj)

	// topology to define service template
	case "Topology":
		//master.handleTopologyUpdate(contextObj)

	// fog function that includes a pair of topology and intent
	case "FogFunction":
		//master.handleFogFunctionUpdate(contextObj)

	// service orchestration
	case "ServiceIntent":
		//master.serviceMgr.handleServiceIntentUpdate(contextObj)

	// task orchestration
	case "TaskIntent":
		//master.taskMgr.handleTaskIntentUpdate(contextObj)
	}
}*/

//
// to handle the registry of operator
//

/*func (master *Master) handleOperatorRegistration(operatorCtxObj *ContextObject) {
	INFO.Println(operatorCtxObj)

	if operatorCtxObj.IsEmpty() {
		// does not handle the removal of operator
		return
	}

	var operator = Operator{}
	jsonText, _ := json.Marshal(operatorCtxObj.Attributes["operator"].Value.(map[string]interface{}))
	fmt.Println("**** jsonText *****",jsonText)
	err := json.Unmarshal(jsonText, &operator)
	fmt.Println("**** Unmarshal Operator inside function *****",operator)
	if err != nil {
		ERROR.Println("failed to read the given operator")
	} else {
		master.operatorList_lock.Lock()
		master.operatorList[operator.Name] = operator
		master.operatorList_lock.Unlock()
	}
}*/

//start: Pradumn
func (master *Master) handleOperatorRegistration(msg json.RawMessage) {
	INFO.Println(string(msg))
	//fmt.Println(len(msg))
	var operator = Operator{}
        err := json.Unmarshal(msg, &operator)
	if(len(msg) <= 2 || len(operator.Name) == 0) {
		//does not handle the removal of operator
                return
        }
	INFO.Println("Operator : ",&operator)

	if err!=nil {
		ERROR.Println("failed to read the given operator")
	}else {
		master.operatorList_lock.Lock()
                master.operatorList[operator.Name] = operator
                master.operatorList_lock.Unlock()
	}
}
//End: Pradumn


//
// to handle the management of docker images
//
/*func (master *Master) handleDockerImageRegistration(dockerImageCtxObj *ContextObject) {
	INFO.Println(dockerImageCtxObj)

	if dockerImageCtxObj.IsEmpty() {
		// does not handle the removal of operator
		return
	}

	dockerImage := DockerImage{}
	dockerImage.OperatorName = dockerImageCtxObj.Attributes["operator"].Value.(string)
	dockerImage.ImageName = dockerImageCtxObj.Attributes["image"].Value.(string)
	dockerImage.ImageTag = dockerImageCtxObj.Attributes["tag"].Value.(string)
	dockerImage.TargetedHWType = dockerImageCtxObj.Attributes["hwType"].Value.(string)
	dockerImage.TargetedOSType = dockerImageCtxObj.Attributes["osType"].Value.(string)
	dockerImage.Prefetched = dockerImageCtxObj.Attributes["prefetched"].Value.(bool)
	
	fmt.Println("********* dockerImage *********",dockerImage)
	
	master.dockerImageList_lock.Lock()
	master.dockerImageList[dockerImage.OperatorName] = append(master.dockerImageList[dockerImage.OperatorName], dockerImage)
	master.dockerImageList_lock.Unlock()

	if dockerImage.Prefetched == true {
		// inform all workers to prefetch this docker image in advance
		master.prefetchDockerImages(dockerImage)
	}
}*/

//Start: Pradumn
func (master *Master) handleDockerImageRegistration(msg json.RawMessage) {
	if(len(msg) <=2) {
		//does not handle the removal of dockerImage
		return
	}
	INFO.Println(string(msg))

	var dockerImage = DockerImage{}
        err := json.Unmarshal(msg, &dockerImage)
	fmt.Println("******* Docker Image ********",dockerImage.ImageName)
	fmt.Println("***** dockerImage operator name ****",dockerImage.OperatorName)
	if(len(dockerImage.OperatorName) == 0 || len(dockerImage.ImageName) == 0 || len(dockerImage.ImageTag) == 0 || len(dockerImage.TargetedHWType) == 0 || len(dockerImage.TargetedOSType) == 0) {
                //does not handle the removal of dockerImage
                return
        }
        INFO.Println("dockerImage : ",&dockerImage)

        if err!=nil {
                ERROR.Println("failed to read the given dockerImage")
	}else {
		master.dockerImageList_lock.Lock()
	        master.dockerImageList[dockerImage.OperatorName] = append(master.dockerImageList[dockerImage.OperatorName], dockerImage)
		master.dockerImageList_lock.Unlock()
	}
	if dockerImage.Prefetched == true {
                // inform all workers to prefetch this docker image in advance
                master.prefetchDockerImages(dockerImage)
        }
  }


func (master *Master) prefetchDockerImages(image DockerImage) {
	master.workerList_lock.RLock()
	defer master.workerList_lock.RUnlock()

	for _, worker := range master.workers {
		workerID := worker.WID
		taskMsg := SendMessage{Type: "PREFETCH_IMAGE", RoutingKey: workerID + ".", From: master.id, PayLoad: image}
		INFO.Println(taskMsg)
		master.communicator.Publish(&taskMsg)
	}
}

//
// to update the fog function list
//
/*func (master *Master) handleFogFunctionUpdate(fogfunctionCtxObj *ContextObject) {
	INFO.Println(fogfunctionCtxObj)

	// the fog function is going to be deleted
	if fogfunctionCtxObj.IsEmpty() {
		var eid = fogfunctionCtxObj.Entity.ID

		master.fogfunctionList_lock.RLock()
		fogfunction := master.fogfunctionList[eid]
		master.fogfunctionList_lock.RUnlock()

		DEBUG.Printf("%+v\r\n", fogfunction)

		// remove the service intent
		master.serviceMgr.removeServiceIntent(fogfunction.Intent.ID)

		// remove the service topology
		topology := fogfunction.Topology
		master.topologyList_lock.Lock()
		master.topologyList[topology.Name] = &topology
		master.topologyList_lock.Unlock()

		// remove this fog function entity
		master.fogfunctionList_lock.Lock()
		delete(master.fogfunctionList, eid)
		master.fogfunctionList_lock.Unlock()

		return
	}

	topology := Topology{}

	topologyJsonText, err := json.Marshal(fogfunctionCtxObj.Attributes["topology"].Value.(map[string]interface{}))
	fmt.Println("*********topologyJsonText*************",topologyJsonText)
	if err != nil {
		ERROR.Println("the topology object is not defined")
		return
	}
	err = json.Unmarshal(topologyJsonText, &topology)
	fmt.Println("*********topology*************",topology)
	if err != nil {
		ERROR.Println("the topology object is not correctly defined")
		return
	}

	intent := ServiceIntent{}

	intentJsonText, err := json.Marshal(fogfunctionCtxObj.Attributes["intent"].Value.(map[string]interface{}))
	if err != nil {
		ERROR.Println("the intent object is not defined")
		return
	}
	err = json.Unmarshal(intentJsonText, &intent)
	fmt.Println("*********intent*************",intent)
	if err != nil {
		ERROR.Println("the intent object is not correctly defined")
		return
	}

	// allow the ID of this service intent
	intent.ID = fogfunctionCtxObj.Entity.ID

	fogfunction := FogFunction{}

	fogfunction.Id = fogfunctionCtxObj.Entity.ID
	fogfunction.Name = fogfunctionCtxObj.Attributes["name"].Value.(string)
	fogfunction.Topology = topology
	fogfunction.Intent = intent

	// add the service topology
	master.topologyList_lock.Lock()
	master.topologyList[topology.Name] = &topology
	master.topologyList_lock.Unlock()

	// handle the associated service intent
	master.serviceMgr.handleServiceIntent(&fogfunction.Intent)

	fmt.Println("******* inside  handle fog function ***********",fogfunction)
	// create or update this fog function
	master.fogfunctionList_lock.Lock()
	master.fogfunctionList[fogfunction.Id] = &fogfunction
	master.fogfunctionList_lock.Unlock()

	INFO.Println(fogfunction)
}*/


func (master *Master) handleFogFunctionUpdate(msg json.RawMessage) {
	//INFO.Println(msg)
	var fogfunction = FogFunction{}
        err := json.Unmarshal(msg, &fogfunction)
	fogfunction.Intent.ID = fogfunction.Id
	fmt.Println("***** Intent.ID *********",fogfunction.Intent.ID)
	//fmt.Println("********* msg *******",msg, fogfunction)

	if(fogfunction.Action == "DELETE"){
		var eid = fogfunction.Id
		
		master.fogfunctionList_lock.RLock()
		fogfunction := master.fogfunctionList[eid]
		master.fogfunctionList_lock.RUnlock()

		DEBUG.Printf("%+v\r\n", fogfunction)
		master.serviceMgr.removeServiceIntent(fogfunction.Intent.ID)
		
		topology := fogfunction.Topology
		master.topologyList_lock.Lock()
		master.topologyList[topology.Name] = &topology
		master.topologyList_lock.Unlock()

		// remove this fog function entity
		master.fogfunctionList_lock.Lock()
		delete(master.fogfunctionList, eid)
		master.fogfunctionList_lock.Unlock()

		return

	}

        fmt.Println("&&&&&&&&& topology and name &&&&&&&&",&fogfunction.Topology, &fogfunction.Topology.Name)
        master.topologyList_lock.Lock()
        master.topologyList[fogfunction.Topology.Name] = &fogfunction.Topology
        master.topologyList_lock.Unlock()
        master.serviceMgr.handleServiceIntent(&fogfunction.Intent)
        fmt.Println("&&&&&&&&& topology from fogfunction and error &&&&&&&&",&fogfunction,err)

	// create or update this fog function
	master.fogfunctionList_lock.Lock()
	master.fogfunctionList[fogfunction.Id] = &fogfunction
	master.fogfunctionList_lock.Unlock()

	INFO.Println(fogfunction)

}


//
// to update the topology list
//
/*func (master *Master) handleTopologyUpdate(topologyCtxObj *ContextObject) {
	INFO.Println(topologyCtxObj)

	if topologyCtxObj.IsEmpty() {
		// remove this service topology entity
		master.topologyList_lock.Lock()

		var eid = topologyCtxObj.Entity.ID

		// find which one has this id
		for _, topology := range master.topologyList {
			if topology.Id == eid {
				var name = topology.Name
				delete(master.topologyList, name)
				break
			}
		}

		master.topologyList_lock.Unlock()

		return
	}

	// create or update this service topology
	topology := Topology{}
	jsonText, _ := json.Marshal(topologyCtxObj.Attributes["template"].Value.(map[string]interface{}))
	err := json.Unmarshal(jsonText, &topology)
	if err == nil {
		INFO.Println(topology)

		topology.Id = topologyCtxObj.Entity.ID
		fmt.Println("****** topology.Id *****",topology.Id)
		fmt.Println("****** topogoly ko handle karna hai *********",&topology)

		master.topologyList_lock.Lock()
		master.topologyList[topology.Name] = &topology
		master.topologyList_lock.Unlock()

		INFO.Println(topology)
	}

}*/

//start:Pradumn
func (master *Master) handleTopologyUpdate(msg json.RawMessage) {
	INFO.Println(string(msg))

	topology := Topology{}
        err := json.Unmarshal(msg, &topology)
	fmt.Println("***** len(topology.Tasks)*****",len(topology.Tasks))
	fmt.Println("******* unmarshalled topology********",&topology)

	if(topology.Action == "DELETE") {
                 master.topologyList_lock.Lock()

                var eid = "Topology." + topology.Id

                // find which one has this id
                for _, topologyToCheck := range master.topologyList {
                        if topologyToCheck.Id == eid {
                                var name = topologyToCheck.Name
                                delete(master.topologyList, name)
				INFO.Println(name," this topology is deleted ~~~~~~~~~~",master.topologyList)
                                break
                        }
                }

                master.topologyList_lock.Unlock()

                return
        }


	if err == nil {
                INFO.Println(topology)

                topology.Id = "Topology." + topology.Name
                fmt.Println("****** topology.Id *****",topology.Id)
                fmt.Println("****** topogoly ko handle karna hai *********",&topology)

                master.topologyList_lock.Lock()
		fmt.Println("******** topology list ****",master.topologyList)
                master.topologyList[topology.Name] = &topology
                master.topologyList_lock.Unlock()

                INFO.Println(topology)
        }


}

func (master *Master) getTopologyByName(name string) *Topology {
	// find the required topology object
	master.topologyList_lock.RLock()
	defer master.topologyList_lock.RUnlock()

	topology := master.topologyList[name]
	return topology
}

func (master *Master) queryWorkers() []*ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.Type = "Worker"
	entity.IsPattern = true
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	ctxObjects, err := client.QueryContext(&query)
	if err != nil {
		ERROR.Println(err)
	}

	return ctxObjects
}

func (master *Master) onReceiveContextAvailability(notifyCtxAvailReq *NotifyContextAvailabilityRequest) {
	INFO.Println("===========RECEIVE CONTEXT AVAILABILITY=========")
	DEBUG.Println(notifyCtxAvailReq)
	subID := notifyCtxAvailReq.SubscriptionId

	var action string
	switch notifyCtxAvailReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	}

	for _, registrationResp := range notifyCtxAvailReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration
		//entityRegistration := EntityRegistration{}
		for _, entity := range registration.EntityIdList {
			// convert context registration to entity registration
			fmt.Println("entity.MsgFormat", entity.MsgFormat)
			if entity.MsgFormat == "NGSILD" {
				entityRegistration := master.ldContextRegistration2EntityRegistration(&entity, &registration)
				go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
			} else {
				entityRegistration := master.contextRegistration2EntityRegistration(&entity, &registration)
				go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
			}
			//go master.taskMgr.HandleContextAvailabilityUpdate(subID, action, entityRegistration)
		}
	}
}

func (master *Master) RetrieveContextLdEntity(eid string, fsp string) interface{} {
	query := LDQueryContextRequest{}

	query.Entities = make([]EntityId, 0)
	query.Type = "Query"
	entity := EntityId{}
	idSplit := strings.Split(eid, "@")
	entity.ID = idSplit[0]
	entity.IsPattern = false
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	ctxObjects, err := client.QueryLdContext(&query, idSplit[1], fsp)
	if err == nil && ctxObjects != nil && len(ctxObjects) > 0 {
		return ctxObjects[0]
	} else {
		if err != nil {
			ERROR.Println("error occured when retrieving a context entity :", err)
		}

		return nil
	}
}

func (master *Master) ldContextRegistration2EntityRegistration(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	ctxObj := master.RetrieveContextLdEntity(entityId.ID, entityId.FiwareServicePath)
	if ctxObj == nil {
		entityRegistration.ID = entityId.ID
		entityRegistration.Type = entityId.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
	} else {
		ldCtcObj := ctxObj.(map[string]interface{})
		entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)
		entityRegistration.MetadataList = make(map[string]ContextMetadata)
		entityRegistration.MsgFormat = entityId.MsgFormat
		for key, attr := range ldCtcObj {
			if key != "modifiedAt" && key != "createdAt" && key != "observationSpace" && key != "operationSpace" && key != "@context" && key != "fiwareServicePath" {
				if key == "id" {
					entityRegistration.ID = entityId.ID
				} else if key == "type" {
					entityRegistration.Type = ldCtcObj[key].(string)
				} else if key == "FiwareServicePath" {
					entityRegistration.FiwareServicePath = ldCtcObj[key].(string)
				} else {
					attrmap := attr.(map[string]interface{})
					if attrmap["type"] != "GeoProperty" {
						attributeRegistration := ContextRegistrationAttribute{}
						attributeRegistration.Name = key
						attributeRegistration.Type = attrmap["type"].(string)
						entityRegistration.AttributesList[key] = attributeRegistration
					} else {
						metaData := attr.(map[string]interface{})
						cm := ContextMetadata{}
						cm.Name = key
						matadataCordinate := metaData["value"].(map[string]interface{})
						typ, points := GetNGSIV1DomainMetaData(matadataCordinate["type"].(string), matadataCordinate["coordinates"])
						cm.Type = typ
						cm.Value = points
						entityRegistration.MetadataList[key] = cm
					}
				}
			}
		}
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) contextRegistration2EntityRegistration(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	ctxObj := master.RetrieveContextEntity(entityId.ID)
	if ctxObj == nil {
		entityRegistration.ID = entityId.ID
		entityRegistration.Type = entityId.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
	} else {
		entityRegistration.ID = ctxObj.Entity.ID
		entityRegistration.Type = ctxObj.Entity.Type
		entityRegistration.FiwareServicePath = entityId.FiwareServicePath
		entityRegistration.MsgFormat = entityId.MsgFormat
		entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)
		for attrName, attrValue := range ctxObj.Attributes {
			attributeRegistration := ContextRegistrationAttribute{}
			attributeRegistration.Name = attrName
			attributeRegistration.Type = attrValue.Type
			entityRegistration.AttributesList[attrName] = attributeRegistration
		}

		entityRegistration.MetadataList = make(map[string]ContextMetadata)
		for metaname, ctxmeta := range ctxObj.Metadata {
			cm := ContextMetadata{}
			cm.Name = metaname
			cm.Type = ctxmeta.Type
			cm.Value = ctxmeta.Value

			entityRegistration.MetadataList[metaname] = cm
		}
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) contextRegistration2EntityRegistrationNew(entityId *EntityId, ctxRegistration *ContextRegistration) *EntityRegistration {
	entityRegistration := EntityRegistration{}

	entityRegistration.ID = entityId.ID
	entityRegistration.Type = entityId.Type

	entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)

	for _, attribute := range ctxRegistration.ContextRegistrationAttributes {
		attributeRegistration := ContextRegistrationAttribute{}
		attributeRegistration.Name = attribute.Name
		attributeRegistration.Type = attribute.Type

		entityRegistration.AttributesList[attribute.Name] = attributeRegistration
	}

	entityRegistration.MetadataList = make(map[string]ContextMetadata)
	for _, ctxmeta := range ctxRegistration.Metadata {
		cm := ContextMetadata{}
		cm.Name = ctxmeta.Name
		cm.Type = ctxmeta.Type
		cm.Value = ctxmeta.Value

		entityRegistration.MetadataList[ctxmeta.Name] = cm
	}

	entityRegistration.ProvidingApplication = ctxRegistration.ProvidingApplication

	DEBUG.Printf("REGISTERATION OF ENTITY CONTEXT AVAILABILITY: %+v\r\n", entityRegistration)

	return &entityRegistration
}

func (master *Master) subscribeContextAvailability(availabilitySubscription *SubscribeContextAvailabilityRequest) string {

	availabilitySubscription.Reference = master.myURL + "/notifyContextAvailability"

	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL(), SecurityCfg: &master.cfg.HTTPS}
	subscriptionId, err := client.SubscribeContextAvailability(availabilitySubscription)
	if err != nil {
		ERROR.Println(err)
		return ""
	}

	return subscriptionId
}

func (master *Master) unsubscribeContextAvailability(sid string) {
	client := NGSI9Client{IoTDiscoveryURL: master.cfg.GetDiscoveryURL(), SecurityCfg: &master.cfg.HTTPS}
	err := client.UnsubscribeContextAvailability(sid)
	if err != nil {
		ERROR.Println(err)
	}
}

//
// to deal with the communication between master and workers via rabbitmq
//
func (master *Master) Process(msg *RecvMessage) error {
	switch msg.Type {
	case "heart_beat":
		profile := WorkerProfile{}
		err := json.Unmarshal(msg.PayLoad, &profile)
		if err == nil {
			master.onHeartbeat(msg.From, &profile)
		}

	case "task_update":
		update := TaskUpdate{}
		err := json.Unmarshal(msg.PayLoad, &update)
		if err == nil {
			master.onTaskUpdate(msg.From, &update)
		}

	case "Operator":
		master.handleOperatorRegistration(msg.PayLoad)

	case "DockerImage":
		master.handleDockerImageRegistration(msg.PayLoad)

	case "FogFunction":
		master.handleFogFunctionUpdate(msg.PayLoad)


	case "Topology":
                master.handleTopologyUpdate(msg.PayLoad)

	case "ServiceIntent": 
		master.serviceMgr.handleServiceIntentUpdate(msg.PayLoad)

	}

	return nil
}

func (master *Master) onHeartbeat(from string, profile *WorkerProfile) {
	master.workerList_lock.Lock()

	workerID := profile.WID
	fmt.Println("**** workerID and profile ******",workerID,profile)
	if worker, exist := master.workers[workerID]; exist {
		worker.Capacity = profile.Capacity
	} else {
		profile.Workload = 0
		master.workers[workerID] = profile
	}

	master.workerList_lock.Unlock()
}

func (master *Master) onTaskUpdate(from string, update *TaskUpdate) {
	INFO.Println("==task update=========")
	INFO.Println(update)

}

func (master *Master) DeployTask(taskInstance *ScheduledTaskInstance) {
	master.counter_lock.Lock()
	master.curNumOfTasks = master.curNumOfTasks + 1
	master.counter_lock.Unlock()

	taskMsg := SendMessage{Type: "ADD_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
	INFO.Println(taskMsg)

	go master.communicator.Publish(&taskMsg)

	// update the workload of this worker
	workerID := taskInstance.WorkerID

	master.workerList_lock.Lock()
	workerProfile := master.workers[workerID]
	workerProfile.Workload = workerProfile.Workload + 1
	master.workerList_lock.Unlock()
}

func (master *Master) TerminateTask(taskInstance *ScheduledTaskInstance) {
	taskMsg := SendMessage{Type: "REMOVE_TASK", RoutingKey: taskInstance.WorkerID + ".", From: master.id, PayLoad: *taskInstance}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)

	// update the workload of this worker
	workerID := taskInstance.WorkerID

	master.workerList_lock.Lock()
	workerProfile := master.workers[workerID]
	workerProfile.Workload = workerProfile.Workload - 1
	master.workerList_lock.Unlock()
}

func (master *Master) AddInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "ADD_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.id, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

func (master *Master) RemoveInputEntity(flowInfo FlowInfo) {
	taskMsg := SendMessage{Type: "REMOVE_INPUT", RoutingKey: flowInfo.WorkerID + ".", From: master.id, PayLoad: flowInfo}
	INFO.Println(taskMsg)
	master.communicator.Publish(&taskMsg)
}

//
// the shared functions for function manager and topology manager to call
//

func (master *Master) RetrieveContextEntity(eid string) *ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.ID = eid
	entity.IsPattern = false
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: master.BrokerURL, SecurityCfg: &master.cfg.HTTPS}
	ctxObjects, err := client.QueryContext(&query)
	if err == nil && ctxObjects != nil && len(ctxObjects) > 0 {
		return ctxObjects[0]
	} else {
		if err != nil {
			ERROR.Println("error occured when retrieving a context entity :", err)
		}

		return nil
	}
}

//
// to select the right docker image of an operator for the selected worker
//
func (master *Master) DetermineDockerImage(operatorName string, wID string) string {
	INFO.Println("select a suitable image to execute on the selected worker")

	master.workerList_lock.RLock()
	wProfile := master.workers[wID]
	master.workerList_lock.RUnlock()

	if wProfile == nil {
		ERROR.Println("could not find this worker from the curent worker list: ", wID)
		return ""
	}

	//select a suitable image to execute on the selected worker
	selectedDockerImageName := ""

	master.dockerImageList_lock.RLock()
	for _, image := range master.dockerImageList[operatorName] {
		fmt.Println("*****image*******",image)
		DEBUG.Println(wProfile)

		hwType := "X86"
		osType := "Linux"

		if wProfile.HWType == "arm" {
			hwType = "ARM"
		}

		if wProfile.OSType == "linux" {
			osType = "Linux"
		}

		if image.TargetedOSType == osType && image.TargetedHWType == hwType {
			selectedDockerImageName = image.ImageName + ":" + image.ImageTag
			break
		}
	}

	master.dockerImageList_lock.RUnlock()

	DEBUG.Println(selectedDockerImageName)

	return selectedDockerImageName
}

func (master *Master) GetOperatorParamters(operatorName string) []Parameter {
	master.operatorList_lock.RLock()

	operator := master.operatorList[operatorName]
	parameters := make([]Parameter, len(operator.Parameters))
	copy(parameters, operator.Parameters)

	master.operatorList_lock.RUnlock()

	return parameters
}

//
// to select the worker that is closest to the given points
//
func (master *Master) SelectWorker(locations []Point) string {
	master.workerList_lock.RLock()
	defer master.workerList_lock.RUnlock()
	fmt.Println("&&&& len(locations) &&&&&&&&&",len(locations))
	if len(locations) == 0 {
		for _, worker := range master.workers {
			return worker.WID
		}
		return ""
	}

	DEBUG.Printf("points: %+v\r\n", locations)
	fmt.Println("&&&& master.workers &&&&&&",master.workers)

	// select the workers with the closest distance and also the worker is currently not overloaded
	closestWorkerID := ""
	closestTotalDistance := uint64(18446744073709551615)
	for _, worker := range master.workers {
		fmt.Println("***** master.worker *******",worker)
		INFO.Printf("check worker %+v\r\n", worker)

		// if this worker is already overloaded, check the next one
		if worker.IsOverloaded() == true {
			continue
		}

		wp := Point{}
		wp.Latitude = worker.PLocation.Latitude
		wp.Longitude = worker.PLocation.Longitude

		totalDistance := uint64(0)

		for _, location := range locations {
			if location.IsEmpty() == true {
				continue
			}

			distance := Distance(&wp, &location)
			totalDistance += distance
			INFO.Printf("distance = %d between %+v, %+v\r\n", distance, wp, location)
		}

		if totalDistance < closestTotalDistance {
			closestWorkerID = worker.WID
			closestTotalDistance = totalDistance
		}

		INFO.Println("closest worker ", closestWorkerID, " with the closest distance ", closestTotalDistance)
	}

	// select the one with lowest capacity if there are more than one with the closest distance

	return closestWorkerID
}
