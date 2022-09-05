package main

import (
	"sync"

	"github.com/google/uuid"

	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type ServiceMgr struct {
	master *Master

	// list of all service intents and also the mapping between service intent and task intents
	serviceIntentMap map[string]*ServiceIntent
	service2TaskMap  map[string][]*TaskIntent
	intentList_lock  sync.RWMutex
}

func NewServiceMgr(myMaster *Master) *ServiceMgr {
	return &ServiceMgr{master: myMaster}
}

func (sMgr *ServiceMgr) Init() {
	sMgr.serviceIntentMap = make(map[string]*ServiceIntent)
	sMgr.service2TaskMap = make(map[string][]*TaskIntent)
}

func (sMgr *ServiceMgr) handleServiceIntentUpdate(sIntent *ServiceIntent) {
	if sIntent.Action == "DELETE" {
		sMgr.removeServiceIntent(sIntent.ID)
	} else {
		sMgr.handleServiceIntent(sIntent)
	}
}

func (sMgr *ServiceMgr) updateServiceIntentStatus(eid string, status string, reason string) {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = eid
	ctxObj.Entity.IsPattern = false

	ctxObj.Metadata = make(map[string]ValueObject)

	ctxObj.Metadata["status"] = ValueObject{Type: "string", Value: status}
	ctxObj.Metadata["reason"] = ValueObject{Type: "string", Value: reason}

	client := NGSI10Client{IoTBrokerURL: sMgr.master.BrokerURL, SecurityCfg: &sMgr.master.cfg.HTTPS}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		ERROR.Println(err)
	}
}

//
// to break down the service intent from the service level into the task level
//
func (sMgr *ServiceMgr) handleServiceIntent(serviceIntent *ServiceIntent) {
	INFO.Println("[Service Intent]: ", serviceIntent.TopologyName)

	sMgr.intentList_lock.Lock()
	_, existFlag := sMgr.serviceIntentMap[serviceIntent.ID]
	sMgr.intentList_lock.Unlock()

	if existFlag == true {
		sMgr.updateExistingServiceIntent(serviceIntent)
	} else {
		sMgr.createNewServiceIntent(serviceIntent)
	}
}

//
// to break down the service intent from the service level into the task level
//
func (sMgr *ServiceMgr) updateExistingServiceIntent(serviceIntent *ServiceIntent) {
	INFO.Println("updating an existing intent")

	var topologyObject = sMgr.master.getTopologyByName(serviceIntent.TopologyName)
	if topologyObject == nil {
		ERROR.Println("failed to find the associated topology")
		return
	}

	serviceIntent.TopologyObject = topologyObject

	listTaskIntent := make([]*TaskIntent, 0)

	for _, task := range serviceIntent.TopologyObject.Tasks {
		// to handle the task intent directly
		taskIntent := TaskIntent{}

		// new random uid
		u1, _ := uuid.NewUUID()
		rid := u1.String()
		taskIntent.ID = rid

		taskIntent.GeoScope = serviceIntent.GeoScope
		taskIntent.Priority = serviceIntent.Priority
		//taskIntent.SType = serviceIntent.SType
		taskIntent.QoS = serviceIntent.QoS
		taskIntent.TopologyName = serviceIntent.TopologyName
		taskIntent.TaskObject = task

		INFO.Printf("%+v\n", taskIntent)

		sMgr.master.taskMgr.handleTaskIntent(&taskIntent)

		listTaskIntent = append(listTaskIntent, &taskIntent)
	}

	sMgr.intentList_lock.Lock()
	defer sMgr.intentList_lock.Unlock()

	// to record the task intents for this high level service intent
	sMgr.service2TaskMap[serviceIntent.ID] = listTaskIntent
	// record the service intent
	sMgr.serviceIntentMap[serviceIntent.ID] = serviceIntent
}

//
// to break down the service intent from the service level into the task level
//
func (sMgr *ServiceMgr) createNewServiceIntent(serviceIntent *ServiceIntent) {
	var topologyObject = sMgr.master.getTopologyByName(serviceIntent.TopologyName)
	if topologyObject == nil {
		ERROR.Println("failed to find the associated topology")
		return
	}

	serviceIntent.TopologyObject = topologyObject

	listTaskIntent := make([]*TaskIntent, 0)

	for _, task := range serviceIntent.TopologyObject.Tasks {
		// to handle the task intent directly
		taskIntent := TaskIntent{}

		// new random uid
		u1, _ := uuid.NewUUID()
		rid := u1.String()
		taskIntent.ID = rid

		taskIntent.GeoScope = serviceIntent.GeoScope
		taskIntent.Priority = serviceIntent.Priority
		//taskIntent.SType = serviceIntent.SType
		taskIntent.QoS = serviceIntent.QoS
		taskIntent.TopologyName = serviceIntent.TopologyName
		taskIntent.TaskObject = task
		taskIntent.ServiceIntentID = serviceIntent.ID

		sMgr.master.taskMgr.handleTaskIntent(&taskIntent)

		listTaskIntent = append(listTaskIntent, &taskIntent)
	}

	sMgr.intentList_lock.Lock()
	defer sMgr.intentList_lock.Unlock()

	// to record the task intents for this high level service intent
	sMgr.service2TaskMap[serviceIntent.ID] = listTaskIntent
	// record the service intent
	sMgr.serviceIntentMap[serviceIntent.ID] = serviceIntent
}

func (sMgr *ServiceMgr) removeServiceIntent(id string) {
	INFO.Printf("the master is going to remove the requested service intent %s\n", id)

	sMgr.intentList_lock.Lock()
	defer sMgr.intentList_lock.Unlock()

	// remove all related task intents
	listTaskIntent := sMgr.service2TaskMap[id]
	for _, taskIntent := range listTaskIntent {
		sMgr.master.taskMgr.removeTaskIntent(taskIntent)
	}

	// remove this service intent
	delete(sMgr.service2TaskMap, id)
	delete(sMgr.serviceIntentMap, id)
}
