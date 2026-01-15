package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type JenisFileCreatedEvent struct {
	EventID       uuid.UUID
	JenisFileUUID uuid.UUID
	OccurredOn    time.Time
}

func (e JenisFileCreatedEvent) ID() string {
	return e.EventID.String()
}

func (e JenisFileCreatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
