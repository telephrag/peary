package cmd_deploy

import (
	"discordgo"
	"fmt"
	"kubinka/bot_errors"
	"kubinka/config"
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
		return fmt.Errorf(bot_errors.ErrFailedGiveRole+": %w", err)
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
		return fmt.Errorf(bot_errors.ErrFailedTakeRole+": %w", err)
	}

	return nil
}
