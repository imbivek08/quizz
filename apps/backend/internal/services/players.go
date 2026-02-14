package services

import "time"

type PlayerStatus string

const (
	PlayerStatusWaiting  PlayerStatus = "waiting"
	PlayerStatusPlaying  PlayerStatus = "playing"
	PlayerStatusFinished PlayerStatus = "finished"
)

type Player struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	Status       PlayerStatus `json:"status"`
	Score        int          `json:"score"`
	TotalTime    int64        `json:"total_time"`
	CorrectCount int          `json:"correct_count"`
	JoinedAt     time.Time    `json:"joined_at"`
}

func NewPlayer(id, username string) *Player {
	return &Player{
		ID:           id,
		Username:     username,
		Status:       PlayerStatusWaiting,
		Score:        0,
		TotalTime:    0,
		CorrectCount: 0,
		JoinedAt:     time.Now(),
	}
}

func (p *Player) CalculateScore(correct bool, timeMs int64, maxTimeMs int64) int {
	if !correct {
		return 0
	}

	basePoints := 100

	// Speed bonus: 50 points if instant, 0 if at time limit
	speedBonus := 0
	if timeMs < maxTimeMs {
		speedBonus = int(50 * (maxTimeMs - timeMs) / maxTimeMs)
	}

	return basePoints + speedBonus
}
