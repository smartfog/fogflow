package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type EntityRepository struct {
	//connection to the backend database
	db *sql.DB

	dbLock sync.RWMutex
}

func (er *EntityRepository) Init(config *DatabaseCfg) {
	var dbExist = false

	for {
		exist, err := checkDatabase(config)

		if err == nil {
			dbExist = exist
			break
		} else {
			ERROR.Println("having some problem to connect to postgresql ", err)
			time.Sleep(2 * time.Second)
		}
	}

	//create the database if not exist
	if dbExist == false {
		createDatabase(config)
	} else {
		if config.DBReset == true {
			resetDatabase(config)
			createDatabase(config)
		}
	}

	//open the database
	er.db = openDatabase(config)

	INFO.Println("connected to postgresql")
}

func (er *EntityRepository) Close() {
	//close the database
	er.db.Close()
	INFO.Println("close the connection to postgresql")
}

func (er *EntityRepository) updateEntity(entity EntityId, registration *ContextRegistration) {
	er.dbLock.Lock()
	defer er.dbLock.Unlock()

	statements := make([]string, 0)

	// update the entity table
	queryStatement := fmt.Sprintf("SELECT entity_tab.eid, entity_tab.type, entity_tab.providerurl FROM entity_tab WHERE eid = '%s'", entity.ID)
	rows, err := er.query(queryStatement)
	if err != nil {
		return
	}
	if rows.Next() == false {
		// insert new entity
		insertEntity := fmt.Sprintf("INSERT INTO entity_tab(eid, type, isPattern, providerURL) VALUES('%s', '%s', '%t', '%s');",
			entity.ID, entity.Type, entity.IsPattern,
			registration.ProvidingApplication)
		statements = append(statements, insertEntity)
	}
	rows.Close()

	// update attribute table
	for _, attr := range registration.ContextRegistrationAttributes {
		queryStatement := fmt.Sprintf("SELECT * FROM attr_tab WHERE attr_tab.eid = '%s' AND attr_tab.name = '%s';",
			entity.ID, attr.Name)
		rows, err := er.query(queryStatement)
		if err == nil {
			if rows.Next() == false {
				// insert as new attribute
				statement := fmt.Sprintf("INSERT INTO attr_tab(eid, name, type, isDomain) VALUES('%s', '%s', '%s', '%t');",
					entity.ID, attr.Name, attr.Type, attr.IsDomain)

				statements = append(statements, statement)
			} else {
				// update as existing attribute
				statement := fmt.Sprintf("UPDATE attr_tab SET type = '%s', isDomain = '%t' WHERE attr_tab.eid = '%s' AND attr_tab.name = '%s';",
					attr.Type, attr.IsDomain, entity.ID, attr.Name)

				statements = append(statements, statement)
			}
		}
		rows.Close()
	}

	// update metadata table
	for _, meta := range registration.Metadata {
		switch strings.ToLower(meta.Type) {
		case "circle":
			circle := meta.Value.(Circle)
			queryStatement := fmt.Sprintf("SELECT * FROM geo_circle_tab WHERE geo_circle_tab.eid = '%s' AND geo_circle_tab.name = '%s';",
				entity.ID, meta.Name)
			rows, err := er.query(queryStatement)
			if err == nil {
				if rows.Next() == false {
					// insert as new attribute
					statement := fmt.Sprintf("INSERT INTO geo_circle_tab(eid, name, type, center, radius) VALUES ('%s', '%s', '%s', ST_SetSRID(ST_MakePoint(%f, %f), 4326), %f);",
						entity.ID, meta.Name, meta.Type, circle.Longitude, circle.Latitude, circle.Radius)
					statements = append(statements, statement)
				} else {
					// update as existing attribute
					statement := fmt.Sprintf("UPDATE geo_circle_tab SET center = ST_SetSRID(ST_MakePoint(%f, %f), 4326) AND radius = %f WHERE geo_circle_tab.eid = '%s' AND geo_circle_tab.name = '%s';",
						circle.Longitude, circle.Latitude, circle.Radius, entity.ID, meta.Name)

					statements = append(statements, statement)
				}
			}
			rows.Close()

		case "point":
			point := meta.Value.(Point)
			queryStatement := fmt.Sprintf("SELECT * FROM geo_box_tab WHERE geo_box_tab.eid = '%s' AND geo_box_tab.name = '%s';",
				entity.ID, meta.Name)
			rows, err := er.query(queryStatement)
			if err == nil {
				if rows.Next() == false {
					// insert as new attribute
					statement := fmt.Sprintf("INSERT INTO geo_box_tab(eid, name, type, box) VALUES ('%s', '%s', '%s', ST_SetSRID(ST_MakePoint(%f, %f), 4326));",
						entity.ID, meta.Name, meta.Type, point.Longitude, point.Latitude)
					statements = append(statements, statement)
				} else {
					// update as existing attribute
					statement := fmt.Sprintf("UPDATE geo_box_tab SET box = ST_SetSRID(ST_MakePoint(%f, %f), 4326) WHERE geo_box_tab.eid = '%s' AND geo_box_tab.name = '%s';",
						point.Longitude, point.Latitude, entity.ID, meta.Name)

					statements = append(statements, statement)
				}
			}
			rows.Close()

		case "polygon":
			polygon := meta.Value.(Polygon)
			locationText := ""
			for k, point := range polygon.Vertices {
				if k > 0 {
					locationText = locationText + ", "
				}
				locationText = locationText + fmt.Sprintf("%f %f", point.Longitude, point.Latitude)
			}

			queryStatement := fmt.Sprintf("SELECT * FROM geo_box_tab WHERE geo_box_tab.eid = '%s' AND geo_box_tab.name = '%s';",
				entity.ID, meta.Name)
			rows, err := er.query(queryStatement)
			if err == nil {
				if rows.Next() == false {
					// insert as new attribute
					statement := fmt.Sprintf("INSERT INTO geo_box_tab(eid, name, type, box) VALUES ('%s', '%s', '%s', ST_MakePolygon(ST_GeomFromText('POLYGON((%s))', 4326)));",
						entity.ID, meta.Name, meta.Type, locationText)
					statements = append(statements, statement)
				} else {
					// update as existing attribute
					statement := fmt.Sprintf("UPDATE geo_box_tab SET box = ST_MakePolygon(ST_GeomFromText('POLYGON((%s))', 4326)) WHERE geo_box_tab.eid = '%s' AND geo_box_tab.name = '%s';",
						locationText, entity.ID, meta.Name)

					statements = append(statements, statement)
				}
			}
			rows.Close()

		default:
			queryStatement := fmt.Sprintf("SELECT * FROM metadata_tab WHERE metadata_tab.eid = '%s' AND metadata_tab.name = '%s';",
				entity.ID, meta.Name)
			rows, err := er.query(queryStatement)
			if err == nil {
				if rows.Next() == false {
					// insert as new attribute
					statement := fmt.Sprintf("INSERT INTO metadata_tab(eid, name, type, value) VALUES('%s', '%s', '%s', '%s');",
						entity.ID, meta.Name, meta.Type, meta.Value)

					statements = append(statements, statement)
				} else {
					// update as existing attribute
					statement := fmt.Sprintf("UPDATE metadata_tab SET type = '%s', value = '%s' WHERE metadata_tab.eid = '%s' AND metadata_tab.name = '%s';",
						meta.Type, meta.Value, entity.ID, meta.Name)

					statements = append(statements, statement)
				}
			}
			rows.Close()
		}
	}

	// apply the update once for the entire registration request, within a transaction
	er.exec(statements)
}

