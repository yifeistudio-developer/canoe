package route

import (
	. "canoe/internal/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func SetupRoutes(app *iris.Application) {
	party := app.Party("/canoe/api")
	party.Get("/", func(ctx iris.Context) {
		_ = ctx.JSON(Success(nil))
	})
	users := party.Party("/users")
	UserRoutes(users)
	sessions := party.Party("/sessions")
	SessionRoutes(sessions)
	// 普通聊天
	party.Get("/chat/{accessToken:string}", handleWs)
	// 视频互动
	party.Get("/live/{accessToken:string}", handleWs)
}

func handleWs(ctx *context.Context) {
	accessToken := ctx.Params().GetString("accessToken")
	server := wsServer(accessToken)
	conn, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
	defer conn.Close()
	if err != nil {
		return
	}
}
