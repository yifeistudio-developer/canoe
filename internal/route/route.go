package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func SetupRoutes(app *iris.Application, db *gorm.DB) {
	party := app.Party("/canoe/api")

	party.Get("/", func(ctx iris.Context) {
		_ = ctx.JSON(Success(nil))
	})

	users := party.Party("/users")
	UserRoutes(users, &service.UserService{Db: db})

	sessions := party.Party("/sessions")
	SessionRoutes(sessions, &service.SessionService{Db: db})

	party.Get("/channel/{accessToken:string}", handler.HandleChannel())
}
