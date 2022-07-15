package bot_errors

import (
	"discordgo"
	"fmt"
)

func NotifyUser(s *discordgo.Session, i *discordgo.InteractionCreate, e error) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occured during your request: %s", e.Error()),
		},
	})
	if err != nil {
		return ErrFailedSendResponse
	}

	return nil
}
