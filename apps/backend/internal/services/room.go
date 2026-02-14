package services

import (
	"log"
	"sync"
	"time"
)

type RoomStatus string

const (
	RoomStatusWaiting  RoomStatus = "waiting"
	RoomStatusPlaying  RoomStatus = "playing"
	RoomStatusFinished RoomStatus = "finished"
)

type Room struct {
	ID                string             `json:"id"`
	Players           map[string]*Player `json:"players"`
	Status            RoomStatus         `json:"status"`
	Questions         []Question         `json:"-"`
	CurrentQuestion   int                `json:"current_question"`
	QuestionStartTime time.Time          `json:"-"`
	Answers           map[string]int     `json:"-"` // playerID -> answer index
	MaxPlayers        int                `json:"max_players"`
	MinPlayers        int                `json:"min_players"`
	CreatedAt         time.Time          `json:"created_at"`
	mu                sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:              id,
		Players:         make(map[string]*Player),
		Status:          RoomStatusWaiting,
		Questions:       GetSampleQuestions(),
		CurrentQuestion: -1,
		Answers:         make(map[string]int),
		MaxPlayers:      10,
		MinPlayers:      1, // Set to 1 for testing, change to 2+ for production
		CreatedAt:       time.Now(),
	}
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(player *Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= r.MaxPlayers {
		return false
	}

	if r.Status != RoomStatusWaiting {
		return false
	}

	r.Players[player.ID] = player
	log.Printf("Player %s joined room %s", player.Username, r.ID)
	return true
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if player, exists := r.Players[playerID]; exists {
		delete(r.Players, playerID)
		log.Printf("Player %s left room %s", player.Username, r.ID)
	}
}

// CanStart checks if the room can start the game
func (r *Room) CanStart() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Status == RoomStatusWaiting && len(r.Players) >= r.MinPlayers
}

// StartGame starts the quiz game
func (r *Room) StartGame() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Status = RoomStatusPlaying
	r.CurrentQuestion = 0

	for _, player := range r.Players {
		player.Status = PlayerStatusPlaying
	}

	log.Printf("Room %s started game with %d players", r.ID, len(r.Players))
}

// GetCurrentQuestion returns the current question (safe for clients)
func (r *Room) GetCurrentQuestion() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.CurrentQuestion < 0 || r.CurrentQuestion >= len(r.Questions) {
		return nil
	}

	return r.Questions[r.CurrentQuestion].GetSafeQuestion()
}

// SubmitAnswer handles a player's answer
func (r *Room) SubmitAnswer(playerID string, answerIndex int, submittedAt time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if answer already submitted
	if _, exists := r.Answers[playerID]; exists {
		return false
	}

	// Check if player exists
	player, exists := r.Players[playerID]
	if !exists {
		return false
	}

	// Store answer
	r.Answers[playerID] = answerIndex

	// Calculate time taken
	timeTaken := submittedAt.Sub(r.QuestionStartTime).Milliseconds()

	// Check if correct
	currentQ := r.Questions[r.CurrentQuestion]
	isCorrect := answerIndex == currentQ.CorrectAnswer

	// Calculate score
	maxTimeMs := int64(currentQ.TimeLimit * 1000)
	points := player.CalculateScore(isCorrect, timeTaken, maxTimeMs)

	// Update player stats
	player.Score += points
	player.TotalTime += timeTaken
	if isCorrect {
		player.CorrectCount++
	}

	log.Printf("Player %s answered question %d: %v (score: +%d)",
		player.Username, r.CurrentQuestion, isCorrect, points)

	return true
}

// AllPlayersAnswered checks if all players have answered
func (r *Room) AllPlayersAnswered() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.Answers) >= len(r.Players)
}

// NextQuestion moves to the next question
func (r *Room) NextQuestion() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clear previous answers
	r.Answers = make(map[string]int)

	// Move to next question
	r.CurrentQuestion++

	// Check if quiz is finished
	if r.CurrentQuestion >= len(r.Questions) {
		r.Status = RoomStatusFinished
		for _, player := range r.Players {
			player.Status = PlayerStatusFinished
		}
		log.Printf("Room %s finished quiz", r.ID)
		return false
	}

	r.QuestionStartTime = time.Now()
	log.Printf("Room %s moved to question %d", r.ID, r.CurrentQuestion)
	return true
}

// GetLeaderboard returns sorted list of players by score
func (r *Room) GetLeaderboard() []*Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]*Player, 0, len(r.Players))
	for _, player := range r.Players {
		players = append(players, player)
	}

	// Sort by score (descending), then by time (ascending)
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			if players[i].Score < players[j].Score ||
				(players[i].Score == players[j].Score && players[i].TotalTime > players[j].TotalTime) {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	return players
}

// GetRoomInfo returns safe room info for clients
func (r *Room) GetRoomInfo() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	playerList := make([]map[string]interface{}, 0, len(r.Players))
	for _, player := range r.Players {
		playerList = append(playerList, map[string]interface{}{
			"id":       player.ID,
			"username": player.Username,
			"status":   player.Status,
		})
	}

	return map[string]interface{}{
		"id":           r.ID,
		"status":       r.Status,
		"player_count": len(r.Players),
		"max_players":  r.MaxPlayers,
		"players":      playerList,
	}
}
