package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"io"
	"net"
	"os/exec"
	"time"
)

type udpConn struct {
	conn *net.UDPConn
	port int
}

func processLive(username string, track *webrtc.TrackRemote, pc *webrtc.PeerConnection) {
	streamURL := fmt.Sprintf("%s/%s", "rtmp://localhost:1935/stream", username)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a local addr
	var laddr *net.UDPAddr
	var err error
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		cancel()
		return
	}
	udpConns := map[webrtc.RTPCodecType]*udpConn{
		webrtc.RTPCodecTypeAudio: {port: 4000},
		webrtc.RTPCodecTypeVideo: {port: 4002},
	}
	for _, c := range udpConns {
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			fmt.Println(err)
			cancel()
			return
		}
		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			fmt.Println(err)
			cancel()
			return
		}
	}

	err = startFFmpeg(ctx, streamURL)
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
		n, _, err := track.Read(b)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			cancel()
			break
		}
		// Write
		if _, err = c.conn.Write(b[:n]); err != nil {
			fmt.Println(err)
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
