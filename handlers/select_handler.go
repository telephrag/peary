package handlers

import (
	"discordgo"
	"fmt"
)

func Select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h, ok := HandlerToCmd[i.ApplicationCommandData().Name]
	if !ok {
		fmt.Println(
			"Coudn't retrieve handler for command: ",
			i.ApplicationCommandData().Name,
		)
		return
	}

	h(s, i)
}
