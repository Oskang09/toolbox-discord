package lib

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) File() (string, func() command) {
	return "file", func() command {
		return command{
			Registry: discordgo.ApplicationCommand{
				Name:        "file",
				Description: "file serving utility",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "serve",
						Description: "serve file with short link",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "path",
								Description: "path of the file to be serve",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "fixed",
								Description: "fixed key to be serve",
							},
						},
					},
				},
			},
			Handler: map[string]commandHandler{
				"serve": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					filePath := i.ApplicationCommandData().Options[0].Options[0].StringValue()
					fixed := time.Now().Format("06010215040500")
					if len(i.ApplicationCommandData().Options[0].Options) == 2 {
						fixed = i.ApplicationCommandData().Options[0].Options[1].StringValue()
					}

					cfg.Data.Files[fixed] = filePath
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "File Link",
									Fields: []*discordgo.MessageEmbedField{
										{
											Name:   "File",
											Value:  filePath,
											Inline: true,
										},
										{
											Name:   "Link",
											Value:  cfg.Domain + "/file/" + fixed,
											Inline: true,
										},
									},
								},
							},
						},
					})
				},
			},
		}
	}
}
