package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	User *UserAdapter
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(postgres.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", openErr)
	}
	err := db.AutoMigrate(
		&User{},
		&UserSession{},
		&Group{},
		&GroupMember{},
		&Session{},
		&Message{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate order: %w", err)
	}
	return &Adapter{
		User: &UserAdapter{db},
	}, nil

}
