package service

import (
	"bufio"
	. "canoe/internal/model"
	"canoe/internal/remote"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	grl "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"io"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

type WebsocketService struct {
	peers *sync.Map
}

type wsCtx struct {
	ctx     context.Context
	profile AlpsUserProfile
	cancel  context.CancelFunc
	conn    *neffos.NSConn
	msg     neffos.Message
}

type MsgHandler func(ctx *wsCtx) error

func ChatMsgHandler(ctx *wsCtx) error {
	var evp Envelope
	message := ctx.msg
	conn := ctx.conn
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

// SignalingMessage 信令消息
type SignalingMessage struct {
	Type      string                  `json:"type"`
	SDP       string                  `json:"sdp,omitempty"`
	Candidate webrtc.ICECandidateInit `json:"candidate,omitempty"`
	Intent    string                  `json:"intent,omitempty"`
}

func (s *WebsocketService) DialMsgHandler(ctx *wsCtx) error {
	var msg SignalingMessage
	message := ctx.msg
	conn := ctx.conn
	profile := ctx.profile
	peers := s.peers
	err := message.Unmarshal(&msg)
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
		intent := msg.Intent
		go func() {
			if intent == "_anyone_" {
				processLive(ctx.ctx, ctx.cancel, profile.Username, track, pc)
			} else if intent != "" {
				s.processDialog(intent)
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
		if value, ok := peers.Load(profile.Username); ok {
			candidate := webrtc.ICECandidateInit{Candidate: msg.Candidate.Candidate}
			pc := value.(*webrtc.PeerConnection)
			if err := pc.AddICECandidate(candidate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *WebsocketService) NewWsServer(token string, handler MsgHandler) *neffos.Server {
	ctx, cancel := context.WithCancel(context.Background())
	profile := remote.GetUserProfile(token)
	upgrader := gorilla.Upgrader(grl.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	ws := websocket.New(upgrader, websocket.Events{websocket.OnNativeMessage: func(conn *neffos.NSConn, message neffos.Message) error {
		wx := &wsCtx{
			ctx:    ctx,
			cancel: cancel,
			msg:    message,
			conn:   conn,
		}
		return handler(wx)
	}})

	// 当连接建立
	// 初始化用户会话信息
	ws.OnConnect = func(conn *neffos.Conn) error {
		logger.Infof("got new connection: access-token = %s", token)
		return nil
	}

	// 清理会话信息
	ws.OnDisconnect = func(c *neffos.Conn) {
		logger.Infof("disconnected: access-token = %s", token)
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("recover from error: %v", r)
			}
		}()
		s.peers.Delete(profile.Username)
		cancel()
	}
	return ws
}

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

func (s *WebsocketService) processDialog(username string) {
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

func NewWebSocketService() *WebsocketService {
	service := &WebsocketService{
		peers: &sync.Map{},
	}
	return service
}

type udpConn struct {
	conn *net.UDPConn
	port int
}

var udpConns = map[webrtc.RTPCodecType]*udpConn{
	webrtc.RTPCodecTypeAudio: {port: 4000},
	webrtc.RTPCodecTypeVideo: {port: 4002},
}

func initUDP() {
	// Create a local addr
	var laddr *net.UDPAddr
	var err error
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		logger.Errorf("resolve udp addr err: %v", err)
		return
	}
	for _, c := range udpConns {
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			logger.Errorf("handle resolve udp addr error: %v", err)
			return
		}
		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			logger.Errorf("handle dial udp error: %v", err)
			return
		}
	}
}

func processLive(ctx context.Context, cancel context.CancelFunc, username string, track *webrtc.TrackRemote, pc *webrtc.PeerConnection) {

	streamURL := fmt.Sprintf("%s/%s", "rtmp://localhost:1935/stream", username)
	err := startFFmpeg(ctx, streamURL)
	if err != nil {
		logger.Errorf("start ffmpeg error: %v", err)
		cancel()
		return
	}
	// Retrieve udp connection
	c, ok := udpConns[track.Kind()]
	if !ok {
		return
	}

	// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			if rtcpErr := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); rtcpErr != nil {
				logger.Errorf("handle write interval rtp error: %v", rtcpErr)
			}
			if errors.Is(context.Canceled, ctx.Err()) {
				break
			}
		}
	}()

	b := make([]byte, 1500)
	for {
		// Read
		n, _, err := track.Read(b)
		if err != nil && err != io.EOF {
			logger.Errorf("handle read rtp error: %v", err)
			cancel()
			break
		}
		// Write
		if _, err = c.conn.Write(b[:n]); err != nil {
			logger.Errorf("handle write rtp error: %v", err)
			if errors.Is(context.Canceled, ctx.Err()) {
				break
			}
		}
	}
}

func startFFmpeg(ctx context.Context, streamURL string) error {
	// Create a ffmpeg process that consumes MKV via stdin, and broadcasts out to Stream URL
	cmd := exec.CommandContext(ctx,
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
			if errors.Is(context.Canceled, ctx.Err()) {
				break
			}
		}
	}()
	return nil
}
