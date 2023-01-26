package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"peary/config"
	"peary/errconst"
	"peary/service"
	"peary/strg"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/telephrag/errlist"
)

func getLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			log.Fatalln(errlist.New(err))
		}
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
		log.Fatal(errlist.New(errors.New("could not create session")))
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

func createCommands(ds *discordgo.Session, appID string, guildID string) {
	var err error
	for i, cmd := range service.CmdDef {
		service.CmdDef[i], err = ds.ApplicationCommandCreate(
			appID,
			guildID,
			cmd,
		)
		if err != nil {
			if i > 0 {
				deleteCommands(ds, guildID)
			}
			log.Fatal(errlist.New(fmt.Errorf("failed to create command /%s", cmd.Name)).Wrap(err))
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
			log.Fatal(errlist.New(fmt.Errorf("could not delete %s command", cmd.Name)).Wrap(err))
		}
	}
}

// ds.GuildRoleDelete(guildId, roleId) to rollback
func createRole(ds *discordgo.Session, name, guildID string, color int) (roleID string, err error) {
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
			return errlist.New(err).Set("event", errconst.StartupRoleReissue).Set("session", id)
		}
		log.Print(errlist.New(nil).Set("session", id).Set("event", errconst.StartupRoleReissue))
	}

	return nil
}

func main() {
	logFile := getLogFile("/data/" + config.NAME + ".log")
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)
	log.Print(errlist.New(nil).Set("event", "SESSION STARTUP"))
	shutdownLogRec := errlist.New(nil).Set("event", "SESSION SHUTDOWN")
	defer log.Print(shutdownLogRec)

	token, ok := syscall.Getenv("BOT_TOKEN")
	if !ok {
		shutdownLogRec.Wrap(fmt.Errorf("no BOT_TOKEN provided"))
		return
	}
	appID, ok := syscall.Getenv("BOT_APP_ID")
	if !ok {
		shutdownLogRec.Wrap(fmt.Errorf("no BOT_APP_ID provided"))
		return
	}

	db, err := strg.Connect("/data/"+config.NAME, config.DB_PLAYERS_BUCKET_NAME)
	if err != nil {
		shutdownLogRec.Wrap(err).Set("event", "startup_db_connect")
	} else {
		defer db.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ds, masterHandler := newDiscordSession(ctx, cancel, token, db)
	defer masterHandler.Cancel()
	defer masterHandler.HaltUntilAllDone()
	defer ds.Close()

	guildID := ds.State.Ready.Guilds[0].ID

	createCommands(ds, appID, guildID)
	defer deleteCommands(ds, guildID) // Removing commands on bot shutdown.

	// Each time we create new role to save user a hustle getting ID of precreated one
	// using developer mode.
	config.BOT_ROLE_ID, err = createRole(ds, config.BOT_ROLE_NAME, guildID, config.BOT_ROLE_COLOR)
	if err != nil {
		shutdownLogRec.Wrap(
			errlist.New(fmt.Errorf("failed to create role: %w", err)).
				Set("event", "startup_role_create"))
		return
	}
	defer func() { // Said role is deleted on shutdown.
		if err := ds.GuildRoleDelete(guildID, config.BOT_ROLE_ID); err != nil {
			shutdownLogRec.Wrap(
				errlist.New(fmt.Errorf("failed to delete role")).
					Set("event", "shutdown_role_delete"))
		}
	}()

	// Since we delete roles on shutdown users lose their roles which is a good thing
	// because, they won't receive pings while not being able to get rid of the role.
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
