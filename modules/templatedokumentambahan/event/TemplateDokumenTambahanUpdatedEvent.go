package event

import (
	"time"

	"github.com/google/uuid"
	// "UnpakSiamida/common/domain"
)

type TemplateDokumenTambahanUpdatedEvent struct {
	EventID                     uuid.UUID
	TemplateDokumenTambahanUUID uuid.UUID
	OccurredOn                  time.Time
}

func (e TemplateDokumenTambahanUpdatedEvent) ID() string {
	return e.EventID.String()
}

func (e TemplateDokumenTambahanUpdatedEvent) OccurredOnUTC() time.Time {
	return e.OccurredOn
}
