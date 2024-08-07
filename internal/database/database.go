package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(cfg struct {
	Host     string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port     uint64 `envconfig:"DATABASE_PORT" default:"5432"`
	Username string `envconfig:"DATABASE_USERNAME" default:"canoe"`
	Password string `envconfig:"DATABASE_PASSWORD" default:"canoe110930008"`
	DbName   string `envconfig:"DATABASE_NAME" default:"canoe"`
}) *sql.DB {
	param := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbName, cfg.Username, cfg.Password)
	conn, err := sql.Open("postgres", param)
	if err != nil {
		panic(err.Error())
	}
	err = conn.Ping()
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func Disconnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err.Error())
	}
}
