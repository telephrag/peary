package command

import (
	"context"
	"discordgo"
	"kubinka/bot_errors"
	"kubinka/strg"
)

type Command interface {
	Handle(ctx context.Context) *bot_errors.Nested
	Event() string
}

type Env struct {
	DiscordSession           *discordgo.Session
	DiscordInteractionCreate *discordgo.InteractionCreate
	DBConn                   *strg.BoltConn
}
