package types

import (
	"Shipyard/internal/docker"
	"github.com/gorilla/websocket"
)

type EnvType string

const (
	EnvTypeLocal  EnvType = "local"
	EnvTypeRemote EnvType = "remote"
)

type EnvInterface interface {
	GetName() string
	SetName(name string)

	GetEnvType() EnvType
	GetEnvDescription() EnvDescription
}

type LocalEnvironment interface {
	EnvInterface
	ScanContainers()
	GetContainers() map[string]*docker.Container
	GetContainer(id string) *docker.Container
	GetContainerCount() int

	ScanImages()
	GetImages() map[string]*docker.Image
	GetImage(id string) *docker.Image
	GetImageCount() int

	ScanNetworks()
	GetNetworks() map[string]*docker.Network
	GetNetwork(idOrName string) *docker.Network
	GetNetworkCount() int

	ScanStacks()
	GetStacks() map[string]*docker.Stack
	GetStack(idOrName string) *docker.Stack
	GetStackCount() int
}

type RemoteEnvironment interface {
	EnvInterface
	Heartbeat()
	HasHeartbeat() bool

	Connect(conn *websocket.Conn)
	IsConnected() bool
	Disconnect()

	Need()
	IsNeeded() bool

	GetResponse(path string) (RequestResponse, error)
	PostResponse(path string, body string) (RequestResponse, error)

	SendMessage(message map[string]interface{}, key string) error
	SendMessageWaitForResponse(key string, messageType string, data map[string]interface{}) (string, error)

	AddMessageChan(key string, channel chan []byte)
	RemoveMessageChan(key string)
}
