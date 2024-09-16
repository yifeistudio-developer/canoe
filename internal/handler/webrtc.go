package handler

import (
	"encoding/json"
	"github.com/kataras/neffos"
	"github.com/pion/webrtc/v4"
	"sync"
)

// SignalingMessage 信令消息
type SignalingMessage struct {
	Type      string                  `json:"type"`
	SDP       string                  `json:"sdp,omitempty"`
	Candidate webrtc.ICECandidateInit `json:"candidate,omitempty"`
	Intent    string                  `json:"intent,omitempty"`
}

var (
	peers = make(map[string]*webrtc.PeerConnection)
	mu    sync.Mutex
)

func handleWebRTC(conn *neffos.NSConn, username string, message neffos.Message) error {
	var msg SignalingMessage
	err := message.Unmarshal(&msg)
	if err != nil {
		return err
	}
	pc, err := initPeerConnection()
	if err != nil {
		return err
	}
	peers[username] = pc
	// Allow us to receive 1 audio track, and 1 video track
	if _, err := pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		return err
	} else if _, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return err
	}
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		intent := msg.Intent
		go func() {
			if intent == "_anyone_" {
				processLive(username, track, pc)
			} else if intent != "" {
				processDialog(username, track)
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
	switch msg.Type {
	case "offer":
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
		conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, resp))
	case "candidate":
		if pc, ok := peers[username]; ok {
			candidate := webrtc.ICECandidateInit{Candidate: msg.Candidate.Candidate}
			if err := pc.AddICECandidate(candidate); err != nil {
				return err
			}
		}
	}
	return nil
}

// 初始化连接

func initPeerConnection() (*webrtc.PeerConnection, error) {
	// Create a MediaEngine object to configure the supported codec
	m := webrtc.MediaEngine{}
	// Setup the codecs you want to use.
	h264Codec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeH264,
			ClockRate: 90000,
		},
	}
	err := m.RegisterCodec(h264Codec, webrtc.RTPCodecTypeVideo)
	if err != nil {
		return nil, err
	}
	opusCodec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeOpus,
			ClockRate: 48000,
		},
	}
	err = m.RegisterCodec(opusCodec, webrtc.RTPCodecTypeAudio)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	return api.NewPeerConnection(config)
}

//视频通话

func processDialog(username string, track *webrtc.TrackRemote) {

}
