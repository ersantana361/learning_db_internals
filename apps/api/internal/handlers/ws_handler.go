package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// MessageHandler processes incoming messages from clients
type MessageHandler func(clientID string, message []byte)

// DisconnectHandler is called when a client disconnects
type DisconnectHandler func(clientID string)

// HandleWebSocket creates a WebSocket handler with message routing
func HandleWebSocket(hub *Hub, onMessage MessageHandler, onDisconnect DisconnectHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := hub.Register(conn)
		log.Printf("Client connected: %s", client.ID)

		go writePump(client)
		go readPumpWithHandler(client, onMessage, onDisconnect)
	}
}

// readPumpWithHandler reads messages and routes them through the handler
func readPumpWithHandler(client *Client, onMessage MessageHandler, onDisconnect DisconnectHandler) {
	defer func() {
		log.Printf("Client disconnected: %s", client.ID)
		if onDisconnect != nil {
			onDisconnect(client.ID)
		}
		client.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if onMessage != nil {
			onMessage(client.ID, message)
		}
	}
}

func writePump(client *Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
