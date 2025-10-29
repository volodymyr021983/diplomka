package signaling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coder/websocket"
	"github.com/pion/webrtc/v3"
)

func MarshalSignalingMsg(msgType string, user_id *string, payload interface{}) (*[]byte, error) {
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("error while marshaling payload")
	}
	result, err := json.Marshal(SignalingMsg{
		MsgType: msgType,
		UserID:  user_id,
		Payload: payloadMarshal,
	})

	if err != nil {
		return nil, errors.New("error while marshaling signal msg")
	}
	return &result, nil
}

func UnmarshalSignalMsg(msg []byte) (*SignalingMsg, error) {
	var signalMsg SignalingMsg
	err := json.Unmarshal(msg, &signalMsg)
	if err != nil {
		return nil, errors.New("error while unmarshaling signal msg")
	}
	return &signalMsg, nil
}
func sendSignalMsg(msg *[]byte, client *Client) {
	client.ws_conn.Write(context.Background(), websocket.MessageText, *msg)
}

func (channel *Channel) addUser(client *Client) error {
	channel.mu.Lock()
	defer channel.mu.Unlock()
	user := channel.users[client.user_id]

	if user != nil {
		return errors.New("user already in channel")
	}
	channel.users[client.user_id] = client
	return nil
}

// add peer connection to the map of peer connections
// This peer connections represent another user peer conns
// for example if user connect to the channel server create new peer connections from every existing user in this channel
// to stream their streams to the new user
func (client *Client) addPeerConnection(peerConnection *webrtc.PeerConnection, peerKey string) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	serverPeerConnection := client.RTCcons[peerKey]
	if peerKey == "server" && serverPeerConnection == nil {
		client.RTCcons = make(map[string]*webrtc.PeerConnection)
	}
	if serverPeerConnection != nil {
		//it is needs to be finalized. If user already have a conn from another user
		//for example if user already present in voice chan and try to connect to another voice chan
		return fmt.Errorf("this user: %s already have connection", peerKey)
	}
	client.RTCcons[peerKey] = peerConnection
	return nil
}
