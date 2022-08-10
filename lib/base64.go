package lib

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) Base64() (string, func() command) {
	return "base64", func() command {
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
				Name:        "base64",
				Description: "base64 encode & decode utility",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "encode",
						Description: "get encoded 'value' with base64 encode.",
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
						Description: "get decoded 'value' with base64 decode.",
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
					value := []byte(i.ApplicationCommandData().Options[0].Options[0].StringValue())
					encodedBytes := base64.RawStdEncoding.EncodeToString(value)

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{generateEmbed("Base64 Encode", string(value), string(encodedBytes))},
						},
					})
				},
				"decode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					response := ""
					value := i.ApplicationCommandData().Options[0].Options[0].StringValue()
					decodedBytes, err := base64.RawStdEncoding.DecodeString(value)
					if err != nil {
						response = "Error: " + err.Error()
					} else {
						response = string(decodedBytes)
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{generateEmbed("Base64 Decode", string(value), string(response))},
						},
					})
				},
			},
		}
	}
}
