package remote_environment

import (
	"Shipyard/internal/local_environment"
	"Shipyard/internal/shared_config"
	"Shipyard/internal/utils"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// RemoteEnvironment - this is a remote for the controller (the main host), not the agent
type RemoteEnvironment struct {
	local_environment.LocalEnvironment

	EnvKey        string
	heartbeat     bool // has the heartbeat been received
	lastHeartbeat time.Time
	requested     bool // has the remote been requested by the user
	lastRequest   time.Time
	mutex         sync.RWMutex // - for various container operations

	Connection *websocket.Conn
}

func (r *RemoteEnvironment) GetEnvDescription() utils.EnvDescription {
	return utils.EnvDescription{
		Name:      r.Name,
		EnvType:   r.EnvType,
		EnvKey:    r.EnvKey,
		Heartbeat: r.heartbeat,
	}
}

func (r *RemoteEnvironment) Request() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.requested = true
	r.lastRequest = time.Now()
}
func (r *RemoteEnvironment) IsRequested() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.requested
}
func (r *RemoteEnvironment) IsConnected() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.Connection != nil
}

//#region Remote functions

func (r *RemoteEnvironment) Heartbeat() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.heartbeat = true
	r.lastHeartbeat = time.Now()
}
func (r *RemoteEnvironment) IsHeartbeat() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.heartbeat
}

func (r *RemoteEnvironment) heartbeatKiller() {
	for range time.Tick(shared_config.RemoteHeartbeatInterval) {
		r.mutex.Lock()
		if time.Since(r.lastHeartbeat) > shared_config.RemoteHeartbeatInterval*2 {
			r.heartbeat = false
		}
		if time.Since(r.lastRequest) > shared_config.RemoteRequestedTimeout {
			r.requested = false
		}
		r.mutex.Unlock()
	}
}

//#endregion

func NewRemoteEnvironment(envKey string) *RemoteEnvironment {
	env := &RemoteEnvironment{
		LocalEnvironment: local_environment.LocalEnvironment{
			Name:    "remote",
			EnvType: "remote",
		},
		EnvKey: envKey,
	}

	go env.heartbeatKiller()

	return env
}
