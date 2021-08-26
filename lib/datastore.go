package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/bwmarrin/discordgo"
)

func (cfg *config) Datastore() command {

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

	keyPat := regexp.MustCompile(`^/([^\.]*),(\d+)$`)
	parseKey := func(s string) (*datastore.Key, error) {
		m := keyPat.FindStringSubmatch(s)
		i := strings.Index(s, ",")
		if i < 0 {
			return nil, errors.New("bad format")
		}
		n, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			return nil, err
		}
		return datastore.IDKey(m[1], n, nil), nil
	}

	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "datastore",
			Description: "datastore utility",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "encode",
					Description: "get encoded 'value' with datastore key encode.",
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
					Description: "get decoded 'value' with datastore key decode.",
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

				key, err := parseKey(value)
				if err != nil {
					response = "Error: " + err.Error()
				} else {
					response = key.Encode()
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{generateEmbed("DatastoreKey Encode", string(value), string(response))},
					},
				})
			},
			"decode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				response := ""
				value := i.ApplicationCommandData().Options[0].Options[0].StringValue()
				key, err := datastore.DecodeKey(value)
				if err != nil {
					response = "Error: " + err.Error()
				} else {
					response = key.String()
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{generateEmbed("DatastoreKey Decode", string(value), string(response))},
					},
				})
			},
		},
	}
}
