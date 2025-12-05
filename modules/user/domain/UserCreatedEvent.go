package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type UserCreatedEvent struct {
	EventID    uuid.UUID
	UserUUID    uuid.UUID
	OccurredOn time.Time
}

func (e UserCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e UserCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}