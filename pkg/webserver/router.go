package webserver

import (
	"github.com/gorilla/mux"
)

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", login)
	r.HandleFunc("/auth/discord/callback", callback)
	r.HandleFunc("/guild/add", addToGuild)
	return r
}
