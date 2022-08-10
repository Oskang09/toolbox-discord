package lib

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) URL() (string, func() command) {
	return "url", func() command {
		codeBlock := func(text string) string {
			return "`" + text + "`"
		}

		return command{
			Registry: discordgo.ApplicationCommand{
				Name:        "url",
				Description: "url utility",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "parse",
						Description: "read url in better way",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "value",
								Description: "value to parse",
								Required:    true,
							},
						},
					},
				},
			},
			Handler: map[string]commandHandler{
				"parse": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					rurl := i.ApplicationCommandData().Options[0].Options[0].Value.(string)
					parsedUrl, err := url.Parse(rurl)
					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds: []*discordgo.MessageEmbed{
									{
										Title: "URL Parse",
										Fields: []*discordgo.MessageEmbedField{
											{
												Name:  "Error: ",
												Value: "Error: " + err.Error(),
											},
										},
									},
								},
							},
						})
						return
					}

					contents := make([]string, 0)
					contents = append(contents, "**Host**        | "+codeBlock(parsedUrl.Host))
					contents = append(contents, "**Path**        | "+codeBlock(parsedUrl.Path))
					contents = append(contents, "**Scheme**  | "+codeBlock(parsedUrl.Scheme))
					queryBytes, _ := json.MarshalIndent(parsedUrl.Query(), "", "  ")
					if string(queryBytes) != "{}" {
						contents = append(contents, "**Query**     |")
						contents = append(contents, "```json")
						contents = append(contents, string(queryBytes))
						contents = append(contents, "```")
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: strings.Join(contents, "\n"),
						},
					})
				},
			},
		}
	}
}
