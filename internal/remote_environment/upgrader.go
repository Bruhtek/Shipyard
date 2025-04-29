package remote_environment

import (
	"Shipyard/internal/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	remote := r.Context().Value("remote").(*RemoteEnvironment)

	conn, err := upgrader.Upgrade(w, r, nil)
	utils.IFErr(err, "Websocket upgrade error")

	remote.mutex.Lock()
	defer remote.mutex.Unlock()
	remote.Connection = conn
}
