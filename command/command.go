package command

import (
	"context"
	"discordgo"
)

type Command interface {
	Init(ctx context.Context) interface{}
	Handle(ds *discordgo.Session, i *discordgo.InteractionCreate)
	Recover(ctx context.Context) error
	GetErr() error
}
