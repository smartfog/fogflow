package main

import (
	"strings"
	"fmt"
)

type cType func(interface{})interface{}

func getEntityId(id interface{}) string {
	Id := id.(string)
	return Id
}

func getType(typ []interface{}) interface{} {
	return typ
}


func MultiLineStringHandler(coordinates interface{}) interface{} {
	var value interface{}
	return value
}

func LineStringHandler(coordinates interface{})interface{} {
	var value interface{}
        return value
}

func MultiPolygonHandler(coordinates interface{}) interface{} {
	var value interface{}
        return value
}

func PolygonHandler(coordinates interface{})interface{} {
	var value interface{}
        return value

}

// MultiPoint Handler

func MultiPointHandler(coordinates interface{}) interface{} {
        coordinatesArray := coordinates.([]interface{})
	Points := make([]interface{},0)
	for index, _ := range coordinatesArray {
		point := PointHandler(coordinatesArray[index])
		Points = append(Points,point)
	}
	return Points
}

// Point Handler

func PointHandler(coordinates interface{})interface{} {
	coordinateArray := coordinates.([]interface{})
	point := make([]interface{},0)
	for _, value := range coordinateArray {
		p  := value.(map[string]interface{})
		point = append(point,p["@value"])

	}
	fmt.Println("point",point)
        return point

}
func getCoordinateType(typ []interface{}) cType {
	coordinateType := typ[0].(string)
	var  functionType cType
	fmt.Println("coordinateType",coordinateType)
	if strings.Contains(coordinateType, "MultiPoint"){
		functionType = MultiPointHandler
	} else if strings.Contains(coordinateType, "MultiLineString") {
		functionType = MultiLineStringHandler
	} else if strings.Contains(coordinateType, "LineString") {
		functionType = LineStringHandler
	} else if strings.Contains(coordinateType, "MultiPolygon") {
		functionType = MultiPolygonHandler
	} else if strings.Contains(coordinateType, "Polygon") {
		functionType = PolygonHandler
	} else {
		functionType = PointHandler
	}
	fmt.Println("coordinatetype",functionType)
	return functionType
}

func getValue(val []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	var coordinateHandler cType
	valueMap := val[0].(map[string]interface{})
	if val, ok := valueMap["@type"]; ok {
		coordinateHandler = getCoordinateType(val.([]interface{}))
	}
	for key,val := range valueMap {
		if strings.Contains(key, "@type") {
			result[key] = getType(val.([]interface{}))
		} else if strings.Contains(key, "coordinates") {
			result[key] = coordinateHandler(val.(interface{}))
		}
	}
    return result
}
