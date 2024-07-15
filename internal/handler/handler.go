package handler

import (
	. "canoe/internal/model"
	"encoding/json"
	"errors"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

func HandleChannel() context.Handler {
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(conn *neffos.NSConn, message neffos.Message) error {
			var msg Message
			if err := json.Unmarshal(message.Body, &msg); err != nil {
				rlt := Result{Code: 400, Msg: "bad request: message format is illegal."}
				str, _ := json.Marshal(rlt)
				conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, str))
				return err
			}
			return nil
		},
	})
	// 连接建立
	ws.OnConnect = func(c *neffos.Conn) error {
		socket := c.Socket()
		request := socket.Request()
		uri := request.RequestURI
		accessToken := uri[strings.LastIndex(uri, "/"):]
		if accessToken == "" {
			c.Close()
			return errors.New("access-token is empty")
		}
		if strings.HasSuffix(accessToken, "/") {
			accessToken = accessToken[1:]
		}
		return nil
	}
	// 连接断开
	ws.OnDisconnect = func(c *neffos.Conn) {

	}
	return websocket.Handler(ws)
}
