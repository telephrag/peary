package cmd_deploy

import (
	"discordgo"
	"errors"
	"fmt"
	"kubinka/bot_errors"
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
		return fmt.Errorf(bot_errors.ErrFailedSendResponse+": %w", err)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return errors.New(bot_errors.ErrFailedToRecover)
}
