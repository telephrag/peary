module kubinka

go 1.17

require (
    github.com/gorilla/websocket v1.4.0 // indirect
    golang.org/x/crypto v0.0.0-20181030102418-4d3f4d9ffa16 // indirect
)

// Also requires discordgo placed into GOPATH/src since at the moment 
// functionality to make use of slash commands wasn't released 
// thus "go get" pulled outdated version without said functionality
