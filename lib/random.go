package lib

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) Random() command {
	rand.Seed(time.Now().UnixNano())

	var specialAlphanumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var alphanumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var numeric = []rune("1234567890")

	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "random",
			Description: "random utility",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "number",
					Description: "get a random number between 'min' and 'max'",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "min",
							Description: "how much minimum value can be?",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "max",
							Description: "how much maximum value can be?",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "alphabet",
					Description: "generate an alphabet string",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "length",
							Description: "how long is the string",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "type",
							Description: "what type of generation?",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Uppercase Only",
									Value: "uppercase",
								},
								{
									Name:  "Lowercase Only",
									Value: "lowercase",
								},
							},
						},
					},
				},
				{
					Name:        "numeric",
					Description: "generate an numeric string",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "length",
							Description: "how long is the string",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "alphanumeric",
					Description: "generate an alphanumeric string",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "length",
							Description: "how long is the string",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "type",
							Description: "what type of generation?",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Uppercase Only",
									Value: "uppercase",
								},
								{
									Name:  "Lowercase Only",
									Value: "lowercase",
								},
							},
						},
					},
				},
				{
					Name:        "special-alphanumeric",
					Description: "generate an special alphanumeric string",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "length",
							Description: "how long is the string",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "type",
							Description: "what type of generation?",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Uppercase Only",
									Value: "uppercase",
								},
								{
									Name:  "Lowercase Only",
									Value: "lowercase",
								},
							},
						},
					},
				},
			},
		},
		Handler: map[string]commandHandler{
			"number": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				min := i.Data.Options[0].Options[0].IntValue()
				max := i.Data.Options[0].Options[1].IntValue()
				if min >= max {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionApplicationCommandResponseData{
							Content: "Error: 'min' value is bigger than 'max' value'",
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf("%v", rand.Int63n(max-min)+min),
					},
				})
			},
			"alphabet": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				n := i.Data.Options[0].Options[0].IntValue()

				b := make([]rune, n)
				for i := range b {
					b[i] = alpha[rand.Intn(len(alpha))]
				}

				val := string(b)
				if len(i.Data.Options[0].Options) >= 2 {
					typ := i.Data.Options[0].Options[1].StringValue()
					switch typ {

					case "uppercase":
						val = strings.ToUpper(val)

					case "lowercase":
						val = strings.ToLower(val)

					}
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf("%v", val),
					},
				})
			},
			"numeric": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				n := i.Data.Options[0].Options[0].IntValue()

				b := make([]rune, n)
				for i := range b {
					b[i] = numeric[rand.Intn(len(numeric))]
				}

				val := string(b)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf("%v", val),
					},
				})
			},
			"alphanumeric": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				n := i.Data.Options[0].Options[0].IntValue()

				b := make([]rune, n)
				for i := range b {
					b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
				}

				val := string(b)
				if len(i.Data.Options[0].Options) >= 2 {
					typ := i.Data.Options[0].Options[1].StringValue()
					switch typ {

					case "uppercase":
						val = strings.ToUpper(val)

					case "lowercase":
						val = strings.ToLower(val)

					}
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf("%v", val),
					},
				})
			},
			"special-alphanumeric": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				n := i.Data.Options[0].Options[0].IntValue()

				b := make([]rune, n)
				for i := range b {
					b[i] = specialAlphanumeric[rand.Intn(len(specialAlphanumeric))]
				}

				val := string(b)
				if len(i.Data.Options[0].Options) >= 2 {
					typ := i.Data.Options[0].Options[1].StringValue()
					switch typ {

					case "uppercase":
						val = strings.ToUpper(val)

					case "lowercase":
						val = strings.ToLower(val)

					}
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf("%v", val),
					},
				})
			},
		},
	}
}
