package domain

import (
	"time"
	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type IndikatorRenstraUpdatedEvent struct {
	EventID    uuid.UUID
	IndikatorRenstraUUID    uuid.UUID
	OccurredOn time.Time
}

func (e IndikatorRenstraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e IndikatorRenstraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}