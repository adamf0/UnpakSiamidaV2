package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type DokumenProkerCreatedEvent struct {
	EventID           uuid.UUID
	DokumenProkerUUID uuid.UUID
	OccurredOn        time.Time
}

func (e DokumenProkerCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e DokumenProkerCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
