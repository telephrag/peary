package dsc

import (
	"discordgo"
	"log"
)

func Master(s *discordgo.Session, i *discordgo.InteractionCreate) {
	init, ok := HandlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		log.Println("Couldn't retreive command Init(): ", i.ApplicationCommandData().Name)
		return
	}

	cmd := init()

	cmd.Handle(s, i)
}
