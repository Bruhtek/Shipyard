package websocket

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/terminals"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func Handler(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				log.Err(err).Msg("[WS] Panic while handling message:")
			} else {
				log.Err(err).Msg("[WS] Panic while handling message - unable to cast to error")
			}
		}
	}()

	log.Debug().Str("message", string(message)).Msg("[WS] Received message")

	var msg map[string]interface{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Err(err).
			Str("message", string(message)).
			Msg("[WS] Error unmarshalling message:")
		return
	}

	// basic validation
	envName, ok1 := msg["Environment"].(string)
	objectType, ok2 := msg["Object"].(string)
	action, ok3 := msg["Action"].(string)
	objectId, ok4 := msg["ObjectId"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Error().
			Str("message", string(message)).
			Msg("[WS] Invalid message format")
		return
	}

	env := env_manager.EnvManager.GetEnv(envName)
	if env == nil {
		log.Error().
			Str("environment", envName).
			Msg("[WS] Environment not found")
		return
	}

	log.Debug().
		Str("objectType", objectType).
		Str("action", action).
		Str("envName", envName).
		Msg("[WS] Received message:")

	cmd := GetDockerCommand(objectType, action, objectId)
	if len(cmd) == 0 {
		log.Error().
			Str("objectType", objectType).
			Str("action", action).
			Str("objectId", objectId).
			Str("envName", envName).
			Msg("[WS] Invalid command")
		return
	}

	log.Debug().
		Strs("command", cmd).
		Str("environment", envName).
		Msg("[WS] Running command")

	broadcaster := Broadcaster{
		BroadcastFn:     ConnectionManager.BroadcastActionOutput,
		BroadcastMetaFn: ConnectionManager.BroadcastActionMetadata,
		BroadcastMiscFn: ConnectionManager.BroadcastActionMisc,
	}
	actionObj := NewBroadcastAction(cmd, &broadcaster, envName, objectType, action, objectId)

	runner := terminals.Runner{
		Command:      cmd,
		Ctx:          actionObj.Ctx,
		OutputFn:     actionObj.HandleOutput,
		OutputMetaFn: actionObj.HandleMetadata,
		DeleteFn:     actionObj.HandleDelete,
	}

	ActionManager.createAction(actionObj)

	go runner.Run()
}
