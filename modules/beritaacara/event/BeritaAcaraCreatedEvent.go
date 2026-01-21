package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type BeritaAcaraCreatedEvent struct {
	EventID         uuid.UUID
	BeritaAcaraUUID uuid.UUID
	OccurredOn      time.Time
}

func (e BeritaAcaraCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e BeritaAcaraCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
