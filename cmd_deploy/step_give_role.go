package cmd_deploy

import (
	"fmt"
	"kubinka/config"
	"kubinka/errlist"

	"github.com/bwmarrin/discordgo"
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
		config.BOT_GUILD_ID,
		s.InteractionCreate.Member.User.ID,
		config.BOT_ROLE_ID,
	)
	if err != nil {
		// TODO: Check if actual error is outputted not just bytes
		return errlist.New(fmt.Errorf("%s: %w", errlist.ErrFailedGiveRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errlist.CmdDeployDo)
	}

	return nil
}

func (s *GiveRoleStep) Rollback() error {
	err := s.DiscordSession.GuildMemberRoleRemove(
		config.BOT_GUILD_ID,
		s.InteractionCreate.Member.User.ID,
		config.BOT_ROLE_ID,
	)
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errlist.ErrFailedTakeRole, err)).
			Set("session", s.InteractionCreate.Member.User.ID).
			Set("event", errlist.CmdDeployRollback)
	}

	return nil
}
