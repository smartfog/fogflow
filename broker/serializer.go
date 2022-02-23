package main

import (
	"errors"
	. "fogflow/common/constants"
	. "fogflow/common/ldContext"
	. "fogflow/common/ngsi"
	"strings"
	"fmt"
	"time"
)

type Serializer struct{}

type fName func(map[string]interface{}) (map[string]interface{}, error)

func (sz Serializer) geoHandler(geoMap map[string]interface{}) (map[string]interface{}, error) {
	geoResult := make(map[string]interface{})
	var err error
	fmt.Println("test username")
	geoValue := false
	for key, val := range geoMap {
		switch key {
		case NGSI_LD_TYPE:
			if val != nil {
				geoResult[key] = getType(val.([]interface{}))
			}
		case NGSI_LD_CREATEDAT:
			if val != nil {
				geoResult[key] = getCreatedTime(val.([]interface{}))
			}
		case NGSI_LD_OBSERVED_AT:
			if val != nil {
				geoResult[key] = getObservedTime(val.([]interface{}))
			}
		case NGSI_LD_MODIFIEDAT:
			if val != nil {
				geoResult[key] = getModifiedTime(val.([]interface{}))
			}
		case NGSI_LD_HAS_VALUE:
			if val != nil {
				geoValue = true
				geoResult[key] = getPropertyValue(val.([]interface{}))
			}
		case NGSI_LD_DATASETID:
			if val != nil {
				geoResult[key] = getDataSetID(val.([]interface{}))
			}
		case NGSI_LD_INSTANCEID:
			if val != nil {
				geoResult[key] = getInstanceID(val.([]interface{}))
			}
		case NGSILD_UniCode:
			if val != nil {
				geoResult[key] = val
			}
		default:
			var interfaceArray []interface{}
			switch val.(type) {
			case []interface{}:
				interfaceArray = val.([]interface{})
			default:
				interfaceArray = make([]interface{}, 0)
				geoResult[key] = val
			}
			if len(interfaceArray) > 0 {
				attrHandler, err := sz.getAttrType(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return geoResult, err
				}
				geoResult[key], err = attrHandler(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return geoResult, err
				}
			}
		}
	}
	if geoValue == false {
		err = errors.New("value field is mandatory!")
	}
	return geoResult, err
}

func (sz Serializer) proprtyHandler(propertyMap map[string]interface{}) (map[string]interface{}, error) {
	propertyResult := make(map[string]interface{})
	var err error
	propertyValue := false
	for key, val := range propertyMap {
		switch key {
		case NGSI_LD_TYPE:
			if val != nil {
				//propertyResult[key] = getType(val.([]interface{}))
				propertyResult[key] = getType(val)
			}
		case NGSI_LD_OBSERVED_AT:
			if val != nil {
				propertyResult[key] = getObservedTime(val.([]interface{}))
			}
		case NGSI_LD_CREATEDAT:
			if val != nil {
				propertyResult[key] = getCreatedTime(val.([]interface{}))
			}
		case NGSI_LD_MODIFIEDAT:
			if val != nil {
				propertyResult[key] = getModifiedTime(val.([]interface{}))
			}
		case NGSI_LD_HAS_VALUE:
			if val != nil {
				propertyValue = true
				propertyResult[key] = getPropertyValue(val.([]interface{}))
			}
		case NGSI_LD_DATASETID:
			if val != nil {
				propertyResult[key] = getDataSetID(val.([]interface{}))
			}
		case NGSI_LD_INSTANCEID:
			if val != nil {
				propertyResult[key] = getInstanceID(val.([]interface{}))
			}
		case NGSILD_UniCode:
			if val != nil {
				propertyResult[key] = val
			}
		default:
			interfaceArray := make([]interface{}, 0)
			switch val.(type) {
			case []interface{}:
				interfaceArray = val.([]interface{})
			default:
				propertyValue = true
				propertyResult[key] = val
			}
			if len(interfaceArray) > 0 {
				attrHandler, err := sz.getAttrType(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return propertyResult, err
				}
				propertyResult[key], err = attrHandler(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return propertyResult, err
				}
			}
		}
	}
	if propertyValue == false {
		err = errors.New("value field is mandatory!")
	}
	return propertyResult, err
}

