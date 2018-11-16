package main

import (
	"math"
	"regexp"
	"strings"

	. "github.com/smartfog/fogflow/common/ngsi"
)

func matchEntityId(entity EntityId, subscribedEntity EntityId) bool {
	if subscribedEntity.IsPattern == true {
		if subscribedEntity.ID != "" {
			matched, _ := regexp.MatchString(subscribedEntity.ID, entity.ID)
			if matched == false {
				return false
			}
		}

		if subscribedEntity.Type != "" {
			matched, _ := regexp.MatchString(subscribedEntity.Type, entity.Type)
			if matched == false {
				return false
			}
		}
	} else {
		if subscribedEntity.Type != "" {
			matched := subscribedEntity.Type == entity.Type && subscribedEntity.ID == entity.ID
			if matched == false {
				return false
			}
		} else {
			matched := subscribedEntity.ID == entity.ID
			if matched == false {
				return false
			}
		}
	}

	return true
}

func matchAttributes(registeredAttributes []ContextRegistrationAttribute, requiredAttributeNames []string) bool {
	for _, attrName := range requiredAttributeNames {
		exist := false
		for _, attribute := range registeredAttributes {
			if attribute.Name == attrName {
				exist = true
				break
			}
		}

		if exist == false {
			return false
		}
	}

	return true
}

func matchMetadatas(metadatas []ContextMetadata, restriction Restriction) bool {
	for _, scope := range restriction.Scopes {
		switch scope.Type {
		case "circle": // check if the location metadata belongs to the circle
			for _, meta := range metadatas {
				if meta.Type == "point" {
					point := meta.Value.(Point)
					circle := scope.Value.(Circle)

					if PointInCircle(&point, &circle) == false {
						return false
					}
				}
			}

		case "polygon": // check if the location metadata belongs to the polygon
			for _, meta := range metadatas {
				if meta.Type == "point" {
					point := meta.Value.(Point)
					polygon := scope.Value.(Polygon)

					if PointInPolygon(&point, &polygon) == false {
						return false
					}
				}
			}

		case "stringQuery": // check if the other metadatas fit the query statement
			queryString := scope.Value.(string)
			constraints := strings.Split(queryString, ";")

			for _, constraint := range constraints {
				items := strings.Split(constraint, "=")
				attrName := items[0]
				attrValue := items[1]

				found := false
				for _, meta := range metadatas {
					if meta.Name == attrName {
						if meta.Value != attrValue {
							return false
						}

						found = true
					}
				}

				if found == false {
					return false
				}
			}
		}
	}

	return true
}

// Returns whether or not the current Polygon contains the passed in Point.
func PointInPolygon(point *Point, polygon *Polygon) bool {
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

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Longitude > end.Longitude {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.Longitude == start.Longitude || point.Longitude == end.Longitude {
		newLng := math.Nextafter(point.Longitude, math.Inf(1))
		point = &Point{Latitude: point.Latitude, Longitude: newLng}
	}

	// If we are outside of the polygon, indicate so.
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

func PointInCircle(point *Point, circle *Circle) bool {

	dist := Distance(point.Latitude, point.Longitude, circle.Latitude, circle.Longitude)

	if dist <= circle.Radius {
		return true
	} else {
		return false
	}
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
