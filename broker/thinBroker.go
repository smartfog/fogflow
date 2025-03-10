package main

import (
	"strings"
	"sync"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"

	"github.com/ant0ine/go-json-rest/rest"
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
	tb.MyURL = cfg.GetBrokerURL()
	tb.IoTDiscoveryURL = cfg.GetDiscoveryURL()

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
	tb.myProfile.MyURL = cfg.GetExternalBrokerURL()

	// register itself to the IoT discovery
	tb.registerMyself()
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
			if subscription.Subscriber.RequireReliability && len(subscription.Subscriber.NotifyCache) > 0 {
				hasCachedNotification = true
			}
		}
		tb.subscriptions_lock.Unlock()

		if hasCachedNotification {
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

	// send the first heartbeat message
	tb.sendHeartBeat()

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
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println(" TO REMOVE ENTITY ", eid)
	}

	//remove it from the local entity map
	tb.entities_lock.Lock()
	delete(tb.entities, eid)
	tb.entities_lock.Unlock()

	// inform the subscribers that this entity is deleted by sending a empty context element without any attribute, metadata
	emptyElement := ContextElement{}
	emptyElement.Entity.ID = eid
	tb.notifySubscribers(&emptyElement, "", false)

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

// handle context updates from external applications/devices

func (tb *ThinBroker) handleInternalUpdateContext(updateCtxReq *UpdateContextRequest) {

	// DEBUG.Println("updateCtxReq", updateCtxReq)

	switch strings.ToUpper(updateCtxReq.UpdateAction) {
	case "UPDATE":
		for _, ctxElem := range updateCtxReq.ContextElements {
			tb.UpdateContext2LocalSite(&ctxElem, updateCtxReq.Correlator)
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
			brokerURL := tb.queryOwnerOfEntity(ctxElem.Entity.ID)
			if brokerURL == tb.myProfile.MyURL {
				tb.UpdateContext2LocalSite(&ctxElem, updateCtxReq.Correlator, w)
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

func (tb *ThinBroker) queryOwnerOfEntity(eid string) string {
	inLocalBroker := true

	tb.entities_lock.RLock()
	_, exist := tb.entities[eid]
	inLocalBroker = exist
	tb.entities_lock.RUnlock()

	if inLocalBroker {
		return tb.myProfile.MyURL
	}

	// ask the discovery service which broker is hosting this entity
	client := NGSI9Client{IoTDiscoveryURL: tb.IoTDiscoveryURL, SecurityCfg: tb.SecurityCfg}
	brokerURL, _ := client.GetProviderURL(eid)
	if brokerURL == "" {
		return tb.myProfile.MyURL
	}

	return brokerURL
}

func (tb *ThinBroker) UpdateContext2LocalSite(ctxElem *ContextElement, correlator string, params ...rest.ResponseWriter) {

	// DEBUG.Println("elements", ctxElem)

	// register the entity if there is any changes on attribute list, domain metadata
	tb.entities_lock.Lock()
	eid := ctxElem.Entity.ID
	hasUpdatedMetadata := hasUpdatedMetadata(ctxElem, tb.entities[eid])
	tb.entities_lock.Unlock()

	if hasUpdatedMetadata {
		tb.registerContextElement(ctxElem)
	}

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("before updating ctxElem: ", ctxElem)
	}
	// apply the new update to the entity in the entity map
	tb.updateContextElement(ctxElem)

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("after updating ctxElem: ", ctxElem)
	}

	// propogate this update to its subscribers
	tb.notifySubscribers(ctxElem, correlator, true)
}

func (tb *ThinBroker) UpdateContext2RemoteSite(ctxElem *ContextElement, updateAction string, brokerURL string) {
	switch updateAction {
	case "UPDATE":
		// INFO.Println(brokerURL)
		client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: tb.SecurityCfg}
		client.UpdateContext(ctxElem)

	case "DELETE":
		client := NGSI10Client{IoTBrokerURL: brokerURL, SecurityCfg: tb.SecurityCfg}
		client.DeleteContext(&ctxElem.Entity)
	}
}

func (tb *ThinBroker) notifySubscribers(ctxElem *ContextElement, correlator string, checkSelectedAttributes bool) {
	eid := ctxElem.Entity.ID

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("elements to check against subscriptions: ", ctxElem, ", checkSelectedAttributes: ", checkSelectedAttributes)
	}

	tb.e2sub_lock.RLock()
	defer tb.e2sub_lock.RUnlock()
	subscriberList := tb.entityId2Subcriptions[eid]
	subscriberList = append(subscriberList, tb.entityId2Subcriptions[ctxElem.GetTypeWildCard()]...)

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("subscriberList: ", subscriberList)
	}
	//send this context element to the subscriber
	for _, sid := range subscriberList {

		if isProsumerSubscription(correlator, sid) {
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("notification comes from the prosumer, thus I avoid to notify back the prosumer, sid", sid, ", correlator ", correlator)
			}
			continue
		}

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("elements to check against subscription:", sid, ", checkSelectedAttributes: ", checkSelectedAttributes, "correlator", correlator)
		}

		elements := make([]ContextElement, 0)

		beTheSame := false

		// check if both the producer and subscriber of this update is the same originator
		tb.subscriptions_lock.RLock()
		if subscription, exist := tb.subscriptions[sid]; exist {
			originator := subscription.Subscriber.Correlator
			if correlator != "" && originator != "" && correlator == originator {
				beTheSame = true
			}
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("session ID from producer ", correlator, ", subscriber ", originator, "beTheSame ", beTheSame)
			}
		}
		tb.subscriptions_lock.RUnlock()

		if beTheSame {
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println(" ======= producer and subscriber are the same ===========")
			}
			continue
		}

		if checkSelectedAttributes {

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("subscription to check: ", sid, " against attribute: ", tb.subscriptions[sid].Attributes)
			}

			selectedAttributes := make([]string, 0)

			tb.subscriptions_lock.RLock()
			if subscription, exist := tb.subscriptions[sid]; exist {
				if subscription.Attributes != nil {
					selectedAttributes = append(selectedAttributes, tb.subscriptions[sid].Attributes...)
				}
			}
			tb.subscriptions_lock.RUnlock()

			tb.entities_lock.RLock()
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("context Entity: ", tb.entities[eid])
			}
			element := tb.entities[eid].CloneWithSelectedAttributes(selectedAttributes)
			tb.entities_lock.RUnlock()

			elements = append(elements, *element)
		} else {
			elements = append(elements, *ctxElem)
		}

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("elements", elements, "sid", sid)
		}
		go tb.sendReliableNotify(elements, sid)
	}
}

