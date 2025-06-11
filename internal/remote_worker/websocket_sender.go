package remote_worker

import "github.com/rs/zerolog/log"

func (c *ConnectionManager) SendResponse(key string, data map[string]interface{}) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	type Response struct {
		Key  string
		Data map[string]interface{}
	}

	resp := Response{
		Key:  key,
		Data: data,
	}
	_, ok := data["Body"]
	if !ok {
		log.Info().Interface("response", resp).Send()
	}

	err := c.conn.WriteJSON(resp)
	if err != nil {
		log.Err(err).
			Str("key", key).
			Msg("Cannot send a WS response")
		return
	}
}
