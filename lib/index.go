package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/extension"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate)

type command struct {
	Registry discordgo.ApplicationCommand
	Handler  map[string]commandHandler
	Closer   func()
}

type config struct {
	Domain   string                 `json:"domain"`
	Services map[string]bool        `json:"services"`
	State    map[string]interface{} `json:"-"`
	Data     *data                  `json:"data"`
	Discord  struct {
		BotToken string `json:"botToken"`
		ServerID string `json:"serverId"`
	} `json:"discord"`
	Ngrok struct {
		Type  string   `json:"type"`
		Port  string   `json:"port"`
		Token string   `json:"authtoken"`
		Args  []string `json:"args"`
	} `json:"ngrok"`
	Shortlink struct {
		Authenticate bool   `json:"auth"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	} `json:"shortlink"`
	FileServer struct {
		Authenticate bool   `json:"auth"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	} `json:"fileServer"`
}

func (cfg config) hasServiceEnabled(service string) bool {
	value, ok := cfg.Services[service]
	if ok {
		return value
	}
	return false
}

type data struct {
	Shortcuts map[string]string `json:"shortcuts"`
	Files     map[string]string `json:"files"`
}

func New() (*config, map[string]command, func()) {
	file, err := os.OpenFile("config.json", os.O_RDWR, 0644)
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

	dataFile, err := os.OpenFile("data.json", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	dataBytes, err := ioutil.ReadAll(dataFile)
	if err != nil {
		panic(err)
	}

	data := new(data)
	if string(dataBytes) != "" {
		if err := json.Unmarshal(dataBytes, data); err != nil {
			panic(err)
		}
	}

	cfg.Data = data
	if cfg.Data.Shortcuts == nil {
		cfg.Data.Shortcuts = make(map[string]string)
	}

	if cfg.Data.Files == nil {
		cfg.Data.Files = make(map[string]string)
	}

	cfg.State = make(map[string]interface{})

	template := template.Must(template.ParseGlob("html/*.html"))
	handler := new(extension.RegexpHandler)

	if cfg.hasServiceEnabled("ngrok") {
		if cfg.hasServiceEnabled("shortcut") {
			exp, _ := regexp.Compile("/shortcut/\\w+")
			handler.HandleFunc(exp, func(w http.ResponseWriter, r *http.Request) {
				if cfg.Shortlink.Authenticate {
					responseCode, invalidPassword := extension.BasicAuth(cfg.Shortlink.Username, cfg.Shortlink.Password, w, r)
					if invalidPassword {
						template.ExecuteTemplate(w, "status.html", http.StatusUnauthorized)
						return
					}

					if responseCode != 0 {
						w.WriteHeader(responseCode)
						return
					}
				}

				path := strings.TrimLeft(r.URL.Path, "/shortcut/")
				if url, ok := cfg.Data.Shortcuts[path]; ok {
					http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				} else {
					template.ExecuteTemplate(w, "status.html", http.StatusNotFound)
				}
			})
		}

		if cfg.hasServiceEnabled("file") {
			exp, _ := regexp.Compile("/file/\\w+")
			handler.HandleFunc(exp, func(w http.ResponseWriter, r *http.Request) {
				if cfg.FileServer.Authenticate {
					responseCode, invalidPassword := extension.BasicAuth(cfg.FileServer.Username, cfg.FileServer.Password, w, r)
					if invalidPassword {
						template.ExecuteTemplate(w, "status.html", http.StatusUnauthorized)
						return
					}

					if responseCode != 0 {
						w.WriteHeader(responseCode)
						return
					}
				}

				path := strings.TrimLeft(r.URL.Path, "/file/")
				if file, ok := cfg.Data.Files[path]; ok {
					if strings.HasSuffix(file, ".apk") {
						w.Header().Add("Content-Type", "application/vnd.android.package-archive")
					}
					http.ServeFile(w, r, file)
				} else {
					template.ExecuteTemplate(w, "status.html", http.StatusNotFound)
				}
			})
		}

		go http.ListenAndServe(":"+cfg.Ngrok.Port, handler)
	}

	return cfg, getCommands(cfg), func() {
		defer file.Close()
		defer dataFile.Close()

		bytes, err := json.Marshal(cfg.Data)
		if err != nil {
			panic(err)
		}
		dataFile.Truncate(0)
		dataFile.Seek(0, 0)
		dataFile.Write(bytes)
	}
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
