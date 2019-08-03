package discord

import (
	"github.com/quorauk/cholibot/pkg/models"
	"github.com/quorauk/dmux"
	"testing"
)

var discordBotUser dmux.Session
var testDB models.ModelSource

func TestGetRole(t *testing.T) {
	var mockRoles []dmux.Role
	var mockRole = dmux.MockRole{RoleID: "1234", RoleName: "Test Role"}
	mockRoles = append(mockRoles, mockRole)
	mockGuild := dmux.MockGuild{SessionRoles: mockRoles}
	mockSession := dmux.MockSession{Guild: mockGuild}
	role, err := GetRole(mockSession, mockGuild, "Test Role")
	if err != nil {
		t.Errorf("GetRole threw an error, " + err.Error())
	}
	if role.Name() != "Test Role" {
		t.Errorf("Incorrect Role returned")
	}
}

/*
func TestIam(t *testing.T) {
	configor.Load(&Config, "../../config/test.yaml")
	discordInstance, _ = dmux.Router(Config.DiscordBot.AuthToken)
	discordBotUser, _ = dmux.Router(Config.DiscordBot.TestAccount)
	discordInstance.AddHandler(
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!iam (?P<role>.*)$`,
			HandlerFn:          IamHandler,
			HandlerName:        "!iam",
			HandlerDescription: "Give yourself a role",
		},
	)
	role, err := discordInstance.GuildRoleCreate("guild", "Test Role")
	if err != nil {
		t.Errorf("Couldn't create role for integration test, " + err.Error())
	}
	t.Logf(role.Name())
	testDB = models.GetClient()
	testDB.SetRole(Config.DiscordBot.TestGuild, role.Name(), role.ID())
	discordBotUser.MessageChannel(dmux.CreateDiscordChannel(role.ID()), dmux.DiscordMessageString("!iam TestRole"))
}*/