func (tb *ThinBroker) notifyOneSubscriberWithCurrentStatus(entities []EntityId, sid string) {
	elements := make([]ContextElement, 0)
	// check if the subscription still exists; if yes, then find out the selected attribute list
	tb.subscriptions_lock.RLock()

	subscription, ok := tb.subscriptions[sid]
	if !ok {
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
	Send the Notification to NGSIV1 subscriber
*/

func (tb *ThinBroker) sendReliableNotifyToSubscriber(elements []ContextElement, sid string) {
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("elements", elements, "sid", sid)
	}
	tb.subscriptions_lock.Lock()
	subscription, ok := tb.subscriptions[sid]
	if !ok {
		tb.subscriptions_lock.Unlock()
	}
	subscriberURL := subscription.Reference

	notifyVersion := subscription.Subscriber.DestinationType
	Tenant := subscription.Subscriber.Tenant

	if subscription.Subscriber.RequireReliability && len(subscription.Subscriber.NotifyCache) > 0 {
		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("resend notify:  ", len(subscription.Subscriber.NotifyCache))
		}
		for _, pCtxElem := range subscription.Subscriber.NotifyCache {
			elements = append(elements, *pCtxElem)
		}
		subscription.Subscriber.NotifyCache = make([]*ContextElement, 0)
	}
	tb.subscriptions_lock.Unlock()

	//INFO.Println("NOTIFY: ", len(elements), ", ", sid, ", ", subscriberURL, ", ", DestinationBroker)
	// DEBUG.Println(elements)

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("elements:", elements, "sid:", sid, "subscriberURL:", subscriberURL, "DestinationBroker:", notifyVersion, "Tenant,", Tenant)
	}

	if len(elements) > 0 {

		err := postNotifyContext(elements, sid, subscriberURL, notifyVersion, Tenant, tb.SecurityCfg)

		if err != nil {
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("NOTIFY is not received by the subscriber, ", subscriberURL)
			}

			tb.subscriptions_lock.Lock()
			if subscription, exist := tb.subscriptions[sid]; exist {
				if subscription.Subscriber.RequireReliability {
					for _, ctxElem := range elements {
						subscription.Subscriber.NotifyCache = append(subscription.Subscriber.NotifyCache, &ctxElem)
					}

					tb.tmpNGSI10NotifyCache = append(tb.tmpNGSI10NotifyCache, sid)
				}
			}
			tb.subscriptions_lock.Unlock()
		}
	}

}

