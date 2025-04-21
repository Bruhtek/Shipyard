package env_manager

import (
	"Shipyard/database"
	"Shipyard/local_environment"
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

var EnvManager *EnvManagerStruct = NewEnvManager()

func NewEnvManager() *EnvManagerStruct {
	envManager := &EnvManagerStruct{
		env:   make(map[string]EnvInterface),
		mutex: sync.RWMutex{},
	}

	envManager.mutex.Lock()
	defer envManager.mutex.Unlock()

	environments := database.LoadEnvironments()
	for _, env := range environments {
		if env.EnvType == "local" {
			envManager.env[env.Name] = local_environment.LocalEnv

			envManager.env[env.Name].SetEnvType(env.EnvType)
			envManager.env[env.Name].SetName(env.Name)
		}
		// TODO: Handle remotes as well
	}

	return envManager
}
