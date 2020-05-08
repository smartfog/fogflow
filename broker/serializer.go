package main

import (
        "fmt"
		"errors"
        . "github.com/smartfog/fogflow/common/constants"
        . "github.com/smartfog/fogflow/common/ngsi"
)

type Serializer struct{}

func (sz Serializer) SerializeEntity(expanded []interface{}) (LDContextElement, error) {
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
								if v != nil {
									entity.Id = sz.getId(v.(interface{}))
								} else {
									err := errors.New("Id can not be nil!")
									return entity, err
								}
                                break
                        case TYPE:
                                fmt.Println("case @type.......")
								if v != nil {
									entity.Type = sz.getType(v.([]interface{}))
								} else {
									err := errors.New("Type can not be nil!")
									return entity, err
								}
                                break
                        case CREATED_AT:
                                fmt.Println("case createdAt.......")
								if v != nil {
									entity.CreatedAt = sz.getCreatedAt(v.([]interface{}))
								}                                
                                break
                        case LOCATION:
                                fmt.Println("case location.......")
								if v != nil {
									entity.Location = sz.getLocation(v.([]interface{}))
								}                                
                                break
                        default: // default cases like property, relationship here.
                                interfaceArray := v.([]interface{})
                                if len(interfaceArray) > 0 {
                                        mp := interfaceArray[0].(map[string]interface{})
                                        typ := mp[TYPE].([]interface{})
                                        if len(typ) > 0 {
                                                if typ[0].(string) == PROPERTY {
                                                        fmt.Println("It is a property....")
														property, err := sz.getProperty(k, mp)
														if err != nil {
															fmt.Println("Errored Property!")
															return entity, err
														} else {
															entity.Properties = append(entity.Properties, property)
														}
                                                } else if typ[0].(string) == RELATIONSHIP {
                                                        fmt.Println("It is a relationship....")
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
        return entity, nil
}

func (sz Serializer) SerializeRegistration(expanded []interface{}) (CSourceRegistrationRequest, error) {
        fmt.Println("Inside SerializeRegistration.......")
                registration := CSourceRegistrationRequest{}
                for _, val := range expanded {
            fmt.Println("Inside outer for loop.......")
            stringsMap := val.(map[string]interface{})
            for k, v := range stringsMap {
                                fmt.Println("Inside inner for loop.......")
                switch k {
                                        case ID:
												fmt.Println("case ID.......")
												if v != nil { 
													registration.Id = sz.getId(v.(interface{}))
												}
                                                break
                                        case TYPE:
												fmt.Println("case TYPE.......")
												if v != nil {
													registration.Type = sz.getType(v.([]interface{}))
												} else {
													err := errors.New("Type can not be nil!")
													return registration, err
												}
                                                break
                                        case TIMESTAMP:
                                                fmt.Println("case TIMESTAMP.......")
												if v != nil {
													//---------------------registration.Expires = 
													sz.getDateAndTimeValue(v.([]interface{}))
												}
                                                break
                                        case DESCRIPTION:
                                                fmt.Println("case DESCRIPTION.......")
												if v != nil { 
													registration.Description = sz.getStringValue(v.([]interface{}))
												}
                                                break
                                        case ENDPOINT:
                                                fmt.Println("case ENDPOINT.......")
												if v != nil {
													registration.Endpoint = sz.getStringValue(v.([]interface{}))
												} else {
													err := errors.New("Endpoint value can not be nil!")
													return registration, err
												}                                                
                                                break
                                        case EXPIRES:
                                                fmt.Println("case EXPIRES.......")
												if v != nil { 
													registration.Expires = sz.getDateAndTimeValue(v.([]interface{}))
												}
                                                break
                                        case INFORMATION:
                                                fmt.Println("case INFORMATION.......")
												if v != nil {
													fmt.Println("case INFORMATION.......1")
													registration.Information = sz.getInformation(v.([]interface{}))
													fmt.Println("case INFORMATION.......1-done")
												} else {
													fmt.Println("case INFORMATION.......2")
													err := errors.New("Information value can not be nil!")
													return registration, err
												} 
												fmt.Println("case INFORMATION.......3")
                                                break
                                        case LOCATION:
                                                fmt.Println("case LOCATION.......")
												if v != nil { 
													registration.Location = sz.getStringValue(v.([]interface{}))
												}
                                                break
                                        case NAME:
                                                fmt.Println("case NAME.......")
												if v != nil { 
													registration.Name = sz.getStringValue(v.([]interface{}))
												}
                                                break
                                        default:
                                                fmt.Println("case default.......\n", v)
                                                break
                                }
                        }
                }
                return registration, nil
}

func (sz Serializer) SerializeSubscription(expanded []interface{}) (LDSubscriptionRequest, error) {
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
												if v != nil { 
													subscription.Id = sz.getIdFromArray(v.([]interface{}))
												}
                                                break
                                        case TYPE:
                                                fmt.Println("case TYPE.......")
												if v != nil {
													subscription.Type = sz.getType(v.([]interface{}))
												} else {
													err := errors.New("Information value can not be nil!")
													return subscription, err
												}
                                                
                                                break										
                                        case NAME:
                                                fmt.Println("case NAME.......")
												if v != nil { 
													subscription.Name = sz.getValue(v.([]interface{})).(string)
												}
                                                break
                                        case ENTITIES:
                                                fmt.Println("case ENTITIES.......")
												if v != nil { 
													subscription.Entities = sz.getEntities(v.([]interface{}))
												}
                                                break
                                        case DESCRIPTION:
                                                fmt.Println("case DESCRIPTION.......")
												if v != nil { 
													subscription.Description = sz.getValue(v.([]interface{})).(string)
												}
                                                break
                                        case NOTIFICATION:
                                                fmt.Println("case NOTIFICATION.......")
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
                                                fmt.Println("case WATCHED_ATTRIBUTES.......")
												if v != nil { 
													subscription.WatchedAttributes = sz.getArrayOfIds(v.([]interface{}))
												}
                                                break
                                        default:
                                                fmt.Println("case default.......")
                                                break
                                }
                        }
                }
                return subscription, nil
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
		fmt.Println("Inside Property:   1")
        Property := Property{}

		fmt.Println("Inside Property:   2")
		if propertyName != "" {
			Property.Id = propertyName
			Property.Name = propertyName
		}
		fmt.Println("Inside Property:   3")
		Property.Type = PROPERTY
		
		fmt.Println("Inside Property:   4")
		if propertyMap[HAS_VALUE] != nil {
			Property.Value = sz.getValuefromArray(propertyMap[HAS_VALUE].([]interface{}))
		} else {
			err := errors.New("Property Value can not be nil!")
			return Property, err
		}        
		fmt.Println("Inside Property:   5")
		if propertyMap[OBSERVED_AT] != nil {
			Property.ObservedAt = sz.getDateAndTimeValue(propertyMap[OBSERVED_AT].([]interface{}))
		}
        
		fmt.Println("Inside Property:   6")
		if propertyMap[DATASET_ID] != nil {
			Property.DatasetId = sz.getDatasetId(propertyMap[DATASET_ID].([]interface{}))
		}
        
		fmt.Println("Inside Property:   7")
		if propertyMap[INSTANCE_ID] != nil {
			Property.InstanceId = sz.getInstanceId(propertyMap[INSTANCE_ID].([]interface{}))
		}
        
		fmt.Println("Inside Property:   8")
		if propertyMap[CREATED_AT] != nil {
			Property.CreatedAt = sz.getCreatedAt(propertyMap[CREATED_AT].([]interface{}))
		}
        
		fmt.Println("Inside Property:   9")
		if propertyMap[MODIFIED_AT] != nil {
			Property.ModifiedAt = sz.getModifiedAt(propertyMap[MODIFIED_AT].([]interface{}))
		}
        
		fmt.Println("Inside Property:   10")
		if propertyMap[UNIT_CODE] != nil {
			Property.UnitCode = sz.getUnitCode(propertyMap[UNIT_CODE].(interface{}))
		}
        
		fmt.Println("Inside Property:   11")
        return Property, nil
}

