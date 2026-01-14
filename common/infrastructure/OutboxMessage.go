package infrastructure

import "time"

type OutboxMessage struct {
	ID            string    `gorm:"primaryKey;size:36"`
	Type          string    `gorm:"size:255;index"`
	Payload       string    `gorm:"type:longtext"`
	OccurredOnUTC time.Time
	ProcessedAt   *time.Time
	Error         *string   `gorm:"type:text"`
}
