package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type MataProgramUpdatedEvent struct {
	EventID         uuid.UUID
	MataProgramUUID uuid.UUID
	OccurredOn      time.Time
}

func (e MataProgramUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e MataProgramUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
