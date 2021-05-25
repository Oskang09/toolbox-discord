![image](https://user-images.githubusercontent.com/15674107/119238117-3503b280-bb73-11eb-9e58-bceca156d728.png)

# toolbox

personal discord toolbox bot mainly for speed up problem fixing by adding some utility that always have to use when development & solving issues.

# Features

1. Generate Public & Private Key with RSA - 1024, 2048 and 4096
2. Base64 Encode, Decode
3. Datastore Key Encode, Decode
4. URL Parsing ( to have a better view )
5. Ngrok IP ( get hosted public ip )
6. Random characters generator ( alphabet, numeric, symbol )
7. Keygen UUID and TUID ( timestmap based uid )
8. Shortlink Creator
9. File Serving

# Make your own

1. Create a `config.json` with structure below.

```go
type config struct {
	Discord struct {
		BotToken string `json:"botToken"` // Discord bot token
		ServerID string `json:"serverId"` // Your personal server id
	} `json:"discord"`
	Ngrok struct {
		Type  string   `json:"type"` // Start port in? http, tcp
		Port  string   `json:"port"` // Port numbr
		Token string   `json:"authtoken"`// AuthToken
		Args  []string `json:"args"` // Extra arguments for setup -auth, -region, 
	} `json:"ngrok"`
	Shortlink struct {
		Authenticate bool `json:"auth"` // use authenticate
		Username  string `json:"username"` // Auth username
		Password  string `json:"password"` // Auth password
	} `json:"shortlink"`
	FileServer struct {
		Authenticate bool `json:"auth"` // use authenticate
		Username  string `json:"username"` // Auth username
		Password  string `json:"password"` // Auth password
	} `json:"fileServer"`
}
```

2. Create `data.json` and leave it empty.
3. Start your bot by `go run .` or build binary with `go build`.

# Extra: Startup Application

Start bot when computer startup, for Windows 10 users you can use "Windows + R" and type "shell:startup". After folder popup, just put built binrary shortcut inside. If you want move binary directly inside, you must move `lib` folder too, since it's depend on the `ngrok.exe` cli aslo `config.json` and `data.json`.