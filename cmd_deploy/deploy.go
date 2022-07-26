package cmd_deploy

import (
	"context"
	"kubinka/bot_errors"
	"kubinka/command"
	"kubinka/config"
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
		eventName: bot_errors.CmdDeploy,
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
					doErr = bot_errors.New(cmd.session, cmd.eventName, bot_errors.ErrHandlerTimeout)
				}
				break do
			case <-ctx.Done():
				// /deploy should not continue execution since its completion in context of failure...
				if doErr == nil { // can break state
					doErr = bot_errors.New(cmd.session, cmd.eventName, bot_errors.ErrSomewhereElse)
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
					rbErr = bot_errors.New(cmd.session, cmd.eventName, bot_errors.ErrHandlerTimeout)
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
		doErr.(*bot_errors.Err).Nest(rbErr)
	}

	return doErr
}

func (cmd *DeployCmd) Event() string {
	return cmd.eventName
}
