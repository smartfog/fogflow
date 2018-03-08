package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"

	. "fogflow/common/ngsi"
)

type FastDiscovery struct {
	//backend entity repository
	repository EntityRepository

	//mapping from subscriptionID to subscription
	subscriptions      map[string]*SubscribeContextAvailabilityRequest
	subscriptions_lock sync.RWMutex
}

func (fd *FastDiscovery) Init(cfg *DatabaseCfg) {
	fd.subscriptions = make(map[string]*SubscribeContextAvailabilityRequest)
	fd.repository.Init(cfg)
}

func (fd *FastDiscovery) Stop() {
	fd.repository.Close()
}

func (fd *FastDiscovery) RegisterContext(w rest.ResponseWriter, r *rest.Request) {
	registerCtxReq := RegisterContextRequest{}
	err := r.DecodeJsonPayload(&registerCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if registerCtxReq.RegistrationId == "" {
		u1, err := uuid.NewV4()
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		registrationID := u1.String()
		registerCtxReq.RegistrationId = registrationID
	}

	// update context registration
	go fd.updateRegistration(&registerCtxReq)

	// send out the response
	registerCtxResp := RegisterContextResponse{}
	registerCtxResp.RegistrationId = registerCtxReq.RegistrationId
	registerCtxResp.Duration = registerCtxReq.Duration
	registerCtxResp.ErrorCode.Code = 200
	registerCtxResp.ErrorCode.ReasonPhrase = "OK"
	w.WriteJson(&registerCtxResp)
}

func (fd *FastDiscovery) notifySubscribers(registration *ContextRegistration, updateAction string) {
	fd.subscriptions_lock.RLock()
	defer fd.subscriptions_lock.RUnlock()

	providerURL := registration.ProvidingApplication
	for _, subscription := range fd.subscriptions {
		// find out the updated entities matched with this subscription
		entities := fd.matchingWithSubscription(registration, subscription)
		if len(entities) == 0 {
			continue
		}

		subscriberURL := subscription.Reference
		subID := subscription.SubscriptionId

		entityMap := make(map[string][]EntityId)
		entityMap[providerURL] = entities

		// send out AvailabilityNotify to subscribers
		go fd.sendNotify(subID, subscriberURL, entityMap, updateAction)
	}
}

func (fd *FastDiscovery) matchingWithSubscription(registration *ContextRegistration, subscription *SubscribeContextAvailabilityRequest) []EntityId {
	matchedEntities := make([]EntityId, 0)

	for _, entity := range registration.EntityIdList {
		// check entityId part
		atLeastOneMatched := false
		for _, tmp := range subscription.Entities {
			matched := matchEntityId(entity, tmp)
			if matched == true {
				atLeastOneMatched = true
				break
			}
		}
		if atLeastOneMatched == false {
			continue
		}

		// check attribute set
		matched := matchAttributes(registration.ContextRegistrationAttributes, subscription.Attributes)
		if matched == false {
			continue
		}

		// check metadata set
		matched = matchMetadatas(registration.Metadata, subscription.Restriction)
		if matched == false {
			continue
		}

		// if matched, add it into the list
		if matched == true {
			matchedEntities = append(matchedEntities, entity)
		}
	}

	return matchedEntities
}

func (fd *FastDiscovery) updateRegistration(registReq *RegisterContextRequest) {
	for _, registration := range registReq.ContextRegistrations {
		for _, entity := range registration.EntityIdList {

			INFO.Printf("registration:%+v\r\n", entity)

			fd.repository.updateEntity(entity, &registration)

			// inform the associated subscribers after updating the repository
			updatedRegistration := fd.repository.retrieveRegistration(entity.ID)
			if updatedRegistration != nil {
				fd.notifySubscribers(updatedRegistration, "UPDATE")
			}
		}
	}
}

func (fd *FastDiscovery) deleteRegistration(registration *ContextRegistration) {
	for _, entity := range registration.EntityIdList {
		// notify the affected subscribers before deleting the entity
		registration := fd.repository.retrieveRegistration(entity.ID)
		if registration != nil {
			fd.notifySubscribers(registration, "DELETE")
		}

		fd.repository.deleteEntity(entity.ID)
	}
}

func (fd *FastDiscovery) DiscoverContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	discoverCtxReq := DiscoverContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&discoverCtxReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// query database to get the result
	result := fd.handleQueryCtxAvailability(&discoverCtxReq)

	// send out the response
	discoverCtxResp := DiscoverContextAvailabilityResponse{}
	if result == nil {
		discoverCtxResp.ErrorCode.Code = 500
		discoverCtxResp.ErrorCode.ReasonPhrase = "database is too overloaded"
	} else {
		discoverCtxResp.ContextRegistrationResponses = *result
		discoverCtxResp.ErrorCode.Code = 200
		discoverCtxResp.ErrorCode.ReasonPhrase = "OK"
	}
	w.WriteJson(&discoverCtxResp)
}

func (fd *FastDiscovery) handleQueryCtxAvailability(req *DiscoverContextAvailabilityRequest) *[]ContextRegistrationResponse {
	//fmt.Println("************** query availability ****************")
	entityMap := fd.repository.queryEntities(req.Entities, req.Attributes, req.Restriction)

	// prepare the response
	registrationList := make([]ContextRegistrationResponse, 0)

	for url, entity := range entityMap {
		resp := ContextRegistrationResponse{}
		resp.ContextRegistration.ProvidingApplication = url
		resp.ContextRegistration.EntityIdList = entity

		resp.ErrorCode.Code = 200
		resp.ErrorCode.ReasonPhrase = "OK"

		registrationList = append(registrationList, resp)
	}

	return &registrationList
}

