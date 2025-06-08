package remote_controller

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *RemoteEnvironment) Heartbeat() {
	r.LastHeartbeat = time.Now()
}
func (r *RemoteEnvironment) HasHeartbeat() bool {
	return time.Now().Sub(r.LastHeartbeat) < MAX_HEARTBEAT_INTERVAL
}

func (r *RemoteEnvironment) Connect(conn *websocket.Conn) {
	r.Connection = conn
	log.Info().
		Str("remote_environment", r.Name).
		Msg("Connected to remote environment")
}
func (r *RemoteEnvironment) IsConnected() bool {
	return r.Connection != nil
}
func (r *RemoteEnvironment) Disconnect() {
	if r.Connection != nil {
		r.Connection.Close()
	}
	r.Connection = nil
}

func (r *RemoteEnvironment) Need() {
	r.LastNeeded = time.Now()
}

func (r *RemoteEnvironment) IsNeeded() bool {
	return time.Now().Sub(r.LastNeeded) < NEED_TIMEOUT
}
