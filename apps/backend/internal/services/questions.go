package services

type Question struct {
	ID            int      `json:"id"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correct_answer"` // index of correct option
	TimeLimit     int      `json:"time_limit"`     // seconds
}

// GetSafeQuestion returns question without the correct answer
func (q *Question) GetSafeQuestion() map[string]interface{} {
	return map[string]interface{}{
		"id":         q.ID,
		"text":       q.Text,
		"options":    q.Options,
		"time_limit": q.TimeLimit,
	}
}

// Sample questions for testing
func GetSampleQuestions() []Question {
	return []Question{
		{
			ID:            1,
			Text:          "What is the capital of France?",
			Options:       []string{"London", "Paris", "Berlin", "Madrid"},
			CorrectAnswer: 1,
			TimeLimit:     10,
		},
		{
			ID:            2,
			Text:          "Which planet is known as the Red Planet?",
			Options:       []string{"Venus", "Mars", "Jupiter", "Saturn"},
			CorrectAnswer: 1,
			TimeLimit:     10,
		},
		{
			ID:            3,
			Text:          "What is 2 + 2?",
			Options:       []string{"3", "4", "5", "6"},
			CorrectAnswer: 1,
			TimeLimit:     5,
		},
		{
			ID:            4,
			Text:          "Who wrote 'Romeo and Juliet'?",
			Options:       []string{"Charles Dickens", "William Shakespeare", "Mark Twain", "Jane Austen"},
			CorrectAnswer: 1,
			TimeLimit:     10,
		},
		{
			ID:            5,
			Text:          "What is the largest ocean on Earth?",
			Options:       []string{"Atlantic", "Indian", "Arctic", "Pacific"},
			CorrectAnswer: 3,
			TimeLimit:     10,
		},
	}
}