func (fd *FastDiscovery) SubscribeContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	subscribeCtxAvailabilityReq := SubscribeContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&subscribeCtxAvailabilityReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate a new subscription id
	u1, err := uuid.NewV4()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subID := u1.String()

	subscribeCtxAvailabilityReq.SubscriptionId = subID

	// add the new subscription
	fd.subscriptions_lock.Lock()
	fd.subscriptions[subID] = &subscribeCtxAvailabilityReq
	fd.subscriptions_lock.Unlock()

	// send out the response
	subscribeCtxAvailabilityResp := SubscribeContextAvailabilityResponse{}
	subscribeCtxAvailabilityResp.SubscriptionId = subID
	subscribeCtxAvailabilityResp.Duration = subscribeCtxAvailabilityReq.Duration
	subscribeCtxAvailabilityResp.ErrorCode.Code = 200
	subscribeCtxAvailabilityResp.ErrorCode.ReasonPhrase = "OK"

	w.WriteJson(&subscribeCtxAvailabilityResp)

	// trigger the process to send out the matched context availability infomation to the subscriber
	go fd.handleSubscrieCtxAvailability(&subscribeCtxAvailabilityReq)
}

func (fd *FastDiscovery) handleSubscrieCtxAvailability(subReq *SubscribeContextAvailabilityRequest) {
	entityMap := fd.repository.queryEntities(subReq.Entities, subReq.Attributes, subReq.Restriction)

	if len(entityMap) > 0 {
		fd.sendNotify(subReq.SubscriptionId, subReq.Reference, entityMap, "CREATE")
	}
}

func (fd *FastDiscovery) sendNotify(subID string, subscriberURL string, entityMap map[string][]EntityId, action string) {
	notifyReq := NotifyContextAvailabilityRequest{}
	notifyReq.SubscriptionId = subID

	// carry the actions via the code number
	switch action {
	case "CREATE":
		notifyReq.ErrorCode.Code = 201
	case "UPDATE":
		notifyReq.ErrorCode.Code = 301
	case "DELETE":
		notifyReq.ErrorCode.Code = 410
	}

	notifyReq.ErrorCode.ReasonPhrase = "OK"

	// prepare the response
	registrationList := make([]ContextRegistrationResponse, 0)

	for url, entity := range entityMap {
		resp := ContextRegistrationResponse{}
		resp.ContextRegistration.ProvidingApplication = url
		resp.ContextRegistration.EntityIdList = entity

		resp.ErrorCode.Code = 200

		resp.ErrorCode.ReasonPhrase = "OK"

		registrationList = append(registrationList, resp)
	}

	notifyReq.ContextRegistrationResponseList = registrationList

	body, err := json.Marshal(notifyReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(string(body))
	//fmt.Println("send to subscriber at ", subscriberURL)

	req, err := http.NewRequest("POST", subscriberURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err2 := client.Do(req)
	defer resp.Body.Close()
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	text, _ := ioutil.ReadAll(resp.Body)

	notifyCtxAvailResp := NotifyContextAvailabilityResponse{}
	err = json.Unmarshal(text, &notifyCtxAvailResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	if notifyCtxAvailResp.ResponseCode.Code != 200 {
		fmt.Println(notifyCtxAvailResp.ResponseCode.ReasonPhrase)
	}
}

func (fd *FastDiscovery) UnsubscribeContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	unsubscribeCtxAvailabilityReq := UnsubscribeContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&unsubscribeCtxAvailabilityReq)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subID := unsubscribeCtxAvailabilityReq.SubscriptionId

	fmt.Println("unsubscribe context availability, ", subID)

	// remove the subscription
	fd.subscriptions_lock.Lock()
	delete(fd.subscriptions, subID)
	fd.subscriptions_lock.Unlock()

	// send out the response
	unsubscribeCtxAvailabilityResp := UnsubscribeContextAvailabilityResponse{}
	unsubscribeCtxAvailabilityResp.SubscriptionId = unsubscribeCtxAvailabilityReq.SubscriptionId
	unsubscribeCtxAvailabilityResp.StatusCode.Code = 200
	unsubscribeCtxAvailabilityResp.StatusCode.Details = "OK"

	w.WriteJson(&unsubscribeCtxAvailabilityResp)
}

func (fd *FastDiscovery) getRegisteredEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	registration := fd.repository.retrieveRegistration(eid)
	w.WriteJson(registration)
}

func (fd *FastDiscovery) deleteRegisteredEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	DEBUG.Printf("delete the context availability %s\r\n", eid)

	registration := fd.repository.retrieveRegistration(eid)
	if registration != nil {
		fd.deleteRegistration(registration)
	}

	w.WriteHeader(200)
}

func (fd *FastDiscovery) getSubscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	fd.subscriptions_lock.RLocker()
	defer fd.subscriptions_lock.RUnlock()

	subscription := fd.subscriptions[sid]
	w.WriteJson(subscription)
}

func (fd *FastDiscovery) getSubscriptions(w rest.ResponseWriter, r *rest.Request) {
	fd.subscriptions_lock.RLock()
	defer fd.subscriptions_lock.RUnlock()

	w.WriteJson(fd.subscriptions)
}

func (fd *FastDiscovery) getStatus(w rest.ResponseWriter, r *rest.Request) {
	w.WriteHeader(200)
}
