package main

import (
	"Shipyard/cmd/remote/internal"
	"os"
	"strings"
)

func main() {
	main_host := os.Getenv("MAIN_HOST")
	remote_token := os.Getenv("REMOTE_TOKEN")
	secure := strings.ToLower(os.Getenv("SECURE"))
	var useHttps bool

	if main_host == "" || remote_token == "" {
		panic("MAIN_HOST and REMOTE_TOKEN variables MUST be set for this remote to operate")
	}
	if secure == "" {
		useHttps = false
	} else {
		if secure == "true" {
			useHttps = true
		} else if secure == "false" {
			useHttps = false
		} else {
			panic("SECURE variable must be set to true or false")
		}
	}

	connectionString := ""
	if useHttps {
		connectionString = "https://"
	} else {
		connectionString = "http://"
	}
	connectionString += main_host + "/api/remote/" + remote_token + "/heartbeat"

	internal.RemoteEnv.MainHost = main_host
	internal.RemoteEnv.RemoteToken = remote_token
	internal.RemoteEnv.UseHttps = useHttps
	internal.RemoteEnv.ConnectionString = connectionString
	internal.RemoteEnv.HasSuccessfullyConnected = false

	internal.HeartbeatLoop()
}
