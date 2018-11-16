package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	. "github.com/smartfog/fogflow/common/config"
)

var DBTables = []string{
	"CREATE TABLE entity_tab (eid text,  type text, isPattern text, providerURL text);",
	"CREATE TABLE attr_tab (eid text, name text, type text, isDomain text);",
	"CREATE TABLE metadata_tab (eid text, name text, type text, value text);",
	"CREATE TABLE geo_box_tab (eid text, name text, type text, box geometry);",
	"CREATE TABLE geo_circle_tab (eid text, name text, type text, center geometry, radius float);"}

// check if the specified database exists
func checkDatabase(dbcfg *DatabaseCfg) (bool, error) {
	connURL := fmt.Sprintf("host=%s  port=%d user=%s password=%s sslmode=disable",
		dbcfg.Host, dbcfg.Port, dbcfg.Username, dbcfg.Password)

	dbconn, err := sql.Open("postgres", connURL)
	if dbconn != nil {
		defer dbconn.Close()
	}

	if err != nil {
		return false, err
	}

	rows, err2 := dbconn.Query("select count(*) from pg_catalog.pg_database where datname = '" + dbcfg.DBname + "'")
	if err2 != nil {
		return false, err2
	}

	bExist := false

	for rows.Next() {
		var count int
		rows.Scan(&count)
		if count == 1 {
			bExist = true
			break
		}
	}

	return bExist, nil
}

// create a new database and all tables
func createDatabase(dbcfg *DatabaseCfg) {
	// (1) create the database
	connURL := fmt.Sprintf("host=%s  port=%d user=%s password=%s sslmode=disable",
		dbcfg.Host, dbcfg.Port, dbcfg.Username, dbcfg.Password)

	dbconn, err := sql.Open("postgres", connURL)
	if err != nil {
		panic(err)
	}

	_, err = dbconn.Exec("CREATE DATABASE " + dbcfg.DBname)
	if err != nil {
		panic(err)
	}

	dbconn.Close()

	// (2) create all tables
	connURL = fmt.Sprintf("host=%s  port=%d user=%s password=%s  dbname=%s sslmode=disable",
		dbcfg.Host, dbcfg.Port, dbcfg.Username, dbcfg.Password, dbcfg.DBname)

	dbconn, err = sql.Open("postgres", connURL)
	if err != nil {
		panic(err)
	}

	_, err = dbconn.Exec("CREATE EXTENSION postgis")
	if err != nil {
		panic(err)
	}

	for _, SQLstatement := range DBTables {
		_, err = dbconn.Exec(SQLstatement)
		if err != nil {
			panic(err)
		}
	}

	dbconn.Close()
}

// delete the database in order to start from scratch
func resetDatabase(dbcfg *DatabaseCfg) {
	connURL := fmt.Sprintf("host=%s  port=%d user=%s password=%s sslmode=disable",
		dbcfg.Host, dbcfg.Port, dbcfg.Username, dbcfg.Password)

	dbconn, err := sql.Open("postgres", connURL)
	if err != nil {
		panic(err)
	}

	_, err = dbconn.Exec("DROP DATABASE " + dbcfg.DBname)
	if err != nil {
		panic(err)
	}

	dbconn.Close()
}

// create a new database and all tables
func openDatabase(dbcfg *DatabaseCfg) *sql.DB {
	connURL := fmt.Sprintf("host=%s  port=%d user=%s password=%s dbname=%s  sslmode=disable",
		dbcfg.Host, dbcfg.Port, dbcfg.Username, dbcfg.Password, dbcfg.DBname)

	dbconn, err := sql.Open("postgres", connURL)
	if err != nil {
		panic(err)
	}

	dbconn.SetMaxIdleConns(10)
	dbconn.SetMaxIdleConns(5)

	return dbconn
}
