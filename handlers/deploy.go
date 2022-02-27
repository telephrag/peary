package handlers

import (
	"discordgo"
	"fmt"
	"kubinka/config"
	"kubinka/db"
	"time"

	"log"
)

func getDeployDuration(i *discordgo.InteractionCreate) time.Duration {
	opt := i.ApplicationCommandData().Options
	m := opt[0].IntValue()
	if len(opt) > 1 {
		m += opt[1].IntValue() * 60
	}

	return time.Duration(m)
}

func Deploy(s *discordgo.Session, i *discordgo.InteractionCreate) { // TODO: Reduce err handling boilerplate

	d := getDeployDuration(i)
	respContent := fmt.Sprintf("You have been deployed till %v", d)
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: respContent,
		},
	})
	var err error = nil
	defer logCommand(i, err)

	err = db.Instance.InsertPlayer(i.Member.User, d)
	if err != nil {
		log.Println(err, " ", i.ApplicationCommandData())
		respContent = fmt.Sprint(err)
		return
	}

	err = s.GuildMemberRoleAdd(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)

	if err != nil {
		db.Instance.DeletePlayer(i.Member.User) // delete player from db cause we failed to give him role
		log.Println(err, " ", i.ApplicationCommandData())
		respContent = fmt.Sprint(err)
		return
	}

}
