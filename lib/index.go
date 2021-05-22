package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate)

type command struct {
	Registry discordgo.ApplicationCommand
	Handler  map[string]commandHandler
	Closer   func()
}

type config struct {
	Token  string `json:"token"`
	Server string `json:"server"`
	Ngrok  struct {
		Type   string `json:"type"`
		Port   string `json:"port"`
		Region string `json:"region"`
		Token  string `json:"authtoken"`
	} `json:"ngrok"`
}

func New() (*config, map[string]command, func()) {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	cfg := new(config)
	if err := json.Unmarshal(fileBytes, cfg); err != nil {
		panic(err)
	}

	return cfg, getCommands(cfg), func() {
		defer file.Close()

		bytes, err := json.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		file.Write(bytes)
	}
}

func getCommands(cfg *config) map[string]command {
	commandMap := make(map[string]command, 0)
	rValue := reflect.ValueOf(cfg)
	for i := 0; i < rValue.NumMethod(); i++ {
		method := rValue.Method(i)
		rResponse := method.Call(nil)[0].Interface().(command)
		commandMap[rResponse.Registry.Name] = rResponse
	}
	return commandMap
}
