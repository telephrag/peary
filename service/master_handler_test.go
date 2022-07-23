package service

import (
	"errors"
	"kubinka/bot_errors"
	"testing"
	"time"
)

func TestMasterHandler(t *testing.T) {
	mh := NewMasterHandler(nil)
	go func() {
		time.Sleep(time.Millisecond * 5)
		mh.Cancel()
	}()

	mhErr := bot_errors.Err{
		Session: "1337",
		Event:   bot_errors.CmdDeploy,
	}

	time.Sleep(time.Millisecond * 10)

	select {
	case <-mh.Ctx.Done():
		mhErr.Nest(&bot_errors.Nested{
			Event: bot_errors.CmdDeploy,
			Err:   errors.New(bot_errors.ErrSomewhereElse),
		})

		// err := bot_errors.NotifyUser(nil, i, bot_errors.ErrSomewhereElse)
		// if err != nil {
		// 	mhErr.Nest(err)
		// 	mh.Cancel()
		// }

		return
	default:
	}

}
