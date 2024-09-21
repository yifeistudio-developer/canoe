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
		Register().
		Handle(new(rootController))

	// users api
	root.Party("/users").
		Register(NewUserService()).
		Handle(new(userController))

	// sessions api
	root.Party("/sessions").
		Register(NewSessionService()).
		Handle(new(sessionController))

	// websocket api
	root.Party("/ws").
		Register(NewWebSocketService()).
		Handle(new(webSocketController))

}
