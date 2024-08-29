package handler

import (
	. "canoe/internal/model"
	"encoding/json"
	"errors"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"net/http"
	"strings"
)

func HandleChannel() context.Handler {
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	ws := websocket.New(upgrader, websocket.Events{
		websocket.OnNativeMessage: func(conn *neffos.NSConn, message neffos.Message) error {
			var msg Message
			if err := json.Unmarshal(message.Body, &msg); err != nil {
				rlt := Result{Code: 400, Msg: "bad request: message format is illegal."}
				str, _ := json.Marshal(rlt)
				conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, str))
				return err
			}
			payload := msg.Payload
			str := payload.(string)
			conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, []byte(str)))
			return nil
		},
	})
	// 连接建立
	ws.OnConnect = func(conn *neffos.Conn) error {
		socket := conn.Socket()
		request := socket.Request()
		uri := request.RequestURI
		accessToken := uri[strings.LastIndex(uri, "/"):]
		if accessToken == "" {
			conn.Close()
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
