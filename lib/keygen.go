package lib

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cfg config) Keygen() command {

	generateSuccess := func(title string) *discordgo.MessageEmbed {
		return &discordgo.MessageEmbed{
			Title:       title,
			Description: "keys generated successfully as attachment below.",
		}
	}

	generateError := func(title string, err error) *discordgo.MessageEmbed {
		return &discordgo.MessageEmbed{
			Title: title,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error: ",
					Value: "Error: " + err.Error(),
				},
			},
		}
	}

	return command{
		Registry: discordgo.ApplicationCommand{
			Name:        "keygen",
			Description: "keygen utility",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "rsa",
					Description: "generate public & private key with rsa algorithm",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "size",
							Description: "how much bit size to generate key?",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "1024",
									Value: 1024,
								},
								{
									Name:  "2048",
									Value: 2048,
								},
								{
									Name:  "4096",
									Value: 4096,
								},
							},
						},
					},
				},
			},
		},
		Handler: map[string]commandHandler{
			"rsa": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				bitSize := i.Data.Options[0].Options[0].Value.(float64)
				privKey, err := rsa.GenerateKey(rand.Reader, int(bitSize))
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionApplicationCommandResponseData{
							Embeds: []*discordgo.MessageEmbed{generateError("RSA "+fmt.Sprintf("%v", bitSize)+" Keys", err)},
						},
					})
					return
				}

				pubKey := privKey.PublicKey

				privPem := pem.EncodeToMemory(
					&pem.Block{
						Type:  "RSA PRIVATE KEY",
						Bytes: x509.MarshalPKCS1PrivateKey(privKey),
					},
				)

				pubPem := pem.EncodeToMemory(
					&pem.Block{
						Type:  "RSA PUBLIC KEY",
						Bytes: x509.MarshalPKCS1PublicKey(&pubKey),
					},
				)

				pubFileBuffer := new(bytes.Buffer)
				pubFileBuffer.Write(pubPem)
				fileBuffer := new(bytes.Buffer)
				fileBuffer.Write(privPem)

				s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
					Files: []*discordgo.File{
						{
							Name:   "publickey.pub",
							Reader: pubFileBuffer,
						},
						{
							Name:   "privatekey",
							Reader: fileBuffer,
						},
					},
				})
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Embeds: []*discordgo.MessageEmbed{generateSuccess("RSA " + fmt.Sprintf("%v", bitSize) + " Keys")},
					},
				})
			},
		},
	}
}
