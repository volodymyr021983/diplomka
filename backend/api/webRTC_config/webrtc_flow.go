package webrtcconfig

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
	"github.com/pion/webrtc/v3"
)

func joinRoom(roomID string, userID string, ws *websocket.Conn) {
	roomsLock.Lock()
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &Room{Channel_id: roomID, Peers: make(map[string]*Peer)}
	}
	roomsLock.Unlock()

	newPeer := &Peer{
		User_id:         userID,
		Conn:            ws,
		PublishedTracks: make(map[string]*webrtc.TrackLocalStaticRTP),
		Subscriptions:   make(map[string]*webrtc.PeerConnection),
	}

	room := rooms[roomID]
	room.PeersLock.Lock()
	room.Peers[userID] = newPeer
	room.PeersLock.Unlock()

	log.Printf("Peer %s joined Room %s", userID, roomID)
	send(context.Background(), ws, WebSocketMessage{
		Type:    "room_joined",
		Payload: json.RawMessage(string(mustMarshal(SignalPayload{MyPeerID: userID}))),
	})

	room.PeersLock.RLock()
	defer room.PeersLock.RUnlock()
	for existingPeerID, existingPeer := range room.Peers {
		if existingPeerID != userID {
			for _, track := range existingPeer.PublishedTracks {
				createSubscription(newPeer, existingPeer, track)
			}
		}
	}
}

// HANDLING OFFER FROM THE USER
func handleOffer(payload SignalPayload, userID string) {
	room := findRoomForPeer(userID)
	if room == nil {
		return
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Printf("handleOffer NewPeerConnection error: %v", err)
		return
	}

	// Set up the OnTrack handler BEFORE acquiring any locks.
	pc.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		localTrack, trackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), remoteTrack.StreamID())
		if trackErr != nil {
			log.Printf("OnTrack NewTrackLocalStaticRTP error: %v", trackErr)
			return
		}

		// Acquire lock only to update the shared PublishedTracks map.
		room.PeersLock.Lock()
		room.Peers[userID].PublishedTracks[localTrack.ID()] = localTrack
		room.PeersLock.Unlock()

		go func() { // Goroutine to forward RTP packets
			rtpBuf := make([]byte, 1500)
			for {
				i, _, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					return
				}
				if _, writeErr := localTrack.Write(rtpBuf[:i]); writeErr != nil {
					return
				}
			}
		}()

		// Acquire a read lock to safely iterate and create subscriptions for others.
		room.PeersLock.RLock()
		// We use a temporary slice to avoid holding the lock during the createSubscription call
		peersToSubscribe := []*Peer{}
		for existingPeerID, existingPeer := range room.Peers {
			if existingPeerID != userID {
				peersToSubscribe = append(peersToSubscribe, existingPeer)
			}
		}
		room.PeersLock.RUnlock()

		// Now create subscriptions without holding the lock
		for _, existingPeer := range peersToSubscribe {
			createSubscription(existingPeer, room.Peers[userID], localTrack)
		}
	})

	// Now that callbacks are set up, acquire lock to update the peer state.
	room.PeersLock.Lock()
	peer := room.Peers[userID]
	peer.PeerConnection = pc
	room.PeersLock.Unlock()

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		send(context.Background(), peer.Conn, WebSocketMessage{
			Type:    "candidate",
			Payload: json.RawMessage(string(mustMarshal(SignalPayload{Candidate: c.ToJSON(), RemotePeerID: peer.User_id}))),
		})
	})

	if err := pc.SetRemoteDescription(payload.SDP); err != nil {
		log.Printf("handleOffer SetRemoteDescription error: %v", err)
		return
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Printf("handleOffer CreateAnswer error: %v", err)
		return
	}

	if err := pc.SetLocalDescription(answer); err != nil {
		log.Printf("handleOffer SetLocalDescription error: %v", err)
		return
	}

	send(context.Background(), peer.Conn, WebSocketMessage{
		Type:    "answer",
		Payload: json.RawMessage(string(mustMarshal(SignalPayload{SDP: answer, RemotePeerID: peer.User_id}))),
	})
}
func handleAnswer(payload SignalPayload, peerID string) {
	room := findRoomForPeer(peerID)
	if room == nil {
		return
	}

	room.PeersLock.RLock()
	defer room.PeersLock.RUnlock()
	subscriber := room.Peers[peerID]

	if subPC, ok := subscriber.Subscriptions[payload.RemotePeerID]; ok {
		log.Printf("Applying answer for sub %s to stream from %s. Current state: %s", subscriber.User_id, payload.RemotePeerID, subPC.SignalingState().String())
		if err := subPC.SetRemoteDescription(payload.SDP); err != nil {
			log.Printf("error setting remote description for subscription: %v", err)
		}
	} else {
		log.Printf("Could not find subscription for peer %s for remote %s", peerID, payload.RemotePeerID)
	}
}

func handleCandidate(payload SignalPayload, userID string) {
	room := findRoomForPeer(userID)
	if room == nil {
		return
	}

	room.PeersLock.RLock()
	defer room.PeersLock.RUnlock()
	peer := room.Peers[userID]

	if payload.RemotePeerID == userID {
		if peer.PeerConnection != nil {
			peer.PeerConnection.AddICECandidate(payload.Candidate)
		}
	} else {
		if subPC, ok := peer.Subscriptions[payload.RemotePeerID]; ok {
			subPC.AddICECandidate(payload.Candidate)
		}
	}
}

func createSubscription(subscriber *Peer, publisher *Peer, track *webrtc.TrackLocalStaticRTP) {
	downlinkPC, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Printf("createSubscription NewPeerConnection error: %v", err)
		return
	}

	// Acquire lock to safely update the subscriber's map
	subscriberRoom := findRoomForPeer(subscriber.User_id)
	subscriberRoom.PeersLock.Lock()
	subscriber.Subscriptions[publisher.User_id] = downlinkPC
	subscriberRoom.PeersLock.Unlock()

	downlinkPC.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		send(context.Background(), subscriber.Conn, WebSocketMessage{
			Type: "candidate",
			Payload: json.RawMessage(string(mustMarshal(SignalPayload{
				Candidate:    c.ToJSON(),
				RemotePeerID: publisher.User_id,
			}))),
		})
	})

	if _, err = downlinkPC.AddTrack(track); err != nil {
		log.Printf("createSubscription AddTrack error: %v", err)
		return
	}

	offer, err := downlinkPC.CreateOffer(nil)
	if err != nil {
		log.Printf("createSubscription CreateOffer error: %v", err)
		return
	}

	// EDITED: Added crucial error handling. This is the most likely fix.
	if err := downlinkPC.SetLocalDescription(offer); err != nil {
		log.Printf("createSubscription SetLocalDescription error: %v", err)
		return
	}
	log.Printf("Downlink PC state for sub %s after SetLocalDescription: %s", subscriber.User_id, downlinkPC.SignalingState().String())

	send(context.Background(), subscriber.Conn, WebSocketMessage{
		Type: "offer",
		Payload: json.RawMessage(string(mustMarshal(SignalPayload{
			SDP:          offer,
			RemotePeerID: publisher.User_id,
		}))),
	})
}
