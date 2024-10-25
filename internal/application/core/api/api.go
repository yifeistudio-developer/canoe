package api

import (
	"github.com/yifeistudio-developer/canoe/internal/adapters/db"
)

type Application struct {
	User *UserService
}

func NewApplication(dbAdapter *db.Adapter) *Application {
	return &Application{
		User: &UserService{db: dbAdapter.User},
	}
}
