package handlers

import (
	"discordgo"
	"kubinka/config"
	"log"
)

func Deploy(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.GuildMemberRoleAdd(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)

	if err != nil {
		log.Println(err, " ", i.ApplicationCommandData())
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been deployed. No timeframe yet :(",
		},
	})

	log.Println(
		i.ApplicationCommandData().Name,
		i.Member.User.ID,
		i.Member.User.Username,
	)
}
