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
	"os"
	"os/exec"
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
		fmt.Printf("handle track %v, %v\n", track.Codec().MimeType, track.Kind().String())
		go func() {
			codec := track.Codec().MimeType
			if codec == webrtc.MimeTypeVP8 {
				pushVideoToFFmpeg(track)
			} else if codec == webrtc.MimeTypeOpus {
				pushAudioToFFmpeg(track)
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

// handle track
func passBackStream(track *webrtc.TrackRemote, audioTrack *webrtc.TrackLocalStaticRTP, videoTrack *webrtc.TrackLocalStaticRTP) {
	buf := make([]byte, 1500)
	for {
		i, _, readErr := track.Read(buf)
		if readErr != nil {
			break
		}
		if track.Kind() == webrtc.RTPCodecTypeAudio {
			_, _ = audioTrack.Write(buf[:i])
		} else if track.Kind() == webrtc.RTPCodecTypeVideo {
			_, _ = videoTrack.Write(buf[:i])
		}
	}
}

func pushRtmp(track *webrtc.TrackRemote) {
	if track.Kind() == webrtc.RTPCodecTypeAudio {
		return
	}
	ffmpegCmd := exec.Command("ffmpeg",
		"-i", "pipe:0", // 使用 pipe 输入
		"-c:v", "libx264", // 编码格式为 h264
		"-f", "flv", // 输出为 FLV 格式
		"rtmp://localhost:1935/stream/canoe") // 推送到 RTMP 地址
	stdin, err := ffmpegCmd.StdinPipe()
	if err != nil {
		fmt.Println("error: ", err)
	}
	go func() {
		err := ffmpegCmd.Start()
		if err != nil {
			panic(err)
		}
		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				panic(err)
			}
			// Write RTP packet to FFmpeg
			size, err := stdin.Write(packet.Payload)
			if err != nil {
				return
			}
			fmt.Println(size)
		}
	}()
}

func pushVideoToFFmpeg(track *webrtc.TrackRemote) {
	// Start FFmpeg process to push video stream to RTMP server
	ffmpegArgs := []string{
		"-i", "pipe:0", // Input from pipe (video)
		"-c:v", "libx264", // Re-encode using x264
		"-f", "flv", // Output format FLV
		"rtmp://localhost:1935/stream/canoe", // RTMP server URL
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		err := cmd.Start()
		if err != nil {
			fmt.Println("start ffmpeg error", err)
		}

		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				fmt.Printf("read rpt error: %v\n", err)
				break
			}
			// Write RTP packet to FFmpeg
			_, err = stdin.Write(packet.Payload)
			if err != nil {
				fmt.Printf("write rpt error: %v\n", err)
				return
			}
		}
	}()
}

func pushAudioToFFmpeg(track *webrtc.TrackRemote) {
	// Similar process for audio stream (Opus to AAC or MP3 conversion)
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0", // Input from stdin
		"-c:a", "aac", // Audio codec (convert Opus to AAC)
		"-f", "flv", // Output format RTMP/FLV
		"rtmp://localhost:1935/stream/canoe", // RTMP URL
	)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("error: ", err)
	}
	go func() {
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				fmt.Printf("Read RTP error: %v\n", err)
				break
			}
			_, err = stdin.Write(packet.Payload)
			if err != nil {
				return
			}
		}
	}()
}
