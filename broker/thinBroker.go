package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	. "fogflow/common/config"
	. "fogflow/common/constants"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/google/uuid"
	"github.com/piprate/json-gold/ld"
)

type ThinBroker struct {
	id              string
	MyLocation      PhysicalLocation
	MyURL           string
	IoTDiscoveryURL string
	SecurityCfg     *HTTPS

	myEntityId string

	myProfile BrokerProfile

	//mapping from subscriptionID to subscription
	subscriptions        map[string]*SubscribeContextRequest
	tmpNGSI10NotifyCache []string
	subscriptions_lock   sync.RWMutex

	tmpNGSIv2NotifyCache []string
	v2subscriptions      map[string]*SubscriptionRequest
	tmpNGSIV2NotifyCache map[string]*Notifyv2ContextAvailabilityRequest
	v2subscriptions_lock sync.RWMutex

	//mapping from main subscription to other related subscriptions
	main2Other              map[string][]string
	availabilitySub2MainSub map[string]string
	tmpNGSI9NotifyCache     map[string]*NotifyContextAvailabilityRequest
	subLinks_lock           sync.RWMutex

	//list of all updated context entities
	entities      map[string]*ContextElement //latest view of context entities
	entities_lock sync.RWMutex

	//Southbound feature addition
	fiwareData      map[string]*FiwareData
	fiwareData_lock sync.RWMutex
	//mapping from entityID to subscriptionID
	entityId2Subcriptions  map[string][]string
	e2sub_lock             sync.RWMutex
	ev2sub_lock            sync.RWMutex
	LDe2sub_lock           sync.RWMutex
	entityIdv2Subcriptions map[string][]string
	e2subv2_lock           sync.RWMutex

	//counter of heartbeat
	counter int64

	//NGSI-LD feature addition
	ldEntities      map[string]interface{} // to map Entity Id with LDContextElement.
	ldEntities_lock sync.RWMutex

	ldContextRegistrations      map[string]CSourceRegistrationRequest // to map Registration Id with CSourceRegistrationRequest.
	ldContextRegistrations_lock sync.RWMutex

	ldEntityID2RegistrationID      map[string]string //to map the Entity IDs with their registration id.
	ldEntityID2RegistrationID_lock sync.RWMutex

	ldSubscriptions      map[string]*LDSubscriptionRequest // to map Subscription Id with LDSubscriptionRequest.
	ldSubscriptions_lock sync.RWMutex

	tmpNGSIldNotifyCache    []string
	tmpNGSILDNotifyCache    map[string]*NotifyContextAvailabilityRequest
	entityId2LDSubcriptions map[string][]string
	entityTypeTOEntityId    map[string][]string
}

func (tb *ThinBroker) Start(cfg *Config) {
	tb.MyURL = cfg.GetBrokerURL()
	tb.IoTDiscoveryURL = cfg.GetDiscoveryURL()

	tb.myEntityId = tb.id

	tb.SecurityCfg = &cfg.HTTPS

	tb.MyLocation = cfg.Location

	tb.subscriptions = make(map[string]*SubscribeContextRequest)
	tb.tmpNGSI10NotifyCache = make([]string, 0)

	tb.tmpNGSIV2NotifyCache = make(map[string]*Notifyv2ContextAvailabilityRequest)

	tb.v2subscriptions = make(map[string]*SubscriptionRequest)
	tb.tmpNGSIv2NotifyCache = make([]string, 0)

	tb.entities = make(map[string]*ContextElement)
	tb.entityId2Subcriptions = make(map[string][]string)
	//Southbound feature addition
	tb.fiwareData = make(map[string]*FiwareData)

	tb.entityIdv2Subcriptions = make(map[string][]string)

	tb.availabilitySub2MainSub = make(map[string]string)
	tb.tmpNGSI9NotifyCache = make(map[string]*NotifyContextAvailabilityRequest)
	tb.main2Other = make(map[string][]string)

	tb.myProfile.BID = tb.myEntityId
	tb.myProfile.MyURL = tb.MyURL

	// NGSI-LD feature addition
	tb.ldEntities = make(map[string]interface{})
	tb.ldContextRegistrations = make(map[string]CSourceRegistrationRequest)
	tb.ldEntityID2RegistrationID = make(map[string]string)
	tb.ldSubscriptions = make(map[string]*LDSubscriptionRequest)
	tb.tmpNGSILDNotifyCache = make(map[string]*NotifyContextAvailabilityRequest)
	tb.entityId2LDSubcriptions = make(map[string][]string)
	tb.entityTypeTOEntityId = make(map[string][]string)
	// register itself to the IoT discovery
	tb.registerMyself()

	// send the first heartbeat message
	tb.sendHeartBeat()
}

func (tb *ThinBroker) Stop() {
	// deregister myself to IoT Discovery
	tb.deregisterMyself()

	// cancel all subscriptions that have been issues to outside
}

func (tb *ThinBroker) OnTimer() { // for every 2 second
	tb.subscriptions_lock.Lock()
	remainItems := tb.tmpNGSI10NotifyCache
	tb.tmpNGSI10NotifyCache = make([]string, 0)
	tb.subscriptions_lock.Unlock()
	for _, sid := range remainItems {
		hasCachedNotification := false
		tb.subscriptions_lock.Lock()
		if subscription, exist := tb.subscriptions[sid]; exist {
			if subscription.Subscriber.RequireReliability == true && len(subscription.Subscriber.NotifyCache) > 0 {
				hasCachedNotification = true
			}
		}
		tb.subscriptions_lock.Unlock()

		if hasCachedNotification == true {
			elements := make([]ContextElement, 0)
			tb.sendReliableNotify(elements, sid)
		}
	}

	// send heartbeat to IoT Discovery
	if tb.counter >= 5 {
		//every 10 seconds
		tb.sendHeartBeat()
		tb.counter = 0
	}
	tb.counter++

}

func (tb *ThinBroker) NGSILDOnTimer() {
	tb.ldSubscriptions_lock.Lock()
	remainNGSILDItems := tb.tmpNGSIldNotifyCache
	tb.tmpNGSIldNotifyCache = make([]string, 0)
	tb.ldSubscriptions_lock.Unlock()

	for _, sid := range remainNGSILDItems {
		hasCachedNGSILDNotification := false
		tb.ldSubscriptions_lock.Lock()
		if ldSubscription, exist := tb.ldSubscriptions[sid]; exist {
			if ldSubscription.Subscriber.RequireReliability == true && len(ldSubscription.Subscriber.LDNotifyCache) > 0 {
				hasCachedNGSILDNotification = true
			}
		}
		tb.ldSubscriptions_lock.Unlock()

		if hasCachedNGSILDNotification == true {
			elements := make([]map[string]interface{}, 0)
			tb.sendReliableNotifyToNgsiLDSubscriber(elements, sid)
		}
	}

}

func (tb *ThinBroker) sendHeartBeat() {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.SendHeartBeat(&tb.myProfile)
	if err != nil {
		ERROR.Println("failed to send my heartbeat info")
	}
}

func (tb *ThinBroker) registerMyself() bool {
	registerCtxReq := RegisterContextRequest{}
	registerCtxReq.ContextRegistrations = make([]ContextRegistration, 0)

	registration := ContextRegistration{}

	entities := make([]EntityId, 0)
	entity := EntityId{ID: tb.myEntityId, Type: "IoTBroker", IsPattern: false}
	entities = append(entities, entity)
	registration.EntityIdList = entities

	metadataList := make([]ContextMetadata, 0)

	metadata := ContextMetadata{}
	metadata.Name = "location"
	metadata.Type = "point"
	location := Point{Latitude: tb.MyLocation.Latitude, Longitude: tb.MyLocation.Longitude}
	metadata.Value = location
	metadataList = append(metadataList, metadata)

	registration.Metadata = metadataList

	registration.ProvidingApplication = tb.MyURL

	registerCtxReq.ContextRegistrations = append(registerCtxReq.ContextRegistrations, registration)
	registerCtxReq.Duration = "PT10M"

	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	_, err := client.RegisterContext(&registerCtxReq)
	if err != nil {
		ERROR.Println("not able to register myself to IoT Discovery: ", tb.myEntityId, ", error information: ", err)
		return false
	}

	INFO.Println("already registered myself to IoT Discovery: ", tb.myEntityId, " , ", tb.IoTDiscoveryURL)
	return true
}

func (tb *ThinBroker) deregisterMyself() {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.UnregisterEntity(tb.myEntityId)
	if err != nil {
		ERROR.Println(err)
	}

	INFO.Println("deregister myself to IoT Discovery: ", tb.myEntityId)
}

func (tb *ThinBroker) getEntities() []ContextElement {
	tb.entities_lock.RLock()
	defer tb.entities_lock.RUnlock()

	entities := make([]ContextElement, 0)

	for _, entity := range tb.entities {
		entities = append(entities, *entity)
	}

	return entities
}

func (tb *ThinBroker) getEntity(eid string) *ContextElement {
	tb.entities_lock.RLock()
	defer tb.entities_lock.RUnlock()

	if entity, exist := tb.entities[eid]; exist {
		element := ContextElement{}

		element.Entity = entity.Entity
		//element.AttributeDomainName = entity.AttributeDomainName
		element.Attributes = make([]ContextAttribute, len(entity.Attributes))
		copy(element.Attributes, entity.Attributes)
		element.Metadata = make([]ContextMetadata, len(entity.Metadata))
		copy(element.Metadata, entity.Metadata)

		return &element
	}

	return nil
}

func (tb *ThinBroker) deleteEntity(eid string) error {
	DEBUG.Println(" TO REMOVE ENTITY ", eid)

	//remove it from the local entity map
	tb.entities_lock.Lock()
	delete(tb.entities, eid)
	tb.entities_lock.Unlock()

	// inform the subscribers that this entity is deleted by sending a empty context element without any attribute, metadata
	emptyElement := ContextElement{}
	emptyElement.Entity.ID = eid
	tb.notifySubscribers(&emptyElement, false)

	//unregister this entity from IoT Discovery
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.UnregisterEntity(eid)
	if err != nil {
		ERROR.Println(err)
		return err
	}

	return nil
}

func (tb *ThinBroker) getAttribute(eid string, attrname string) *ContextAttribute {
	tb.entities_lock.RLock()
	defer tb.entities_lock.RUnlock()

	if entity, exist := tb.entities[eid]; exist {
		for _, attribute := range entity.Attributes {
			if attribute.Name == attrname {
				return &attribute
			}
		}
	}

	return nil
}

func (tb *ThinBroker) getSubscriptions() map[string]SubscribeContextRequest {
	tb.subscriptions_lock.RLock()
	defer tb.subscriptions_lock.RUnlock()

	subscriptions := make(map[string]SubscribeContextRequest)

	for sid, sub := range tb.subscriptions {
		subscriptions[sid] = *sub
	}

	return subscriptions
}

/*
	Get all NGSIV2  subscriptions
*/
func (tb *ThinBroker) getv2Subscriptions() map[string]SubscriptionRequest {
	tb.v2subscriptions_lock.RLock()
	defer tb.v2subscriptions_lock.RUnlock()
	v2subscriptions := make(map[string]SubscriptionRequest)
	for sid, sub := range tb.v2subscriptions {
		v2subscriptions[sid] = *sub
	}
	return v2subscriptions
}

func (tb *ThinBroker) getSubscription(sid string) *SubscribeContextRequest {
	tb.subscriptions_lock.RLock()
	defer tb.subscriptions_lock.RUnlock()

	if sub, exist := tb.subscriptions[sid]; exist {
		found := *sub
		return &found
	}

	return nil
}

/*
	Get NGSIV2 subscription by subscription id
*/

func (tb *ThinBroker) getv2Subscription(sid string) *SubscriptionRequest {
	tb.v2subscriptions_lock.RLock()
	defer tb.v2subscriptions_lock.RUnlock()
	if sub, exist := tb.v2subscriptions[sid]; exist {
		found := *sub
		return &found
	}
	return nil
}

