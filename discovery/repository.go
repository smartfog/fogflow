package main

import (
	. "github.com/smartfog/fogflow/common/ngsi"
	"sort"
	"strings"
	"sync"
)

type Candidate struct {
	ProviderURL string
	ID          string
	Type        string
	Distance    uint64
}

type EntityRepository struct {
	// cache all received registration in the memory for the performance reason
	//ctxRegistrationList      map[string]*ContextRegistration
	ctxRegistrationList      map[string]*EntityRegistration
	ctxRegistrationList_lock sync.RWMutex

	// lock to control the update of database
	dbLock sync.RWMutex
}

func (er *EntityRepository) Init() {
	// initialize the registration list
	er.ctxRegistrationList = make(map[string]*EntityRegistration)
}

//
// update the registration in the repository and also
// return a flag to indicate if there is anything in the repository before
//
func (er *EntityRepository) updateEntity(entity EntityId, registration *ContextRegistration) *EntityRegistration {
	updatedRegistration := er.updateRegistrationInMemory(entity, registration)

	// return the latest view of the registration for this entity
	return updatedRegistration
}

//
// for the performance purpose, we still keep the latest view of all registrations
//
func (er *EntityRepository) updateRegistrationInMemory(entity EntityId, registration *ContextRegistration) *EntityRegistration {
	er.ctxRegistrationList_lock.Lock()
	defer er.ctxRegistrationList_lock.Unlock()

	var eid string
	if entity.ID != "" {
		eid = entity.ID
	} else if entity.IdPattern != "" { // NGSI-LD feature update, support for registration based on IdPattern added.
		eid = entity.IdPattern
	}

	if existRegistration, exist := er.ctxRegistrationList[eid]; exist {
		// update existing entity type
		if entity.Type != "" {
			existRegistration.Type = entity.Type
		}

		// update existing attribute table
		for _, attr := range registration.ContextRegistrationAttributes {
			existRegistration.AttributesList[attr.Name] = attr
		}

		// update existing metadata table
		for _, meta := range registration.Metadata {
			existRegistration.MetadataList[meta.Name] = meta
		}

		// update existing providerURL
		if len(registration.ProvidingApplication) > 0 {
			existRegistration.ProvidingApplication = registration.ProvidingApplication
		}
	} else {
		entityRegistry := EntityRegistration{}

		entityRegistry.ID = eid
		entityRegistry.Type = entity.Type

		entityRegistry.AttributesList = make(map[string]ContextRegistrationAttribute)
		entityRegistry.MetadataList = make(map[string]ContextMetadata)

		for _, attr := range registration.ContextRegistrationAttributes {
			entityRegistry.AttributesList[attr.Name] = attr
		}

		// update existing metadata table
		for _, meta := range registration.Metadata {
			entityRegistry.MetadataList[meta.Name] = meta
		}

		// update existing providerURL
		if len(registration.ProvidingApplication) > 0 {
			entityRegistry.ProvidingApplication = registration.ProvidingApplication
		}

		er.ctxRegistrationList[eid] = &entityRegistry
	}

	return er.ctxRegistrationList[eid]
}

func (er *EntityRepository) queryEntities(entities []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	return er.queryEntitiesInMemory(entities, attributes, restriction)
}

func (er *EntityRepository) queryEntitiesInMemory(entities []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	nearby := restriction.GetNearbyFilter()

	candidates := make([]Candidate, 0)
	for _, registration := range er.ctxRegistrationList {
		if matchingWithFilters(registration, entities, attributes, restriction) == true {
			candidate := Candidate{}
			candidate.ID = registration.ID
			candidate.Type = registration.Type
			candidate.ProviderURL = registration.ProvidingApplication

			if nearby != nil {
				landmark := Point{}
				landmark.Longitude = nearby.Longitude
				landmark.Latitude = nearby.Latitude

				location := registration.GetLocation()

				candidate.Distance = Distance(&location, &landmark)
			}

			candidates = append(candidates, candidate)
		}
	}
	if nearby != nil {
		if len(candidates) > nearby.Limit {
			// for the nearby query, just select the closest n matched entities
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].Distance < candidates[j].Distance
			})

			candidates = candidates[0:nearby.Limit]
		}

		DEBUG.Println("number of returned entities: ", nearby.Limit)
	}

	// return the final result
	entityMap := make(map[string][]EntityId, 0)

	for _, candidate := range candidates {
		entity := EntityId{}
		entity.ID = candidate.ID
		entity.Type = candidate.Type
		entity.IsPattern = false

		providerURL := candidate.ProviderURL
		entityMap[providerURL] = append(entityMap[providerURL], entity)
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

//NGSI-LD feature update, registration can now be retrieved based on IdPattern.
func (er *EntityRepository) retrieveRegistration(entityID string) *EntityRegistration {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	if _, ok := er.ctxRegistrationList[entityID]; ok == true {
		return er.ctxRegistrationList[entityID]
	} else {
		for id, registration := range er.ctxRegistrationList {
			if strings.Contains(id, ".*") && strings.Contains(id, "*.") {
				id = strings.Trim(id, ".*")
				id = strings.Trim(id, "*.")
				if strings.Contains(entityID, id) {
					return registration
				}

			} else if strings.Contains(id, ".*") {
				id = strings.Trim(id, ".*")
				if strings.HasPrefix(entityID, id) {
					return registration
				}

			} else if strings.Contains(id, "*.") {
				id = strings.Trim(id, "*.")
				if strings.HasSuffix(entityID, id) {
					return registration
				}
			}
		}
	}
	return nil
}
