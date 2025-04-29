package shared_config

import "time"

const (
	RemoteHeartbeatInterval = 5 * time.Second
	RemoteRequestedTimeout  = 5 * time.Minute
)
