package route

import (
	"canoe/internal/service"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var logger *golog.Logger

func SetupRoutes(app *iris.Application) {

	logger = app.Logger()

	// root api
	api := app.Party("/canoe/api")
	mvc.New(api).
		Handle(&rootController{})

	// users api
	userApi := api.Party("/users")
	mvc.New(userApi).
		Register().
		Handle(&userController{})

	// sessions api
	sessionApi := api.Party("/sessions")
	mvc.New(sessionApi).
		Handle(&sessionController{})

	// websocket api
	wsApi := api.Party("/ws")
	mvc.New(wsApi).
		Register(service.NewWebSocketService()).
		Handle(new(webSocketController))

}