func (sz Serializer) relHandler(relmap map[string]interface{}) (map[string]interface{}, error) {
	relResult := make(map[string]interface{})
	var err error
	ralationobject := false
	for key, val := range relmap {
		switch key {
		case NGSI_LD_TYPE:
			if val != nil {
				relResult[key] = getType(val.([]interface{}))
			}
		case NGSI_LD_CREATEDAT:
			if val != nil {
				relResult[key] = getCreatedTime(val.([]interface{}))
			}
		case NGSI_LD_OBSERVED_AT:
			if val != nil {
				relResult[key] = getObservedTime(val.([]interface{}))
			}
		case NGSI_LD_MODIFIEDAT:
			if val != nil {
				relResult[key] = getModifiedTime(val.([]interface{}))
			}
		case NGSI_LD_HAS_OBJECT:
			if val != nil {
				ralationobject = true
				relResult[key] = getPropertyValue(val.([]interface{}))
			}
		case NGSI_LD_DATASETID:
			if val != nil {
				relResult[key] = getDataSetID(val.([]interface{}))
			}
		case NGSI_LD_INSTANCEID:
			if val != nil {
				relResult[key] = getInstanceID(val.([]interface{}))
			}
		case NGSILD_UniCode:
			if val != nil {
				relResult[key] = val
			}
		default:
			interfaceArray := make([]interface{}, 0)
			switch val.(type) {
			case []interface{}:
				interfaceArray = val.([]interface{})
			default:
				ralationobject = true
				relResult[key] = val
			}
			if len(interfaceArray) > 0 {
				attrHandler, err := sz.getAttrType(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return relResult, err
				}
				relResult[key], err = attrHandler(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return relResult, err
				}
			}
		}
	}
	if ralationobject == false {
		err = errors.New("object field is mandatory!")
	}
	return relResult, err
}

func (sz Serializer) handler(ExpEntity interface{}) (map[string]interface{}, error) {
	ExpEntityMap := ExpEntity.(map[string]interface{})
	resultEntity := make(map[string]interface{})
	for key, val := range ExpEntityMap {
		switch key {
		case NGSI_LD_ID:
			if val != nil {
				resultEntity[key] = getEntityId(val.(interface{}))
			}
		case NGSI_LD_TYPE:
			if val != nil {
				resultEntity[key] = getType(val.([]interface{}))
			}
		case NGSI_LD_CREATEDAT:
			if val != nil {
				resultEntity[key] = getCreatedTime(val.([]interface{}))
			}
		case NGSI_LD_MODIFIEDAT:
			if val != nil {
				resultEntity[key] = getModifiedTime(val.([]interface{}))
			}
		default:
			interfaceArray := val.([]interface{})
			if len(interfaceArray) > 0 {
				attrHandler, err := sz.getAttrType(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return resultEntity, err
				}
				resultEntity[key], err = attrHandler(interfaceArray[0].(map[string]interface{}))
				if err != nil {
					return resultEntity, err
				}
			}

		}
	}
	return resultEntity, nil
}

func (sz Serializer) getTypeValue(typ string) (fName, error) {
	var funcName fName
	var err error
	switch typ {
	case LD_RELATIONSHIP:
		funcName = sz.relHandler
	case LD_PRPERTY:
		funcName = sz.proprtyHandler
	case LD_GEOPROPERTY:
		funcName = sz.geoHandler
	default:
		err = errors.New("Unknown Type !")
	}
	return funcName, err
}

func (sz Serializer) getAttrType(attr map[string]interface{}) (fName, error) {
	var funcName fName
	var err error
	if _, okey := attr["@type"]; okey == false {
		err := errors.New("attribute type can not be nil!")
		return funcName, err
	}
	var resType interface{}
	var tValue int
	switch attr["@type"].(type) {
	case []interface{}:
		resType = attr["@type"]
		tValue = 1
	case string:
		resType = attr["@type"]
		tValue = 2
	default:
		err := errors.New("Unknown Type!")
		return funcName, err
	}
	if tValue == 1 {
		resType1 := resType.([]interface{})
		if len(resType1) > 0 {
			Type1 := resType1[0].(string)
			funcName, err = sz.getTypeValue(Type1)
		}
	} else if tValue == 2 {
		resType2 := resType.(string)
		funcName, err = sz.getTypeValue(resType2)
	} else {
		err := errors.New("Unknown Type!")
		return funcName, err
	}
	return funcName, err
}

func (sz Serializer) getId(id interface{}) string {
	Id := id.(string)
	return Id
}

func (sz Serializer) DeSerializeEntity(expEntities []interface{}) (map[string]interface{}, error) {
	expEntity := expEntities[0]
	result, err := sz.handler(expEntity.(map[string]interface{}))
	return result, err
}

