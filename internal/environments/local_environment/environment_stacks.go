package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"

	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) ScanStacks() {
	e.stackMutex.Lock()
	defer e.stackMutex.Unlock()

	out, err := terminals.RunSimpleCommand("docker compose ls --format json")
	if err != nil {
		log.Err(err).Msg("Error listing stacks")
		return
	}

	composes := ParseComposeJson(&out)
	if composes == nil {
		return
	}

	e.stacks = make(map[string]*docker.Stack)
	for _, compose := range composes {
		compose.Containers = e.GetContainersInStack(compose.ConfigFiles)
		compose.Networks = e.GetNetworksInStack(compose.Name)

		e.stacks[compose.ConfigFiles] = compose
	}
}

func (e *LocalEnvironment) GetContainersInStack(configFile string) []*docker.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	containers := make([]*docker.Container, 0)
	for _, container := range e.containers {
		if container.Labels["com.docker.compose.project.config_files"] == configFile {
			containers = append(containers, container)
		}
	}

	return containers
}

func (e *LocalEnvironment) GetNetworksInStack(name string) []*docker.Network {
	e.networkMutex.RLock()
	defer e.networkMutex.RUnlock()

	networks := make([]*docker.Network, 0)
	for _, network := range e.networks {
		if network.Labels["com.docker.compose.project"] == name {
			networks = append(networks, network)
		}
	}

	return networks
}

func (e *LocalEnvironment) GetStacks() map[string]*docker.Stack {
	e.stackMutex.RLock()
	defer e.stackMutex.RUnlock()

	return e.stacks
}
func (e *LocalEnvironment) GetStack(configFile string) *docker.Stack {
	e.stackMutex.RLock()
	defer e.stackMutex.RUnlock()

	stack, ok := e.stacks[configFile]
	if ok {
		return stack
	}
	return nil
}
func (e *LocalEnvironment) GetStackCount() int {
	e.stackMutex.RLock()
	defer e.stackMutex.RUnlock()

	return len(e.stacks)
}
