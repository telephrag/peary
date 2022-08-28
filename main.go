package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"peary/config"
	"peary/errlist"
	"peary/service"
	"peary/strg"
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
	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Could not create session.\n\n\n")
	}
	ds.SyncEvents = false
	masterHandler := service.NewMasterHandler(ctx, cancel, c)
	ds.AddHandler(masterHandler.Handle)

	err = ds.Open()
	if err != nil {
		log.Fatal(errlist.New(fmt.Errorf("could not open connection: %w", err)))
	}

	if len(ds.State.Ready.Guilds) > 1 {
		log.Fatal(errlist.New(fmt.Errorf("attempt to use bot in more than one guild simultaneously")))
	}

	return ds, masterHandler
}

func createCommands(ds *discordgo.Session, appId string, guildID string) {
	var err error
	for i, cmd := range service.CmdDef {
		service.CmdDef[i], err = ds.ApplicationCommandCreate(
			appId,
			guildID,
			cmd,
		)
		if err != nil {
			if i > 0 {
				deleteCommands(ds, guildID)
			}
			log.Fatalf("Failed to create command %s:\n %s\n\n\n", cmd.Name, err)
		}
	}
}

func deleteCommands(ds *discordgo.Session, guildID string) {
	for _, cmd := range service.CmdDef {
		err := ds.ApplicationCommandDelete(
			ds.State.User.ID,
			guildID,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete %q command: %v\n\n\n", cmd.Name, err)
		}
	}
}

// ds.GuildRoleDelete(guildId, roleId) to rollback
func createRole(ds *discordgo.Session, name, guildID string, color int) (roleId string, err error) {
	var perm int64 = 0
	var yes bool = true
	st, err := ds.GuildRoleCreate(
		guildID,
		&discordgo.RoleParams{
			Name:        name,
			Color:       &color,
			Hoist:       &yes,
			Permissions: &perm,
			Mentionable: &yes,
		},
	)
	if err != nil {
		return "", err
	}

	return st.ID, nil
}

func reissueRoles(ds *discordgo.Session, db *strg.BoltConn, guildID, roleID string) error {
	idsWithRoles := db.GetPlayerIDs()
	for _, id := range idsWithRoles {
		err := ds.GuildMemberRoleAdd(guildID, id, roleID)
		if err != nil {
			return errlist.New(err).Set("event", errlist.StartupRoleReissue).Set("session", id)
		}
		log.Print(errlist.New(nil).Set("session", id).Set("event", errlist.StartupRoleReissue))
	}

	return nil
}

func main() {
	logFile := getLogFile(config.LOG_FILE_NAME)
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)
	log.Print(errlist.New(nil).Set("event", "SESSION STARTUP"))
	shutdownLogRec := errlist.New(nil).Set("event", "SESSION SHUTDOWN")
	defer log.Print(shutdownLogRec)

	db, err := strg.Connect(config.DB_NAME, config.DB_PLAYERS_BUCKET_NAME)
	if err != nil {
		shutdownLogRec.Wrap(err).Set("event", "startup_db_connect")
	} else {
		defer db.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ds, masterHandler := newDiscordSession(ctx, cancel, config.BOT_TOKEN, db)
	defer masterHandler.Cancel()
	defer masterHandler.HaltUntilAllDone()
	defer ds.Close()

	// app, err := ds.Application(config.BOT_APP_ID)
	// if err != nil {
	// 	shutdownLogRec.Wrap(err).Set("event", "startup_appdata_retrieve")
	// 	return
	// }
	// config.BOT_GUILD_ID = app.GuildID

	// config.BOT_GUILD_ID, err = ds.State.GuildOneID()
	// if err != nil {
	// 	shutdownLogRec.Wrap(err).Set("event", "startup_appdata_retrieve")
	// 	return
	// }

	guildID := ds.State.Ready.Guilds[0].ID

	createCommands(ds, config.BOT_APP_ID, guildID)
	defer deleteCommands(ds, guildID) // Removing commands on bot shutdown

	config.BOT_ROLE_ID, err = createRole(ds, "Waiting deploy", guildID, 307015)
	if err != nil {
		shutdownLogRec.Wrap(
			errlist.New(fmt.Errorf("failed to create role: %w", err)).
				Set("event", "startup_role_create"))
		return
	}
	defer func() {
		if err := ds.GuildRoleDelete(guildID, config.BOT_ROLE_ID); err != nil {
			shutdownLogRec.Wrap(
				errlist.New(fmt.Errorf("failed to delete role")).
					Set("event", "shutdown_role_delete"))
		}
	}()

	if err := reissueRoles(ds, db, guildID, config.BOT_ROLE_ID); err != nil {
		shutdownLogRec.Wrap(err)
	}

	go func() {
		err := db.WatchExpirations(ctx, ds)
		shutdownLogRec.Wrap(err)
		masterHandler.Cancel()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
halt:
	for {
		select {
		case <-interrupt:
			shutdownLogRec.Wrap(errlist.New(nil).Set("event", "execution stopped by user"))
			masterHandler.Cancel()
			break halt
		case <-masterHandler.Ctx.Done():
			shutdownLogRec.Wrap(errlist.New(nil).Set("event", ctx.Err().Error()))
			break halt
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
