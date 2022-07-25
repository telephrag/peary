package bot_errors

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func NotifyUser(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) *Nested {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occured during your request: %s", errMsg),
		},
	})
	if err != nil {
		return &Nested{
			Event: NotifyUsr,
			Err:   fmt.Errorf(ErrFailedSendResponse+": %w", err),
		}
	}

	return nil
}
