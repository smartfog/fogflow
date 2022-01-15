package geolocation 

import (
	"math"
	"strings"
	"fmt"
	. "fogflow/common/ngsi"
	. "fogflow/common/constants"
)

const (
	earthRadiusMi = 3958
	earthRaidusKm = 6371
)


type Coord struct {
	Lat float64
	Lon float64
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}


func LDDistance(p, q Coord) (mi, km float64) {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	mi = c * earthRadiusMi
	km = c * earthRaidusKm
	return mi, km
}

func DistForPoint(typ string, metadatas interface{},  loc interface{}) (float64, float64){
	var mi , km float64
	var lag, log float64
	fmt.Println("1-Point")
	switch metadatas.(type) {
		case []float64:
			meta := metadatas.([]interface{})
			if len(meta) == 2 {
				lag = meta[0].(float64)
				log = meta[1].(float64)
			}
		case Point :
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
			mi , km = LDDistance(p1,p2)
        }
        return mi, km
}

func FindDistForPoint(typ string , metaData interface{}, loc interface{}) (float64, float64) {
	var mi, km float64
	fmt.Println("1-point")
	switch strings.ToLower(typ) {
		case "point" :
			mi, km = DistForPoint(typ,metaData,loc)
		case "polygon":
		default : 
	}
	return mi, km
}

type geoPoint struct {
	x float64
	y float64
}

func onSegment(p, q,  r geoPoint ) bool { 
    if (q.x <= math.Max(p.x, r.x) && q.x >= math.Min(p.x, r.x) &&  q.y <= math.Max(p.y, r.y) && q.y >= math.Min(p.y, r.y))  {
        return true; 
     }
    return false; 
} 

func orientation(p,q, r geoPoint) int {
	val := (q.y - p.y) * (r.x - q.x) -
			(q.x - p.x) * (r.y - q.y);

	if val == 0 {
	    return 0
	}
	var res int 
	if val > 0 {
		res = 1
	} else {
		res = 2
	}
	return res 
}



func doIntersect(p1, q1, p2, q2 geoPoint) bool {
	 o1 := orientation(p1, q1, p2)
	 o2 := orientation(p1, q1, q2)
	 o3 := orientation(p2, q2, p1)
	 o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}

	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	if o3 == 0 && onSegment(p2, p1, q2) {
		 return true
	}

	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

func isInside(polygon []geoPoint, n int , p geoPoint) bool {
	if n < 3 {
		return false 
	}
	extreme := geoPoint {
			x: INF,
			y: p.y,
		}
	count := 0
	i := 0
	for i !=0 {
		next := (i +1 )%n
		if doIntersect(polygon[i], polygon[next], p, extreme) {
			if orientation(polygon[i], p, polygon[next]) == 0 {
				return onSegment(polygon[i], p, polygon[next])
			}
		count  = count +1
		}
		i = next
	}
	var res bool 
	if count&1 == 1 {
		res = true
	} else {
		res = false
	}
	return res
}

func convertInStructure(coor interface{}) []geoPoint {
	coorA := coor.([]interface{})
	res := make([]geoPoint, 0)
	for _ , val := range coorA {
		valA := val.([]interface{})
		var geo geoPoint
		geo.x = valA[0].(float64)
		geo.y = valA[1].(float64)
		res = append(res, geo)
	}
	return res
}

func commonConverter(entityP interface{}, queryP interface{}) ([]geoPoint,[]geoPoint) {
	entityMeta := make([]geoPoint,0)
	queryMeta := make([]geoPoint,0)
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
	return entityMeta,queryMeta
}

func checkEquals(entityP interface{}, queryP interface{}) bool {
	entityMeta,queryMeta := commonConverter(entityP,queryP)
	fmt.Println(entityMeta)
	fmt.Println(queryMeta)
	/*if entityP == queryP {
		fmt.Println("matched")
		return true 
	}*/
	return false
}

func checkDisjoint(entityP interface{}, queryP interface{}) bool {
	entityMeta,queryMeta := commonConverter(entityP,queryP)
	var  disjoint bool
	disjoint = true
	for _ , val := range queryMeta {
		size := len(entityMeta)
		status := isInside(entityMeta,size,val)
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
	for _ , val := range queryMeta {
                size := len(entityMeta)
                status := isInside(entityMeta,size,val)
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
        for _ , val := range entityMeta {
                size := len(queryMeta)
                status := isInside(queryMeta,size,val)
                if status == false {
                        contain = false
                        break
                }
        }
        return contain
}

func checkPoint(meta Point, queryP interface{}) bool {
	metaPoint :=geoPoint{}
	metaPoint.x = meta.Latitude
	metaPoint.y = meta.Longitude
	_,queryMeta := commonConverter(nil,queryP)
	//inside := false
	//fmt.Println("MetaPoint, queryMeta",metaPoint,queryMeta)
	//for _, val := range queryMeta {
	size := len(queryMeta)
	inside := isInside(queryMeta,size,metaPoint)
		/*if status == true {
			inside = true
			break 
		}*/
	return inside

}
func FindDistForPolygon(typ string , metaData interface{}, res Restriction) (bool) {
	//var mi, km float64
	var status bool
	geoRel := strings.ReplaceAll(res.Georel, " ", "")
	switch strings.ToLower(typ) {
		case "point":
			status = checkPoint(metaData.(Point),res.Cordinates)
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
