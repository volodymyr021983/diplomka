package webrtcconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func WebRTCConnectToVoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("here")
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Printf("websocket accept error: %v", err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "internal server error")
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		user_id := sessionContainer.GetUserID()
		log.Printf("New client connected: %s", user_id)

		for {
			var msg WebSocketMessage

			err := wsjson.Read(r.Context(), c, &msg)
			if err != nil {
				// Check if the error is a normal closure.
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					log.Printf("client %s disconnected normally", user_id)
				} else {
					log.Printf("read error for peer %s: %v", user_id, err)
				}
				break
			}
			switch msg.Type {
			case "join_room":
				var payload JoinRoomPayload
				json.Unmarshal(msg.Payload, &payload)
				joinRoom(payload.Channel_id, user_id, c) // Pass the nhooyr connection object

			case "offer": // This is an offer from a client wanting to PUBLISH their stream.
				var payload SignalPayload
				json.Unmarshal(msg.Payload, &payload)
				handleOffer(payload, user_id)
			case "candidate":
				var payload SignalPayload
				json.Unmarshal(msg.Payload, &payload)
				handleCandidate(payload, user_id)
			case "answer":
				var payload SignalPayload
				json.Unmarshal(msg.Payload, &payload)
				handleAnswer(payload, user_id)
			}
		}
	}

}
