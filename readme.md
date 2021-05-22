![image](https://user-images.githubusercontent.com/15674107/119238101-18677a80-bb73-11eb-9fee-483807b9d3d0.png)

# toolbox

personal discord toolbox bot mainly for speed up problem fixing by adding some utility that always have to use when development & solving issues.

# Features

1. Generate Public & Private Key with RSA - 1024, 2048 and 4096
2. Base64 Encode, Decode
3. Datastore Key Encode, Decode
4. URL Parsing ( to have a better view )
5. Ngrok IP ( get hosted public ip )

# Make your own

1. Create a `config.json` with structure below.

```go
type config struct {
	Token   string `json:"token"` // Discord bot token
	Server  string `json:"server"` // Your personal server id
	Ngrok   struct { 
		Type   string `json:"type"` // Start port in? http, tcp
		Port   string `json:"port"` // Port numbr
		Region string `json:"region"` // Region
		Token  string `json:"authtoken"` // AuthToken
	} `json:"ngrok"`
}
```

2. Start your bot by `go run .` or build binary with `go build`.
