package cmd_deploy

import (
	"context"
	"errors"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/config"
	"kubinka/step"
	"time"
)

type DeployCmd struct {
	steps     *step.Saga
	eventName string
	err       error
}

func Init(env *command.Env) command.Command {
	return &DeployCmd{
		steps: step.NewSaga([]step.Step{
			NewGiveRoleStep(env.DiscordSession, env.DiscordInteractionCreate),
			NewMsgResponseStep(env.DiscordSession, env.DiscordInteractionCreate),
		}),
		eventName: bot_errors.CmdDeploy,
	}
}

func (cmd *DeployCmd) Handle(ctx context.Context) *bot_errors.Nested {

	timeout := time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
	dErr := bot_errors.Nested{
		Event: bot_errors.CmdDeployDo,
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
			case <-ctx.Done():
				if dErr.Err == nil {
					dErr.Err = errors.New(bot_errors.ErrSomewhereElse)
				}
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
		Event: bot_errors.CmdDeployRollback,
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

func (cmd *DeployCmd) Event() string {
	return cmd.eventName
}
