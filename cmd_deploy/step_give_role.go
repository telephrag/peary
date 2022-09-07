package cmd_deploy

import (
	"fmt"
	"peary/config"
	"peary/errconst"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
)

type GiveRoleStep struct {
	DiscordSession    *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
}

func NewGiveRoleStep(s *discordgo.Session, i *discordgo.InteractionCreate) *GiveRoleStep {
	return &GiveRoleStep{
		DiscordSession:    s,
		InteractionCreate: i,
	}
}

func (s *GiveRoleStep) Do() error {
	err := s.DiscordSession.GuildMemberRoleAdd( // TODO: discordgo.RESTError?
		s.InteractionCreate.GuildID,
		s.InteractionCreate.Member.User.ID,
		config.BOT_ROLE_ID,
	)
	if err != nil {
		// TODO: Check if actual error is outputted not just bytes
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedGiveRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errconst.CmdDeployDo)
	}

	return nil
}

func (s *GiveRoleStep) Rollback() error {
	err := s.DiscordSession.GuildMemberRoleRemove(
		s.InteractionCreate.GuildID,
		s.InteractionCreate.Member.User.ID,
		config.BOT_ROLE_ID,
	)
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedTakeRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errconst.CmdDeployRollback)
	}

	return nil
}
