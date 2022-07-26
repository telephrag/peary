package bot_errors

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
		return New(i.Member.User.ID, NotifyUsr, fmt.Errorf("%s: %w", ErrFailedSendResponse, err))
	}

	return nil
}
