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

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln("Failed to open file for logging.")
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("<<<<< SESSION STARTUP >>>>>")

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalln("Could not create session.")
		return
	}
	discord.SyncEvents = true
	discord.AddHandler(handlers.Select)

	err = discord.Open()
	if err != nil {
		log.Fatalln("Could not open connection.")
	}
	defer discord.Close()

	for _, cmd := range commands.Commands {
		_, err := discord.ApplicationCommandCreate(
			config.AppID,
			config.GuildID,
			&cmd,
		)
		if err != nil {
			log.Panicln(err, " while creating command: ", cmd.Name)
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	log.Print("<<<<< SESSION SHUTDOWN >>>>>\n\n\n")
}
