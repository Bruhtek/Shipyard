package websocket

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/remote_environment"
	"Shipyard/internal/utils"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

func Handler(connectionId string, msg map[string]interface{}) {
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

	env := env_manager.EnvManager.GetEnv(envName)
	if env == nil {
		log.Println("[WS] Environment not found:", envName)
		return
	}

	if env.GetEnvDescription().EnvType == "remote" {
		log.Println("[WS] Remote environment detected, passing action along")
		remote, ok := env.(*remote_environment.RemoteEnvironment)
		if !ok {
			log.Println("[WS] Failed to cast environment to remote environment")
			return
		}
		if !remote.IsConnected() {
			log.Println("[WS] Remote environment is not connected")
			return
		}

		remote.Request()
		remote.PassWSMessage(msg)

		actionObj := Action{
			Environment:   envName,
			Object:        objectType,
			Action:        action,
			ObjectId:      objectId,
			ActionId:      actionId,
			InitializedBy: connectionId,
			StartedAt:     time.Now(),
			FinishedAt:    time.Time{},
			Status:        Pending,
			Output:        "",
			Command:       cmd,
			Ctx:           ctx,
			CancelFunc:    cancelFunc,
			Mutex:         sync.RWMutex{},
		}

		runner := Runner{
			Command:  cmd,
			ActionId: actionObj.ActionId,
			Action:   &actionObj,
			Ctx:      ctx,
		}

		ActionManager.createAction(&actionObj)

		return
	}

	println("[WS] Received message:", objectType, action, envName)

	cmd := GetDockerCommand(objectType, action, objectId)
	if len(cmd) == 0 {
		log.Println("[WS] Invalid command:", objectType, action, objectId)
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	actionId := utils.RandString(32)

	actionObj := Action{
		Environment:   envName,
		Object:        objectType,
		Action:        action,
		ObjectId:      objectId,
		ActionId:      actionId,
		InitializedBy: connectionId,
		StartedAt:     time.Now(),
		FinishedAt:    time.Time{},
		Status:        Pending,
		Output:        "",
		Command:       cmd,
		Ctx:           ctx,
		CancelFunc:    cancelFunc,
		Mutex:         sync.RWMutex{},
	}

	runner := Runner{
		Command:  cmd,
		ActionId: actionObj.ActionId,
		Action:   &actionObj,
		Ctx:      ctx,
	}

	ActionManager.createAction(&actionObj)

	go runner.Run()
}

func HandlerString(connectionId string, message []byte) {
	log.Println("Received message:", string(message))

	var msg map[string]interface{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("[WS] Error unmarshalling message:", err)
		return
	}

	Handler(connectionId, msg)
}
