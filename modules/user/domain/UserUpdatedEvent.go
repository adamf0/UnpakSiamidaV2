package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type UserUpdatedEvent struct {
	EventID    uuid.UUID
	UserUUID    uuid.UUID
	OccurredOn time.Time
}

func (e UserUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e UserUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}