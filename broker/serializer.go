package main

import (
	"encoding/json" //to be removed later, used for printing only
	"errors"
	"fmt"
	. "github.com/smartfog/fogflow/common/constants"
	. "github.com/smartfog/fogflow/common/ngsi"
	"time"
)

type Serializer struct{}

func (sz Serializer) DeSerializeEntity(expanded []interface{}) (LDContextElement, error) {
	entity := LDContextElement{}
	for _, val := range expanded {
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			switch k {
			case ID:
				if v != nil {
					entity.Id = sz.getId(v.(interface{}))
				} else {
					err := errors.New("Id can not be nil!")
					return entity, err
				}
				break
			case TYPE:
				if v != nil {
					entity.Type = sz.getType(v.([]interface{}))
				} else {
					err := errors.New("Type can not be nil!")
					return entity, err
				}
				break
			case LOCATION:
				if v != nil {
					entity.Location = sz.getLocation(v.([]interface{}))
				}
				break
			case CREATED_AT:
				break
			default: // default cases like property, relationship here.
				interfaceArray := v.([]interface{})
				if len(interfaceArray) > 0 {
					mp := interfaceArray[0].(map[string]interface{})
					typ := mp[TYPE].([]interface{})
					if len(typ) > 0 {
						if typ[0].(string) == PROPERTY {
							property, err := sz.getProperty(k, mp)
							if err != nil {
								fmt.Println("Errored Property!")
								return entity, err
							} else {
								entity.Properties = append(entity.Properties, property)
								fmt.Println("Created Property...")
							}
						} else if typ[0].(string) == RELATIONSHIP {
							relationship, err := sz.getRelationship(k, mp)
							if err != nil {
								fmt.Println("Errored Relationship!")
								return entity, err
							} else {
								entity.Relationships = append(entity.Relationships, relationship)
							}
						}
					}
				}
				break
			}
		}

	}
	entity.CreatedAt = time.Now().String()
	return entity, nil
}

func (sz Serializer) DeSerializeRegistration(expanded []interface{}) (CSourceRegistrationRequest, error) {
	registration := CSourceRegistrationRequest{}
	for _, val := range expanded {
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			switch k {
			case ID:
				if v != nil {
					registration.Id = sz.getId(v.(interface{}))
				}
				break
			case TYPE:
				if v != nil {
					registration.Type = sz.getType(v.([]interface{}))
				} else {
					err := errors.New("Type can not be nil!")
					return registration, err
				}
				break
			// timestamp from payload is taken as observationInterval in given datatypes in spec:
			case TIMESTAMP:
				if v != nil {
					registration.ObservationInterval = sz.getTimeStamp(v.([]interface{}))
				}
				break
			case DESCRIPTION:
				if v != nil {
					registration.Description = sz.getStringValue(v.([]interface{}))
				}
				break
			case ENDPOINT:
				if v != nil {
					registration.Endpoint = sz.getStringValue(v.([]interface{}))
				} else {
					err := errors.New("Endpoint value can not be nil!")
					return registration, err
				}
				break
			case EXPIRES:
				if v != nil {
					registration.Expires = sz.getDateAndTimeValue(v.([]interface{}))
				}
				break
			case INFORMATION:
				if v != nil {
					registration.Information = sz.getInformation(v.([]interface{}))
				} else {
					err := errors.New("Information value can not be nil!")
					return registration, err
				}
				break
			case LOCATION:
				if v != nil {
					registration.Location = sz.getStringValue(v.([]interface{}))
				}
				break
			case NAME:
				if v != nil {
					registration.Name = sz.getStringValue(v.([]interface{}))
				}
				break
			default:
				break
			}
		}
	}
	return registration, nil
}

func (sz Serializer) DeSerializeSubscription(expanded []interface{}) (LDSubscriptionRequest, error) {
	subscription := LDSubscriptionRequest{}
	for _, val := range expanded {
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			switch k {
			case ID:
				if v != nil {
					subscription.Id = sz.getId(v.(interface{}))
				}
				break
			case TYPE:
				if v != nil {
					subscription.Type = sz.getType(v.([]interface{}))
				} else {
					err := errors.New("Information value can not be nil!")
					return subscription, err
				}

				break
			case NAME:
				if v != nil {
					subscription.Name = sz.getValue(v.([]interface{})).(string)
				}
				break
			case ENTITIES:
				if v != nil {
					subscription.Entities = sz.getEntities(v.([]interface{}))
				}
				break
			case DESCRIPTION:
				if v != nil {
					subscription.Description = sz.getValue(v.([]interface{})).(string)
				}
				break
			case NOTIFICATION:
				if v != nil {
					notification, err := sz.getNotification(v.([]interface{}))
					if err != nil {
						return subscription, err
					} else {
						subscription.Notification = notification
					}
				} else {
					err := errors.New("Information value can not be nil!")
					return subscription, err
				}
				break
			case WATCHED_ATTRIBUTES:
				if v != nil {
					subscription.WatchedAttributes = sz.getArrayOfIds(v.([]interface{}))
				}
				break
			default:
				break
			}
		}
	}
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

