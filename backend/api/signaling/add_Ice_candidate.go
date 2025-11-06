package signaling

import (
	"github.com/pion/webrtc/v3"
)

func AddIceCandidate(iceCandidate *webrtc.ICECandidateInit, candidate_userid string, client *Client) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.PCconn.AddICECandidate(*iceCandidate)
}
