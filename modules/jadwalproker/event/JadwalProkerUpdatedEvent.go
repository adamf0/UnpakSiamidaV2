package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type JadwalProkerUpdatedEvent struct {
	EventID          uuid.UUID
	JadwalProkerUUID uuid.UUID
	OccurredOn       time.Time
}

func (e JadwalProkerUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e JadwalProkerUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
