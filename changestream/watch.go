package changestream

import (
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func iterateChangeStream(stream *mongo.ChangeStream, ctx context.Context, cancel context.CancelFunc) {
	defer stream.Close(ctx)

	for stream.Next(ctx) {
		var event bson.M
		err := stream.Decode(&event)
		if err != nil {
			log.Panic(err)
		}

		ra := reflect.ValueOf(event["operationType"])
		opTypte, ok := ra.Interface().(string)
		if !ok {
			log.Panic("string expected")
		}

		handlerToEvent[opTypte](ctx, cancel)
	}
}

func WatchEvents(collection *mongo.Collection, ctx context.Context, cancel context.CancelFunc) {

	pipeline := mongo.Pipeline{
		bson.D{{
			"$match",
			bson.D{{
				"$or", bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "delete"}},
					bson.D{{"operationType", "invalidate"}},
				},
			}},
		}},
	}

	stream, err := collection.Watch(ctx, pipeline)
	if err != nil {
		log.Panic(err)
	}
	defer stream.Close(context.TODO())

	iterateChangeStream(stream, ctx, cancel)
}
