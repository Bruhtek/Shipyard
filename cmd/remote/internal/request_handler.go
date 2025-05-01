package internal

import "encoding/json"

func RequestHandler(message []byte) {
	var msg map[string]interface{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		println("ERROR: Failed to unmarshal message:", err.Error())
		return
	}

	actionType, ok := msg["Type"].(string)
	if !ok {
		println("ERROR: ActionType not found in message")
		return
	}

	actionData, ok := msg["Data"].(map[string]interface{})
	if !ok {
		println("ERROR: ActionData not found in message")
		return
	}

	switch actionType {
	case "Websocket":
		println("INFO: Websocket action received")
		println("INFO: Action Data:", actionData)

		WSHandler(actionData)

		break
	case "REST":
		println("INFO: REST action received")
		break
	}
}
