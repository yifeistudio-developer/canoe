package handler

import (
	. "canoe/internal/model"
	"canoe/internal/remote"
	"encoding/json"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"net/http"
)

const (
	WebSocket = iota
	WebRTC
)

type MsgHandler func(conn *neffos.NSConn, username string, message neffos.Message) error

var HandlerMap = make(map[uint]func(conn *neffos.NSConn, username string, message neffos.Message) error)

func init() {
	HandlerMap[WebSocket] = handleWebSocket
	HandlerMap[WebRTC] = handleWebRTC
}

func handleWebSocket(conn *neffos.NSConn, username string, message neffos.Message) error {
	var evp Envelope
	err := message.Unmarshal(&evp)
	if err != nil {
		rlt := Result{Code: 400, Msg: "bad request: message format is illegal."}
		str, _ := json.Marshal(rlt)
		conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, str))
		return err
	}
	payload := evp.Payload
	str := payload.(string)
	conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, []byte(str)))
	return nil
}

// 处理websocket 连接

func HandleWS(accessToken string, handler MsgHandler) *neffos.Server {
	profile := remote.GetUserProfile(accessToken)
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	ws := websocket.New(upgrader, websocket.Events{websocket.OnNativeMessage: func(conn *neffos.NSConn, message neffos.Message) error {
		return handler(conn, profile.Username, message)
	}})

	// 当连接建立
	// 初始化用户会话信息
	ws.OnConnect = func(conn *neffos.Conn) error {
		logger.Infof("got new connection: access-token = %s", accessToken)
		return nil
	}

	// 清理会话信息
	ws.OnDisconnect = func(c *neffos.Conn) {
		delete(peers, profile.Username)
		logger.Infof("disconnected: access-token = %s", accessToken)
	}
	return ws
}
