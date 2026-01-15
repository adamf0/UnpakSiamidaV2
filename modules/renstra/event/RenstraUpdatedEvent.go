package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type RenstraUpdatedEvent struct {
	EventID     uuid.UUID
	RenstraUUID uuid.UUID
	OccurredOn  time.Time
}

func (e RenstraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e RenstraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