func (sz Serializer) DeSerializeSubscription(expanded []interface{}) (LDSubscriptionRequest, error) {
	subscription := LDSubscriptionRequest{}
	for _, val := range expanded {
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			if strings.Contains(k, "@id") {
				if v != nil {
					subscription.Id = sz.getId(v.(interface{}))
				}
			} else if strings.Contains(k, "@type") {
				if v != nil {
					subscription.Type = sz.getType(v.([]interface{}))
				}
			} else if strings.Contains(k, "description") {
				if v != nil {
					subscription.Description = sz.getValue(v.([]interface{})).(string)
				}
			} else if strings.Contains(k, "notification") {
				if v != nil {
					notification, err := sz.getNotification(v.([]interface{}))
					if err != nil {
						return subscription, err
					} else {
						subscription.Notification = notification
					}
				}
			} else if strings.Contains(k, "entities") {
				if v != nil {
					subscription.Entities = sz.getEntities(v.([]interface{}))
				}
			} else if strings.Contains(k, "name") {
				if v != nil {
					subscription.Name = sz.getValue(v.([]interface{})).(string)
				}
			} else if strings.Contains(k, "watchedAttributes") {
				if v != nil {
					subscription.WatchedAttributes = sz.getArrayOfIds(v.([]interface{}))
				}
			} else if strings.Contains(k, "georel") {
				if v != nil {
					switch v.(type) {
					case []interface{}:
						data := v.([]interface{})
						dataMap := data[0].(map[string]interface{})
						dataMap["@context"] = DEFAULT_CONTEXT
						resolved, _ := compactData(dataMap, DEFAULT_CONTEXT)
						subscription.Restriction, _ = sz.assignRestriction(resolved.(map[string]interface{}))
					default:
						err := errors.New("Unknown Type!")
						return subscription, err
					}
				}
			} else {
			}
		}
	}
	subscription.ModifiedAt = time.Now().String()
	subscription.IsActive = true
	return subscription, nil
}

func (sz Serializer) DeSerializeType(attrPayload []interface{}) string {
	var attr string
	if len(attrPayload) > 0 {
		attrMap := attrPayload[0].(map[string]interface{})
		attrs := attrMap["@type"].([]interface{})
		attr = attrs[0].(string)
	}
	return attr
}

/*func (sz Serializer) getId(id interface{}) string {
	Id := id.(string)
	return Id
}*/

func (sz Serializer) getType(typ []interface{}) string {
	var Type, Type1 string
	if len(typ) > 0 {
		Type1 = typ[0].(string)
		if strings.Contains(Type1, "GeoProperty") || strings.Contains(Type1, "geoproperty") {
			Type = "GeoProperty"
		} else if strings.Contains(Type1, "Point") || strings.Contains(Type1, "point") {
			Type = "Point"
		} else if strings.Contains(Type1, "Relationship") || strings.Contains(Type1, "relationship") {
			Type = "Relationship"
		} else if strings.Contains(Type1, "Property") || strings.Contains(Type1, "property") {
			Type = "Property"
		} else if strings.Contains(Type1, "person") || strings.Contains(Type1, "Person") {
			Type = "Person"
		} else {
			Type = typ[0].(string)
		}
	}
	return Type
}

func (sz Serializer) getValue(hasValue []interface{}) interface{} {

	//Value := make(map[string]interface{})
	var Value interface{}
	if len(hasValue) > 0 {
		val := hasValue[0].(map[string]interface{})
		Value = val["@value"]
	}
	return Value
}

func (sz Serializer) getValueFromArray(hasValue []interface{}) interface{} {
	Value := make(map[string]interface{})
	value := make([]interface{}, 0)
	if len(hasValue) == 0 {
		return value
	} else if len(hasValue) > 0 {
		for _, oneValue := range hasValue {
			if val := oneValue.(map[string]interface{}); val != nil {
				if val["@type"] != nil {
					Value["Type"] = val["@type"].(string)
					Value["Value"] = val["@value"].(interface{})
					return Value
				}
				if val["@value"] != nil {
					value = append(value, val["@value"].(interface{}))
				}
			}
		}
	}
	return value
}

func (sz Serializer) getIdFromArray(object []interface{}) string {
	var Id string
	if len(object) > 0 {
		hasObject := object[0].(map[string]interface{})
		Id = hasObject["@id"].(string)
	}
	return Id
}

