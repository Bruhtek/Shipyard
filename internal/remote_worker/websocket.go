package remote_worker

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
	"sync"
)

type ConnectionManager struct {
	conn *websocket.Conn

	writeMutex *sync.Mutex
}

var CManager *ConnectionManager = &ConnectionManager{
	writeMutex: &sync.Mutex{},
}

func (c *ConnectionManager) IsConnected() bool {
	return c.conn != nil
}

func (c *ConnectionManager) ConnectToController(uri string, key string) {
	if c.conn != nil {
		return
	}
	log.Info().
		Msg("Controller requested connection")

	protocol := "ws"
	if strings.HasPrefix(uri, "https") {
		protocol = "wss"
	}

	host := uri[strings.Index(uri, "://")+3:]

	u := url.URL{
		Scheme:   protocol,
		Host:     host,
		Path:     "/api/remote/ws",
		RawQuery: "key=" + key,
	}

	log.Info().
		Str("host", host).
		Str("path", u.Path).
		Msg("Connecting to controller")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("host", host).
			Str("path", u.Path).
			Msg("Unable to connect to controller")
	}

	c.useConnection(conn)

	return
}

func (c *ConnectionManager) useConnection(conn *websocket.Conn) {
	c.writeMutex.Lock()
	c.conn = conn
	c.writeMutex.Unlock()

	go func() {
		defer c.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error().
					Err(err).
					Msg("Unable to read message. Disconnecting from controller")
				return
			}
			c.HandleMessage(message)
		}
	}()
}

func (c *ConnectionManager) Close() {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = nil
}
