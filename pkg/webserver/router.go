package webserver

import (
	"github.com/gorilla/mux"
)

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", login)
	s := r.PathPrefix("/server").Subrouter()
	s.HandleFunc("/{id}", server)
	s.HandleFunc("/{id}/roles", roles)
	s.HandleFunc("/{id}/roles/{roleid}", role)
	s.Use(allowedToModifyServer)
	r.HandleFunc("/auth/discord/callback", callback)
	r.HandleFunc("/", home)
	return r
}
