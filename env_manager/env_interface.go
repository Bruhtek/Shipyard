package env_manager

import (
	"Shipyard/docker"
	"Shipyard/utils"
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

	GetEnvDescription() utils.EnvDescription
}
