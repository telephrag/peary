package cmd_return

import (
	"kubinka/bot_errors"
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
	return bot_errors.New(
		s.InteractionCreate.Member.User.ID,
		bot_errors.CmdReturnRollback,
		bot_errors.ErrFailedToRecover,
	)
}
