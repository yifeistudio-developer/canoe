package api

import (
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

type Application struct {
	user *UserApi
}

func (app *Application) UserApiPort() ports.UserApiPort {
	return app.user
}

func NewApplication(dbPort ports.DbPort) *Application {
	return &Application{
		user: &UserApi{db: dbPort.UserDbPort()},
	}
}
