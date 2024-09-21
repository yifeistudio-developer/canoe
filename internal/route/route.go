package route

import (
	. "canoe/internal/service"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var logger *golog.Logger

func SetupRoutes(app *iris.Application) {
	logger = app.Logger()

	// root api
	api := app.Party("/canoe/api")
	root := mvc.New(api).
		Register(NewWebSocketService(), NewUserService(), NewSessionService()).
		Handle(new(rootController))

	// users api
	root.Party("/users").
		Handle(new(userController))

	// sessions api
	root.Party("/sessions").
		Handle(new(sessionController))

	// websocket api
	root.Party("/ws").
		Handle(new(webSocketController))

}
