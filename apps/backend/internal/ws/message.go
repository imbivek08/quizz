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

	// Errors
	MessageTypeError MessageType = "error"
	// Game events
	MessageTypeGameStart    MessageType = "game_start"
	MessageTypeGameEnd      MessageType = "game_end"
	MessageTypeQuestion     MessageType = "question"
	MessageTypeAnswer       MessageType = "submit_answer"
	MessageTypeAnswerResult MessageType = "answer_result"
	MessageTypeLeaderboard  MessageType = "leaderboard"
	MessageTypeNextQuestion MessageType = "next_question"
)

// Message is the wrapper for all WebSocket messages
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type AnswerPayload struct {
	QuestionID int   `json:"question_id"`
	Answer     int   `json:"answer"`
	Timestamp  int64 `json:"timestamp"`
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
