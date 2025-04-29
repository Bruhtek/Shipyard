package internal

import "sync"

type RemoteEnvironment struct {
	MainHost         string
	RemoteToken      string
	ConnectionString string
	UseHttps         bool

	HasSuccessfullyConnected bool

	Mutex sync.RWMutex
}

var RemoteEnv *RemoteEnvironment = NewRemoteEnvironment()

func NewRemoteEnvironment() *RemoteEnvironment {
	return &RemoteEnvironment{}
}
