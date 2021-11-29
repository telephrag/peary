package handlers

import (
	"kubinka/config"

	"discordgo"
)

func Message(s *discordgo.Session, mc *discordgo.MessageCreate) {

	ch := mc.Message.ChannelID
	if ch != config.ChanID {
		return
	}

	//deploy, err := models.NewDeployFromMessage(mc.Message)

}
