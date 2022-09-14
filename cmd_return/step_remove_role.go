package cmd_return

import (
	"fmt"
	"peary/config"
	"peary/errconst"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
)

type RemoveRoleStep struct {
	DiscordSession    *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
}

func NewRemoveRoleStep(s *discordgo.Session, i *discordgo.InteractionCreate) *RemoveRoleStep {
	return &RemoveRoleStep{
		DiscordSession:    s,
		InteractionCreate: i,
	}
}

func (s *RemoveRoleStep) Do() error {
	err := s.DiscordSession.GuildMemberRoleRemove(
		s.InteractionCreate.GuildID,
		s.InteractionCreate.Member.User.ID,
		config.BOT_ROLE_ID,
	)
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedTakeRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errconst.CmdReturnDo)
	}

	return nil
}

func (s *RemoveRoleStep) Rollback() error {
	// if we removed role already, better leave it like this even if user gets no response
	// which is better than receiving pings you didn't sign for
	return errlist.New(errconst.ErrRecoveryImpossible).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errconst.CmdReturnRollback)
}
