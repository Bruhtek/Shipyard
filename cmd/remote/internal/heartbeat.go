package internal

import (
	"Shipyard/internal/shared_config"
	"net/http"
	"time"
)

func Heartbeat() {
	RemoteEnv.Mutex.RLock()
	connString := RemoteEnv.ConnectionString
	hasConnected := RemoteEnv.HasSuccessfullyConnected
	RemoteEnv.Mutex.RUnlock()

	res, err := http.Get(connString)
	if err != nil {
		if hasConnected {
			println("WARNING: Lost connection to main server")
		} else {
			panic("ERROR: Unable to connect to main server: " + err.Error())
		}
		return
	}
	if res.StatusCode == http.StatusOK {
		if !hasConnected {
			println("INFO: Successfully connected to main server")

			RemoteEnv.Mutex.Lock()
			defer RemoteEnv.Mutex.Unlock()

			RemoteEnv.HasSuccessfullyConnected = true
		}
		if RemoteEnv.IsConnected() {
			println("INFO: Remote environment is no longer needed")
			RemoteEnv.Disconnect()
		}
	} else if res.StatusCode == http.StatusAccepted {
		if !RemoteEnv.IsConnected() {
			println("INFO: Remote requested connection")

			go RemoteEnv.Connect()
		}
	} else {
		println("WARNING: Unexpected response from main server: " + res.Status)
	}
}

func HeartbeatLoop() {
	Heartbeat()
	for range time.Tick(shared_config.RemoteHeartbeatInterval) {
		Heartbeat()
	}
}