/*
	Identify the subscriber by using SubscriptionId
*/

func (tb *ThinBroker) sendReliableNotify(elements []ContextElement, sid string) {
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("elements", elements, "sid", sid)
	}
	tb.subscriptions_lock.Lock()
	_, ok := tb.subscriptions[sid]
	if ok {
		tb.subscriptions_lock.Unlock()
		tb.sendReliableNotifyToSubscriber(elements, sid)
	} else {
		tb.subscriptions_lock.Unlock()
	}
}

func createSet() map[string]struct{} {
	return make(map[string]struct{}) // we use this as a set of string
}

func addToSet(set map[string]struct{}, key string) {
	if !setContains(set, key) {
		set[key] = struct{}{} //Add element to the set. We use struct{}{} as value in the map because it takes up 0 memory
	}
}

func setContains(set map[string]struct{}, key string) bool {
	_, exists := set[key]
	return exists
}

func (tb *ThinBroker) updateContextElement(ctxElem *ContextElement) {
	//look up who already subscribed to this context element
	eid := ctxElem.Entity.ID

	tb.entities_lock.Lock()
	defer tb.entities_lock.Unlock()

	// update its value in the entity map
	if curElement, exist := tb.entities[eid]; exist {
		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("curElement: ", curElement)
		}
		updateAttributes(&curElement.Attributes, ctxElem.Attributes)
		// updatedAttributeNames := createSet()
		// for _, attr := range ctxElem.Attributes {
		// 	if setContains(updatedAttributeNames, attr.Name) { //if already updated once
		// 		curElement.Attributes = append(curElement.Attributes, attr)
		// 	} else {
		// 		addToSet(updatedAttributeNames, attr.Name)
		// 		updateAttribute(&attr, curElement)
		// 	}
		// }
		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("curElement after updating attribute: ", curElement)
		}

		for _, metadata := range ctxElem.Metadata {
			updateDomainMetadata(&metadata, curElement)
		}
		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("tb.entities[", eid, "], ", tb.entities[eid])
		}
	} else {
		newContextElement := *ctxElem
		tb.entities[eid] = &newContextElement
		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("tb.entities[", eid, "], ", tb.entities[eid])
		}
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
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("ngsi9sub ", subscriptionId, " availabilitySubscription ", availabilitySubscription)
		DEBUG.Println(tb.entityId2Subcriptions)
	}
	if subscriptionId != "" {
		tb.subLinks_lock.Lock()
		tb.main2Other[sid] = append(tb.main2Other[sid], subscriptionId)
		tb.availabilitySub2MainSub[subscriptionId] = sid
		notifyMessage, alreadyBack := tb.tmpNGSI9NotifyCache[subscriptionId]
		tb.subLinks_lock.Unlock()
		if alreadyBack {
			INFO.Println("========forward the availability notify that arrive earlier===========")
			tb.handleNGSI9Notify(sid, notifyMessage, false)

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

func stringsContains(slice []string, e string) bool {
	for _, sliceElement := range slice {
		if sliceElement == e {
			return true
		}
	}
	return false
}

func (tb *ThinBroker) handleNGSI9Notify(mainSubID string, notifyContextAvailabilityReq *NotifyContextAvailabilityRequest, isProsumer bool) {

	var action string
	switch notifyContextAvailabilityReq.ErrorCode.Code {
	case 201:
		action = "CREATE"
	case 301:
		action = "UPDATE"
	case 410:
		action = "DELETE"
	default:
		action = "UPDATE" //often the status code is 200
	}

	var contextSubscription *SubscribeContextRequest
	if !isProsumer {
		tb.subscriptions_lock.RLock()
		contextSubscription = tb.subscriptions[mainSubID]
		tb.subscriptions_lock.RUnlock()

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println(action, " subID ", mainSubID, " subscription", contextSubscription, "isSimplyByType ", contextSubscription.IsSimplyByType())
			DEBUG.Println(tb.entityId2Subcriptions)
		}
	}

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println(action, " subID ", mainSubID, " subscription", contextSubscription)
		DEBUG.Println(tb.entityId2Subcriptions)
	}

	for _, registrationResp := range notifyContextAvailabilityReq.ContextRegistrationResponseList {
		registration := registrationResp.ContextRegistration

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("registration ", registration, "isProsumer", isProsumer)
		}

		// INFO.Println(registration.ProvidingApplication, ", ", tb.MyURL)
		// INFO.Println("TO ngsi10 subscription, ", mainSubID)
		// INFO.Printf("entity list: %+v\r\n", registration.EntityIdList)

		//
		// Here we send out the subscription to the context provider
		//
		if registration.ProvidingApplication == tb.MyURL {
			//for matched entities provided by myself
			if action == "CREATE" || action == "UPDATE" {
				tb.notifyOneSubscriberWithCurrentStatus(registration.EntityIdList, mainSubID)
			}
		} else {

			// Here we subscribe for data to the providing application

			// // this check is to subscribe to the data only for complex subscription (i.e., not simply by type)
			// if !contextSubscription.IsSimplyByType() || stringsContains(tb.entityId2Subcriptions["*"], mainSubID) {

			// 	//for matched entities provided by other IoT Brokers
			// 	newSubscription := SubscribeContextRequest{}
			// 	if contextSubscription.IsSimplyByType() {
			// 		// This loop is to only get pattern entity request (by Type) of the actual matching types
			// 		for _, subEntity := range contextSubscription.Entities {
			// 			for _, regEntity := range registration.EntityIdList {
			// 				if subEntity.Type == regEntity.Type {
			// 					newSubscription.Entities = append(newSubscription.Entities, subEntity)
			// 					break
			// 				}
			// 			}
			// 		}
			// 	} else {
			// 		newSubscription.Entities = registration.EntityIdList
			// 	}
			// 	newSubscription.Reference = tb.MyURL
			// 	newSubscription.Subscriber.BrokerURL = registration.ProvidingApplication

			// 	if action == "CREATE" || action == "UPDATE" {
			// 		sid, err := subscribeContextProvider(&newSubscription, registration.ProvidingApplication, tb.SecurityCfg)
			// 		if err == nil {
			// 			// INFO.Println("issue a new subscription ", sid)

			// 			tb.subscriptions_lock.Lock()
			// 			tb.subscriptions[sid] = &newSubscription
			// 			tb.subscriptions_lock.Unlock()

			// 			tb.subLinks_lock.Lock()
			// 			tb.main2Other[mainSubID] = append(tb.main2Other[mainSubID], sid)
			// 			tb.subLinks_lock.Unlock()
			// 		}
			// 	}
			// }

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("contextSubscription", contextSubscription, "| registration ", registration)
			}

			newSubscription := SubscribeContextRequest{}
			// this check is to subscribe to the data only for complex subscription (i.e., not simply by type)
			if !isProsumer && contextSubscription.IsSimplyByType() {

				// This loop is to only get pattern entity request (by Type) of the actual matching types
				for _, subEntity := range contextSubscription.Entities {
					for _, regEntity := range registration.EntityIdList {
						if LoggerIsEnabled(DEBUG) {
							DEBUG.Println("subEntity.Type", subEntity.Type, "| regEntity.Type ", regEntity.Type)
						}
						if subEntity.Type == regEntity.Type {
							if LoggerIsEnabled(DEBUG) {
								DEBUG.Println("regEntity.GetTypeWildCard() ", regEntity.GetTypeWildCard(), "tb.entityId2Subcriptions", tb.entityId2Subcriptions[regEntity.GetTypeWildCard()], "| mainSubID ", mainSubID)
							}
							// Let's make the subscription only if there was not another subscription by Type for the same type
							if !stringsContains(tb.entityId2Subcriptions[regEntity.GetTypeWildCard()], mainSubID) {
								newSubscription.Entities = append(newSubscription.Entities, subEntity)
							}
							break
						}
					}
				}

			} else {
				newSubscription.Entities = registration.EntityIdList
			}

			ngsi_version := checkNGSIversion(registration.Metadata)

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("newSubscription.Entities ", len(newSubscription.Entities), newSubscription.Entities, ngsi_version)
			}

			if len(newSubscription.Entities) > 0 {
				newSubscription.Reference = tb.MyURL
				newSubscription.Subscriber.BrokerURL = registration.ProvidingApplication

				if action == "CREATE" || action == "UPDATE" {
					// #########################
					// Send the subscription
					// ########################
					var sid string
					var err error
					if ngsi_version == "NGSI-LD" {
						sid, err = subscribeContextProviderNGSILD(&newSubscription, registration.ProvidingApplication, notifyContextAvailabilityReq.SubscriptionId, tb.SecurityCfg)
					} else {
						sid, err = subscribeContextProvider(&newSubscription, registration.ProvidingApplication, tb.SecurityCfg)
					}
					if err == nil {
						// INFO.Println("issue a new subscription ", sid)

						tb.subscriptions_lock.Lock()
						tb.subscriptions[sid] = &newSubscription
						tb.subscriptions_lock.Unlock()

						tb.subLinks_lock.Lock()
						tb.main2Other[mainSubID] = append(tb.main2Other[mainSubID], sid)
						tb.subLinks_lock.Unlock()

						if LoggerIsEnabled(DEBUG) {
							DEBUG.Println("tb.subscriptions ", tb.subscriptions)
							DEBUG.Println("tb.main2Other ", tb.main2Other)
						}
					}
				}
			}
		}

		//
		// In this loop we associate the entityIds (or a wildcard) to the subscription
		//
		for _, eid := range registration.EntityIdList {

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("===> ", eid, " , ", mainSubID)
			}

			tb.e2sub_lock.Lock()

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("action ", action)
			}

			if action == "CREATE" {
				if isProsumer {
					var id string
					if eid.IsSimplyByType() {
						id = eid.GetTypeWildCard()
					} else {
						id = eid.ID
					}
					if !stringsContains(tb.entityId2Subcriptions[id], mainSubID) {
						tb.entityId2Subcriptions[id] = append(tb.entityId2Subcriptions[id], mainSubID)
					}
				} else {
					if contextSubscription.IsSimplyByType() {
						wildCards := contextSubscription.GetTypeWildCards(&eid)
						for _, wildCard := range wildCards {
							if !stringsContains(tb.entityId2Subcriptions[wildCard], mainSubID) {
								tb.entityId2Subcriptions[wildCard] = append(tb.entityId2Subcriptions[wildCard], mainSubID)
							}
						}
					} else {
						tb.entityId2Subcriptions[eid.ID] = append(tb.entityId2Subcriptions[eid.ID], mainSubID)
					}
				}
			} else if action == "DELETE" {
				subList := tb.entityId2Subcriptions[eid.ID]
				for i, id := range subList {
					if id == mainSubID {
						tb.entityId2Subcriptions[eid.ID] = append(subList[:i], subList[i+1:]...)
						break
					}
				}
			} else if action == "UPDATE" {

				if isProsumer {

					var id string
					if eid.IsSimplyByType() {
						id = eid.GetTypeWildCard()
					} else {
						id = eid.ID
					}
					if !stringsContains(tb.entityId2Subcriptions[id], mainSubID) {
						tb.entityId2Subcriptions[id] = append(tb.entityId2Subcriptions[id], mainSubID)
					}
					// if eid.IsSimplyByType() {
					// 	if !stringsContains(tb.entityId2Subcriptions[eid.GetTypeWildCard()], mainSubID) {
					// 		tb.entityId2Subcriptions[eid.GetTypeWildCard()] = append(tb.entityId2Subcriptions[eid.GetTypeWildCard()], mainSubID)
					// 	}
					// } else {
					// 	if !stringsContains(tb.entityId2Subcriptions[eid.ID], mainSubID) {
					// 		tb.entityId2Subcriptions[eid.ID] = append(tb.entityId2Subcriptions[eid.ID], mainSubID)
					// 	}
					// }
				} else {
					if contextSubscription.IsSimplyByType() {
						wildCards := contextSubscription.GetTypeWildCards(&eid)
						if LoggerIsEnabled(DEBUG) {
							DEBUG.Println("wildCards:", wildCards, "contextSubscription:", contextSubscription, "eid:", eid)
						}
						for _, wildCard := range wildCards {
							if !stringsContains(tb.entityId2Subcriptions[wildCard], mainSubID) {
								tb.entityId2Subcriptions[wildCard] = append(tb.entityId2Subcriptions[wildCard], mainSubID)
							}
							if LoggerIsEnabled(DEBUG) {
								DEBUG.Println("wildCard", wildCard, "tb.entityId2Subcriptions[wildCard] ", tb.entityId2Subcriptions[wildCard], "mainSubID", mainSubID)
							}
						}
					} else {
						if !stringsContains(tb.entityId2Subcriptions[eid.ID], mainSubID) {
							tb.entityId2Subcriptions[eid.ID] = append(tb.entityId2Subcriptions[eid.ID], mainSubID)
						}
					}
				}
			}

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("tb.entityId2Subcriptions ", tb.entityId2Subcriptions)
			}

			tb.e2sub_lock.Unlock()
		}

	}
}

