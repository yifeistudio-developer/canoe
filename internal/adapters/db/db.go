package db

import (
	"fmt"
	"github.com/yifeistudio-developer/canoe/config"
	"github.com/yifeistudio-developer/canoe/internal/ports"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	user *UserAdapter
}

func (a *Adapter) UserDbPort() ports.UserDbPort {
	return a.user
}

func NewAdapter() (*Adapter, error) {
	dataSourceURL := config.GetDataSourceURL()
	db, openErr := gorm.Open(postgres.Open(dataSourceURL), &gorm.Config{})
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
		user: &UserAdapter{db},
	}, nil

}