func (tb *ThinBroker) deleteSubscription(sid string) error {
	tb.subscriptions_lock.Lock()
	defer tb.subscriptions_lock.Unlock()
	tb.subLinks_lock.RLock()
	defer tb.subLinks_lock.RUnlock()

	//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
	for index, otherSubID := range tb.main2Other[sid] {
		if index == 0 {
			tb.UnsubscribeContextAvailability(otherSubID)
		} else {
			unsubscribeContextProvider(otherSubID, tb.subscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
		}
	}

	// remove the subscription from the map
	delete(tb.subscriptions, sid)

	return nil
}

/*
	Delete subscription by subscriptionId
*/

func (tb *ThinBroker) deletev2Subscription(sid string) error {
	tb.v2subscriptions_lock.Lock()
	defer tb.v2subscriptions_lock.Unlock()
	tb.subLinks_lock.RLock()
	defer tb.subLinks_lock.RUnlock()

	//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
	for index, otherSubID := range tb.main2Other[sid] {
		if index == 0 {
			tb.Unsubscribev2ContextAvailability(otherSubID)
		} else {
			unsubscribev2ContextProvider(otherSubID, tb.v2subscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
		}
	}

	// remove the subscription from the map
	if _, oK := tb.v2subscriptions[sid]; oK {
		delete(tb.v2subscriptions, sid)
	}

	return nil
}

func (tb *ThinBroker) QueryContext(w rest.ResponseWriter, r *rest.Request) {
	queryCtxReq := QueryContextRequest{}
	err := r.DecodeJsonPayload(&queryCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matchedCtxElement := make([]ContextElement, 0)

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		// handle the query from another broker
		for _, eid := range queryCtxReq.Entities {
			tb.entities_lock.RLock()
			if element, exist := tb.entities[eid.ID]; exist {
				matchedCtxElement = append(matchedCtxElement, *element)
			}
			tb.entities_lock.RUnlock()
		}
	} else { // handle the query from an external consumer
		// discover the availability of all matched entities
		entityMap := tb.discoveryEntities(queryCtxReq.Entities, queryCtxReq.Attributes, queryCtxReq.Restriction)

		// fetch all matched entities from their providers
		for providerURL, entityList := range entityMap {
			if providerURL == tb.MyURL {
				for _, eid := range entityList {
					tb.entities_lock.RLock()
					if element, exist := tb.entities[eid.ID]; exist {
						returnedElement := element.CloneWithSelectedAttributes(queryCtxReq.Attributes)
						matchedCtxElement = append(matchedCtxElement, *returnedElement)
					}
					tb.entities_lock.RUnlock()
				}
			} else {
				elements := tb.fetchEntities(entityList, providerURL)
				matchedCtxElement = append(matchedCtxElement, elements...)
			}
		}
	}

	// send out the response
	queryCtxResp := QueryContextResponse{}

	ContextResponses := make([]ContextElementResponse, 0)
	for _, ctxElem := range matchedCtxElement {
		ctxElemResp := ContextElementResponse{}
		ctxElemResp.StatusCode.Code = 200
		ctxElemResp.ContextElement = ctxElem

		ContextResponses = append(ContextResponses, ctxElemResp)
	}
	queryCtxResp.ContextResponses = ContextResponses

	queryCtxResp.ErrorCode.Code = 200
	queryCtxResp.ErrorCode.ReasonPhrase = "OK"
	w.WriteJson(&queryCtxResp)
}

func (tb *ThinBroker) discoveryEntities(ids []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	discoverCtxAvailabilityReq := DiscoverContextAvailabilityRequest{}
	discoverCtxAvailabilityReq.Entities = ids
	discoverCtxAvailabilityReq.Attributes = attributes
	discoverCtxAvailabilityReq.Restriction = restriction

	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	registrationList, _ := client.DiscoverContextAvailability(&discoverCtxAvailabilityReq)

	result := make(map[string][]EntityId)
	for _, registration := range registrationList {
		reference := registration.ProvidingApplication
		entities := registration.EntityIdList
		if entityList, exist := result[reference]; exist {
			result[reference] = append(result[reference], entityList...)
		} else {
			result[reference] = make([]EntityId, 0)
			result[reference] = append(result[reference], entities...)
		}
	}

	return result
}

func (tb *ThinBroker) LDQueryContext(w rest.ResponseWriter, r *rest.Request) {
	reqBytes, _ := ioutil.ReadAll(r.Body)
	var LDqueryCtxReq interface{}
	err := json.Unmarshal(reqBytes, &LDqueryCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var fiwareService, fiwareServicePath string

	if r.Header.Get("fiware-service") != "" {
		fiwareService = r.Header.Get("fiware-service")
	} else {
		fiwareService = "default"
	}
	if r.Header.Get("fiware-servicepath") != "" {
		fiwareServicePath = r.Header.Get("fiware-servicepath")
	} else {
		fiwareServicePath = "default"
	}

	//fsp := r.Header.Get("fiware-servicepath")
	cType := r.Header.Get("Content-Type")
	link := r.Header.Get("Link")
	context, contextInpayload := extractcontext(cType, link)
	resolved, err := tb.ExpandPayload(LDqueryCtxReq, context, contextInpayload)
	LDQueryContext := LDQueryContextRequest{}
	var resolveError error
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		sz := Serializer{}
		LDQueryContext, resolveError = sz.uploadQueryContext(resolved, fiwareService, context)
		if resolveError != nil {
			rest.Error(w, resolveError.Error(), 400)
			return
		}
	}
	matchedCtxElement := make([]interface{}, 0)

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		tb.ldEntities_lock.Lock()
		for _, eid := range LDQueryContext.Entities {
			EID := eid.ID
			if element, exist := tb.ldEntities[EID]; exist {
				matchedCtxElement = append(matchedCtxElement, element)
			}
		}
		tb.ldEntities_lock.Unlock()
	} else {
		entityMap := tb.ldDiscoveryEntities(LDQueryContext)
		for providerURL, entityList := range entityMap {
			if providerURL == tb.MyURL {
				for _, eid := range entityList {
					tb.ldEntities_lock.Lock()
					if element, exist := tb.ldEntities[eid.ID]; exist {
						matchedCtxElement = append(matchedCtxElement, element)
					}
					tb.ldEntities_lock.Unlock()
				}
			} else {
				elements := tb.fetchLDEntities(entityList, providerURL, fiwareService, fiwareServicePath)
				matchedCtxElement = append(matchedCtxElement, elements...)
			}
		}

	}
	queryContextResponse := make([]interface{}, 0)
	for _, val := range matchedCtxElement {
		responseEle := val.(map[string]interface{})
		delete(responseEle, "fiwareServicePath")
		fmt.Println("responseEle",responseEle)
		if _ , ok := responseEle["@id"] ; ok == true {
			eid := responseEle["@id"].(string)
			actualEid := strings.Split(eid, "@")
			responseEle["@id"] = actualEid[0]
		}
		returnvalue, err := compactData(responseEle, responseEle["@context"])
		if err != nil {
			continue
		}
		queryContextResponse = append(queryContextResponse, returnvalue)
	}
	w.WriteHeader(200)
	w.WriteJson(queryContextResponse)
}

func (tb *ThinBroker) ldDiscoveryEntities(ldQueryContext LDQueryContextRequest) map[string][]EntityId {
	discoverCtxAvailabilityReq := DiscoverContextAvailabilityRequest{}
	discoverCtxAvailabilityReq.Entities = ldQueryContext.Entities
	discoverCtxAvailabilityReq.Attributes = ldQueryContext.Attributes
	discoverCtxAvailabilityReq.Restriction = ldQueryContext.Restriction
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	registrationList, _ := client.DiscoverContextAvailability(&discoverCtxAvailabilityReq)

	result := make(map[string][]EntityId)
	for _, registration := range registrationList {
		reference := registration.ProvidingApplication
		entities := registration.EntityIdList
		if entityList, exist := result[reference]; exist {
			result[reference] = append(result[reference], entityList...)
		} else {
			result[reference] = make([]EntityId, 0)
			result[reference] = append(result[reference], entities...)
		}
	}

	return result
}

func (tb *ThinBroker) fetchLDEntities(ids []EntityId, providerURL string, fs string, fsp string) []interface{} {
	newEntityList := make([]EntityId, 0)
	for _, entity := range ids {
		id := entity.ID
		idSplit := strings.Split(id, "@")
		entity.ID = idSplit[0]
		newEntityList = append(newEntityList, entity)
	}
	queryCtxLDReq := LDQueryContextRequest{}
	queryCtxLDReq.Entities = newEntityList
	queryCtxLDReq.Type = "Query"
	client := NGSI10Client{IoTBrokerURL: providerURL, SecurityCfg: tb.SecurityCfg}
	ctxElementList, _ := client.InternalLDQueryContext(&queryCtxLDReq, fs, fsp)
	return ctxElementList
}

func (tb *ThinBroker) fetchEntities(ids []EntityId, providerURL string) []ContextElement {
	queryCtxReq := QueryContextRequest{}
	queryCtxReq.Entities = ids

	client := NGSI10Client{IoTBrokerURL: providerURL, SecurityCfg: tb.SecurityCfg}
	ctxElementList, _ := client.InternalQueryContext(&queryCtxReq)
	return ctxElementList
}

func (tb *ThinBroker) UpdateContext(w rest.ResponseWriter, r *rest.Request) {
	updateCtxReq := UpdateContextRequest{}
	err := r.DecodeJsonPayload(&updateCtxReq)
	if err != nil {
		DEBUG.Println("not able to decode the orion updates")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		tb.handleInternalUpdateContext(&updateCtxReq)
	}

	//Southbound feature addition
	if r.Header.Get("fiware-service") != "" && r.Header.Get("fiware-servicepath") != "" {
		fs := r.Header.Get("fiware-service")
		fsp := r.Header.Get("fiware-servicepath")
		tb.handleExternalUpdateContext(w, &updateCtxReq, true, fs, fsp)
	} else {
		tb.handleExternalUpdateContext(w, &updateCtxReq, false)
	}
}

// handle context updates from external applications/devices

func (tb *ThinBroker) handleInternalUpdateContext(updateCtxReq *UpdateContextRequest) {
	switch updateCtxReq.UpdateAction {
	case "UPDATE":
		for _, ctxElem := range updateCtxReq.ContextElements {
			tb.UpdateContext2LocalSite(&ctxElem)
		}
	case "DELETE":
		for _, ctxElem := range updateCtxReq.ContextElements {
			tb.deleteEntity(ctxElem.Entity.ID)
		}
	}
}

// handle context updates forwarded by IoT Discovery
func (tb *ThinBroker) handleExternalUpdateContext(w rest.ResponseWriter, updateCtxReq *UpdateContextRequest, fiwareHeadersExist bool, params ...string) {
	// perform the update action accordingly
	switch strings.ToUpper(updateCtxReq.UpdateAction) {
	case "UPDATE", "APPEND":
		for _, ctxElem := range updateCtxReq.ContextElements {
			// just in case this is orion ngsi v1
			ctxElem.SetEntityID()
			// params[0] has FiwareService header and params[1] has FiwareServicePath
			if fiwareHeadersExist {
				tb.updateIdWithFiwareHeaders(&ctxElem, params[0], params[1])
			}

			brokerURL := tb.queryOwnerOfEntity(ctxElem.Entity.ID)

			if brokerURL == tb.myProfile.MyURL {
				tb.UpdateContext2LocalSite(&ctxElem, w)
			} else {
				tb.UpdateContext2RemoteSite(&ctxElem, updateCtxReq.UpdateAction, brokerURL)
			}
		}

	case "DELETE":
		for _, ctxElem := range updateCtxReq.ContextElements {
			brokerURL := tb.queryOwnerOfEntity(ctxElem.Entity.ID)
			if brokerURL == tb.myProfile.MyURL {
				tb.deleteEntity(ctxElem.Entity.ID)
			} else {
				tb.UpdateContext2RemoteSite(&ctxElem, updateCtxReq.UpdateAction, brokerURL)
			}
		}
	}
	//Send out the response
	w.WriteHeader(200)
	updateCtxResp := UpdateContextResponse{}
	w.WriteJson(&updateCtxResp)
}

//Southbound feature addition
func (tb *ThinBroker) handleSouthboundCommand(w rest.ResponseWriter, ctxElem *ContextElement) {
	rid := ctxElem.Entity.ID
	//Get Provider IoT Agent for the registered device on local broker
	tb.fiwareData_lock.RLock()
	fiData := tb.fiwareData[rid]
	tb.fiwareData_lock.RUnlock()

	if fiData != nil {
		providerURL := fiData.ProviderIoTAgent + "/ngsi10"
		fs := fiData.FiwareService
		fsp := fiData.FiwareServicePath
		//Extract actual Element ID before sending the Context Element to IoT Agent
		tb.removeFiwareHeadersFromId(ctxElem, fs, fsp)
		tb.FogflowToFiwareContextElement(ctxElem)
		DEBUG.Println("Handling command update through local broker.")
		DEBUG.Println(providerURL)
		client := NGSI10Client{IoTBrokerURL: providerURL, SecurityCfg: tb.SecurityCfg}
		client.SouthboundUpdateContext(ctxElem, fs, fsp)
	} else {
		rest.Error(w, "The device registration was not found!", 404)
		return
	}
}

//Southbound Feature addition
func (tb *ThinBroker) FogflowToFiwareContextElement(ctxElem *ContextElement) {
	ctxElem.ID = ctxElem.Entity.ID
	ctxElem.Type = ctxElem.Entity.Type

	if ctxElem.Entity.IsPattern == true {
		ctxElem.IsPattern = "true"
	} else if ctxElem.Entity.IsPattern == false {
		ctxElem.IsPattern = "false"
	}
	ctxElem.Entity = EntityId{}
}

func (tb *ThinBroker) queryOwnerOfEntity(eid string) string {
	inLocalBroker := true

	tb.entities_lock.RLock()
	_, exist := tb.entities[eid]
	inLocalBroker = exist
	tb.entities_lock.RUnlock()

	if inLocalBroker == true {
		return tb.myProfile.MyURL
	} else {
		client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
		brokerURL, _ := client.GetProviderURL(eid)
		if brokerURL == "" {
			return tb.myProfile.MyURL
		}
		return brokerURL
	}
}

func (tb *ThinBroker) UpdateContext2LocalSite(ctxElem *ContextElement, params ...rest.ResponseWriter) {
	command := false
	//If any of the attributes is of type "command", all the attributes will be considered as commands
	for _, attr := range ctxElem.Attributes {
		if attr.Type == "command" {
			command = true
			break
		}
	}

	if command == true {
		tb.handleSouthboundCommand(params[0], ctxElem)
	} else {
		tb.entities_lock.Lock()
		eid := ctxElem.Entity.ID
		hasUpdatedMetadata := hasUpdatedMetadata(ctxElem, tb.entities[eid])
		tb.entities_lock.Unlock()

		// apply the new update to the entity in the entity map
		tb.updateContextElement(ctxElem)

		// propogate this update to its subscribers
		go tb.notifySubscribers(ctxElem, true)

		//propagate this update to its ngsiv2 subscribers
		go tb.notifySubscribersV2(ctxElem, true)

		// register the entity if there is any changes on attribute list, domain metadata
		if hasUpdatedMetadata == true {
			tb.registerContextElement(ctxElem)
		}
	}
}

func (tb *ThinBroker) UpdateContext2RemoteSite(ctxElem *ContextElement, updateAction string, brokerURL string) {
	switch updateAction {
	case "UPDATE":
		INFO.Println(brokerURL)
		client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: tb.SecurityCfg}
		client.UpdateContext(ctxElem)

	case "DELETE":
		client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: tb.SecurityCfg}
		client.DeleteContext(&ctxElem.Entity)
	}
}

