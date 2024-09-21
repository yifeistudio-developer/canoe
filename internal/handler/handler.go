package handler

import (
	. "canoe/internal/model"
	"canoe/internal/remote"
	"context"
	"encoding/json"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"net/http"
)

const (
	WebSocket = iota
	WebRTC
)

var logger *golog.Logger

func SetLogger(l *golog.Logger) {
	logger = l
	initUDP()
}

type MsgHandler func(ctx context.Context,
	cancel context.CancelFunc,
	conn *neffos.NSConn,
	username string,
	message neffos.Message) error

var HandlerMap = make(map[uint]func(ctx context.Context, cancel context.CancelFunc, conn *neffos.NSConn, username string, message neffos.Message) error)

func init() {
	HandlerMap[WebSocket] = handleWebSocket
	HandlerMap[WebRTC] = handleWebRTC
}

func handleWebSocket(ctx context.Context, cancel context.CancelFunc, conn *neffos.NSConn, username string, message neffos.Message) error {
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
	ctx, cancel := context.WithCancel(context.Background())
	profile := remote.GetUserProfile(accessToken)
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	ws := websocket.New(upgrader, websocket.Events{websocket.OnNativeMessage: func(conn *neffos.NSConn, message neffos.Message) error {
		return handler(ctx, cancel, conn, profile.Username, message)
	}})

	// 当连接建立
	// 初始化用户会话信息
	ws.OnConnect = func(conn *neffos.Conn) error {
		logger.Infof("got new connection: access-token = %s", accessToken)
		return nil
	}

	// 清理会话信息
	ws.OnDisconnect = func(c *neffos.Conn) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("recover from error: %v", r)
			}
		}()
		peers.Delete(profile.Username)
		cancel()
		logger.Infof("disconnected: access-token = %s", accessToken)
	}
	return ws
}
