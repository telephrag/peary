package cmd_return

import (
	"context"
	"discordgo"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/config"
)

type ReturnCmd struct {
	err error
	ctx context.Context
}

func Init() command.Command {
	return &ReturnCmd{}
}

func (cmd *ReturnCmd) Handle(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.GuildMemberRoleRemove(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		return bot_errors.ErrFailedTakeRole
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have returned from deployment.",
		},
	})
	if err != nil {
		return bot_errors.ErrFailedSendResponse
	}

	return nil
}

func (cmd *ReturnCmd) Recover(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bot_errors.ErrFailedToRecover
}

func (cmd *ReturnCmd) GetErr() error {
	return cmd.err
}
