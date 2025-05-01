package internal

import (
	"Shipyard/internal/api/websocket"
	"Shipyard/internal/utils"
	"context"
	"log"
	"sync"
	"time"
)

func WSHandler(msg map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WS] Recovered from panic while handling: %v", r)
		}
	}()

	// basic validation
	envName, ok1 := msg["Environment"].(string)
	objectType, ok2 := msg["Object"].(string)
	action, ok3 := msg["Action"].(string)
	objectId, ok4 := msg["ObjectId"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Println("[WS] Invalid message format")
		return
	}

	println("[WS] Received message:", objectType, action, envName)

	cmd := websocket.GetDockerCommand(objectType, action, objectId)
	if len(cmd) == 0 {
		log.Println("[WS] Invalid command:", objectType, action, objectId)
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	actionId := utils.RandString(32)

	actionObj := websocket.Action{
		Environment:   envName,
		Object:        objectType,
		Action:        action,
		ObjectId:      objectId,
		ActionId:      actionId,
		InitializedBy: "", // TODO: set this to the user id
		StartedAt:     time.Now(),
		FinishedAt:    time.Time{},
		Status:        websocket.Pending,
		Output:        "",
		Command:       cmd,
		Ctx:           ctx,
		CancelFunc:    cancelFunc,
		Mutex:         sync.RWMutex{},
	}

	runner := websocket.Runner{
		Command:  cmd,
		ActionId: actionObj.ActionId,
		Action:   &actionObj,
		Ctx:      ctx,
	}

	go runner.Run()
}
