package ws

import (
	"encoding/json"
)

func (service *WebSocketService) StartServer(server WebSocketServer) {
	// Initiall add empty ws client list
	service.ClientList = make(map[string]*WebSocketClient, 0)

	// Add call back functions
	server.OnNewClient(func(c *WebSocketClient) {
		service.Logger.Infow("client connected", "connection", "")
	})

	server.OnNewMessage(func(client *WebSocketClient, message string) {
		// new message received
		decodeMessage := Message{}
		json.Unmarshal([]byte(message), &decodeMessage)

		service.Logger.Infow("received the message", "remote_ip", decodeMessage)

	})

	server.OnClientConnectionClosed(func(client *WebSocketClient, err error) {
		// connection with client lost
		service.Logger.Infow("client disconnected", "remote_ip")
	})

}
