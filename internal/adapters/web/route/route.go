package route

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/middleware"
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

func Register(party iris.Party, app ports.ApiPort) {
	index := mvc.New(party)
	index.Handle(new(indexController))
	index.Party("/users").
		Register(app.UserApiPort()).
		Handle(new(userController))
	index.Party("/ws").
		Register(middleware.NewSocketServer()).
		Handle(new(websocketController))
}

type indexController struct {
}

func (*indexController) Get() domain.Result {
	return domain.SuccessNil()
}
