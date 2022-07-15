package cmd_deploy

import (
	"discordgo"
	"time"
)

func getDeployDuration(i *discordgo.InteractionCreate) time.Time {
	opt := i.ApplicationCommandData().Options
	m := opt[0].IntValue()
	if len(opt) > 1 {
		m += opt[1].IntValue() * 60
	}

	return time.Now().Add(time.Minute * time.Duration(m)).UTC()
}
