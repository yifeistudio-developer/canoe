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
	"io"
	"log"
	"net/http"
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

		// 消息转播
		//go passBackStream(track, audioTrack, videoTrack)
		go pushRtmp(track)
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
	// Start FFmpeg process
	//cmd := exec.Command("ffmpeg", "-i", "-", "-f", "flv", "rtmp://localhost:1935/stream/canoe")
	// Create a pipe for passing data between Go and ffmpeg
	r, w := io.Pipe()

	// Run ffmpeg to push the stream to RTMP
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-f", "h264", // Indicate H.264 input format
		"-i", "pipe:0", // Input from stdin (pipe)
		"-c:v", "libx264", // Use x264 codec
		"-f", "flv", // Output format
		"rtmp://localhost:1935/stream/canoe", // RTMP URL
	)

	// Capture stderr to diagnose any issues
	ffmpegCmd.Stderr = log.Writer() // Redirect ffmpeg stderr to Go's logger

	ffmpegCmd.Stdin = r // Pipe RTP stream to ffmpeg

	// Start ffmpeg in a goroutine
	go func() {
		if err := ffmpegCmd.Run(); err != nil {
			log.Printf("Failed to run ffmpeg: %v", err)
		}
	}()

	// Buffer to store the incoming RTP packets
	packet := make([]byte, 1500)

	for {
		// Read RTP packet from the WebRTC track
		n, _, readErr := track.Read(packet)
		if readErr != nil {
			log.Printf("Error reading from track: %v", readErr)
			break
		}
		// Write the raw RTP packet (H.264) to the ffmpeg pipe
		if _, writeErr := w.Write(packet[:n]); writeErr != nil {
			log.Printf("Error writing to ffmpeg pipe: %v", writeErr)
			break
		}
	}
	// Close the pipe when done
	err := w.Close()
	if err != nil {
		return
	}
}
