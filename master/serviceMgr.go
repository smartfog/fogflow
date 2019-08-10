package main

import (
	"encoding/json"
	"sync"

	"github.com/satori/go.uuid"

	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type TaskIntentRecord struct {
	taskIntent TaskIntent
	site       SiteInfo
}

type ServiceMgr struct {
	master *Master

	// list of all service intents and also the mapping between service intent and task intents
	serviceIntentMap map[string]*ServiceIntent
	service2TaskMap  map[string][]*TaskIntentRecord
	intentList_lock  sync.RWMutex
}

func NewServiceMgr(myMaster *Master) *ServiceMgr {
	return &ServiceMgr{master: myMaster}
}

func (sMgr *ServiceMgr) Init() {
	sMgr.serviceIntentMap = make(map[string]*ServiceIntent)
	sMgr.service2TaskMap = make(map[string][]*TaskIntentRecord)
}

func (sMgr *ServiceMgr) handleServiceIntentUpdate(intentCtxObj *ContextObject) {
	INFO.Println("handle intent update")
	INFO.Println(intentCtxObj)

	if intentCtxObj.IsEmpty() == true {
		sMgr.removeServiceIntent(intentCtxObj.Entity.ID)
	} else {
		sIntent := ServiceIntent{}
		jsonText, _ := json.Marshal(intentCtxObj.Attributes["intent"].Value.(map[string]interface{}))
		err := json.Unmarshal(jsonText, &sIntent)
		if err == nil {
			sIntent.ID = intentCtxObj.Entity.ID
			INFO.Println(sIntent)
			sMgr.handleServiceIntent(&sIntent)
		} else {
			ERROR.Println(err)
			//sMgr.updateServiceIntentStatus(intentCtxObj.Entity.ID, "NOT_ACTIVATED", "intent object is not properly defined")
		}
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
	INFO.Println("receive a service intent")
	INFO.Println(serviceIntent)

	var topologyObject = sMgr.master.getTopologyByName(serviceIntent.TopologyName)
	if topologyObject == nil {
		ERROR.Println("failed to find the associated topology")
		//sMgr.updateServiceIntentStatus(serviceIntent.ID, "NOT_ACTIVATED", "failed to find the associated topology")
		return
	}

	serviceIntent.TopologyObject = topologyObject

	listTaskIntent := make([]*TaskIntentRecord, 0)

	for _, task := range serviceIntent.TopologyObject.Tasks {
		// to handle the task intent directly
		taskIntent := TaskIntent{}

		// new random uid
		u1, _ := uuid.NewV4()
		rid := u1.String()
		taskIntent.ID = rid

		taskIntent.GeoScope = serviceIntent.GeoScope
		taskIntent.Priority = serviceIntent.Priority
		taskIntent.QoS = serviceIntent.QoS
		taskIntent.ServiceName = serviceIntent.TopologyName
		taskIntent.TaskObject = task

		INFO.Printf("%+v\n", taskIntent)

		sMgr.master.taskMgr.handleTaskIntent(&taskIntent)

		record := TaskIntentRecord{}
		record.taskIntent = taskIntent
		record.site.IsLocalSite = true

		listTaskIntent = append(listTaskIntent, &record)
	}

	sMgr.intentList_lock.Lock()
	defer sMgr.intentList_lock.Unlock()

	// to record the task intents for this high level service intent
	sMgr.service2TaskMap[serviceIntent.ID] = listTaskIntent
	// record the service intent
	sMgr.serviceIntentMap[serviceIntent.ID] = serviceIntent

	//sMgr.updateServiceIntentStatus(serviceIntent.ID, "ACTIVATED", "scheduled")
}

//
// to divide the task intent for all sites in this geoscope
//
func (sMgr *ServiceMgr) intentPartition(taskIntent *TaskIntent) []*TaskIntentRecord {
	var geoscope = taskIntent.GeoScope

	listTaskIntent := make([]*TaskIntentRecord, 0)

	client := NGSI9Client{IoTDiscoveryURL: sMgr.master.discoveryURL, SecurityCfg: &sMgr.master.cfg.HTTPS}
	siteList, err := client.QuerySiteList(geoscope)
	if err != nil {
		ERROR.Println("error happens when querying the site list from IoT Discovery")
		ERROR.Println(err)
	} else {
		DEBUG.Printf("%+v\n", siteList)

		for _, site := range siteList {
			if site.IsLocalSite == true {
				DEBUG.Printf("%+v is a local site\n", site)

				intent := TaskIntent{}

				// new random uid
				u1, _ := uuid.NewV4()
				rid := u1.String()
				intent.ID = rid

				intent.GeoScope = geoscope
				intent.Priority = taskIntent.Priority
				intent.QoS = taskIntent.QoS
				intent.ServiceName = taskIntent.ServiceName
				intent.TaskObject = taskIntent.TaskObject

				// handle a sub-intent locally
				sMgr.master.taskMgr.handleTaskIntent(&intent)

				record := TaskIntentRecord{}
				record.taskIntent = *taskIntent
				record.site = site

				listTaskIntent = append(listTaskIntent, &record)
			} else {
				DEBUG.Printf("%+v is a remote site\n", site)

				// forward a sub-intent to the remote site
				intent := TaskIntent{}

				// new random uid
				u1, _ := uuid.NewV4()
				rid := u1.String()
				intent.ID = rid

				intent.GeoScope = geoscope
				intent.Priority = taskIntent.Priority
				intent.QoS = taskIntent.QoS
				intent.ServiceName = taskIntent.ServiceName
				intent.TaskObject = taskIntent.TaskObject

				sMgr.ForwardIntentToRemoteSite(&intent, site)

				record := TaskIntentRecord{}
				record.taskIntent = *taskIntent
				record.site = site

				listTaskIntent = append(listTaskIntent, &record)
			}
		}
	}

	return listTaskIntent
}

func (sMgr *ServiceMgr) ForwardIntentToRemoteSite(taskIntent *TaskIntent, site SiteInfo) {
	brokerURL := "http://" + site.ExternalAddress + "/proxy"

	ctxElem := ContextElement{}
	ctxElem.Entity.ID = "TaskIntent." + taskIntent.ID
	ctxElem.Entity.Type = "TaskIntent"

	ctxElem.Attributes = make([]ContextAttribute, 0)

	attribute := ContextAttribute{}
	attribute.Type = "object"
	attribute.Name = "intent"
	attribute.Value = taskIntent

	ctxElem.Attributes = append(ctxElem.Attributes, attribute)

	client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: &sMgr.master.cfg.HTTPS}
	client.UpdateContext(&ctxElem)
}

func (sMgr *ServiceMgr) RemoveIntentFromRemoteSite(taskIntentRecord *TaskIntentRecord) {
	brokerURL := "http://" + taskIntentRecord.site.ExternalAddress + "/proxy"

	ctxElem := ContextElement{}
	ctxElem.Entity.ID = "TaskIntent." + taskIntentRecord.taskIntent.ID

	client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: &sMgr.master.cfg.HTTPS}
	client.UpdateContext(&ctxElem)
}

func (sMgr *ServiceMgr) removeServiceIntent(id string) {
	INFO.Printf("the master is going to remove the requested service intent %s\n", id)

	sMgr.intentList_lock.Lock()
	defer sMgr.intentList_lock.Unlock()

	// remove all related task intents
	listTaskIntentRecord := sMgr.service2TaskMap[id]
	for _, taskIntentRecord := range listTaskIntentRecord {
		if taskIntentRecord.site.IsLocalSite == true {
			// remove this task intent for the local site
			sMgr.master.taskMgr.removeTaskIntent(&taskIntentRecord.taskIntent)
		} else {
			// issue a request to delete this task intent that has been handled by a remote site
			sMgr.RemoveIntentFromRemoteSite(taskIntentRecord)
		}
	}

	// remove this service intent
	delete(sMgr.service2TaskMap, id)
}
