package lib

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) Shortcut() command {

	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "shortcut",
			Description: "shortcut redirect utility",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "redirect",
					Description: "redirect with shorter link",
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
			"redirect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				link := i.ApplicationCommandData().Options[0].Options[0].StringValue()
				fixed := time.Now().Format("06010215040500")
				if len(i.ApplicationCommandData().Options[0].Options) == 2 {
					fixed = i.ApplicationCommandData().Options[0].Options[1].StringValue()
				}

				cfg.Data.Shortcuts[fixed] = link
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title: "Shortcut Link",
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:   "Link",
										Value:  link,
										Inline: true,
									},
									{
										Name:   "Shorter Link",
										Value:  cfg.Domain + "/shortcut/" + fixed,
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
