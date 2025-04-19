package websocket

import (
	"Shipyard/utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type CMStruct struct {
	connections map[*websocket.Conn]ConnectionData
	mutex       sync.RWMutex
}

type ConnectionData struct {
	id string
}

var ConnectionManager = newConnectionManager()

func newConnectionManager() *CMStruct {
	return &CMStruct{
		connections: make(map[*websocket.Conn]ConnectionData),
		mutex:       sync.RWMutex{},
	}
}

func (m *CMStruct) ConnectionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.connections)
}

func (m *CMStruct) TryAddConnection(conn *websocket.Conn) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// if it already exists, don't add it again
	if _, ok := m.connections[conn]; ok {
		return false
	}

	id := utils.RandString(32)

	data := ConnectionData{
		id: id,
	}

	m.connections[conn] = data

	go func() {
		defer m.RemoveConnection(conn)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			Handler(data, conn, message)
		}
	}()
	return true
}

func (m *CMStruct) RemoveConnection(conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if data, ok := m.connections[conn]; ok {
		log.Println("Removing connection with ID:", data.id)
		conn.Close()
	}

	delete(m.connections, conn)
}

func (m *CMStruct) Broadcast(taskId string, message interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for conn := range m.connections {
		type TaskMessage struct {
			TaskId  string
			Message interface{}
		}

		taskMessage := TaskMessage{
			TaskId:  taskId,
			Message: message,
		}

		message, err := json.Marshal(taskMessage)
		if err != nil {
			log.Println("Error marshalling message:", err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			m.RemoveConnection(conn)
		}
	}
}
