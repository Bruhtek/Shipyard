package remote_controller

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *RemoteEnvironment) Heartbeat() {
	r.heartbeatMutex.Lock()
	defer r.heartbeatMutex.Unlock()

	r.LastHeartbeat = time.Now()
}
func (r *RemoteEnvironment) HasHeartbeat() bool {
	r.heartbeatMutex.RLock()
	defer r.heartbeatMutex.RUnlock()

	return time.Now().Sub(r.LastHeartbeat) < MAX_HEARTBEAT_INTERVAL
}

func (r *RemoteEnvironment) Connect(conn *websocket.Conn) {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	r.Connection = conn
	log.Info().
		Str("remote_environment", r.Name).
		Msg("Connected to remote environment")

	go func() {
		interval := 5 * time.Second
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if r.IsNeeded() == false {
					log.Info().
						Str("remote_environment", r.Name).
						Msg("Remote environment is no longer needed. Disconnecting...")
					r.Disconnect()
					return
				}
				if r.IsConnected() == false {
					return
				}
			}
		}
	}()
}
func (r *RemoteEnvironment) IsConnected() bool {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	return r.Connection != nil
}
func (r *RemoteEnvironment) Disconnect() {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	if r.Connection != nil {
		r.Connection.Close()
	}
	r.Connection = nil
}

func (r *RemoteEnvironment) Need() {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	r.LastNeeded = time.Now()
}

func (r *RemoteEnvironment) IsNeeded() bool {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	return time.Now().Sub(r.LastNeeded) < NEED_TIMEOUT
}
