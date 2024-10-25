package route

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web"
	"github.com/yifeistudio-developer/canoe/internal/application/core/api"
)

func Register(party iris.Party, app *api.Application) {
	root := mvc.New(party)
	root.Party("/users").
		Register(app.User).
		Handle(new(userController))
	root.Party("/ws").
		Register(web.NewSocketServer()).
		Handle(new(websocketController))
}
