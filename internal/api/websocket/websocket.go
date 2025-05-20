package websocket

import (
	"Shipyard/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections by default TODO: Change this
	},
}

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logger.HandleSimpleRecoverPanic(r, "[WS] Error while handling a new websocket connection")
		}
	}()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).
			Str("request", r.URL.String()).
			Msg("[WS] Error while upgrading connection")
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	if err != nil {
		log.Err(err).
			Str("request", r.URL.String()).
			Msg("[WS] Error while sending a 'Connected' message")
		return
	}

	ConnectionManager.TryAddConnection(conn)
}
