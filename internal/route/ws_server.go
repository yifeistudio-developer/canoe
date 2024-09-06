package route

import (
	. "canoe/internal/model"
	"encoding/json"
	"fmt"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/pion/webrtc/v4"
	"net/http"
	"sync"
)

var (
	w = webrtc.NewAPI(webrtc.WithSettingEngine(webrtc.SettingEngine{
		// 配置 WebRTC 的设置
	}))
	peers = make(map[string]*webrtc.PeerConnection)
	mu    sync.Mutex
)

// SignalingMessage 信令消息
type SignalingMessage struct {
	Type      string                  `json:"type"`
	SDP       string                  `json:"sdp,omitempty"`
	Candidate webrtc.ICECandidateInit `json:"candidate,omitempty"`
}

func handleLiveMsg(conn *neffos.NSConn, message neffos.Message) error {
	var msg SignalingMessage
	err := message.Unmarshal(&msg)
	if err != nil {
		return err
	}
	switch msg.Type {
	case "offer":
		pc, err := w.NewPeerConnection(webrtc.Configuration{})
		if err != nil {
			return err
		}
		peers["client"] = pc
		offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: msg.SDP}
		if err := pc.SetRemoteDescription(offer); err != nil {
			return err
		}
		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			return err
		}
		if err := pc.SetLocalDescription(answer); err != nil {
			return err
		}
		answerMsg := SignalingMessage{
			Type: webrtc.SDPTypeAnswer.String(),
			SDP:  pc.LocalDescription().SDP,
		}
		resp, err := json.Marshal(answerMsg)
		if err != nil {
			return err
		}
		pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
			fmt.Printf("handle track message: %s \n", track.Codec().MimeType)
			localTrack, err2 := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
			if err2 != nil {
				fmt.Printf("get track err: %s", err2)
			}
			rtpSender, err2 := pc.AddTrack(localTrack)
			if err2 != nil {
				fmt.Println("add track error:", err2)
			}
			go func() {
				for {
					rtpPacket, _, err := track.ReadRTP()
					if err != nil {
						fmt.Printf("handle error: %s\n", err)
					}
					err = localTrack.WriteRTP(rtpPacket)
					if err != nil {
						fmt.Printf("write error: %s\n", err)
					}
				}
			}()
			defer func(rtpSender *webrtc.RTPSender) {
				err := rtpSender.Stop()
				if err != nil {
					fmt.Printf("stop err: %s\n", err)
				}
			}(rtpSender)
		})
		pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			fmt.Println("ICE Connection State has changed. ", state.String())
		})
		pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate == nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			candidateJson := candidate.ToJSON()
			candidateMsg := SignalingMessage{
				Type:      "candidate",
				Candidate: candidateJson,
			}
			candidateStr, _ := json.Marshal(candidateMsg)
			conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, candidateStr))
		})
		conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, resp))
	case "candidate":
		if pc, ok := peers["client"]; ok {
			candidate := webrtc.ICECandidateInit{Candidate: msg.Candidate.Candidate}
			if err := pc.AddICECandidate(candidate); err != nil {
				return err
			}
		}
	}
	return nil
}

// 处理普通聊天消息
func handleChatMsg(conn *neffos.NSConn, message neffos.Message) error {
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
func wsServer(accessToken string, handler neffos.MessageHandlerFunc) *neffos.Server {
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	ws := websocket.New(upgrader, websocket.Events{websocket.OnNativeMessage: handler})
	// 连接建立
	ws.OnConnect = func(conn *neffos.Conn) error {
		return nil
	}
	// 连接断开
	ws.OnDisconnect = func(c *neffos.Conn) {
	}
	return ws
}
