module kubinka

go 1.17

require (
	github.com/pkg/errors v0.9.1
	go.mongodb.org/mongo-driver v1.8.3
)

require (
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
	golang.org/x/text v0.3.5 // indirect
)

// Also requires discordgo placed into GOPATH/src since at the moment
// functionality to make use of slash commands wasn't released
// thus "go get" pulled outdated version without said functionality
