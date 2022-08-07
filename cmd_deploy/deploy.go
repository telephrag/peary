package cmd_deploy

import (
	"context"
	"kubinka/command"
	"kubinka/config"
	"kubinka/errlist"
	"kubinka/step"
	"time"
)

type DeployCmd struct {
	steps     *step.Saga
	eventName string
	session   string
}

func Init(env *command.Env) command.Command {
	return &DeployCmd{
		steps: step.NewSaga([]step.Step{
			NewGiveRoleStep(env.DiscordSession, env.DiscordInteractionCreate),
			NewAddToDBStep(env.DBConn, env.DiscordInteractionCreate),
			NewMsgResponseStep(env.DiscordSession, env.DiscordInteractionCreate),
		}),
		eventName: errlist.CmdDeploy,
		session:   env.DiscordInteractionCreate.Member.User.ID,
	}
}

func (cmd *DeployCmd) Handle(ctx context.Context) error {

	var doErr error
	timeout := time.After(time.Second * config.CMD_HANDLER_TIMEOUT_SECONDS)
do: // iterate all steps in command
	for cmd.steps.Next() != nil {
		s := cmd.steps.GetStep()
	retry_do:
		for {
			select {
			case <-timeout:
				if doErr == nil {
					doErr = errlist.New(errlist.ErrHandlerTimeout).
						Set("session", cmd.session).
						Set("event", cmd.eventName)
				}
				break do
			case <-ctx.Done():
				// /deploy should not continue execution since its completion in context of failure...
				if doErr == nil { // can break state
					doErr = errlist.New(errlist.ErrHandlerTimeout).
						Set("session", cmd.session).
						Set("event", cmd.eventName)
				}
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
					rbErr = errlist.New(errlist.ErrHandlerTimeout).
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
		doErr.(*errlist.ErrNode).Wrap(rbErr)
	}

	return doErr
}

func (cmd *DeployCmd) Event() string {
	return cmd.eventName
}
