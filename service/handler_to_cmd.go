package service

import (
	"peary/cmd_deploy"
	"peary/cmd_return"
	"peary/command"

	"github.com/bwmarrin/discordgo"
)

var CmdDef = []*discordgo.ApplicationCommand{
	{
		Name:        "play",
		Description: "Gives temporary role. Other players can ping you by it.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "hours",
				Description: "Time in hours.",
				MaxValue:    12,
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "minutes",
				Description: "Time in minutes.",
				MaxValue:    59,
				Required:    false,
			},
		},
	},

	{
		Name:        "return",
		Description: "Takes away role received after using /play.",
	},
}

var handlerToCmd = map[string]func(*command.Env) command.Command{
	"play":   cmd_deploy.Init,
	"return": cmd_return.Init,
}
