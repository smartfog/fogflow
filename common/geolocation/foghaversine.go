package geolocation 

import (
	"math"
	"strings"
	"fmt"
	. "fogflow/common/ngsi"
	//. "fogflow/common/constants"
)

const (
	earthRadiusMi = 3958
	earthRaidusKm = 6371
)


type Coord struct {
	Lat float64
	Lon float64
}

/*This routine calculates the distance between two points (given the    
latitude/longitude of those points). It is being used to calculate    
the distance between two locations using GeoDataSource (TM) products */

//The main implementation is here https://www.geodatasource.com/developers/go

func LDDistance(p, q Coord, unit ...string) float64 {
	const PI float64 = 3.141592653589793
	lat1 := p.Lat
	lat2 := q.Lat
	lng1 := p.Lon
	lng2 := q.Lon
	fmt.Println(lat1,lat2, lng1, lng2)
	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)
	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)
	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func DistForPoint(typ string, metadatas interface{},  loc interface{}) float64{
	fmt.Println("metadatas",metadatas)
	fmt.Println("loc",loc)
	var mi  float64
	var lag, log float64
	switch metadatas.(type) {
		case []interface{}:
			fmt.Println("comming in this part--1")
			meta := metadatas.([]interface{})
			if len(meta) == 2 {
				lag = meta[0].(float64)
				log = meta[1].(float64)
			}
		case Point :
			fmt.Println("Comming in this part ---2")
			meta := metadatas.(Point)
			lag = meta.Latitude
			log = meta.Longitude
	}
        switch strings.ToLower(typ) {
                case "point":
                        locn := loc.([]interface{})
			p1 := Coord {
				Lat : lag,
				Lon : log,
			}
			p2 := Coord{
                                Lat: locn[0].(float64),
                                Lon: locn[1].(float64),
                        }
			fmt.Println(p1,p2)
			mi  = LDDistance(p1,p2, "M")
        }
        return mi
}

func FindDistForPoint(typ string , metaData interface{}, loc interface{}) float64 {
	var mi float64
	fmt.Println("1-point")
	switch strings.ToLower(typ) {
		case "point" :
			mi = DistForPoint(typ,metaData,loc)
		case "polygon":
		default : 
	}
	return mi
}

type geoPoint struct {
	Latitude float64
	Longitude float64
}

type Poly struct {
	Vertices []geoPoint `json:"vertices"`
}

func convertInStructure(coor interface{}) Poly {
	coorA := coor.([]interface{})
	res := make([]geoPoint, 0)
	for _ , val := range coorA {
		valA := val.([]interface{})
		var geo geoPoint
		geo.Latitude = valA[0].(float64)
		geo.Longitude = valA[1].(float64)
		res = append(res, geo)
	}
	polygon := Poly{}
	polygon.Vertices = res
	return polygon
}

func commonConverter(entityP interface{}, queryP interface{}) (Poly,Poly) {
	entityMeta := Poly{}
	queryMeta := Poly{}
	if entityP != nil {
		entityA := entityP.([]interface{})
		for _ , val := range entityA {
			entityMeta = convertInStructure(val)
		}
	}
	if queryP != nil {
		queryA := queryP.([]interface{})
		for _ , val := range queryA {
			queryMeta = convertInStructure(val)
		}
	}
	return entityMeta, queryMeta
}

func checkEquals(entityP interface{}, queryP interface{}) bool {
	entityMeta,queryMeta := commonConverter(entityP,queryP)
	fmt.Println(entityMeta)
	fmt.Println(queryMeta)
	equal := true
	if len(entityMeta.Vertices) != len(queryMeta.Vertices) {
		return false
	}
	hashMap := make(map[geoPoint]bool)
	for _ , val := range queryMeta.Vertices {
		hashMap[val] = true
	}
	for _ , val := range entityMeta.Vertices {
		if hashMap[val] == false {
			return false
		}
	}
	return equal
}



