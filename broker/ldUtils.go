package main

import (
	"strings"
	//"fmt"
)

type cType func(interface{}) interface{}

func getEntityId(id interface{}) string {
	Id := id.(string)
	return Id
}

func getType(typ []interface{}) interface{} {
	return typ
}

/*func collectfourPoints(frist, seecond ,third , forth interface{}) interface{}{
	points := make([]interface{},0)
	points = append(points,frist)
	points = append(points,seecond)
	points = append(points,third)
	points = append(points,forth)
	return points
}
func MultiLineStringHandler(coordinates interface{}) interface{} {
	mulLinePoints := make([]interface{},0)
	coordinateArray := coordinates.([]interface{})
	for p:=0 ; p<len(coordinateArray); p = p+4 {
		points := collectfourPoints(coordinateArray[p],coordinateArray[p+1],coordinateArray[p+2],coordinateArray[p+3])
		linePoints := MultiPointHandler(points)
		mulLinePoints = append(mulLinePoints,linePoints)
	}
	fmt.Println("mulLinePoints",mulLinePoints)
	return mulLinePoints
}

func LineStringHandler(coordinates interface{})interface{} {
	points := MultiPointHandler(coordinates)
        return points
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
	for p := 0; p < len(coordinatesArray) ; p = p + 2 {
		point := pointExtract(coordinatesArray[p],coordinatesArray[p + 1])
		Points = append(Points,point)
	}
	return Points
}

// Point Handler

func pointExtract(latitude, logitude interface{}) []interface{}{
	point := make([]interface{},0)
	latitudeMap := latitude.(map[string]interface{})
	logitudeMap := logitude.(map[string]interface{})
	point = append(point,latitudeMap["@value"])
	point = append(point,logitudeMap["@value"])
	return point
}
func PointHandler(coordinates interface{})interface{} {
	coordinateArray := coordinates.([]interface{})
	points := pointExtract(coordinateArray[0],coordinateArray[1])
        return points

}*/
/*func getCoordinateType(typ []interface{}) cType {
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
}*/

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
	typeArray  := typ.([]interface{})
	return typeArray[0].(string)
}

func reslice(slice []string, s int) []string {
	slice[s] = slice[len(slice)-1]
        return slice[:len(slice)-1]
}

func checkCondition(k string) bool {
	if k == "@id" {
		return false
	} else if  k == "@type" {
		return false
	} else if k == "@context"{
		return false
	} else if strings.Contains(k, "modifiedAt") {
		return false
	} else if strings.Contains(k, "createdAt") {
                return false
        } else if strings.Contains(k, "observationSpace") {
                return false
        } else if strings.Contains(k, "operationSpace") {
                return false
        } else {
		return true
	}
}  
