package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type AktivitasProkerCreatedEvent struct {
	EventID             uuid.UUID
	AktivitasProkerUUID uuid.UUID
	OccurredOn          time.Time
}

func (e AktivitasProkerCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e AktivitasProkerCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
