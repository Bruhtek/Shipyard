package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminal_simple"
	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) ScanStacks() {
	e.stackMutex.Lock()
	defer e.stackMutex.Unlock()

	out, err := terminal_simple.RunSimpleCommand("docker compose ls --format json")
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
		e.stacks[compose.ConfigFiles] = compose
	}
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
