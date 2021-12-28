package geolocation 

import (
	"math"
	"strings"
	. "fogflow/common/ngsi"
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

func FindDist(typ string, metadatas interface{},  loc interface{}) (float64, float64){
	var mi , km float64
	var lag, log float64
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


