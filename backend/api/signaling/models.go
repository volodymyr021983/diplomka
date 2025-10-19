package signaling

import (
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
)

type ExistingChannels struct {
	channels map[string]*Channel
	mu       sync.Mutex
}
type Client struct {
	user_id    string
	ws_conn    *websocket.Conn
	is_ws_conn bool
}
type Channel struct {
	channel_id string
	users      map[string]*Client
	mu         sync.Mutex
}
type SignalingMsg struct {
	MsgType string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
type ConnectedPayload struct {
	User_id string `json:"user_id"`
}
