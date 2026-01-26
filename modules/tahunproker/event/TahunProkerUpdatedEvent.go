package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TahunProkerUpdatedEvent struct {
	EventID         uuid.UUID
	TahunProkerUUID uuid.UUID
	OccurredOn      time.Time
}

func (e TahunProkerUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e TahunProkerUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
