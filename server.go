package main

import (
	"fmt"
	"main/lib"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config, commands, closer := lib.New()
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands[i.Data.Name]; ok {
			h.Handler[i.Data.Options[0].Name](s, i)
		}
	})
	if err := discord.Open(); err != nil {
		panic(err)
	}

	registeredCommandList := make([]*discordgo.ApplicationCommand, 0)
	for _, cmd := range commands {
		registeredCmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.Server, &cmd.Registry)
		if err != nil {
			panic(err)
		}
		registeredCommandList = append(registeredCommandList, registeredCmd)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("Shutting down ...")

	for _, cmd := range commands {
		if cmd.Closer != nil {
			cmd.Closer()
		}
	}

	for _, cmd := range registeredCommandList {
		err := discord.ApplicationCommandDelete(cmd.ApplicationID, config.Server, cmd.ID)
		if err != nil {
			panic(err)
		}
	}

	closer()
	discord.Close()
}
