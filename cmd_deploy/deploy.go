package cmd_deploy

import (
	"context"
	"discordgo"
	"fmt"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/config"
	"log"
	"time"
)

type DeployCmd struct {
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

func (cmd *DeployCmd) Handle(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	d := getDeployDuration(i)

	// err = db.Instance.InsertPlayer(i.Member.User, d)
	// if err != nil {
	// 	log.Println(errors.Errorf("Failed to add record to db: %w", err))
	// 	return
	// }

	select {
	case <-ctx.Done():
		bot_errors.NotifyUser(s, i, bot_errors.ErrSomewhereElse)
		log.Println(bot_errors.ErrSomewhereElse)
		return bot_errors.ErrSomewhereElse
	default:
	}

	err := s.GuildMemberRoleAdd(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		return bot_errors.ErrFailedGiveRole
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You have been deployed till %02d:%02d", d.Hour(), d.Minute()),
		},
	})
	if err != nil {
		return bot_errors.ErrFailedSendResponse
	}

	return err
}

func (cmd *DeployCmd) Recover(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.GuildMemberRoleRemove(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		return bot_errors.ErrFailedTakeRole
	}
	return nil
}

func (cmd *DeployCmd) GetErr() error {
	return cmd.err
}
