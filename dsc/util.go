package dsc

import (
	"discordgo"
	"log"
)

func logCommand(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	log.Println(
		i.ApplicationCommandData().Name,
		i.Member.User.ID,
		i.Member.User.Username,
		err,
	)

	// r := recover()
	// if r != nil {
	// 	log.Printf("%s failed, retrying...\n", i.ApplicationCommandData().Name)
	// 	//HandlerToCmd[i.ApplicationCommandData().Name](, i) // TODO: Make discord session a singleton?
	// 	// initialization cycle error? wtf
	// }
}
