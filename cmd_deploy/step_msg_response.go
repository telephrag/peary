package cmd_deploy

import (
	"fmt"
	"peary/errconst"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
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
				Content: fmt.Sprintf("You have been deployed till %02d:%02d UTC", d.Hour(), d.Minute()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedSendResponse, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errconst.CmdDeployDo)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return errlist.New(errconst.ErrFailedToRecover).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errconst.CmdDeployRollback)
}
