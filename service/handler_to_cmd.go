package service

import (
	"kubinka/cmd_deploy"
	"kubinka/cmd_return"
	"kubinka/command"

	"github.com/bwmarrin/discordgo"
)

var CmdDef = []*discordgo.ApplicationCommand{
	{
		Name:        "deploy",
		Description: "Gives temporary role by which you can be pinged by other people who want to play.",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "minutes",
				Description: "Time in minutes. Will be converted to hours automatically if > 60",
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "hours",
				Description: "Time in hours.",
				Required:    false,
			},
		},
	},

	{
		Name:        "return",
		Description: "Takes away role received after using /deploy.",
	},
}

var handlerToCmd = map[string]func(*command.Env) command.Command{
	"deploy": cmd_deploy.Init,
	"return": cmd_return.Init,
}
