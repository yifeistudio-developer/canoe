package middleware

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"time"

	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
)

type WebsocketServer struct {
	peers  *sync.Map
	ctx    context.Context
	cancel context.CancelFunc
}

type SignalingMessage struct {
	Type      string                  `json:"type"`
	SDP       string                  `json:"sdp,omitempty"`
	Candidate webrtc.ICECandidateInit `json:"candidate,omitempty"`
	Intent    string                  `json:"intent,omitempty"`
}
type udpConn struct {
	conn *net.UDPConn
	port int
}

type MsgHandler func(profile domain.AlpsUserProfile, conn *neffos.NSConn, msg neffos.Message) error

var udpConns = map[webrtc.RTPCodecType]*udpConn{
	webrtc.RTPCodecTypeAudio: {port: 4000},
	webrtc.RTPCodecTypeVideo: {port: 4002},
}

func ChatMsgHandler(profile domain.AlpsUserProfile, conn *neffos.NSConn, msg neffos.Message) error {
	var evp domain.Envelope
	err := msg.Unmarshal(&evp)
	if err != nil {
		rlt := domain.Result{Code: 400, Msg: "bad request: message format is illegal."}
		str, _ := json.Marshal(rlt)
		conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, str))
		return err
	}
	payload := evp.Payload
	str := payload.(string)
	conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, []byte(str)))
	return nil
}

func (s *WebsocketServer) DialMsgHandler(profile domain.AlpsUserProfile, conn *neffos.NSConn, msg neffos.Message) error {
	var sm SignalingMessage
	peers := s.peers
	err := msg.Unmarshal(&sm)
	if err != nil {
		return err
	}
	pc, err := initPeerConnection()
	if err != nil {
		return err
	}
	peers.Store(profile.Username, pc)
	// Allow us to receive 1 audio track, and 1 video track
	if _, err := pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		return err
	} else if _, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return err
	}
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		intent := sm.Intent
		go func() {
			if intent == "_anyone_" {
				s.handleLive(profile.Username, track, pc)
			} else if intent != "" {
				s.handleDialog(intent)
			}
		}()
	})
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		candidateJson := candidate.ToJSON()
		candidateMsg := SignalingMessage{
			Type:      "candidate",
			Candidate: candidateJson,
		}
		candidateStr, _ := json.Marshal(candidateMsg)
		conn.Conn.Write(conn.Conn.DeserializeMessage(neffos.TextMessage, candidateStr))
	})
	switch sm.Type {
	case "offer":
		offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sm.SDP}
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
		if value, ok := peers.Load(profile.Username); ok {
			candidate := webrtc.ICECandidateInit{Candidate: sm.Candidate.Candidate}
			pc := value.(*webrtc.PeerConnection)
			if err := pc.AddICECandidate(candidate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *WebsocketServer) Handle(token string, handler MsgHandler) (*neffos.Server, error) {
	// todo get alps user profile.
	upgrader := gorilla.Upgrader(grl.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }})
	ws := websocket.New(upgrader, websocket.Events{websocket.OnNativeMessage: func(conn *neffos.NSConn, msg neffos.Message) error {
		return handler(domain.AlpsUserProfile{}, conn, msg)
	}})
	ws.OnConnect = func(conn *neffos.Conn) error {
		return nil
	}
	ws.OnDisconnect = func(c *neffos.Conn) {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		// todo username
		s.peers.Delete("")
		s.cancel()
	}
	return ws, nil
}

func initPeerConnection() (*webrtc.PeerConnection, error) {
	m := webrtc.MediaEngine{}
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

func (s *WebsocketServer) handleDialog(username string) {
	peers := s.peers
	value, ok := peers.Load(username)
	if !ok {
		return
	}
	peer := value.(*webrtc.PeerConnection)
	if peer == nil {
		return
	}
}

func initUDP() {
	var laddr *net.UDPAddr
	var err error
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		return
	}
	for _, c := range udpConns {
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			return
		}
		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			return
		}
	}
}

func (s *WebsocketServer) handleLive(username string, track *webrtc.TrackRemote, pc *webrtc.PeerConnection) {
	streamURL := fmt.Sprintf("%s/%s", "rtmp://localhost:1935/stream", username)
	err := s.startFFmpeg(streamURL)
	if err != nil {
		s.cancel()
		return
	}
	initUDP()
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			if rtcpErr := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); rtcpErr != nil {
			}
			if errors.Is(context.Canceled, s.ctx.Err()) {
				break
			}
		}
	}()
	c, ok := udpConns[track.Kind()]
	if !ok {
		return
	}
	b := make([]byte, 1500)
	for {
		// Read
		n, _, err := track.Read(b)
		if err != nil && err != io.EOF {
			s.cancel()
			break
		}
		// Write
		if _, err = c.conn.Write(b[:n]); err != nil {
			if errors.Is(context.Canceled, s.ctx.Err()) {
				break
			}
		}
	}
}

func (s *WebsocketServer) startFFmpeg(streamURL string) error {
	// Create a ffmpeg process that consumes MKV via stdin, and broadcasts out to Stream URL
	cmd := exec.CommandContext(s.ctx,
		"ffmpeg",
		"-protocol_whitelist", "file,udp,rtp",
		"-i", "/tmp/rtp-forwarder.sdp",
		"-c:v", "copy",
		"-c:a", "aac",
		"-f", "flv",
		"-strict", "-2",
		streamURL) //nolint
	_, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	ffmpegOut, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		scanner := bufio.NewScanner(ffmpegOut)
		for scanner.Scan() {
			if errors.Is(context.Canceled, s.ctx.Err()) {
				break
			}
		}
	}()
	return nil
}
