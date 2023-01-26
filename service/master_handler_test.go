package service

import (
	"context"
	"errors"
	"peary/command"
	"peary/step"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var (
	errFailedDo      = errors.New("failed")
	errFailedRecover = errors.New("failed but recovered")
)

const (
	cmdFailName    = "cmd_fail"
	cmdFailSession = "failedSession"
)

//
// command that fails
//
type failCmd struct {
	steps     *step.Saga
	eventName string
	session   string
}

func initFailCmd(env *command.Env) command.Command {

	return &failCmd{
		steps:     step.NewSaga([]step.Step{&failingStep{}}),
		eventName: cmdFailName,
		session:   cmdFailSession,
	}
}

func (cmd *failCmd) Handle(ctx context.Context) error { return nil } // TODO

func (cmd *failCmd) Event() string { return cmd.eventName }

//
// step that fails and doesn't recover
//
type failingStep struct{}

func (s *failingStep) Do() error { return errFailedDo }

func (s *failingStep) Rollback() error { return errFailedRecover }

//
// step that recovers after failure
//
type recoveringStep struct{}

func (s *recoveringStep) Do() error { return errFailedDo }

func (s *recoveringStep) Rollback() error { return nil }

//
// substitutes discordgo.InteractionCreate
//
func fakeInteractionCreate(cmdName, usrID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: cmdName,
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID: usrID,
				},
			},
		},
	}
}

func TestFailCmd(t *testing.T) {
	handlerToCmd = map[string]func(*command.Env) command.Command{
		"cmd_fail": initFailCmd,
	}

	ctx, cancel := context.WithCancel(context.Background())
	master := NewMasterHandler(ctx, cancel, nil)

	// log.Print(fakeInteractionCreate(cmdFailName, cmdFailSession).ApplicationCommandData().Name)

	master.Handle(nil, fakeInteractionCreate(cmdFailName, cmdFailSession))

	// TODO: make io.Writer that will log shit into variable you can compare expected results against
}
