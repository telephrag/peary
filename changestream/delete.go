package changestream

import (
	"context"
	"discordgo"
	"kubinka/config"
	"kubinka/models"
	"log"

	"github.com/pkg/errors"
)

func Delete(ds *discordgo.Session, ctx context.Context, cancel context.CancelFunc) {

	var err error

	player, ok := ctx.Value("doc").(models.Player)
	if !ok {
		log.Print(errors.Errorf("Failed to retrieve Player: %w\n", err))
		cancel()
	}
	log.Println(player)

	err = ds.GuildMemberRoleRemove(
		config.GuildID,
		player.DiscordID,
		config.RoleID,
	)
	if err != nil {
		log.Print(errors.Errorf("Failed to remove role: %w\n", err))
		cancel()
	}
}
