package lib

import (
	"fmt"
	"strings"

	"github.com/RevenueMonster/sqlike/types"
	"github.com/bwmarrin/discordgo"
)

func (cfg *config) Sqlike() (string, func() command) {
	return "sqlike", func() command {

		generateEmbed := func(title string, before string, after string) *discordgo.MessageEmbed {
			contents := []string{
				"**Before**",
				fmt.Sprintf("```%v```", before),
				"",
				"**After**",
				fmt.Sprintf("```%v```", after),
			}
			return &discordgo.MessageEmbed{
				Title:       title,
				Description: strings.Join(contents, "\n"),
			}
		}

		return command{
			Registry: discordgo.ApplicationCommand{
				Name:        "sqlike",
				Description: "sqlike utility",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "encode",
						Description: "get encoded 'value' with types key encode.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "value",
								Description: "value to encode",
								Required:    true,
							},
						},
					},
					{
						Name:        "decode",
						Description: "get decoded 'value' with types key decode.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "value",
								Description: "value to decode",
								Required:    true,
							},
						},
					},
				},
			},
			Handler: map[string]commandHandler{
				"encode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					response := ""
					value := i.ApplicationCommandData().Options[0].Options[0].StringValue()
					if !strings.HasPrefix(value, "/") {
						value = "/" + value
					}

					key, err := types.ParseKey(value)
					if err != nil {
						response = "Error: " + err.Error()
					} else {
						response = key.Encode()
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{generateEmbed("TypesKey Encode", string(value), string(response))},
						},
					})
				},
				"decode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					response := ""
					value := i.ApplicationCommandData().Options[0].Options[0].StringValue()

					key, err := types.DecodeKey(value)
					if err != nil {
						response = "Error: " + err.Error()
					} else {
						response = key.String()
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{generateEmbed("TypesKey Decode", string(value), string(response))},
						},
					})
				},
			},
		}
	}
}
