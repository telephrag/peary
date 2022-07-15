package orchestrator

import (
	"context"
	"discordgo"
	"fmt"
	"kubinka/bot_errors"
	"kubinka/config"
	"log"
	"sync/atomic"
	"time"
)

type CmdOrchestrator struct {
	Ctx          context.Context
	Cancel       context.CancelFunc
	RunningCount int32
}

func New() *CmdOrchestrator {
	ctx, cancel := context.WithCancel(context.Background())
	return &CmdOrchestrator{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func (co *CmdOrchestrator) Orchestrate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	init, ok := handlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		log.Println("Couldn't retreive command Init(): ", i.ApplicationCommandData().Name)
		return
	}

	cmd := init()

	atomic.AddInt32(&co.RunningCount, 1)
	defer atomic.AddInt32(&co.RunningCount, -1)

	ctx, _ := context.WithTimeout(co.Ctx, time.Second*config.CMD_HANDLER_TIMEOUT_SECONDS)
	err := cmd.Handle(ctx, s, i)
	if err != nil {
		log.Printf("Err in command %s: %s\n", i.ApplicationCommandData().Name, err.Error())
		err = bot_errors.NotifyUser(s, i, err)
		if err != nil {
			log.Println(err)
			co.Cancel()
		}

		err = cmd.Recover(s, i)
		if err != nil {
			err = fmt.Errorf("%s: %w", bot_errors.ErrFailedToRecover.Error(), err)
			co.Cancel()
		}
	}
}

func (co *CmdOrchestrator) HaltUntilAllDone() {
	for atomic.LoadInt32(&co.RunningCount) != 0 {
		time.Sleep(time.Millisecond * 100)
	}
}
