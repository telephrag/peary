package cmd_return

import (
	"peary/errconst"
	"peary/strg"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
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
	return errlist.New(errconst.ErrRecoveryImpossible).
		Set("session", s.InteractionCreate.Member.User.ID).
		Set("event", errconst.CmdReturnRollback)
}
