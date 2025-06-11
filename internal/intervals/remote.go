package intervals

import (
	"Shipyard/internal/remote_worker"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func SetupHeartbeat() {
	host := os.Getenv("CONTROLLER_HOST")
	key := os.Getenv("CONTROLLER_KEY")

	if host == "" || key == "" {
		log.Fatal().
			Str("CONTROLLER_HOST", host).
			Str("CONTROLLER_KEY", key).
			Msg("CONTROLLER_HOST and CONTROLLER_KEY are required")
	}

	log.Info().
		Str("CONTROLLER_HOST", host).
		Str("CONTROLLER_KEY", key).
		Msg("Setting up heartbeat")

	// initial heartbeat should be immediate
	DoHeartbeat(host, key)

	go func() {
		interval := 5 * time.Second
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				DoHeartbeat(host, key)
			}
		}
	}()
}

func DoHeartbeat(host string, key string) {
	resp, err := http.Get(host + "/api/remote/heartbeat?key=" + key)
	if err != nil {
		log.Error().
			Str("host", host).
			Str("key", key).
			Err(err).
			Msg("Error while sending heartbeat")
		return
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode == http.StatusOK {
		log.Debug().
			Msg("Successfully sent heartbeat")
		return
	}

	if statusCode == http.StatusAccepted {
		log.Debug().
			Msg("Successfully sent heartbeat")
		go remote_worker.CManager.ConnectToController(host, key)
		return
	}

	log.Error().
		Str("host", host).
		Str("key", key).
		Int("statusCode", statusCode).
		Msg("Error while sending heartbeat")
}
