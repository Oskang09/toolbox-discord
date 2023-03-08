package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate)

type command struct {
	Registry discordgo.ApplicationCommand
	Handler  map[string]commandHandler
	Closer   func()
}

type config struct {
	Services map[string]bool        `json:"services"`
	State    map[string]interface{} `json:"-"`
	Discord  struct {
		BotToken string `json:"botToken"`
		ServerID string `json:"serverId"`
	} `json:"discord"`
}

func (cfg config) hasServiceEnabled(service string) bool {
	value, ok := cfg.Services[service]
	if ok {
		return value
	}
	return false
}

func New() (*config, map[string]command, func()) {
	var bytes []byte
	base64Config := os.Getenv("TOOLBOX_CONFIG")
	if base64Config != "" {
		configBytes, err := base64.RawStdEncoding.DecodeString(base64Config)
		if err != nil {
			panic(err)
		}
		bytes = configBytes
	} else {
		file, err := os.OpenFile("config.json", os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		bytes, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	cfg := new(config)
	if err := json.Unmarshal(bytes, cfg); err != nil {
		panic(err)
	}

	cfg.State = make(map[string]interface{})
	return cfg, getCommands(cfg), func() {}
}

func getCommands(cfg *config) map[string]command {
	commandMap := make(map[string]command)
	rValue := reflect.ValueOf(cfg)

	fmt.Println("Registered Services :-")
	fmt.Println("")
	count := 0
	for i := 0; i < rValue.NumMethod(); i++ {
		method := rValue.Method(i)
		rResponse := method.Call(nil)
		serviceString := rResponse[0].Interface().(string)
		if cfg.hasServiceEnabled(serviceString) {
			fmt.Println("[" + strconv.Itoa(count) + "] " + serviceString)

			initializer := rResponse[1].Interface().(func() command)
			cmd := initializer()
			count += 1
			commandMap[cmd.Registry.Name] = cmd
		}
	}
	return commandMap
}
