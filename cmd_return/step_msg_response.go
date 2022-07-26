package cmd_return

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
	err := s.DiscordSession.InteractionRespond(
		s.InteractionCreate.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have returned from deployment.",
			},
		})
	if err != nil {
		return bot_errors.New(
			s.InteractionCreate.Member.User.ID,
			bot_errors.CmdReturnDo,
			fmt.Errorf("%s: %w", bot_errors.ErrFailedSendResponse, err),
		)
	}

	return nil
}

func (s *MsgResponseStep) Rollback() error {
	return bot_errors.New(
		s.InteractionCreate.Member.User.ID,
		bot_errors.CmdReturnRollback,
		bot_errors.ErrFailedToRecover,
	)
}
