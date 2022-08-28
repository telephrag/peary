package command

import (
	"context"
	"peary/strg"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	/* Cancelled context means that state has been broken somewhere. Execution of some commands means to risk breaking state even more. In these commands `<-ctx.Done()` should be called before every Do() (see pkg `step`) */
	Handle(ctx context.Context) error

	/* Returns unique event identifier that is used in errors to provide context
	   of where error has occured. */
	Event() string
}

type Env struct {
	DiscordSession           *discordgo.Session
	DiscordInteractionCreate *discordgo.InteractionCreate
	DBConn                   *strg.BoltConn
}