func (tb *ThinBroker) NotifyLdContext(w rest.ResponseWriter, r *rest.Request) {
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		var context []interface{}
		context = append(context, DEFAULT_CONTEXT)
		//notifyElement, _ := tb.getStringInterfaceMap(r)
		reqBytes, _ := ioutil.ReadAll(r.Body)
		fmt.Println("reqBytes",string(reqBytes))
		var notifyRequest interface{}

		err := json.Unmarshal(reqBytes, &notifyRequest)
		var fiwareService string
		//var fiwareServicePath string
		if r.Header.Get("fiware-service") != "" {
			fiwareService = r.Header.Get("fiware-service")
		} else {
			fiwareService = "default"
		}
		if err != nil {
			err := errors.New("not able to decode  orion update")
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notifyElement := notifyRequest.(map[string]interface{})
		notifyElemtData := notifyElement["data"]
		var notifyEleDatamap []interface{}
		if notifyElemtData != nil {
			notifyEleDatamap = notifyElemtData.([]interface{})
		} else {
			fmt.Println("nil data repersent entity has been deleted")
		}
		notifyCtxResp := NotifyContextResponse{}
		w.WriteJson(&notifyCtxResp)
		Link := DEFAULT_CONTEXT
		for _, data := range notifyEleDatamap {
			notifyData := data.(map[string]interface{})
			notifyData["@context"] = context
			expand, _ := tb.ExpandData(notifyData)
			sz := Serializer{}
			deSerializedEntity, err := sz.DeSerializeEntity(expand)
			if err != nil {
				continue
			} else {
				deSerializedEntity["@id"] = getIoTID(deSerializedEntity["@id"].(string), fiwareService)
				var fsp string
				if r.Header.Get("fiware-servicePath") != "" {
					fsp = r.Header.Get("fiware-servicePath")
				} else {
					//deSerializedEntity["fiwareServicePath"] = "default"
					fsp = "default"
				}
				deSerializedEntity["fiwareServicePath"] = fsp

				deSerializedEntity["@context"] = context
				//deSerializedEntity["createdAt"] = time.Now().String()
				eid := deSerializedEntity["@id"].(string)
				tb.LDe2sub_lock.RLock()
				subscriberList := tb.entityId2LDSubcriptions[eid]
				tb.LDe2sub_lock.RUnlock()
				if len(subscriberList) > 0 {
					//send the notification to subscriber
					go tb.LDNotifySubscribers(deSerializedEntity, false)
				} else {
					tb.handleLdExternalUpdateContext(deSerializedEntity, Link)
				}
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}
func (tb *ThinBroker) NotifyContext(w rest.ResponseWriter, r *rest.Request) {
	notifyCtxReq := NotifyContextRequest{}
	err := r.DecodeJsonPayload(&notifyCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send out the response
	notifyCtxResp := NotifyContextResponse{}
	w.WriteJson(&notifyCtxResp)

	// inform its subscribers
	for _, ctxResp := range notifyCtxReq.ContextResponses {
		go tb.notifySubscribers(&ctxResp.ContextElement, false)
	}
}

/*
	This function will return true if there will be any update in updateContext element for subscription condition attributes
*/

func (tb *ThinBroker) checkMatchedAttr(ctxElemAttrs []string, sid string) bool {
	tb.v2subscriptions_lock.RLock()
	_, res := tb.v2subscriptions[sid]
	if res == false {
		tb.v2subscriptions_lock.RUnlock()
		return false
	}
	conditionList := tb.v2subscriptions[sid].Subject.Conditions.Attrs
	tb.v2subscriptions_lock.RUnlock()
	matchedAtleastOnce := false
	for _, attrs1 := range ctxElemAttrs {
		for _, attrs2 := range conditionList {
			if attrs1 == attrs2 {
				matchedAtleastOnce = true
				break
			}
			if matchedAtleastOnce == true {
				break
			}
		}
	}
	return matchedAtleastOnce
}

/*
	Clone the attribute for new update
	Send the ReliableNotify if there will be any change in condition attribute of subscription request
*/

func (tb *ThinBroker) notifySubscribersV2(ctxElem *ContextElement, checkSelectedAttributes bool) {
	eid := ctxElem.Entity.ID
	subscriberList := tb.entityIdv2Subcriptions[eid]
	ctxAttrsName := make([]string, 0)
	ctxAttrs := ctxElem.Attributes
	for _, ctxAttrsEle := range ctxAttrs {
		ctxAttrsName = append(ctxAttrsName, ctxAttrsEle.Name)
	}
	for _, sid := range subscriberList {
		elements := make([]ContextElement, 0)
		checkCondition := tb.checkMatchedAttr(ctxAttrsName, sid)
		if checkSelectedAttributes == true && checkCondition == true {
			selectedAttributes := make([]string, 0)
			tb.v2subscriptions_lock.RLock() // change the lock here
			if v2subscription, exist := tb.v2subscriptions[sid]; exist {
				if v2subscription.Notification.Attrs != nil {
					selectedAttributes = append(selectedAttributes, tb.v2subscriptions[sid].Notification.Attrs...)
				}
			}
			tb.v2subscriptions_lock.RUnlock()
			tb.entities_lock.RLock()
			element := tb.entities[eid].CloneWithSelectedAttributes(selectedAttributes)
			tb.entities_lock.RUnlock()
			elements = append(elements, *element)
		} else {
			elements = append(elements, *ctxElem)
		}
		if checkCondition == true {
			go tb.sendReliableNotify(elements, sid)
		}
	}
}

func (tb *ThinBroker) notifySubscribers(ctxElem *ContextElement, checkSelectedAttributes bool) {
	eid := ctxElem.Entity.ID
	tb.e2sub_lock.RLock()
	defer tb.e2sub_lock.RUnlock()
	subscriberList := tb.entityId2Subcriptions[eid]
	//send this context element to the subscriber
	for _, sid := range subscriberList {
		elements := make([]ContextElement, 0)

		if checkSelectedAttributes == true {
			selectedAttributes := make([]string, 0)

			tb.subscriptions_lock.RLock()

			if subscription, exist := tb.subscriptions[sid]; exist {
				if subscription.Attributes != nil {
					selectedAttributes = append(selectedAttributes, tb.subscriptions[sid].Attributes...)
				}
			}

			tb.subscriptions_lock.RUnlock()

			tb.entities_lock.RLock()
			element := tb.entities[eid].CloneWithSelectedAttributes(selectedAttributes)
			tb.entities_lock.RUnlock()

			elements = append(elements, *element)
		} else {
			elements = append(elements, *ctxElem)
		}
		go tb.sendReliableNotify(elements, sid)
	}
}

/*func (tb *ThinBroker) notifyOneSubscriberWithCurrentStatus(entities []EntityId, sid string) {
	// check if the subscription still exists; if yes, then find out the selected attribute list
	tb.subscriptions_lock.RLock()

	v1Subscription, ok := tb.subscriptions[sid]
	if ok == false {
		tb.subscriptions_lock.RUnlock()
		tb.ldSubscriptions_lock.RLock()
		ldSubscription, ldOK := tb.ldSubscriptions[sid]

		if ldOK == false {
			tb.ldSubscriptions_lock.RUnlock()
			return
		}
		selectedAttributes := ldSubscription.WatchedAttributes
		tb.ldSubscriptions_lock.RUnlock()
		tb.notifyOneSubscriberWithCurrentStatusOfLD(entities, sid, selectedAttributes)
	} else {
		selectedAttributes := v1Subscription.Attributes
		tb.subscriptions_lock.RUnlock()
		tb.notifyOneSubscriberWithCurrentStatusOfV1(entities, sid, selectedAttributes)
	}
}*/

func (tb *ThinBroker) notifyOneSubscriberWithCurrentStatus(entities []EntityId, sid string) {
	elements := make([]ContextElement, 0)
	// check if the subscription still exists; if yes, then find out the selected attribute list
	tb.subscriptions_lock.RLock()

	subscription, ok := tb.subscriptions[sid]
	if ok == false {
		tb.subscriptions_lock.RUnlock()
		return
	}
	selectedAttributes := subscription.Attributes
	tb.subscriptions_lock.RUnlock()

	tb.entities_lock.Lock()
	for _, entity := range entities {
		if element, exist := tb.entities[entity.ID]; exist {
			returnedElement := element.CloneWithSelectedAttributes(selectedAttributes)
			elements = append(elements, *returnedElement)
		}
	}
	tb.entities_lock.Unlock()
	go tb.sendReliableNotify(elements, sid)
}

func (tb *ThinBroker) notifyOneLDSubscriberWithCurrentStatus(entities []EntityId, sid string) {
	elements := make([]map[string]interface{}, 0)
	tb.ldSubscriptions_lock.RLock()
	ldSubscription, oK := tb.ldSubscriptions[sid]
	if oK == false {
		tb.ldSubscriptions_lock.RUnlock()
		return
	}
	selectedAttributes := ldSubscription.Notification.Attributes
	tb.ldSubscriptions_lock.RUnlock()
	tb.ldEntities_lock.Lock()
	for _, entity := range entities {
		if element, exist := tb.ldEntities[entity.ID]; exist {
			elementMap := element.(map[string]interface{})
			returnedElement := ldCloneWithSelectedAttributes(elementMap, selectedAttributes)
			fmt.Println("returnedElement",returnedElement)
			elements = append(elements, returnedElement)
		}
	}
	tb.ldEntities_lock.Unlock()
	go tb.sendReliableNotifyToNgsiLDSubscriber(elements, sid)
}
func (tb *ThinBroker) notifyOneSubscriberWithCurrentStatusOfV1(entities []EntityId, sid string, selectedAttributes []string) {
	// Create NGSIv1 Context Element
	elements := make([]ContextElement, 0)

	tb.entities_lock.Lock()
	for _, entity := range entities {
		if element, exist := tb.entities[entity.ID]; exist {
			returnedElement := element.CloneWithSelectedAttributes(selectedAttributes)
			elements = append(elements, *returnedElement)
		}
	}
	tb.entities_lock.Unlock()
	go tb.sendReliableNotify(elements, sid)
}

/*
	Send the current status of exiting entity to the new subscriber
*/

func (tb *ThinBroker) notifyOneSubscriberv2WithCurrentStatus(entities []EntityId, sid string) {
	elements := make([]ContextElement, 0)
	// check if the subscription still exists; if yes, then find out the selected attribute list
	tb.v2subscriptions_lock.RLock()

	v2subscription, ok := tb.v2subscriptions[sid]
	if ok == false {
		tb.v2subscriptions_lock.RUnlock()
		return
	}
	selectedAttributes := v2subscription.Notification.Attrs
	tb.v2subscriptions_lock.RUnlock()

	tb.entities_lock.Lock()
	for _, entity := range entities {
		if element, exist := tb.entities[entity.ID]; exist {
			returnedElement := element.CloneWithSelectedAttributes(selectedAttributes)
			elements = append(elements, *returnedElement)
		}
	}
	tb.entities_lock.Unlock()
	go tb.sendReliableNotify(elements, sid)
}

/*
	Send the Notification to NGSIV1 subscriber
*/

func (tb *ThinBroker) sendReliableNotifyToNgsiv1Subscriber(elements []ContextElement, sid string) {
	tb.subscriptions_lock.Lock()
	subscription, ok := tb.subscriptions[sid]
	if ok == false {
		tb.subscriptions_lock.Unlock()
	}
	subscriberURL := subscription.Reference
	IsOrionBroker := subscription.Subscriber.IsOrion
	if subscription.Subscriber.RequireReliability == true && len(subscription.Subscriber.NotifyCache) > 0 {
		DEBUG.Println("resend notify:  ", len(subscription.Subscriber.NotifyCache))
		for _, pCtxElem := range subscription.Subscriber.NotifyCache {
			elements = append(elements, *pCtxElem)
		}
		subscription.Subscriber.NotifyCache = make([]*ContextElement, 0)
	}
	tb.subscriptions_lock.Unlock()
	err := postNotifyContext(elements, sid, subscriberURL, IsOrionBroker, tb.SecurityCfg)
	INFO.Println("NOTIFY: ", len(elements), ", ", sid, ", ", subscriberURL, ", ", IsOrionBroker)
	if err != nil {
		INFO.Println("NOTIFY is not received by the subscriber, ", subscriberURL)

		tb.subscriptions_lock.Lock()
		if subscription, exist := tb.subscriptions[sid]; exist {
			if subscription.Subscriber.RequireReliability == true {
				for _, ctxElem := range elements {
					subscription.Subscriber.NotifyCache = append(subscription.Subscriber.NotifyCache, &ctxElem)
				}

				tb.tmpNGSI10NotifyCache = append(tb.tmpNGSI10NotifyCache, sid)
			}
		}
		tb.subscriptions_lock.Unlock()
	}

}

/*
	Send the Notification to NGSIV2 subscriber
*/

func (tb *ThinBroker) sendReliableNotifyToNgsiv2Subscriber(elements []ContextElement, sid string) {
	tb.v2subscriptions_lock.Lock()
	v2subscription, ok := tb.v2subscriptions[sid]
	if ok == false {
		tb.v2subscriptions_lock.Unlock()
	}
	subscriberURL := v2subscription.Notification.Http.Url
	if v2subscription.Subscriber.RequireReliability == true && len(v2subscription.Subscriber.NotifyCache) > 0 {
		DEBUG.Println("resend notify:  ", len(v2subscription.Subscriber.NotifyCache))
		for _, pCtxElem := range v2subscription.Subscriber.NotifyCache {
			elements = append(elements, *pCtxElem)
		}
		v2subscription.Subscriber.NotifyCache = make([]*ContextElement, 0)
	}
	tb.v2subscriptions_lock.Unlock()
	err := postNotifyContext(elements, sid, subscriberURL, true, tb.SecurityCfg)
	INFO.Println("NOTIFY: ", len(elements), ", ", sid, ", ", subscriberURL, ", ", true)
	if err != nil {
		INFO.Println("NOTIFY is not received by the subscriber, ", subscriberURL)

		tb.v2subscriptions_lock.Lock()
		if v2subscription, exist := tb.v2subscriptions[sid]; exist {
			if v2subscription.Subscriber.RequireReliability == true {
				for _, ctxElem := range elements {
					v2subscription.Subscriber.NotifyCache = append(v2subscription.Subscriber.NotifyCache, &ctxElem)
				}

				tb.tmpNGSIv2NotifyCache = append(tb.tmpNGSIv2NotifyCache, sid)
			}
		}
		tb.v2subscriptions_lock.Unlock()
	}

}

/*
	Identify the subscriber(NGSIV1 or NGSIV2) by using SubscriptionId
*/

func (tb *ThinBroker) sendReliableNotify(elements []ContextElement, sid string) {
	tb.subscriptions_lock.Lock()
	_, ok := tb.subscriptions[sid]
	if ok == true {
		tb.subscriptions_lock.Unlock()
		tb.sendReliableNotifyToNgsiv1Subscriber(elements, sid)
	} else {
		tb.subscriptions_lock.Unlock()
	}
	tb.v2subscriptions_lock.Lock()
	_, ok = tb.v2subscriptions[sid]
	if ok == true {
		tb.v2subscriptions_lock.Unlock()
		tb.sendReliableNotifyToNgsiv2Subscriber(elements, sid)
	} else {
		tb.v2subscriptions_lock.Unlock()
	}
}

func (tb *ThinBroker) updateContextElement(ctxElem *ContextElement) {
	//look up who already subscribed to this context element
	eid := ctxElem.Entity.ID

	tb.entities_lock.Lock()
	defer tb.entities_lock.Unlock()

	// update its value in the entity map
	if curElement, exist := tb.entities[eid]; exist {
		for _, attr := range ctxElem.Attributes {
			updateAttribute(&attr, curElement)
		}

		for _, metadata := range ctxElem.Metadata {
			updateDomainMetadata(&metadata, curElement)
		}
	} else {
		newContextElement := *ctxElem
		tb.entities[eid] = &newContextElement
	}
}

func (tb *ThinBroker) SubscribeContext(w rest.ResponseWriter, r *rest.Request) {
	subReq := SubscribeContextRequest{}
	subReq.Attributes = make([]string, 0)

	err := r.DecodeJsonPayload(&subReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// new SubscriptionID
	u1, err := uuid.NewUUID()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subID := u1.String()

	// send out the response
	subResp := SubscribeContextResponse{}
	subResp.SubscribeResponse.SubscriptionId = subID
	subResp.SubscribeError.SubscriptionId = subID
	w.WriteJson(&subResp)

	// check the request header
	if r.Header.Get("Destination") == "orion-broker" {
		subReq.Subscriber.IsOrion = true
	} else {
		subReq.Subscriber.IsOrion = false
	}

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		subReq.Subscriber.IsInternal = true
	} else {
		subReq.Subscriber.IsInternal = false
	}

	// check the required semantics of message delivery
	if r.Header.Get("Require-Reliability") == "true" {
		subReq.Subscriber.RequireReliability = true
		subReq.Subscriber.NotifyCache = make([]*ContextElement, 0)
	} else {
		subReq.Subscriber.RequireReliability = false
	}
	subReq.Subscriber.BrokerURL = tb.MyURL

	INFO.Printf("NEW subscription: %v\n", subReq)

	// add it into the subscription map
	tb.subscriptions_lock.Lock()
	tb.subscriptions[subID] = &subReq
	tb.subscriptions_lock.Unlock()
	// take actions
	if subReq.Subscriber.IsInternal == true {
		INFO.Println("internal subscription coming from another broker")

		for _, entity := range subReq.Entities {
			tb.e2sub_lock.Lock()
			tb.entityId2Subcriptions[entity.ID] = append(tb.entityId2Subcriptions[entity.ID], subID)
			tb.e2sub_lock.Unlock()
		}
		tb.notifyOneSubscriberWithCurrentStatus(subReq.Entities, subID)
	} else {
		tb.SubscribeContextAvailability(subID)
	}
}

/*
	NGSIV2 subscription request handler
*/

func (tb *ThinBroker) Subscriptionv2Context(w rest.ResponseWriter, r *rest.Request) {
	subReqv2 := SubscriptionRequest{}

	err := r.DecodeJsonPayload(&subReqv2)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u1, err := uuid.NewUUID()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subID := u1.String()
	// send out the response
	subRespv2 := Subscribev2Response{}
	subRespv2.SubscriptionResponse.SubscriptionId = subID
	subRespv2.SubscriptionError.SubscriptionId = subID
	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&subRespv2)

	subReqv2.Subscriber.BrokerURL = tb.MyURL

	INFO.Printf("NEW subscription: %v\n", subReqv2)

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		subReqv2.Subscriber.IsInternal = true
	} else {
		subReqv2.Subscriber.IsInternal = false
	}

	tb.v2subscriptions_lock.Lock()
	ctxEle := &subReqv2
	ctxEle.Subject.SetIDpattern()
	tb.v2subscriptions[subID] = ctxEle
	tb.v2subscriptions_lock.Unlock()
	if subReqv2.Subscriber.IsInternal == true {
		INFO.Println("internal subscription coming from another broker")
		for _, entity := range subReqv2.Subject.Entities {
			tb.e2sub_lock.Lock()
			tb.entityIdv2Subcriptions[entity.ID] = append(tb.entityIdv2Subcriptions[entity.ID], subID)
			tb.e2sub_lock.Unlock()
		}
		tb.notifyOneSubscriberv2WithCurrentStatus(subReqv2.Subject.Entities, subID)
	} else {
		tb.Subscribev2ContextAvailability(subID)
	}
}

/*
	Send request to discovery to check the SubscribeContextAvailability
*/

func (tb *ThinBroker) SubscribeContextAvailability(sid string) error {
	availabilitySubscription := SubscribeContextAvailabilityRequest{}

	tb.subscriptions_lock.RLock()
	availabilitySubscription.Entities = tb.subscriptions[sid].Entities
	availabilitySubscription.Attributes = tb.subscriptions[sid].Attributes
	availabilitySubscription.Duration = tb.subscriptions[sid].Duration
	availabilitySubscription.Restriction = tb.subscriptions[sid].Restriction
	tb.subscriptions_lock.RUnlock()

	availabilitySubscription.Reference = tb.MyURL + "/notifyContextAvailability"

	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	subscriptionId, err := client.SubscribeContextAvailability(&availabilitySubscription)
	if subscriptionId != "" {
		tb.subLinks_lock.Lock()
		tb.main2Other[sid] = append(tb.main2Other[sid], subscriptionId)
		tb.availabilitySub2MainSub[subscriptionId] = sid
		notifyMessage, alreadyBack := tb.tmpNGSI9NotifyCache[subscriptionId]
		tb.subLinks_lock.Unlock()
		if alreadyBack == true {
			INFO.Println("========forward the availability notify that arrive earlier===========")
			tb.handleNGSI9Notify(sid, notifyMessage)

			tb.subLinks_lock.Lock()
			delete(tb.tmpNGSI9NotifyCache, subscriptionId)
			tb.subLinks_lock.Unlock()
		}

		return nil
	} else {
		INFO.Println("failed to subscribe the availability of requested entities ", err)
		return err
	}
}

/*
	Send request to discovery to check the ContextAvailability for NGSIV2 subscriber
*/

func (tb *ThinBroker) Subscribev2ContextAvailability(sid string) error {
	availabilitySubscriptionv2 := SubscribeContextAvailabilityRequest{}

	tb.v2subscriptions_lock.RLock()
	availabilitySubscriptionv2.Entities = tb.v2subscriptions[sid].Subject.Entities
	availabilitySubscriptionv2.Attributes = tb.v2subscriptions[sid].Subject.Conditions.Attrs
	availabilitySubscriptionv2.Attributes = append(availabilitySubscriptionv2.Attributes, tb.v2subscriptions[sid].Notification.Attrs...)
	availabilitySubscriptionv2.Duration = tb.v2subscriptions[sid].Expires
	tb.v2subscriptions_lock.RUnlock()

	availabilitySubscriptionv2.Reference = tb.MyURL + "/notifyContextAvailabilityv2"
	client := NGSIV2Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	subscriptionId, err := client.Subscribev2ContextAvailability(&availabilitySubscriptionv2)
	if subscriptionId != "" {
		tb.subLinks_lock.Lock()
		tb.main2Other[sid] = append(tb.main2Other[sid], subscriptionId)
		tb.availabilitySub2MainSub[subscriptionId] = sid
		notifyMessage, alreadyBack := tb.tmpNGSIV2NotifyCache[subscriptionId]
		tb.subLinks_lock.Unlock()

		if alreadyBack == true {
			INFO.Println("========forward the availability notify that arrive earlier===========")
			tb.handleNGSIV2Notify(sid, notifyMessage)

			tb.subLinks_lock.Lock()
			delete(tb.tmpNGSIV2NotifyCache, subscriptionId)
			tb.subLinks_lock.Unlock()
		}

		return nil
	} else {
		INFO.Println("failed to subscribe the availability of requested entities ", err)
		return err
	}
}

func (tb *ThinBroker) UnsubscribeContextAvailability(sid string) error {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.UnsubscribeContextAvailability(sid)
	return err
}

func (tb *ThinBroker) UnsubscribeLDContext(w rest.ResponseWriter, r *rest.Request) {
	unsubscribeCtxReq := UnsubscribeContextRequest{}
	err := r.DecodeJsonPayload(&unsubscribeCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subID := unsubscribeCtxReq.SubscriptionId

	// send out the response
	unsubscribeCtxResp := UnsubscribeContextResponse{}
	unsubscribeCtxResp.StatusCode.Code = 200
	unsubscribeCtxResp.StatusCode.ReasonPhrase = "OK"
	w.WriteJson(&unsubscribeCtxResp)

	tb.ldSubscriptions_lock.Lock()
	defer tb.ldSubscriptions_lock.Unlock()
	tb.subLinks_lock.Lock()
	defer tb.subLinks_lock.Unlock()

	// check the request header
	if r.Header.Get("User-Agent") != "lightweight-iot-broker" {
		//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
		for index, otherSubID := range tb.main2Other[subID] {
			if index == 0 {
				tb.UnsubscribeContextAvailability(otherSubID)
			} else {
				unsubscribeContextProvider(otherSubID, tb.ldSubscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
			}
		}
	}

	// remove the subscription from the map
	delete(tb.ldSubscriptions, subID)
}

/*
	Unsubscribe to the NGSIV2 subscriber
*/

func (tb *ThinBroker) Unsubscribev2ContextAvailability(sid string) error {
	client := NGSIV2Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.Unsubscribev2ContextAvailability(sid)
	return err
}

func (tb *ThinBroker) UnsubscribeContext(w rest.ResponseWriter, r *rest.Request) {
	unsubscribeCtxReq := UnsubscribeContextRequest{}
	err := r.DecodeJsonPayload(&unsubscribeCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subID := unsubscribeCtxReq.SubscriptionId

	// send out the response
	unsubscribeCtxResp := UnsubscribeContextResponse{}
	unsubscribeCtxResp.StatusCode.Code = 200
	unsubscribeCtxResp.StatusCode.ReasonPhrase = "OK"
	w.WriteJson(&unsubscribeCtxResp)

	tb.subscriptions_lock.Lock()
	defer tb.subscriptions_lock.Unlock()
	tb.subLinks_lock.RLock()
	defer tb.subLinks_lock.RUnlock()

	// check the request header
	fmt.Println("======testing the code=====")
	if r.Header.Get("User-Agent") != "lightweight-iot-broker" {
		//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
		for index, otherSubID := range tb.main2Other[subID] {
			if index == 0 {
				tb.UnsubscribeContextAvailability(otherSubID)
			} else {
				if _, found := tb.subscriptions[otherSubID]; found {
					unsubscribeContextProvider(otherSubID, tb.subscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
				}
				if _, found := tb.ldSubscriptions[otherSubID]; found {
					unsubscribeContextProvider(otherSubID, tb.ldSubscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
				}
			}
		}
	}

	// remove the ngsiv1 subscription from the map
	if _, found := tb.subscriptions[subID]; found {
		delete(tb.subscriptions, subID)
		// remove the ngsild subscription from map
	} else if _, found := tb.ldSubscriptions[subID]; found {
		delete(tb.ldSubscriptions, subID)
	}

}

func (tb *ThinBroker) NotifyLDContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	notifyLDContextAvailabilityReq := NotifyContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&notifyLDContextAvailabilityReq)
	if err != nil {
		ERROR.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send out the response
	notifyLDContextAvailabilityResp := NotifyContextAvailabilityResponse{}
	notifyLDContextAvailabilityResp.ResponseCode.Code = 200
	notifyLDContextAvailabilityResp.ResponseCode.ReasonPhrase = "OK"
	w.WriteJson(&notifyLDContextAvailabilityResp)

	subID := notifyLDContextAvailabilityReq.SubscriptionId

	//map it to the main subscription
	tb.subLinks_lock.Lock()
	mainSubID, exist := tb.availabilitySub2MainSub[subID]
	if exist == false {
		DEBUG.Println("put it into the tempCache and handle it later")
		tb.tmpNGSILDNotifyCache[subID] = &notifyLDContextAvailabilityReq
	}
	tb.subLinks_lock.Unlock()

	if exist == true {
		tb.handleNGSILDNotify(mainSubID, &notifyLDContextAvailabilityReq)
	}
}

func (tb *ThinBroker) NotifyContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	notifyContextAvailabilityReq := NotifyContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&notifyContextAvailabilityReq)
	if err != nil {
		ERROR.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send out the response
	notifyContextAvailabilityResp := NotifyContextAvailabilityResponse{}
	notifyContextAvailabilityResp.ResponseCode.Code = 200
	notifyContextAvailabilityResp.ResponseCode.ReasonPhrase = "OK"
	w.WriteJson(&notifyContextAvailabilityResp)

	subID := notifyContextAvailabilityReq.SubscriptionId

	//map it to the main subscription
	tb.subLinks_lock.Lock()
	mainSubID, exist := tb.availabilitySub2MainSub[subID]
	if exist == false {
		DEBUG.Println("put it into the tempCache and handle it later")
		tb.tmpNGSI9NotifyCache[subID] = &notifyContextAvailabilityReq
	}
	tb.subLinks_lock.Unlock()

	if exist == true {
		tb.handleNGSI9Notify(mainSubID, &notifyContextAvailabilityReq)
	}
}

func (tb *ThinBroker) Notifyv2ContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	notifyv2ContextAvailabilityReq := Notifyv2ContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&notifyv2ContextAvailabilityReq)
	if err != nil {
		ERROR.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send out the response
	notifyv2ContextAvailabilityResp := Notifyv2ContextAvailabilityResponse{}
	notifyv2ContextAvailabilityResp.ResponseCode.Code = 200
	notifyv2ContextAvailabilityResp.ResponseCode.ReasonPhrase = "OK"
	w.WriteJson(&notifyv2ContextAvailabilityResp)

	subID := notifyv2ContextAvailabilityReq.SubscriptionId
	//map it to the main subscription
	tb.subLinks_lock.Lock()
	mainSubID, exist := tb.availabilitySub2MainSub[subID]
	if exist == false {
		DEBUG.Println("put it into the tempCache and handle it later")
		tb.tmpNGSIV2NotifyCache[subID] = &notifyv2ContextAvailabilityReq
	}
	tb.subLinks_lock.Unlock()

	if exist == true {
		tb.handleNGSIV2Notify(mainSubID, &notifyv2ContextAvailabilityReq)
	}
}

//handleNGSILDNotify

func (tb *ThinBroker) handleNGSILDNotify(mainSubID string, notifyContextAvailabilityReq *NotifyContextAvailabilityRequest) {
	var action string
	notifyContextAvailabilityReq.ErrorCode.Code = 301
	switch notifyContextAvailabilityReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	}

	INFO.Println(action, " subID ", mainSubID)

	for _, registrationResp := range notifyContextAvailabilityReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration
		for _, eid := range registration.EntityIdList {
			INFO.Println("===> ", eid, " , ", mainSubID)

			tb.LDe2sub_lock.Lock()
			if action == "CREATE" {
				tb.entityId2LDSubcriptions[eid.ID] = append(tb.entityId2LDSubcriptions[eid.ID], mainSubID) // change here
			} else if action == "DELETE" {
				subList := tb.entityId2LDSubcriptions[eid.ID] //change here
				for i, id := range subList {
					if id == mainSubID {
						tb.entityId2LDSubcriptions[eid.ID] = append(subList[:i], subList[i+1:]...) // change here
						break
					}
				}
			} else if action == "UPDATE" {
				existFlag := false
				for _, subID := range tb.entityId2LDSubcriptions[eid.ID] { //change here
					if subID == mainSubID {
						existFlag = true
						break
					}
				}
				if existFlag == false {
					tb.entityId2LDSubcriptions[eid.ID] = append(tb.entityId2LDSubcriptions[eid.ID], mainSubID) //chage here
				}
			}
			tb.LDe2sub_lock.Unlock()
		}

		INFO.Println(registration.ProvidingApplication, ", ", tb.MyURL)
		INFO.Println("TO ngsi10 subscription, ", mainSubID)
		INFO.Printf("entity list: %+v\r\n", registration.EntityIdList)
		if registration.ProvidingApplication == tb.MyURL {
			//for matched entities provided by myself
			if action == "CREATE" || action == "UPDATE" {
				tb.notifyOneLDSubscriberWithCurrentStatus(registration.EntityIdList, mainSubID)
			}
		} else {
			//for matched entities provided by other IoT Brokers
			newLDSubscription := LDSubscriptionRequest{}
			newLDSubscription.Entities = registration.EntityIdList
			newLDSubscription.Type = "Subscription"
			MyLDURL := strings.TrimSuffix(tb.MyURL, "/ngsi10")
			MyLDURL = MyLDURL + "/ngsi-ld/v1/notifyContext/"
			newLDSubscription.Notification.Endpoint.URI = MyLDURL
			newLDSubscription.Subscriber.BrokerURL = registration.ProvidingApplication

			if action == "CREATE" || action == "UPDATE" {
				registration.ProvidingApplication = strings.TrimSuffix(registration.ProvidingApplication, "/ngsi10")
				sid, err := subscriptionLDContextProvider(&newLDSubscription, registration.ProvidingApplication, tb.SecurityCfg)
				if err == nil {
					INFO.Println("issue a new subscription ", sid)

					tb.ldSubscriptions_lock.Lock()
					tb.ldSubscriptions[sid] = &newLDSubscription
					tb.ldSubscriptions_lock.Unlock()

					tb.subLinks_lock.Lock()
					tb.main2Other[mainSubID] = append(tb.main2Other[mainSubID], sid)
					tb.subLinks_lock.Unlock()
				}
			}
		}
	}
}
func (tb *ThinBroker) handleNGSI9Notify(mainSubID string, notifyContextAvailabilityReq *NotifyContextAvailabilityRequest) {
	var action string
	notifyContextAvailabilityReq.ErrorCode.Code = 301
	switch notifyContextAvailabilityReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	}
	INFO.Println(action, " subID ", mainSubID)
	for _, registrationResp := range notifyContextAvailabilityReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration
		for _, eid := range registration.EntityIdList {
			INFO.Println("===> ", eid, " , ", mainSubID)
			tb.e2sub_lock.Lock()

			if action == "CREATE" {
				tb.entityId2Subcriptions[eid.ID] = append(tb.entityId2Subcriptions[eid.ID], mainSubID)
			} else if action == "DELETE" {
				subList := tb.entityId2Subcriptions[eid.ID]
				for i, id := range subList {
					if id == mainSubID {
						tb.entityId2Subcriptions[eid.ID] = append(subList[:i], subList[i+1:]...)
						break
					}
				}
			} else if action == "UPDATE" {
				existFlag := false
				for _, subID := range tb.entityId2Subcriptions[eid.ID] {
					if subID == mainSubID {
						existFlag = true
						break
					}
				}
				if existFlag == false {
					tb.entityId2Subcriptions[eid.ID] = append(tb.entityId2Subcriptions[eid.ID], mainSubID)
				}
			}

			tb.e2sub_lock.Unlock()
		}

		INFO.Println(registration.ProvidingApplication, ", ", tb.MyURL)
		INFO.Println("TO ngsi10 subscription, ", mainSubID)
		INFO.Printf("entity list: %+v\r\n", registration.EntityIdList)

		if registration.ProvidingApplication == tb.MyURL {
			//for matched entities provided by myself
			if action == "CREATE" || action == "UPDATE" {
				tb.notifyOneSubscriberWithCurrentStatus(registration.EntityIdList, mainSubID)
			}
		} else {
			//for matched entities provided by other IoT Brokers
			newSubscription := SubscribeContextRequest{}
			newSubscription.Entities = registration.EntityIdList
			newSubscription.Reference = tb.MyURL
			newSubscription.Subscriber.BrokerURL = registration.ProvidingApplication

			if action == "CREATE" || action == "UPDATE" {
				sid, err := subscribeContextProvider(&newSubscription, registration.ProvidingApplication, tb.SecurityCfg)
				if err == nil {
					INFO.Println("issue a new subscription ", sid)

					tb.subscriptions_lock.Lock()
					tb.subscriptions[sid] = &newSubscription
					tb.subscriptions_lock.Unlock()

					tb.subLinks_lock.Lock()
					tb.main2Other[mainSubID] = append(tb.main2Other[mainSubID], sid)
					tb.subLinks_lock.Unlock()
				}
			}
		}
	}
}

