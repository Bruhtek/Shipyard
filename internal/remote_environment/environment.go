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
	lastHeartbeat time.Time
	lastRequest   time.Time
	mutex         sync.RWMutex // - for various container operations

	Connection *websocket.Conn
}

func (r *RemoteEnvironment) GetEnvDescription() utils.EnvDescription {
	return utils.EnvDescription{
		Name:      r.Name,
		EnvType:   r.EnvType,
		EnvKey:    r.EnvKey,
		Heartbeat: r.IsHeartbeat(),
	}
}

func (r *RemoteEnvironment) Request() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.lastRequest = time.Now()
}
func (r *RemoteEnvironment) IsRequested() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	timeSince := time.Since(r.lastRequest)
	return timeSince < shared_config.RemoteRequestedTimeout
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

	r.lastHeartbeat = time.Now()
}
func (r *RemoteEnvironment) IsHeartbeat() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	timeSince := time.Since(r.lastHeartbeat)
	return timeSince < shared_config.RemoteHeartbeatInterval*2
}

func (r *RemoteEnvironment) ClearConnection() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Connection != nil {
		r.Connection.Close()
	}

	r.Connection = nil
}

func (r *RemoteEnvironment) SetConnection(conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Connection != nil {
		r.Connection.Close()
	}

	r.Connection = conn

	go func() {
		defer r.ClearConnection()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			go WebsocketMessageHandler(message)
		}
	}()
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

	return env
}
