package cmd_return

import (
	"context"
	"errors"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/config"
	"kubinka/step"
	"time"
)

type ReturnCmd struct {
	steps     *step.Saga
	eventName string
	err       error
}

func Init(env *command.Env) command.Command {
	return &ReturnCmd{
		steps: step.NewSaga([]step.Step{
			NewRemoveRoleStep(env.DiscordSession, env.DiscordInteractionCreate),
			NewDeleteFromDBStep(env.DBConn, env.DiscordInteractionCreate),
			NewMsgResponseStep(env.DiscordSession, env.DiscordInteractionCreate),
		}),
		eventName: bot_errors.CmdReturn,
	}
}

func (cmd *ReturnCmd) Handle(ctx context.Context) *bot_errors.Nested {

	timeout := time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
	dErr := bot_errors.Nested{
		Event: bot_errors.CmdReturnDo,
	}
do:
	for cmd.steps.Next() != nil {
		s := cmd.steps.GetStep()
	retry_do:
		for {
			select {
			case <-timeout:
				if dErr.Err == nil {
					dErr.Err = errors.New(bot_errors.ErrHandlerTimeout)
				}
				break do
			default:
				dErr.Err = s.Do()
				if dErr.Err == nil {
					break retry_do
				}
			}
		}
	}
	if dErr.Err == nil {
		return nil
	}

	timeout = time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
	rErr := bot_errors.Nested{
		Event: bot_errors.CmdReturnRolback,
	}
rollback:
	for cmd.steps.Prev() != nil {
		s := cmd.steps.GetStep()
	retry_rb:
		for {
			select {
			case <-timeout:
				if rErr.Err == nil {
					rErr.Err = errors.New(bot_errors.ErrHandlerTimeout)
				}
				break rollback
			default:
				rErr.Err = s.Rollback()
				if rErr.Err == nil {
					break retry_rb
				}
			}
		}
	}

	if rErr.Err != nil {
		dErr.Next = &rErr
	}

	return &dErr
}

func (cmd *ReturnCmd) Event() string {
	return cmd.eventName
}
