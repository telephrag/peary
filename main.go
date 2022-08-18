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

func deleteCommands(ds *discordgo.Session, guildId string) {
	for _, cmd := range service.CmdDef {
		err := ds.ApplicationCommandDelete(
			ds.State.User.ID,
			guildId,
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
				deleteCommands(ds, config.BOT_GUILD_ID)
			}
			log.Fatalf("Failed to create command %s:\n %s\n\n\n", cmd.Name, err)
		}
	}
}

func reissueRoles(ds *discordgo.Session, db *strg.BoltConn, guildId, roleId string) error {
	idsWithRoles := db.GetPlayerIDs()
	for _, id := range idsWithRoles {
		err := ds.GuildMemberRoleAdd(guildId, id, roleId)
		if err != nil {
			return errlist.New(err).Set("event", errlist.StartupRoleReissue).Set("session", id)
		}
	}

	return nil
}

func takeRoles(ds *discordgo.Session, db *strg.BoltConn, guildId, roleId string) {
	idsWithRoles := db.GetPlayerIDs()
	for _, id := range idsWithRoles {
		err := ds.GuildMemberRoleRemove(guildId, id, roleId)
		if err != nil {
			log.Print(errlist.New(err).Set("event", errlist.ShutdownRoleRemove).Set("session", id))
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
	defer deleteCommands(ds, config.BOT_GUILD_ID) // Removing commands on bot shutdown

	db.RemoveExpired(ds)
	reissueRoles(ds, db, config.BOT_GUILD_ID, config.BOT_ROLE_ID)
	// roles are removed on shutdown at the end of main, see bellow
	go func() {
		err := db.WatchExpirations(ctx, ds)
		log.Print(err)
		masterHandler.Cancel()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	shutdownLogRec := errlist.New(nil).Set("event", "SESSION SHUTDOWN")
halt:
	for {
		select {
		case <-interrupt:
			log.Print(shutdownLogRec.Set("cause", "execution stopped by user"))
			masterHandler.Cancel()
			log.Println(db.GetPlayerIDs())
			break halt
		case <-masterHandler.Ctx.Done():
			log.Print(shutdownLogRec.Set("cause", ctx.Err().Error()))
			break halt
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	// can't defer, see why at NOTES 04
	takeRoles(ds, db, config.BOT_GUILD_ID, config.BOT_ROLE_ID)
}
