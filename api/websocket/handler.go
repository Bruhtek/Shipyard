package websocket

import (
	"Shipyard/env_manager"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
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
	actionId, ok4 := msg["ActionId"].(string)
	objectId, ok5 := msg["ObjectId"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		log.Println("[WS] Invalid message format")
		return
	}

	env := env_manager.EnvManager.GetEnv(envName)
	if env == nil {
		log.Println("[WS] Environment not found:", envName)
		return
	}

	println("[WS] Received message:", objectType, action, envName)

	ctx, _ := context.WithCancel(context.Background())

	actionObj := Action{
		Environment:   envName,
		Object:        objectType,
		Action:        action,
		ObjectId:      objectId,
		ActionId:      actionId,
		InitializedBy: ConnectionManager.GetConnectionId(conn),
		StartedAt:     time.Now(),
		FinishedAt:    time.Time{},
		Status:        Pending,
		Output:        "",
		ctx:           ctx,
		Mutex:         sync.RWMutex{},
	}

	cmd := GetDockerCommand(objectType, action, objectId)
	if len(cmd) == 0 {
		log.Println("[WS] Invalid command:", objectType, action, objectId)
		return
	}

	runner := Runner{
		Command:  cmd,
		ActionId: actionObj.ActionId,
		Action:   &actionObj,
		Ctx:      ctx,
	}

	ActionManager.createAction(&runner, &actionObj)

	go runner.Run()
}