func (er *EntityRepository) queryEntities(entities []EntityId, attributes []string, restriction Restriction) map[string][]EntityId {
	er.dbLock.RLock()
	defer er.dbLock.RUnlock()

	entityMap := make(map[string][]EntityId)

	for _, entity := range entities {
		// three steps to construct the SQL statement to query the result
		queryStatement := "SELECT entity_tab.eid, entity_tab.type, entity_tab.ispattern, entity_tab.providerurl FROM entity_tab "

		// (1) consider attribute list
		for i, attr := range attributes {
			queryStatement = queryStatement + fmt.Sprintf(" INNER JOIN attr_tab at%d  ON entity_tab.eid = at%d.eid AND at%d.name = '%s' ",
				i+1, i+1, i+1, attr)
		}

		// (2) apply scopes to metadata
		boxTabFilter := ""
		orderBy := ""
		var num_of_geo_scopes int
		num_of_geo_scopes = 0
		for _, scope := range restriction.Scopes {
			switch strings.ToLower(scope.Type) {
			case "nearby":
				nearby := scope.Value.(NearBy)
				orderBy = fmt.Sprintf("  ST_Distance(geo_box_tab.box, ST_SetSRID(ST_MakePoint(%f, %f), 4326)) LIMIT %d ",
					nearby.Longitude, nearby.Latitude, nearby.Limit)

			case "circle":
				circle := scope.Value.(Circle)
				if num_of_geo_scopes == 0 {
					boxTabFilter = boxTabFilter + "( "
				} else {
					boxTabFilter = boxTabFilter + " OR "
				}
				boxTabFilter = boxTabFilter + fmt.Sprintf(" ST_DWithin(geo_box_tab.box, ST_SetSRID(ST_MakePoint(%f, %f), 4326), %f, true) ",
					circle.Longitude, circle.Latitude, circle.Radius)
				num_of_geo_scopes = num_of_geo_scopes + 1

			case "simplegeolocation":
				value := scope.Value.(Segment)
				segment := value.Converter()
				locationText := ""
				locationText = locationText + fmt.Sprintf("%f %f,", segment.NW_Corner.Longitude, segment.NW_Corner.Latitude)
				locationText = locationText + fmt.Sprintf("%f %f,", segment.NW_Corner.Longitude, segment.SE_Corner.Latitude)
				locationText = locationText + fmt.Sprintf("%f %f,", segment.SE_Corner.Longitude, segment.SE_Corner.Latitude)
				locationText = locationText + fmt.Sprintf("%f %f,", segment.SE_Corner.Longitude, segment.NW_Corner.Latitude)
				locationText = locationText + fmt.Sprintf("%f %f", segment.NW_Corner.Longitude, segment.NW_Corner.Latitude)

				if num_of_geo_scopes == 0 {
					boxTabFilter = boxTabFilter + "( "
				} else {
					boxTabFilter = boxTabFilter + " OR "
				}
				boxTabFilter = boxTabFilter + fmt.Sprintf(" ST_Within(geo_box_tab.box, ST_GeomFromText('POLYGON((%s))', 4326)) ",
					locationText)
				num_of_geo_scopes = num_of_geo_scopes + 1

			case "polygon":
				polygon := scope.Value.(Polygon)
				locationText := ""
				for k, point := range polygon.Vertices {
					if k > 0 {
						locationText = locationText + ", "
					}
					locationText = locationText + fmt.Sprintf("%f %f", point.Longitude, point.Latitude)
				}
				if num_of_geo_scopes == 0 {
					boxTabFilter = boxTabFilter + "( "
				} else {
					boxTabFilter = boxTabFilter + " OR "
				}
				boxTabFilter = boxTabFilter + fmt.Sprintf(" ST_Within(geo_box_tab.box, ST_GeomFromText('POLYGON((%s))', 4326)) ",
					locationText)
				num_of_geo_scopes = num_of_geo_scopes + 1

			case "stringquery":
				queryString := scope.Value.(string)
				constraints := strings.Split(queryString, ";")
				for i, constraint := range constraints {
					items := strings.Split(constraint, "=")
					queryStatement = queryStatement + fmt.Sprintf(" INNER JOIN metadata_tab md%d ON entity_tab.eid = md%d.eid and md%d.name = '%s' and md%d.value = '%s' ",
						i+1, i+1, i+1, items[0], i+1, items[1])
				}
			}
		}

		// (3) apply geo-scopes
		if boxTabFilter != "" {
			queryStatement = queryStatement + fmt.Sprintf(" INNER JOIN geo_box_tab ON entity_tab.eid = geo_box_tab.eid and %s) ", boxTabFilter)
		} else if orderBy != "" {
			queryStatement = queryStatement + fmt.Sprintf(" INNER JOIN geo_box_tab ON entity_tab.eid = geo_box_tab.eid ")
		}

		// (4) consider entity_id
		if entity.IsPattern == true {
			if entity.Type != "" && entity.ID != "" {
				queryStatement = queryStatement + fmt.Sprintf(" WHERE entity_tab.eid like '%s' AND entity_tab.type like '%s'",
					strings.Replace(entity.ID, ".*", "%", -1), strings.Replace(entity.Type, ".*", "%", -1))
			} else if entity.Type != "" {
				queryStatement = queryStatement + fmt.Sprintf(" WHERE entity_tab.type like '%s'",
					strings.Replace(entity.Type, ".*", "%", -1))
			} else if entity.ID != "" {
				queryStatement = queryStatement + fmt.Sprintf(" WHERE entity_tab.eid like '%s' ",
					strings.Replace(entity.ID, ".*", "%", -1))
			}
		} else {
			queryStatement = queryStatement + fmt.Sprintf(" WHERE entity_tab.eid = '%s' ", entity.ID)
		}

		// (5) consider sorting based on geo-distance
		if orderBy != "" {
			queryStatement = queryStatement + fmt.Sprintf(" ORDER BY %s ", orderBy)
		}

		DEBUG.Println(queryStatement)

		// perform the query
		rows, err := er.query(queryStatement)
		if err != nil {
			return nil
		}

		// prepare the result according the returned dataset
		for rows.Next() {
			var eid, etype, ispattern, providerURL string
			rows.Scan(&eid, &etype, &ispattern, &providerURL)

			var bIsPattern bool
			if ispattern == "true" {
				bIsPattern = true
			} else {
				bIsPattern = false
			}
			e := EntityId{ID: eid, Type: etype, IsPattern: bIsPattern}
			entityMap[providerURL] = append(entityMap[providerURL], e)
		}
		rows.Close()
	}

	return entityMap
}

