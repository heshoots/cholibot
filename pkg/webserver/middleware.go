package webserver

import (
	"github.com/gorilla/mux"
	"github.com/heshoots/cholibot/pkg/discord"
	"log"
	"net/http"
)

func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", 307)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func allowedToModifyServer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api_key, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", 307)
			return
		}
		id := mux.Vars(r)["id"]
		guild, err := discord.GetGuild(id)

		if guild.OwnerID == getUser(api_key.Value).ID {
			next.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println("Guild not found, unauthenticated: " + id)
		w.WriteHeader(http.StatusUnauthorized)
	})
}
