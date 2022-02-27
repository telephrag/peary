package handlers

import (
	"discordgo"
	"kubinka/config"
	"kubinka/db"
	"log"
)

func Return(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.GuildMemberRoleRemove(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	defer logCommand(i, err)
	if err != nil {
		log.Println(err, " ", i.ApplicationCommandData())
		return
	}

	db.Instance.DeletePlayer(i.Member.User)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have returned from deployment.",
		},
	})

}
