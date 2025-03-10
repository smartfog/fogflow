package main

import (
	"encoding/json"
	. "fogflow/common/config"
	. "fogflow/common/ngsi"
	"io/ioutil"
	"sort"
	"sync"
	"time"
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
	dbLock                        sync.RWMutex
	storeToFileScheduled          bool
	delayStoreRegistrationsOnFile int
	registrationsOnDisk           bool

	dbFile string
}

func (er *EntityRepository) Init(config *Config) {
	// initialize the registration list
	er.ctxRegistrationList = make(map[string]*EntityRegistration)

	er.storeToFileScheduled = false
	er.registrationsOnDisk = config.Discovery.StoreOnDisk
	//INFO.Println("config.Discovery.DelayStoreRegistrationsOnFile ", config.Discovery.DelayStoreStoreOnFile)
	er.delayStoreRegistrationsOnFile = config.Discovery.DelayStoreOnFile

	er.dbFile = "discoveryDB/registrations.json"

	if er.registrationsOnDisk {
		er.readRegistrationsFromDisk()
	}
}

// update the registration in the repository and also
// return a flag to indicate if there is anything in the repository before
func (er *EntityRepository) updateEntity(entity EntityId, registration *ContextRegistration) *EntityRegistration {
	updatedRegistration := er.updateRegistrationInMemory(entity, registration)

	// If requested, update registrations on disk
	if er.registrationsOnDisk {
		er.updateRegistrationsOnDisk()
	}

	// return the latest view of the registration for this entity
	return updatedRegistration
}

// return all available entity types
func (er *EntityRepository) GetEntityTypes() []string {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	typeMap := make(map[string]bool)
	for _, registration := range er.ctxRegistrationList {
		etype := registration.Type
		if _, found := typeMap[etype]; !found {
			typeMap[etype] = true
		}
	}

	typeList := make([]string, 0, len(typeMap))
	for k := range typeMap {
		typeList = append(typeList, k)
	}

	return typeList
}

// for the performance purpose, we still keep the latest view of all registrations
func (er *EntityRepository) updateRegistrationInMemory(entity EntityId, registration *ContextRegistration) *EntityRegistration {
	er.ctxRegistrationList_lock.Lock()
	defer er.ctxRegistrationList_lock.Unlock()
	eid := entity.ID
	if existRegistration, exist := er.ctxRegistrationList[eid]; exist {
		// update existing entity type
		if entity.Type != "" {
			existRegistration.Type = entity.Type
		}

		attrilist := make(map[string]ContextRegistrationAttribute)
		// update existing attribute table
		for _, attr := range registration.ContextRegistrationAttributes {
			//existRegistration.AttributesList[:0]
			existRegistration.AttributesList[attr.Name] = attr
			attrilist[attr.Name] = attr
		}

		for _, attributeOld := range existRegistration.AttributesList {
			found := false
			for _, attributeNew := range attrilist {
				if attributeNew.Name == attributeOld.Name {
					found = true
					break
				}
			}
			if !found {
				delete(existRegistration.AttributesList, attributeOld.Name)
			}
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
		entityRegistration := EntityRegistration{}

		entityRegistration.ID = eid
		entityRegistration.Type = entity.Type

		entityRegistration.AttributesList = make(map[string]ContextRegistrationAttribute)
		entityRegistration.MetadataList = make(map[string]ContextMetadata)

		for _, attr := range registration.ContextRegistrationAttributes {
			entityRegistration.AttributesList[attr.Name] = attr
		}

		// update existing metadata table
		for _, meta := range registration.Metadata {
			entityRegistration.MetadataList[meta.Name] = meta
		}

		// update existing providerURL
		if len(registration.ProvidingApplication) > 0 {
			entityRegistration.ProvidingApplication = registration.ProvidingApplication
		}

		er.ctxRegistrationList[eid] = &entityRegistration
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
		if matchingWithFilters(registration, entities, attributes, restriction) {
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

		providerURL := candidate.ProviderURL
		entityMap[providerURL] = append(entityMap[providerURL], entity)
	}

	return entityMap
}

func (er *EntityRepository) deleteEntity(eid string) {
	er.ctxRegistrationList_lock.Lock()
	delete(er.ctxRegistrationList, eid)
	er.ctxRegistrationList_lock.Unlock()

	// If requested, update registrations on disk
	if er.registrationsOnDisk {
		er.updateRegistrationsOnDisk()
	}

}

func (er *EntityRepository) ProviderLeft(providerURL string) {
	er.ctxRegistrationList_lock.Lock()
	for eid, registration := range er.ctxRegistrationList {
		if registration.ProvidingApplication == providerURL {
			delete(er.ctxRegistrationList, eid)
		}
	}
	er.ctxRegistrationList_lock.Unlock()

	// If requested, update registrations on disk
	if er.registrationsOnDisk {
		er.updateRegistrationsOnDisk()
	}
}

func (er *EntityRepository) retrieveRegistration(entityID string) *EntityRegistration {
	er.ctxRegistrationList_lock.RLock()
	defer er.ctxRegistrationList_lock.RUnlock()

	return er.ctxRegistrationList[entityID]
}

func (er *EntityRepository) updateRegistrationsOnDisk() {

	// This initial code is to avoid to write on file for every new subscriptions if
	// they are arriving too close with each other in terms of time
	// So here we check if the dblock was already taken without trying to take it
	// if it is already taken then just return because the new subscription will be written
	// by the already scheduled write to file
	// If the lock is not taken, then take it and wait for 3 seconds.
	// Then for the next 3 seconds other pursuer of storing on file will see that
	// somebody has the lock

	if er.storeToFileScheduled {
		INFO.Println("A store on file for registrations is already scheduled")
		return
	}

	er.dbLock.Lock()

	er.storeToFileScheduled = true

	time.Sleep(time.Duration(er.delayStoreRegistrationsOnFile) * time.Second)

	er.ctxRegistrationList_lock.RLock()

	INFO.Println("Writing registration into file")

	defer er.ctxRegistrationList_lock.RUnlock()
	defer er.dbLock.Unlock()

	//...................................
	//Writing struct type to a JSON file
	//...................................
	content, err := json.Marshal(er.ctxRegistrationList)
	if err != nil {
		ERROR.Println(err)
	}
	// err = ioutil.WriteFile("registrations.json", content, 0644)
	err = ioutil.WriteFile(er.dbFile, content, 0644)
	if err != nil {
		ERROR.Println(err)
	}

	er.storeToFileScheduled = false

}

func (er *EntityRepository) readRegistrationsFromDisk() {

	INFO.Println("Reading registrations from file")

	er.dbLock.Lock()
	defer er.dbLock.Unlock()

	// content, err := ioutil.ReadFile("registrations.json")
	content, err := ioutil.ReadFile(er.dbFile)
	if err != nil {
		ERROR.Println(err)
	}

	err = json.Unmarshal(content, &er.ctxRegistrationList)
	if err != nil {
		ERROR.Println(err)
	}
}
