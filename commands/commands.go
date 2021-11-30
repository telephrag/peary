package commands

import (
	"discordgo"
)

var Commands = []discordgo.ApplicationCommand{
	{
		Name:        "deploy",
		Description: "Выдаёт на время специальную роль по которой вас смогут пинговать в #поиск-игроков.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "hours",
				Description: "Время в часах, 0..12",
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "minutes",
				Description: "Время в минутах, 0..59",
				Required:    false,
			},
		},
	},

	{
		Name:        "return",
		Description: "Забирает у вас роль, данную командой /deploy.",
	},
}
