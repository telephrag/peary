package cmd_return

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
	err := s.DiscordSession.InteractionRespond(
		s.InteractionCreate.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have returned from deployment.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedSendResponse, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errconst.CmdReturnDo)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return errlist.New(errconst.ErrRecoveryImpossible).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errconst.CmdReturnRollback)
}
