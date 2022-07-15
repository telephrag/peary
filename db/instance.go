package db

import (
	"context"
	"kubinka/config"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Instance *MongoInstance = New(
	config.MONGO_URI,
	config.MONGO_DB_NAME,
	config.MONGO_COLLECTION_NAME,
)

type MongoInstance struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New(uri, dbName, colName string) *MongoInstance {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic(err)
	}

	collection := client.Database(dbName).Collection(colName)

	// create ttl index that'll handle automatic deletion on deployment end
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"expire", 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},

		// {
		// 	Keys:    bson.D{{"discord_id", 1}},
		// 	Options: nil,
		// },
	}
	_, err = collection.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		log.Panic(err)
	}

	ms := &MongoInstance{
		Client:     client,
		Collection: collection,
	}

	return ms
}
