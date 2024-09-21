package route

import (
	. "canoe/internal/service"
	"github.com/kataras/iris/v12"
)

type webSocketController struct {
	Ws *WebSocketService
}

// 聊天

func (c *webSocketController) GetChatBy(accessToken string, ctx iris.Context) {
	service := c.Ws
	server := service.NewWsServer(accessToken, ChatMsgHandler)
	_, err := server.Upgrade(
		ctx.ResponseWriter(),
		ctx.Request(),
		nil,
		nil,
	)
	if err != nil {
		logger.Errorf("websocket service upgrade err: %v", err)
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}

// 视频通话

func (c *webSocketController) GetDialBy(accessToken string, ctx iris.Context) {
	service := c.Ws
	server := service.NewWsServer(accessToken, service.DialMsgHandler)
	_, err := server.Upgrade(
		ctx.ResponseWriter(),
		ctx.Request(),
		nil,
		nil,
	)
	if err != nil {
		logger.Errorf("websocket service upgrade err: %v", err)
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}
