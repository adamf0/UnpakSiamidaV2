package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type GenerateDokumenTambahanUpdatedEvent struct {
	EventID                     uuid.UUID
	GenerateDokumenTambahanUUID uuid.UUID
	OccurredOn                  time.Time
}

func (e GenerateDokumenTambahanUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e GenerateDokumenTambahanUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
