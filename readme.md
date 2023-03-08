![image](https://user-images.githubusercontent.com/15674107/119238117-3503b280-bb73-11eb-9e58-bceca156d728.png)
![image](https://user-images.githubusercontent.com/15674107/131013217-fd7a3664-47df-418c-b6ae-44f2bcd79e72.png)

# toolbox

personal discord toolbox bot mainly for speed up problem fixing by adding some utility that always have to use when development & solving issues.

# Features

1. Generate Public & Private Key with RSA - 1024, 2048 and 4096
2. Base64 Encode, Decode
3. Datastore Key Encode, Decode
4. URL Parsing ( to have a better view )
5. Random characters generator ( alphabet, numeric, symbol )
6. Keygen UUID and TUID ( timestmap based uid )
7. Sqlike `types.Key` Encode, Decode

# Make your own

1. Create a `config.json` with structure below.

```go
type config struct {
	Discord struct {
		BotToken string `json:"botToken"` // Discord bot token
		ServerID string `json:"serverId"` // Your personal server id
	} `json:"discord"`
	Services    map[string]bool        `json:"services"`
}
```

2. Create `data.json` and leave it empty.
3. Start your bot by `go run .` or build binary with `go build`.


# Services

| Service   | Description                                           | Config Key   | Dependecies & CLI |
| --------- | ----------------------------------------------------- | ------------ | ----------------- |
| base64    | base64 encode, decode                                 |              |                   |
| datastore | datastore encode, decode                              |              |                   |
| file      | file serving                                          | fileServer   | ngrok,ngrok.exe   |
| keygen    | generate public, private key with RSA1024, 2048, 4096 |              |                   |
| music     | music player bot                                      |              | ffmpeg.exe        |
| ngrok     | hosted tunnel public ip                               | ngrok,domain | ngrok.exe         |
| random    | random characters generator                           |              |                   |
| shortcut  | short link creator                                    | shortlink    | ngrok,ngrok.exe   |
| url       | url parsing reader                                    |              |                   |
| sqlike    | sqlike encode, decode                                 |              |                   |

# Example Configuration

```json
{
    "services": {
        "base64": true,
        "datastore": true,
        "file": true,
        "keygen": true,
        "random": true,
        "url": true,
        "sqlike": true
    },
    "discord": {
        "botToken": "", // bot token
        "serverId": "" // server id
    }
}
```

# Extra: Startup Application

Start bot when computer startup, for Windows 10 users you can use "Windows + R" and type "shell:startup". After folder popup, just put built binrary shortcut inside. Since some services required `cli`, aslo `config.json` and `data.json`.


# Docker Build

```
$ docker buildx build . -t IMAGE_NAME:VERSION --platform=linux/amd64
```
