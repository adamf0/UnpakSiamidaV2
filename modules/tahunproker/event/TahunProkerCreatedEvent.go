package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TahunProkerCreatedEvent struct {
	EventID         uuid.UUID
	TahunProkerUUID uuid.UUID
	OccurredOn      time.Time
}

func (e TahunProkerCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e TahunProkerCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
