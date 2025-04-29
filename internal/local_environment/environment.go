package local_environment

import (
	docker2 "Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"Shipyard/internal/utils"
	"fmt"
	"log"
	"strings"
	"sync"
)

type LocalEnvironment struct {
	EnvType        string
	Name           string
	containers     map[string]*docker2.Container
	containerMutex sync.RWMutex

	images     map[string]*docker2.Image
	imageMutex sync.RWMutex
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

func (e *LocalEnvironment) ScanContainers() {
	e.containerMutex.Lock()
	defer e.containerMutex.Unlock()

	out, err := terminals.RunSimpleCommand("docker ps -a --format json --no-trunc")
	if err != nil {
		log.Println("Failed to list containers: ", err)
		return
	}

	containers := ParsePsJson([]byte(out))
	for id, container := range containers {
		currentContainer, ok := e.containers[container.ID]
		// image ID is immutable, so we can skip expensive inspect command if we already have it
		if ok {
			containers[id].ImageID = currentContainer.ImageID
		} else {
			out, err = terminals.RunSimpleCommand(
				fmt.Sprintf("docker container inspect --format '{{.Image}}' %s", container.ID))
			if err != nil {
				log.Println("Failed to inspect container: ", err)
				continue
			}
			containers[id].ImageID = strings.Trim(strings.TrimSpace(out), "'")
		}
	}

	e.containers = make(map[string]*docker2.Container)
	for _, container := range containers {
		e.containers[container.ID] = &container
	}
}

func (e *LocalEnvironment) ScanImages() {
	e.imageMutex.Lock()
	defer e.imageMutex.Unlock()

	out, err := terminals.RunSimpleCommand("docker image ls --format json --no-trunc")
	if err != nil {
		log.Println("Failed to list images: ", err)
		return
	}

	// TODO: check if an image is dangling
	//_, err = terminals.RunSimpleCommand("docker images -f dangling=true -q --no-trunc")
	//if err != nil {
	//	log.Println("Failed to list dangling images: ", err)
	//	return
	//}

	images := ParseImageLsJson([]byte(out))
	for num, image := range images {
		currentImage, ok := e.images[image.ID]
		if ok && currentImage.RepoDigests != nil {
			images[num].RepoDigests = currentImage.RepoDigests
		} else {
			out, err = terminals.RunSimpleCommand(
				fmt.Sprintf("docker image inspect --format {{.RepoDigests}} %s", image.ID))
			if err != nil {
				log.Println("Failed to inspect image: ", err)
				continue
			}
			processedOut := strings.Split(strings.Trim(strings.TrimSpace(out), "[]"), ",")
			images[num].RepoDigests = make([]string, len(processedOut))
			for i, digest := range processedOut {
				images[num].RepoDigests[i] = strings.Trim(strings.TrimSpace(digest), "'\"")
			}
		}
	}

	e.images = make(map[string]*docker2.Image)
	for _, image := range images {
		e.images[image.ID] = &image
	}

	ids := make([]string, 0)
	for id := range e.images {
		ids = append(ids, id)
	}
	usedIds := e.getUsedImageIds(ids)
	for _, id := range usedIds {
		e.images[id].Used = true
	}

	//danglignIds = strings.TrimSpace(danglignIds)
	//danglignIdsList := strings.Split(danglignIds, "\n")
	//for _, id := range danglignIdsList {
	//	id = strings.Trim(strings.TrimSpace(id), "'")
	//	if id == "" {
	//		continue
	//	}
	//
	//}
}

func (e *LocalEnvironment) GetImages() map[string]*docker2.Image {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	return e.images
}

func (e *LocalEnvironment) GetImage(id string) *docker2.Image {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	image, ok := e.images[id]
	if !ok {
		return nil
	}

	return image
}

func (e *LocalEnvironment) GetImageCount() int {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	return len(e.images)
}

func (e *LocalEnvironment) GetContainers() map[string]*docker2.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers
}

func (e *LocalEnvironment) GetContainer(id string) *docker2.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers[id]
}

func (e *LocalEnvironment) GetContainerCount() int {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return len(e.containers)
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
		containers:     make(map[string]*docker2.Container),
		containerMutex: sync.RWMutex{},
	}

	return env
}
