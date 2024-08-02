package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"github.com/kataras/iris/v12"
)

func SetupRoutes(app *iris.Application) {
	party := app.Party("/canoe")
	party.Get("/", func(ctx iris.Context) {
		_ = ctx.JSON(Success(nil))
	})
	party.Get("/channel/{accessToken:string}", handler.HandleChannel())
}
