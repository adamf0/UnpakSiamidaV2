package application

import "time"

type VerifyEmailHistoryCommand struct {
	ID uint
}

type VerifyEmailHistoryResult struct {
	ID         uint      `json:"id"`
	MessageID  string    `json:"message_id"`
	IsSent     bool      `json:"is_sent"`
	Labels     []string  `json:"labels"`
	Snippet    string    `json:"snippet"`
	VerifiedAt time.Time `json:"verified_at"`
}
