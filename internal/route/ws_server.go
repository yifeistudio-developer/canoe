package route

import (
	. "canoe/internal/model"
	"encoding/json"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/pion/webrtc/v4"
	"log"
	"net/http"
)

func wsServer(accessToken string) *neffos.Server {
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})

	ws := websocket.New(upgrader, websocket.Events{
		// 普通消息
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
		// webrtc 消息
		"ice-candidate": func(conn *neffos.NSConn, msg neffos.Message) error {
			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(msg.Body, &candidate); err != nil {
				return err
			}
			peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
			if err != nil {
				return err
			}
			return peerConnection.AddICECandidate(candidate)
		},
		// webrtc 消息
		"sdp": func(conn *neffos.NSConn, msg neffos.Message) error {
			var offer webrtc.SessionDescription
			if err := json.Unmarshal(msg.Body, &offer); err != nil {
				return err
			}
			peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
			if err != nil {
				return err
			}
			// 设置远程描述
			if err := peerConnection.SetRemoteDescription(offer); err != nil {
				return err
			}
			// 创建应答
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				return err
			}
			// 设置本地描述
			if err := peerConnection.SetLocalDescription(answer); err != nil {
				return err
			}
			// 发送应答
			answerJSON, err := json.Marshal(peerConnection.LocalDescription())
			if err != nil {
				log.Println("Failed to marshal answer:", err)
				return nil
			}
			conn.Emit("sdp", answerJSON)
			return nil
		},
	})
	// 连接建立
	ws.OnConnect = func(conn *neffos.Conn) error {
		println(accessToken)
		return nil
	}
	// 连接断开
	ws.OnDisconnect = func(c *neffos.Conn) {

	}
	return ws
}
