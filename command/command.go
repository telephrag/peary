package command

import (
	"context"
	"discordgo"
)

type Command interface {
	Handle(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error
	Recover(s *discordgo.Session, i *discordgo.InteractionCreate) error
	GetErr() error
}
