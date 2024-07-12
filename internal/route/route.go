package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"github.com/kataras/iris/v12"
)

func SetupRoutes(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(Success(nil))
	})
	app.Get("/channel/{accessToken:string}", handler.HandleChannel())
}
