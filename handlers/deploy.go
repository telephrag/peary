package handlers

import (
	"discordgo"
	"kubinka/config"
	"kubinka/models"
	"time"
)

func Deploy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// d := i.Interaction.Data

	if i.ApplicationCommandData().Name != "deploy" {
		return
	}

	p := models.Player{
		UserID:   config.Skif,
		UserName: i.Member.User.Username,
		Duration: 1,
		Begin:    time.Now(),
		End:      time.Now(),
	}

	err := s.GuildMemberRoleAdd(
		config.GuildID,
		p.UserID,
		config.RoleID,
	)
	if err != nil {
		//log.Panic(err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been deployed. No timeframe yet :(",
		},
	})
}
