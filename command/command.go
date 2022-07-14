package command

import (
	"context"
	"discordgo"
)

type Command interface {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
	Recover(ctx context.Context) error
	GetErr() error
}
