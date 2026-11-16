package infrastructure

import "time"

type OutboxMessage struct {
	ID            string     `gorm:"primaryKey;size:36"`
	Type          string     `gorm:"size:255;index"`
	SchemaVersion int        `gorm:"default:1"`
	AggregateType *string    `gorm:"size:64;index"`
	AggregateID   *string    `gorm:"size:64;index"`
	Payload       string     `gorm:"type:longtext"`
	Status        string     `gorm:"type:enum('PENDING','PROCESSED','FAILED');default:PENDING"`
	RetryCount    int        `gorm:"default:0"`
	OccurredOnUTC time.Time
	ProcessedAt   *time.Time
	Error         *string    `gorm:"type:text"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
}
