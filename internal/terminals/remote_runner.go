package terminals

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/logger"
	"Shipyard/internal/utils"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"time"
)

type RemoteRunner struct {
	Command      []string
	Ctx          context.Context    `json:"-"`
	CancelFunc   context.CancelFunc `json:"-"`
	ID           string
	Env          env_manager.RemoteEnvironment
	OutputFn     func(string)
	OutputMetaFn func(status utils.ActionStatus)
	DeleteFn     func()

	finished bool
}

func (r *RemoteRunner) Cancel() {
	if r.finished {
		return
	}

	message := map[string]interface{}{
		"Type":   "Runner",
		"Action": "Cancel",
	}

	err := r.Env.SendMessage(message, r.ID)
	if err != nil {
		log.Err(err).
			Strs("cmd", r.Command).
			Msg("Failed to cancel remote runner")
	}
}

func (r *RemoteRunner) Delete() {
	message := map[string]interface{}{
		"Type":   "Runner",
		"Action": "Delete",
	}

	err := r.Env.SendMessage(message, r.ID)
	if err != nil {
		log.Err(err).
			Strs("cmd", r.Command).
			Msg("Failed to delete remote runner")
	}
}

func (r *RemoteRunner) Retry() {
	message := map[string]interface{}{
		"Type":   "Runner",
		"Action": "Retry",
	}

	r.OutputMetaFn(utils.Pending)

	// add buffering, as else it is prone to skip messages
	messageChannel := make(chan []byte, 100)
	r.Env.AddMessageChan(r.ID, messageChannel)

	go r.WaitForWSMessages(messageChannel)

	err := r.Env.SendMessage(message, r.ID)
	if err != nil {
		log.Err(err).
			Strs("cmd", r.Command).
			Msg("Failed to retry remote runner")
	}
}

func (r *RemoteRunner) Run() {
	// react to cancels
	go func() {
		<-r.Ctx.Done()
		r.Cancel()
	}()

	message := map[string]interface{}{
		"Type":    "Runner",
		"Action":  "Create",
		"Command": r.Command,
	}

	r.OutputMetaFn(utils.Pending)

	// add buffering, as else it is prone to skip messages
	messageChannel := make(chan []byte, 100)
	r.Env.AddMessageChan(r.ID, messageChannel)

	go r.WaitForWSMessages(messageChannel)

	err := r.Env.SendMessage(message, r.ID)
	if err != nil {
		r.OutputFn("Failed to start action\r\n")
		r.OutputMetaFn(utils.Failed)
	}
}

func (r *RemoteRunner) WaitForWSMessages(messageChannel chan []byte) {
	toBeDismissed := make(chan bool)

	for {
		select {
		case msg := <-messageChannel:
			r.HandleWSMessage(msg)
		case <-r.Ctx.Done():
			go func() {
				time.Sleep(5 * time.Second)
				toBeDismissed <- true
			}()
		case <-toBeDismissed:
			r.Env.RemoveMessageChan(r.ID)
			return
		}
	}
}

func (r *RemoteRunner) HandleWSMessage(msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.HandleSimpleRecoverPanic(r, "Error handling a Remote Runner message")
		}
	}()

	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		// ignore non-json messages
		return
	}

	key, ok := message["Key"].(string)
	if !ok || key != r.ID {
		return
	}
	log.Info().
		Str("Msg", string(msg)).
		Msg("Received message")

	data, ok := message["Data"].(map[string]interface{})
	if !ok {
		return
	}

	msgType := data["Type"].(string)
	switch msgType {
	case "OutputFn":
		r.OutputFn(data["Output"].(string))
		break
	case "DeleteFn":
		go func() {
			time.Sleep(3 * time.Second)
			r.Env.RemoveMessageChan(r.ID)
		}()

		r.CancelFunc()
		r.DeleteFn()
		go func() {
			time.Sleep(10 * time.Second)
			r.Delete()
		}()
		break
	case "OutputMetaFn":
		status := utils.ActionStatus(int64(data["ActionStatus"].(float64)))
		r.OutputMetaFn(status)
		if status == utils.Success {
			r.finished = true
		}
		break
	default:
		log.Warn().
			Str("msgType", msgType).
			Msg("Unhandled message type")
	}
}
