package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type DokumenProkerUpdatedEvent struct {
	EventID           uuid.UUID
	DokumenProkerUUID uuid.UUID
	OccurredOn        time.Time
}

func (e DokumenProkerUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e DokumenProkerUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
