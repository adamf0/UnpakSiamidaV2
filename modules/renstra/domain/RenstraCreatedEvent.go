package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type RenstraCreatedEvent struct {
	EventID    uuid.UUID
	RenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e RenstraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e RenstraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}