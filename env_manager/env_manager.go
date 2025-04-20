package env_manager

import (
	"Shipyard/database"
	"Shipyard/local_environment"
)

type EnvManagerStruct struct {
	Env map[string]EnvInterface
}

var EnvManager *EnvManagerStruct = NewEnvManager()

func NewEnvManager() *EnvManagerStruct {
	envManager := &EnvManagerStruct{
		Env: make(map[string]EnvInterface),
	}

	environments := database.LoadEnvironments()
	for _, env := range environments {
		if env.EnvType == "local" {
			envManager.Env[env.Name] = local_environment.LocalEnv
			
			envManager.Env[env.Name].SetEnvType(env.EnvType)
			envManager.Env[env.Name].SetName(env.Name)
		}
		// TODO: Handle remotes as well
	}

	return envManager
}