func (sz Serializer) getDateAndTimeValue(dateTimeValue []interface{}) string {
	var DateTimeValue string
	if len(dateTimeValue) > 0 {
		observedAtMap := dateTimeValue[0].(map[string]interface{})
		if strings.Contains(observedAtMap["@type"].(string), "DateTime") {
			DateTimeValue = observedAtMap["@value"].(string)
		}
	}
	return DateTimeValue
}

func (sz Serializer) getProvidedBy(providedBy []interface{}) ProvidedBy {
	ProvidedBy := ProvidedBy{}
	if len(providedBy) > 0 {
		providedByMap := providedBy[0].(map[string]interface{})
		for k, v := range providedByMap {
			if strings.Contains(k, "@type") {
				ProvidedBy.Type = sz.getType(v.([]interface{}))
			} else if strings.Contains(k, "hasObject") {
				ProvidedBy.Object = sz.getIdFromArray(v.([]interface{}))
			}
		}
	}
	return ProvidedBy
}

//DATASET_ID
func (sz Serializer) getDatasetId(datasetId []interface{}) string {
	var DatasetId string
	if len(datasetId) > 0 {
		datasetIdMap := datasetId[0].(map[string]interface{})
		DatasetId = datasetIdMap["@id"].(string)
	}
	return DatasetId
}

//INSTANCE_ID
func (sz Serializer) getInstanceId(instanceId []interface{}) string {
	var InstanceId string
	if len(instanceId) > 0 {
		instanceIdMap := instanceId[0].(map[string]interface{})
		InstanceId = instanceIdMap["@id"].(string)
	}
	return InstanceId
}

//UNIT_CODE
func (sz Serializer) getUnitCode(unitCode []interface{}) string {
	var UnitCode string
	if len(unitCode) > 0 {
		unitCodeMap := unitCode[0].(map[string]interface{})
		UnitCode = unitCodeMap["@value"].(string)
	}
	return UnitCode
}

//LOCATION
func (sz Serializer) getLocation(location []interface{}) LDLocation {
	Location := LDLocation{}
	if len(location) > 0 {
		locationMap := location[0].(map[string]interface{})
		for k, v := range locationMap {
			if strings.Contains(k, "@type") {
				Location.Type = sz.getType(v.([]interface{}))
			} else if strings.Contains(k, "hasValue") {
				Location.Value = sz.getLocationValue(v.([]interface{}))
			}
		}
	}
	return Location
}

func (sz Serializer) getLocationValue(locationValue []interface{}) interface{} {
	if len(locationValue) > 0 {
		locationValueMap := locationValue[0].(map[string]interface{})
		if locationValueMap["@value"] != nil {
			valueScalar := locationValueMap["@value"].(interface{})
			stringValue := valueScalar.(string)
			return stringValue
		} else if locationValueMap["@type"] != nil {
			LocationValue := LDLocationValue{}
			for k, v := range locationValueMap {
				if strings.Contains(k, "@type") {
					LocationValue.Type = sz.getType(v.([]interface{}))
				}
			}
			for k, v := range locationValueMap {
				if strings.Contains(k, "coordinates") {
					if v != nil {
						if strings.Contains(LocationValue.Type, "Point") {
							LocationValue.Coordinates = sz.getPointLocation(v.([]interface{}))
						} else if strings.Contains(LocationValue.Type, "GeometryCollection") {
							LocationValue.Geometries = sz.getGeometryCollectionLocation(v.([]interface{}))
						} else if strings.Contains(LocationValue.Type, "LineString") || strings.Contains(LocationValue.Type, "Polygon") || strings.Contains(LocationValue.Type, "MultiPoint") || strings.Contains(LocationValue.Type, "MultiLineString") || strings.Contains(LocationValue.Type, "MultiPolygon") {
							LocationValue.Coordinates = sz.getArrayofCoordinates(v.([]interface{}))
						}
					}
				}
			}
			return LocationValue
		}
	}
	return nil
}

func (sz Serializer) getPointLocation(coordinates []interface{}) []float64 {
	var Coordinates []float64 //contains longitude & latitude values in order.

	for _, v := range coordinates {
		coord := v.(map[string]interface{})
		Coordinates = append(Coordinates, coord["@value"].(float64))
	}
	return Coordinates
}

