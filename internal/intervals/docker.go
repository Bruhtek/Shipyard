package intervals

import (
	"Shipyard/internal/actions"
	"Shipyard/internal/environments"
	"Shipyard/internal/environments/types"
	"Shipyard/internal/remote_worker"
	"time"

	"github.com/rs/zerolog/log"
)

func SetupScanning(isRemote bool) {
	go func() {
		scanEnvs() // initial scan should be done immediately, but not block

		interval := 5 * time.Second
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		slowdownWhenIdle := 60              // times the interval
		slowdownCounter := slowdownWhenIdle // initially, scan it instantly even without any connections

		for {
			select {
			case <-ticker.C:
				connectionCount := 0
				if !isRemote {
					connectionCount = actions.ConnectionManager.ConnectionCount()
				} else {
					if remote_worker.CManager.IsConnected() {
						connectionCount = 1
					}
				}

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
	envs := environments.EnvManager.GetEnvs()
	for _, envI := range envs {
		env, ok := envI.(types.LocalEnvironment)
		if !ok {
			continue
		}

		log.Debug().
			Str("env", env.GetName()).
			Msg("Scanning environment data")
		env.ScanImages()
		env.ScanContainers()
		env.ScanNetworks()
		env.ScanStacks()
	}
}
