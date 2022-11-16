package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	. "fogflow/common/ngsi"
)

// ============= update ====================

func (tb *ThinBroker) NGSILD_UpdateContext(w rest.ResponseWriter, r *rest.Request) {
	ngsildUpsert := []map[string]interface{}{}

	err := r.DecodeJsonPayload(&ngsildUpsert)
	if err != nil {
		DEBUG.Println("not able to decode the received message")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateCtxReq := UpdateContextRequest{}
	numUpdates := updateCtxReq.ReadFromNGSILD(ngsildUpsert)

	if numUpdates > 0 {
		tb.handleInternalUpdateContext(&updateCtxReq)
	}
}

func (tb *ThinBroker) NGSILD_CreateEntity(w rest.ResponseWriter, r *rest.Request) {
	entity := map[string]interface{}{}

	err := r.DecodeJsonPayload(&entity)
	if err != nil {
		DEBUG.Println("not able to decode the received message")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ngsildUpsert := make([]map[string]interface{}, 0)
	ngsildUpsert = append(ngsildUpsert, entity)

	updateCtxReq := UpdateContextRequest{}
	numUpdates := updateCtxReq.ReadFromNGSILD(ngsildUpsert)

	if numUpdates > 0 {
		tb.handleInternalUpdateContext(&updateCtxReq)
	}
}

func (tb *ThinBroker) NGSILD_DeleteEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	err := tb.deleteEntity(eid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

func (tb *ThinBroker) NGSILD_DeleteAttribute(w rest.ResponseWriter, r *rest.Request) {

}

// ============= query ====================

func (tb *ThinBroker) NGSILD_QueryByPostedFilters(w rest.ResponseWriter, r *rest.Request) {

}

func (tb *ThinBroker) NGSILD_QueryById(w rest.ResponseWriter, r *rest.Request) {
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

	// send out the matched context elements in the form of NGSI-LD
	if len(matchedCtxElement) > 0 {
		w.WriteJson(toNGSILDPayload(matchedCtxElement)[0])
	} else {
		message := map[string]interface{}{
			"error":       "NotFound",
			"description": "The requested entity has not been found. Check type and id",
		}
		w.WriteJson(message)
	}
}

func (tb *ThinBroker) NGSILD_QueryByParameters(w rest.ResponseWriter, r *rest.Request) {

}

// ============= subscribe and notify ====================

func (tb *ThinBroker) NGSILD_SubcribeContext(w rest.ResponseWriter, r *rest.Request) {

}

func (tb *ThinBroker) NGSILD_NotifyContext(w rest.ResponseWriter, r *rest.Request) {
	msg := map[string]interface{}{}

	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		DEBUG.Println("not able to decode the received message")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data, exist := msg["data"]; exist {
		updateList := data.([]interface{})

		ngsildUpsert := make([]map[string]interface{}, 0)
		for _, update := range updateList {
			entity := update.(map[string]interface{})
			ngsildUpsert = append(ngsildUpsert, entity)
		}

		updateCtxReq := UpdateContextRequest{}
		numUpdates := updateCtxReq.ReadFromNGSILD(ngsildUpsert)

		if numUpdates > 0 {
			tb.handleInternalUpdateContext(&updateCtxReq)
		}
	}
}

func (tb *ThinBroker) NGSILD_UnsubscribeLDContext(w rest.ResponseWriter, r *rest.Request) {

}
