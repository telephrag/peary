package cmd_deploy

import (
	"kubinka/models"
	"kubinka/strg"

	"github.com/bwmarrin/discordgo"
)

type AddToDBStep struct {
	DBConn            *strg.BoltConn
	InteractionCreate *discordgo.InteractionCreate
}

func NewAddToDBStep(dbConn *strg.BoltConn, i *discordgo.InteractionCreate) *AddToDBStep {
	return &AddToDBStep{
		DBConn:            dbConn,
		InteractionCreate: i,
	}
}

func (s *AddToDBStep) Do() error {
	return s.DBConn.Insert(&models.Player{
		DiscordID: s.InteractionCreate.Member.User.ID,
		Expire:    getDeployDuration(s.InteractionCreate),
	})
}

func (s *AddToDBStep) Rollback() error {
	return s.DBConn.Delete(s.InteractionCreate.Member.User.ID)
}
