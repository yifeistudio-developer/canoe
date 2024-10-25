package route

import (
	"github.com/kataras/iris/v12"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web"
)

type websocketController struct {
	Socket *web.SocketServer
}

func (c *websocketController) GetChatBy(accessToken string, ctx iris.Context) {
	socket := c.Socket
	server, err := socket.NewWsServer(accessToken, socket.DialMsgHandler)
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
	server, err := socket.NewWsServer(accessToken, web.ChatMsgHandler)
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
