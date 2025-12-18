package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type GenerateRenstraUpdatedEvent struct {
	EventID    uuid.UUID
	GenerateRenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e GenerateRenstraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e GenerateRenstraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}