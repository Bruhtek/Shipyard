package remote_controller

import (
	"Shipyard/internal/utils"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type RequestData struct {
	Path   string
	Method string
	Body   string
}

type RequestResponse struct {
	Body string
	Code int
}

func (r *RemoteEnvironment) GetResponse(path string) (RequestResponse, error) {
	r.Need()
	connected := r.waitForConnection()
	if !connected {
		log.Warn().
			Str("path", path).
			Str("remote", r.GetName()).
			Msg("Failed to connect to the remote")
		return RequestResponse{}, errors.New("connection timed out")
	}

	// for remote, the required environment is the local one
	path = strings.ReplaceAll(path, r.GetName(), "local")

	key := utils.RandString(32)
	message := map[string]interface{}{
		"Path":   path,
		"Method": "GET",
	}

	res, err := r.SendMessageWithData(key, "API", message)
	if err != nil {
		return RequestResponse{}, err
	}

	parsed := parseResponseData(res)

	return parsed, err
}
func (r *RemoteEnvironment) PostResponse(path string, body string) (RequestResponse, error) {
	r.Need()
	connected := r.waitForConnection()
	if !connected {
		log.Warn().
			Str("path", path).
			Str("remote", r.GetName()).
			Msg("Failed to connect to the remote")
		return RequestResponse{}, errors.New("connection timed out")
	}

	return RequestResponse{}, nil
}

func parseResponseData(data string) (res RequestResponse) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Stack().Msgf("%v", err)

			res.Code = 500
			res.Body = "Internal Server Error"
		}
	}()

	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(data), &parsed)
	if err != nil {
		return RequestResponse{
			Code: 500,
			Body: "Internal Server Error",
		}
	}

	resData := parsed["Data"].(map[string]interface{})
	res.Body = resData["Body"].(string)
	res.Code = int(resData["Code"].(float64))
	return res
}

const (
	CONNECTION_CHECK_EVERY = time.Millisecond * 100
	WAIT_UNTIL_TIMEOUT     = time.Second * 10
)

func (r *RemoteEnvironment) waitForConnection() bool {
	connected := r.IsConnected()
	if connected {
		return true
	}

	start := time.Now()
	for range time.Tick(CONNECTION_CHECK_EVERY) {
		connected = r.IsConnected()
		if connected {
			return true
		}
		if time.Since(start) > WAIT_UNTIL_TIMEOUT {
			return false
		}
	}

	return false
}
