package main

import (
	"fmt"
	"kubinka/config"
	"kubinka/handlers"
	"os"
	"os/signal"
	"syscall"

	"discordgo"
)

func main() {

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	discord.SyncEvents = true

	discord.AddHandler(handlers.Message)

	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	discord.ChannelMessageSend(config.ChanID, "Privet, mudila.")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
}
