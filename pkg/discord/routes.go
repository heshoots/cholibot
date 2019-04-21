package discord

import (
	"github.com/quorauk/dmux"
)

func Handlers() []*dmux.DiscordRegexMessageHandler {
	var handlers = []*dmux.DiscordRegexMessageHandler{
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `!showroles$`,
			HandlerFn:          showRolesHandler,
			HandlerName:        "!showroles",
			HandlerDescription: "Show available roles",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!addrole (?P<role>.*)$`,
			HandlerFn:          AddRole,
			HandlerName:        "!addrole",
			RequiresAdmin:      true,
			HandlerDescription: "Add role to available roles",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!removerole (?P<role>.*)$`,
			HandlerFn:          RemoveRole,
			HandlerName:        "!removerole",
			RequiresAdmin:      true,
			HandlerDescription: "Remove role to available roles",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!iam (?P<role>.*)$`,
			HandlerFn:          IamHandler,
			HandlerName:        "!iam",
			HandlerDescription: "Give yourself a role",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!iamn (?P<role>.*)$`,
			HandlerFn:          IamnHandler,
			HandlerName:        "!iamn",
			HandlerDescription: "Remove a role from yourself",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!challonge (?P<name>\w+) (?P<game>.*)$`,
			HandlerFn:          ChallongeHandler,
			HandlerName:        "!challonge",
			HandlerDescription: "create challonge tournament",
			RequiresAdmin:      true,
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!setcommand (?P<command>\w+) (?P<response>(.*))$`,
			HandlerFn:          SetCommandHandler,
			HandlerName:        "!setcommand",
			RequiresAdmin:      true,
			HandlerDescription: "Set a custom response to a command",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern:     `^!(listcommands|help)$`,
			HandlerFn:          ListCommands,
			HandlerName:        "!help",
			HandlerDescription: "Get help",
		},
		&dmux.DiscordRegexMessageHandler{
			HandlerPattern: `^!(?P<command>\w+)$`,
			HandlerFn:      CustomCommand,
			HandlerName:    "server custom commands",
		},
	}
	return handlers
}
