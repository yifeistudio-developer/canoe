package handler

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"
)

func HandleChannel() context.Handler {
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{})
	return websocket.Handler(ws)
}
