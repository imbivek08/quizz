package ws

import (
	"encoding/json"
	"log"
)

// ClientMessage wraps a client and their message
type ClientMessage struct {
	Client  *Client
	Message *Message
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	Clients map[*Client]bool

	// Inbound messages from clients
	HandleMessage chan *ClientMessage

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		HandleMessage: make(chan *ClientMessage),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Clients:       make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client registered: %s", client.ID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %s", client.ID)
			}

		case clientMsg := <-h.HandleMessage:
			h.handleClientMessage(clientMsg)
		}
	}
}

func (h *Hub) handleClientMessage(cm *ClientMessage) {
	switch cm.Message.Type {
	case MessageTypeJoinRoom:
		h.handleJoinRoom(cm)
	case MessageTypeLeaveRoom:
		h.handleLeaveRoom(cm)
	default:
		log.Printf("Unknown message type: %s", cm.Message.Type)
	}
}

func (h *Hub) handleJoinRoom(cm *ClientMessage) {
	var payload JoinRoomPayload
	if err := json.Unmarshal(cm.Message.Payload, &payload); err != nil {
		log.Printf("error unmarshaling join room payload: %v", err)
		return
	}

	cm.Client.Username = payload.Username
	cm.Client.RoomID = payload.RoomID

	log.Printf("Client %s joined room %s as %s", cm.Client.ID, payload.RoomID, payload.Username)

	// Broadcast to room that user joined (we'll improve this later)
	h.BroadcastToRoom(payload.RoomID, MessageTypeRoomUsers, map[string]string{
		"message": payload.Username + " joined the room",
	})
}

func (h *Hub) handleLeaveRoom(cm *ClientMessage) {
	log.Printf("Client %s left room %s", cm.Client.ID, cm.Client.RoomID)
	cm.Client.RoomID = ""
}

// BroadcastToRoom sends a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID string, msgType MessageType, payload interface{}) {
	message, err := NewMessage(msgType, payload)
	if err != nil {
		log.Printf("error creating message: %v", err)
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	for client := range h.Clients {
		if client.RoomID == roomID {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}
