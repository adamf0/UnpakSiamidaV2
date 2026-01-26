package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type MataProgramCreatedEvent struct {
	EventID         uuid.UUID
	MataProgramUUID uuid.UUID
	OccurredOn      time.Time
}

func (e MataProgramCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e MataProgramCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
