package changestream

import (
	"context"
	"log"
)

var handlerToEvent = map[string]func(ctx context.Context, cancel context.CancelFunc){
	"insert": func(ctx context.Context, cancel context.CancelFunc) {
		// TODO: Consider notifying about someones deploy
		log.Println("insert handling...")
	},
	"delete": func(ctx context.Context, cancel context.CancelFunc) {
		// TODO: Retrieve deleted document and take away role from specific member
		log.Println("delete handling...")
	},
	"invalidate": func(ctx context.Context, cancel context.CancelFunc) {
		// TODO: Take away roles from everyone
		log.Println("invalidate handling...")
		cancel()
	},
}
