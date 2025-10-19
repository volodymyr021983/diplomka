package signaling

import (
	"encoding/json"
	"errors"
)

func MarshalSignalingMsg(msgType string, payload interface{}) (*[]byte, error) {
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("error while marshaling payload")
	}
	result, err := json.Marshal(SignalingMsg{
		MsgType: msgType,
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

func (channel *Channel) addUser(client Client) error {
	channel.mu.Lock()
	defer channel.mu.Unlock()
	user := channel.users[client.user_id]

	if user != nil {
		return errors.New("user already in channel")
	}
	channel.users[client.user_id] = &client
	return nil
}
