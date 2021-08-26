package lib

import (
	"errors"
	"fmt"
	"main/lib/music"

	"github.com/bwmarrin/discordgo"
)

var player *music.Player = music.NewPlayer()

var (
	errNotInVoiceChannel = errors.New("sorry, you're not in any voice channel now")
)

func (cfg *config) Music() command {
	var replyError = func(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       15158332,
						Title:       "Something wrong when invoke action",
						Description: err.Error(),
					},
				},
			},
		})
	}

	var replySuccess = func(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       5763719,
						Description: message,
					},
				},
			},
		})
	}

	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "music",
			Description: "music youtube player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "now",
					Description: "view current song",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "skip",
					Description: "skip current song",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "stop",
					Description: "stop current song",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "play",
					Description: "play current song",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "shuffle",
					Description: "shuffle queue songs",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "repeat",
					Description: "repeat current song",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "list",
					Description: "list all queue songs",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
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
							Required: true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "url",
							Description: "youtube playlist or video url",
							Required:    true,
						},
					},
				},
			},
		},
		Handler: map[string]commandHandler{
			"now": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				current := player.Current()
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Thumbnail: &discordgo.MessageEmbedThumbnail{
									URL:    current.Display.URL,
									Width:  current.Display.Width,
									Height: current.Display.Height,
								},
								Author: &discordgo.MessageEmbedAuthor{
									Name: current.Author,
								},
								Color:       5763719,
								Title:       current.Title,
								Description: current.Description,
								Timestamp:   current.Duration.String(),
							},
						},
					},
				})
			},
			"shuffle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				player.Shuffle()
				replySuccess(s, i, "Current song have been shuffled")
			},
			"repeat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				repeat := player.Repeat()
				if repeat {
					replySuccess(s, i, "Repeat Mode: ON")
				} else {
					replySuccess(s, i, "Repeat Mode: OFF")
				}
			},
			"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				player.Skip()
				replySuccess(s, i, "Current song have been skipped")
			},
			"skip": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				player.Skip()
				replySuccess(s, i, "Current song have been skipped")
			},
			"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				player.Stop()
				replySuccess(s, i, "I'm not talking anymore :(")
			},
			"list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				current, queue := player.ListQueueAndCurrent()
				if len(queue) >= 10 {
					queue = queue[:10]
				}

				lists := "Current: " + current.Title
				for _, q := range queue {
					lists += "\n" + "1." + q.Title
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: lists,
					},
				})
			},
			"clear": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				player.ClearQueue()
				replySuccess(s, i, "Queued has been cleared")
			},
			"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if !player.CanContinue(s, i) {
					replyError(s, i, errNotInVoiceChannel)
					return
				}

				addType := i.ApplicationCommandData().Options[0].Options[0].StringValue()
				addLink := i.ApplicationCommandData().Options[0].Options[1].StringValue()
				switch addType {

				case "song":
					video, err := player.AddSong(addLink)
					if err != nil {
						replyError(s, i, err)
						return
					}

					replySuccess(s, i, "Queued "+video.Title)

				case "playlist":
					count, err := player.AddPlayList(addLink)
					if err != nil {
						replyError(s, i, err)
						return
					}

					replySuccess(s, i, "Queued "+fmt.Sprintf("%d", count)+" tracks")
				}
			},
		},
	}
}
