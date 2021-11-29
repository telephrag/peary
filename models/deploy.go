package models

import (
	"errors"
	"strings"

	"discordgo"
)

type Deploy struct {
	UserID   string
	UserName string
	Duration int
	Begin    discordgo.Timestamp
	End      discordgo.Timestamp
	// timestamp discordgo.Timestamp?
}

func New(duration int, userId, userName string, begin, end discordgo.Timestamp) *Deploy {
	return &Deploy{
		Duration: duration,
		UserID:   userId,
		UserName: userName,
		Begin:    begin,
		End:      end,
	}
}

func NewDeployFromMessage(m *discordgo.Message) (*Deploy, error) {
	content := strings.Split(m.Content, " ")

	if content[0] != "!deploy" {
		return nil, errors.New("")
	}

	return nil, nil
}
