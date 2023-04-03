package ws

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WebSocketClient struct {
	Conn   net.Conn
	Server *WebSocketServer
}

// Web Socket server
type WebSocketServer struct {
	address                  string // Address to open connection: localhost:9999
	config                   *tls.Config
	onNewClientCallback      func(c *WebSocketClient)
	onClientConnectionClosed func(c *WebSocketClient, err error)
	onNewMessage             func(c *WebSocketClient, message string)
}

// Called right after server starts listening new client
func (s *WebSocketServer) OnNewClient(callback func(c *WebSocketClient)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *WebSocketServer) OnClientConnectionClosed(callback func(c *WebSocketClient, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *WebSocketServer) OnNewMessage(callback func(c *WebSocketClient, message string)) {
	s.onNewMessage = callback
}

func NewWebSocketServer(address string) *WebSocketServer {
	fmt.Println("Creating server with address", address)
	server := &WebSocketServer{
		address: address,
	}

	// Add empty call backs
	server.OnNewClient(func(c *WebSocketClient) {})
	server.OnNewMessage(func(c *WebSocketClient, message string) {})
	server.OnClientConnectionClosed(func(c *WebSocketClient, err error) {})

	return server
}

func (client *WebSocketClient) listen() {
	client.Server.onNewClientCallback(client)
	for {
		msg, op, err := wsutil.ReadClientData(client.Conn)
		if err != nil {
			fmt.Println("unknown message", "error", err.Error())
		}
		switch op {
		case ws.OpPing:
			if err := wsutil.WriteServerMessage(client.Conn, ws.OpPong, msg); err != nil {
				fmt.Println("failed to write pong message ", "error", err.Error())
				return
			}
		case ws.OpText:
			fmt.Println("receving message from client", "message", msg)
			client.Server.onNewMessage(client, string(msg))
		case ws.OpClose:
			fmt.Println("client is disconnected")
			client.Conn.Close()
		}

	}
}

func (s *WebSocketServer) Listen() {
	var ln net.Listener
	var err error
	if s.config == nil {
		ln, err = net.Listen("tcp", "localhost:8080")
	} else {
		ln, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.\r\n", err)
	}
	defer ln.Close()

	for {
		// Accept waits for and returns the next connection to the listener (Blocking)
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error")
		}
		_, err = ws.Upgrade(conn)
		if err != nil {
			fmt.Println("Connection upgrade error")
		}

		client := WebSocketClient{
			Conn:   conn,
			Server: s,
		}

		go client.listen()
	}

}
