package main

import (
	"canoe/internal/model/data"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func connectDB(cfg struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     uint64 `envconfig:"DB_PORT" default:"5432"`
	Username string `envconfig:"DB_USERNAME" default:"canoe"`
	Password string `envconfig:"DB_PASSWORD" default:"canoe110930008"`
	DbName   string `envconfig:"DB_NAME" default:"canoe"`
}) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbName, cfg.Username, cfg.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// migrate
	err = db.AutoMigrate(&data.User{})
	err = db.AutoMigrate(&data.Group{})
	err = db.AutoMigrate(&data.GroupMember{})
	err = db.AutoMigrate(&data.Session{})
	err = db.AutoMigrate(&data.UserSession{})
	err = db.AutoMigrate(&data.Message{})
	if err != nil {
		return nil
	}
	Database = db
	return db
}
