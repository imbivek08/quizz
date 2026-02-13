package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/imbivek08/quizz/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development (restrict in production!)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	// Create a new client
	client := &ws.Client{
		ID:   uuid.New().String(),
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// Register the client with the hub
	client.Hub.Register <- client

	// Start the client's read and write pumps
	go client.WritePump()
	go client.ReadPump()

	log.Printf("New WebSocket connection established: %s", client.ID)
}
