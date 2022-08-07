package service

import (
	"context"
	"kubinka/command"
	"kubinka/errlist"
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

	select {
	case <-mh.Ctx.Done(): // cancellation of context means breakage of state somewhere...
		mhErr := errlist.New(errlist.ErrSomewhereElse).
			Set("session", i.Member.User.ID).
			Set("event", cmd.Event()).
			Wrap(errlist.NotifyUser(s, i, errlist.ErrSomewhereElse.Error()))
		log.Print(mhErr)
		return // ... so, do not handle any more commands to not risk breaking state even more
	default:
	}

	atomic.AddInt32(&mh.RunningCount, 1)
	defer atomic.AddInt32(&mh.RunningCount, -1)

	handlerErr := cmd.Handle(mh.Ctx)
	if handlerErr != nil {
		err := errlist.NotifyUser(s, i, handlerErr.Error())
		if err != nil {
			handlerErr.(*errlist.ErrNode).Wrap(err)
			mh.Cancel()
		}
		log.Print(handlerErr)
		return
	}

	log.Print(errlist.New(nil).Set("session", i.Member.User.ID).Set("event", cmd.Event()))
}

func (mh *MasterHandler) HaltUntilAllDone() {
	for atomic.LoadInt32(&mh.RunningCount) != 0 {
		time.Sleep(time.Millisecond * 100)
	}
}
