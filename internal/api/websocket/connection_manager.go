package websocket

import (
	"Shipyard/internal/utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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

func (m *CMStruct) GetConnection(id string) *websocket.Conn {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for conn, data := range m.connections {
		if data.id == id {
			return conn
		}
	}
	return nil
}

func (m *CMStruct) GetConnectionId(conn *websocket.Conn) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if data, ok := m.connections[conn]; ok {
		return data.id
	}
	return ""
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

	log.Debug().
		Str("connection-id", id).
		Msg("[WS] Adding new WS connection")

	go func() {
		defer m.RemoveConnection(conn)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			Handler(message)
		}
	}()
	return true
}

func (m *CMStruct) RemoveConnection(conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if data, ok := m.connections[conn]; ok {
		log.Debug().Str("connection-id", data.id).Msg("[WS] Removing WS connection")
		conn.Close()
	}

	delete(m.connections, conn)
}

func (m *CMStruct) BroadcastActionOutput(actionId string, message interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for conn := range m.connections {
		type TaskMessage struct {
			ActionId string
			Type     string
			Message  interface{}
		}

		taskMessage := TaskMessage{
			ActionId: actionId,
			Type:     "ActionOutput",
			Message:  message,
		}

		message, err := json.Marshal(taskMessage)
		if err != nil {
			log.Err(err).Msg("[WS] Error marshalling message")
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			m.RemoveConnection(conn)
		}
	}
}

func (m *CMStruct) BroadcastActionMetadata(action *Action) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for conn := range m.connections {
		type ActionMetadata struct {
			ActionId string
			Type     string
			Metadata interface{}
		}

		actionMetadata := ActionMetadata{
			ActionId: action.ActionId,
			Type:     "ActionMetadata",
			Metadata: action,
		}

		message, err := json.Marshal(actionMetadata)
		if err != nil {
			log.Err(err).Msg("[WS] Error marshalling message")
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			m.RemoveConnection(conn)
		}
	}
}

func (m *CMStruct) BroadcastActionMisc(actionId string, messageKey string, message interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for conn := range m.connections {
		messageData := map[string]interface{}{
			"ActionId": actionId,
			"Type":     "Action" + messageKey,
			messageKey: message,
		}

		msg, err := json.Marshal(messageData)
		if err != nil {
			log.Err(err).Msg("[WS] Error marshalling message")
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			m.RemoveConnection(conn)
		}
	}
}
