package main

import (
	"strings"
	"errors"
	. "fogflow/common/constants"
)

type cType func(interface{}) interface{}

func getEntityId(id interface{}) string {
	Id := id.(string)
	return Id
}

func getType(typ interface{}) interface{} {
	return typ
}

func getObservedTime(observedTime []interface{}) interface{} {
	return observedTime
}

func getCreatedTime(createdTime []interface{}) interface{} {
	return createdTime
}

func getModifiedTime(modifiedTime []interface{}) interface{} {
	return modifiedTime
}

func getDataSetID(dataSetID []interface{}) interface{} {
	return dataSetID
}
func getUniCode(dataSetID []interface{}) interface{} {
	return dataSetID
}

func getInstanceID(instanceID []interface{}) interface{} {
	return instanceID
}

func getGeoValue(val []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	valueMap := val[0].(map[string]interface{})
	for key, val := range valueMap {
		if strings.Contains(key, "@type") {
			result[key] = getType(val.([]interface{}))
		} else if strings.Contains(key, "coordinates") {
			result[key] = val
		}
	}
	return result
}

func getPropertyValue(val []interface{}) interface{} {
	return val
}

//getType

func getRegistrationType(typ interface{}) string {
	var result string
	typeArray := make([]interface{}, 0)
	switch typ.(type) {
	case []interface{}:
		typeArray = typ.([]interface{})
		result = typeArray[0].(string)
	case string:
		result = typ.(string)
	}
	return result
}

func reslice(slice []string, s int) []string {
	slice[s] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func checkCondition(k string) bool {
	if k == "@id" {
		return false
	} else if k == "@type" {
		return false
	} else if k == "@context" {
		return false
	} else if k == NGSI_LD_MODIFIEDAT {
		return false
	} else if k == NGSI_LD_CREATEDAT {
		return false
	} else if k == NGSI_LD_OBSERVED_AT {
		return false
	} else if k == NGSI_LD_OPERATIONSPACE {
		return false
	} else {
		return true
	}
}

/*
	check aatribute type
*/

func checkAttributeType(typ interface{}) string {
	attrType := getRegistrationType(typ)
	var typeResult string
	if attrType == LD_GEOPROPERTY {
		typeResult = "GeoProperty"
	} else if attrType == LD_RELATIONSHIP {
		typeResult = "Relationship"
	} else if attrType == LD_PRPERTY {
		typeResult = "Property"
	} else {
		typeResult = ""
	}
	return typeResult
}

// update entity

func update(prev, curr map[string]interface{}) map[string]interface{} {
	for key, value := range curr {
		if _, ok := prev[key]; ok == true {
			switch value.(type) {
			case map[string]interface{}:
				valueMap := value.(map[string]interface{})
				if value, ok := valueMap[NGSI_LD_TYPE]; ok == true {
					typ := checkAttributeType(value)
					if typ != "" {
						newPrev := prev[key].(map[string]interface{})
						newCurr := curr[key].(map[string]interface{})
						prev[key] = update(newPrev, newCurr)
					} else {
						prev[key] = curr[key]
					}
				} else {
					prev[key] = curr[key]
				}
			default:
				prev[key] = curr[key]
			}
		} else {
			prev[key] = value
		}
	}
	return prev
}

func getLDobject(attr, context interface{}) map[string]interface{} {
	ldAttributes := make(map[string]interface{})
	ldAttributes[attr.(string)] = ""
	ldAttributes["@context"] = context
	return ldAttributes
}

func getAttribute(attributes interface{}) string {
	attrs := attributes.([]interface{})
	attr := attrs[0].(map[string]interface{})
	for key, _ := range attr {
		return key
	}
	return ""
}

/*
	Subscription related functions
*/

func getSubscriptionID(id interface{}) string {
	Id := id.(string)
	return Id
}

func getSubscriptionType(typ interface{}) (interface{}, error) {
	var err error
	switch typ.(type) {
	case []interface{}:
		subTyp := typ.([]interface{})
		if subTyp[0].(string) != NGSILD_SUBSCRIPTION {
			err = errors.New("Type not allowed!")
		}
	case string:
		if typ != NGSILD_SUBSCRIPTION {
			err = errors.New("Type not allowed!")
		}
	default:
	}
	return typ, err
}

func getWatchedAttribute(expWattr []interface{}) ( []string ,error){
	var err error
	wAttr := make([]string,0)
	if len(expWattr) > 0 {
		watchedAttr := expWattr[0].(map[string]interface{})
		for _, value := range watchedAttr {
			wAttr = append(wAttr,value.(string))
		}
	} else {
		 err = errors.New("Zero leanth watched attribute not allowed!")
	}
	return wAttr, err
}

func getStringValue(value []interface{}) string {
        var Value string
        if len(value) > 0 {
                val := value[0].(map[string]interface{})
                Value = val["@value"].(string)
        }
        return Value
}


