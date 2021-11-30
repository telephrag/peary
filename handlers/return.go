package handlers

import (
	"discordgo"
	"kubinka/config"
)

func Return(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.ApplicationCommandData().Name != "return" {
		return
	}

	err := s.GuildMemberRoleRemove(
		config.GuildID,
		config.Skif,
		config.RoleID,
	)
	if err != nil {
		// handle
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have returned from deployment :)",
		},
	})
}
