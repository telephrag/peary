package cmd_deploy

import (
	"context"
	"discordgo"
	"fmt"
	"kubinka/command"
	"kubinka/config"
	"log"
	"time"

	"github.com/pkg/errors"
)

type DeployCmd struct {
	command.Command
	err error
	ctx context.Context
}

func Init() command.Command {
	return &DeployCmd{}
}

func getDeployDuration(i *discordgo.InteractionCreate) time.Time {
	opt := i.ApplicationCommandData().Options
	m := opt[0].IntValue()
	if len(opt) > 1 {
		m += opt[1].IntValue() * 60
	}

	return time.Now().Add(time.Minute * time.Duration(m)).UTC()
}

func (cmd *DeployCmd) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	d := getDeployDuration(i)

	// err = db.Instance.InsertPlayer(i.Member.User, d)
	// if err != nil {
	// 	log.Println(errors.Errorf("Failed to add record to db: %w", err))
	// 	return
	// }

	err := s.GuildMemberRoleAdd(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		log.Println(errors.Errorf("Failed to issue role: %w", err))
		cmd.err = err
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You have been deployed till %02d:%02d", d.Hour(), d.Minute()),
		},
	})
	if err != nil {
		msg := err.(discordgo.RESTError).Message
		log.Println(errors.Errorf("Failed to respond to the player: %v", msg))
		cmd.err = err
		return
	}
}

func (cmd *DeployCmd) Recover(ctx context.Context) error {
	return nil
}

func (cmd *DeployCmd) Log(ctx context.Context) error {
	return nil
}

func (cmd *DeployCmd) GetErr() error {
	return cmd.err
}
