package db

import (
	"context"
	"discordgo"
	"kubinka/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mi *MongoInstance) InsertPlayer(u *discordgo.User, t time.Time) error {

	expTimePrimitive := primitive.NewDateTimeFromTime(t)

	player := models.Player{
		DiscordID: u.ID,
		Expire:    expTimePrimitive,
	}

	_, err := mi.Collection.InsertOne(context.Background(), player)

	return err
}
