package route

import (
	"github.com/kataras/iris/v12"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/middleware"
)

type websocketController struct {
	Socket *middleware.WebsocketServer
}

func (c *websocketController) GetChatBy(accessToken string, ctx iris.Context) {
	socket := c.Socket
	server, err := socket.Handle(accessToken, socket.DialMsgHandler)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, err = server.Upgrade(
		ctx.ResponseWriter(),
		ctx.Request(),
		nil,
		nil,
	)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}

func (c *websocketController) GetDialBy(accessToken string, ctx iris.Context) {
	socket := c.Socket
	server, err := socket.Handle(accessToken, middleware.ChatMsgHandler)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, err = server.Upgrade(
		ctx.ResponseWriter(),
		ctx.Request(),
		nil,
		nil,
	)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}