func (sz Serializer) getRelationship(relationshipName string, relationshipMap map[string]interface{}) (Relationship, error) {
		fmt.Println("Inside Relationship:   1")
        Relationship := Relationship{}
		
		fmt.Println("Inside Relationship:   2")
		if relationshipName != "" {
			Relationship.Id = relationshipName
			Relationship.Name = relationshipName
		}
        
		fmt.Println("Inside Relationship:   3")
        Relationship.Type = RELATIONSHIP
		
		fmt.Println("Inside Relationship:   4")
		if relationshipMap[HAS_OBJECT] != nil {
			Relationship.Object = sz.getIdFromArray(relationshipMap[HAS_OBJECT].([]interface{}))
		} else {
			err := errors.New("Relationship Object value can not be nil!")
			return Relationship, err
		}          
		
		fmt.Println("Inside Relationship:   5")
		if relationshipMap[OBSERVED_AT] != nil {
			Relationship.ObservedAt = sz.getDateAndTimeValue(relationshipMap[OBSERVED_AT].([]interface{}))
		}
        		
		fmt.Println("Inside Relationship:   6")
		if relationshipMap[PROVIDED_BY] != nil {
			Relationship.ProvidedBy = sz.getProvidedBy(relationshipMap[PROVIDED_BY].([]interface{}))
		}        
		
		fmt.Println("Inside Relationship:   7")
		if relationshipMap[DATASET_ID] != nil {
			Relationship.DatasetId = sz.getDatasetId(relationshipMap[DATASET_ID].([]interface{}))
		}        
		
		fmt.Println("Inside Relationship:   8")
		if relationshipMap[INSTANCE_ID] != nil {
			Relationship.InstanceId = sz.getInstanceId(relationshipMap[INSTANCE_ID].([]interface{}))
		}        
		
		fmt.Println("Inside Relationship:   9")
		if relationshipMap[CREATED_AT] != nil {
			Relationship.CreatedAt = sz.getCreatedAt(relationshipMap[CREATED_AT].([]interface{}))
		}        
		
		fmt.Println("Inside Relationship:   10")
		if relationshipMap[MODIFIED_AT] != nil {
			Relationship.ModifiedAt = sz.getModifiedAt(relationshipMap[MODIFIED_AT].([]interface{}))
		}       
		
		fmt.Println("Inside Relationship:   11")
        return Relationship, nil
}

