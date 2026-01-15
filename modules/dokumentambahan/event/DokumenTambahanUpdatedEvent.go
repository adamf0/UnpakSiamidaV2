package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type DokumenTambahanUpdatedEvent struct {
	EventID             uuid.UUID
	DokumenTambahanUUID uuid.UUID
	OccurredOn          time.Time
}

func (e DokumenTambahanUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e DokumenTambahanUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
