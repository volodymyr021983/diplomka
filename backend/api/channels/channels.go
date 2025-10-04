package channels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test/discord/db"
	"test/discord/db/models"

	"github.com/olahol/melody"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func WSConnectToChannel(m *melody.Melody, dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server_id := r.PathValue("server_id")
		channel_id := r.PathValue("channel_id")

		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		user_id := sessionContainer.GetUserID()
		username := getUserUsernameUsingId(user_id, dbContainer)

		if username == nil {
			w.WriteHeader(403)
			return
		}

		sessionKeys := map[string]any{
			"server_id":  server_id,
			"channel_id": channel_id,
			"username":   *username,
		}
		err := m.HandleRequestWithKeys(w, r, sessionKeys)
		if err != nil {
			return
		}

	})
}

func WSHandleMessage(m *melody.Melody) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		username, _ := s.Get("username")

		data := fmt.Sprintf("%s: \n%s", username, msg)
		m.BroadcastFilter([]byte(data), func(q *melody.Session) bool {
			server_id, _ := s.Get("server_id")
			channel_id, _ := s.Get("channel_id")
			user_serv, _ := q.Get("server_id")
			user_chan, _ := q.Get("channel_id")
			if server_id == user_serv && channel_id == user_chan {
				return true
			}
			return false
		})
	}
}

/*
	func WSConn(s *melody.Session) {
		server_id := s.Request.PathValue("server_id")
		channel_id := s.Request.PathValue("channel_id")

		s.Set("server_id", server_id)
		s.Set("channel_id", channel_id)
	}
*/
func CreateChannel(dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//sessionContainer := session.GetSessionFromRequestContext(r.Context())

		byteBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
		}

		type ChannelReq struct {
			ChannelName string
			ServerId    string
			ChannelType string
		}

		var channelBody ChannelReq
		err = json.Unmarshal(byteBody, &channelBody)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		channel_id, err := GetNewChannelId(dbContainer)
		server_id := r.PathValue("server_id")

		if err != nil {
			w.WriteHeader(409)
			return
		}

		createChannel := models.Channels{OwnServerId: server_id, ChannelId: *channel_id,
			Channelname: channelBody.ChannelName, ChannelType: channelBody.ChannelType}

		err = CreateTextChannel(createChannel, dbContainer)

		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})

}
func GetChannels(dbContainer *db.DbContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server_id := r.PathValue("server_id")
		channels := GetServerChannels(server_id, dbContainer)

		jsonResp, err := json.Marshal(channels)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(jsonResp)
	}
}

func GetFirstChannelHandler(dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server_id := r.PathValue("server_id")
		channel_id, err := GetFirstChannel(server_id, dbContainer)

		if err != nil {
			w.WriteHeader(404)
			return
		}

		bodyResponse, err := json.Marshal(map[string]string{"channel_id": *channel_id})
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(bodyResponse)
	})
}

/*
func WSHandleMessage(s *melody.Session, msg []byte) {
	m.Broadcast(msg)
	s.Write([]byte("TEEEEST"))
}
*/
