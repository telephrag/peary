package cmd_deploy

import (
	"fmt"
	"kubinka/bot_errors"

	"github.com/bwmarrin/discordgo"
)

type MsgResponseStep struct {
	DiscordSession    *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
}

func NewMsgResponseStep(s *discordgo.Session, i *discordgo.InteractionCreate) *MsgResponseStep {
	return &MsgResponseStep{
		DiscordSession:    s,
		InteractionCreate: i,
	}
}

func (s *MsgResponseStep) Do() error {
	// wont match with real expiration time but that's not critical
	d := getDeployDuration(s.InteractionCreate)

	err := s.DiscordSession.InteractionRespond(
		s.InteractionCreate.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You have been deployed till %02d:%02d", d.Hour(), d.Minute()),
			},
		})
	if err != nil {
		return bot_errors.New(
			s.InteractionCreate.Member.User.ID,
			bot_errors.CmdDeployDo,
			fmt.Errorf("%s: %w", bot_errors.ErrFailedSendResponse, err),
		)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return bot_errors.New(
		s.InteractionCreate.Member.User.ID,
		bot_errors.CmdDeployRollback,
		bot_errors.ErrFailedToRecover,
	)
}
