package domain

import "time"

type HistoryEmail struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OutboxID       string     `gorm:"size:64;index" json:"outbox_id"`
	UserUUID       string     `gorm:"size:64;index" json:"user_uuid"`
	Username       string     `gorm:"size:100" json:"username"`
	RecipientEmail string     `gorm:"size:150;not null" json:"recipient_email"`
	Subject        string     `gorm:"size:255" json:"subject"`
	Status         string     `gorm:"size:50;not null;default:'PENDING';index" json:"status"` // PENDING, SENDING, SENT, FAILED
	Provider       string     `gorm:"size:50;default:'GOOGLE_SMTP'" json:"provider"`          // GOOGLE_SMTP, SMTP
	MessageID      string     `gorm:"size:255" json:"message_id"`
	ThreadID       string     `gorm:"size:255" json:"thread_id"`
	ErrorMessage   *string    `gorm:"type:text" json:"error_message"`
	SentAt         *time.Time `json:"sent_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (HistoryEmail) TableName() string {
	return "history_email"
}
