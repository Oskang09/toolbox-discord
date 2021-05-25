package main

import (
	"fmt"
	"main/lib"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const banner = `
████████╗ ██████╗  ██████╗ ██╗     ██████╗  ██████╗ ██╗  ██╗
╚══██╔══╝██╔═══██╗██╔═══██╗██║     ██╔══██╗██╔═══██╗╚██╗██╔╝
   ██║   ██║   ██║██║   ██║██║     ██████╔╝██║   ██║ ╚███╔╝ 
   ██║   ██║   ██║██║   ██║██║     ██╔══██╗██║   ██║ ██╔██╗ 
   ██║   ╚██████╔╝╚██████╔╝███████╗██████╔╝╚██████╔╝██╔╝ ██╗
   ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝
`

func main() {

	fmt.Println(banner)

	config, commands, closer := lib.New()
	discord, err := discordgo.New("Bot " + config.Discord.BotToken)
	if err != nil {
		panic(err)
	}
	defer closer()
	defer discord.Close()

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateGameStatus(0, "with Tools")
	})

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
		registeredCmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.Discord.ServerID, &cmd.Registry)
		if err != nil {
			panic(err)
		}
		registeredCommandList = append(registeredCommandList, registeredCmd)
	}

	fmt.Println("")
	fmt.Println(">> Service is now ready! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("Shutting down services ...")

	for _, cmd := range commands {
		if cmd.Closer != nil {
			cmd.Closer()
		}
	}

	for _, cmd := range registeredCommandList {
		err := discord.ApplicationCommandDelete(cmd.ApplicationID, config.Discord.ServerID, cmd.ID)
		if err != nil {
			panic(err)
		}
	}
}
