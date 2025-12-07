package signaling

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

func handleClientOffer(offer webrtc.SessionDescription, channel_id string, client *Client) {

	serverPeerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})

	if err != nil {
		fmt.Println("error while creating new peer connection")
		return
	}
	clientChannel, err := existingChannels.getChannelById(channel_id)
	if err != nil {
		fmt.Println("error during finding channel")
		return
	}
	err = serverPeerConnection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Println("error while setting remote description")
		return
	}
	client.PCconn = serverPeerConnection
	serverPeerConnection.OnNegotiationNeeded(func() {
		offer, err := serverPeerConnection.CreateOffer(&webrtc.OfferOptions{})
		if err != nil {
			fmt.Println("error during offer creation")
			return
		}
		err = serverPeerConnection.SetLocalDescription(offer)
		if err != nil {
			fmt.Println("error during setting local description")
			return
		}
		signalMsg, err := MarshalSignalingMsg("conn_offer", nil, offer)
		if err != nil {
			fmt.Println("error during marshaling offer")
			return
		}
		sendSignalMsg(signalMsg, client)
		fmt.Println("negotiation nedded")
	})
	serverPeerConnection.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Println("track arrives from the client start to creating offers")
		clientChannel.mu.Lock()
		for _, forwarder := range clientChannel.remoteTrackForwarders {
			forwarder.AddSubscriber(serverPeerConnection)
		}

		forwarder := NewTrackForwarder(remote)
		clientChannel.remoteTrackForwarders[remote.ID()] = forwarder

		client.mu.Lock()
		client.RemoteTrackIds[remote.ID()] = remote.ID()
		client.mu.Unlock()

		for _, user := range clientChannel.users {
			user.mu.Lock()
			if user.user_id != client.user_id {
				forwarder.AddSubscriber(user.PCconn)
			}
			user.mu.Unlock()
		}
		clientChannel.mu.Unlock()

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
