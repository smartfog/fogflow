package main

import (
	"errors"
	. "github.com/smartfog/fogflow/common/ngsi"
	"strings"
	"time"
)

type Serializer struct{}

func (sz Serializer) DeSerializeEntity(expanded []interface{}) (map[string]interface{}, error) {
	entity := make(map[string]interface{})
	for _, val := range expanded {
		stringsMap := val.(map[string]interface{})
		for k, v := range stringsMap {
			if strings.Contains(k, "id") {
				if v != nil {
					entity["id"] = sz.getId(v.(interface{}))
				}
			} else if strings.Contains(k, "type") {
				if v != nil {
					entity["type"] = sz.getType(v.([]interface{}))
				}
			} else if strings.Contains(k, "location") {
				if v != nil {
					entity["location"] = sz.getLocation(v.([]interface{}))
				}
			} else if strings.Contains(k, "createdAt") {
				continue
			} else if strings.Contains(k, "modifiedAt") {
				continue
			} else {
				interfaceArray := v.([]interface{})
				if len(interfaceArray) > 0 {
					mp := interfaceArray[0].(map[string]interface{})
					typ := mp["@type"].([]interface{})
					if len(typ) > 0 {
						if strings.Contains(typ[0].(string), "Property") {
							property, err := sz.getProperty(mp)
							if err != nil {
								return entity, err
							} else {
								entity[k] = property
							}
						} else if strings.Contains(typ[0].(string), "Relationship") {

							relationship, err := sz.getRelationship(mp)
							if err != nil {
								return entity, err
							} else {
								entity[k] = relationship
							}
						}
					}
				}
			}
		}

	}
	entity["modifiedAt"] = time.Now().String()
	return entity, nil
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
			} else {
				// other subscription fields
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

func (sz Serializer) getId(id interface{}) string {
	Id := id.(string)
	return Id
}

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

func (sz Serializer) getProperty(propertyMap map[string]interface{}) (map[string]interface{}, error) {

	Property := make(map[string]interface{})
	for propertyField, fieldValue := range propertyMap {
		if strings.Contains(propertyField, "@type") {
			if fieldValue != nil {
				Property["type"] = sz.getType(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "hasValue") {
			if fieldValue != nil {
				Property["value"] = sz.getValueFromArray(fieldValue.([]interface{}))
				if Property["value"] == "nil" || Property["value"] == "" {
					err := errors.New("Property value can not be nil!")
					return Property, err
				}
			} else {
				err := errors.New("Property value can not be nil!")
				return Property, err
			}
		} else if strings.Contains(propertyField, "observedAt") {
			if fieldValue != nil {
				Property["observedAt"] = sz.getDateAndTimeValue(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "datasetId") {
			if fieldValue != nil {
				Property["datasetId"] = sz.getDatasetId(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "instanceId") {
			if fieldValue != nil {
				Property["instanceId"] = sz.getInstanceId(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "unitCode") {
			if fieldValue != nil {
				Property["unitCode"] = sz.getUnitCode(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "providedBy") {
			if fieldValue != nil {
				Property["providedBy"] = sz.getProvidedBy(fieldValue.([]interface{}))
			}
		} else if strings.Contains(propertyField, "createdAt") {
			continue
		} else if strings.Contains(propertyField, "modifiedA") {
			continue
		} else { // Nested property or relationship

			var typ string
			nested := fieldValue.([]interface{})
			for _, val := range nested {
				mp := val.(map[string]interface{})
				typInterface := mp["@type"].([]interface{})
				typ = typInterface[0].(string)
				if strings.Contains(typ, "Property") {
					property, err := sz.getProperty(mp)
					if err != nil {
						return Property, err
					} else {
						Property[propertyField] = property
					}
				} else if strings.Contains(typ, "Relationship") {
					relationship, err := sz.getRelationship(mp)
					if err != nil {
						return Property, err
					} else {
						Property[propertyField] = relationship
					}
				}
			}
		}
	}
	//Property["modifiedAt"] = time.Now().String()
	return Property, nil
}

func (sz Serializer) getRelationship(relationshipMap map[string]interface{}) (map[string]interface{}, error) {
	Relationship := make(map[string]interface{})
	for relationshipField, fieldValue := range relationshipMap {
		if strings.Contains(relationshipField, "@type") {
			if fieldValue != nil {
				Relationship["type"] = sz.getType(fieldValue.([]interface{}))
			}
		} else if strings.Contains(relationshipField, "hasObject") {
			if fieldValue != nil {
				Relationship["object"] = sz.getIdFromArray(fieldValue.([]interface{}))
				if Relationship["object"] == "nil" || Relationship["object"] == "" {
					err := errors.New("Relationship Object value can not be nil!")
					return Relationship, err
				}
			} else {
				err := errors.New("Relationship Object value can not be nil!")
				return Relationship, err
			}
		} else if strings.Contains(relationshipField, "Object") {
			if fieldValue != nil {
				Relationship["object"] = sz.getValueFromArray(fieldValue.([]interface{})).(string)
			} else {
				err := errors.New("Relationship Object value can not be nil!")
				return Relationship, err
			}
		} else if strings.Contains(relationshipField, "observedAt") {
			if fieldValue != nil {
				Relationship["observedAt"] = sz.getDateAndTimeValue(fieldValue.([]interface{}))
			}
		} else if strings.Contains(relationshipField, "providedBy") {
			if fieldValue != nil {
				Relationship["providedBy"] = sz.getProvidedBy(fieldValue.([]interface{}))
			}
		} else if strings.Contains(relationshipField, "datasetId") {
			if fieldValue != nil {
				Relationship["datasetId"] = sz.getDatasetId(fieldValue.([]interface{}))
			}
		} else if strings.Contains(relationshipField, "instanceId") {
			if fieldValue != nil {
				Relationship["instanceId"] = sz.getInstanceId(fieldValue.([]interface{}))
			}
		} else if strings.Contains(relationshipField, "createdAt") {
			continue
		} else if strings.Contains(relationshipField, "modifiedA") {
			continue
		} else { // Nested property or relationship
			var typ string
			nested := fieldValue.([]interface{})
			for _, val := range nested {
				mp := val.(map[string]interface{})
				typInterface := mp["@type"].([]interface{})
				typ = typInterface[0].(string)

				if strings.Contains(typ, "Property") {
					property, err := sz.getProperty(mp)
					if err != nil {
						return Relationship, err
					} else {
						Relationship[relationshipField] = property
					}
				} else if strings.Contains(typ, "Relationship") {
					relationship, err := sz.getRelationship(mp)
					if err != nil {
						return Relationship, err
					} else {
						Relationship[relationshipField] = relationship
					}
				}
			}
		}
	}
	//Relationship["modifiedAt"] = time.Now().String()
	return Relationship, nil
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
				} //Value is not  overwritten, in case of multiple values in payload, value array never return
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
