package remote_worker

import (
	"Shipyard/internal/logger"
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func (c *ConnectionManager) HandleMessage(message []byte) {

	if message == nil || string(message) == "" || len(message) == 0 {
		return
	}

	if string(message) == "Connected" {
		log.Info().
			Msg("Controller acknowledged WS Connection")
		return
	}

	var msg map[string]interface{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Err(err).
			Str("message", string(message)).
			Msg("Error unmarshalling controller message")
		return
	}

	log.Debug().
		Interface("msg", msg).Msg("Received message from controller")

	key, ok := msg["Key"].(string)
	if !ok || key == "" {
		log.Warn().
			Msg("Received message from controller without a Key. Ignoring.")
		return
	}

	if msg["Type"] == "API" {
		c.HandleAPI(key, msg)
	}

	if msg["Type"] == "Runner" {
		c.HandleRunner(key, msg)
	}
}

func (c *ConnectionManager) HandleAPI(key string, msg map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.HandleSimpleRecoverPanic(err, "Error handling API message")
		}
	}()

	data := msg["Data"].(map[string]interface{})
	var req *http.Request
	pathEscaped := strings.ReplaceAll(url.PathEscape(data["Path"].(string)), "%2F", "/")

	if data["Body"] != nil {
		req = httptest.NewRequest(data["Method"].(string), pathEscaped, bytes.NewBuffer([]byte(data["Body"].(string))))
	} else {
		req = httptest.NewRequest(data["Method"].(string), pathEscaped, nil)
	}
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	//log.Info().
	//	Str("body", w.Body.String()).
	//	Str("method", req.Method).
	//	Str("path", req.URL.Path).
	//	Msg("Result of evaluating an API request")

	responseData := map[string]interface{}{
		"Body": w.Body.String(),
		"Code": w.Result().StatusCode,
	}

	c.SendResponse(key, responseData)
}

func (c *ConnectionManager) HandleRunner(key string, msg map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.HandleSimpleRecoverPanic(err, "Error handling Runner message")
		}
	}()

	msgType := msg["Action"].(string)
	id := key

	switch msgType {
	case "Create":
		cmdInterface := msg["Command"].([]interface{})
		cmd := make([]string, len(cmdInterface))
		for i, v := range cmdInterface {
			cmd[i] = v.(string)
		}
		RunnerManager.NewRunner(id, cmd)
		break
	case "Retry":
		RunnerManager.RetryRunner(id)
		break
	case "Cancel":
		RunnerManager.CancelRunner(id)
		break
	case "Delete":
		RunnerManager.DeleteRunner(id)
		break
	default:
		log.Warn().
			Str("type", msgType).
			Msg("Unhandled Runner message type")
	}
}
