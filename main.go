package main

import (
	"kubinka/commands"
	"kubinka/config"
	"kubinka/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discordgo"
)

func main() {

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Could not create session.")
		return
	}
	discord.SyncEvents = true

	discord.AddHandler(handlers.Select)

	err = discord.Open()
	if err != nil {
		log.Fatal("Could not open connection.")
	}

	for _, cmd := range commands.Commands {
		_, err := discord.ApplicationCommandCreate(
			config.AppID,
			config.GuildID,
			&cmd,
		)
		if err != nil {
			log.Panic(err, " while creating command: ", cmd.Name)
		}

	}

	defer discord.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
}
