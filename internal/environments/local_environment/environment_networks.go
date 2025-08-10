package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) ScanNetworks() {
	e.networkMutex.Lock()
	defer e.networkMutex.Unlock()

	out, err := terminals.RunSimpleCommand(NetworkLsCommand)
	if err != nil {
		log.Err(err).Msg("Error listing networks")
		return
	}

	networks := ParseNetworkLsJson(&out)
	e.networks = make(map[string]*docker.Network)

	containers := e.GetContainers()

	for _, network := range networks {
		curr, ok := e.networks[network.ID]
		if ok {
			curr.UpdateNetworkContainers(containers)
			continue
		}

		network.UpdateNetworkContainers(containers)
		e.networks[network.ID] = &network
	}
}

func (e *LocalEnvironment) GetNetworks() map[string]*docker.Network {
	e.networkMutex.RLock()
	defer e.networkMutex.RUnlock()

	return e.networks
}

func (e *LocalEnvironment) GetNetwork(idOrName string) *docker.Network {
	e.networkMutex.RLock()
	defer e.networkMutex.RUnlock()

	network, ok := e.networks[idOrName]
	if !ok {
		for _, net := range e.networks {
			if net.Name == idOrName {
				return net
			}
		}
		return nil
	}
	return network
}

func (e *LocalEnvironment) GetNetworkCount() int {
	e.networkMutex.RLock()
	defer e.networkMutex.RUnlock()

	return len(e.networks)
}
