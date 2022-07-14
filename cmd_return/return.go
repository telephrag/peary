package cmd_return

import (
	"context"
	"discordgo"
	"kubinka/command"
	"kubinka/config"
	"log"

	"github.com/pkg/errors"
)

type ReturnCmd struct {
	err error
	ctx context.Context
}

func Init() command.Command {
	return &ReturnCmd{}
}

func (cmd *ReturnCmd) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.GuildMemberRoleRemove(
		config.GuildID,
		i.Member.User.ID,
		config.RoleID,
	)
	if err != nil {
		msg := err.(discordgo.RESTError).Message
		log.Print(errors.Errorf("Failed to remove role: %v\n", msg))
		cmd.err = err
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have returned from deployment.",
		},
	})
	if err != nil {
		msg := err.(discordgo.RESTError).Message
		log.Println(errors.Errorf("Failed to respond to the player: %v", msg))
		cmd.err = err
		return
	}
}

func (cmd *ReturnCmd) Recover(ctx context.Context) error {
	return nil
}
func (cmd *ReturnCmd) GetErr() error {
	return cmd.err
}
