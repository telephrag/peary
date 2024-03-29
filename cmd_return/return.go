package cmd_return

import (
	"context"
	"peary/command"
	"peary/config"
	"peary/errconst"
	"peary/step"
	"time"

	"github.com/telephrag/errlist"
)

type ReturnCmd struct {
	steps     *step.Saga
	eventName string
	session   string
}

func Init(env *command.Env) command.Command {
	return &ReturnCmd{
		steps: step.NewSaga([]step.Step{
			NewRemoveRoleStep(env.DiscordSession, env.DiscordInteractionCreate),
			NewDeleteFromDBStep(env.DBConn, env.DiscordInteractionCreate),
			NewMsgResponseStep(env.DiscordSession, env.DiscordInteractionCreate),
		}),
		eventName: errconst.CmdReturn,
		session:   env.DiscordInteractionCreate.Member.User.ID,
	}
}

// /return completion is beneficial since bot won't be left in a broken state with someone still
// being having role. Hence `ctx` here is not used.
func (cmd *ReturnCmd) Handle(ctx context.Context) error {

	var doErr error
	timeout := time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
do: // iterate all steps of command
	for cmd.steps.Next() != nil {
		s := cmd.steps.GetStep()
	retry_do:
		for {
			select {
			case <-timeout:
				if doErr == nil {
					doErr = errlist.New(errconst.ErrHandlerTimeout).
						Set("session", cmd.session).
						Set("event", cmd.eventName)
				}
				break do
			default:
				doErr = s.Do()
				if doErr == nil {
					break retry_do
				}
			}
		}
	}
	if doErr == nil {
		return nil
	}

	timeout = time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
	var rbErr error
rollback: // reverse iterate from point of failure
	for cmd.steps.Prev() != nil {
		s := cmd.steps.GetStep()
	retry_rb:
		for {
			select {
			case <-timeout:
				if rbErr == nil {
					rbErr = errlist.New(errconst.ErrHandlerTimeout).
						Set("session", cmd.session).
						Set("event", cmd.eventName)
				}
				break rollback
			default:
				rbErr = s.Rollback()
				if rbErr == nil {
					break retry_rb
				}
			}
		}
	}

	if rbErr != nil {
		doErr.(*errlist.ErrNode).Wrap(errconst.ErrFailedToRecover).Wrap(rbErr)
	}

	return doErr
}

func (cmd *ReturnCmd) Event() string {
	return cmd.eventName
}
