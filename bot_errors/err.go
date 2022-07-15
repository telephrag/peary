package bot_errors

import "discordgo"

type Err struct {
	Session string // doubles as session id from taking role to losing it by one way or another
	Next    *Nested
}

func NewBotErr(i *discordgo.InteractionCreate) *Err {
	return &Err{Session: i.Member.User.ID}
}

func (e *Err) Nest(n *Nested) {
	tail := e.Next
	for tail != nil {
		tail = tail.Next
	}
	tail = n
}

type Nested struct {
	Event string
	Err   error
	Next  *Nested
}
