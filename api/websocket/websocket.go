package websocket

import (
	"Shipyard/utils"
	"github.com/gorilla/websocket"
	"log"
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
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	conn, err := upgrader.Upgrade(w, r, nil)
	utils.IFErr(err, "Websocket upgrade error")

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	utils.IFErr(err, "Websocket write error")

	ConnectionManager.TryAddConnection(conn)
}
