package cmd_deploy

import (
	"fmt"
	"kubinka/bot_errors"
	"kubinka/config"

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
		return bot_errors.New(
			s.InteractionCreate.Member.User.ID,
			bot_errors.CmdDeployDo,
			fmt.Errorf("%s: %w", bot_errors.ErrFailedGiveRole, err),
		)
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
		return bot_errors.New(
			s.InteractionCreate.Member.User.ID,
			bot_errors.CmdDeployRollback,
			fmt.Errorf("%s: %w", bot_errors.ErrFailedTakeRole, err),
		)
	}

	return nil
}
