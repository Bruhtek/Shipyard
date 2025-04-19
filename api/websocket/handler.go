package websocket

import (
	"Shipyard/env_manager"
	"context"
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

	if !ok1 || !ok2 || !ok3 {
		log.Println("[WS] Invalid message format")
		return
	}

	_, ok := env_manager.EnvManager.Env[envName]
	if !ok {
		log.Println("[WS] Environment not found:", envName)
		return
	}

	println("[WS] Received message:", objectType, action, envName)

	ctx, _ := context.WithCancel(context.Background())

	runner := Runner{
		Command: []string{"docker", "pull", "linuxserver/apprise-api"},
		TaskId:  "test-task",
		Ctx:     ctx,
	}

	go runner.Run()
}
