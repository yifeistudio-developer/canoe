package route

import (
	. "canoe/internal/model"
	"github.com/kataras/iris/v12"
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
	party.Get("/chat/{accessToken:string}", func(ctx iris.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := wsServer(accessToken, handleChatMsg)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		//defer conn.Close()
		if err != nil {
			return
		}
	})
	// 视频互动
	party.Get("/live/{accessToken:string}", func(ctx iris.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := wsServer(accessToken, handleLiveMsg)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		//defer conn.Close()
		if err != nil {
			return
		}
	})
}
