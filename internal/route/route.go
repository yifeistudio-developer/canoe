package route

import (
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
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

	// 普通聊天
	party.Get("/chat/{accessToken:string}", func(ctx *context.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := wsServer(accessToken)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			return
		}
	})

	// 视频互动
	party.Get("/live/{accessToken:string}", func(ctx *context.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := wsServer(accessToken)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			return
		}
	})
}
