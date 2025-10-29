package signaling

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

func handleClientOffer(offer webrtc.SessionDescription, peerKey string, client *Client) {

	serverPeerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		fmt.Println("error while creating new peer connection")
		return
	}
	err = serverPeerConnection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Println("error while setting remote description")
		return
	}
	err = client.addPeerConnection(serverPeerConnection, peerKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	serverPeerConnection.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Println("track arrives from the client")
	})
	serverAnswer, err := serverPeerConnection.CreateAnswer(&webrtc.AnswerOptions{})
	if err != nil {
		fmt.Println("error during answer creation")
		return
	}
	serverPeerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		signalMsg, err := MarshalSignalingMsg("new_ice_candidate", &client.user_id, candidate)
		if err != nil {
			fmt.Println("error during candidate marshaling")
		}
		sendSignalMsg(signalMsg, client)
	})

	serverPeerConnection.SetLocalDescription(serverAnswer)

	signalMsg, err := MarshalSignalingMsg("conn_answer", nil, serverAnswer)
	if err != nil {
		fmt.Println("error during marshaling answer")
		return
	}
	sendSignalMsg(signalMsg, client)

}
