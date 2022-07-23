package models

import (
	"time"
)

type Player struct {
	DiscordID string    `json:"discord_id"`
	Expire    time.Time `json:"expire"`
}
