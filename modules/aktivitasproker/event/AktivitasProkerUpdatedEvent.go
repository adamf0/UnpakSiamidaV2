package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type AktivitasProkerUpdatedEvent struct {
	EventID             uuid.UUID
	AktivitasProkerUUID uuid.UUID
	OccurredOn          time.Time
}

func (e AktivitasProkerUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e AktivitasProkerUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
