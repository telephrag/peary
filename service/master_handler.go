package service

import (
	"context"
	"discordgo"
	"errors"
	"kubinka/bot_errors"
	"kubinka/command"
	"log"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MasterHandler struct {
	Ctx          context.Context
	Cancel       context.CancelFunc
	RunningCount int32
	DBConn       *mongo.Client
}

func NewMasterHandler() *MasterHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &MasterHandler{
		Ctx:    ctx,
		Cancel: cancel,
		// DBConn: nil, *TODO*
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
	mhErr := bot_errors.NewBotErr(i)

	select {
	case <-mh.Ctx.Done():
		mhErr.Nest(&bot_errors.Nested{
			Event: cmd.Event(),
			Err:   errors.New(bot_errors.ErrSomewhereElse),
		})
		// move SessionID into command.Command?
		// make anything that returns error modify BotError instead to set field values?
		// preserve bot_errors.ErrSomewhereElse here
		// ErrSomewhereElse -- error in command itself
		//  L Error inside NotifyUser()
		err := bot_errors.NotifyUser(s, i, bot_errors.ErrSomewhereElse)
		if err != nil {
			mhErr.Nest(err)
			// TODO
			mh.Cancel()
			return
		}
	default:
	}

	atomic.AddInt32(&mh.RunningCount, 1)
	defer atomic.AddInt32(&mh.RunningCount, -1)

	err := cmd.Handle(mh.Ctx)
	if err != nil {
		mhErr.Nest(err)
		// log.Printf("Err in command %s: %s\n", i.ApplicationCommandData().Name, err.Error())
		err := bot_errors.NotifyUser(s, i, err.Next.Err.Error())
		if err != nil {
			mhErr.Nest(err)
			log.Println(err)
			mh.Cancel()
		}
	}
}

func (mh *MasterHandler) HaltUntilAllDone() {
	for atomic.LoadInt32(&mh.RunningCount) != 0 {
		time.Sleep(time.Millisecond * 100)
	}
	mh.Cancel()
}
