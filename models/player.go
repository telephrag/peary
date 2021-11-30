package models

import (
	"time"
)

type Player struct {
	UserID   string
	UserName string
	Duration int
	Begin    time.Time
	End      time.Time
	// timestamp discordgo.Timestamp?
}
