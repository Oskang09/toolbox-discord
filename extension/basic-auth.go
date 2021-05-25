package extension

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuth(username string, password string, w http.ResponseWriter, r *http.Request) (int, bool) {

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return 401, false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return 401, false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return 401, false
	}

	if pair[0] != username || pair[1] != password {
		return 401, true
	}

	return 0, false
}