func (sz Serializer) getValue(hasValue []interface{}) interface{} {
		fmt.Println("Inside getValue:   1")

        Value := make(map[string]interface{})
        if len(hasValue) > 0 {
				fmt.Println("Inside getValue:   2")
                val := hasValue[0].(map[string]interface{})
				fmt.Println("Inside getValue:   3")

                if val[TYPE] != nil {
						fmt.Println("Inside getValue:   4")	
						Value["Type"] = val[TYPE].(string)
						fmt.Println("Inside getValue:   5")		
                        Value["Value"] = val[VALUE].(interface{})
                } else {
						fmt.Println("Inside getValue:   6")
                        Value["Value"] = hasValue[0]
                }
        }
		fmt.Println("Inside getValue:   7")
        return Value
}

func (sz Serializer) getValuefromArray(hasValue []interface{}) interface{} {
		fmt.Println("Inside getValuefromArray:   1")
		Value := make(map[string]interface{})  
		var value interface{}
        if len(hasValue) > 0 {
				fmt.Println("Inside getValuefromArray:   2")
                val := hasValue[0].(map[string]interface{})
				fmt.Println("Inside getValuefromArray:   3")
				
                if val[TYPE] != nil {
						fmt.Println("Inside getValuefromArray:   4")	
						Value["Type"] = val[TYPE].(string)
						fmt.Println("Inside getValuefromArray:   5")
						Value["Value"] = val[VALUE].(interface{})
						return Value
				}
				fmt.Println("Inside getValuefromArray:   6")
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
		fmt.Println("Inside getProvidedBy:   1")
        ProvidedBy := ProvidedBy{}
        if len(providedBy) > 0 {
		fmt.Println("Inside getProvidedBy:   2")
                providedByMap := providedBy[0].(map[string]interface{})
		fmt.Println("Inside getProvidedBy:   3")
                ProvidedBy.Type = sz.getType(providedByMap[TYPE].([]interface{}))
		fmt.Println("Inside getProvidedBy:   4")
                ProvidedBy.Object = sz.getIdFromArray(providedByMap[HAS_OBJECT].([]interface{}))
        }
		fmt.Println("Inside getProvidedBy:   5")
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
		//var coordinates

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
		fmt.Println("Inside getInformation....1")
        regInfoArray := []RegistrationInfo{}
        for _,val := range information {				
				fmt.Println("Inside getInformation....2")
                infoVal := val.(map[string]interface{})
                regInfo := RegistrationInfo{}
				fmt.Println("Inside getInformation....3")
                for k,v := range infoVal {
                        switch k {
                                case PROPERTIES:
										if v != nil { 
											regInfo.Properties = sz.getArrayOfIds(v.([]interface{}))
											fmt.Println("Inside getInformation....4")
										}
										fmt.Println("case properties done")                            
                                        break
                                case RELATIONSHIPS:
										if v != nil { 
											regInfo.Relationships = sz.getArrayOfIds(v.([]interface{}))
											fmt.Println("Inside getInformation....5")
										} 
										fmt.Println("case relationships done")
                                        break
                                case ENTITIES:
										if v != nil { 
											regInfo.Entities = sz.getEntities(v.([]interface{}))
											fmt.Println("Inside getInformation....6")
										} 
										fmt.Println("case entities done")
                                        break
                        }
                }
                regInfoArray = append(regInfoArray, regInfo)
				fmt.Println("Inside getInformation....7")
        }
        return regInfoArray
}

func (sz Serializer) getArrayOfIds(arrayOfIds []interface{}) []string {
        var ArrayOfIds []string
        for _,v := range arrayOfIds {
                idValue := v.(map[string]interface{})
                id := idValue[ID].(string)
                ArrayOfIds = append(ArrayOfIds, id)
        }
        return ArrayOfIds
}

func (sz Serializer) getEntities(entitiesArray []interface{}) []EntityInfo {
		fmt.Println("Inside getEntities.......1")
        entities := []EntityInfo{}
        for _,val := range entitiesArray {
                entityInfo := EntityInfo{}				
				fmt.Println("Inside getEntities.......2")
                entityFields := val.(map[string]interface{})
				
				fmt.Println("Inside getEntities.......3")
                for k,v := range entityFields {
						
						fmt.Println("Inside getEntities.......4")
                        switch k {
                                case ID:										
										fmt.Println("Inside getEntities.......case ID, v is...\n", v)
                                        entityInfo.Id = sz.getId(v.(string))
                                        break
                                case TYPE:										
										fmt.Println("Inside getEntities.......case TYPE, v is...\n", v)
                                        entityInfo.Type = sz.getType(v.([]interface{}))
                                        break
                                case ID_PATTERN:
										fmt.Println("Inside getEntities.......case ID_PATTERN, v is...\n", v)
                                        entityInfo.IdPattern = sz.getStringValue(v.([]interface{}))
                                        break
                        }
                }
                entities = append(entities, entityInfo)
				fmt.Println("Inside getEntities.......5")
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
        for _,val := range notificationArray {
			notificationFields := val.(map[string]interface{})
                for k,v := range notificationFields {
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
                                        fmt.Println("case default.......")
                                        break
                        }
				}
        }
        return notification, nil
}

func (sz Serializer) getEndpoint(endpointArray []interface{}) (Endpoint, error) {
        endpoint := Endpoint{}
        for _,val := range endpointArray {			
			endpointFields := val.(map[string]interface{})
                for k,v := range endpointFields {
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
