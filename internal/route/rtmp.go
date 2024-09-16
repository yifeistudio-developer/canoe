package route

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"net"
	"os/exec"
	"time"
)

type udpConn struct {
	conn *net.UDPConn
	port int
}

func t(pc *webrtc.PeerConnection) {
	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a local addr
	var laddr *net.UDPAddr
	var err error
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		fmt.Println(err)
		cancel()
	}

	// Prepare udp conns
	udpConns := map[string]*udpConn{
		"audio": {port: 4000},
		"video": {port: 4002},
	}
	for _, c := range udpConns {
		// Create remote addr
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			fmt.Println(err)
			cancel()
		}

		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			fmt.Println(err)
			cancel()
		}
		//defer func(conn net.PacketConn) {
		//	if closeErr := conn.Close(); closeErr != nil {
		//		fmt.Println(closeErr)
		//	}
		//}(c.conn)
	}
	streamURL := fmt.Sprintf("%s/%s", "rtmp://localhost:1935/stream", "canoe")
	startFFmpeg(ctx, streamURL)

	// Set a handler for when a new remote track starts, this handler will forward data to
	// our UDP listeners.
	// In your application this is where you would handle/process audio/video
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {

		// Retrieve udp connection
		c, ok := udpConns[track.Kind().String()]
		if !ok {
			return
		}

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			for range ticker.C {
				if rtcpErr := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); rtcpErr != nil {
					fmt.Println(rtcpErr)
				}
				if errors.Is(context.Canceled, ctx.Err()) {
					break
				}
			}
		}()

		b := make([]byte, 1500)
		for {
			// Read
			n, _, readErr := track.Read(b)
			if readErr != nil {
				fmt.Println(readErr)
			}

			// Write
			if _, err = c.conn.Write(b[:n]); err != nil {
				fmt.Println(err)
				if errors.Is(context.Canceled, ctx.Err()) {
					break
				}
			}
		}
	})

}

func startFFmpeg(ctx context.Context, streamURL string) {
	// Create a ffmpeg process that consumes MKV via stdin, and broadcasts out to Stream URL
	ffmpeg := exec.CommandContext(ctx, "ffmpeg", "-protocol_whitelist", "file,udp,rtp", "-i", "/tmp/rtp-forwarder.sdp", "-c:v", "copy", "-c:a", "aac", "-f", "flv", "-strict", "-2", streamURL) //nolint
	_, err := ffmpeg.StdinPipe()
	if err != nil {
		return
	}
	ffmpegOut, _ := ffmpeg.StderrPipe()
	if err := ffmpeg.Start(); err != nil {
		panic(err)
	}
	go func() {
		scanner := bufio.NewScanner(ffmpegOut)
		for scanner.Scan() {
			if errors.Is(context.Canceled, ctx.Err()) {
				break
			}
		}
	}()
}
