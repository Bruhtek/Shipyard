package env_manager

import (
	"Shipyard/database"
	"Shipyard/internal/local_environment"
	"Shipyard/internal/remote_environment"
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
func (e *EnvManagerStruct) GetEnvByKey(key string) EnvInterface {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	for _, env := range e.env {
		if env.GetEnvDescription().EnvKey == key {
			return env
		}
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
			envManager.env[env.Name] = local_environment.NewLocalEnv()

			envManager.env[env.Name].SetEnvType(env.EnvType)
			envManager.env[env.Name].SetName(env.Name)
			continue
		}

		if env.EnvType == "remote" {
			envManager.env[env.Name] = remote_environment.NewRemoteEnvironment(env.EnvKey)
			envManager.env[env.Name].SetEnvType(env.EnvType)
			envManager.env[env.Name].SetName(env.Name)
		}
		// TODO: Handle remotes as well
	}

	return envManager
}
