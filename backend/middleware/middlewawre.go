package middleware

import (
	"net/http"
	"test/discord/api/channels"
	"test/discord/api/servers"
	"test/discord/db"

	"github.com/rs/cors"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func CORS(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	})

	return c.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
	}))
}

func ValidateConnectionToServer(next http.Handler, dbContainer *db.DbContainer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()

		server_id := r.PathValue("server_id")

		server := servers.FindServerById(server_id, dbContainer)
		if server == nil {
			w.WriteHeader(404)
			return
		}
		isItMember := servers.IsMember(userID, server_id, dbContainer)
		if !isItMember {
			w.WriteHeader(403)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ValidateConnectionToChannel(next http.Handler, dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		channel_id := r.PathValue("channel_id")
		server_id := r.PathValue("server_id")

		channel := channels.FindChannelById(channel_id, dbContainer)
		if channel == nil {
			w.WriteHeader(404)
			return
		}

		if channel.OwnServerId != server_id {
			w.WriteHeader(404)
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}

	})
}

func ValidateServerOwner(next http.Handler, dbContainer *db.DbContainer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server_id := r.PathValue("server_id")

		server := servers.FindServerById(server_id, dbContainer)

		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()

		if server == nil {
			w.WriteHeader(404)
			return
		}

		if server.OwnerId != userID {
			w.WriteHeader(403)
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}

	})
}
