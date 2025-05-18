package env_manager

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/utils"
)

type EnvInterface interface {
	GetName() string
	SetName(name string)

	GetEnvType() string
	SetEnvType(envType string)

	ScanContainers()
	GetContainers() map[string]*docker.Container
	GetContainer(id string) *docker.Container
	GetContainerCount() int

	ScanImages()
	GetImages() map[string]*docker.Image
	GetImage(id string) *docker.Image
	GetImageCount() int

	ScanNetworks()
	GetNetworks() map[string]*docker.Network
	GetNetwork(idOrName string) *docker.Network
	GetNetworkCount() int

	GetEnvDescription() utils.EnvDescription
}
