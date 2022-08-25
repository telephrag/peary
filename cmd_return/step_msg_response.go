package cmd_return

import (
	"fmt"
	"kubinka/errlist"

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
	err := s.DiscordSession.InteractionRespond(
		s.InteractionCreate.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have returned from deployment.",
				Flags:   uint64(1 << 6),
			},
		})
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errlist.ErrFailedSendResponse, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errlist.CmdReturnDo)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return errlist.New(errlist.ErrFailedToRecover).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errlist.CmdReturnRollback)
}