/*
	Send the subscription based notification
*/

func (tb *ThinBroker) handleNGSIV2Notify(mainSubID string, notifyv2ContextAvailabilityReq *Notifyv2ContextAvailabilityRequest) {
	var action string
	switch notifyv2ContextAvailabilityReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	}

	INFO.Println(action, " subID ", mainSubID)

	for _, registrationResp := range notifyv2ContextAvailabilityReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration
		for _, eid := range registration.EntityIdList {
			INFO.Println("===> ", eid, " , ", mainSubID)

			tb.ev2sub_lock.Lock()
			if action == "CREATE" {
				tb.entityIdv2Subcriptions[eid.ID] = append(tb.entityIdv2Subcriptions[eid.ID], mainSubID)
			} else if action == "DELETE" {
				subList := tb.entityIdv2Subcriptions[eid.ID]
				for i, id := range subList {
					if id == mainSubID {
						tb.entityIdv2Subcriptions[eid.ID] = append(subList[:i], subList[i+1:]...)
						break
					}
				}
			} else if action == "UPDATE" {
				existFlag := false
				for _, subID := range tb.entityIdv2Subcriptions[eid.ID] {
					if subID == mainSubID {
						existFlag = true
						break
					}
				}
				if existFlag == false {
					tb.entityIdv2Subcriptions[eid.ID] = append(tb.entityIdv2Subcriptions[eid.ID], mainSubID)
				}
			}
			tb.ev2sub_lock.Unlock()
		}

		INFO.Println(registration.ProvidingApplication, ", ", tb.MyURL)
		INFO.Println("TO ngsi10 subscription, ", mainSubID)
		INFO.Printf("entity list: %+v\r\n", registration.EntityIdList)

		if registration.ProvidingApplication == tb.MyURL {
			//for matched entities provided by myself
			if action == "CREATE" || action == "UPDATE" {
				tb.notifyOneSubscriberv2WithCurrentStatus(registration.EntityIdList, mainSubID)
			}
		} else {
			//for matched entities provided by other IoT Brokers
			newv2Subscription := SubscriptionRequest{}
			newv2Subscription.Subject.Entities = registration.EntityIdList
			newv2Subscription.Notification.Http.Url = tb.MyURL
			newv2Subscription.Subscriber.BrokerURL = registration.ProvidingApplication

			if action == "CREATE" || action == "UPDATE" {
				registration.ProvidingApplication = strings.TrimSuffix(registration.ProvidingApplication, "/ngsi10")
				sid, err := subscriptionProvider(&newv2Subscription, registration.ProvidingApplication, tb.SecurityCfg)
				if err == nil {
					INFO.Println("issue a new subscription ", sid)

					tb.v2subscriptions_lock.Lock()
					tb.v2subscriptions[sid] = &newv2Subscription
					tb.v2subscriptions_lock.Unlock()

					tb.subLinks_lock.Lock()
					tb.main2Other[mainSubID] = append(tb.main2Other[mainSubID], sid)
					tb.subLinks_lock.Unlock()
				}
			}
		}
	}
}

func (tb *ThinBroker) registerContextElement(element *ContextElement) {
	registration := ContextRegistration{}

	entities := make([]EntityId, 0)
	entities = append(entities, element.Entity)
	registration.EntityIdList = entities

	attributes := make([]ContextRegistrationAttribute, 0)
	for _, item := range element.Attributes {
		attr := ContextRegistrationAttribute{}
		attr.Name = item.Name
		attr.Type = item.Type
		attr.IsDomain = false
		attributes = append(attributes, attr)
	}
	registration.ContextRegistrationAttributes = attributes
	registration.Metadata = element.Metadata
	registration.ProvidingApplication = tb.MyURL
	registration.MsgFormat = "NGSIV1"

	// create or update registered context
	registerCtxReq := RegisterContextRequest{}
	registerCtxReq.RegistrationId = ""
	registerCtxReq.ContextRegistrations = []ContextRegistration{registration}
	registerCtxReq.Duration = "PT10M"
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	_, err := client.RegisterContext(&registerCtxReq)
	if err != nil {
		ERROR.Println(err)
	}
}

func (tb *ThinBroker) deregisterContextElements(ContextElements []ContextElement) {
	registrationList := make([]ContextRegistration, 0)

	for _, element := range ContextElements {
		registration := ContextRegistration{}

		entities := make([]EntityId, 0)
		entities = append(entities, element.Entity)
		registration.EntityIdList = entities

		registration.ProvidingApplication = tb.MyURL

		registrationList = append(registrationList, registration)
	}

	// issue a contextRegistration to remove their availability information based on entity id
	registerCtxReq := RegisterContextRequest{}
	registerCtxReq.RegistrationId = ""
	registerCtxReq.ContextRegistrations = registrationList
	registerCtxReq.Duration = "0"

	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	_, err := client.RegisterContext(&registerCtxReq)
	if err != nil {
		ERROR.Println(err)
	}
}

