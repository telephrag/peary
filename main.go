package main

import (
	"context"
	"kubinka/config"
	"kubinka/service"
	"kubinka/strg"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
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

func newDiscordSession(
	ctx context.Context,
	cancel context.CancelFunc,
	token string,
	c *strg.BoltConn,
) (*discordgo.Session, *service.MasterHandler) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Could not create session.\n\n\n")
	}
	discord.SyncEvents = false
	masterHandler := service.NewMasterHandler(ctx, cancel, c)
	discord.AddHandler(masterHandler.Handle)

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
	logFile := getLogFile(config.LOG_FILE_NAME)
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Print("SESSION STARTUP\n")
	defer log.Print("SESSION SHUTDOWN\n\n\n")

	db, err := strg.Connect(config.DB_NAME, config.DB_PLAYERS_BUCKET_NAME)
	if err != nil {
		log.Panicf("failed to connect to db: %v", err)
	}
	defer db.Close() // works fine, wtf ???
	// defer func() {
	// 	if err := db.Close(); err != nil { // will close normally in debug mode
	// 		log.Panicf("error closing db conn: %v", err) // will stuck otherwise
	// 	}
	// }()

	ctx, cancel := context.WithCancel(context.Background())

	ds, masterHandler := newDiscordSession(ctx, cancel, config.BOT_TOKEN, db)
	defer masterHandler.Cancel()
	// defer masterHandler.HaltUntilAllDone()
	defer ds.Close()

	createCommands(ds, config.BOT_APP_ID, config.BOT_GUILD_ID)
	defer deleteCommands(ds) // Removing commands on bot shutdown

	go func() {
		err = db.WatchExpirations(ctx, ds)
		if err != nil {
			log.Printf("error while watching deployments expirations: %s", err.Error())
			masterHandler.Cancel()
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-interrupt:
			log.Println("Execution stopped by user")
			masterHandler.Cancel()
			return // why return doesn't work here?
			// or does it? cause I've seen break work only in debug mode
		case <-masterHandler.Ctx.Done():
			log.Println("ctx cancelled")
			return
			// default:
			// 	time.Sleep(time.Millisecond * 100)
		}
	}
}
