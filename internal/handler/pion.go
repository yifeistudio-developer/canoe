package handler

import (
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"log"
)

func HandleWebRTConnection(ws *websocket.Conn) {
	api := webrtc.NewAPI()
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("Failed to create PeerConnection:", err)
		return
	}
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		candidate := c.ToJSON()
		err := ws.WriteJSON(candidate)
		if err != nil {
			log.Println("Failed to write ICE candidate:", err)
		}
	})
	for {
		var offer webrtc.SessionDescription
		err := ws.ReadJSON(&offer)
		if err != nil {
			log.Println("Failed to read offer:", err)
			break
		}

		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			log.Println("Failed to set remote description:", err)
			break
		}

		// Create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Println("Failed to create answer:", err)
			break
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Println("Failed to set local description:", err)
			break
		}

		err = ws.WriteJSON(answer)
		if err != nil {
			log.Println("Failed to write answer:", err)
			break
		}
	}
}
