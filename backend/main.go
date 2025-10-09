package main

import (
	//"fmt"
	"net/http"
	"test/discord/api"
	"test/discord/auth"
	"test/discord/db"
	"test/discord/middleware"

	//"github.com/supertokens/supertokens-golang/recipe/session"
	//"github.com/olahol/melody"
	"github.com/olahol/melody"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	mux := http.NewServeMux()
	//m := melody.New()
	dbContainer := db.ConnectToDb()

	auth.Supertokens_init(dbContainer)
	/*
		mux.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
			m.HandleRequest(w, req)
		})

		m.HandleMessage(func(s *melody.Session, msg []byte) {
			m.Broadcast(msg)
		})
	*/
	m := melody.New()

	api.Api_init(m, mux, dbContainer)
	http.ListenAndServeTLS(":443", "./certificate/localhost+2.pem", "./certificate/localhost+2-key.pem", middleware.CORS(supertokens.Middleware(mux)))

}
