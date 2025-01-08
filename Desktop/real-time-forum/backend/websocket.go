package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket connection upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Client represents a single WebSocket connection
type Client struct {
	ID    string
	Conn  *websocket.Conn
	Mutex sync.Mutex
	Send  chan []byte
}

// Global map to keep track of connected clients
var clients = make(map[string]*Client)
var clientsMutex sync.Mutex

// WebSocketHandler upgrades HTTP requests to WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Generate a unique ID for the client
	clientID := r.RemoteAddr
	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte),
	}

	clientsMutex.Lock()
	clients[clientID] = client
	clientsMutex.Unlock()

	log.Println("Client connected:", clientID)

	// Start listening for messages from this client
	go handleClientMessages(client)

	// Start a writer goroutine for this client
	go handleClientWrites(client)
}

// handleClientMessages listens for messages from a client
func handleClientMessages(client *Client) {
	defer func() {
		clientsMutex.Lock()
		delete(clients, client.ID)
		clientsMutex.Unlock()
		client.Conn.Close()
		log.Println("Client disconnected:", client.ID)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Message received from %s: %s\n", client.ID, message)

		// Broadcast the message to all clients
		broadcastMessage(client.ID, message)
	}
}

// handleClientWrites listens for outgoing messages to the client
func handleClientWrites(client *Client) {
	for message := range client.Send {
		client.Mutex.Lock()
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		client.Mutex.Unlock()
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

// broadcastMessage sends a message to all connected clients
func broadcastMessage(senderID string, message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for clientID, client := range clients {
		if clientID != senderID {
			select {
			case client.Send <- message:
			default:
				log.Printf("Client %s is not receiving messages\n", clientID)
			}
		}
	}
}
