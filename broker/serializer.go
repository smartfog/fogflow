package main

import (
	"fmt"
	. "github.com/smartfog/fogflow/common/constants"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type Serializer struct{}

func (sz Serializer) SerializeEntity(expanded []interface{}) LDContextElement {
	fmt.Println("Inside SerializeEntity.......")
	entity := LDContextElement{}
	for _, val := range expanded {
		fmt.Println("Inside outer for loop.......")
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			fmt.Println("Inside inner for loop.......")
			switch k {
			case ID:
				fmt.Println("case @id.......")
				entity.Id = sz.getId(v.(interface{}))
				break
			case TYPE:
				fmt.Println("case @type.......")
				entity.Type = sz.getType(v.([]interface{}))
				break
			case CREATED_AT:
				fmt.Println("case createdAt.......")
				entity.CreatedAt = sz.getCreatedAt(v.([]interface{}))
				break
			case LOCATION:
				fmt.Println("case location.......")
				entity.Location = sz.getLocation(v.([]interface{}))
				break
			default: // default cases like property, relationship here.
				interfaceArray := v.([]interface{})
				if len(interfaceArray) > 0 {
					mp := interfaceArray[0].(map[string]interface{})
					typ := mp[TYPE].([]interface{})
					if len(typ) > 0 {
						if typ[0].(string) == PROPERTY {
							fmt.Println("It is a property....")
							entity.Properties = append(entity.Properties, sz.getProperty(k, mp))
						} else if typ[0].(string) == RELATIONSHIP {
							fmt.Println("It is a relationship....")
							entity.Relationships = append(entity.Relationships, sz.getRelationship(k, mp))
						}
					}
				}
				break
			}
		}

	}
	return entity
}

func (sz Serializer) SerializeRegistration(expanded []interface{}) CSourceRegistrationRequest {
	fmt.Println("Inside SerializeRegistration.......")
	registration := CSourceRegistrationRequest{}
	for _, val := range expanded {
		fmt.Println("Inside outer for loop.......")
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			fmt.Println("Inside inner for loop.......")
			switch k {
			case W3_ID:
				fmt.Println("case W3_ID.......")
				registration.Id = sz.getIdFromArray(v.([]interface{}))
				break
			case W3_TYPE:
				fmt.Println("case W3_TYPE.......")
				registration.Type = sz.getIdFromArray(v.([]interface{}))
				break
			case TIMESTAMP:
				fmt.Println("case TIMESTAMP.......")
				//---------------------registration.Expires = sz.getDateAndTimeValue(v.([]interface{}))
				break
			case DESCRIPTION:
				fmt.Println("case DESCRIPTION.......")
				registration.Description = sz.getValue(v.([]interface{})).(string)
				break
			case ENDPOINT:
				fmt.Println("case ENDPOINT.......")
				registration.Endpoint = sz.getValue(v.([]interface{})).(string)
				break
			case EXPIRES:
				fmt.Println("case EXPIRES.......")
				registration.Expires = sz.getDateAndTimeValue(v.([]interface{}))
				break
			case INFORMATION:
				fmt.Println("case INFORMATION.......")
				registration.Information = sz.getInformation(v.([]interface{}))
				break
			case LOCATION:
				fmt.Println("case LOCATION.......")
				registration.Location = sz.getLocation(v.([]interface{}))
				break
			case NAME:
				fmt.Println("case NAME.......")
				registration.Name = sz.getValue(v.([]interface{})).(string)
				break
			default:
				fmt.Println("case default.......")
				break
			}
		}
	}
	return registration
}

func (sz Serializer) SerializeSubscription(expanded []interface{}) LDSubscriptionRequest {
	fmt.Println("Inside SerializeSubscription.......")
	subscription := LDSubscriptionRequest{}
	for _, val := range expanded {
		fmt.Println("Inside outer for loop.......")
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			fmt.Println("Inside inner for loop.......")
			switch k {
			case ID:
				fmt.Println("case ID.......")
				subscription.Id = sz.getIdFromArray(v.([]interface{}))
				break
			case TYPE:
				fmt.Println("case TYPE.......")
				subscription.Type = sz.getType(v.([]interface{}))
				break
			case NAME:
				fmt.Println("case NAME.......")
				subscription.Name = sz.getValue(v.([]interface{})).(string)
				break
			case ENTITIES:
				fmt.Println("case ENTITIES.......")
				subscription.Entities = sz.getEntities(v.([]interface{}))
				break
			case DESCRIPTION:
				fmt.Println("case DESCRIPTION.......")
				subscription.Description = sz.getValue(v.([]interface{})).(string)
				break
			case NOTIFICATION:
				fmt.Println("case NOTIFICATION.......")
				subscription.Notification = sz.getNotification(v.([]interface{}))
				break
			case WATCHED_ATTRIBUTES:
				fmt.Println("case WATCHED_ATTRIBUTES.......")
				subscription.WatchedAttributes = sz.getArrayOfIds(v.([]interface{}))
				break
			default:
				fmt.Println("case default.......")
				break
			}
		}
	}
	return subscription
}

func (sz Serializer) getId(id interface{}) string {
	Id := id.(string)
	return Id
}

