package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type BeritaAcaraUpdatedEvent struct {
	EventID         uuid.UUID
	BeritaAcaraUUID uuid.UUID
	OccurredOn      time.Time
}

func (e BeritaAcaraUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e BeritaAcaraUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
