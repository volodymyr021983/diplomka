package webrtcconfig

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func send(ctx context.Context, conn *websocket.Conn, msg WebSocketMessage) {
	// wsjson.Write is a helper that marshals the `msg` struct to JSON and sends it.
	// It takes a context, which allows the write operation to be cancelled.
	err := wsjson.Write(ctx, conn, msg)
	// If there's an error sending the message, we log it.
	if err != nil {
		log.Printf("error sending message: %v", err)
	}
}

// mustMarshal is a helper function that marshals an interface to a JSON byte slice.
// It's a convenience that avoids error handling for types we know will marshal correctly.
func mustMarshal(v interface{}) []byte {
	// json.Marshal converts the Go struct `v` into a JSON byte slice.
	b, _ := json.Marshal(v)
	// It returns the byte slice.
	return b
}

// findRoomForPeer is a utility to find which room a given peer belongs to.
func findRoomForPeer(user_id string) *Room {
	roomsLock.RLock()
	defer roomsLock.RUnlock()
	for _, room := range rooms {
		room.PeersLock.RLock()
		if _, ok := room.Peers[user_id]; ok {
			room.PeersLock.RUnlock()
			return room
		}
		room.PeersLock.RUnlock()
	}
	return nil
}
