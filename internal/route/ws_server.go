package route

import (
	. "canoe/internal/model"
	"canoe/internal/remote"
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

func onOfferMsg(conn *neffos.NSConn, msg *SignalingMessage) error {
	pc, err := w.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return err
	}
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType:  webrtc.MimeTypeOpus,
		ClockRate: 48000,
		Channels:  2,
	}, "audio", "pion-audio")
	if err != nil {
		return err
	}
	_, err = pc.AddTrack(audioTrack)
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType:  webrtc.MimeTypeVP8,
		ClockRate: 90000,
	}, "video", "pion-video")
	if err != nil {
		return err
	}
	_, err = pc.AddTrack(videoTrack)
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
		// 消息转播
		fmt.Printf("接收到轨道: %s\n", track.Kind().String())
		fmt.Println(track.Codec().RTPCodecCapability)
		go func() {
			buf := make([]byte, 1500)
			for {
				i, _, readErr := track.Read(buf)
				if readErr != nil {
					panic(readErr)
				}
				if track.Kind() == webrtc.RTPCodecTypeAudio {
					_, err = audioTrack.Write(buf[:i])
				} else if track.Kind() == webrtc.RTPCodecTypeVideo {
					_, err = videoTrack.Write(buf[:i])
				}
				if err != nil {
				}
			}
		}()
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
	return nil
}

// 处理视频消息
func handleLiveMsg(conn *neffos.NSConn, message neffos.Message) error {
	var msg SignalingMessage
	err := message.Unmarshal(&msg)
	if err != nil {
		return err
	}
	switch msg.Type {
	case "offer":
		err := onOfferMsg(conn, &msg)
		if err != nil {
			return err
		}
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

	// 当连接建立
	// 初始化用户回话信息
	ws.OnConnect = func(conn *neffos.Conn) error {
		logger.Infof("got new connection: access-token = %s", accessToken)
		profile := remote.GetUserProfile(accessToken)
		fmt.Println(profile)
		return nil
	}

	// 清理回话信息
	ws.OnDisconnect = func(c *neffos.Conn) {
		logger.Infof("disconnected: access-token = %s", accessToken)
	}
	return ws
}
