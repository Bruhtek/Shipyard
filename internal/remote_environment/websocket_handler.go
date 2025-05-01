package remote_environment

func WebsocketMessageHandler(message []byte) {
	println("Received message:", string(message))
}
