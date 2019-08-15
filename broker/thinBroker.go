package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/satori/go.uuid"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
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

	//mapping from main subscription to other related subscriptions
	main2Other              map[string][]string
	availabilitySub2MainSub map[string]string
	tmpNGSI9NotifyCache     map[string]*NotifyContextAvailabilityRequest
	subLinks_lock           sync.RWMutex

	//list of all updated context entities
	entities      map[string]*ContextElement //latest view of context entities
	entities_lock sync.RWMutex

	//mapping from entityID to subscriptionID
	entityId2Subcriptions map[string][]string
	e2sub_lock            sync.RWMutex

	//counter of heartbeat
	counter int64
}

func (tb *ThinBroker) Start(cfg *Config) {
	if cfg.HTTPS.Enabled == true {
		tb.MyURL = "https://" + cfg.ExternalIP + ":" + strconv.Itoa(cfg.Broker.HTTPSPort) + "/ngsi10"
		tb.IoTDiscoveryURL = cfg.GetDiscoveryURL(true)
	} else {
		tb.MyURL = "http://" + cfg.ExternalIP + ":" + strconv.Itoa(cfg.Broker.HTTPPort) + "/ngsi10"
		tb.IoTDiscoveryURL = cfg.GetDiscoveryURL(false)
	}

	tb.myEntityId = tb.id

	tb.SecurityCfg = &cfg.HTTPS

	tb.MyLocation = cfg.Location

	tb.subscriptions = make(map[string]*SubscribeContextRequest)
	tb.tmpNGSI10NotifyCache = make([]string, 0)

	tb.entities = make(map[string]*ContextElement)
	tb.entityId2Subcriptions = make(map[string][]string)

	tb.availabilitySub2MainSub = make(map[string]string)
	tb.tmpNGSI9NotifyCache = make(map[string]*NotifyContextAvailabilityRequest)
	tb.main2Other = make(map[string][]string)

	tb.myProfile.BID = tb.myEntityId
	tb.myProfile.MyURL = tb.MyURL

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

func (tb *ThinBroker) getSubscription(sid string) *SubscribeContextRequest {
	tb.subscriptions_lock.RLock()
	defer tb.subscriptions_lock.RUnlock()

	if sub, exist := tb.subscriptions[sid]; exist {
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

	// send out the response
	w.WriteHeader(200)

	updateCtxResp := UpdateContextResponse{}
	// updateCtxResp.ErrorCode.Code = 200
	// updateCtxResp.ErrorCode.ReasonPhrase = "OK"
	w.WriteJson(&updateCtxResp)

	if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
		tb.handleInternalUpdateContext(&updateCtxReq)
	} else {
		tb.handleExternalUpdateContext(&updateCtxReq)
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
func (tb *ThinBroker) handleExternalUpdateContext(updateCtxReq *UpdateContextRequest) {
	// perform the update action accordingly
	switch strings.ToUpper(updateCtxReq.UpdateAction) {
	case "UPDATE", "APPEND":
		for _, ctxElem := range updateCtxReq.ContextElements {
			// just in case this is orion ngsi v1
			ctxElem.SetEntityID()

			brokerURL := tb.queryOwnerOfEntity(ctxElem.Entity.ID)

			if brokerURL == tb.myProfile.MyURL {
				tb.UpdateContext2LocalSite(&ctxElem)
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
		brokerURL := client.GetProviderURL(eid)
		if brokerURL == "" {
			return tb.myProfile.MyURL
		}
		return brokerURL
	}
}

func (tb *ThinBroker) UpdateContext2LocalSite(ctxElem *ContextElement) {
	tb.entities_lock.Lock()
	eid := ctxElem.Entity.ID
	hasUpdatedMetadata := hasUpdatedMetadata(ctxElem, tb.entities[eid])
	tb.entities_lock.Unlock()

	// apply the new update to the entity in the entity map
	tb.updateContextElement(ctxElem)

	// propogate this update to its subscribers
	go tb.notifySubscribers(ctxElem, true)

	// register the entity if there is any changes on attribute list, domain metadata
	if hasUpdatedMetadata == true {
		tb.registerContextElement(ctxElem)
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

func (tb *ThinBroker) sendReliableNotify(elements []ContextElement, sid string) {
	tb.subscriptions_lock.Lock()
	subscription, ok := tb.subscriptions[sid]
	if ok == false {
		tb.subscriptions_lock.Unlock()
		return
	}

	subscriberURL := subscription.Reference
	IsOrionBroker := subscription.Subscriber.IsOrion

	//check if there is any element that has not been received
	if subscription.Subscriber.RequireReliability == true && len(subscription.Subscriber.NotifyCache) > 0 {
		DEBUG.Println("resend notify:  ", len(subscription.Subscriber.NotifyCache))

		for _, pCtxElem := range subscription.Subscriber.NotifyCache {
			elements = append(elements, *pCtxElem)
		}

		subscription.Subscriber.NotifyCache = make([]*ContextElement, 0)
	}

	tb.subscriptions_lock.Unlock()

	INFO.Println("NOTIFY: ", len(elements), ", ", sid, ", ", subscriberURL, ", ", IsOrionBroker)

	err := postNotifyContext(elements, sid, subscriberURL, IsOrionBroker, tb.SecurityCfg)
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
	u1, err := uuid.NewV4()
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

func (tb *ThinBroker) UnsubscribeContextAvailability(sid string) error {
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	err := client.UnsubscribeContextAvailability(sid)
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
	if r.Header.Get("User-Agent") != "lightweight-iot-broker" {
		//for external subscription, we need to cancel all subscriptions to IoT Discovery and other Brokers
		for index, otherSubID := range tb.main2Other[subID] {
			if index == 0 {
				tb.UnsubscribeContextAvailability(otherSubID)
			} else {
				unsubscribeContextProvider(otherSubID, tb.subscriptions[otherSubID].Subscriber.BrokerURL, tb.SecurityCfg)
			}
		}
	}

	// remove the subscription from the map
	delete(tb.subscriptions, subID)
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

func (tb *ThinBroker) handleNGSI9Notify(mainSubID string, notifyContextAvailabilityReq *NotifyContextAvailabilityRequest) {
	var action string
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