func (sz Serializer) getProperty(propertyName string, propertyMap map[string]interface{}) (Property, error) {
	Property := Property{}
	if propertyName != "" {
		Property.Id = propertyName
		Property.Name = propertyName
	}
	Property.Type = PROPERTY
	if propertyMap[HAS_VALUE] != nil {
		Property.Value = sz.getValueFromArray(propertyMap[HAS_VALUE].([]interface{}))
	} /*else {
		err := errors.New("Property Value can not be nil!")
		return Property, err
	}*/
	if propertyMap[OBSERVED_AT] != nil {
		Property.ObservedAt = sz.getDateAndTimeValue(propertyMap[OBSERVED_AT].([]interface{}))
	}

	if propertyMap[DATASET_ID] != nil {
		Property.DatasetId = sz.getDatasetId(propertyMap[DATASET_ID].([]interface{}))
	}

	if propertyMap[INSTANCE_ID] != nil {
		Property.InstanceId = sz.getInstanceId(propertyMap[INSTANCE_ID].([]interface{}))
	}

	Property.ModifiedAt = time.Now().String()

	if propertyMap[UNIT_CODE] != nil {
		Property.UnitCode = sz.getUnitCode(propertyMap[UNIT_CODE].(interface{}))
	}

	return Property, nil
}

func (sz Serializer) getRelationship(relationshipName string, relationshipMap map[string]interface{}) (Relationship, error) {
	Relationship := Relationship{}

	if relationshipName != "" {
		Relationship.Id = relationshipName
		Relationship.Name = relationshipName
	}

	Relationship.Type = RELATIONSHIP

	if relationshipMap[HAS_OBJECT] != nil {
		Relationship.Object = sz.getIdFromArray(relationshipMap[HAS_OBJECT].([]interface{}))
	} else if relationshipMap[OBJECT] != nil {
		Relationship.Object = sz.getValueFromArray(relationshipMap[OBJECT].([]interface{})).(string)
	} else {
		err := errors.New("Relationship Object value can not be nil!")
		return Relationship, err
	}

	if relationshipMap[OBSERVED_AT] != nil {
		Relationship.ObservedAt = sz.getDateAndTimeValue(relationshipMap[OBSERVED_AT].([]interface{}))
	}

	if relationshipMap[PROVIDED_BY] != nil {
		Relationship.ProvidedBy = sz.getProvidedBy(relationshipMap[PROVIDED_BY].([]interface{}))
	}

	if relationshipMap[DATASET_ID] != nil {
		Relationship.DatasetId = sz.getDatasetId(relationshipMap[DATASET_ID].([]interface{}))
	}

	if relationshipMap[INSTANCE_ID] != nil {
		Relationship.InstanceId = sz.getInstanceId(relationshipMap[INSTANCE_ID].([]interface{}))
	}

	Relationship.ModifiedAt = time.Now().String()

	return Relationship, nil
}

func (sz Serializer) getValue(hasValue []interface{}) interface{} {

	Value := make(map[string]interface{})
	if len(hasValue) > 0 {
		val := hasValue[0].(map[string]interface{})

		if val[TYPE] != nil {
			Value["Type"] = val[TYPE].(string)
			Value["Value"] = val[VALUE].(interface{})
		} else {
			Value["Value"] = hasValue[0]
		}
	}
	return Value
}

func (sz Serializer) getValueFromArray(hasValue []interface{}) interface{} {
	Value := make(map[string]interface{})
	var value interface{}
	if len(hasValue) > 0 {
		val := hasValue[0].(map[string]interface{})

		if val[TYPE] != nil {
			Value["Type"] = val[TYPE].(string)
			Value["Value"] = val[VALUE].(interface{})
			return Value
		}
		value = val[VALUE].(interface{})
	}
	return value
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
		providedByMap := providedBy[0].(map[string]interface{})
		ProvidedBy.Type = sz.getType(providedByMap[TYPE].([]interface{}))
		ProvidedBy.Object = sz.getIdFromArray(providedByMap[HAS_OBJECT].([]interface{}))
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
			LocationValue.Coordinates = sz.getArrayofCoordinates(locationValueMap[COORDINATES].([]interface{}))
			break
		case POLYGON:
			LocationValue.Coordinates = sz.getArrayofCoordinates(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_POINT:
			LocationValue.Coordinates = sz.getArrayofCoordinates(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_LINE_STRING:
			LocationValue.Coordinates = sz.getArrayofCoordinates(locationValueMap[COORDINATES].([]interface{}))
			break
		case MULTI_POLYGON:
			LocationValue.Coordinates = sz.getArrayofCoordinates(locationValueMap[COORDINATES].([]interface{}))
			break
		case GEOMETRY_COLLECTION:
			LocationValue.Geometries = sz.getGeometryCollectionLocation(locationValueMap[GEOMETRIES].([]interface{}))
			break
		}
	}
	return LocationValue
}