func (sz Serializer) getArrayofCoordinates(coordinates []interface{}) [][]float64 {
	var Coordinates [][]float64 //Array contains point coordinates with longitude & latitude values in order
	for i := 0; i < len(coordinates); i = i + 2 {
		var coord []float64
		fCor := coordinates[i].(map[string]interface{})
		sCor := coordinates[i+1].(map[string]interface{})
		coord = append(coord, fCor["@value"].(float64))
		coord = append(coord, sCor["@value"].(float64))
		Coordinates = append(Coordinates, coord)
	}
	return Coordinates
}

func (sz Serializer) getGeometryCollectionLocation(geometries []interface{}) []Geometry {
	Geometries := []Geometry{}
	for _, val := range geometries {
		geometry := Geometry{}
		geometryValueMap := val.(map[string]interface{})
		for k, v := range geometryValueMap {
			if strings.Contains(k, "@Type") {
				geometry.Type = sz.getType(v.([]interface{}))
			} else if strings.Contains(k, "coordinates") {
				if strings.Contains(geometry.Type, "Point") {
					geometry.Coordinates = sz.getPointLocation(v.([]interface{}))
				} else {
					geometry.Coordinates = sz.getArrayofCoordinates(v.([]interface{}))
				}
			}
		}
		Geometries = append(Geometries, geometry)
	}
	return Geometries
}

func (sz Serializer) getArrayOfIds(arrayOfIds []interface{}) []string {
	var ArrayOfIds []string
	for _, v := range arrayOfIds {
		idValue := v.(map[string]interface{})
		id := idValue["@id"].(string)
		ArrayOfIds = append(ArrayOfIds, id)
	}
	return ArrayOfIds
}

func (sz Serializer) getEntities(entitiesArray []interface{}) []EntityId {
	entities := []EntityId{}
	for _, val := range entitiesArray {
		entityId := EntityId{}
		entityFields := val.(map[string]interface{})
		for k, v := range entityFields {
			if strings.Contains(k, "@id") {
				entityId.ID = sz.getId(v.(string))
			} else if strings.Contains(k, "@type") {
				entityId.Type = sz.getType(v.([]interface{}))
			} else if strings.Contains(k, "idPattern") {
				entityId.IdPattern = sz.getStringValue(v.([]interface{}))
			}
		}
		entities = append(entities, entityId)
	}
	return entities
}

func (sz Serializer) getStringValue(value []interface{}) string {
	var Value string
	if len(value) > 0 {
		val := value[0].(map[string]interface{})
		Value = val["@value"].(string)
	}
	return Value
}

func (sz Serializer) getNotification(notificationArray []interface{}) (NotificationParams, error) {
	notification := NotificationParams{}
	for _, val := range notificationArray {
		notificationFields := val.(map[string]interface{})
		for k, v := range notificationFields {
			if strings.Contains(k, "attributes") {
				notification.Attributes = sz.getArrayOfIds(v.([]interface{}))
			} else if strings.Contains(k, "endpoint") {
				endpoint, err := sz.getEndpoint(v.([]interface{}))
				if err != nil {
					return notification, err
				} else {
					notification.Endpoint = endpoint
				}
			} else if strings.Contains(k, "format") {
				notification.Format = sz.getStringValue(v.([]interface{}))
			}
		}
	}
	return notification, nil
}

func (sz Serializer) getEndpoint(endpointArray []interface{}) (Endpoint, error) {
	endpoint := Endpoint{}
	for _, val := range endpointArray {
		endpointFields := val.(map[string]interface{})
		for k, v := range endpointFields {
			if strings.Contains(k, "accept") {
				if v != nil {
					endpoint.Accept = sz.getStringValue(v.([]interface{}))
				}
			} else if strings.Contains(k, "uri") {
				if v != nil {
					endpoint.URI = sz.getStringValue(v.([]interface{}))
				} else {
					err := errors.New("URI can not be nil!")
					return endpoint, err
				}
			}
		}
	}
	return endpoint, nil
}

func (sz Serializer) afterString(str string, markingStr string) string {
	// Get sub-string after markingStr string
	li := strings.LastIndex(str, markingStr)
	liAdjusted := li + len(markingStr)
	if liAdjusted >= len(str) {
		return ""
	}
	return str[liAdjusted:len(str)]
}

// get NGSILD type

func (sz Serializer) getQueryType(QueryData map[string]interface{}) (string, error) {
	var typ string
	var err error
	if val, ok := QueryData["@type"]; ok == true {
		valueResult := val.([]interface{})
		typ = valueResult[0].(string)
	} else if val, ok := QueryData["type"]; ok == true {
		valueResult := val.([]interface{})
		typ = valueResult[0].(string)
	} else {
		err = errors.New("type can not be Empty")
	}
	return typ, err
}

