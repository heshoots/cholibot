package discord

import (
	"errors"
	"github.com/heshoots/cholibot/pkg/challonge"
	"github.com/heshoots/cholibot/pkg/models"
	"github.com/heshoots/dmux"
	log "github.com/sirupsen/logrus"
)

func SendMessage(s dmux.Session, c dmux.Channel, m string) (dmux.Message, error) {
	return s.MessageChannel(c, dmux.DiscordMessageString(m))
}

func GetRole(s dmux.Session, guild dmux.Guild, roleName string) (dmux.Role, error) {
	roles, err := s.GuildRoles(guild)
	if err != nil {
		return nil, err
	}
	for _, guildrole := range roles {
		if guildrole.Name() == roleName {
			return guildrole, nil
		}
	}
	return nil, errors.New("Role not found")
}

func ChallongeHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	_, ctx := context.MessageContext()
	tournamentName := context.Groups()["name"]
	game := context.Groups()["game"]
	apikey, subdomain, err := challonge.GetCredentials(ctx)
	if err != nil {
		log.Error(err)
		SendMessage(s, ctx.Channel(), "Couldn't find challonge credentials for this server")
		return
	}
	url, err := challonge.CreateTournament(apikey, subdomain, tournamentName, game)
	log.Println(url)
	if err != nil {
		log.Error(err)
		SendMessage(s, ctx.Channel(), "Couldn't create tournament")
		return
	}
	SendMessage(s, ctx.Channel(), url)
}

func AddRole(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		discordRole, err := GetRole(s, ctx.Guild(), role)
		if err != nil {
			log.Error(err)
			SendMessage(s, ctx.Channel(), "couldn't find role in server")
			return
		}
		db := models.GetClient()
		err = db.SetRole(ctx.Guild().ID(), role, discordRole.ID())
		if err != nil {
			log.Error(err)
			SendMessage(s, ctx.Channel(), "couldn't add new role")
			return
		}
		SendMessage(s, ctx.Channel(), role+" added to roles")
	}
}

func RemoveRole(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		db := models.GetClient()
		err := db.RemoveRole(ctx.Guild().ID(), role)
		if err != nil {
			log.Error(err)
			SendMessage(s, ctx.Channel(), "couldn't remove role")
			return
		}
		SendMessage(s, ctx.Channel(), "role removed")
	}
}

func IamHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		db := models.GetClient()
		dbrole, _ := db.GetRole(ctx.Guild().ID(), role)
		err := s.GuildMemberRoleAdd(ctx.Guild(), ctx.User(), dmux.CreateDiscordRole(dbrole))
		if err != nil {
			log.Error(err)
			SendMessage(s, ctx.Channel(), "couldn't add role")
			return
		}
		SendMessage(s, ctx.Channel(), "role added")
	}
}

func IamnHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		log.Info(role)
		db := models.GetClient()
		dbrole, _ := db.GetRole(ctx.Guild().ID(), role)
		err := s.GuildMemberRoleRemove(ctx.Guild(), ctx.User(), dmux.CreateDiscordRole(dbrole))
		if err != nil {
			log.Error(err)
			SendMessage(s, ctx.Channel(), "couldn't remove role")
			return
		}
		SendMessage(s, ctx.Channel(), "role removed")
	}
}

func showRolesHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		db := models.GetClient()
		roles, err := db.GetRoles(ctx.Guild().ID())
		if err != nil {
			log.Error(err)
		}
		out := `
To get a role use !iam Role
To remove a role use !iamn Role

Roles ending in "Fighters" can be @ mentioned

Available Roles
-----------
`
		for _, role := range roles {
			out += "\n!iam " + role
		}
		_, err = SendMessage(s, ctx.Channel(), out)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func SetCommandHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		db := models.GetClient()
		groups := context.Groups()
		err := db.SetCustomCommand(ctx.Guild().ID(), groups["command"], groups["response"])
		if err != nil {
			SendMessage(s, ctx.Channel(), "Couldn't set command")
			log.Error(err)
			return
		}
		SendMessage(s, ctx.Channel(), "Command set")
		return
	}
}

func CustomCommand(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		command := context.Groups()["command"]
		db := models.GetClient()
		resp, err := db.GetCustomCommand(ctx.Guild().ID(), command)
		if err != nil {
			return
		}
		SendMessage(s, ctx.Channel(), resp)
	}
}

func ListCommands(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		isAdmin, err := ctx.User().Admin(s, ctx.Channel())
		if err != nil {
			log.Error("Couldn't determine if admin user")
			isAdmin = false
		}
		commands := "Commands\n---------------"
		for _, handler := range Handlers() {
			if !handler.NeedsAdmin() {
				commands += "\n" + handler.Name()
			}
			if handler.NeedsAdmin() && isAdmin {
				commands += "\n" + handler.Name()
			}
		}
		db := models.GetClient()
		customCommands, err := db.GetCustomCommands(ctx.Guild().ID())
		if err != nil {
			log.Error("Couldn't get custom commands for server")
			SendMessage(s, ctx.Channel(), commands)
			return
		}
		for _, command := range customCommands {
			commands += "\n!" + command
		}
		SendMessage(s, ctx.Channel(), commands)
		return
	}
}