func checkDisjoint(entityP interface{}, queryP interface{}) bool {
	entityMeta,queryMeta := commonConverter(entityP,queryP)
	var  disjoint bool
	disjoint = true
	for _ , val := range queryMeta.Vertices {
		//size := len(entityMeta.Vertices)
		status := isInside(&val,&entityMeta)
		if status == true {
			disjoint = false
			break
		}
	}
	return disjoint
}

func checkWithin(entityP interface{}, queryP interface{}) bool {
	entityMeta,queryMeta := commonConverter(entityP,queryP)
	within := true
	for _ , val := range queryMeta.Vertices {
                //size := len(entityMeta)
                status := isInside(&val,&entityMeta)
                if status == false {
                        within = false
                        break
                }
        }
        return within
}

func checkContains(entityP interface{}, queryP interface{}) bool {
        entityMeta,queryMeta := commonConverter(entityP,queryP)
        contain:= true
        for _ , val := range entityMeta.Vertices {
		status := isInside(&val ,&queryMeta)
                if status == false {
                        contain = false
                        break
                }
        }
        return contain
}

func checkPoint(meta geoPoint, queryP interface{}) bool {
	metaPoint :=geoPoint{}
	metaPoint.Latitude = meta.Latitude
	metaPoint.Longitude = meta.Longitude
	_,queryMeta := commonConverter(nil,queryP)
	inside := isInside(&metaPoint,&queryMeta)
	return inside

}

func FindDistForPolygon(typ string , metaData interface{}, res Restriction) (bool) {
	//var mi, km float64
	var status bool
	geoRel := strings.ReplaceAll(res.Georel, " ", "")
	switch strings.ToLower(typ) {
		case "point":
			status = checkPoint(metaData.(geoPoint),res.Cordinates)
		case "polygon":
			if geoRel == "equals" {
				status = checkEquals(metaData,res.Cordinates)
				fmt.Println(status)
			} else if geoRel == "disjoint" {
				status =  checkDisjoint(metaData,res.Cordinates)
			} else if geoRel == "intersects" {
				fmt.Println("To be implemented latter")
			} else if geoRel == "within" {
				status = checkWithin(metaData,res.Cordinates)
			} else if geoRel == "contains" {
				status = checkContains(metaData,res.Cordinates)
			} else if geoRel == "overlaps" {
			}
		default:
	}
	return status
}

// Returns whether or not the current Polygon contains the passed in Point.
func isInside(point *geoPoint, polygon *Poly) bool {
	start := len(polygon.Vertices) - 1
	end := 0

	contains := intersectsWithRaycast(point, &polygon.Vertices[start], &polygon.Vertices[end])

	for i := 1; i < len(polygon.Vertices); i++ {
		if intersectsWithRaycast(point, &polygon.Vertices[i-1], &polygon.Vertices[i]) {
			contains = !contains
		}
	}

	return contains
}

func intersectsWithRaycast(point *geoPoint, start *geoPoint, end *geoPoint) bool {
	if start.Longitude > end.Longitude {

		start, end = end, start

	}

	for point.Longitude == start.Longitude || point.Longitude == end.Longitude {
		newLng := math.Nextafter(point.Longitude, math.Inf(1))
		point = &geoPoint{Latitude: point.Latitude, Longitude: newLng}
	}

	if point.Longitude < start.Longitude || point.Longitude > end.Longitude {
		return false
	}

	if start.Latitude > end.Latitude {
		if point.Latitude > start.Latitude {
			return false
		}
		if point.Latitude < end.Latitude {
			return true
		}

	} else {
		if point.Latitude > end.Latitude {
			return false
		}
		if point.Latitude < start.Latitude {
			return true
		}
	}

	raySlope := (point.Longitude - start.Longitude) / (point.Latitude - start.Latitude)
	diagSlope := (end.Longitude - start.Longitude) / (end.Latitude - start.Latitude)

	return raySlope >= diagSlope
}

