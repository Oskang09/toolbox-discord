package lib

import (
	"main/lib/music"

	"github.com/bwmarrin/discordgo"
)

var player *music.Player

func (cfg *config) Music() command {
	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "music",
			Description: "music youtube player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "clear",
					Description: "clear current queue",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "add",
					Description: "add song / playlist current queue",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "type",
							Description: "type of youtube url",
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "playlist",
									Value: "playlist",
								},
								{
									Name:  "song",
									Value: "song",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "url",
							Description: "youtube playlist or video url",
						},
					},
				},
			},
		},
		Handler: map[string]commandHandler{
			"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if player == nil {
					player = music.NewPlayer(s, cfg.Discord.ServerID, cfg.MusicPlayer.VoiceChannel)
				}

				addType := i.Data.Options[0].Options[0].StringValue()
				addLink := i.Data.Options[0].Options[1].StringValue()

				switch addType {

				case "song":
					player.AddSong(addLink)

				case "playlist":
					player.AddPlayList(addLink)

				}
			},
		},
	}
}
