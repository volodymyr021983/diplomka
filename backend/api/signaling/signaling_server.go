package signaling

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

var existingChannels = ExistingChannels{
	channels: make(map[string]*Channel),
}

func AcceptConnection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client_conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			fmt.Println("error while upgrading")
			return
		}
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		user_id := sessionContainer.GetUserID()
		channel_id := r.PathValue("channel_id")

		client := Client{
			user_id:    user_id,
			ws_conn:    client_conn,
			is_ws_conn: false,
		}
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		marshaled, _ := MarshalSignalingMsg("connected", ConnectedPayload{
			User_id: client.user_id,
		})
		client.ws_conn.Write(ctx, 1, *marshaled)

		for {

			derCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			_, msgByte, err := client.ws_conn.Read(derCtx)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			signalMsg, err := UnmarshalSignalMsg(msgByte)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			switch signalMsg.MsgType {
			case "connected":
				fmt.Println("web socket connection established")
				client.is_ws_conn = true
			case "join_channel":
				JoinRoom(channel_id, client)
			case "get_clients":
				GetClients(channel_id)
			}
		}
	})
}
