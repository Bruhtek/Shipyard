package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"maps"
	"slices"

	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) ScanNetworks() {
	out, err := terminals.RunSimpleCommand(NetworkLsCommand)
	if err != nil {
		log.Err(err).Msg("Error listing networks")
		return
	}

	networks := ParseNetworkLsJson(&out)
	containers := e.GetContainers()

	for i, network := range networks {
		curr := e.GetNetwork(network.ID)
		if curr != nil {
			curr.UpdateNetworkContainers(containers)
			networks[i] = *curr
			continue
		}

		network.UpdateNetworkContainers(containers)
	}

	e.networkMutex.Lock()
	defer e.networkMutex.Unlock()
	e.networks = make(map[string]*docker.Network)
	for _, network := range networks {
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

		id := selectIdPrefixFromList(slices.Collect(maps.Keys(e.networks)), idOrName)
		return e.networks[id]
	}
	return network
}

func (e *LocalEnvironment) GetNetworkCount() int {
	e.networkMutex.RLock()
	defer e.networkMutex.RUnlock()

	return len(e.networks)
}
