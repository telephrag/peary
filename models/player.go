package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	DiscordID string             `bson:"discord_id"`
	Expire    primitive.DateTime `bson:"expire"`
}
