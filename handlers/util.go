package handlers

import (
	"discordgo"
	"log"
)

func logCommand(i *discordgo.InteractionCreate, err error) {
	log.Println(
		i.ApplicationCommandData().Name,
		i.Member.User.ID,
		i.Member.User.Username,
		err,
	)
}