//Southbound feature addition- Device Registration starts here.
func (tb *ThinBroker) RegisterContext(w rest.ResponseWriter, r *rest.Request) {

	// IoTAgent sends isPattern in RegisterContext Request as a string while Fogflow accepts it as bool
	RegCtxReq, err := tb.handleIoTRegisterContext(r)

	if err != nil {
		DEBUG.Println("Not able to decode the registration!")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RegisterCtxReq := *RegCtxReq

	FiwareService := r.Header.Get("Fiware-Service")
	FiwareServicePath := r.Header.Get("Fiware-ServicePath")
	if FiwareService == "" || FiwareServicePath == "" {
		rest.Error(w, "Bad Request! Fiware-Service and/or Fiware-ServicePath Headers are Missing!", 400)
		return
	} else {
		for _, registration := range RegisterCtxReq.ContextRegistrations {
			reg := ContextRegistration{}

			fiData := FiwareData{}
			fiData.ProviderIoTAgent = registration.ProvidingApplication
			fiData.FiwareService = FiwareService
			fiData.FiwareServicePath = FiwareServicePath

			//creating separate registration for each entity
			for _, entity := range registration.EntityIdList {
				reg.EntityIdList = nil
				RegID := tb.createIdWithFiwareHeaders(entity.ID, FiwareService, FiwareServicePath)
				errString := tb.createFiwareData(RegID, fiData)
				if errString != "" {
					rest.Error(w, errString, 409)
					continue
				} else {
					entity.ID = RegID
					reg.EntityIdList = append(reg.EntityIdList, entity)
				}
				//Creating registration request for discovery
				registration.EntityIdList = reg.EntityIdList

				registration.ProvidingApplication = tb.MyURL

				RegCtxReq := RegisterContextRequest{}
				RegCtxReq.ContextRegistrations = append(RegCtxReq.ContextRegistrations, registration)

				RegCtxReq.Duration = RegisterCtxReq.Duration

				DEBUG.Println("Sending following registration to Discovery:")
				DEBUG.Println(RegCtxReq)

				client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
				RegisterCtxResp, err := client.RegisterContext(&RegCtxReq)
				// send out the response
				if err != nil {
					w.WriteJson(err)
				} else {
					w.WriteJson(&RegisterCtxResp)
				}
			}
		}
	}
}

func (tb *ThinBroker) handleIoTRegisterContext(r *rest.Request) (*RegisterContextRequest, error) {
	RegisterCtxReq := RegisterContextRequest{}
	//Decode IoT Agent Payload having isPattern of type string
	RegisterCtxReq1 := RegisterContextRequest1{}
	err := r.DecodeJsonPayload(&RegisterCtxReq1)
	DEBUG.Println("JSON payload decoded....")
	if err != nil {
		return nil, err
	}

	for _, registration := range RegisterCtxReq1.ContextRegistrations {
		contextRegistration := ContextRegistration{}
		for _, fiwareEntity := range registration.EntityIdList {

			fogflowEntity := EntityId{}

			fogflowEntity.ID = fiwareEntity.ID
			fogflowEntity.Type = fiwareEntity.Type
			if fiwareEntity.IsPattern == "true" || fiwareEntity.IsPattern == "True" {
				fogflowEntity.IsPattern = true
			} else if fiwareEntity.IsPattern == "false" || fiwareEntity.IsPattern == "False" {
				fogflowEntity.IsPattern = false
			}

			contextRegistration.EntityIdList = append(contextRegistration.EntityIdList, fogflowEntity)
		}
		contextRegistration.ContextRegistrationAttributes = registration.ContextRegistrationAttributes
		contextRegistration.Metadata = registration.Metadata
		contextRegistration.ProvidingApplication = registration.ProvidingApplication
		RegisterCtxReq.ContextRegistrations = append(RegisterCtxReq.ContextRegistrations, contextRegistration)
	}

	RegisterCtxReq.Duration = RegisterCtxReq1.Duration
	return &RegisterCtxReq, nil

}

func (tb *ThinBroker) createFiwareData(RegID string, fiData FiwareData) string {
	errString := ""
	if tb.getRegistration(RegID) != nil {
		errString = "Registration already exists for this Entity ID!"
	} else {
		//Storing FiwareData
		tb.fiwareData_lock.Lock()
		tb.fiwareData[RegID] = &fiData
		tb.fiwareData_lock.Unlock()
	}
	return errString
}

func (tb *ThinBroker) getRegistration(eid string) *EntityRegistration {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	_, registration := client.GetProviderURL(eid)
	if registration.ID == "" {
		DEBUG.Println("Registration not found!")
		return nil
	} else {
		return registration
	}
}

func (tb *ThinBroker) deleteRegistration(rid string) error {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	client.UnregisterEntity(rid)
	return nil
}

func (tb *ThinBroker) updateIdWithFiwareHeaders(ctxElem *ContextElement, fiwareService string, fiwareServicePath string) {
	ctxElem.Entity.ID = tb.createIdWithFiwareHeaders(ctxElem.Entity.ID, fiwareService, fiwareServicePath)
}

func (tb *ThinBroker) createIdWithFiwareHeaders(eid string, fiwareService string, fiwareServicePath string) string {
	eid = eid + "." + fiwareService + "." + fiwareServicePath
	eid = strings.ReplaceAll(eid, "/", "~")
	return eid
}

func (tb *ThinBroker) removeFiwareHeadersFromId(ctxElem *ContextElement, fiwareService string, fiwareServicePath string) {
	cutStr := "." + fiwareService + "." + fiwareServicePath
	cutStr = strings.ReplaceAll(cutStr, "/", "~")
	ctxElem.Entity.ID = strings.TrimRight(ctxElem.Entity.ID, cutStr)
}

// NGSI-LD starts from here.

func (tb *ThinBroker) LDUpdateContext(w rest.ResponseWriter, r *rest.Request) {
	err := contentTypeValidator(r.Header.Get("Content-Type"))
	var fiwareService, fsp string
	//var fiwareServicePath string
	if r.Header.Get("fiware-service") != "" {
		fiwareService = r.Header.Get("fiware-service")
	} else {
		fiwareService = "default"
	}
	if r.Header.Get("fiware-servicePath") != "" {
		fsp = r.Header.Get("fiware-servicePath")
	} else {
		fsp = "default"
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		reqBytes, _ := ioutil.ReadAll(r.Body)
		var LDupdateCtxReq []interface{}

		err = json.Unmarshal(reqBytes, &LDupdateCtxReq)

		if err != nil {
			err := errors.New("This interface only supports arrays of entities")
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := ResponseError{}
		fmt.Println("LDupdateCtxReq",LDupdateCtxReq)
		for _, ctx := range LDupdateCtxReq {
			var context []interface{}
			contextInPayload := false
			cType := r.Header.Get("Content-Type")
			cTypeInLower := strings.ToLower(cType)
			Link := r.Header.Get("Link")
			if cTypeInLower == "application/ld+json" {
				contextInPayload = true
			} else {
				if link := r.Header.Get("Link"); link != "" {
					link := extractLinkHeaderFields(link)
					if link == "default" {
						context = append(context, DEFAULT_CONTEXT)
					} else {
						context = append(context, link)
					}
				} else {
					context = append(context, DEFAULT_CONTEXT)
				}
			}
			ctxEle := ctx.(map[string]interface{})
			_, ok1 := ctxEle["id"]
			_, ok2 := ctxEle["@id"]
			var Eid, brokerURL string
			if ok1 == true {
				Eid = ctxEle["id"].(string)
			} else if ok2 == true {
				Eid = ctxEle["@id"].(string)
			}
			if Eid == "" {
				brokerURL = tb.myProfile.MyURL
			} else {
				eid := Eid + "@" + fiwareService
				brokerURL = tb.queryOwnerOfLDEntity(eid)
			}
			if brokerURL != tb.myProfile.MyURL {
				updateCtxReq := ctx.(map[string]interface{})
				updateCtxReq["id"] = updateCtxReq["id"].(string) + "@" + fiwareService
				updateCtxReq["fiwareServicePath"] = fsp
				tb.UpdateLdContext2RemoteSite(updateCtxReq, brokerURL, Link)
			} else {
				resolved, err := tb.ExpandPayload(ctx, context, contextInPayload)
				if err != nil {
					problemSet := ProblemDetails{}
					problemSet.Details = err.Error()
					res.Errors = append(res.Errors, problemSet)
					continue
				} else {
					sz := Serializer{}

					// Deserialize the payload here.
					deSerializedEntity, err := sz.DeSerializeEntity(resolved)
					if err != nil {
						problemSet := ProblemDetails{}
						problemSet.Details = err.Error()
						res.Errors = append(res.Errors, problemSet)
						continue
					} else {
						//Update createdAt value.
						if !strings.HasPrefix(deSerializedEntity[NGSI_LD_ID].(string), "urn:ngsi-ld:") {
							problemSet := ProblemDetails{}
							problemSet.Details = "Entity id must contain uri!"
							res.Errors = append(res.Errors, problemSet)
							continue
						}
						deSerializedEntity[NGSI_LD_ID] = getIoTID(deSerializedEntity[NGSI_LD_ID].(string), fiwareService)
						//deSerializedEntity["createdAt"] = time.Now().String()
						// Store Context
						deSerializedEntity["@context"] = context
						deSerializedEntity["fiwareServicePath"] = fsp
						res.Success = append(res.Success, deSerializedEntity[NGSI_LD_ID].(string))
						tb.UpdateLdContext2LocalSite(deSerializedEntity)
					}
				}
			}
		}
		if res.Errors != nil && res.Success == nil {
			w.WriteHeader(404)
			w.WriteJson(&res)
		}
		if res.Errors != nil && res.Success != nil {
			w.WriteHeader(207)
			w.WriteJson(&res)
		}
		if res.Errors == nil && res.Success != nil {
			w.WriteHeader(http.StatusNoContent)
		}

	}
}

// Create an NGSI-LD Entity

func (tb *ThinBroker) LDCreateEntity(w rest.ResponseWriter, r *rest.Request) {
	//Also allow the header to json+ld for specific cases
	if ctype, accept := r.Header.Get("Content-Type"), r.Header.Get("Accept"); (ctype == "application/json" || ctype == "application/ld+json") && accept == "application/ld+json" {
		var fiwareService, fsp string
		//var fiwareServicePath string
		if r.Header.Get("fiware-service") != "" {
			fiwareService = r.Header.Get("fiware-service")
		} else {
			fiwareService = "default"
		}

		if r.Header.Get("fiware-servicePath") != "" {
			fsp = r.Header.Get("fiware-servicePath")
		} else {
			//deSerializedEntity["fiwareServicePath"] = "default"
			fsp = "default"
		}

		var context []interface{}
		contextInPayload := false
		cType := r.Header.Get("Content-Type")
		cTypeInLower := strings.ToLower(cType)
		Link := r.Header.Get("Link")
		if cTypeInLower == "application/ld+json" {
			contextInPayload = true
		} else {
			context = append(context, DEFAULT_CONTEXT)
		}

		//context = append(context, DEFAULT_CONTEXT)
		reqBytes, _ := ioutil.ReadAll(r.Body)
		var LDupdateCtxReq interface{}

		err := json.Unmarshal(reqBytes, &LDupdateCtxReq)

		if err != nil {
			err := errors.New("Unable to decode payload/message !")
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Get a resolved object ([]interface object)
		ctxEle := LDupdateCtxReq.(map[string]interface{})
		_, ok1 := ctxEle["id"]
		_, ok2 := ctxEle["@id"]
		var Eid, brokerURL string
		if ok1 == true {
			Eid = ctxEle["id"].(string)
		} else if ok2 == true {
			Eid = ctxEle["@id"].(string)
		}
		if Eid == "" {
			brokerURL = tb.myProfile.MyURL
		} else {
			eid := Eid + fiwareService
			brokerURL = tb.queryOwnerOfLDEntity(eid)
		}
		if brokerURL != tb.myProfile.MyURL {
			w.WriteHeader(201)
			updateCtxReq := ctxEle
			updateCtxReq["id"] = updateCtxReq["id"].(string) + "@" + fiwareService
			updateCtxReq["fiwareServicePath"] = fsp
			tb.UpdateLdContext2RemoteSite(updateCtxReq, brokerURL, Link)
		} else {
			resolved, err := tb.ExpandPayload(LDupdateCtxReq, context, contextInPayload)
			if err != nil {

				if err.Error() == "EmptyPayload!" {
					rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
					return
				}
				/*if err.Error() == "AlreadyExists!" {
					rest.Error(w, "AlreadyExists!", 409)
					return
				}*/
				if err.Error() == "Id can not be nil!" {
					rest.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err.Error() == "Type can not be nil!" {
					rest.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				sz := Serializer{}

				// Deserialize the payload here.
				deSerializedEntity, err := sz.DeSerializeEntity(resolved)
				if err != nil {
					rest.Error(w, err.Error(), 400)
					return
				} else {
					//Update createdAt value.
					//deSerializedEntity["createdAt"] = time.Now().String()
					// Store Context
					deSerializedEntity["@context"] = context

					if !strings.HasPrefix(deSerializedEntity[NGSI_LD_ID].(string), "urn:ngsi-ld:") {
						rest.Error(w, "Entity id must contain uri!", 400)
						return
					}
					deSerializedEntity[NGSI_LD_ID] = getIoTID(deSerializedEntity[NGSI_LD_ID].(string), fiwareService)
					deSerializedEntity["@context"] = context
					w.Header().Set("Location", "/ngis-ld/v1/entities/"+deSerializedEntity[NGSI_LD_ID].(string))
					w.WriteHeader(201)
					deSerializedEntity["fiwareServicePath"] = fsp
					tb.UpdateLdContext2LocalSite(deSerializedEntity)
				}
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

func (tb *ThinBroker) updateCtxElemet(elem map[string]interface{}, eid string) error {
	entity := tb.ldEntities[eid]
	entityMap := entity.(map[string]interface{})
	for k, v := range elem {
		isLdJsonOject := checkCondition(k)
		if isLdJsonOject == true && k != "mgsFormat" && k != "fiwareServicePath" {
			if _, ok := entityMap[k]; ok == true {
				updatedResult := make(map[string]interface{})
				prevEleAttr := entityMap[k].(map[string]interface{})
				currEleAttr := elem[k].(map[string]interface{})
				updatedResult = update(prevEleAttr, currEleAttr)
				updatedResult["modifiedAt"] = time.Now().String()
				entityMap[k] = updatedResult
			} else {
				if isLdJsonOject == true {
					entityMap[k] = v
				}
			}
		}
	}
	//entityMap["modifiedAt"] = time.Now().String()
	//tb.ldEntities[eid] = entityMap
	return nil
}

func (tb *ThinBroker) updateLdContextElement(ctxEle map[string]interface{}) {
	tb.ldEntities_lock.Lock()
	defer tb.ldEntities_lock.Unlock()
	eid := getId(ctxEle)
	if _, exist := tb.ldEntities[eid]; exist {
		tb.updateCtxElemet(ctxEle, eid)
	} else {
		typ := getRegistrationType(ctxEle["@type"])
		tb.entityTypeTOEntityId[typ] = append(tb.entityTypeTOEntityId[typ], eid)
		tb.ldEntities[eid] = ctxEle
	}
}
func (tb *ThinBroker) UpdateLdContext2LocalSite(updateCtxReq map[string]interface{}) {
	tb.ldEntities_lock.Lock()
	eid := getId(updateCtxReq)
	hasLdUpdatedMetadata := hasLdUpdatedMetadata(updateCtxReq, tb.ldEntities[eid])
	tb.ldEntities_lock.Unlock()

	tb.updateLdContextElement(updateCtxReq)

	go tb.LDNotifySubscribers(updateCtxReq, true)
	if hasLdUpdatedMetadata == true {
		tb.registerLDContextElement(updateCtxReq)
	}
}

func (tb *ThinBroker) UpdateLdContext2RemoteSite(updateCtxReq map[string]interface{}, brokerURL string, link string) {
	INFO.Println(brokerURL)
	//client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: tb.SecurityCfg}
	LDBrokerURL := strings.TrimSuffix(brokerURL, "/ngsi10")
	if link != "" {
		delete(updateCtxReq, "@context")
	}
	client := NGSI10Client{IoTBrokerURL: LDBrokerURL, SecurityCfg: tb.SecurityCfg}
	client.CreateLDEntityOnRemote(updateCtxReq, link)
}

func (tb *ThinBroker) handleLdExternalUpdateContext(updateCtxReq map[string]interface{}, link string) {
	eid := getId(updateCtxReq)
	brokerURL := tb.queryOwnerOfLDEntity(eid)
	if brokerURL == tb.myProfile.MyURL {
		tb.UpdateLdContext2LocalSite(updateCtxReq)
	} else {
		tb.UpdateLdContext2RemoteSite(updateCtxReq, brokerURL, link)
	}
}

func (tb *ThinBroker) updateLDspecificAttributeValues2RemoteSite(req map[string]interface{}, remoteURL string, eid string, attr string) (error, int) {
	client := NGSI10Client{IoTBrokerURL: remoteURL, SecurityCfg: tb.SecurityCfg}
	err, code := client.UpdateLDEntityspecificAttributeOnRemote(req, eid, attr)

	if err != nil {
		return err, code
	}
	return nil, code
}

func (tb *ThinBroker) updateLDAttributeValues2RemoteSite(req map[string]interface{}, remoteURL string, eid string) (error, int) {
	client := NGSI10Client{IoTBrokerURL: remoteURL, SecurityCfg: tb.SecurityCfg}
	err, code := client.UpdateLDEntityAttributeOnRemote(req, eid)

	if err != nil {
		return err, code
	}
	return nil, code
}

func (tb *ThinBroker) updateLDAttribute2RemoteSite(req map[string]interface{}, remoteURL string, eid string) error {
	client := NGSI10Client{IoTBrokerURL: remoteURL, SecurityCfg: tb.SecurityCfg}
	err := client.AppendLDEntityOnRemote(req, eid)

	if err != nil {
		return err
	}
	return nil
}

func (tb *ThinBroker) updateLDContextElement2RemoteSite(req map[string]interface{}, remoteURL string, link string) error {
	client := NGSI10Client{IoTBrokerURL: remoteURL, SecurityCfg: tb.SecurityCfg}
	err := client.CreateLDEntityOnRemote(req, link)

	if err != nil {
		return err
	}
	return nil
}

// Register a new context entity on Discovery
func (tb *ThinBroker) registerLDContextElement(elem map[string]interface{}) {
	fmt.Println("elem for register", elem)
	registerCtxReq := RegisterContextRequest{}
	entities := make([]EntityId, 0)
	entityId := EntityId{}
	entityId.ID = elem["@id"].(string)
	_, fs := FiwareId(elem["@id"].(string))
	//}
	//fmt.Println("Fs", Fs)
	entityId.Type = getRegistrationType(elem["@type"])
	//entityId.Type = elem["type"].(string)
	entities = append(entities, entityId)

	ctxRegistrations := make([]ContextRegistration, 0)

	ctxReg := ContextRegistration{}
	ctxReg.EntityIdList = entities
	ctxRegAttr := ContextRegistrationAttribute{}
	ctxRegAttrs := make([]ContextRegistrationAttribute, 0)
	ctxMetadatas := make([]ContextMetadata,0)
	for k, attr := range elem { // considering properties and relationships as attributes
		isLdJsonObject := checkCondition(k)
		ctxMetaData := ContextMetadata{}
		if isLdJsonObject == true {
			if k == "fiwareServicePath" {
				ctxReg.FiwareServicePath = attr.(string)
				continue
			}
			attrValue := attr.(map[string]interface{})
			ctxRegAttr.Name = k
			typ := getRegistrationType(attrValue["@type"])
			if strings.Contains(typ,"GeoProperty") || strings.Contains(typ, "geoProperty") {
				if strings.Contains(k,"location") {
					ctxMetaData.Name = "location"
				} else {
					ctxMetaData.Name = k
				}
				attrValue["@context"] = elem["@context"]
				resolved, err := compactData(attrValue, DEFAULT_CONTEXT)
				if err != nil {
					continue
				}
				resolvedMap := resolved.(map[string]interface{})
				value := resolvedMap["value"].(map[string]interface{})
				ctxMetaData.Type = value["type"].(string)
				ctxMetaData.Cordinates = value["coordinates"]
				ctxMetadatas =append(ctxMetadatas, ctxMetaData)

			} else if strings.Contains(typ, "Property") || strings.Contains(typ, "property") {
				ctxRegAttr.Type = "Property"
			} else if strings.Contains(typ, "Relationship") || strings.Contains(typ, "relationship") {
				ctxRegAttr.Type = "Relationship"
			}
			ctxRegAttrs = append(ctxRegAttrs, ctxRegAttr)
		}
	}
	ctxReg.ContextRegistrationAttributes = ctxRegAttrs
	ctxReg.Metadata  = ctxMetadatas
	ctxReg.ProvidingApplication = tb.MyURL
	ctxReg.MsgFormat = "NGSILD"
	ctxReg.FiwareService = fs
	ctxRegistrations = append(ctxRegistrations, ctxReg)

	registerCtxReq.ContextRegistrations = ctxRegistrations

	// Send the registration to discovery
	fmt.Println("&registerCtxReq",&registerCtxReq)
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	_, err := client.RegisterContext(&registerCtxReq)
	if err != nil {
		ERROR.Println(err)
	}
}

// Store the NGSI-LD Entities  at local broker
func (tb *ThinBroker) saveEntity(ctxElem map[string]interface{}) {
	eid := ctxElem["id"].(string)
	tb.ldEntities_lock.Lock()
	tb.ldEntities[eid] = ctxElem
	tb.ldEntities_lock.Unlock()
}

// GET API method for entity
func (tb *ThinBroker) ldGetEntity(eid string) interface{} {
	tb.ldEntities_lock.RLock()
	if entity := tb.ldEntities[eid]; entity != nil {
		tb.ldEntities_lock.RUnlock()
		entityMap := entity.(map[string]interface{})
		compactEntity := tb.createOriginalPayload(entityMap)
		resultEntity := compactEntity.(map[string]interface{})
		actualId := getActualEntity(resultEntity)
		resultEntity["@id"] = actualId
		delete(resultEntity, "fiwareServicePath")
		return resultEntity
	} else {
		tb.ldEntities_lock.RUnlock()
		return nil
	}
}

// Creating original payload as provided by user from FogFlow Data Structure
func (tb *ThinBroker) createOriginalPayload(entity interface{}) interface{} {
	entityMap := entity.(map[string]interface{})

	// Expanding the entity to get uniformly expanded entity which was missing in internal representation
	expandedEntity, err := tb.ExpandData(entityMap)

	if err != nil {
		DEBUG.Println("Error while expanding:", err)
		return nil
	}

	// Compacting the expanded entity.
	entity1 := expandedEntity[0].(map[string]interface{})
	compactEntity, err := compactData(entity1, entityMap["@context"])
	if err != nil {
		DEBUG.Println("Error while compacting:", err)
		return nil
	}
	return compactEntity
}

// Compacting data to display to user in original form.
/*func (tb *ThinBroker) compactData(entity map[string]interface{}, context interface{}) (interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	compacted, err := proc.Compact(entity, context, options)
	return compacted, err
}*/

func (tb *ThinBroker) LDCreateSubscription(w rest.ResponseWriter, r *rest.Request) {
	//Also allow the header to json+ld for specific cases
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		var fiwareService string
		var fiwareServicePath string
		if r.Header.Get("fiware-service") != "" {
			fiwareService = r.Header.Get("fiware-service")
		} else {
			fiwareService = "default"
		}
		if r.Header.Get("fiware-servicePath") != "" {
			fiwareServicePath = r.Header.Get("fiware-servicePath")
		} else {
			fiwareServicePath = "default"
		}
		fmt.Println(fiwareService, fiwareServicePath)
		reqBytes, _ := ioutil.ReadAll(r.Body)
		var LDSubscribeCtxReq interface{}

		err := json.Unmarshal(reqBytes, &LDSubscribeCtxReq)
		fmt.Println("err", err)
		if err != nil {
			err := errors.New("Unable to decode payload/message !")
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var context []interface{}
		contextInPayload := false
		cType := r.Header.Get("Content-Type")
		cTypeInLower := strings.ToLower(cType)
		if cTypeInLower == "application/ld+json" {
			contextInPayload = true
		} else {
			if link := r.Header.Get("Link"); link != "" {
				link := extractLinkHeaderFields(link)
				if link == "default" {
					context = append(context, DEFAULT_CONTEXT)
				} else {
					context = append(context, link)
				}
			} else {
				context = append(context, DEFAULT_CONTEXT)
			}
		}

		resolved, err := tb.ExpandPayload(LDSubscribeCtxReq, context, contextInPayload)

		fmt.Println("err", err)
		if err != nil {
			if err.Error() == "EmptyPayload!" {
				rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
				return
			}
			if err.Error() == "Type can not be nil!" {
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			sz := Serializer{}
			deSerializedSubscription, err := sz.DeSerializeSubscription(resolved)

			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				deSerializedSubscription.CreatedAt = time.Now().String()
				// Create Subscription Id, if missing

				if deSerializedSubscription.Id == "" {
					u1, err := uuid.NewUUID()
					if err != nil {
						rest.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					sid := u1.String()
					deSerializedSubscription.Id = sid
				}
				if !strings.HasSuffix(deSerializedSubscription.Type, "Subscription") && !strings.HasSuffix(deSerializedSubscription.Type, "subscription") {
					rest.Error(w, "Type not allowed!", http.StatusBadRequest)
					return
				}
				if len(deSerializedSubscription.Entities) == 0 {
					rest.Error(w, "Missing entites and its parameter!", http.StatusBadRequest)
					return
				}
				// send response
				w.WriteHeader(http.StatusCreated)
				subResp := SubscribeContextResponse{}
				subResp.SubscribeResponse.SubscriptionId = deSerializedSubscription.Id
				subResp.SubscribeError.SubscriptionId = deSerializedSubscription.Id
				w.WriteJson(&subResp)
				//deSerializedSubscription.Id = getIoTID(deSerializedSubscription.Id,fiwareService)
				if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
					deSerializedSubscription.Subscriber.IsInternal = true
				} else {
					deSerializedSubscription.Subscriber.IsInternal = false
				}

				if r.Header.Get("Require-Reliability") == "true" {
					deSerializedSubscription.Subscriber.RequireReliability = true
					deSerializedSubscription.Subscriber.LDNotifyCache = make([]map[string]interface{}, 0)
				} else {
					deSerializedSubscription.Subscriber.RequireReliability = false
				}

				deSerializedSubscription.Status = "active"                  // others allowed: paused, expired
				deSerializedSubscription.Notification.Format = "normalized" // other allowed: keyValues
				deSerializedSubscription.Subscriber.BrokerURL = tb.MyURL
				//tb.createEntityID2SubscriptionsIDMap(&deSerializedSubscription)

				// For integration with NGSILD broker
				if r.Header.Get("Integration") != "" {
					deSerializedSubscription.Subscriber.Integration = r.Header.Get("Integration")
				}
				// for integration With IoT Agent
				if r.Header.Get("fiware-service") != "" {
					deSerializedSubscription.Subscriber.FiwareService = r.Header.Get("fiware-service")
				} else {
					deSerializedSubscription.Subscriber.FiwareService = ""
				}
				if r.Header.Get("fiware-servicepath") != "" {
					deSerializedSubscription.Subscriber.FiwareServicePath = r.Header.Get("fiware-servicepath")
				} else {
					deSerializedSubscription.Subscriber.FiwareServicePath = ""
				}

				// save subscription
				for key, entity := range deSerializedSubscription.Entities {
					if entity.ID != "" && strings.Contains(entity.ID, "@") == false {
						entity.ID = entity.ID + "@" + fiwareService
					}
					if entity.IdPattern == "" && entity.ID == "" {
						entity.IsPattern = true
					}
					deSerializedSubscription.Entities[key] = entity
				}
				tb.createSubscription(&deSerializedSubscription)
				if deSerializedSubscription.Subscriber.IsInternal == true {
					INFO.Println("internal subscription coming from another broker")
					for _, entity := range deSerializedSubscription.Entities {
						tb.LDe2sub_lock.Lock()
						tb.entityId2LDSubcriptions[entity.ID] = append(tb.entityId2LDSubcriptions[entity.ID], deSerializedSubscription.Id)
						tb.LDe2sub_lock.Unlock()
					}
					tb.notifyOneLDSubscriberWithCurrentStatus(deSerializedSubscription.Entities, deSerializedSubscription.Id)
				} else {
					tb.SubscribeLDContextAvailability(&deSerializedSubscription, fiwareService)
				}
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

// Subscribe to Discovery for context availabiltiy
func (tb *ThinBroker) SubscribeLDContextAvailability(subReq *LDSubscriptionRequest, fiwareService string) error {
	ctxAvailabilityRequest := SubscribeContextAvailabilityRequest{}
	for key, entity := range subReq.Entities {
		if entity.IdPattern != "" {
			entity.IsPattern = true
		}
		subReq.Entities[key] = entity
	}
	ctxAvailabilityRequest.FiwareService = fiwareService
	ctxAvailabilityRequest.Entities = subReq.Entities
	ctxAvailabilityRequest.Attributes = subReq.WatchedAttributes
	ctxAvailabilityRequest.Attributes = append(ctxAvailabilityRequest.Attributes, subReq.Notification.Attributes...)
	ctxAvailabilityRequest.Reference = tb.MyURL + "/notifyLDContextAvailability"
	ctxAvailabilityRequest.Duration = subReq.Expires
	ctxAvailabilityRequest.Restriction = subReq.Restriction
	// Subscribe to discovery
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	AvailabilitySubID, err := client.SubscribeContextAvailability(&ctxAvailabilityRequest)

	if AvailabilitySubID != "" {
		//tb.createSubscriptionIdMappings(subReq.Id, AvailabilitySubID)
		tb.subLinks_lock.Lock()
		subID := subReq.Id
		tb.main2Other[subID] = append(tb.main2Other[subID], AvailabilitySubID)
		tb.availabilitySub2MainSub[AvailabilitySubID] = subID
		notifyMessage, alreadyBack := tb.tmpNGSILDNotifyCache[AvailabilitySubID]
		tb.subLinks_lock.Unlock()
		if alreadyBack == true {
			INFO.Println("========forward the availability notify that arrived earlier===========")
			tb.handleNGSI9Notify(subReq.Id, notifyMessage)

			tb.subLinks_lock.Lock()
			delete(tb.tmpNGSILDNotifyCache, AvailabilitySubID)
			tb.subLinks_lock.Unlock()
		}
		return nil
	} else {
		INFO.Println("failed to subscribe the availability of requested entities ", err)
		return err
	}
}

// Store in SubID - SubscriptionPayload Map
func (tb *ThinBroker) createSubscription(subscription *LDSubscriptionRequest) {
	subscription.Subscriber.LDNotifyCache = make([]map[string]interface{}, 0)
	tb.ldSubscriptions_lock.Lock()
	subscription.SetLdIdPattern()
	tb.ldSubscriptions[subscription.Id] = subscription
	tb.ldSubscriptions_lock.Unlock()
}

// Expand the payload
func (tb *ThinBroker) ExpandPayload(ctx interface{}, context []interface{}, contextInPayload bool) ([]interface{}, error) {
	//get map[string]interface{} of reqBody
	itemsMap, err := tb.getStringInterfaceMap(ctx)
	if err != nil {
		return nil, err
	} else {
		// Check the type of payload: Entity, registration or Subscription
		var payloadType string
		if _, ok := itemsMap["type"]; ok == true {
			payloadType = itemsMap["type"].(string)
		} else if _, ok := itemsMap["@type"]; ok == true {
			if reflect.ValueOf(itemsMap["@type"]).Kind() == reflect.String {
				payloadType = itemsMap["@type"].(string)
			} else {
				typ := itemsMap["@type"].([]interface{})
				payloadType = typ[0].(string)
			}
		}
		if payloadType == "" {
			err := errors.New("Type can not be nil!")
			return nil, err
		}
		if payloadType != "ContextSourceRegistration" && payloadType != "Subscription" && payloadType != "Query" {
			// Payload is of Entity Type
			// Check if some other broker is registered for providing this entity or not
			var entityId string
			if _, ok := itemsMap["id"]; ok == true {
				entityId = itemsMap["id"].(string)
			} else if _, ok := itemsMap["@id"]; ok == true {
				entityId = itemsMap["@id"].(string)
			}

			if entityId == "" {
				err := errors.New("Id can not be nil!")
				return nil, err
			}
			if contextInPayload == true {
				if Context := itemsMap["@context"]; Context == nil {
					err := errors.New("@context is Empty")
					return nil, err
				}
				if Context := itemsMap["@context"].([]interface{}); len(Context) == 0 {
					err := errors.New("@context is Empty")
					return nil, err
				}
			}
		}

		// Update Context in itemMap
		if contextInPayload == true && itemsMap["@context"] != nil {
			contextItems := itemsMap["@context"].([]interface{})
			context = append(context, contextItems...)
		}
		itemsMap["@context"] = context

		if expanded, err := tb.ExpandData(itemsMap); err != nil {
			return nil, err
		} else {

			return expanded, nil
		}
	}
}

func (tb *ThinBroker) ExpandAttributePayload(r *rest.Request, context []interface{}, params ...string) ([]interface{}, error) {
	//eid := params[0]
	itemsMap, err := tb.getStringInterfaceMap(r)
	context = append(context, DEFAULT_CONTEXT)
	//get map[string]interface{} of reqBody
	//itemsMap, err := tb.getStringInterfaceMap(r)
	if err != nil {
		return nil, err
	} else {
		// Update Context in itemMap
		if itemsMap["@context"] != nil {
			contextItem := itemsMap["@context"]
			context = append(context, contextItem)
		}
		itemsMap["@context"] = context

		//Add attribute to payload, if found in params, case: Partial update
		if params != nil {
			eid := params[0]
			attrName := params[1]

			tb.ldEntities_lock.Lock()
			// Check if the attribute exists
			entity := tb.ldEntities[eid]
			entityMap := entity.(map[string]interface{})
			attrFound := false
			attrType := "" // To record whether it is Property or Relationship
			for attr, attrVal := range entityMap {
				if strings.HasSuffix(attr, "/"+attrName) {
					attrFound = true
					// Check the type of attribute (Property or Relationship)
					attrMp := attrVal.(map[string]interface{})
					attrType = attrMp["type"].(string)
					// Prepare attribute payload from partial payload
					mp := make(map[string]interface{})
					for key, val := range itemsMap {
						switch key {
						case "@context":
							continue
						default:
							mp[key] = val
							delete(itemsMap, key)
						}
					}
					mp["type"] = attrType
					itemsMap[attrName] = mp
					break
				}
			}
			if attrFound != true {
				tb.ldEntities_lock.Unlock()
				err := errors.New("Attribute not found!")
				return nil, err
			}
			tb.ldEntities_lock.Unlock()
		}
		if expanded, err := tb.ExpandData(itemsMap); err != nil {
			return nil, err
		} else {
			return expanded, nil
		}
	}
}

func (tb *ThinBroker) getTypeResolved(link string, typ string) string {
	newLink := extractLinkHeaderFields(link) // Keys in returned map are: "link", "rel" and "type"
	fmt.Println("newLink", newLink)
	var context []interface{}

	if newLink == "default" {
		newLink = DEFAULT_CONTEXT
	}

	context = append(context, newLink)

	itemsMap := make(map[string]interface{})
	fmt.Println("itemsMap", itemsMap)
	itemsMap["@context"] = context
	itemsMap["type"] = typ //Error, when entire slice typ is assigned :  invalid type value: @type value must be a string or array of strings
	resolved, err := tb.ExpandData(itemsMap)

	if err != nil {
		DEBUG.Println("Error: ", err)
		return ""
	}

	sz := Serializer{}
	typ = sz.DeSerializeType(resolved)
	return typ
}

// Expand the NGSI-LD Data with context
func (tb *ThinBroker) ExpandData(v interface{}) ([]interface{}, error) {
	dl := Expand_once()
	proc := ld.NewJsonLdProcessor()
	opts := ld.NewJsonLdOptions("")
	opts.ProcessingMode = ld.JsonLd_1_1
	opts.DocumentLoader = dl
	expanded, err := proc.Expand(v, opts)
	return expanded, err
}

//Get string-interface{} map from request body
func (tb *ThinBroker) getStringInterfaceMap(ctx interface{}) (map[string]interface{}, error) {
	// Get bite array of request body
	/*	reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		// Unmarshal using a generic interface
		var req interface{}
		err = json.Unmarshal(reqBytes, &req)
		if err != nil {
			DEBUG.Println("Invalid Request.")
			return nil, err
		}*/
	// Parse the JSON object into a map with string keys
	itemsMap := ctx.(map[string]interface{})

	if len(itemsMap) != 0 {
		return itemsMap, nil
	} else {
		return nil, errors.New("EmptyPayload!")
	}
}

/*func (tb *ThinBroker) extractLinkHeaderFields(link string) map[string]string {
	mp := make(map[string]string)
	linkArray := strings.Split(link, ";")

	for i, arrValue := range linkArray {
		linkArray[i] = strings.Trim(arrValue, " ")
		if strings.HasPrefix(arrValue, "<{{link}}>") {
			continue // TBD, context link
		} else if strings.HasPrefix(arrValue, "http") {
			mp["link"] = arrValue
		} else if strings.HasPrefix(arrValue, " rel=") {
			mp["rel"] = arrValue[6 : len(arrValue)-1] // Trimmed `rel="` and `"`
		} else if strings.HasPrefix(arrValue, " type=") {
			mp["type"] = arrValue[7 : len(arrValue)-1] // Trimmed `type="` and `"`
		}
	}

	return mp
}*/

func (tb *ThinBroker) queryOwnerOfLDEntity(eid string) string {
	inLocalBroker := true
	tb.ldEntities_lock.RLock()
	_, exist := tb.ldEntities[eid]
	inLocalBroker = exist
	tb.ldEntities_lock.RUnlock()

	if inLocalBroker == true {
		return tb.myProfile.MyURL
	} else {
		client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
		brokerURL, _ := client.GetProviderURL(eid)
		if brokerURL == "" {
			return tb.myProfile.MyURL
		}
		return brokerURL
	}
}

func (tb *ThinBroker) checkMatcheLdAttr(ctxElemAttrs []string, sid string) bool {
	tb.ldSubscriptions_lock.RLock()
	_, exist := tb.ldSubscriptions[sid]
	if exist == false {
		tb.ldSubscriptions_lock.RUnlock()
		return false
	}
	conditionList := tb.ldSubscriptions[sid].WatchedAttributes
	tb.ldSubscriptions_lock.RUnlock()
	if len(conditionList) == 0 {
		return true
	}
	matchedAtleastOnce := false
	for _, attrs1 := range ctxElemAttrs {
		for _, attrs2 := range conditionList {
			if attrs1 == attrs2 {
				matchedAtleastOnce = true
				break
			}
			if matchedAtleastOnce == true {
				break
			}
		}
	}
	return matchedAtleastOnce
}

func (tb *ThinBroker) LDNotifySubscribers(ctxElem map[string]interface{}, checkSelectedAttributes bool) {
	eid := ctxElem["@id"].(string)
	tb.LDe2sub_lock.RLock()
	defer tb.LDe2sub_lock.RUnlock()
	subscriberList := tb.entityId2LDSubcriptions[eid]
	/*for k, _ := range tb.entityId2LDSubcriptions {
		matched := tb.matchPattern(k, eid) // (pattern, id) to check if the current eid lies in the pattern given in the key.
		if matched == true {
			list := tb.entityId2Subcriptions[k]
			subscriberList = append(subscriberList, list...)
		}
	}*/
	ldAttr := make([]string, 0)
	for k, _ := range ctxElem {
		if k != "@id" && k != "@type" && k != "modifiedAt" && k != "createdAt" && k != "observationSpace" && k != "operationSpace" && k != "location" && k != "@context" {
			ldAttr = append(ldAttr, k)
		}
	}
	//send this context element to the subscriber
	for _, sid := range subscriberList {
		elements := make([]map[string]interface{}, 0)
		checkCondition := tb.checkMatcheLdAttr(ldAttr, sid)
		if checkSelectedAttributes == true && checkCondition == true {
			selectedAttributes := make([]string, 0)
			tb.ldSubscriptions_lock.RLock()
			if subscription, exist := tb.ldSubscriptions[sid]; exist {
				if subscription.Notification.Attributes != nil {
					selectedAttributes = append(selectedAttributes, subscription.Notification.Attributes...)
				}
			}
			tb.ldSubscriptions_lock.RUnlock()
			tb.ldEntities_lock.RLock()
			element := tb.ldEntities[eid]
			elementMap := element.(map[string]interface{})
			clonedElement := ldCloneWithSelectedAttributes(elementMap, selectedAttributes)
			//element := tb.ldEntities[eid]
			tb.ldEntities_lock.RUnlock()
			//elementMap := element.(map[string]interface{})
			elements = append(elements, clonedElement)
		} else {
			elements = append(elements, ctxElem)
		}
		if checkCondition == true {
			go tb.sendReliableNotifyToNgsiLDSubscriber(elements, sid)
		}
	}
}

func (tb *ThinBroker) sendReliableNotifyToNgsiLDSubscriber(elements []map[string]interface{}, sid string) {
	tb.ldSubscriptions_lock.Lock()
	ldSubscription, ok := tb.ldSubscriptions[sid]
	if ok == false {
		tb.ldSubscriptions_lock.Unlock()
	}
	subscriberURL := ldSubscription.Notification.Endpoint.URI
	if ldSubscription.Subscriber.RequireReliability == true && len(ldSubscription.Subscriber.LDNotifyCache) > 0 {
		DEBUG.Println("resend notify:  ", len(ldSubscription.Subscriber.LDNotifyCache))
		elements = append(elements, ldSubscription.Subscriber.LDNotifyCache...)
		ldSubscription.Subscriber.LDNotifyCache = make([]map[string]interface{}, 0)
	}
	tb.ldSubscriptions_lock.Unlock()
	FiwareService := ldSubscription.Subscriber.FiwareService
	FiwareServicePath := ldSubscription.Subscriber.FiwareServicePath
	err := ldPostNotifyContext(elements, sid, subscriberURL, ldSubscription.Subscriber.Integration, FiwareService, FiwareServicePath, tb.SecurityCfg)
	notifyTime := time.Now().String()
	//INFO.Println("NOTIFY: ", len(elements), ", ", sid, ", ", subscriberURL)
	if err != nil {
		INFO.Println("NOTIFY is not received by the subscriber, ", subscriberURL)

		tb.ldSubscriptions_lock.Lock()
		if ldSubscription, exist := tb.ldSubscriptions[sid]; exist {
			if ldSubscription.Subscriber.RequireReliability == true {
				ldSubscription.Subscriber.LDNotifyCache = append(ldSubscription.Subscriber.LDNotifyCache, elements...)
				ldSubscription.Notification.LastFailure = notifyTime
				ldSubscription.Notification.Status = "failed"
				tb.tmpNGSIldNotifyCache = append(tb.tmpNGSIldNotifyCache, sid)
			}
		}
		tb.ldSubscriptions_lock.Unlock()
		return
	}
	tb.updateLastSuccessParameters(notifyTime, sid)
	INFO.Println("NOTIFY is sent to the subscriber, ", subscriberURL)
}

func (tb *ThinBroker) updateLastSuccessParameters(time string, sid string) {
	tb.ldSubscriptions_lock.Lock()
	if ldSubscription, exist := tb.ldSubscriptions[sid]; exist {
		ldSubscription.Notification.LastNotification = time
		ldSubscription.Notification.LastSuccess = time
		ldSubscription.Notification.TimeSent += 1
		ldSubscription.Notification.Status = "ok"
	}
	tb.ldSubscriptions_lock.Unlock()
}

//PATCH
func (tb *ThinBroker) LDUpdateEntityAttributes(w rest.ResponseWriter, r *rest.Request) {
	var context []interface{}
	eid := r.PathParam("eid")
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		tb.ldEntities_lock.RLock()
		if _, ok := tb.ldEntities[eid]; ok == true {
			tb.ldEntities_lock.RUnlock()
			//Get a resolved object ([]interface object)
			resolved, err := tb.ExpandAttributePayload(r, context)
			if err != nil {
				if err.Error() == "EmptyPayload!" {
					rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
					return
				}
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				// Deserialize the resolved payload
				sz := Serializer{}
				deSerializedAttributePayload, err := sz.DeSerializeEntity(resolved)
				if err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					err := tb.updateAttributes(deSerializedAttributePayload, eid)
					if err != nil {
						rest.Error(w, err.Error(), 207)
						return
					}
					w.WriteHeader(204)
				}
			}
		} else {
			tb.ldEntities_lock.RUnlock()
			ownerURL := tb.queryOwnerOfLDEntity(eid)
			if ownerURL != tb.MyURL {
				ownerURL = strings.TrimSuffix(ownerURL, "/ngsi10")
				reqCxt, _ := tb.getStringInterfaceMap(r)
				//link := r.Header.Get("Link") // Pick link header if present
				//fmt.Println("Here 1..., link sending to remote broker:", link, "\nOwner URL:", ownerURL, "\nMy URL:", tb.MyURL)
				_, code := tb.updateLDAttributeValues2RemoteSite(reqCxt, ownerURL, eid)
				if code == 207 {
					//rest.Error(w, err.Error(), 404)
					//ERROR.Println(err)
					rest.Error(w, "The attribute was not found!", 404)
					return
				}
				w.WriteHeader(204)

				//return nil, err
			} else {
				ERROR.Println("The entity was not found!")
				rest.Error(w, "The entity was not found!", 404)
				return
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

//POST
func (tb *ThinBroker) LDAppendEntityAttributes(w rest.ResponseWriter, r *rest.Request) {
	var context []interface{}
	eid := r.PathParam("eid")
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		tb.ldEntities_lock.RLock()
		if _, ok := tb.ldEntities[eid]; ok == true {
			tb.ldEntities_lock.RUnlock()
			//Get a resolved object ([]interface object)
			resolved, err := tb.ExpandAttributePayload(r, context)

			if err != nil {
				if err.Error() == "EmptyPayload!" {
					rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
					return
				}
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				// Deserialize the resolved payload
				sz := Serializer{}
				deSerializedAttributePayload, err := sz.DeSerializeEntity(resolved)
				deSerializedAttributePayload["@context"] = context
				if err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					//Update createdAt for each new attribute
					for key, _ := range deSerializedAttributePayload {
						if key != "@context" && key != "modifiedAt" {
							attr := deSerializedAttributePayload[key].(map[string]interface{})
							attr["createdAt"] = time.Now().String()
							deSerializedAttributePayload[key] = attr
						}
					}

					// Write entity to tb.ldEntities
					tb.ldEntities_lock.Lock()
					entity := tb.ldEntities[eid]
					entityMap := entity.(map[string]interface{})
					multiStatus := false
					for k, attr := range deSerializedAttributePayload {
						if k != "@context" && k != "modifiedAt" {
							if _, ok := entityMap[k]; ok == true {
								multiStatus = true // atleast one duplicate attribute found
							} else {
								entityMap[k] = attr
							}
						}
					}
					entityMap["modifiedAt"] = time.Now().String()

					// Update context in entity in tb.ldEntities
					ctxList := entityMap["@context"].([]interface{})
					ctxList = append(ctxList, context...)
					entityMap["@context"] = ctxList

					tb.ldEntities[eid] = entityMap

					tb.registerLDContextElement(entityMap)
					tb.ldEntities_lock.Unlock()

					// Update Registration on Broker
					tb.appendLDAttributes(deSerializedAttributePayload, eid)
					if multiStatus == true {
						rest.Error(w, "Some duplicate attributes were found!", 207)
					} else {
						w.WriteHeader(204)
					}
				}
			}
		} else {
			tb.ldEntities_lock.RUnlock()
			ownerURL := tb.queryOwnerOfLDEntity(eid)
			if ownerURL != tb.MyURL {
				ownerURL = strings.TrimSuffix(ownerURL, "/ngsi10")
				reqCxt, _ := tb.getStringInterfaceMap(r)
				//link := r.Header.Get("Link") // Pick link header if present
				//fmt.Println("Here 1..., link sending to remote broker:", link, "\nOwner URL:", ownerURL, "\nMy URL:", tb.MyURL)
				tb.updateLDAttribute2RemoteSite(reqCxt, ownerURL, eid)
				w.WriteHeader(204)
				//return nil, err
			} else {

				rest.Error(w, "The entity was not found!", 404)
				return
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

func (tb *ThinBroker) appendLDAttributes(elem map[string]interface{}, eid string) {
	tb.ldEntityID2RegistrationID_lock.Lock()
	if rid, ok := tb.ldEntityID2RegistrationID[eid]; ok == true {
		tb.ldEntityID2RegistrationID_lock.Unlock()

		tb.ldContextRegistrations_lock.Lock()
		for k, info := range tb.ldContextRegistrations[rid].Information {
			for _, entity := range info.Entities {
				if entity.ID == eid {
					for key, attr := range elem {
						if key != "@context" {
							attrValue := attr.(map[string]interface{})
							if strings.Contains(attrValue["@type"].(string), "Property") {
								for _, existingProperty := range tb.ldContextRegistrations[rid].Information[k].Properties {
									if existingProperty == key {
										continue
									}
									tb.ldContextRegistrations[rid].Information[k].Properties = append(tb.ldContextRegistrations[rid].Information[k].Properties, key)
								}
							} else if strings.Contains(attrValue["@type"].(string), "Relationship") {
								for _, existingRelationship := range tb.ldContextRegistrations[rid].Information[k].Relationships {
									if existingRelationship == key {
										continue
									}
									tb.ldContextRegistrations[rid].Information[k].Relationships = append(tb.ldContextRegistrations[rid].Information[k].Relationships, key)
								}
							}
						}
					}
				}
			}
		}
		tb.ldContextRegistrations_lock.Unlock()
	} else {
		tb.ldEntityID2RegistrationID_lock.Unlock()
	}
}

//PATCH: Partial update, Attr name in URL, value in payload
func (tb *ThinBroker) LDUpdateEntityByAttribute(w rest.ResponseWriter, r *rest.Request) {
	var context []interface{}
	eid := r.PathParam("eid")
	attr := r.PathParam("attr")
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		tb.ldEntities_lock.RLock()
		if _, ok := tb.ldEntities[eid]; ok == true {
			tb.ldEntities_lock.RUnlock()
			//Get a resolved object ([]interface object)
			resolved, err := tb.ExpandAttributePayload(r, context, eid, attr)

			if err != nil {
				if err.Error() == "EmptyPayload!" {
					rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
					return
				}
				if err.Error() == "Attribute not found!" {
					rest.Error(w, "Attribute not found!", 404)
					return
				}
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				// Deserialize the resolved payload
				sz := Serializer{}
				deSerializedAttributePayload, err := sz.DeSerializeEntity(resolved)
				if err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					tb.updateAttributes(deSerializedAttributePayload, eid)
					w.WriteHeader(204)
				}
			}
		} else {
			tb.ldEntities_lock.RUnlock()
			ownerURL := tb.queryOwnerOfLDEntity(eid)
			if ownerURL != tb.MyURL {
				ownerURL = strings.TrimSuffix(ownerURL, "/ngsi10")
				reqCxt, _ := tb.getStringInterfaceMap(r)
				//link := r.Header.Get("Link") // Pick link header if present
				//fmt.Println("Here 1..., link sending to remote broker:", link, "\nOwner URL:", ownerURL, "\nMy URL:", tb.MyURL)
				_, code := tb.updateLDspecificAttributeValues2RemoteSite(reqCxt, ownerURL, eid, attr)
				if code == 404 {

					rest.Error(w, "The attribute was not found!", 404)
					return
				}
				w.WriteHeader(204)
				//return nil, err
			} else {

				ERROR.Println("The entity was not found!")
				rest.Error(w, "The entity was not found!", 404)
				return
			}
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

func (tb *ThinBroker) updateAttributes(elem map[string]interface{}, eid string) error {
	tb.ldEntities_lock.Lock()
	entity := tb.ldEntities[eid]
	entityMap := entity.(map[string]interface{})
	missing := false
	for k, _ := range elem {
		if k != "@context" && k != "modifiedAt" {
			if _, ok := entityMap[k]; ok == true {
				entityAttrMap := entityMap[k].(map[string]interface{}) // existing
				attrMap := elem[k].(map[string]interface{})            // to be updated as
				if strings.Contains(attrMap["type"].(string), "Property") {
					if attrMap["value"] != nil {
						entityAttrMap["value"] = attrMap["value"]
					}
					if attrMap["observedAt"] != nil {
						entityAttrMap["observedAt"] = attrMap["observedAt"]
					}
					if attrMap["datasetId"] != nil {
						entityAttrMap["datasetId"] = attrMap["datasetId"]
					}
					if attrMap["instanceId"] != nil {
						entityAttrMap["instanceId"] = attrMap["instanceId"]
					}
					if attrMap["unitCode"] != nil {
						entityAttrMap["unitCode"] = attrMap["unitCode"]
					}
				} else if strings.Contains(attrMap["type"].(string), "Relationship") {
					if attrMap["object"] != nil {
						entityAttrMap["object"] = attrMap["object"]
					}
					if attrMap["providedBy"] != nil {
						entityAttrMap["providedBy"] = attrMap["providedBy"]
					}
					if attrMap["datasetId"] != nil {
						entityAttrMap["datasetId"] = attrMap["datasetId"]
					}
					if attrMap["instanceId"] != nil {
						entityAttrMap["instanceId"] = attrMap["instanceId"]
					}
				}
				entityAttrMap["modifiedAt"] = time.Now().String()
				entityMap[k] = entityAttrMap
			} else {
				missing = true
				ERROR.Println("Attribute", k, "was not found in the entity!")
			}
		}
	}
	entityMap["modifiedAt"] = time.Now().String()
	tb.ldEntities[eid] = entityMap

	// registration of entity is not required on discovery while attribute updation
	//tb.registerLDContextElement(entityMap)

	// send notification to the subscriber

	go tb.LDNotifySubscribers(entityMap, true)

	tb.ldEntities_lock.Unlock()

	if missing == true {
		err := errors.New("Some attributes were not found!")
		return err
	}
	return nil
}

func (tb *ThinBroker) ldDeleteEntity(eid string) error {
	tb.ldEntities_lock.Lock()
	if tb.ldEntities[eid] != nil {
		delete(tb.ldEntities, eid)
	} else {
		tb.ldEntities_lock.Unlock()
		ERROR.Println("Entity not found!")
		err := errors.New("Entity not found!")
		return err
	}

	// Delete registration from Broker
	tb.ldEntityID2RegistrationID_lock.Lock()
	tb.ldContextRegistrations_lock.Lock()

	rid := tb.ldEntityID2RegistrationID[eid]
	delete(tb.ldContextRegistrations, rid)
	delete(tb.ldEntityID2RegistrationID, eid)

	tb.ldContextRegistrations_lock.Unlock()
	tb.ldEntityID2RegistrationID_lock.Unlock()

	// Unregister entity from Discovery
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.UnregisterEntity(eid)
	if err != nil {
		ERROR.Println(err)
	}

	tb.ldEntities_lock.Unlock()
	return nil
}

func (tb *ThinBroker) LDDeleteEntityAttribute(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")
	var attr = r.PathParam("attr")
	var newEid string
	//if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
	if r.Header.Get("fiware-service") != "" {
		newEid = eid + "@" + r.Header.Get("fiware-service")
		w.Header().Set("fiware-service", r.Header.Get("fiware-service"))
	} else {
		newEid = eid + "@" + "default"
	}
	//}
	err := tb.ldDeleteEntityAttribute(newEid, attr /*req*/)

	if err == nil {
		w.WriteHeader(204)
	} else {
		rest.Error(w, err.Error(), 404)
		return
	}
}

func (tb *ThinBroker) ldDeleteEntityAttribute(eid string, attr string /*req interface{}*/) error {
	tb.ldEntities_lock.Lock()
	if tb.ldEntities[eid] != nil {
		entityMap := tb.ldEntities[eid].(map[string]interface{})

		attrExists := false
		//for i := 0; i < len(tb.ldEntities[eid].Properties); i++ {
		for attrN, _ := range entityMap {
			if strings.HasSuffix(attrN, "/"+attr) {
				attrExists = true

				delete(entityMap, attrN)
				tb.ldEntities[eid] = entityMap

			}
		}
		if attrExists == false {
			tb.ldEntities_lock.Unlock()
			ERROR.Println("Attribute not found!")
			err := errors.New("Attribute not found!")
			return err
		}
		// Deleting attribute from registration at Broker: Get rid at broker
		rid := ""
		tb.ldEntityID2RegistrationID_lock.RLock()
		if _, ok := tb.ldEntityID2RegistrationID[eid]; ok == true {
			rid = tb.ldEntityID2RegistrationID[eid]

		}
		tb.ldEntityID2RegistrationID_lock.RUnlock()

		// Deleting attribute from registration at Broker: Update registration at broker, if found
		if rid != "" { // Registration is present at Broker; for registrations created explicitly at FogFlow.
			tb.ldContextRegistrations_lock.Lock()
			attrType := ""
			// update registration at broker here.
			for k, info := range tb.ldContextRegistrations[rid].Information {
				for _, entity := range info.Entities {
					if entity.ID == eid {
						if strings.Contains(attrType, "Property") {
							for key, property := range tb.ldContextRegistrations[rid].Information[k].Properties {
								if property == attr {
									tb.ldContextRegistrations[rid].Information[k].Properties = append(tb.ldContextRegistrations[rid].Information[k].Properties[:key], tb.ldContextRegistrations[rid].Information[k].Properties[key+1:]...)
									break
								}
							}
						} else if strings.Contains(attrType, "Relationship") {
							for key, relationship := range tb.ldContextRegistrations[rid].Information[k].Relationships {
								if relationship == attr {
									tb.ldContextRegistrations[rid].Information[k].Relationships = append(tb.ldContextRegistrations[rid].Information[k].Relationships[:key], tb.ldContextRegistrations[rid].Information[k].Relationships[key+1:]...)
									break
								}
							}
						}
					}
				}
			}
			tb.ldContextRegistrations_lock.Unlock()
		}
		// Update Registration at Discovery

		tb.registerLDContextElement(entityMap)

	} else {
		tb.ldEntities_lock.Unlock()

		ERROR.Println("Entity not found!")
		err := errors.New("Entity not found!")
		return err

	}
	tb.ldEntities_lock.Unlock()
	//              }
	/*else {
	        tb.ldEntities_lock.Unlock()
	        ERROR.Println("Entity not found!")
	        err := errors.New("Entity not found!")
	        return err
	}*/
	//              }
	//      }
	return nil
}

func (tb *ThinBroker) ldEntityGetByAttribute(attrs []string, fiwareService string) []interface{} {
	var entities []interface{}
	tb.ldEntities_lock.Lock()
	for _, entity := range tb.ldEntities {
		entityMap := entity.(map[string]interface{})
		if strings.HasSuffix(entityMap["@id"].(string), fiwareService) == true {
			allExist := true
			for _, attr := range attrs {
				if _, ok := entityMap[attr]; ok != true {
					allExist = false
				}
			}
			if allExist == true {
				compactEntity := tb.createOriginalPayload(entityMap)
				resultEntity := compactEntity.(map[string]interface{})
				actualEId := getActualEntity(resultEntity)
				resultEntity["id"] = actualEId
				//compactEntity := tb.createOriginalPayload(entityMap)
				delete(resultEntity, "fiwareServicePath")
				entities = append(entities, resultEntity)
			}
		}
	}
	tb.ldEntities_lock.Unlock()
	return entities
}

func (tb *ThinBroker) ldEntityGetById(eids []string, typ []string, fiwareService string) []interface{} {
	var newEid string
	tb.ldEntities_lock.Lock()
	var entities []interface{}
	for index, eid := range eids {
		newEid = eid + "@" + fiwareService
		if entity, ok := tb.ldEntities[newEid]; ok == true {
			entityMap := entity.(map[string]interface{})
			typeArray := entityMap["@type"].([]interface{})
			newType := typeArray[0].(string)
			if newType == typ[index] {
				compactEntity := tb.createOriginalPayload(entityMap)
				resultEntity := compactEntity.(map[string]interface{})
				actualId := getActualEntity(resultEntity)
				resultEntity["@id"] = actualId
				delete(resultEntity, "fiwareServicePath")
				entities = append(entities, resultEntity)
			}
		}
	}
	tb.ldEntities_lock.Unlock()
	return entities
}

func (tb *ThinBroker) ldEntityGetByType(typs []string, link string, fiwareService string) ([]interface{}, error) {
	var entities []interface{}
	typ := typs[0]
	if link != "" {
		typ = tb.getTypeResolved(link, typ)
		if typ == "" {
			err := errors.New("Type not resolved!")
			return nil, err
		}
	}
	tb.ldEntities_lock.Lock()
	/*for _, entity := range tb.ldEntities {
		entityMap := entity.(map[string]interface{})
		if strings.HasSuffix(entityMap["id"].(string), fiwareService) == true {
			if entityMap["type"] == typ {
				compactEntity := tb.createOriginalPayload(entityMap)
				compactEntityMap := compactEntity.(map[string]interface{})
				compactEntityMap["id"], _ = FiwareId(compactEntityMap["id"].(string))
				delete(compactEntityMap, "fiwareServicePath")
				entities = append(entities, compactEntityMap)
			}
		}
	}*/
	for index, value := range tb.entityTypeTOEntityId[typ] {
		fmt.Println("value", tb.ldEntities[value])
		if result, okey := tb.ldEntities[value]; okey == true {
			fmt.Println("result", result)
			compactEntity := tb.createOriginalPayload(result)
			compactEntityMap := compactEntity.(map[string]interface{})
			compactEntityMap["id"], _ = FiwareId(compactEntityMap["id"].(string))
			delete(compactEntityMap, "fiwareServicePath")
			entities = append(entities, compactEntityMap)
		} else { // update tb.entityTypeTOEntityId for deleted entity
			// reslice the slice to manage delete entity
			tb.entityTypeTOEntityId[typ] = reslice(tb.entityTypeTOEntityId[typ], index)
		}
	}
	tb.ldEntities_lock.Unlock()
	return entities, nil
}

func (tb *ThinBroker) ldEntityGetByIdPattern(idPatterns []string, typ []string) []interface{} {
	var entities []interface{}

	for eid, entity := range tb.ldEntities {
		entityMap := entity.(map[string]interface{})
		for index, idPattern := range idPatterns {
			if strings.Contains(idPattern, ".*") && strings.Contains(idPattern, "*.") {
				idPattern = strings.Trim(idPattern, ".*")
				idPattern = strings.Trim(idPattern, "*.")
				if strings.Contains(eid, idPattern) {
					if entityMap["type"] == typ[index] {
						compactEntity := tb.createOriginalPayload(entity)
						resultEntity2 := compactEntity.(map[string]interface{})
						actualEId2 := getActualEntity(resultEntity2)
						resultEntity2["id"] = actualEId2
						//compactEntity := tb.createOriginalPayload(entityMap)
						delete(resultEntity2, "fiwareServicePath")
						entities = append(entities, resultEntity2)
						break
					}
				}
			}
			if strings.Contains(idPattern, ".*") {
				idPattern = strings.Trim(idPattern, ".*")
				if strings.HasPrefix(eid, idPattern) {
					if entityMap["type"] == typ[index] {
						compactEntity := tb.createOriginalPayload(entity)
						resultEntity3 := compactEntity.(map[string]interface{})
						actualEId3 := getActualEntity(resultEntity3)
						resultEntity3["id"] = actualEId3
						//compactEntity := tb.createOriginalPayload(entityMap)
						delete(resultEntity3, "fiwareServicePath")
						entities = append(entities, resultEntity3)
						break
					}
				}
			}
			if strings.Contains(idPattern, "*.") {
				idPattern = strings.Trim(idPattern, "*.")
				if strings.HasSuffix(eid, idPattern) {
					if entityMap["type"] == typ[index] {
						compactEntity := tb.createOriginalPayload(entity)
						resultEntity4 := compactEntity.(map[string]interface{})
						actualEId4 := getActualEntity(resultEntity4)
						resultEntity4["id"] = actualEId4
						//compactEntity := tb.createOriginalPayload(entityMap)
						delete(resultEntity4, "fiwareServicePath")
						entities = append(entities, resultEntity4)
						break
					}
				}
			}
		}
	}
	return entities
}

// Subscription
func (tb *ThinBroker) UpdateLDSubscription(w rest.ResponseWriter, r *rest.Request) {
	//var context []interface{}
	sid := r.PathParam("sid")
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		/*if link := r.Header.Get("Link"); link != "" {
			link := extractLinkHeaderFields(link) // Keys in returned map are: "link", "rel" and "type"
			if link != DEFAULT_CONTEXT {
				context = append(context, linkMap["rel"]) // Make use of "link" and "type" also
			}
		}
		context = append(context, DEFAULT_CONTEXT)*/
		var context []interface{}
		contextInPayload := false
		cType := r.Header.Get("Content-Type")
		cTypeInLower := strings.ToLower(cType)
		//Link := r.Header.Get("Link")
		if cTypeInLower == "application/ld+json" {
			contextInPayload = true
		} else {
			if link := r.Header.Get("Link"); link != "" {
				link := extractLinkHeaderFields(link)
				if link == "default" {
					context = append(context, DEFAULT_CONTEXT)
				} else {
					context = append(context, link)
				}
			} else {
				context = append(context, DEFAULT_CONTEXT)
			}
		}
		tb.ldSubscriptions_lock.Lock()
		if _, ok := tb.ldSubscriptions[sid]; ok == true {
			tb.ldSubscriptions_lock.Unlock()
			reqBytes, _ := ioutil.ReadAll(r.Body)
			var LDUpdateSubscribeCtxReq interface{}

			err := json.Unmarshal(reqBytes, &LDUpdateSubscribeCtxReq)

			if err != nil {
				err := errors.New("not able to decode orion update")
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resolved, err := tb.ExpandPayload(LDUpdateSubscribeCtxReq, context, contextInPayload) // Context in Link header

			if err != nil {
				if err.Error() == "EmptyPayload!" {
					rest.Error(w, "Empty payloads are not allowed in this operation!", 400)
					return
				}
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				sz := Serializer{}
				deSerializedSubscription, err := sz.DeSerializeSubscription(resolved)

				if err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					/*tb.UpdateSubscriptionInMemory(deSerializedSubscription, sid)

					// Update in discovery here.
					tb.ldSubscriptions_lock.RLock()
					subReq := tb.ldSubscriptions[sid]
					tb.ldSubscriptions_lock.RUnlock()

					if err := tb.SubscribeLDContextAvailability(subReq); err != nil {
						rest.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(204)*/
					tb.ldSubscriptions_lock.RLock()
					//subReq := tb.ldSubscriptions[sid]
					tb.ldSubscriptions_lock.RUnlock()
					err := tb.UpdateLDContextAvailability(deSerializedSubscription, sid)
					if err != nil {
						rest.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						// update in broker memory
						tb.UpdateSubscriptionInMemory(deSerializedSubscription, sid)
					}
					w.WriteHeader(204)
				}
			}
		} else {
			tb.ldSubscriptions_lock.Unlock()
			rest.Error(w, "Resource not found!", 404)
			return
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", 400)
		return
	}
}

// update subscription context availability in discovery
func (tb *ThinBroker) UpdateLDContextAvailability(subReq LDSubscriptionRequest, sid string) error {
	ctxAvailabilityRequest := SubscribeContextAvailabilityRequest{}

	for key, entity := range subReq.Entities {
		if entity.IdPattern != "" {
			entity.IsPattern = true
		}
		subReq.Entities[key] = entity
	}
	ctxAvailabilityRequest.Entities = subReq.Entities
	ctxAvailabilityRequest.Attributes = subReq.WatchedAttributes
	//copy(ctxAvailabilityRequest.Attributes, subReq.Notification.Attributes)
	ctxAvailabilityRequest.Reference = tb.MyURL + "/notifyContextAvailability"
	ctxAvailabilityRequest.Duration = subReq.Expires
	eid := ""
	for key, value := range tb.availabilitySub2MainSub {
		value = tb.availabilitySub2MainSub[key]
		if value == sid {
			eid = key
			break
		}
	}
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	AvailabilitySubID, err := client.UpdateLDContextAvailability(&ctxAvailabilityRequest, eid)

	if AvailabilitySubID != "" {
		tb.subLinks_lock.Lock()
		notifyMessage, alreadyBack := tb.tmpNGSILDNotifyCache[AvailabilitySubID]
		tb.subLinks_lock.Unlock()
		if alreadyBack == true {
			INFO.Println("========forward the availability notify that arrived earlier===========")
			tb.handleNGSI9Notify(subReq.Id, notifyMessage)

			tb.subLinks_lock.Lock()
			delete(tb.tmpNGSILDNotifyCache, AvailabilitySubID)
			tb.subLinks_lock.Unlock()
		}
		return nil
	} else {
		INFO.Println("failed to subscribe the availability of requested entities ", err)
		return err
	}
}
func (tb *ThinBroker) UpdateSubscriptionInMemory(sub LDSubscriptionRequest, sid string) {
	tb.ldSubscriptions_lock.Lock()
	subscription := tb.ldSubscriptions[sid]
	if sub.Name != "" {
		subscription.Name = sub.Name
	}
	if sub.Description != "" {
		subscription.Description = sub.Description
	}
	if sub.Expires != "" {
		subscription.Expires = sub.Expires
	}
	if sub.Status != "" {
		subscription.Status = sub.Status
	}
	if sub.IsActive != false {
		subscription.IsActive = sub.IsActive
	}
	if sub.Entities != nil {
		subscription.Entities = sub.Entities
	}
	if sub.WatchedAttributes != nil {
		subscription.WatchedAttributes = sub.WatchedAttributes
	}
	if sub.Notification.Attributes != nil {
		subscription.Notification.Attributes = sub.Notification.Attributes
	}
	if sub.Notification.Format != "" {
		subscription.Notification.Format = sub.Notification.Format
	}
	if sub.Notification.Endpoint.URI != "" {
		subscription.Notification.Endpoint.URI = sub.Notification.Endpoint.URI
	}
	if sub.Notification.Endpoint.Accept != "" {
		subscription.Notification.Endpoint.Accept = sub.Notification.Endpoint.Accept
	}
	if sub.TimeInterval != 0 {
		subscription.TimeInterval = sub.TimeInterval
	}
	if sub.Throttling != 0 {
		subscription.Throttling = sub.Throttling

	}
	if sub.Q != "" {
		subscription.Q = sub.Q
	}
	nilGeoQ := GeoQuery{}
	if sub.GeoQ != nilGeoQ {
		subscription.GeoQ = sub.GeoQ
	}
	nilTemporalQ := TemporalQuery{}
	if sub.TemporalQ != nilTemporalQ {
		subscription.TemporalQ = sub.TemporalQ
	}
	if sub.Csf != "" {
		subscription.Csf = sub.Csf
	}
	tb.ldSubscriptions_lock.Unlock()
}

func (tb *ThinBroker) deleteLDSubscription(sid string) error {
	tb.ldSubscriptions_lock.Lock()
	defer tb.ldSubscriptions_lock.Unlock()
	tb.subLinks_lock.RLock()
	defer tb.subLinks_lock.RUnlock()
	//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
	for index, otherSubID := range tb.main2Other[sid] {
		if index == 0 {
			tb.UnsubscribeContextAvailability(otherSubID)
		} else {
			//fmt.Println(" tb.ldSubscriptions[otherSubID]", tb.ldSubscriptions[otherSubID])
			if _, exist := tb.ldSubscriptions[otherSubID]; exist {
				unsubscribeContextLDProvider(otherSubID, tb.ldSubscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
			} else {
				tb.UnsubscribeContextAvailability(otherSubID)
			}
		}
	}

	// remove the subscription from the map
	if _, oK := tb.ldSubscriptions[sid]; oK {
		delete(tb.ldSubscriptions, sid)
	}
	//delete(tb.subscriptions, sid)

	return nil
}

func (tb *ThinBroker) getLDSubscription(sid string) *LDSubscriptionRequest {
	tb.ldSubscriptions_lock.Lock()
	subscription := tb.ldSubscriptions[sid]
	tb.ldSubscriptions_lock.Unlock()
	return subscription
}

func (tb *ThinBroker) GetLDSubscriptions(w rest.ResponseWriter, r *rest.Request) {
	if accept := r.Header.Get("Accept"); accept == "application/ld+json" {
		tb.ldSubscriptions_lock.RLock()
		defer tb.ldSubscriptions_lock.RUnlock()

		subscriptions := make(map[string]LDSubscriptionRequest)

		for sid, sub := range tb.ldSubscriptions {
			subscriptions[sid] = *sub
		}
		w.WriteHeader(200)
		w.WriteJson(&subscriptions)
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", http.StatusBadRequest)
	}
}

func (tb *ThinBroker) matchPattern(pattern string, id string) bool {
	if strings.Contains(pattern, ".*") && strings.Contains(pattern, "*.") {
		id = strings.Trim(pattern, ".*")
		id = strings.Trim(pattern, "*.")
		if strings.Contains(id, pattern) {
			return true
		}

	} else if strings.Contains(pattern, ".*") {
		id = strings.Trim(pattern, ".*")
		if strings.HasPrefix(id, pattern) {
			return true
		}

	} else if strings.Contains(pattern, "*.") {
		id = strings.Trim(pattern, "*.")
		if strings.HasSuffix(id, pattern) {
			return true
		}
	}
	return false
}
