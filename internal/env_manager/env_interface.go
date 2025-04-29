package env_manager

import (
	docker2 "Shipyard/internal/docker"
	"Shipyard/internal/utils"
)

type EnvInterface interface {
	GetName() string
	SetName(name string)

	GetEnvType() string
	SetEnvType(envType string)

	ScanContainers()
	GetContainers() map[string]*docker2.Container
	GetContainer(id string) *docker2.Container
	GetContainerCount() int

	ScanImages()
	GetImages() map[string]*docker2.Image
	GetImage(id string) *docker2.Image
	GetImageCount() int

	GetEnvDescription() utils.EnvDescription
}
