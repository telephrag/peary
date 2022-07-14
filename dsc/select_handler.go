package dsc

import (
	"context"
	"discordgo"
	"kubinka/command"
	"log"
)

func Master(s *discordgo.Session, i *discordgo.InteractionCreate) {
	init, ok := HandlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		log.Println("Couldn't retreive command Init(): ", i.ApplicationCommandData().Name)
		return
	}

	cmd, ok := init(context.TODO()).(command.Command)
	if !ok {
		log.Println("Couldn't init command state: ", i.ApplicationCommandData().Name)
		return
	}

	cmd.Handle(s, i)
}
