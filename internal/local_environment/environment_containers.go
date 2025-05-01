package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"fmt"
	"log"
	"strings"
)

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

	e.containers = make(map[string]*docker.Container)
	for _, container := range containers {
		e.containers[container.ID] = &container
	}
}

func (e *LocalEnvironment) GetContainers() map[string]*docker.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers
}

func (e *LocalEnvironment) GetContainer(id string) *docker.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers[id]
}

func (e *LocalEnvironment) GetContainerCount() int {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return len(e.containers)
}
