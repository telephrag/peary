package step

import (
	"discordgo"

	"go.mongodb.org/mongo-driver/mongo"
)

type Step interface {
	Do() error
	Rollback() error
}

type WithDB struct {
	DBConn *mongo.Client
}

type WithDiscordSession struct {
	DiscordSession *discordgo.Session
}

type WithDiscordInteractionCreate struct {
	InteractionCreate *discordgo.InteractionCreate
}
