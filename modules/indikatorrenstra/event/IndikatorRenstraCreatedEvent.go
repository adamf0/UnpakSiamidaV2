package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type IndikatorRenstraCreatedEvent struct {
	EventID              uuid.UUID
	IndikatorRenstraUUID uuid.UUID
	OccurredOn           time.Time
}

func (e IndikatorRenstraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e IndikatorRenstraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
