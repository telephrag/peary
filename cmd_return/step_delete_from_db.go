package cmd_return

import (
	"fmt"
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
	return fmt.Errorf(bot_errors.ErrFailedToRecover)
}
