package orchestrator

import (
	"discordgo"
	"kubinka/config"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "deploy",
		Description: "Выдаёт на время специальную роль по которой вас смогут пинговать в #" + config.ChanName + ".",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "minutes",
				Description: "Время в минутах, автоматически конвертируется в часы",
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "hours",
				Description: "Время в часах, 0..12",
				Required:    false,
			},
		},
	},

	{
		Name:        "return",
		Description: "Забирает роль, данную командой /deploy.",
	},
}
