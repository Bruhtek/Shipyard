package intervals

import (
	"Shipyard/internal/api/websocket"
	"Shipyard/internal/env_manager"
	"github.com/rs/zerolog/log"
	"time"
)

func SetupScanning() {
	scanEnvs() // initial scan should be done immediately

	go func() {
		interval := 5 * time.Second
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		slowdownWhenIdle := 60              // times the interval
		slowdownCounter := slowdownWhenIdle // initially, scan it instantly even without any connections

		for {
			select {
			case <-ticker.C:
				connectionCount := websocket.ConnectionManager.ConnectionCount()

				if connectionCount == 0 {
					if slowdownCounter < slowdownWhenIdle {
						slowdownCounter++
						continue
					} else {
						slowdownCounter = 0
					}
				}

				scanEnvs()
			}
		}
	}()
}

func scanEnvs() {
	envs := env_manager.EnvManager.GetEnvs()
	for _, env := range envs {
		log.Debug().
			Str("env", env.GetName()).
			Msg("Scanning environment data")
		env.ScanContainers()
		env.ScanImages()
		env.ScanNetworks()
	}
}
