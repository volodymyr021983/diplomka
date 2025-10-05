package api

import (
	"net/http"
	"test/discord/api/channels"
	"test/discord/api/servers"
	"test/discord/db"

	"test/discord/middleware"

	"github.com/olahol/melody"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func Api_init(m *melody.Melody, mux *http.ServeMux, dbContainer *db.DbContainer) {
	mux.HandleFunc("POST /api/server/create", session.VerifySession(nil, servers.CreateServer(dbContainer)))
	mux.HandleFunc("GET /api/server/getServers", session.VerifySession(nil, servers.GetServers(dbContainer)))
	mux.HandleFunc("GET /api/server/check-connect/{server_id}/{channel_id}", session.VerifySession(
		nil, middleware.ValidateConnectionToServer(
			middleware.ValidateConnectionToChannel(nil, dbContainer), dbContainer)))
	mux.HandleFunc("GET /api/server/connect/{server_id}/{channel_id}", session.VerifySession(
		nil, middleware.ValidateConnectionToServer(
			middleware.ValidateConnectionToChannel(
				channels.WSConnectToChannel(m, dbContainer), dbContainer), dbContainer)))
	mux.HandleFunc("POST /api/server/create-channel/{server_id}", session.VerifySession(nil,
		middleware.ValidateConnectionToServer(middleware.ValidateServerOwner(channels.CreateChannel(dbContainer),
			dbContainer),
			dbContainer)))
	mux.HandleFunc("GET /api/server/getChannels/{server_id}", session.VerifySession(nil, middleware.ValidateConnectionToServer(channels.GetChannels(dbContainer),
		dbContainer)))

	mux.HandleFunc("GET /api/server/invite/{server_id}", session.VerifySession(nil, middleware.ValidateConnectionToServer(middleware.ValidateServerOwner(servers.CreateInviteLink(dbContainer), dbContainer), dbContainer)))
	mux.HandleFunc("GET /api/server/invite/token/{invite_code}", session.VerifySession(nil, servers.AcceptInvitation(dbContainer)))
	mux.HandleFunc("GET /api/server/get-server-channel/{server_id}", session.VerifySession(nil, middleware.ValidateConnectionToServer(channels.GetFirstChannelHandler(dbContainer), dbContainer)))
	mux.HandleFunc("POST /api/server/delete-channel/{server_id}", session.VerifySession(nil, middleware.ValidateConnectionToServer(
		middleware.ValidateServerOwner(channels.DeleteChannel(dbContainer, m), dbContainer), dbContainer)))
	mux.HandleFunc("GET /api/server/delete-server/{server_id}", session.VerifySession(nil, middleware.ValidateConnectionToServer(
		middleware.ValidateServerOwner(servers.DeleteServer(m, dbContainer), dbContainer), dbContainer)))
	//melody package init if have time will be replaced with low level coder(nhooyr) WEBSOCKETS
	m.HandleMessage(channels.WSHandleMessage(m))

	//m.HandleConnect(channels.WSConn)
}
