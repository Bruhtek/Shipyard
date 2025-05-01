package websocket

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/terminals"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

func Handler(data ConnectionData, conn *websocket.Conn, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WS] Recovered from panic while handling: %v", r)
		}
	}()

	log.Println("Received message:", string(message))

	var msg map[string]interface{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("[WS] Error unmarshalling message:", err)
		return
	}

	// basic validation
	envName, ok1 := msg["Environment"].(string)
	objectType, ok2 := msg["Object"].(string)
	action, ok3 := msg["Action"].(string)
	objectId, ok4 := msg["ObjectId"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Println("[WS] Invalid message format")
		return
	}

	env := env_manager.EnvManager.GetEnv(envName)
	if env == nil {
		log.Println("[WS] Environment not found:", envName)
		return
	}

	println("[WS] Received message:", objectType, action, envName)

	cmd := GetDockerCommand(objectType, action, objectId)
	if len(cmd) == 0 {
		log.Println("[WS] Invalid command:", objectType, action, objectId)
		return
	}

	broadcaster := Broadcaster{
		BroadcastFn:     ConnectionManager.BroadcastActionOutput,
		BroadcastMetaFn: ConnectionManager.BroadcastActionMetadata,
		BroadcastMiscFn: ConnectionManager.BroadcastActionMisc,
	}
	actionObj := NewBroadcastAction(cmd, broadcaster, envName, objectType, action, objectId)

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
