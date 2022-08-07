package main

import (
	"context"
	"kubinka/config"
	"kubinka/errlist"
	"kubinka/service"
	"kubinka/strg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getLogFile(fileName string) *os.File {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		_, err := os.Create(fileName)
		if err != nil {
			log.Fatalln("Failed to create or open file for logging.")
		}
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Failed to open just created file.")
		}
		return f
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
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)
	log.Print(errlist.New(nil).Set("event", "SESSION STARTUP"))

	db, err := strg.Connect(config.DB_NAME, config.DB_PLAYERS_BUCKET_NAME)
	if err != nil {
		log.Panicf("failed to connect to db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())

	ds, masterHandler := newDiscordSession(ctx, cancel, config.BOT_TOKEN, db)
	defer masterHandler.Cancel()
	defer masterHandler.HaltUntilAllDone()
	defer ds.Close()

	createCommands(ds, config.BOT_APP_ID, config.BOT_GUILD_ID)
	defer deleteCommands(ds) // Removing commands on bot shutdown

	go func() {
		err := db.WatchExpirations(ctx, ds)
		log.Print(err)
		masterHandler.Cancel()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	shutdownLogRec := errlist.New(nil).Set("event", "SESSION SHUTDOWN")
	for {
		select {
		case <-interrupt:
			log.Print(shutdownLogRec.Set("cause", "execution stopped by user"))
			masterHandler.Cancel()
			return
		case <-masterHandler.Ctx.Done():
			log.Print(shutdownLogRec.Set("cause", ctx.Err().Error()))
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
