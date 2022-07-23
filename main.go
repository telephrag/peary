package main

import (
	"discordgo"
	"kubinka/config"
	"kubinka/service"
	"kubinka/strg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getLogFile(fileName string) *os.File {
	// setting up logging, for some reason logging wont work properly
	// if it was setup inside init()
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Failed to open file for logging.\n\n\n")
	}
	return f
}

func newDiscordSession(token string, c *strg.BoltConn) (*discordgo.Session, *service.MasterHandler) {
	discord, err := discordgo.New("Bot " + token) // TODO: Make to arg
	if err != nil {
		log.Fatal("Could not create session.\n\n\n")
	}
	discord.SyncEvents = false
	masterHandler := service.NewMasterHandler(c)
	discord.AddHandler(masterHandler.Handle) // see "notes 02" in NOTES.md

	err = discord.Open()
	if err != nil {
		log.Fatal("Could not open connection.\n\n\n")
	}

	return discord, masterHandler
}

func deleteCommands(ds *discordgo.Session) { // make stuff passed in as params
	for _, cmd := range service.CmdDef {
		err := ds.ApplicationCommandDelete(
			ds.State.User.ID,
			config.BOT_GUILD_ID,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete %q command: %v\n\n\n", cmd.Name, err)
		}
	}
}

func createCommands(ds *discordgo.Session, appId string, guildId string) {
	var err error
	for i, cmd := range service.CmdDef {
		service.CmdDef[i], err = ds.ApplicationCommandCreate(
			appId,
			guildId,
			cmd,
		)
		if err != nil {
			if i > 0 {
				deleteCommands(ds)
			}
			log.Fatalf("Failed to create command %s:\n %s\n\n\n", cmd.Name, err)
		}
	}
}

func main() {
	log.SetOutput(getLogFile(config.LOG_FILE_NAME))
	log.Print("SESSION STARTUP\n")
	defer log.Print("SESSION SHUTDOWN\n\n\n")

	db, err := strg.Connect(config.DB_NAME, config.PLAYERS_COLLECTION_NAME)
	if err != nil {
		log.Panicf("failed to connect to db: %v", err)
	}

	ds, masterHandler := newDiscordSession(config.BOT_TOKEN, db)
	defer masterHandler.HaltUntilAllDone()
	defer ds.Close()

	createCommands(ds, config.BOT_APP_ID, config.BOT_GUILD_ID)
	defer deleteCommands(ds) // Removing commands on bot shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	// handling invalidation of collection at shutdown
	for {
		select {
		case <-interrupt:
			log.Println("Execution stopped by user")
			return
		case <-masterHandler.Ctx.Done():
			masterHandler.HaltUntilAllDone()
		default:
		}
		time.Sleep(time.Millisecond * 500)
	}
}
