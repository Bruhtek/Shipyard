package internal

import (
	"Shipyard/internal/local_environment"
	"sync"
)

type RemoteEnvironment struct {
	MainHost         string
	RemoteToken      string
	ConnectionString string
	UseHttps         bool

	HasSuccessfullyConnected bool

	Mutex sync.RWMutex

	environment *local_environment.LocalEnvironment
}

var RemoteEnv *RemoteEnvironment = NewRemoteEnvironment()

func NewRemoteEnvironment() *RemoteEnvironment {
	return &RemoteEnvironment{}
}
