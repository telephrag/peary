package cmd_return

import (
	"discordgo"
	"errors"
	"fmt"
	"kubinka/bot_errors"
	"kubinka/config"
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
		config.GuildID,
		s.InteractionCreate.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		return fmt.Errorf(bot_errors.ErrFailedTakeRole+": %w", err)
	}

	return nil
}

func (s *RemoveRoleStep) Rollback() error {
	// if we removed role already, better leave it like this even if user gets no response
	// which is better than receiving pings you didn't sign for
	return errors.New(bot_errors.ErrFailedToRecover)
}
