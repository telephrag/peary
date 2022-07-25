module kubinka

go 1.17

require (
	github.com/bwmarrin/discordgo v0.25.0
	go.etcd.io/bbolt v1.3.6
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)

// Also requires discordgo placed into GOPATH/src since at the moment
// functionality to make use of slash commands wasn't released
// thus "go get" pulled outdated version without said functionality
