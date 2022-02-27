package db

import (
	"context"
	"discordgo"

	"go.mongodb.org/mongo-driver/bson"
)

func (mi *MongoInstance) DeletePlayer(u *discordgo.User) error {

	filter := bson.M{
		"discord_id": u.ID,
	}
	_, err := mi.Collection.DeleteOne(context.Background(), filter)

	return err
}