// get NGSILD attributes
func (sz Serializer) getQueryAttributes(attributes, context []interface{}) ([]string, error) {
	attributesList := make([]string, 0)
	if len(attributes) <= 0 {
		err := errors.New("Zero length attribute list is not allowed")
		return attributesList, err
	}
	for _, attr := range attributes {
		//fmt.Println("attr",attr)
		valueMap := attr.(map[string]interface{})
		ldobject := getLDobject(valueMap["@value"], context)
		ExpandedAttr, _ := ExpandEntity(ldobject)
		attrUri := getAttribute(ExpandedAttr)
		attributesList = append(attributesList, attrUri)
	}
	return attributesList, nil
}
func (sz Serializer) getEntityId(id interface{}, fs string) string {
	ID := id.(string) + "@" + fs
	return ID
}

func (sz Serializer) getEntityType(typ interface{}) string {
	Etype := typ.([]interface{})
	return Etype[0].(string)
}

func (sz Serializer) resolveEntity(entityobj interface{}, fs string) EntityId {
	entity := EntityId{}
	entitymap := entityobj.(map[string]interface{})
	if val, ok := entitymap["@id"]; ok == true {
		entity.ID = sz.getEntityId(val, fs)
	} else if val, ok := entitymap["id"]; ok == true {
		entity.ID = sz.getEntityId(val, fs)
	} else {
		entity.IsPattern = true
	}
	if val, ok := entitymap["@type"]; ok == true {
		entity.Type = sz.getEntityType(val)
	} else if val, ok := entitymap["type"]; ok == true {
		entity.Type = sz.getEntityType(val)
	}
	return entity
}

func (sz Serializer) assignRestriction(restriction map[string]interface{}) (Restriction, error) {
	restrictions := Restriction{}
	if _, ok := restriction["coordinates"]; ok == true {
		restrictions.Cordinates = restriction["coordinates"]
	}
	if _, ok := restriction["geometry"]; ok == true {
		restrictions.Geometry = restriction["geometry"].(string)
	}
	if _, ok := restriction["georel"]; ok == true {
		restrictions.Georel = restriction["georel"].(string)
	}
	restrictions.RestrictionType = "ld"

	return restrictions, nil

}

//get Entities
func (sz Serializer) getQueryEntities(entities []interface{}, fs string) ([]EntityId, error) {
	entitiesList := make([]EntityId, 0)
	if len(entities) <= 0 {
		err := errors.New("Zero length Entity List is not allowed")
		return entitiesList, err
	}
	for _, val := range entities {
		entity := sz.resolveEntity(val, fs)
		if entity.ID == "" && entity.Type == "" {
			continue
		} else {
			entitiesList = append(entitiesList, entity)
		}
	}
	return entitiesList, nil
}

// serialize NGSIld Query
func (sz Serializer) uploadQueryContext(expanded interface{}, fs string, context []interface{}) (LDQueryContextRequest, error) {
	ngsildQueryContext := LDQueryContextRequest{}
	expandedArray := expanded.([]interface{})
	QueryData := expandedArray[0].(map[string]interface{})
	typ, err := sz.getQueryType(QueryData)
	if err != nil {
		return ngsildQueryContext, err
	}
	ngsildQueryContext.Type = typ
	var newErr error
forloop:
	for key, value := range QueryData {
		switch key {
		case NGSI_LD_ATTRS:
			ngsildQueryContext.Attributes, newErr = sz.getQueryAttributes(value.([]interface{}), context)
			if newErr != nil {
				break forloop
			}
		case NGSI_LD_ENTITIES:
			ngsildQueryContext.Entities, newErr = sz.getQueryEntities(value.([]interface{}), fs)
			if newErr != nil {
				break forloop
			}
		case NGSI_LD_GEO_QUERY:
			data := QueryData[key].([]interface{})
			dataMap := data[0].(map[string]interface{})
			dataMap["@context"] = DEFAULT_CONTEXT
			resolved, newErr := compactData(dataMap, DEFAULT_CONTEXT)
			if newErr != nil {
				break forloop
			}
			ngsildQueryContext.Restriction, err = sz.assignRestriction(resolved.(map[string]interface{}))
		case NGSI_LD_QUERY:

		default:
			continue forloop
		}
	}
	return ngsildQueryContext, newErr
}
