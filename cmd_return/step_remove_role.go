package cmd_return

import (
	"fmt"
	"peary/config"
	"peary/errlist"

	"github.com/bwmarrin/discordgo"
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
		return errlist.New(fmt.Errorf("%s: %w", errlist.ErrFailedTakeRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errlist.CmdReturnDo)
	}

	return nil
}

func (s *RemoveRoleStep) Rollback() error {
	// if we removed role already, better leave it like this even if user gets no response
	// which is better than receiving pings you didn't sign for
	return errlist.New(errlist.ErrFailedToRecover).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errlist.CmdReturnRollback)
}
