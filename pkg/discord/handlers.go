package discord

import (
	"errors"
	"github.com/heshoots/cholibot/pkg/challonge"
	"github.com/heshoots/cholibot/pkg/models"
	"github.com/heshoots/dmux"
	log "github.com/sirupsen/logrus"
)

func GetRole(s dmux.Session, guild string, roleName string) (dmux.Role, error) {
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
		s.MessageChannel(ctx.ChannelID(), "Couldn't find challonge credentials for this server")
		return
	}
	url, err := challonge.CreateTournament(apikey, subdomain, tournamentName, game)
	log.Println(url)
	if err != nil {
		log.Error(err)
		s.MessageChannel(ctx.ChannelID(), "Couldn't create tournament")
		return
	}
	s.MessageChannel(ctx.ChannelID(), url)
}

func AddRole(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		discordRole, err := GetRole(s, ctx.GuildID(), role)
		if err != nil {
			log.Error(err)
			_, err = s.MessageChannel(ctx.ChannelID(), "couldn't find role in server")
			return
		}
		db := models.GetClient()
		err = db.SetRole(ctx.GuildID(), role, discordRole.ID())
		if err != nil {
			log.Error(err)
			_, err = s.MessageChannel(ctx.ChannelID(), "couldn't add new role")
			return
		}
		_, err = s.MessageChannel(ctx.ChannelID(), role+" added to roles")
	}
}

func RemoveRole(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		db := models.GetClient()
		err := db.RemoveRole(ctx.GuildID(), role)
		if err != nil {
			log.Error(err)
			_, err = s.MessageChannel(ctx.ChannelID(), "couldn't remove role")
			return
		}
		_, err = s.MessageChannel(ctx.ChannelID(), "role removed")
	}
}

func IamHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		db := models.GetClient()
		dbrole, _ := db.GetRole(ctx.GuildID(), role)
		err := s.GuildMemberRoleAdd(ctx.GuildID(), ctx.UserID(), dbrole)
		if err != nil {
			log.Error(err)
			s.MessageChannel(ctx.ChannelID(), "couldn't add role")
			return
		}
		s.MessageChannel(ctx.ChannelID(), "role added")
	}
}

func IamnHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		role := context.Groups()["role"]
		log.Info(role)
		db := models.GetClient()
		dbrole, _ := db.GetRole(ctx.GuildID(), role)
		err := s.GuildMemberRoleRemove(ctx.GuildID(), ctx.UserID(), dbrole)
		if err != nil {
			log.Error(err)
			s.MessageChannel(ctx.ChannelID(), "couldn't remove role")
			return
		}
		s.MessageChannel(ctx.ChannelID(), "role removed")
	}
}

func showRolesHandler(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		db := models.GetClient()
		roles, err := db.GetRoles(ctx.GuildID())
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
		_, err = s.MessageChannel(ctx.ChannelID(), out)
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
		err := db.SetCustomCommand(ctx.GuildID(), groups["command"], groups["response"])
		if err != nil {
			s.MessageChannel(ctx.ChannelID(), "Couldn't set command")
			log.Error(err)
			return
		}
		s.MessageChannel(ctx.ChannelID(), "Command set")
		return
	}
}

func CustomCommand(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		command := context.Groups()["command"]
		db := models.GetClient()
		resp, err := db.GetCustomCommand(ctx.GuildID(), command)
		if err != nil {
			return
		}
		s.MessageChannel(ctx.ChannelID(), resp)
	}
}

func ListCommands(s dmux.Session, context dmux.RegexHandlerContext) {
	ok, ctx := context.MessageContext()
	if ok {
		isAdmin, err := ctx.UserAdmin(s)
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
		customCommands, err := db.GetCustomCommands(ctx.GuildID())
		if err != nil {
			log.Error("Couldn't get custom commands for server")
			s.MessageChannel(ctx.ChannelID(), commands)
			return
		}
		for _, command := range customCommands {
			commands += "\n!" + command
		}
		s.MessageChannel(ctx.ChannelID(), commands)
		return
	}
}
