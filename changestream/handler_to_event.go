package changestream

import (
	"context"
	"discordgo"
	"log"
)

var handlerToEvent = map[string]func(ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc){
	"insert": func(ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc) {
		// TODO: Consider notifying about someones deploy
		log.Println("insert handling...")
	},
	"delete": Delete,
	"invalidate": func(ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc) {
		// TODO: Take away roles from everyone
		log.Println("invalidate handling...")
		cancel()
	},
}
