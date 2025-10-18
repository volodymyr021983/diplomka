package webrtcconfig

import (
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	roomsLock sync.RWMutex
	rooms     = make(map[string]*Room)
)

type Room struct {
	Channel_id string
	Peers      map[string]*Peer
	PeersLock  sync.RWMutex
}
type Peer struct {
	User_id         string
	Conn            *websocket.Conn
	PeerConnection  *webrtc.PeerConnection
	PublishedTracks map[string]*webrtc.TrackLocalStaticRTP
	Subscriptions   map[string]*webrtc.PeerConnection
}
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type JoinRoomPayload struct {
	Channel_id string `json:"channel_id"`
}
type SignalPayload struct {
	SDP          webrtc.SessionDescription `json:"sdp,omitempty"`
	Candidate    webrtc.ICECandidateInit   `json:"candidate,omitempty"`
	RemotePeerID string                    `json:"remote_peer_id"` // ID of the OTHER peer in the exchange.
	MyPeerID     string                    `json:"my_peer_id,omitempty"`
}