func (sz Serializer) getType(typ []interface{}) string {
	var Type string
	if len(typ) > 0 {
		Type = typ[0].(string)
	}
	return Type
}

func (sz Serializer) getCreatedAt(createdAt []interface{}) string {
	var CreatedAt string
	if len(createdAt) > 0 {
		mp := createdAt[0].(map[string]interface{})
		if mp[TYPE].(string) == DATE_TIME {
			CreatedAt = mp[VALUE].(string)
		}
	}
	return CreatedAt
}

func (sz Serializer) getProperty(propertyName string, propertyMap map[string]interface{}) Property {
	Property := Property{}
	Property.Id = propertyName
	Property.Name = propertyName
	Property.Type = PROPERTY
	Property.Value = sz.getValue(propertyMap[HAS_VALUE].([]interface{}))
	Property.ObservedAt = sz.getDateAndTimeValue(propertyMap[OBSERVED_AT].([]interface{}))
	Property.DatasetId = sz.getDatasetId(propertyMap[DATASET_ID].([]interface{}))
	Property.InstanceId = sz.getInstanceId(propertyMap[INSTANCE_ID].([]interface{}))
	Property.CreatedAt = sz.getCreatedAt(propertyMap[CREATED_AT].([]interface{}))
	Property.ModifiedAt = sz.getModifiedAt(propertyMap[MODIFIED_AT].([]interface{}))
	Property.UnitCode = sz.getUnitCode(propertyMap[UNIT_CODE].(interface{}))
	return Property
}

func (sz Serializer) getRelationship(relationshipName string, relationshipMap map[string]interface{}) Relationship {
	Relationship := Relationship{}
	Relationship.Id = relationshipName
	Relationship.Name = relationshipName
	Relationship.Type = RELATIONSHIP
	Relationship.Object = sz.getIdFromArray(relationshipMap[HAS_OBJECT].([]interface{}))
	Relationship.ObservedAt = sz.getDateAndTimeValue(relationshipMap[OBSERVED_AT].([]interface{}))
	Relationship.ProvidedBy = sz.getProvidedBy(relationshipMap[PROVIDED_BY].([]interface{}))
	Relationship.DatasetId = sz.getDatasetId(relationshipMap[DATASET_ID].([]interface{}))
	Relationship.InstanceId = sz.getInstanceId(relationshipMap[INSTANCE_ID].([]interface{}))
	Relationship.CreatedAt = sz.getCreatedAt(relationshipMap[CREATED_AT].([]interface{}))
	Relationship.ModifiedAt = sz.getModifiedAt(relationshipMap[MODIFIED_AT].([]interface{}))
	return Relationship
}

func (sz Serializer) getValue(hasValue []interface{}) interface{} {
	Value := make(map[string]interface{})
	if len(hasValue) > 0 {
		val := hasValue[0].(map[string]interface{})

		if Value["Type"] = val[TYPE].(string); Value["Type"] != "" {
			Value["Value"] = val[VALUE].(interface{})
		} else {
			Value["Value"] = hasValue[0]
		}
	}
	return Value
}

func (sz Serializer) getIdFromArray(object []interface{}) string {
	var Id string
	if len(object) > 0 {
		hasObject := object[0].(map[string]interface{})
		Id = hasObject[ID].(string)
	}
	return Id
}

func (sz Serializer) getDateAndTimeValue(dateTimeValue []interface{}) string {
	var DateTimeValue string
	if len(dateTimeValue) > 0 {
		observedAtMap := dateTimeValue[0].(map[string]interface{})
		if observedAtMap[TYPE] == DATE_TIME {
			DateTimeValue = observedAtMap[VALUE].(string)
		}
	}
	return DateTimeValue
}

func (sz Serializer) getProvidedBy(providedBy []interface{}) ProvidedBy {
	ProvidedBy := ProvidedBy{}
	if len(providedBy) > 0 {
		providedByMap := providedBy[0].(map[string][]interface{})
		ProvidedBy.Type = sz.getType(providedByMap[TYPE])
		ProvidedBy.Object = sz.getIdFromArray(providedByMap[HAS_OBJECT])
	}
	return ProvidedBy
}

//DATASET_ID
func (sz Serializer) getDatasetId(datasetId []interface{}) string {
	return ""
}

//INSTANCE_ID
func (sz Serializer) getInstanceId(instanceId []interface{}) string {

	return ""
}

//MODIFIED_AT
func (sz Serializer) getModifiedAt(modifiedAt []interface{}) string {

	return ""
}

//UNIT_CODE
func (sz Serializer) getUnitCode(unitCode interface{}) string {

	return ""
}

//LOCATION
func (sz Serializer) getLocation(location []interface{}) LDLocation {
	Location := LDLocation{}
	if len(location) > 0 {
		locationMap := location[0].(map[string]interface{})
		Location.Type = sz.getType(locationMap[TYPE].([]interface{}))
		Location.Value = sz.getLocationValue(locationMap[HAS_VALUE].([]interface{}))
	}
	return Location
}

