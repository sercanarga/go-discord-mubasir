package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mubasirNew/db"
	"mubasirNew/modules"
	"os"
	"os/signal"
	"syscall"
)

var Buffer = make([][]byte, 0)

func init() {
	if *db.BotToken == "" {
		fmt.Println("No token provided.")
		return
	}
}

func main() {
	dg, err := discordgo.New("Bot " + *db.BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(modules.UpdateStatus)
	dg.AddHandler(modules.CallBot)
	dg.AddHandler(modules.UserJoinEvent)
	dg.AddHandler(modules.BotChangeChannelEvent)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("mubasir is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
