package database

import (
	"database/sql"
	"testing"
)

func TestConnect(*testing.T) {

	database := Database{
		Host:     "localhost",
		Username: "canoe",
		Password: "canoe110930008",
		Database: "canoe",
	}
	conn, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	err = conn.Ping()
	if err != nil {
		panic(err)
	}

}
