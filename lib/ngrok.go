package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (cfg *config) NgrokCMD() (string, func() command) {
	return "ngrok", func() command {
		wd, _ := os.Getwd()

		authCmd := exec.Command("ngrok.exe", "authtoken", cfg.Ngrok.Token)
		authCmd.Dir = wd + "/cli"
		if err := authCmd.Start(); err != nil {
			panic(err)
		}

		if err := authCmd.Wait(); err != nil {
			panic(err)
		}

		commands := make([]string, 0)
		commands = append(commands, cfg.Ngrok.Type)
		commands = append(commands, cfg.Ngrok.Port)
		commands = append(commands, cfg.Ngrok.Args...)

		cmd := exec.Command("ngrok.exe", commands...)
		cmd.Dir = wd + "/cli"
		if err := cmd.Start(); err != nil {
			panic(err)
		}

		req, err := http.NewRequest("GET", "http://127.0.0.1:4040/api/tunnels/command_line", nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", "Bearer "+cfg.Ngrok.Token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var httpResponse struct {
			PublicURL string `json:"public_url"`
		}

		if err := json.Unmarshal(bytes, &httpResponse); err != nil {
			panic(err)
		}

		// fallback for those don't have domain will use ngrok public ip
		if cfg.Domain == "" {
			cfg.Domain = strings.TrimRight(httpResponse.PublicURL, "/")
		}

		return command{
			Registry: discordgo.ApplicationCommand{
				Name:        "ngrok",
				Description: "ngrok utility",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "ip",
						Description: "get current ngrok public ip",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
			Handler: map[string]commandHandler{
				"ip": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "Ngrok",
									Fields: []*discordgo.MessageEmbedField{
										{
											Name:   "Domain IP",
											Value:  cfg.Domain,
											Inline: true,
										},
										{
											Name:   "Public IP",
											Value:  httpResponse.PublicURL,
											Inline: true,
										},
										{
											Name:   "WebUI IP",
											Value:  "http://127.0.0.1:4040",
											Inline: true,
										},
									},
								},
							},
						},
					})
				},
			},
			Closer: func() {
				cmd.Process.Kill()
				cmd.Process.Wait()
			},
		}
	}
}
