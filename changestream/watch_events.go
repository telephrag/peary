package changestream

import (
	"context"
	"discordgo"
	"fmt"
	"kubinka/db"
	"kubinka/models"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type csEvent struct {
	OperationType string        `bson:"operationType"`
	FullDocument  models.Player `bson:fullDocument`
}

func iterateChangeStream(stream *mongo.ChangeStream, ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc) {
	defer stream.Close(ctx)
	defer cancel()

	for stream.Next(ctx) {

		var event csEvent
		err := stream.Decode(&event)
		if err != nil {
			log.Print(errors.Errorf("Failed to decode event: %w\n", err))
			cancel()
			return
		}

		// empty ns (namespace), fullDocument, documentKey

		fmt.Println(event.FullDocument)
		handlerCtx := context.WithValue(ctx, "doc", event.FullDocument)
		go handlerToEvent[event.OperationType](ds, handlerCtx, cancel) // TODO: Thread pool here
	}
}

func WatchEvents(ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc) {

	pipeline := mongo.Pipeline{
		bson.D{{
			"$match", bson.D{{
				"$or", bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "delete"}},
					bson.D{{"operationType", "invalidate"}},
				},
			}},
		}},
	}
	// need to create doc than update it to get fullDocument?
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, err := db.Instance.Collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Panic(err)
	}
	defer stream.Close(ctx)

	iterateChangeStream(stream, ds, ctx, cancel)
}
