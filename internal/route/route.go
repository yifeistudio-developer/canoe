package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

var logger *golog.Logger

var chatService *service.ChatService

var sessionService *service.SessionService

func SetupRoutes(app *iris.Application) {
	logger = app.Logger()
	handler.SetLogger(logger)
	chatService = service.NewChatService()
	sessionService = service.NewSessionService()
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
		server := handler.HandleWS(accessToken, handler.HandlerMap[handler.WebSocket])
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			logger.Errorf("websocket upgrade failed: %v", err)
			return
		}
	})

	// 视频/直播互动
	party.Get("/live/{accessToken:string}", func(ctx iris.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := handler.HandleWS(accessToken, handler.HandlerMap[handler.WebRTC])
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			logger.Errorf("websocket upgrade failed: %v", err)
			return
		}
	})
}
