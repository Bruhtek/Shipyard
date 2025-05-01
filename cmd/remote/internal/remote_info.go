package internal

import (
	websocket2 "Shipyard/internal/api/websocket"
	"Shipyard/internal/local_environment"
	"github.com/gorilla/websocket"
	"sync"
)

type RemoteEnvironment struct {
	MainHost                 string
	RemoteToken              string
	ConnectionString         string
	WSConnectionString       string
	UseHttps                 bool
	HasSuccessfullyConnected bool

	Conn *websocket.Conn

	Mutex sync.RWMutex

	environment *local_environment.LocalEnvironment
}

func (r *RemoteEnvironment) Connect() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Conn != nil {
		// already connected
		return
	}

	var err error
	r.Conn, _, err = websocket.DefaultDialer.Dial(r.WSConnectionString, nil)
	if err != nil {
		panic("ERROR: Unable to connect to main server websocket: " + err.Error())
	}

	println("Connected to main server websocket")

	r.Conn.WriteMessage(websocket.TextMessage, []byte("Connected"))

	connections := make(map[*websocket.Conn]websocket2.ConnectionData)
	connections[r.Conn] = websocket2.ConnectionData{
		Id: r.Conn.RemoteAddr().String(),
	}
	websocket2.ConnectionManager.UnsafeOverrideConnections(connections)

	go func() {
		defer r.Disconnect()

		for {
			_, message, err := r.Conn.ReadMessage()
			if err != nil {
				break
			}

			go RequestHandler(message)
		}
	}()
}

func (r *RemoteEnvironment) Disconnect() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Conn != nil {
		r.Conn.Close()
		r.Conn = nil
	}

	websocket2.ConnectionManager.UnsafeOverrideConnections(make(map[*websocket.Conn]websocket2.ConnectionData))
}
func (r *RemoteEnvironment) IsConnected() bool {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Conn != nil
}

var RemoteEnv *RemoteEnvironment = NewRemoteEnvironment()

func NewRemoteEnvironment() *RemoteEnvironment {
	return &RemoteEnvironment{}
}
