package webserver

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/heshoots/cholibot/pkg/discord"
	"github.com/heshoots/cholibot/pkg/models"
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	api_key, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 307)

	}
	guilds := getUsersGuilds(api_key.Value)
	db := models.GetDB()
	var knownguilds []models.DiscordGuild
	db.Find(&knownguilds)

	for i, guild := range guilds {
		for _, knownguild := range knownguilds {
			if guild.ID == knownguild.GuildID {
				guilds[i].BotJoined = true
			}
		}
	}

	var homeData = HomeData{
		KnownGuilds: knownguilds,
		UserGuilds:  guilds,
	}

	tmpl := template.Must(template.ParseFiles("./web/home.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf8")

	tmpl.Execute(w, homeData)
}

func server(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	vars := mux.Vars(r)
	db := models.GetDB()
	var roles []models.DiscordRole
	type ServerData struct {
		ID         string
		Roles      []models.DiscordRole
		GuildRoles []*discordgo.Role
	}
	guildroles, err := discord.GetRolesForGuild(vars["id"])
	if err != nil {
		panic(err)
	}
	db.Find(&roles)
	for i, role := range roles {
		for _, grole := range guildroles {
			if grole.ID == role.RoleID {
				roles[i].Name = grole.Name
			}
		}
	}
	serverData := ServerData{
		ID:         vars["id"],
		Roles:      roles,
		GuildRoles: guildroles,
	}
	tmpl := template.Must(template.ParseFiles("./web/server.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	tmpl.Execute(w, serverData)
}

func roles(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	roleID := r.URL.Query()["id"][0]
	db := models.GetDB()
	var discordRole = models.DiscordRole{}
	db.FirstOrCreate(&discordRole, &models.DiscordRole{RoleID: roleID, GuildID: id})
	log.Println("Adding role: " + roleID)
	http.Redirect(w, r, "/server/"+id, 307)
}

func role(w http.ResponseWriter, r *http.Request) {
	db := models.GetDB()
	vars := mux.Vars(r)
	var model = &models.DiscordRole{RoleID: vars["roleid"]}
	db.First(&model)
	db.Delete(&model)
	http.Redirect(w, r, "/server/"+vars["id"], 307)
}