func (er *EntityRepository) deleteEntity(eid string) {
	er.dbLock.Lock()
	defer er.dbLock.Unlock()

	fmt.Println("==delete entity ", eid)

	// find out the associated entity
	queryStatement := fmt.Sprintf("SELECT entity_tab.eid, entity_tab.type, entity_tab.providerurl FROM entity_tab WHERE eid = '%s'", eid)

	// perform the query
	rows, err := er.query(queryStatement)
	if err != nil {
		return
	}

	statements := make([]string, 0)

	for rows.Next() {
		var entityID, entityType, providerURL string
		rows.Scan(&entityID, &entityType, &providerURL)

		if entityType == "IoTBroker" {
			fmt.Println("IoT Broker left as a context provider")
			er.ProviderLeft(providerURL)
		}

		// remove all attributes related to this entity
		executeStatement := fmt.Sprintf("DELETE FROM attr_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)

		// remove all metadata related to this entity
		executeStatement = fmt.Sprintf("DELETE FROM metadata_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)

		// remove all geo-metadata related to this entity
		executeStatement = fmt.Sprintf("DELETE FROM geo_box_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)

		executeStatement = fmt.Sprintf("DELETE FROM geo_circle_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)
	}
	rows.Close()

	// remove the entity
	executeStatement := fmt.Sprintf("DELETE FROM entity_tab WHERE eid = '%s'", eid)
	statements = append(statements, executeStatement)

	er.exec(statements)
}

func (er *EntityRepository) ProviderLeft(providerURL string) {
	// find out all entities associated with this broker
	queryStatement := fmt.Sprintf("SELECT entity_tab.eid FROM entity_tab WHERE providerurl = '%s'", providerURL)

	fmt.Println(queryStatement)

	// perform the query
	rows, err := er.query(queryStatement)
	if err != nil {
		return
	}

	statements := make([]string, 0)

	for rows.Next() {
		var entityID string
		rows.Scan(&entityID)

		// remove all attributes related to this entity
		executeStatement := fmt.Sprintf("DELETE FROM attr_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)

		// remove all metadata related to this entity
		executeStatement = fmt.Sprintf("DELETE FROM metadata_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)

		// remove all geo-metadata related to this entity
		executeStatement = fmt.Sprintf("DELETE FROM geo_box_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)
		executeStatement = fmt.Sprintf("DELETE FROM geo_circle_tab WHERE eid = '%s'", entityID)
		statements = append(statements, executeStatement)
	}
	rows.Close()

	// remove all entities related to this registration
	executeStatement := fmt.Sprintf("DELETE FROM entity_tab WHERE providerurl = '%s'", providerURL)
	statements = append(statements, executeStatement)

	er.exec(statements)
}

func (er *EntityRepository) retrieveRegistration(entityID string) *ContextRegistration {
	er.dbLock.RLock()
	defer er.dbLock.RUnlock()

	// query all entities associated with this registrationId
	queryStatement := fmt.Sprintf("SELECT eid, type, isPattern, providerURL FROM entity_tab WHERE entity_tab.eid = '%s';", entityID)

	rows, err := er.query(queryStatement)
	if err != nil {
		ERROR.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var eid, etype, epattern, providerURL string
		rows.Scan(&eid, &etype, &epattern, &providerURL)

		ctxRegistration := ContextRegistration{}

		entities := make([]EntityId, 0)

		entity := EntityId{}
		entity.ID = eid
		entity.Type = etype

		if epattern == "true" {
			entity.IsPattern = true
		} else {
			entity.IsPattern = false
		}

		entities = append(entities, entity)

		ctxRegistration.EntityIdList = entities
		ctxRegistration.ProvidingApplication = providerURL

		// query all attributes that belong to those entities
		registeredAttributes := make([]ContextRegistrationAttribute, 0)

		queryStatement = fmt.Sprintf("SELECT name, type, isDomain FROM attr_tab WHERE attr_tab.eid = '%s';", eid)
		results, _ := er.query(queryStatement)
		for results.Next() {
			var name, attributeType, isDomain string
			results.Scan(&name, &attributeType, &isDomain)

			attr := ContextRegistrationAttribute{}
			attr.Name = name
			attr.Type = attributeType

			if isDomain == "true" {
				attr.IsDomain = true
			} else {
				attr.IsDomain = false
			}

			registeredAttributes = append(registeredAttributes, attr)
		}
		results.Close()

		ctxRegistration.ContextRegistrationAttributes = registeredAttributes

		// query all metadatas that belong to those entities
		registeredMetadatas := make([]ContextMetadata, 0)

		queryStatement = fmt.Sprintf("SELECT name, type, value FROM metadata_tab WHERE metadata_tab.eid = '%s';", eid)
		results, _ = er.query(queryStatement)
		for results.Next() {
			var name, mdType, value string
			results.Scan(&name, &mdType, &value)

			metadata := ContextMetadata{}
			metadata.Name = name
			metadata.Type = mdType
			metadata.Value = value

			registeredMetadatas = append(registeredMetadatas, metadata)
		}
		results.Close()

		// query all geo-related metadatas that belong to those entities
		queryStatement = fmt.Sprintf("SELECT name, ST_AsText(box) FROM geo_box_tab WHERE geo_box_tab.eid = '%s';", eid)
		results, _ = er.query(queryStatement)
		for results.Next() {
			var name, mtype, box string
			results.Scan(&name, &box)

			metadata := ContextMetadata{}
			metadata.Name = name
			metadata.Type = mtype

			switch mtype {
			case "point":
				fmt.Println(box)
				var latitude, longitude float64
				_, err := fmt.Scanf(box, "POINT(%f %f)", &longitude, &latitude)
				if err == nil {
					point := Point{}
					point.Latitude = latitude
					point.Longitude = longitude

					metadata.Value = point
				} else {
					metadata.Type = "string"
					metadata.Value = box
				}

				fmt.Println("point")
			}

			metadata.Value = box

			registeredMetadatas = append(registeredMetadatas, metadata)
		}
		results.Close()

		queryStatement = fmt.Sprintf("SELECT name, ST_AsText(center), radius FROM geo_circle_tab WHERE geo_circle_tab.eid = '%s';", eid)
		results, _ = er.query(queryStatement)
		for results.Next() {
			var name, mtype, center string
			var radius float64
			results.Scan(&name, &center, &radius)

			metadata := ContextMetadata{}
			metadata.Name = name
			metadata.Type = mtype

			circle := Circle{}

			var latitude, longitude float64
			_, err := fmt.Scanf(center, "POINT(%f %f)", &longitude, &latitude)
			if err == nil {
				circle.Latitude = latitude
				circle.Longitude = longitude
				circle.Radius = radius

				metadata.Value = circle
			} else {
				metadata.Type = "string"
				metadata.Value = fmt.Sprintf("%s, %f", center, radius)
			}

			fmt.Println("circle")

			metadata.Value = circle

			registeredMetadatas = append(registeredMetadatas, metadata)
		}
		results.Close()

		ctxRegistration.Metadata = registeredMetadatas

		return &ctxRegistration
	}

	return nil
}

func (er *EntityRepository) query(statement string) (*sql.Rows, error) {
	return er.db.Query(statement)
}

func (er *EntityRepository) exec(statements []string) {
	for _, statement := range statements {
		er.db.Exec(statement)
	}
}
