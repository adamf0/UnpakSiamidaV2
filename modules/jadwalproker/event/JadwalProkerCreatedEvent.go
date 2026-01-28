package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type JadwalProkerCreatedEvent struct {
	EventID          uuid.UUID
	JadwalProkerUUID uuid.UUID
	OccurredOn       time.Time
}

func (e JadwalProkerCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e JadwalProkerCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