func (sz Serializer) getPointLocation(coordinates []interface{}) []float64 {
	var Coordinates []float64 //contains longitude & latitude values in order.

	for _, v := range coordinates {
		coord := v.(map[string]interface{})
		Coordinates = append(Coordinates, coord[VALUE].(float64))
	}
	return Coordinates
}

func (sz Serializer) getArrayofCoordinates(coordinates []interface{}) [][]float64 {
	var Coordinates [][]float64 //Array contains point coordinates with longitude & latitude values in order
	for i := 0; i < len(coordinates); i = i + 2 {
		var coord []float64
		coord = append(coord, coordinates[i].(float64))
		coord = append(coord, coordinates[i+1].(float64))
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
			switch k {
			case TYPE:
				geometry.Type = sz.getType(v.([]interface{}))
				break
			case COORDINATES:
				if geometry.Type == POINT {
					geometry.Coordinates = sz.getPointLocation(v.([]interface{}))
				} else {
					geometry.Coordinates = sz.getArrayofCoordinates(v.([]interface{}))
				}
				break
			}
		}
		Geometries = append(Geometries, geometry)
	}

	// Pretty print
	geos, _ := json.MarshalIndent(Geometries, "", " ")
	DEBUG.Println("Geometries:", string(geos))
	return Geometries
}

func (sz Serializer) getInformation(information []interface{}) []RegistrationInfo {
	regInfoArray := []RegistrationInfo{}
	for _, val := range information {
		infoVal := val.(map[string]interface{})
		regInfo := RegistrationInfo{}
		for k, v := range infoVal {
			switch k {
			case PROPERTIES:
				if v != nil {
					regInfo.Properties = sz.getArrayOfIds(v.([]interface{}))
				}
				break
			case RELATIONSHIPS:
				if v != nil {
					regInfo.Relationships = sz.getArrayOfIds(v.([]interface{}))
				}
				break
			case ENTITIES:
				if v != nil {
					regInfo.Entities = sz.getEntities(v.([]interface{}))
				}
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

func (sz Serializer) getEntities(entitiesArray []interface{}) []EntityId {
	entities := []EntityId{}
	for _, val := range entitiesArray {
		entityId := EntityId{}
		entityFields := val.(map[string]interface{})

		for k, v := range entityFields {

			switch k {
			case ID:
				entityId.ID = sz.getId(v.(string))
				break
			case TYPE:
				entityId.Type = sz.getType(v.([]interface{}))
				break
			case ID_PATTERN:
				entityId.IdPattern = sz.getStringValue(v.([]interface{}))
				break
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
		Value = val[VALUE].(string)
	}
	return Value
}

func (sz Serializer) getNotification(notificationArray []interface{}) (NotificationParams, error) {
	notification := NotificationParams{}
	for _, val := range notificationArray {
		notificationFields := val.(map[string]interface{})
		for k, v := range notificationFields {
			switch k {
			case ATTRIBUTES:
				notification.Attributes = sz.getArrayOfIds(v.([]interface{}))
				break
			case ENDPOINT:
				endpoint, err := sz.getEndpoint(v.([]interface{}))
				if err != nil {
					return notification, err
				} else {
					notification.Endpoint = endpoint
				}
				break
			case FORMAT:
				notification.Format = sz.getStringValue(v.([]interface{}))
				break
			default:
				break
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
			switch k {
			case ACCEPT:
				if v != nil {
					endpoint.Accept = sz.getStringValue(v.([]interface{}))
				}
				break
			case URI:
				if v != nil {
					endpoint.URI = sz.getStringValue(v.([]interface{}))
				} else {
					err := errors.New("URI can not be nil!")
					return endpoint, err
				}
				break
			}
		}
	}
	return endpoint, nil
}

func (sz Serializer) getTimeStamp(timestampArray []interface{}) TimeInterval {
        timeInterval := TimeInterval{}
        for _, timestamp := range timestampArray {
                timestampMap := timestamp.(map[string]interface{})
                for k, v := range timestampMap {
                        switch k {
                        case START:
                                timeInterval.Start = sz.getDateAndTimeValue(v.([]interface{}))
                        case END:
                                timeInterval.End = sz.getDateAndTimeValue(v.([]interface{}))
                        }
                }
        }
        return timeInterval
}
