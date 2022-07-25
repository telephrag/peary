package service

import (
	"context"
	"errors"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/strg"
	"log"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MasterHandler struct {
	Ctx          context.Context
	Cancel       context.CancelFunc
	RunningCount int32
	DBConn       *strg.BoltConn
}

func NewMasterHandler(ctx context.Context, cancel context.CancelFunc, dbConn *strg.BoltConn) *MasterHandler {
	return &MasterHandler{
		Ctx:    ctx,
		Cancel: cancel,
		DBConn: dbConn,
	}
}

func (mh *MasterHandler) getEnv(s *discordgo.Session, i *discordgo.InteractionCreate) *command.Env {
	return &command.Env{
		DiscordSession:           s,
		DiscordInteractionCreate: i,
		DBConn:                   mh.DBConn,
	}
}

func (mh *MasterHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	init, ok := handlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		log.Println("Couldn't retreive command Init(): ", i.ApplicationCommandData().Name)
		return
	}

	cmd := init(mh.getEnv(s, i))
	mhErr := bot_errors.Err{
		Session: i.Member.User.ID,
		Event:   cmd.Event(),
	}

	select {
	case <-mh.Ctx.Done(): // cancellation of context means breakage of state somewhere
		mhErr.Nest(&bot_errors.Nested{
			Event: cmd.Event(),
			Err:   errors.New(bot_errors.ErrSomewhereElse),
		})

		err := bot_errors.NotifyUser(s, i, bot_errors.ErrSomewhereElse)
		if err != nil {
			mhErr.Nest(err)
			mh.Cancel()
		}
		log.Print(mhErr.String())
		return // do not handle any more commands to not risk breaking state even more
	default:
	}

	atomic.AddInt32(&mh.RunningCount, 1)
	defer atomic.AddInt32(&mh.RunningCount, -1)

	err := cmd.Handle(mh.Ctx)
	if err != nil {
		mhErr.Nest(err)
		err := bot_errors.NotifyUser(s, i, err.Next.Err.Error())
		if err != nil {
			mhErr.Nest(err)
			mh.Cancel()
		}
	}
	log.Print(mhErr.String())

}

func (mh *MasterHandler) HaltUntilAllDone() {
	for atomic.LoadInt32(&mh.RunningCount) != 0 {
		time.Sleep(time.Millisecond * 100)
	}
}
