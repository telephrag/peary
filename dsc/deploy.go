package dsc

import (
	"discordgo"
	"fmt"
	"kubinka/config"
	"kubinka/db"
	"time"

	"log"

	"github.com/pkg/errors"
)

func getDeployDuration(i *discordgo.InteractionCreate) time.Time {
	opt := i.ApplicationCommandData().Options
	m := opt[0].IntValue()
	if len(opt) > 1 {
		m += opt[1].IntValue() * 60
	}

	return time.Now().Add(time.Minute * time.Duration(m)).UTC()
}

func Deploy(s *discordgo.Session, i *discordgo.InteractionCreate) { // TODO: Reduce err handling boilerplate

	d := getDeployDuration(i)
	var err error = nil
	defer logCommand(s, i, err)

	err = db.Instance.InsertPlayer(i.Member.User, d)
	if err != nil {
		log.Println(errors.Errorf("Failed to add record to db: %w", err))
		return
	}

	err = s.GuildMemberRoleAdd(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		log.Println(errors.Errorf("Failed to issue role: %w", err))

		err = db.Instance.DeletePlayer(i.Member.User)
		if err != nil {
			log.Panic(errors.Errorf("Failed to recover: %w", err))
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You have been deployed till %02d:%02d", d.Hour(), d.Minute()),
		},
	})
	if err != nil {
		log.Panic(errors.Errorf("Failed to respond to the player: %w", err))
	}
}
