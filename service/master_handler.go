package service

import (
	"context"
	"fmt"
	"log"
	"peary/command"
	"peary/errconst"
	"peary/strg"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
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

func notifyUser(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occured during your request: %s", errMsg),
		},
	})
	if err != nil {
		return errlist.New(fmt.Errorf("%s: %w", errconst.ErrFailedSendResponse, err)).
			Set("session", i.Member.User.ID).Set("event", errconst.NotifyUsr)
	}

	return nil
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
		mhErr := errlist.New(errconst.ErrSomewhereElse).
			Set("session", i.Member.User.ID).
			Set("event", cmd.Event()).
			Wrap(notifyUser(s, i, errconst.ErrSomewhereElse.Error()))
		log.Print(mhErr)
		return // ... so, do not handle any more commands to not risk breaking state even more
	default:
	}

	atomic.AddInt32(&mh.RunningCount, 1)
	defer atomic.AddInt32(&mh.RunningCount, -1)

	handlerErr := cmd.Handle(mh.Ctx)
	if handlerErr != nil {
		err := notifyUser(s, i, handlerErr.Error())
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
