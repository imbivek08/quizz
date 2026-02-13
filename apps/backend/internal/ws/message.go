package ws

import "encoding/json"

type MessageType string

const (
	// Connection events
	MessageTypeConnected    MessageType = "connected"
	MessageTypeDisconnected MessageType = "disconnected"

	// Room events
	MessageTypeJoinRoom  MessageType = "join_room"
	MessageTypeLeaveRoom MessageType = "leave_room"
	MessageTypeRoomUsers MessageType = "room_users"

	// Quiz events (we'll expand later)
	MessageTypeAnswer      MessageType = "submit_answer"
	MessageTypeQuestion    MessageType = "question"
	MessageTypeLeaderboard MessageType = "leaderboard"

	// Errors
	MessageTypeError MessageType = "error"
)

// Message is the wrapper for all WebSocket messages
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// JoinRoomPayload - when a user joins a room
type JoinRoomPayload struct {
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

// ErrorPayload - error messages
type ErrorPayload struct {
	Message string `json:"message"`
}

// Helper function to create a message
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    msgType,
		Payload: payloadBytes,
	}, nil
}
