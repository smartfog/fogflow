package main

import (
	"sync"

	. "github.com/smartfog/fogflow/common/ngsi"
)

type EntityRepository struct {
	// cache all received registration in the memory for the performance reason
	ctxRegistrationList      map[string]*ContextRegistration
	ctxRegistrationList_lock sync.RWMutex

	// lock to control the update of database
	dbLock sync.RWMutex
}

func (er *EntityRepository) Init() {
	// initialize the registration list
	er.ctxRegistrationList = make(map[string]*ContextRegistration)
}

//
// update the registration in the repository and also
// return a flag to indicate if there is anything in the repository before
//
func (er *EntityRepository) updateEntity(entity EntityId, registration *ContextRegistration) *ContextRegistration {
	updatedRegistration := er.updateRegistrationInMemory(entity, registration)

	// return the latest view of the registration for this entity
	return updatedRegistration
}

//
// for the performance purpose, we still keep the latest view of all registrations
//
func (er *EntityRepository) updateRegistrationInMemory(entity EntityId, registration *ContextRegistration) *ContextRegistration {
	er.ctxRegistrationList_lock.Lock()
	defer er.ctxRegistrationList_lock.Unlock()

	eid := entity.ID

	if existRegistration, exist := er.ctxRegistrationList[eid]; exist {
		// if the registration already exists, update it with the received update

		// update attribute table
		for _, attr := range registration.ContextRegistrationAttributes {
			for i, existAttr := range existRegistration.ContextRegistrationAttributes {
				if existAttr.Name == attr.Name {
					// remove the old one
					existRegistration.ContextRegistrationAttributes = append(existRegistration.ContextRegistrationAttributes[:i], existRegistration.ContextRegistrationAttributes[i+1:]...)
					break
				}
			}
			// append the new one
			existRegistration.ContextRegistrationAttributes = append(existRegistration.ContextRegistrationAttributes, attr)
		}

		// update metadata table
		for _, meta := range registration.Metadata {
			for i, existMeta := range existRegistration.Metadata {
				if existMeta.Name == meta.Name {
					// remove the old one
					existRegistration.Metadata = append(existRegistration.Metadata[:i], existRegistration.Metadata[i+1:]...)
					break
				}
			}
			// append the new one
			existRegistration.Metadata = append(existRegistration.Metadata, meta)
		}

		// update the provided URL
		if len(registration.ProvidingApplication) > 0 {
			existRegistration.ProvidingApplication = registration.ProvidingApplication
		}
	} else {
		er.ctxRegistrationList[eid] = registration
	}

	return er.ctxRegistrationList[eid]
}

func (er *EntityRepository) queryEntities(entities []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	return er.queryEntitiesInMemory(entities, attributes, restriction)
}

func (er *EntityRepository) queryEntitiesInMemory(entities []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	entityMap := make(map[string][]EntityId)

	for _, registration := range er.ctxRegistrationList {
		entities := matchingWithFilters(registration, entities, attributes, restriction)
		if len(entities) > 0 {
			providerURL := registration.ProvidingApplication
			entityMap[providerURL] = append(entityMap[providerURL], entities...)
		}
	}

	return entityMap
}

func (er *EntityRepository) deleteEntity(eid string) {
	er.ctxRegistrationList_lock.Lock()
	delete(er.ctxRegistrationList, eid)
	er.ctxRegistrationList_lock.Unlock()
}

func (er *EntityRepository) ProviderLeft(providerURL string) {
	er.ctxRegistrationList_lock.Lock()
	for eid, registration := range er.ctxRegistrationList {
		if registration.ProvidingApplication == providerURL {
			delete(er.ctxRegistrationList, eid)
		}
	}
	er.ctxRegistrationList_lock.Unlock()
}

func (er *EntityRepository) retrieveRegistration(entityID string) *ContextRegistration {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	return er.ctxRegistrationList[entityID]
}
