package env_manager

import (
	"Shipyard/database"
	"Shipyard/internal/local_environment"
	"Shipyard/internal/remote_controller"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

type EnvManagerStruct struct {
	env   map[string]EnvInterface
	mutex sync.RWMutex
}

func (e *EnvManagerStruct) GetEnvs() map[string]EnvInterface {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.env
}
func (e *EnvManagerStruct) GetEnv(name string) EnvInterface {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	if env, ok := e.env[name]; ok {
		return env
	}
	return nil
}
func (e *EnvManagerStruct) GetRemoteEnv(key string) EnvInterface {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	for _, env := range e.env {
		if env.GetEnvType() == "remote" {
			remote, ok := env.(*remote_controller.RemoteEnvironment)
			if ok && remote.Key == key {
				return remote
			}
		}
	}

	return nil
}

var EnvManager *EnvManagerStruct

// InitializeEnvManager
// basicOnly bool - controls whether the program should read the database,
// or run in a simple, local-only mode (used for docker_external mode)
func InitializeEnvManager(basicOnly bool) {
	EnvManager = newEnvManager(basicOnly)
}

func newEnvManager(basicOnly bool) *EnvManagerStruct {
	envManager := &EnvManagerStruct{
		env:   make(map[string]EnvInterface),
		mutex: sync.RWMutex{},
	}

	envManager.mutex.Lock()
	defer envManager.mutex.Unlock()

	if basicOnly {
		envManager.env["local"] = local_environment.NewLocalEnv()
		envManager.env["local"].SetName("local")

		return envManager
	}

	environments := database.LoadEnvironments()
	for _, env := range environments {
		switch env.EnvType {
		case "local":
			if os.Getenv("IGNORE_LOCAL") == "true" {
				log.Warn().
					Str("name", env.Name).
					Msg("IGNORE_LOCAL environment variable is set to true. Skipping creating new local environment")
				continue
			}

			envManager.env[env.Name] = local_environment.NewLocalEnv()
			envManager.env[env.Name].SetName(env.Name)
			log.Info().
				Str("name", env.Name).
				Msg("Creating new local environment")
			break
		case "remote":
			envManager.env[env.Name] = remote_controller.NewRemoteEnv(env.Key)
			envManager.env[env.Name].SetName(env.Name)
			log.Info().
				Str("name", env.Name).
				Msg("Creating new remote environment")
			break
		default:
			log.Warn().
				Str("name", env.Name).
				Str("type", env.EnvType).
				Msg("Unrecognized environment type. Skipping.")
		}
	}

	return envManager
}
