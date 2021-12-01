package handlers

import (
	"discordgo"
	"log"
)

func Select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h, ok := HandlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		log.Fatal("Couldn't retreive handler for command: ", i.ApplicationCommandData().Name)
		return
	}
	h(s, i)

	log.Println(
		i.ApplicationCommandData().Name,
		i.Member.User.ID,
		i.Member.User.Username,
	)
}
