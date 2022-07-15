package command

import (
	"context"
	"discordgo"
	"kubinka/bot_errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Command interface {
	Handle(ctx context.Context) *bot_errors.Nested
	Event() string
}

type Env struct {
	DiscordSession           *discordgo.Session
	DiscordInteractionCreate *discordgo.InteractionCreate
	DBConn                   *mongo.Client
}
