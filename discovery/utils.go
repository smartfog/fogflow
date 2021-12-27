package main

import (
	. "fogflow/common/ngsi"
	"math"
	"regexp"
	"strings"
	"fmt"
)

func matchingWithFilters(registration *EntityRegistration, idFilter []EntityId, attrFilter []string, metaFilter Restriction, subFiwareService string, regFiwareService string) bool {

	/*if regFiwareService != "" && subFiwareService != "" && subFiwareService != regFiwareService {
                return false
        }

	// (1) check entityId part
	entity := EntityId{}
	fmt.Println("metaFilter",metaFilter)
	if strings.HasPrefix(registration.Type, "https://uri.etsi.org/ngsi-ld/default-context/") {
		entity.Type = registration.Type
	} else {
		entity.Type = registration.Type
	}
	entity.ID = registration.ID
	//entity.Type = registration.Type
	entity.IsPattern = false
	atLeastOneMatched := false
	for _, tmp := range idFilter {
		matched := matchEntityId(entity, tmp)
		if matched == true {
			atLeastOneMatched = true
			break
		}
	}
	if atLeastOneMatched == false {
		return false
	}

	// (2) check attribute set
	if matchAttributes(registration.AttributesList, attrFilter) == false {
		return false
	}*/

	// (3) check metadata set
	if metaFilter.RestrictionType == "ld" {
		if matchLdMetadatas(registration.MetadataList, metaFilter) == false {
			return false 
		}
	} else { 
		if matchMetadatas(registration.MetadataList, metaFilter) == false {
			return false
		}
	}
	// if all matched, return true
	return true
}


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

func matchAttributes(registeredAttributes map[string]ContextRegistrationAttribute, requiredAttributeNames []string) bool {
	for _, attrName := range requiredAttributeNames {
		if _, exist := registeredAttributes[attrName]; exist == false {
			return false
		}
	}

	return true
}

//matchLdMetadatas

func matchLdMetadatas(metadatas map[string]ContextMetadata, restriction Restriction) bool {
	sp := restriction.Geometry
	//sc := restriction.Cordinates
	fmt.Println("metadatas",metadatas)
	if value , ok := metadatas["location"]; ok == true {
		fmt.Println("value",value)
	}
	switch strings.ToLower(sp) {
		case "point":
			if restriction.Georel != "" {
				gr := restriction.Georel
				contrains := strings.Split(gr,";")
				if len(contrains) > 1 {
					sws := strings.ReplaceAll(contrains[1], " ", "")
					minMax := strings.Split(sws,"==")
					//dist := distance()
					if minMax[0] == "maxDistance" {
						maxDistance := minMax[1]
						fmt.Println("maxDistance",maxDistance)
					} else if minMax[0] == "minDistance" {
						minDistance := minMax[1]
						fmt.Println("minDistance",minDistance)
					} else {
						fmt.Println(minMax[1])
					}
				}
			}
		case "polygon": 
		default :
			fmt.Println("waitting to implement")
	}
	return true
}
func matchMetadatas(metadatas map[string]ContextMetadata, restriction Restriction) bool {
	for _, scope := range restriction.Scopes {
		switch strings.ToLower(scope.Type) {
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

		case "stringquery": // check if the other metadatas fit the query statement
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
	center := Point{}
	center.Longitude = circle.Longitude
	center.Latitude = circle.Latitude

	dist := Distance(point, &center)

	if dist <= uint64(circle.Radius) {
		return true
	} else {
		return false
	}
}
