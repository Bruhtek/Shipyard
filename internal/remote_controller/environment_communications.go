package remote_controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *RemoteEnvironment) HandleMessage(message []byte) {
	r.messageChannelsMutex.RLock()
	defer r.messageChannelsMutex.RUnlock()
	
	for _, channel := range r.MessageChannels {
		select {
		case channel <- message:
			// message sent
		default:
			// no receiver, ignore
		}
	}
}

func (r *RemoteEnvironment) AddMessageChan(key string, channel chan []byte) {
	r.messageChannelsMutex.Lock()
	defer r.messageChannelsMutex.Unlock()

	r.MessageChannels[key] = channel
}

func (r *RemoteEnvironment) RemoveMessageChan(key string) {
	r.messageChannelsMutex.Lock()
	defer r.messageChannelsMutex.Unlock()
	delete(r.MessageChannels, key)
}

func (r *RemoteEnvironment) SendMessageWaitForResponse(key string, messageType string, data map[string]interface{}) (string, error) {
	r.Need()
	type RequestData struct {
		Key  string
		Data map[string]interface{}
		Type string
	}

	req := RequestData{
		Key:  key,
		Data: data,
		Type: messageType,
	}

	marshaled, err := json.Marshal(req)
	if err != nil {
		log.Err(err).
			Interface("data", req).
			Str("remote", r.GetName()).
			Msg("Error marshalling data for remote")
		return "", err
	}

	log.Debug().
		Str("remote", r.GetName()).
		Str("key", key).
		Str("data", string(marshaled)).
		Msg("Sending data to remote")

	messageChan := make(chan []byte)
	r.AddMessageChan(key, messageChan)
	defer r.RemoveMessageChan(key)

	if res := r.waitForConnection(); !res {
		return "", errors.New("remote connection timed out")
	}

	r.connMutex.Lock()
	err = r.Connection.WriteMessage(websocket.TextMessage, marshaled)
	r.connMutex.Unlock()
	if err != nil {
		log.Err(err).
			Interface("data", req).
			Str("remote", r.GetName()).
			Msg("Error sending data to remote")
		return "", err
	}

	response := make(chan string)
	const MAX_MESSAGE_WAIT_TIME = time.Second * 5
	go func() {
		defer close(response)
		for {
			select {
			case <-time.After(MAX_MESSAGE_WAIT_TIME):
				log.Warn().
					Str("remote", r.GetName()).
					Interface("data", req).
					Msg("No response received from remote within the timeout period")
				return
			case message, ok := <-messageChan:
				if !ok {
					log.Err(err).
						Str("remote", r.GetName()).
						Msg("Error reading message from remote")
					return
				}

				var resp map[string]interface{}
				err = json.Unmarshal(message, &resp)
				if err != nil {
					log.Err(err).
						Str("remote", r.GetName()).
						Msg("Error unmarshalling response from remote")
					continue
				}

				if resp["Key"] == key {
					response <- string(message)
					return
				}
			}
		}
	}()

	res, ok := <-response
	if !ok {
		return "", errors.New("no response from remote")
	}
	return res, nil
}

func (r *RemoteEnvironment) SendMessage(message map[string]interface{}, key string) error {
	r.Need()

	message["Key"] = key

	marshaled, err := json.Marshal(message)
	if err != nil {
		log.Err(err).
			Interface("data", message).
			Str("remote", r.GetName()).
			Msg("Error marshalling data for remote")
	}

	if res := r.waitForConnection(); !res {
		return errors.New("remote connection timed out")
	}

	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	err = r.Connection.WriteMessage(websocket.TextMessage, marshaled)
	if err != nil {
		log.Err(err).
			Interface("data", message).
			Str("remote", r.GetName()).
			Msg("Error sending data to remote")
		return err
	}

	return nil
}