func (sz Serializer) getLocationValue(locationValue []interface{}) LDLocationValue {
	LocationValue := LDLocationValue{}
	if len(locationValue) > 0 {
		locationValueMap := locationValue[0].(map[string]interface{})
		LocationValue.Type = sz.getType(locationValueMap[TYPE].([]interface{}))
		switch LocationValue.Type {
		case POINT:
			LocationValue.Coordinates = sz.getPointLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case LINE_STRING:
			LocationValue.Coordinates = sz.getLineStringLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case POLYGON:
			LocationValue.Coordinates = sz.getPolygonLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_POINT:
			LocationValue.Coordinates = sz.getMultiPointLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_LINE_STRING:
			LocationValue.Coordinates = sz.getMultiLineStringLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_POLYGON:
			LocationValue.Coordinates = sz.getMultiPolygonLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		case GEOMETRY_COLLECTION:
			LocationValue.Coordinates = sz.getGeometryCollectionLocation(locationValueMap[COORDINATES].([]interface{}))
			break
		}
	}
	return LocationValue
}

//need to define all thesedata types in ngsi.go
func (sz Serializer) getPointLocation(coordinates []interface{}) []float64 {
	var Coordinates []float64 //contains longitude & latitude values in order.

	for _, v := range coordinates {
		coord := v.(map[string]interface{})
		Coordinates = append(Coordinates, coord[VALUE].(float64))
	}
	return Coordinates
}

func (sz Serializer) getLineStringLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getMultiPointLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getPolygonLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getMultiLineStringLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getMultiPolygonLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getGeometryCollectionLocation(coordinates []interface{}) interface{} {

	var v interface{}
	return v
}

func (sz Serializer) getTimestamp(timestamp []interface{}) {
	fmt.Println("Inside get timestamp...")

}

func (sz Serializer) getInformation(information []interface{}) []RegistrationInfo {
	regInfoArray := []RegistrationInfo{}
	for _, val := range information {
		infoVal := val.(map[string][]interface{})
		regInfo := RegistrationInfo{}
		for k, v := range infoVal {
			switch k {
			case PROPERTIES:
				regInfo.Properties = sz.getArrayOfIds(v)
				break
			case RELATIONSHIPS:
				regInfo.Relationships = sz.getArrayOfIds(v)
				break
			case ENTITIES:
				regInfo.Entities = sz.getEntities(v)
				break
			}
		}
		regInfoArray = append(regInfoArray, regInfo)
	}
	return regInfoArray
}

func (sz Serializer) getArrayOfIds(arrayOfIds []interface{}) []string {
	var ArrayOfIds []string
	for _, v := range arrayOfIds {
		idValue := v.(map[string]interface{})
		id := idValue[ID].(string)
		ArrayOfIds = append(ArrayOfIds, id)
	}
	return ArrayOfIds
}

func (sz Serializer) getEntities(entitiesArray []interface{}) []EntityInfo {
	entities := []EntityInfo{}
	for _, val := range entitiesArray {
		entityInfo := EntityInfo{}
		entityFields := val.(map[string]interface{})
		for k, v := range entityFields {
			switch k {
			case W3_ID:
				entityInfo.Id = sz.getIdFromArray(v.([]interface{}))
				break
			case W3_TYPE:
				entityInfo.Type = sz.getIdFromArray(v.([]interface{}))
				break
			case ID_PATTERN:
				entityInfo.IdPattern = sz.getStringValue(v.([]interface{}))
				break
			}
		}
		entities = append(entities, entityInfo)
	}
	return entities
}

func (sz Serializer) getStringValue(value []interface{}) string {
	var Value string
	if len(value) > 0 {
		val := value[0].(map[string]interface{})
		Value = val[VALUE].(string)
	}
	return Value
}

func (sz Serializer) getNotification(notificationArray []interface{}) NotificationParams {
	notification := NotificationParams{}
	for _, val := range notificationArray {
		notificationFields := val.(map[string]interface{})
		for k, v := range notificationFields {
			switch k {
			case ATTRIBUTES:
				notification.Attributes = sz.getArrayOfIds(v.([]interface{}))
				break
			case ENDPOINT:
				notification.Endpoint = sz.getEndpoint(v.([]interface{}))
				break
			case FORMAT:
				notification.Format = sz.getStringValue(v.([]interface{}))
				break
			default:
				fmt.Println("case default.......")
				break
			}
		}
	}
	return notification
}

func (sz Serializer) getEndpoint(endpointArray []interface{}) Endpoint {
	endpoint := Endpoint{}
	for _, val := range endpointArray {
		endpointFields := val.(map[string]interface{})
		for k, v := range endpointFields {
			switch k {
			case ACCEPT:
				endpoint.Accept = sz.getStringValue(v.([]interface{}))
				break
			case URI:
				endpoint.URI = sz.getStringValue(v.([]interface{}))
				break
			}
		}
	}
	return endpoint
}

// Check the type of values: can be json object or string
func (sz Serializer) checkType(value interface{}) string {
	switch value.(type) {
	case bool:
		return "bool"
	case int:
		return "int"
	case uint:
		return "uint"
	case string:
		return "string"
	case float32:
		return "float32"
	case float64:
		return "float64"
	default:
		return ""
	}
}
