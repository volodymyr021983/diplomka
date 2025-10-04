package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test/discord/api/channels"
	"test/discord/db"
	"test/discord/db/models"

	"github.com/supertokens/supertokens-golang/recipe/session"
)

func CreateServer(dbContainer *db.DbContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		server_id, err := GetNewServerId(dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during getting new server id")
			return
		}

		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		var server_name ServerNameBody

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during reading body r")
			return
		}

		err = json.Unmarshal(body, &server_name)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during unmarshaling body r")
			return
		}

		userID := sessionContainer.GetUserID()
		new_server := models.Servers{ServerId: *server_id, Servername: server_name.Servername, OwnerId: userID}

		err = CreateNewServer(&new_server, dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during creation new server ")
			return
		}

		server_member := models.ServerMembers{
			UserID:   userID,
			ServerID: *server_id,
			UserRole: "owner",
		}

		err = AddNewUser(server_member, dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during adding user in server creation")
			return
		}

		channel_id, err := channels.GetNewChannelId(dbContainer)

		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during getting new channel id")
			return
		}
		channel := models.Channels{
			OwnServerId: *server_id,
			ChannelId:   *channel_id,
			Channelname: "basic",
			ChannelType: "text",
		}

		err = channels.CreateTextChannel(channel, dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during creation of channel")
			return
		}

		w.WriteHeader(200)
		jsonResp, err := json.Marshal(APIServerResponse{OK: true, Err: nil})
		if err == nil {
			w.Write(jsonResp)
		}

	}
}

func GetServers(dbContainer *db.DbContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()
		servers, err := GetUserServers(userID, dbContainer)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(200)
		jsonResp, err := json.Marshal(servers)
		if err == nil {
			w.Write(jsonResp)
		}

	}
}

func CreateInviteLink(dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server_id := r.PathValue("server_id")
		isTokenExists, token := IsServerInviteCodeExists(server_id, dbContainer)
		if isTokenExists {
			err := DeleteInviteCode(token, dbContainer)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println("fail during invite code deletion")
				return
			}
		}

		invite_code, create_at, exp_time, err := CreateInviteCode(server_id)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("fail during code creation")
			return
		}
		token = models.InvitationCodes{
			ServerID:  server_id,
			Token:     invite_code,
			CreatedAt: *create_at,
			ExpiresAt: *exp_time,
		}
		err = SaveInviteToken(token, dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during token saving")
			return
		}
		type InviteCodeResponse struct {
			Invite_code string
		}
		bodyResponse := InviteCodeResponse{
			Invite_code: token.Token,
		}
		byteBody, err := json.Marshal(bodyResponse)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during json encoding")
			return
		}

		w.WriteHeader(200)
		w.Write(byteBody)
	})
}

func AcceptInvitation(dbContainer *db.DbContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()

		invite_code := r.PathValue("invite_code")

		token, err := VerifyInviteCode(invite_code, dbContainer)
		if err != nil {
			w.WriteHeader(404)
			fmt.Println("token error")
			return
		}

		server_id, err := token.Claims.GetSubject()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if IsMember(userID, server_id, dbContainer) {
			type Response struct {
				Server_id  string
				Channel_id *string
			}
			channel_id, err := channels.GetFirstChannel(server_id, dbContainer)
			if err != nil {
				channel_id = nil
			}
			responseBody := Response{
				Server_id:  server_id,
				Channel_id: channel_id,
			}
			byteBody, err := json.Marshal(responseBody)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println("Error during encoding")
				return
			}
			w.WriteHeader(302)
			w.Write(byteBody)
			return
		}

		server_member := models.ServerMembers{
			UserID:   userID,
			ServerID: server_id,
			UserRole: "user",
		}

		err = AddNewUser(server_member, dbContainer)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error during user adding")
			return
		}
		w.WriteHeader(200)
	}
}
