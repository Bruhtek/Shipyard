package internal

import "github.com/gorilla/websocket"

type Status struct {
	Connection *websocket.Conn
}

func (s *Status) IsConnected() bool {
	return s.Connection != nil
}

func NewStatus() *Status {
	return &Status{
		Connection: nil,
	}
}
