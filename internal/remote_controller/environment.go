package remote_controller

import (
	"Shipyard/internal/utils"
	"github.com/gorilla/websocket"
	"time"
)

const (
	MAX_HEARTBEAT_INTERVAL = time.Second * 7
	NEED_TIMEOUT           = time.Second * 30
)

type RemoteEnvironment struct {
	Name string
	Key  string

	LastHeartbeat time.Time
	Connection    *websocket.Conn

	// time when the environment was last requested by a user
	LastNeeded time.Time
}

func (r *RemoteEnvironment) GetName() string {
	return r.Name
}

func (r *RemoteEnvironment) SetName(name string) {
	r.Name = name
}

func (r *RemoteEnvironment) GetEnvType() string {
	return "remote"
}

func (r *RemoteEnvironment) GetEnvDescription() utils.EnvDescription {
	return utils.EnvDescription{
		EnvType:   "remote",
		Name:      r.Name,
		Connected: r.IsConnected(),
	}
}

func NewRemoteEnv(key string) *RemoteEnvironment {
	env := &RemoteEnvironment{
		Key: key,
	}

	return env
}
