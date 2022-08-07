package cmd_return

import (
	"kubinka/errlist"
	"kubinka/strg"

	"github.com/bwmarrin/discordgo"
)

type DeleteFromDBStep struct {
	DBConn            *strg.BoltConn
	InteractionCreate *discordgo.InteractionCreate
}

func NewDeleteFromDBStep(dbConn *strg.BoltConn, i *discordgo.InteractionCreate) *DeleteFromDBStep {
	return &DeleteFromDBStep{
		DBConn:            dbConn,
		InteractionCreate: i,
	}
}

func (s *DeleteFromDBStep) Do() error {
	return s.DBConn.Delete(s.InteractionCreate.Member.User.ID)
}

func (s *DeleteFromDBStep) Rollback() error {
	return errlist.New(errlist.ErrFailedToRecover).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errlist.CmdReturnRollback)
}
