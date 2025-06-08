package remote_worker

import (
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

func ConnectToController(uri string, key string) {
	protocol := "ws"
	if strings.HasPrefix(uri, "https") {
		protocol = "wss"
	}

	host := uri[strings.Index(uri, "://")+3:]

	u := url.URL{
		Scheme: protocol,
		Host:   host,
		Path:   "/api/remote/" + key + "/ws",
	}

	log.Info().
		Str("host", host).
		Str("path", u.Path).
		Msg("Connecting to controller")

	return
}
