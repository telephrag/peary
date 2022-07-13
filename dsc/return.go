package dsc

import (
	"discordgo"
	"kubinka/db"
	"log"

	"github.com/pkg/errors"
)

func Return(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var err error
	defer logCommand(s, i, err)

	err = db.Instance.DeletePlayer(i.Member.User)
	if err != nil {
		log.Panic(errors.Errorf("Failed to delete player from db: %w", err))
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have returned from deployment.",
		},
	})
	if err != nil {
		log.Panic(errors.Errorf("Failed to respond to the player: %w", err))
	}
}
