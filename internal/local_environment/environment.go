package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/utils"
	"sync"
)

type LocalEnvironment struct {
	EnvType string
	Name    string
	
	containers     map[string]*docker.Container
	containerMutex sync.RWMutex
	images         map[string]*docker.Image
	imageMutex     sync.RWMutex
	networks       map[string]*docker.Network
	networkMutex   sync.RWMutex
}

func (e *LocalEnvironment) GetName() string {
	return e.Name
}
func (e *LocalEnvironment) SetName(name string) {
	e.Name = name
}
func (e *LocalEnvironment) GetEnvType() string {
	return e.EnvType
}
func (e *LocalEnvironment) SetEnvType(envType string) {
	e.EnvType = envType
}

func (e *LocalEnvironment) GetEnvDescription() utils.EnvDescription {
	return utils.EnvDescription{
		Name:    e.Name,
		EnvType: e.EnvType,
	}
}

func NewLocalEnv() *LocalEnvironment {
	env := &LocalEnvironment{
		Name:           "Local",
		EnvType:        "local",
		containers:     make(map[string]*docker.Container),
		containerMutex: sync.RWMutex{},
	}

	return env
}
