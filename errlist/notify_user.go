package errlist

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func NotifyUser(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occured during your request: %s", errMsg),
		},
	})
	if err != nil {
		return New(fmt.Errorf("%s: %w", ErrFailedSendResponse, err)).Set("session", i.Member.User.ID).Set("event", NotifyUsr)
	}

	return nil
}
