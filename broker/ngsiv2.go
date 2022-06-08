package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	. "fogflow/common/ngsi"
)

func (tb *ThinBroker) NGSIV2_getEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	// construct a query request based on the received parameter
	queryCtxReq := QueryContextRequest{}

	entityId := EntityId{}
	entityId.ID = eid
	entityId.IsPattern = false

	queryCtxReq.Entities = make([]EntityId, 0)
	queryCtxReq.Entities = append(queryCtxReq.Entities, entityId)

	// discover the availability of all matched entities
	entityMap := tb.discoveryEntities(queryCtxReq.Entities, queryCtxReq.Attributes, queryCtxReq.Restriction)

	matchedCtxElement := make([]ContextElement, 0)

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

	// send out the matched context elements in the form of NGSIv2
	if len(matchedCtxElement) > 0 {
		w.WriteJson(toNGSIv2Payload(matchedCtxElement)[0])
	} else {
		message := map[string]interface{}{
			"error":       "NotFound",
			"description": "The requested entity has not been found. Check type and id",
		}
		w.WriteJson(message)
	}
}

func (tb *ThinBroker) NGSIV2_getEntities(w rest.ResponseWriter, r *rest.Request) {
	queryCtxReq := QueryContextRequest{}

	entityId := EntityId{}
	entityId.ID = "*.*"
	entityId.IsPattern = true

	queryCtxReq.Entities = make([]EntityId, 0)
	queryCtxReq.Entities = append(queryCtxReq.Entities, entityId)

	// discover the availability of all matched entities
	entityMap := tb.discoveryEntities(queryCtxReq.Entities, queryCtxReq.Attributes, queryCtxReq.Restriction)

	matchedCtxElement := make([]ContextElement, 0)

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

	// send out the matched context elements in the form of NGSIv2
	w.WriteJson(toNGSIv2Payload(matchedCtxElement))
}

func (tb *ThinBroker) NGSIV2_queryEntities(w rest.ResponseWriter, r *rest.Request) {
	queryParameters := r.URL.Query()

	// construct a query request based on the received parameter
	queryCtxReq := QueryContextRequest{}

	queryCtxReq.Entities = make([]EntityId, 0)

	entityType := queryParameters.Get("type")
	if len(entityType) > 0 {
		entityId := EntityId{}
		entityId.Type = entityType
		entityId.IsPattern = true
		queryCtxReq.Entities = append(queryCtxReq.Entities, entityId)
	}

	matchedCtxElement := make([]ContextElement, 0)

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

	// send out the matched context elements in the form of NGSIv2
	w.WriteJson(toNGSIv2Payload(matchedCtxElement))
}

func (tb *ThinBroker) NGSIV2_createEntities(w rest.ResponseWriter, r *rest.Request) {
	ngsiv2Upsert := map[string]interface{}{}

	err := r.DecodeJsonPayload(&ngsiv2Upsert)
	if err != nil {
		DEBUG.Println("not able to decode the received message")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateCtxReq := UpdateContextRequest{}
	numUpdates := updateCtxReq.ReadFromNGSIv2(ngsiv2Upsert)

	if numUpdates > 0 {
		tb.handleInternalUpdateContext(&updateCtxReq)
	}
}

func (tb *ThinBroker) NGSIV2_notify(w rest.ResponseWriter, r *rest.Request) {
	msg := map[string]interface{}{}

	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		DEBUG.Println("not able to decode the received message")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data, exist := msg["data"]; exist {
		updateList := data.([]interface{})

		for _, update := range updateList {
			entity := update.(map[string]interface{})

			updateCtxReq := UpdateContextRequest{}
			numUpdates := updateCtxReq.ReadFromNGSIv2(entity)
			if numUpdates > 0 {
				tb.handleInternalUpdateContext(&updateCtxReq)
			}
		}
	}
}
