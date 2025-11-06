package signaling

import (
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
	"github.com/pion/webrtc/v3"
)

type ExistingChannels struct {
	channels map[string]*Channel
	mu       sync.Mutex
}
type Client struct {
	user_id    string
	ws_conn    *websocket.Conn
	is_ws_conn bool
	PCconn     *webrtc.PeerConnection
	mu         sync.Mutex
}
type trackForwarder struct {
	remoteTrack *webrtc.TrackRemote
	localTracks []*webrtc.TrackLocalStaticRTP
	mu          sync.RWMutex
}
type Channel struct {
	channel_id            string
	users                 map[string]*Client
	remoteTrackForwarders map[string]*trackForwarder
	mu                    sync.Mutex
}
type SignalingMsg struct {
	MsgType string          `json:"type"`
	UserID  *string         `json:"userid,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
type ConnectedPayload struct {
	User_id string `json:"user_id"`
}
