package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gorm.io/gorm"
	"net/http"
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
		server := handler.NewChatServer(accessToken)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			return
		}
	})
	// 视频互动
	party.Get("/live/{accessToken:string}", func(ctx *context.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := handler.NewLiveServer(accessToken)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			return
		}
	})

	party.Get("/livex/{accessToken:string}", func(ctx *context.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
		if err != nil {

		}
		for {
			messageType, p, err := conn.ReadMessage()
		}
	})
}
