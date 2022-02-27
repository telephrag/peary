package main

import (
	"context"
	"fmt"
	"kubinka/changestream"
	"kubinka/commands"
	"kubinka/config"
	"kubinka/db"
	"kubinka/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"discordgo"
)

func getLogFile(fileName string) *os.File {
	// setting up logging, for some reason it loggin wont work properly
	// if it was setup inside init()
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Failed to open file for logging.\n\n\n")
	}
	return f
}

func getDiscordSession(token string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + token) // TODO: Make to arg
	if err != nil {
		log.Fatal("Could not create session.\n\n\n")
	}
	discord.SyncEvents = false
	discord.AddHandler(handlers.Select)

	err = discord.Open()
	if err != nil {
		log.Fatal("Could not open connection.\n\n\n")
	}

	return discord
}

func deleteCommands(ds *discordgo.Session) {
	for _, cmd := range commands.Commands {
		err := ds.ApplicationCommandDelete(
			ds.State.User.ID,
			config.GuildID,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete %q command: %v\n\n\n", cmd.Name, err)
		}
	}
}

func createCommands(ds *discordgo.Session, appId string, guildId string) {
	var err error
	for i, cmd := range commands.Commands {
		commands.Commands[i], err = ds.ApplicationCommandCreate(
			appId,
			guildId,
			cmd,
		)
		if err != nil {
			if i > 0 {
				deleteCommands(ds) // deferred code will not run after fatal or panic
			}
			log.Fatalf("Failed to create command %s:\n %s\n\n\n", cmd.Name, err)
		}
	}
}

func main() {
	log.SetOutput(getLogFile(config.LogFileName))
	log.Print("<<<<< SESSION STARTUP >>>>>\n")
	defer log.Print("<<<<< SESSION SHUTDOWN >>>>>\n\n\n")

	discord := getDiscordSession(config.Token)
	defer discord.Close()
	defer deleteCommands(discord) // Removing commands on bot shutdown
	createCommands(discord, config.AppID, config.GuildID)

	ctx, cancel := context.WithCancel(context.Background())
	go changestream.WatchEvents(db.Instance.Collection, ctx, cancel) // Asynchronously watch events

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	// handling invalidation of collection at shutdown
	timeout := time.After(3 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("handled invalidation at shutdown")
			return
		case <-timeout:
			fmt.Println("invalidation wasn't handled or didn't occur")
			return
		default:
		}
		time.Sleep(time.Millisecond * 500)
	}
}
