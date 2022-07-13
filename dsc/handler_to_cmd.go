package dsc

import "discordgo"

var HandlerToCmd = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"deploy": Deploy,
	"return": Return,
}