func checkNGSIversion(metadataList []ContextMetadata) string {

	for _, metadata := range metadataList {
		if strings.ToLower(metadata.Name) == "ngsiversion" || strings.ToLower(metadata.Type) == "ngsiversion" {
			if str, ok := metadata.Value.(string); ok {
				switch ver := strings.ToLower(str); ver {
				case "ngsild", "ngsi-ld":
					return "NGSI-LD"
				case "ngsiv1", "ngsi":
					return "NGSIv1"
				default:
					return "NGSIv1"
				}
			}
		}
	}
	return "NGSIv1"
}

const SUBID_INTERNAL_PROSUMER_SUFFIX = ":ffinternal"

func genProsumerSubID(subID string) string {
	return subID + SUBID_INTERNAL_PROSUMER_SUFFIX
}

func isProsumerSubscription(correlator string, subID string) bool {
	return strings.Replace(correlator, SUBID_INTERNAL_PROSUMER_SUFFIX, "", 1) == strings.Replace(subID, SUBID_INTERNAL_PROSUMER_SUFFIX, "", 1)
}

func (tb *ThinBroker) handleProsumerRegistration(notifyAvail *NotifyContextAvailabilityRequest) {

	subID := notifyAvail.SubscriptionId

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("Handle prosumer subscription ", subID)
	}

	//
	// First subscribe to myself
	//
	subIDInternal := genProsumerSubID(subID)
	for _, registration := range notifyAvail.ContextRegistrationResponseList {

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println("First subscribe to myself", subIDInternal)
		}

		var subscribeRequest SubscribeContextRequest
		subscribeRequest.Entities = registration.ContextRegistration.EntityIdList
		subscribeRequest.Reference = registration.ContextRegistration.ProvidingApplication

		for _, metadata := range registration.ContextRegistration.Metadata {
			switch name := metadata.Name; name {
			case "ngsiversion":
				value, ok := metadata.Value.(string)
				if ok {
					subscribeRequest.Subscriber.DestinationType = value
				}
			case "Ngsild-Tenant":
				value, ok := metadata.Value.(string)
				if ok {
					subscribeRequest.Subscriber.Tenant = value
				}
			case "Fiware-Correlator":
				value, ok := metadata.Value.(string)
				if ok {
					subscribeRequest.Subscriber.Correlator = value
				}
			case "Require-Reliability":
				value, ok := metadata.Value.(bool)
				if ok {
					subscribeRequest.Subscriber.RequireReliability = value
					subscribeRequest.Subscriber.NotifyCache = make([]*ContextElement, 0)
				}
			default:
				continue
			}

			// if r.Header.Get("User-Agent") == "lightweight-iot-broker" {
			// 	subReq.Subscriber.IsInternal = true
			// } else {
			// 	subReq.Subscriber.IsInternal = false
			// }

			//subscribeRequest.Subscriber.BrokerURL = registration.ContextRegistration.ProvidingApplication
			subscribeRequest.Subscriber.BrokerURL = tb.MyURL
			subscribeRequest.Subscriber.IsInternal = true

			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("SubscribeContext ", subIDInternal, " to myself for the prosumer ", subscribeRequest)
			}

			// Here I subscribe to myself on behalf of the providingApplication
			tb.subscribeContext(&subscribeRequest, subIDInternal)

		}
	}

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("Then subscribe to the providing application", subID)
	}
	//
	// Then subscribe to the providing application
	//
	tb.handleNGSI9Notify(subID, notifyAvail, true)

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

// func (tb *ThinBroker) handleStream2Self(notifyContextAvailabilityReq *NotifyContextAvailabilityRequest) {
// 	for _, ctxRegResp := range notifyContextAvailabilityReq.ContextRegistrationResponseList{
// 		ctxReg := ctxRegResp.ContextRegistration

// 	}

// }

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
