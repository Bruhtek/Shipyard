package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) ScanStacks() {
	out, err := terminals.RunSimpleCommand("docker compose ls -a --format json")
	if err != nil {
		log.Err(err).Msg("Error listing stacks")
		return
	}

	composes := ParseComposeJson(&out)
	if composes == nil {
		return
	}

	stacksInFolder := e.ScanStackFolder(false)

	e.stackMutex.Lock()
	defer e.stackMutex.Unlock()
	e.stacks = make(map[string]*docker.Stack)
	for _, compose := range composes {
		compose.Containers = e.GetContainersInStack(compose.ConfigFiles)
		compose.UpToDate = docker.AreContainersUpToDate(compose.Containers)
		compose.Networks = e.GetNetworksInStack(compose.Name)

		e.stacks[compose.ConfigFiles] = compose
	}

	// add any stacks that are in the stacks folder but not in the docker compose ls output
	// these are the "inactive" stacks - they are not running but have a config file
	for _, stackDir := range stacksInFolder {
		_, ok := e.stacks[stackDir.ConfigFiles]
		if ok {
			continue
		} else {
			e.stacks[stackDir.ConfigFiles] = stackDir
		}
	}

}

var noStacksDirWarn = false
var stacksDirLastScan = time.Time{}
var stacksDirScanInterval = 10 * time.Minute

var lastScannedStacks = make([]*docker.Stack, 0)

func (e *LocalEnvironment) ScanStackFolder(forced bool) []*docker.Stack {
	var empty = make([]*docker.Stack, 0)
	if !forced && time.Since(stacksDirLastScan) < stacksDirScanInterval {
		return lastScannedStacks
	}

	folder := os.Getenv("STACKS_DIR")
	if folder == "" {
		if !noStacksDirWarn {
			log.Warn().Msg("STACKS_DIR environment variable is not set - Shipyard will only show active stacks")
			noStacksDirWarn = true
		}
		return empty
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if !noStacksDirWarn {
			log.Warn().Msgf("STACKS_DIR '%s' does not exist - Shipyard will only show active stacks", folder)
			noStacksDirWarn = true
		}
		return empty
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		log.Err(err).Msgf("Error reading stacks directory '%s'", folder)
		return empty
	}

	stacks := make([]*docker.Stack, 0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		filesInside, err := os.ReadDir(folder + "/" + file.Name())
		if err != nil {
			log.Err(err).Msgf("Error reading directory '%s/%s'", folder, file.Name())
			continue
		}
		configFile := ""
		for _, f := range filesInside {
			allowedNames := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yaml", "compose.yml"}
			if f.IsDir() {
				continue
			}
			for _, name := range allowedNames {
				if f.Name() == name {
					configFile = folder + "/" + file.Name() + "/" + f.Name()
					break
				}
			}
		}
		if configFile == "" {
			log.Debug().Msgf("No valid compose file found in stack directory '%s/%s'", folder, file.Name())
			continue
		}
		stack := &docker.Stack{
			Name:        file.Name(),
			ConfigFiles: configFile,
			Status:      "inactive",
			Containers:  make([]*docker.Container, 0),
			Networks:    make([]*docker.Network, 0),
			UpToDate:    docker.Unknown,
		}
		stacks = append(stacks, stack)
	}
	if len(stacks) == 0 {
		log.Debug().Msgf("No stacks found in directory '%s'", folder)
		return empty
	}
	lastScannedStacks = stacks
	stacksDirLastScan = time.Now()
	return stacks
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
