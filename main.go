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
	defer log.Print("\n\n\n")

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

	for i, cmd := range commands.Commands {
		commands.Commands[i], err = discord.ApplicationCommandCreate(
			config.AppID,
			config.GuildID,
			cmd,
		)
		if err != nil {
			log.Panicln(err, " while creating command: ", cmd.Name)
		}
	}

	/*
		go func() {
			check if deployment time expired
			if yes{
				remove player from db
			}
			also need to add players to DB inside Deploy()
			and remove them inside Return()
		}
	*/

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	for _, cmd := range commands.Commands { // Removing commands on bot shutdown
		err := discord.ApplicationCommandDelete(
			discord.State.User.ID,
			config.GuildID,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete %q command: %v", cmd.Name, err)
		}
	}

	log.Print("<<<<< SESSION SHUTDOWN >>>>>")
}
