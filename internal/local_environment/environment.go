package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminal_simple"
	"Shipyard/internal/utils"
	"github.com/rs/zerolog/log"
	"sync"
)

type LocalEnvironment struct {
	Name string

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
	return "local"
}

func (e *LocalEnvironment) GetEnvDescription() utils.EnvDescription {
	return utils.EnvDescription{
		Name:      e.Name,
		EnvType:   "local",
		Connected: true,
	}
}

func NewLocalEnv() *LocalEnvironment {
	// test if we can access docker
	res, err := terminal_simple.RunSimpleCommand("docker -v")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot access docker! Check permissions. " +
			"Is the docker daemon running? If this is intentional (as in: you want " +
			"to only control remotes), set IGNORE_LOCAL=true")
	}

	log.Info().
		Str("docker -v", res).
		Msg("Creating a new local docker environment")

	env := &LocalEnvironment{
		containers:     make(map[string]*docker.Container),
		containerMutex: sync.RWMutex{},
	}

	return env
}
