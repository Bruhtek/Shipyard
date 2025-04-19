package local_environment

import (
	"Shipyard/docker"
	"Shipyard/terminals"
	"encoding/json"
	"log"
	"strings"
	"sync"
)

type LocalEnvironment struct {
	EnvType        string
	Name           string
	containers     map[string]*docker.Container
	containerMutex sync.RWMutex
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

var LocalEnv *LocalEnvironment = newLocalEnv()

func newLocalEnv() *LocalEnvironment {
	env := &LocalEnvironment{
		Name:           "Local",
		EnvType:        "local",
		containers:     make(map[string]*docker.Container),
		containerMutex: sync.RWMutex{},
	}

	return env
}

func ParsePsJson(jsonData []byte) []docker.Container {
	splitData := strings.Split(string(jsonData), "\n")
	containers := make([]docker.Container, 0)

	for _, line := range splitData {
		if line == "" {
			continue
		}

		tempContainer := docker.TempContainer{}
		err := json.Unmarshal([]byte(line), &tempContainer)
		if err != nil {
			continue
		}

		container, err := tempContainer.ToContainer()
		if err != nil {
			log.Printf("Error parsing container: %v", err)
			continue
		}

		containers = append(containers, container)
	}

	return containers
}
